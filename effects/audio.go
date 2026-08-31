// Package effects provides pre-built, type-safe video and audio filters for FFmpeg filtergraphs.
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
func (filter VolumeFilter) Apply(graph *filtergraph.Graph, id string, inputPad *filtergraph.Pad) (*filtergraph.Pad, error) {
	node := graph.NewNode(id, "volume")
	node.SetParam("volume", fmt.Sprintf("%.2f", filter.Volume))
	in := node.AddInput(inputPad.ID, filtergraph.StreamTypeAudio)
	out := node.AddOutput(graph.NextPadID("volume_out"), filtergraph.StreamTypeAudio)
	return out, graph.Connect(inputPad, in)
}

// AcrossFadeFilter cross-fades two audio streams.
type AcrossFadeFilter struct {
	Duration float64 // in seconds
	Curve1   string  // e.g. "tri", "qsin", "esin", "log"
	Curve2   string
}

// Apply connects two audio pads into an acrossfade filter.
func (filter AcrossFadeFilter) Apply(graph *filtergraph.Graph, id string, inputPadA, inputPadB *filtergraph.Pad) (*filtergraph.Pad, error) {
	node := graph.NewNode(id, "acrossfade")
	dur := filter.Duration
	if dur <= 0 {
		dur = 1.0
	}
	node.SetParam("d", fmt.Sprintf("%.2f", dur))
	if filter.Curve1 != "" {
		node.SetParam("c1", filter.Curve1)
	}
	if filter.Curve2 != "" {
		node.SetParam("c2", filter.Curve2)
	}

	inA := node.AddInput(inputPadA.ID, filtergraph.StreamTypeAudio)
	inB := node.AddInput(inputPadB.ID, filtergraph.StreamTypeAudio)
	out := node.AddOutput(graph.NextPadID("across_out"), filtergraph.StreamTypeAudio)

	if err := graph.Connect(inputPadA, inA); err != nil {
		return nil, err
	}
	if err := graph.Connect(inputPadB, inB); err != nil {
		return nil, err
	}
	return out, nil
}
