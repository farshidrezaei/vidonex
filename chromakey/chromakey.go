// Package chromakey provides green/blue screen removal, edge blending, and despill suppression filters.
package chromakey

import (
	"fmt"

	"github.com/farshidrezaei/vidonex/filtergraph"
	"github.com/farshidrezaei/vidonex/types"
)

// DespillMode specifies which color channel to suppress during despill processing.
type DespillMode string

// Despill mode constants.
const (
	// DespillGreen suppresses green light bounce on skin, clothing, and hair.
	DespillGreen DespillMode = "green"
	// DespillBlue suppresses blue light bounce.
	DespillBlue DespillMode = "blue"
)

// Standard studio color key constants.
var (
	// StudioGreenScreen is the standard RGB pure green keying color.
	StudioGreenScreen = types.RGB(0, 255, 0)
	// StudioBlueScreen is the standard RGB pure blue keying color.
	StudioBlueScreen = types.RGB(0, 0, 255)
)

// Options configures chroma keying thresholds, edge softness, and spill suppression.
type Options struct {
	KeyColor      types.Color // The background color to key out (e.g. 0x00FF00)
	Similarity    float64     // Color similarity tolerance (0.01 to 1.0, default 0.3)
	Blend         float64     // Edge softness / alpha feathering (0.0 to 1.0, default 0.1)
	Despill       bool        // Whether to suppress reflected green/blue spill
	DespillType   DespillMode // Channel to suppress (DespillGreen or DespillBlue)
	DespillExpand float64     // Spill suppression multiplier (default 0.0)
}

// DefaultOptions returns standard studio green screen removal settings with despill enabled.
func DefaultOptions() Options {
	return Options{
		KeyColor:      StudioGreenScreen,
		Similarity:    0.30,
		Blend:         0.10,
		Despill:       true,
		DespillType:   DespillGreen,
		DespillExpand: 0.0,
	}
}

// ApplyChromaKey builds and connects the chromakey and optional despill filter nodes.
func ApplyChromaKey(graph *filtergraph.Graph, nodeID string, inputPad *filtergraph.Pad, options Options) (*filtergraph.Pad, error) {
	if options.Similarity <= 0 {
		options.Similarity = 0.30
	}
	if options.Blend < 0 {
		options.Blend = 0.10
	}

	// 1. Chromakey node
	chromaNode := graph.NewNode(nodeID, "chromakey")
	chromaNode.SetParam("color", options.KeyColor.FFmpegColor())
	chromaNode.SetParam("similarity", fmt.Sprintf("%.2f", options.Similarity))
	chromaNode.SetParam("blend", fmt.Sprintf("%.2f", options.Blend))

	chromaInput := chromaNode.AddInput(inputPad.ID, filtergraph.StreamTypeVideo)
	chromaOutput := chromaNode.AddOutput(graph.NextPadID("chroma_out"), filtergraph.StreamTypeVideo)

	if err := graph.Connect(inputPad, chromaInput); err != nil {
		return nil, fmt.Errorf("chromakey: failed connecting input pad: %w", err)
	}

	currentPad := chromaOutput

	// 2. Optional Despill node to clean up green/blue fringe
	if options.Despill {
		despillType := options.DespillType
		if despillType == "" {
			despillType = DespillGreen
		}

		despillNode := graph.NewNode(fmt.Sprintf("%s_despill", nodeID), "despill")
		despillNode.SetParam("type", string(despillType))
		despillNode.SetParam("expand", fmt.Sprintf("%.2f", options.DespillExpand))

		despillInput := despillNode.AddInput(currentPad.ID, filtergraph.StreamTypeVideo)
		despillOutput := despillNode.AddOutput(graph.NextPadID("despill_out"), filtergraph.StreamTypeVideo)

		if err := graph.Connect(currentPad, despillInput); err != nil {
			return nil, fmt.Errorf("chromakey: failed connecting despill node: %w", err)
		}

		currentPad = despillOutput
	}

	return currentPad, nil
}
