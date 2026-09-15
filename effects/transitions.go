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

// MapTransitionToFFmpegXFade maps domain TransitionType to the matching FFmpeg xfade transition identifier.
//
// In FFmpeg's vf_xfade.c filter:
// 1. squeezeh compresses along the Y axis (height), which is a vertical squeeze.
//    squeezev compresses along the X axis (width), which is a horizontal squeeze.
//    Therefore:
//    - TransitionSqueezeVert ("squeezev") maps to FFmpeg "squeezeh".
//    - TransitionSqueezeHorz ("squeezeh") maps to FFmpeg "squeezev".
//
// 2. horzopen uses h2 (height/2), opening along the vertical axis (up and down).
//    vertopen uses w2 (width/2), opening along the horizontal axis (like elevator doors).
//    Therefore, to match human naming ("Horizontal Open" = doors open horizontally):
//    - TransitionHorzOpen ("horzopen") maps to FFmpeg "vertopen".
//    - TransitionVertOpen ("vertopen") maps to FFmpeg "horzopen".
//    - TransitionHorzClose ("horzclose") maps to FFmpeg "vertclose".
//    - TransitionVertClose ("vertclose") maps to FFmpeg "horzclose".
func MapTransitionToFFmpegXFade(transitionType timeline.TransitionType) string {
	switch transitionType {
	case timeline.TransitionSqueezeVert:
		return "squeezeh"
	case timeline.TransitionSqueezeHorz:
		return "squeezev"
	case timeline.TransitionHorzOpen:
		return "vertopen"
	case timeline.TransitionVertOpen:
		return "horzopen"
	case timeline.TransitionHorzClose:
		return "vertclose"
	case timeline.TransitionVertClose:
		return "horzclose"
	case "":
		return string(timeline.TransitionFade)
	default:
		return string(transitionType)
	}
}

// Apply attaches an xfade node connecting sourceOutputPadA and sourceOutputPadB.
func (filter XFadeFilter) Apply(
	graph *filtergraph.Graph,
	nodeIdentifier string,
	sourceOutputPadA, sourceOutputPadB *filtergraph.Pad,
) (*filtergraph.Pad, error) {
	node := graph.NewNode(nodeIdentifier, "xfade")
	transitionName := MapTransitionToFFmpegXFade(filter.Transition)
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
