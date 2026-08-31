package compiler

import (
	"fmt"
	"log/slog"

	"github.com/farshidrezaei/vidonyx/filtergraph"
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
func (r *CompilationResult) Mermaid() (string, error) {
	return visualizer.ToMermaid(r.Graph)
}

// DOT returns the Graphviz DOT representation of the compiled filtergraph DAG.
func (r *CompilationResult) DOT() (string, error) {
	return visualizer.ToDOT(r.Graph)
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
func (c *Compiler) SetEncodingOptions(opts EncodingOptions) *Compiler {
	c.encoding = opts
	return c
}

// Compile translates the given timeline and output path into a complete CompilationResult.
func (c *Compiler) Compile(tl *timeline.Timeline, outputPath string) (*CompilationResult, error) {
	if err := timeline.Validate(tl); err != nil {
		return nil, fmt.Errorf("compiler: timeline validation failed: %w", err)
	}

	g := filtergraph.NewGraph()

	// 1. Discover unique input sources and map to FFmpeg input indices (0, 1, 2, ...)
	inputIndexMap := make(map[string]int)
	uniqueInputs := make([]string, 0)

	for _, track := range tl.Tracks {
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
	for _, track := range tl.Tracks {
		for _, clip := range track.Clips {
			inIdx := inputIndexMap[clip.Source]

			// Handle Video
			if track.Kind == timeline.TrackKindVideo || track.Kind == timeline.TrackKindOverlay {
				rawVideoIn := &filtergraph.Pad{
					ID:         fmt.Sprintf("%d:v", inIdx),
					StreamType: filtergraph.StreamTypeVideo,
					IsInput:    false,
				}

				// Clip video pipeline (trim, speed, opacity, scale)
				clipVPad, err := ProcessClipVideo(g, rawVideoIn, clip, tl)
				if err != nil {
					return nil, fmt.Errorf("compiler: failed processing video clip %q: %w", clip.ID, err)
				}

				// Normalize clip stream to canvas bounds if needed
				normVPad, err := InjectVideoNormalizer(g, clipVPad, tl.Canvas, tl.FPS)
				if err != nil {
					return nil, fmt.Errorf("compiler: failed normalizing video clip %q: %w", clip.ID, err)
				}

				processedVideoPads = append(processedVideoPads, &ClipVideoPad{
					Clip:   clip,
					ZIndex: track.ZIndex,
					Pad:    normVPad,
				})
			}

			// Handle Audio
			if (track.Kind == timeline.TrackKindAudio || track.Kind == timeline.TrackKindVideo) && !track.Muted {
				rawAudioIn := &filtergraph.Pad{
					ID:         fmt.Sprintf("%d:a", inIdx),
					StreamType: filtergraph.StreamTypeAudio,
					IsInput:    false,
				}

				clipAPad, err := ProcessClipAudio(g, rawAudioIn, clip)
				if err != nil {
					return nil, fmt.Errorf("compiler: failed processing audio clip %q: %w", clip.ID, err)
				}

				normAPad, err := InjectAudioNormalizer(g, clipAPad)
				if err != nil {
					return nil, fmt.Errorf("compiler: failed normalizing audio clip %q: %w", clip.ID, err)
				}

				processedAudioPads = append(processedAudioPads, normAPad)
			}
		}
	}

	// 3. Compose final video stream
	finalVideoPad, err := BuildVideoCompositor(g, tl, processedVideoPads)
	if err != nil {
		return nil, fmt.Errorf("compiler: failed building video composition: %w", err)
	}

	// Rename final video pad to out_v
	outVLabel := "out_v"
	finalVideoPad.ID = outVLabel

	// 4. Mix final audio stream
	finalAudioPad, err := BuildAudioMixer(g, tl, processedAudioPads)
	if err != nil {
		return nil, fmt.Errorf("compiler: failed building audio mixer: %w", err)
	}

	// Rename final audio pad to out_a
	outALabel := "out_a"
	finalAudioPad.ID = outALabel

	// 5. Run Graph Optimization Passes
	if err := c.pipeline.Execute(g); err != nil {
		return nil, fmt.Errorf("compiler: graph optimization failed: %w", err)
	}

	// 6. Format filter_complex string
	filterComplexStr, err := g.FormattedFilterComplex()
	if err != nil {
		return nil, fmt.Errorf("compiler: failed serializing filter_complex: %w", err)
	}

	// 7. Emit CLI Arguments
	args := BuildFFmpegArgs(uniqueInputs, filterComplexStr, outVLabel, outALabel, outputPath, c.encoding)

	return &CompilationResult{
		Graph:         g,
		Inputs:        uniqueInputs,
		Args:          args,
		FilterComplex: filterComplexStr,
		OutVLabel:     outVLabel,
		OutALabel:     outALabel,
		OutputPath:    outputPath,
	}, nil
}
