package effects_test

import (
	"testing"

	"github.com/farshidrezaei/vidonex/effects"
	"github.com/farshidrezaei/vidonex/filtergraph"
)

func TestLUT3DFilter_Table(t *testing.T) {
	testCases := []struct {
		name                  string
		filter                effects.LUT3DFilter
		sourceStreamType      filtergraph.StreamType
		expectedInterpolation string
		expectedError         bool
	}{
		{
			name: "Standard 3D LUT with default interpolation",
			filter: effects.LUT3DFilter{
				FilePath: "/path/to/cinematic_teal_orange.cube",
			},
			sourceStreamType:      filtergraph.StreamTypeVideo,
			expectedInterpolation: "tetrahedral",
			expectedError:         false,
		},
		{
			name: "3D LUT with custom trilinear interpolation",
			filter: effects.LUT3DFilter{
				FilePath:      "/path/to/vintage_film.cube",
				Interpolation: "trilinear",
			},
			sourceStreamType:      filtergraph.StreamTypeVideo,
			expectedInterpolation: "trilinear",
			expectedError:         false,
		},
		{
			name: "Empty file path error",
			filter: effects.LUT3DFilter{
				FilePath: "",
			},
			sourceStreamType:      filtergraph.StreamTypeVideo,
			expectedInterpolation: "",
			expectedError:         true,
		},
		{
			name: "Stream type mismatch error (audio pad to video LUT)",
			filter: effects.LUT3DFilter{
				FilePath: "/path/to/lut.cube",
			},
			sourceStreamType:      filtergraph.StreamTypeAudio,
			expectedInterpolation: "",
			expectedError:         true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			graph := filtergraph.NewGraph()
			sourcePad := &filtergraph.Pad{
				ID:         "source_video_pad",
				StreamType: tc.sourceStreamType,
			}

			outputPad, err := tc.filter.Apply(graph, "lut_node", sourcePad)

			if tc.expectedError {
				if err == nil {
					t.Fatalf("expected error for scenario %q, got nil", tc.name)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error for scenario %q: %v", tc.name, err)
			}

			if outputPad == nil {
				t.Fatalf("expected valid output pad, got nil")
			}

			if outputPad.StreamType != filtergraph.StreamTypeVideo {
				t.Errorf("expected video stream type, got %s", outputPad.StreamType)
			}

			node, exists := graph.GetNode("lut_node")
			if !exists || node == nil {
				t.Fatalf("expected lut_node in graph")
			}

			if val, ok := node.GetParam("interp"); !ok || val != tc.expectedInterpolation {
				t.Errorf("expected interp=%s, got %v", tc.expectedInterpolation, val)
			}
		})
	}
}

func TestColorGradingFilter_Table(t *testing.T) {
	testCases := []struct {
		name               string
		filter             effects.ColorGradingFilter
		expectChainedNodes bool
	}{
		{
			name: "Basic contrast and saturation adjustment",
			filter: effects.ColorGradingFilter{
				Contrast:   1.2,
				Brightness: 0.05,
				Saturation: 1.3,
				Gamma:      1.0,
			},
			expectChainedNodes: false,
		},
		{
			name: "Warm temperature color grading with balance chaining",
			filter: effects.ColorGradingFilter{
				Contrast:    1.1,
				Brightness:  0.0,
				Saturation:  1.2,
				Temperature: 0.5,
				Tint:        -0.2,
			},
			expectChainedNodes: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			graph := filtergraph.NewGraph()
			sourcePad := &filtergraph.Pad{
				ID:         "source_video_pad",
				StreamType: filtergraph.StreamTypeVideo,
			}

			outputPad, err := tc.filter.Apply(graph, "grade", sourcePad)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if outputPad == nil {
				t.Fatalf("expected output pad, got nil")
			}

			eqNode, exists := graph.GetNode("grade_eq")
			if !exists || eqNode == nil {
				t.Fatalf("expected grade_eq node in graph")
			}

			if tc.expectChainedNodes {
				tempNode, exists := graph.GetNode("grade_temp")
				if !exists || tempNode == nil {
					t.Fatalf("expected chained grade_temp node in graph")
				}
				if tempNode.FilterName != "colorbalance" {
					t.Errorf("expected colorbalance filter for temperature, got %s", tempNode.FilterName)
				}
			}
		})
	}
}
