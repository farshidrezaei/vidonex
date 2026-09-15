// Package effects provides pre-built, type-safe video and audio filters for FFmpeg filtergraphs.
package effects

import (
	"errors"
	"fmt"

	"github.com/farshidrezaei/vidonex/filtergraph"
)

// LUT3DFilter applies an industry-standard 3D Cube Look-Up Table (.cube) for cinematic color grading.
type LUT3DFilter struct {
	// FilePath specifies the absolute or relative path to the .cube or .3dl file.
	FilePath string
	// Interpolation specifies the 3D LUT interpolation method: "tetrahedral", "trilinear", or "nearest".
	Interpolation string
}

// Apply attaches a lut3d filter node to the given graph.
func (filter LUT3DFilter) Apply(
	graph *filtergraph.Graph,
	nodeIdentifier string,
	sourceOutputPad *filtergraph.Pad,
) (*filtergraph.Pad, error) {
	if filter.FilePath == "" {
		return nil, errors.New("effects: LUT3DFilter requires a non-empty FilePath")
	}

	interpolationMethod := filter.Interpolation
	if interpolationMethod == "" {
		interpolationMethod = "tetrahedral"
	}

	node := graph.NewNode(nodeIdentifier, "lut3d")
	node.SetParam("file", filter.FilePath)
	node.SetParam("interp", interpolationMethod)

	destinationInputPad := node.AddInput(sourceOutputPad.ID, filtergraph.StreamTypeVideo)
	outputPad := node.AddOutput(graph.NextPadID("lut3d_out"), filtergraph.StreamTypeVideo)

	if err := graph.Connect(sourceOutputPad, destinationInputPad); err != nil {
		return nil, err
	}

	return outputPad, nil
}

// ColorBalanceAdjustments holds RGB color balance offsets between -1.0 and 1.0.
type ColorBalanceAdjustments struct {
	Red   float64
	Green float64
	Blue  float64
}

// ColorBalanceFilter adjusts color balance across shadow, midtone, and highlight regions.
type ColorBalanceFilter struct {
	Shadows    ColorBalanceAdjustments
	Midtones   ColorBalanceAdjustments
	Highlights ColorBalanceAdjustments
}

// Apply attaches a colorbalance filter node to the given graph.
func (filter ColorBalanceFilter) Apply(
	graph *filtergraph.Graph,
	nodeIdentifier string,
	sourceOutputPad *filtergraph.Pad,
) (*filtergraph.Pad, error) {
	node := graph.NewNode(nodeIdentifier, "colorbalance")

	if filter.Shadows.Red != 0 {
		node.SetParam("rs", fmt.Sprintf("%.2f", filter.Shadows.Red))
	}
	if filter.Shadows.Green != 0 {
		node.SetParam("gs", fmt.Sprintf("%.2f", filter.Shadows.Green))
	}
	if filter.Shadows.Blue != 0 {
		node.SetParam("bs", fmt.Sprintf("%.2f", filter.Shadows.Blue))
	}

	if filter.Midtones.Red != 0 {
		node.SetParam("rm", fmt.Sprintf("%.2f", filter.Midtones.Red))
	}
	if filter.Midtones.Green != 0 {
		node.SetParam("gm", fmt.Sprintf("%.2f", filter.Midtones.Green))
	}
	if filter.Midtones.Blue != 0 {
		node.SetParam("bm", fmt.Sprintf("%.2f", filter.Midtones.Blue))
	}

	if filter.Highlights.Red != 0 {
		node.SetParam("rh", fmt.Sprintf("%.2f", filter.Highlights.Red))
	}
	if filter.Highlights.Green != 0 {
		node.SetParam("gh", fmt.Sprintf("%.2f", filter.Highlights.Green))
	}
	if filter.Highlights.Blue != 0 {
		node.SetParam("bh", fmt.Sprintf("%.2f", filter.Highlights.Blue))
	}

	destinationInputPad := node.AddInput(sourceOutputPad.ID, filtergraph.StreamTypeVideo)
	outputPad := node.AddOutput(graph.NextPadID("colorbalance_out"), filtergraph.StreamTypeVideo)

	if err := graph.Connect(sourceOutputPad, destinationInputPad); err != nil {
		return nil, err
	}

	return outputPad, nil
}

// ColorGradingFilter controls primary image grading parameters: contrast, brightness, saturation, gamma, temperature, tint, highlights, shadows, whites, blacks, blur, and sharpen.
type ColorGradingFilter struct {
	// Contrast multiplier (e.g. 1.0 = normal, 1.2 = punchier). Default is 1.0.
	Contrast float64
	// Brightness offset (e.g. 0.0 = normal, 0.05 = brighter). Default is 0.0.
	Brightness float64
	// Saturation multiplier (e.g. 1.0 = normal, 0.0 = black & white, 1.3 = vibrant). Default is 1.0.
	Saturation float64
	// Gamma multiplier (e.g. 1.0 = normal). Default is 1.0.
	Gamma float64
	// Temperature shifts warmth (-1.0 to 1.0: negative is cool/blue, positive is warm/orange).
	Temperature float64
	// Tint shifts tint (-1.0 to 1.0: negative is green, positive is magenta).
	Tint float64
	// Highlights shifts high luminance regions (-1.0 to 1.0).
	Highlights float64
	// Shadows shifts low luminance regions (-1.0 to 1.0).
	Shadows float64
	// Whites expands/compresses white peak points (-1.0 to 1.0).
	Whites float64
	// Blacks expands/compresses black floor points (-1.0 to 1.0).
	Blacks float64
	// Blur applies Gaussian blur radius (0.0 to 100.0).
	Blur float64
	// Sharpen applies unsharp mask strength (0.0 to 100.0).
	Sharpen float64
}

// Apply attaches an eq (and optional colorbalance, gblur, unsharp) filter node to the given graph.
func (filter ColorGradingFilter) Apply(
	graph *filtergraph.Graph,
	nodeIdentifier string,
	sourceOutputPad *filtergraph.Pad,
) (*filtergraph.Pad, error) {
	contrastValue := filter.Contrast
	if contrastValue == 0 {
		contrastValue = 1.0
	}

	saturationValue := filter.Saturation
	if saturationValue == 0 {
		saturationValue = 1.0
	}

	gammaValue := filter.Gamma
	if gammaValue == 0 {
		gammaValue = 1.0
	}

	eqNode := graph.NewNode(nodeIdentifier+"_eq", "eq")
	eqNode.SetParam("contrast", fmt.Sprintf("%.2f", contrastValue))
	eqNode.SetParam("brightness", fmt.Sprintf("%.2f", filter.Brightness))
	eqNode.SetParam("saturation", fmt.Sprintf("%.2f", saturationValue))
	eqNode.SetParam("gamma", fmt.Sprintf("%.2f", gammaValue))

	eqDestinationInputPad := eqNode.AddInput(sourceOutputPad.ID, filtergraph.StreamTypeVideo)
	eqOutputPad := eqNode.AddOutput(graph.NextPadID("eq_out"), filtergraph.StreamTypeVideo)

	if err := graph.Connect(sourceOutputPad, eqDestinationInputPad); err != nil {
		return nil, err
	}

	currentPad := eqOutputPad

	// If Temperature, Tint, Highlights, Shadows, Whites, or Blacks are non-zero, chain colorbalance
	hasColorBalance := filter.Temperature != 0 || filter.Tint != 0 ||
		filter.Highlights != 0 || filter.Shadows != 0 ||
		filter.Whites != 0 || filter.Blacks != 0

	if hasColorBalance {
		shadowOffset := filter.Shadows*0.3 - filter.Blacks*0.2
		highlightOffset := filter.Highlights*0.3 + filter.Whites*0.2

		balanceFilter := ColorBalanceFilter{
			Shadows: ColorBalanceAdjustments{
				Red:   shadowOffset,
				Green: shadowOffset,
				Blue:  shadowOffset,
			},
			Midtones: ColorBalanceAdjustments{
				Red:   filter.Temperature * 0.3,
				Blue:  -filter.Temperature * 0.3,
				Green: -filter.Tint * 0.3,
			},
			Highlights: ColorBalanceAdjustments{
				Red:   highlightOffset + (filter.Temperature * 0.2),
				Blue:  highlightOffset - (filter.Temperature * 0.2),
				Green: highlightOffset - (filter.Tint * 0.2),
			},
		}
		balancePad, err := balanceFilter.Apply(graph, nodeIdentifier+"_temp", currentPad)
		if err != nil {
			return nil, err
		}
		currentPad = balancePad
	}

	// If Blur is greater than 0, chain gblur filter
	if filter.Blur > 0 {
		blurNode := graph.NewNode(nodeIdentifier+"_blur", "gblur")
		blurNode.SetParam("sigma", fmt.Sprintf("%.2f", filter.Blur*0.5))
		blurInputPad := blurNode.AddInput(currentPad.ID, filtergraph.StreamTypeVideo)
		blurOutputPad := blurNode.AddOutput(graph.NextPadID("blur_out"), filtergraph.StreamTypeVideo)

		if err := graph.Connect(currentPad, blurInputPad); err != nil {
			return nil, err
		}
		currentPad = blurOutputPad
	}

	// If Sharpen is greater than 0, chain unsharp filter
	if filter.Sharpen > 0 {
		sharpenNode := graph.NewNode(nodeIdentifier+"_sharpen", "unsharp")
		sharpenNode.SetParam("luma_msize_x", "5")
		sharpenNode.SetParam("luma_msize_y", "5")
		sharpenNode.SetParam("luma_amount", fmt.Sprintf("%.2f", filter.Sharpen*0.03))
		sharpenInputPad := sharpenNode.AddInput(currentPad.ID, filtergraph.StreamTypeVideo)
		sharpenOutputPad := sharpenNode.AddOutput(graph.NextPadID("sharpen_out"), filtergraph.StreamTypeVideo)

		if err := graph.Connect(currentPad, sharpenInputPad); err != nil {
			return nil, err
		}
		currentPad = sharpenOutputPad
	}

	return currentPad, nil
}
