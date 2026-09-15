package effects

import (
	"fmt"

	"github.com/farshidrezaei/vidonex/filtergraph"
	"github.com/farshidrezaei/vidonex/timeline"
)

// XFadeFilter applies an FFmpeg xfade transition between two video streams.
type XFadeFilter struct {
	Transition timeline.TransitionType
	Duration   float64 // in seconds
	Offset     float64 // timestamp offset in seconds where transition starts
}

// Apply attaches an xfade node connecting sourceOutputPadA and sourceOutputPadB.
func (filter XFadeFilter) Apply(
	graph *filtergraph.Graph,
	nodeIdentifier string,
	sourceOutputPadA, sourceOutputPadB *filtergraph.Pad,
) (*filtergraph.Pad, error) {
	node := graph.NewNode(nodeIdentifier, "xfade")
	transitionName := string(filter.Transition)
	if transitionName == "" {
		transitionName = string(timeline.TransitionFade)
	}
	node.SetParam("transition", transitionName)
	node.SetParam("duration", fmt.Sprintf("%.2f", filter.Duration))
	node.SetParam("offset", fmt.Sprintf("%.2f", filter.Offset))

	destinationInputPadA := node.AddInput(sourceOutputPadA.ID, filtergraph.StreamTypeVideo)
	destinationInputPadB := node.AddInput(sourceOutputPadB.ID, filtergraph.StreamTypeVideo)
	outputPad := node.AddOutput(graph.NextPadID("xfade_out"), filtergraph.StreamTypeVideo)

	if err := graph.Connect(sourceOutputPadA, destinationInputPadA); err != nil {
		return nil, err
	}
	if err := graph.Connect(sourceOutputPadB, destinationInputPadB); err != nil {
		return nil, err
	}
	return outputPad, nil
}
