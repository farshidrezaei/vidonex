package effects

import (
	"errors"
	"fmt"

	"github.com/farshidrezaei/vidonex/filtergraph"
)

// CropFilter extracts a sub-region of the video stream by applying an FFmpeg crop filter.
type CropFilter struct {
	// Width of the output cropped area. If 0 or negative, it is computed relative to input width (e.g. in_w - left - right).
	Width int
	// Height of the output cropped area. If 0 or negative, it is computed relative to input height (e.g. in_h - top - bottom).
	Height int
	// X coordinate of the left top corner of the crop area.
	X int
	// Y coordinate of the left top corner of the crop area.
	Y int
	// Percentage insets (0.0 to 1.0)
	TopInset    float64
	BottomInset float64
	LeftInset   float64
	RightInset  float64
}

// Apply attaches a crop filter node to the given graph.
func (filter CropFilter) Apply(
	graph *filtergraph.Graph,
	nodeIdentifier string,
	sourceOutputPad *filtergraph.Pad,
) (*filtergraph.Pad, error) {
	if sourceOutputPad.StreamType != filtergraph.StreamTypeVideo {
		return nil, errors.New("effects: CropFilter can only be applied to video streams")
	}

	node := graph.NewNode(nodeIdentifier, "crop")

	if filter.Width > 0 && filter.Height > 0 {
		node.SetParam("w", filter.Width)
		node.SetParam("h", filter.Height)
		node.SetParam("x", filter.X)
		node.SetParam("y", filter.Y)
	} else if filter.TopInset > 0 || filter.BottomInset > 0 || filter.LeftInset > 0 || filter.RightInset > 0 {
		// Calculate crop parameters using FFmpeg expressions
		// in_w - left - right, in_h - top - bottom
		cropWExpr := fmt.Sprintf("in_w-in_w*%.4f-in_w*%.4f", filter.LeftInset, filter.RightInset)
		cropHExpr := fmt.Sprintf("in_h-in_h*%.4f-in_h*%.4f", filter.TopInset, filter.BottomInset)
		cropXExpr := fmt.Sprintf("in_w*%.4f", filter.LeftInset)
		cropYExpr := fmt.Sprintf("in_h*%.4f", filter.TopInset)

		node.SetParam("w", cropWExpr)
		node.SetParam("h", cropHExpr)
		node.SetParam("x", cropXExpr)
		node.SetParam("y", cropYExpr)
	} else {
		// No-op or full dimensions
		node.SetParam("w", "in_w")
		node.SetParam("h", "in_h")
		node.SetParam("x", 0)
		node.SetParam("y", 0)
	}

	destinationInputPad := node.AddInput(sourceOutputPad.ID, filtergraph.StreamTypeVideo)
	outputPad := node.AddOutput(graph.NextPadID("crop_out"), filtergraph.StreamTypeVideo)

	if err := graph.Connect(sourceOutputPad, destinationInputPad); err != nil {
		return nil, err
	}

	return outputPad, nil
}
