// Package compiler translates declarative Timeline AST specifications into optimized FFmpeg filtergraph DAGs and CLI commands.
package compiler

import (
	"fmt"
	"log/slog"

	"github.com/farshidrezaei/vidonyx/filtergraph"
	"github.com/farshidrezaei/vidonyx/subtitles"
	"github.com/farshidrezaei/vidonyx/timeline"
	"github.com/farshidrezaei/vidonyx/visualizer"
)

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
	processedAudioPads := make([]*filtergraph.Pad, 0)

	// 2. Process each track's clips
	for _, track := range compositionTimeline.Tracks {
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

				// Clip video pipeline (trim, speed, opacity, scale)
				clipVideoPad, err := ProcessClipVideo(graph, rawVideoIn, clip, compositionTimeline)
				if err != nil {
					return nil, fmt.Errorf("compiler: failed processing video clip %q: %w", clip.ID, err)
				}

				// Normalize clip stream to canvas bounds if needed
				normalizedVideoPad, err := InjectVideoNormalizer(graph, clipVideoPad, compositionTimeline.Canvas, compositionTimeline.FPS)
				if err != nil {
					return nil, fmt.Errorf("compiler: failed normalizing video clip %q: %w", clip.ID, err)
				}

				trackVideoPads = append(trackVideoPads, normalizedVideoPad)
			}

			// Handle Audio
			if (track.Kind == timeline.TrackKindAudio || track.Kind == timeline.TrackKindVideo) && !track.Muted {
				rawAudioIn := &filtergraph.Pad{
					ID:         fmt.Sprintf("%d:a", inputIndex),
					StreamType: filtergraph.StreamTypeAudio,
					IsInput:    false,
				}

				clipAudioPad, err := ProcessClipAudio(graph, rawAudioIn, clip)
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
			chainedVideoPad, err := ChainTrackTransitionsVideo(graph, trackVideoPads, track.Transitions, track.Clips)
			if err != nil {
				return nil, err
			}
			if len(track.Clips) > 0 {
				processedVideoPads = append(processedVideoPads, &ClipVideoPad{
					Clip:   track.Clips[0],
					ZIndex: track.ZIndex,
					Pad:    chainedVideoPad,
				})
			}
		} else {
			for index, clip := range track.Clips {
				if index < len(trackVideoPads) {
					processedVideoPads = append(processedVideoPads, &ClipVideoPad{
						Clip:   clip,
						ZIndex: track.ZIndex,
						Pad:    trackVideoPads[index],
					})
				}
			}
		}

		if len(track.Transitions) > 0 && len(trackAudioPads) > 1 {
			chainedAudioPad, err := ChainTrackTransitionsAudio(graph, trackAudioPads, track.Transitions)
			if err != nil {
				return nil, err
			}
			processedAudioPads = append(processedAudioPads, chainedAudioPad)
		} else {
			processedAudioPads = append(processedAudioPads, trackAudioPads...)
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
