package effects

import (
	"fmt"

	"github.com/farshidrezaei/vidonyx/filtergraph"
)

// VolumeFilter adjusts audio gain.
type VolumeFilter struct {
	Volume float64 // 1.0 = 100%
}

// Apply attaches a volume node to the graph.
func (f VolumeFilter) Apply(g *filtergraph.Graph, id string, inPad *filtergraph.Pad) (*filtergraph.Pad, error) {
	node := g.NewNode(id, "volume")
	node.SetParam("volume", fmt.Sprintf("%.2f", f.Volume))
	in := node.AddInput(inPad.ID, filtergraph.StreamTypeAudio)
	out := node.AddOutput(g.NextPadID("vol_out"), filtergraph.StreamTypeAudio)
	return out, g.Connect(inPad, in)
}

// AcrossFadeFilter cross-fades two audio streams.
type AcrossFadeFilter struct {
	Duration float64 // in seconds
	Curve1   string  // e.g. "tri", "qsin", "esin", "log"
	Curve2   string
}

// Apply connects two audio pads into an acrossfade filter.
func (f AcrossFadeFilter) Apply(g *filtergraph.Graph, id string, inPadA, inPadB *filtergraph.Pad) (*filtergraph.Pad, error) {
	node := g.NewNode(id, "acrossfade")
	dur := f.Duration
	if dur <= 0 {
		dur = 1.0
	}
	node.SetParam("d", fmt.Sprintf("%.2f", dur))
	if f.Curve1 != "" {
		node.SetParam("c1", f.Curve1)
	}
	if f.Curve2 != "" {
		node.SetParam("c2", f.Curve2)
	}

	inA := node.AddInput(inPadA.ID, filtergraph.StreamTypeAudio)
	inB := node.AddInput(inPadB.ID, filtergraph.StreamTypeAudio)
	out := node.AddOutput(g.NextPadID("across_out"), filtergraph.StreamTypeAudio)

	if err := g.Connect(inPadA, inA); err != nil {
		return nil, err
	}
	if err := g.Connect(inPadB, inB); err != nil {
		return nil, err
	}
	return out, nil
}
