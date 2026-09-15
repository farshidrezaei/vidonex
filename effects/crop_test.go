package effects_test

import (
	"fmt"
	"testing"

	"github.com/farshidrezaei/vidonex/effects"
	"github.com/farshidrezaei/vidonex/filtergraph"
)

func TestCropFilter_Table(t *testing.T) {
	testScenarios := []struct {
		name          string
		cropFilter    effects.CropFilter
		streamType    filtergraph.StreamType
		expectedNode  string
		expectedWExpr string
		expectedXExpr string
		expectError   bool
	}{
		{
			name: "fixed pixel dimensions crop",
			cropFilter: effects.CropFilter{
				Width:  1280,
				Height: 720,
				X:      100,
				Y:      50,
			},
			streamType:    filtergraph.StreamTypeVideo,
			expectedNode:  "crop",
			expectedWExpr: "1280",
			expectedXExpr: "100",
			expectError:   false,
		},
		{
			name: "percentage inset crop",
			cropFilter: effects.CropFilter{
				TopInset:    0.10,
				BottomInset: 0.10,
				LeftInset:   0.05,
				RightInset:  0.05,
			},
			streamType:    filtergraph.StreamTypeVideo,
			expectedNode:  "crop",
			expectedWExpr: "in_w-in_w*0.0500-in_w*0.0500",
			expectedXExpr: "in_w*0.0500",
			expectError:   false,
		},
		{
			name: "error on audio stream type",
			cropFilter: effects.CropFilter{
				Width:  640,
				Height: 480,
			},
			streamType:  filtergraph.StreamTypeAudio,
			expectError: true,
		},
	}

	for _, scenario := range testScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			graph := filtergraph.NewGraph()
			inputNode := graph.NewNode("source", "testsrc")
			sourcePad := inputNode.AddOutput("out_pad", scenario.streamType)

			outputPad, err := scenario.cropFilter.Apply(graph, "crop_node", sourcePad)
			if scenario.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if outputPad == nil {
				t.Fatalf("expected non-nil outputPad")
			}

			node, exists := graph.GetNode("crop_node")
			if !exists || node == nil {
				t.Fatalf("expected crop_node in graph")
			}

			if node.FilterName != scenario.expectedNode {
				t.Errorf("filter = %q, want %q", node.FilterName, scenario.expectedNode)
			}

			if scenario.expectedWExpr != "" {
				wVal, ok := node.GetParam("w")
				if !ok || fmt.Sprintf("%v", wVal) != scenario.expectedWExpr {
					t.Errorf("param 'w' = %v, want %s", wVal, scenario.expectedWExpr)
				}
			}

			if scenario.expectedXExpr != "" {
				xVal, ok := node.GetParam("x")
				if !ok || fmt.Sprintf("%v", xVal) != scenario.expectedXExpr {
					t.Errorf("param 'x' = %v, want %s", xVal, scenario.expectedXExpr)
				}
			}
		})
	}
}
