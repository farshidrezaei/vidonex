// Package waveform provides animated audio waveform and frequency spectrum video generation for podcast audiograms and music visualization.
package waveform

import (
	"fmt"

	"github.com/farshidrezaei/vidonex/filtergraph"
	"github.com/farshidrezaei/vidonex/types"
)

// Mode defines the visual rendering style of the audio waveform.
type Mode string

// Waveform visual mode constants.
const (
	// ModeLine renders a continuous line waveform.
	ModeLine Mode = "line"
	// ModePeakToPeak renders solid peak-to-peak bar waveforms (standard for podcast audiograms).
	ModePeakToPeak Mode = "p2p"
	// ModeCenteredLine renders a center-anchored line waveform.
	ModeCenteredLine Mode = "cline"
	// ModeDot renders discrete frequency dots.
	ModeDot Mode = "dot"
)

// Options configures the visual appearance, dimensions, and color of the waveform.
type Options struct {
	Size     types.Size      // Dimension of the waveform (e.g. 1280x240)
	Mode     Mode            // Visual style (p2p, line, cline, dot)
	Color    types.Color     // Foreground waveform color
	Scale    string          // Amplitude scaling: "lin", "log", "sqrt" (default "sqrt")
	FPS      types.Rational  // Output frame rate
	Position types.Point     // Canvas placement coordinate
	Opacity  float64         // Transparency 0.0 to 1.0 (default 1.0)
}

// DefaultOptions returns standard modern podcast waveform settings.
func DefaultOptions() Options {
	return Options{
		Size:     types.NewSize(1200, 200),
		Mode:     ModePeakToPeak,
		Color:    types.RGB(0, 230, 180), // Neon cyan/green
		Scale:    "sqrt",
		FPS:      types.FPS30,
		Position: types.Point{X: 360, Y: 750},
		Opacity:  1.0,
	}
}

// ApplyWaveformVisualizer converts an audio stream into an animated transparent video stream of the waveform.
func ApplyWaveformVisualizer(graph *filtergraph.Graph, nodeID string, audioInputPad *filtergraph.Pad, options Options) (*filtergraph.Pad, error) {
	if options.Size.IsZero() {
		options.Size = types.NewSize(1200, 200)
	}
	if options.Mode == "" {
		options.Mode = ModePeakToPeak
	}
	if options.Scale == "" {
		options.Scale = "sqrt"
	}
	if options.FPS.IsZero() {
		options.FPS = types.FPS30
	}
	if options.Opacity <= 0 {
		options.Opacity = 1.0
	}

	// 1. showwaves filter: Audio In -> Video Out
	wavesNode := graph.NewNode(nodeID, "showwaves")
	wavesNode.SetParam("s", options.Size.String())
	wavesNode.SetParam("mode", string(options.Mode))
	wavesNode.SetParam("scale", options.Scale)
	wavesNode.SetParam("r", options.FPS.FFmpegString())
	wavesNode.SetParam("colors", options.Color.FFmpegColor())

	inputAudio := wavesNode.AddInput(audioInputPad.ID, filtergraph.StreamTypeAudio)
	outputVideo := wavesNode.AddOutput(graph.NextPadID("waves_raw"), filtergraph.StreamTypeVideo)

	if err := graph.Connect(audioInputPad, inputAudio); err != nil {
		return nil, fmt.Errorf("waveform: failed connecting audio to showwaves: %w", err)
	}

	// 2. format=yuva420p to enable transparent alpha blending onto the canvas
	formatNode := graph.NewNode(fmt.Sprintf("%s_alpha", nodeID), "format")
	formatNode.SetParam("pix_fmts", "yuva420p")
	formatInput := formatNode.AddInput(outputVideo.ID, filtergraph.StreamTypeVideo)
	formatOutput := formatNode.AddOutput(graph.NextPadID("waves_yuva"), filtergraph.StreamTypeVideo)

	if err := graph.Connect(outputVideo, formatInput); err != nil {
		return nil, fmt.Errorf("waveform: failed connecting showwaves to format: %w", err)
	}

	return formatOutput, nil
}
