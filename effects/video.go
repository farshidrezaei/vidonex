package effects

import (
	"fmt"

	"github.com/farshidrezaei/vidonex/filtergraph"
	"github.com/farshidrezaei/vidonex/types"
)

// ScaleFilter configures an FFmpeg scale node.
type ScaleFilter struct {
	Width       int
	Height      int
	AspectRatio string // e.g. "decrease", "increase"
	Flags       string // e.g. "lanczos", "bicubic"
}

// Apply attaches a scale node to the graph.
func (f ScaleFilter) Apply(g *filtergraph.Graph, id string, inPad *filtergraph.Pad) (*filtergraph.Pad, error) {
	node := g.NewNode(id, "scale")
	node.SetParam("w", f.Width).SetParam("h", f.Height)
	if f.AspectRatio != "" {
		node.SetParam("force_original_aspect_ratio", f.AspectRatio)
	}
	if f.Flags != "" {
		node.SetParam("flags", f.Flags)
	}

	in := node.AddInput(inPad.ID, filtergraph.StreamTypeVideo)
	out := node.AddOutput(g.NextPadID("sc_out"), filtergraph.StreamTypeVideo)
	return out, g.Connect(inPad, in)
}

// DrawTextFilter configures dynamic text rendering onto a video stream.
type DrawTextFilter struct {
	Text      string
	FontFile  string
	FontSize  int
	FontColor types.Color
	Box       bool
	BoxColor  types.Color
	X         string
	Y         string
}

// Apply attaches a drawtext node to the graph.
func (f DrawTextFilter) Apply(g *filtergraph.Graph, id string, inPad *filtergraph.Pad) (*filtergraph.Pad, error) {
	node := g.NewNode(id, "drawtext")
	node.SetParam("text", fmt.Sprintf("'%s'", f.Text))
	if f.FontFile != "" {
		node.SetParam("fontfile", fmt.Sprintf("'%s'", f.FontFile))
	}
	if f.FontSize > 0 {
		node.SetParam("fontsize", f.FontSize)
	}
	node.SetParam("fontcolor", f.FontColor.FFmpegColor())
	if f.Box {
		node.SetParam("box", 1)
		node.SetParam("boxcolor", f.BoxColor.FFmpegColor())
	}
	if f.X != "" {
		node.SetParam("x", f.X)
	} else {
		node.SetParam("x", "(w-text_w)/2")
	}
	if f.Y != "" {
		node.SetParam("y", f.Y)
	} else {
		node.SetParam("y", "(h-text_h)/2")
	}

	in := node.AddInput(inPad.ID, filtergraph.StreamTypeVideo)
	out := node.AddOutput(g.NextPadID("dt_out"), filtergraph.StreamTypeVideo)
	return out, g.Connect(inPad, in)
}

// ChromaKeyFilter configures color keying (green/blue screen removal).
type ChromaKeyFilter struct {
	ColorHex   string
	Similarity float64
	Blend      float64
}

// Apply attaches a chromakey node to the graph.
func (f ChromaKeyFilter) Apply(g *filtergraph.Graph, id string, inPad *filtergraph.Pad) (*filtergraph.Pad, error) {
	node := g.NewNode(id, "chromakey")
	if f.ColorHex != "" {
		node.SetParam("color", f.ColorHex)
	} else {
		node.SetParam("color", "0x00FF00") // Green default
	}
	sim := f.Similarity
	if sim <= 0 {
		sim = 0.3
	}
	blend := f.Blend
	if blend <= 0 {
		blend = 0.1
	}
	node.SetParam("similarity", fmt.Sprintf("%.2f", sim))
	node.SetParam("blend", fmt.Sprintf("%.2f", blend))

	in := node.AddInput(inPad.ID, filtergraph.StreamTypeVideo)
	out := node.AddOutput(g.NextPadID("ck_out"), filtergraph.StreamTypeVideo)
	return out, g.Connect(inPad, in)
}
