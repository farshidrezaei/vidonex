// Package compiler translates declarative Timeline AST specifications into optimized FFmpeg filtergraph DAGs and CLI commands.
package compiler

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/farshidrezaei/vidonyx/ducking"
	"github.com/farshidrezaei/vidonyx/filtergraph"
	"github.com/farshidrezaei/vidonyx/subtitles"
	"github.com/farshidrezaei/vidonyx/timeline"
	"github.com/farshidrezaei/vidonyx/visualizer"
	"github.com/farshidrezaei/vidonyx/waveform"
)

// isImageSource reports whether a media source file path is an image.
func isImageSource(sourcePath string) bool {
	ext := strings.ToLower(filepath.Ext(sourcePath))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".bmp", ".gif", ".tiff", ".svg":
		return true
	default:
		return false
	}
}

// CompilationResult contains the compiled artifacts: DAG, CLI arguments, and diagram representations.
type CompilationResult struct {
	Graph         *filtergraph.Graph
	Inputs        []string
	Args          []string
	FilterComplex string
	OutVLabel     string
	OutALabel     string
	OutputPath    string
}

// Mermaid returns the Mermaid.js flowchart representation of the compiled filtergraph DAG.
func (result *CompilationResult) Mermaid() (string, error) {
	return visualizer.ToMermaid(result.Graph)
}

// DOT returns the Graphviz DOT representation of the compiled filtergraph DAG.
func (result *CompilationResult) DOT() (string, error) {
	return visualizer.ToDOT(result.Graph)
}

// Compiler translates a declarative Timeline AST into an optimized Filtergraph DAG and FFmpeg CLI command.
type Compiler struct {
	pipeline *filtergraph.Pipeline
	logger   *slog.Logger
	encoding EncodingOptions
}

// New creates a new Compiler instance.
func New(logger *slog.Logger) *Compiler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Compiler{
		pipeline: filtergraph.DefaultPipeline(logger),
		logger:   logger,
		encoding: DefaultEncodingOptions(),
	}
}

// SetEncodingOptions overrides default encoding parameters.
func (compilerInstance *Compiler) SetEncodingOptions(options EncodingOptions) *Compiler {
	compilerInstance.encoding = options
	return compilerInstance
}

// Compile translates the given timeline and output path into a complete CompilationResult.
func (compilerInstance *Compiler) Compile(compositionTimeline *timeline.Timeline, outputPath string) (*CompilationResult, error) {
	if err := timeline.Validate(compositionTimeline); err != nil {
		return nil, fmt.Errorf("compiler: timeline validation failed: %w", err)
	}

	graph := filtergraph.NewGraph()

	// 1. Discover unique input sources and map to FFmpeg input indices (0, 1, 2, ...)
	inputIndexMap := make(map[string]int)
	uniqueInputs := make([]string, 0)

	for _, track := range compositionTimeline.Tracks {
		for _, clip := range track.Clips {
			if clip.Source != "" {
				if _, exists := inputIndexMap[clip.Source]; !exists {
					inputIndexMap[clip.Source] = len(uniqueInputs)
					uniqueInputs = append(uniqueInputs, clip.Source)
				}
			}
		}
	}

	processedVideoPads := make([]*ClipVideoPad, 0)
	trackAudioPadMap := make(map[string]*filtergraph.Pad)

	// 2. Process each track's clips
	for trackIndex, track := range compositionTimeline.Tracks {
		trackVideoPads := make([]*filtergraph.Pad, 0, len(track.Clips))
		trackAudioPads := make([]*filtergraph.Pad, 0, len(track.Clips))

		for _, clip := range track.Clips {
			inputIndex := inputIndexMap[clip.Source]

			// Handle Video
			if track.Kind == timeline.TrackKindVideo || track.Kind == timeline.TrackKindOverlay {
				rawVideoIn := &filtergraph.Pad{
					ID:         fmt.Sprintf("%d:v", inputIndex),
					StreamType: filtergraph.StreamTypeVideo,
					IsInput:    false,
				}

				clipForVideo := clip
				if len(track.Transitions) > 0 {
					clipForVideoCopy := *clip
					clipForVideoCopy.TimelineStart = 0
					clipForVideo = &clipForVideoCopy
				}

				// 1. Normalize clip stream to canvas bounds, SAR, FPS and yuva420p format
				normalizedVideoPad, err := InjectVideoNormalizer(graph, rawVideoIn, compositionTimeline.Canvas, compositionTimeline.FPS)
				if err != nil {
					return nil, fmt.Errorf("compiler: failed normalizing video clip %q: %w", clip.ID, err)
				}

				// 2. Clip video pipeline (trim, speed, chromakey, fade, opacity, scale, rotation)
				clipVideoPad, err := ProcessClipVideo(graph, normalizedVideoPad, clipForVideo, compositionTimeline)
				if err != nil {
					return nil, fmt.Errorf("compiler: failed processing video clip %q: %w", clip.ID, err)
				}

				trackVideoPads = append(trackVideoPads, clipVideoPad)
			}

			// Handle Audio (Skip for static image sources or clips without audio streams)
			if (track.Kind == timeline.TrackKindAudio || track.Kind == timeline.TrackKindVideo) && !track.Muted && !isImageSource(clip.Source) && clip.HasAudioStream {
				rawAudioIn := &filtergraph.Pad{
					ID:         fmt.Sprintf("%d:a", inputIndex),
					StreamType: filtergraph.StreamTypeAudio,
					IsInput:    false,
				}

				clipForAudio := clip
				if len(track.Transitions) > 0 {
					clipForAudioCopy := *clip
					clipForAudioCopy.TimelineStart = 0
					clipForAudio = &clipForAudioCopy
				}

				clipAudioPad, err := ProcessClipAudio(graph, rawAudioIn, clipForAudio)
				if err != nil {
					return nil, fmt.Errorf("compiler: failed processing audio clip %q: %w", clip.ID, err)
				}

				normalizedAudioPad, err := InjectAudioNormalizer(graph, clipAudioPad)
				if err != nil {
					return nil, fmt.Errorf("compiler: failed normalizing audio clip %q: %w", clip.ID, err)
				}

				trackAudioPads = append(trackAudioPads, normalizedAudioPad)
			}
		}

		// If track has transitions, chain the track's clips together with xfade / acrossfade
		if len(track.Transitions) > 0 && len(trackVideoPads) > 1 {
			chainedVideoPad, chainedDuration, err := ChainTrackTransitionsVideo(graph, trackVideoPads, track.Transitions, track.Clips, compositionTimeline.Canvas)
			if err != nil {
				return nil, err
			}
			if len(track.Clips) > 0 {
				if track.Clips[0].TimelineStart > 0 {
					ptsChainNode := graph.NewNode(fmt.Sprintf("pts_chain_%s", track.ID), "setpts")
					ptsChainNode.SetParam("expr", fmt.Sprintf("PTS-STARTPTS+%.4f/TB", track.Clips[0].TimelineStart.Seconds()))
					ptsIn := ptsChainNode.AddInput(chainedVideoPad.ID, filtergraph.StreamTypeVideo)
					ptsOut := ptsChainNode.AddOutput(graph.NextPadID("pts_chain_out"), filtergraph.StreamTypeVideo)
					_ = graph.Connect(chainedVideoPad, ptsIn)
					chainedVideoPad = ptsOut
				}
				syntheticChainedClip := &timeline.Clip{
					ID:            fmt.Sprintf("chained_%s", track.ID),
					TimelineStart: track.Clips[0].TimelineStart,
					Duration:      chainedDuration,
				}
				processedVideoPads = append(processedVideoPads, &ClipVideoPad{
					Clip:       syntheticChainedClip,
					ZIndex:     track.ZIndex,
					TrackIndex: trackIndex,
					Pad:        chainedVideoPad,
				})
			}
		} else {
			for index, clip := range track.Clips {
				if index < len(trackVideoPads) {
					processedVideoPads = append(processedVideoPads, &ClipVideoPad{
						Clip:       clip,
						ZIndex:     track.ZIndex,
						TrackIndex: trackIndex,
						Pad:        trackVideoPads[index],
					})
				}
			}
		}

		if len(track.Transitions) > 0 && len(trackAudioPads) > 1 {
			chainedAudioPad, err := ChainTrackTransitionsAudio(graph, trackAudioPads, track.Transitions, track.Clips)
			if err != nil {
				return nil, err
			}
			trackAudioPadMap[track.ID] = chainedAudioPad
		} else if len(trackAudioPads) == 1 {
			trackAudioPadMap[track.ID] = trackAudioPads[0]
		} else if len(trackAudioPads) > 1 {
			// Mix intra-track audio clips
			intraMixNode := graph.NewNode(fmt.Sprintf("intra_amix_%s", track.ID), "amix")
			intraMixNode.SetParam("inputs", len(trackAudioPads))
			intraMixNode.SetParam("duration", "longest")
			for _, ap := range trackAudioPads {
				inPad := intraMixNode.AddInput(ap.ID, filtergraph.StreamTypeAudio)
				_ = graph.Connect(ap, inPad)
			}
			trackAudioPadMap[track.ID] = intraMixNode.AddOutput(graph.NextPadID("intra_mix_out"), filtergraph.StreamTypeAudio)
		}
	}

	// Apply sidechain ducking between tracks where requested
	for _, track := range compositionTimeline.Tracks {
		if track.DucksUnderTrackID != "" {
			musicPad, musicExists := trackAudioPadMap[track.ID]
			voicePad, voiceExists := trackAudioPadMap[track.DucksUnderTrackID]
			if musicExists && voiceExists && track.DuckingOptions != nil {
				duckedMusicPad, err := ducking.ApplySidechainDucking(graph, fmt.Sprintf("duck_%s_under_%s", track.ID, track.DucksUnderTrackID), musicPad, voicePad, *track.DuckingOptions)
				if err != nil {
					return nil, fmt.Errorf("compiler: failed applying ducking: %w", err)
				}
				trackAudioPadMap[track.ID] = duckedMusicPad
			}
		}
	}

	// Process animated audio waveforms
	for _, track := range compositionTimeline.Tracks {
		if track.Kind == timeline.TrackKindWaveform && track.WaveformSourceTrackID != "" {
			if sourceAudioPad, exists := trackAudioPadMap[track.WaveformSourceTrackID]; exists && track.WaveformOptions != nil {
				waveVideoPad, err := waveform.ApplyWaveformVisualizer(graph, fmt.Sprintf("waveform_%s", track.ID), sourceAudioPad, *track.WaveformOptions)
				if err != nil {
					return nil, fmt.Errorf("compiler: failed generating waveform: %w", err)
				}
				waveClip := &timeline.Clip{
					ID:            fmt.Sprintf("wave_clip_%s", track.ID),
					TimelineStart: 0,
					Duration:      compositionTimeline.Duration(),
					Position:      track.WaveformOptions.Position,
				}
				processedVideoPads = append(processedVideoPads, &ClipVideoPad{
					Clip:   waveClip,
					ZIndex: track.ZIndex,
					Pad:    waveVideoPad,
				})
			}
		}
	}

	processedAudioPads := make([]*filtergraph.Pad, 0, len(trackAudioPadMap))
	for _, track := range compositionTimeline.Tracks {
		if pad, ok := trackAudioPadMap[track.ID]; ok {
			processedAudioPads = append(processedAudioPads, pad)
		}
	}

	// 3. Compose final video stream
	finalVideoPad, err := BuildVideoCompositor(graph, compositionTimeline, processedVideoPads)
	if err != nil {
		return nil, fmt.Errorf("compiler: failed building video composition: %w", err)
	}

	// 4. Attach subtitle tracks if present
	for _, track := range compositionTimeline.Tracks {
		if track.Kind == timeline.TrackKindSubtitle && track.SubtitleTrack != nil {
			subtitledPad, err := subtitles.AttachSubtitles(graph, finalVideoPad, track.SubtitleTrack, compositionTimeline.Canvas)
			if err != nil {
				return nil, fmt.Errorf("compiler: failed attaching subtitles: %w", err)
			}
			finalVideoPad = subtitledPad
		}
	}

	// Rename final video pad to out_v
	outVLabel := "out_v"
	finalVideoPad.ID = outVLabel

	// 5. Mix final audio stream
	finalAudioPad, err := BuildAudioMixer(graph, compositionTimeline, processedAudioPads)
	if err != nil {
		return nil, fmt.Errorf("compiler: failed building audio mixer: %w", err)
	}

	// Rename final audio pad to out_a
	outALabel := "out_a"
	finalAudioPad.ID = outALabel

	// 6. Run Graph Optimization Passes
	if err := compilerInstance.pipeline.Execute(graph); err != nil {
		return nil, fmt.Errorf("compiler: graph optimization failed: %w", err)
	}

	// 7. Format filter_complex string
	filterComplexStr, err := graph.FormattedFilterComplex()
	if err != nil {
		return nil, fmt.Errorf("compiler: failed serializing filter_complex: %w", err)
	}

	// 8. Emit CLI Arguments
	args := BuildFFmpegArgs(uniqueInputs, filterComplexStr, outVLabel, outALabel, outputPath, compilerInstance.encoding)

	return &CompilationResult{
		Graph:         graph,
		Inputs:        uniqueInputs,
		Args:          args,
		FilterComplex: filterComplexStr,
		OutVLabel:     outVLabel,
		OutALabel:     outALabel,
		OutputPath:    outputPath,
	}, nil
}
