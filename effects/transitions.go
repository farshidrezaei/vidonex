package effects

import (
	"fmt"

	"github.com/farshidrezaei/vidonyx/filtergraph"
	"github.com/farshidrezaei/vidonyx/timeline"
)

// XFadeFilter applies an FFmpeg xfade transition between two video streams.
type XFadeFilter struct {
	Transition timeline.TransitionType
	Duration   float64 // in seconds
	Offset     float64 // timestamp offset in seconds where transition starts
}

// Apply attaches an xfade node connecting inPadA and inPadB.
func (f XFadeFilter) Apply(g *filtergraph.Graph, id string, inPadA, inPadB *filtergraph.Pad) (*filtergraph.Pad, error) {
	node := g.NewNode(id, "xfade")
	trans := string(f.Transition)
	if trans == "" {
		trans = "fade"
	}
	node.SetParam("transition", trans)
	node.SetParam("duration", fmt.Sprintf("%.2f", f.Duration))
	node.SetParam("offset", fmt.Sprintf("%.2f", f.Offset))

	inA := node.AddInput(inPadA.ID, filtergraph.StreamTypeVideo)
	inB := node.AddInput(inPadB.ID, filtergraph.StreamTypeVideo)
	out := node.AddOutput(g.NextPadID("xfade_out"), filtergraph.StreamTypeVideo)

	if err := g.Connect(inPadA, inA); err != nil {
		return nil, err
	}
	if err := g.Connect(inPadB, inB); err != nil {
		return nil, err
	}
	return out, nil
}
