package filtergraph_test

import (
	"strings"
	"testing"

	"github.com/farshidrezaei/vidonex/filtergraph"
)

func TestGraph_TopologicalSortTable(t *testing.T) {
	tests := []struct {
		name         string
		buildGraph   func() *filtergraph.Graph
		wantCount    int
		firstNodeID  string
		lastNodeID   string
		shouldErr    bool
		errSubstring string
	}{
		{
			name: "linear chain 3 nodes",
			buildGraph: func() *filtergraph.Graph {
				g := filtergraph.NewGraph()
				n1 := g.NewNode("n1", "color")
				o1 := n1.AddOutput("p1", filtergraph.StreamTypeVideo)

				n2 := g.NewNode("n2", "scale")
				i2 := n2.AddInput("p1", filtergraph.StreamTypeVideo)
				o2 := n2.AddOutput("p2", filtergraph.StreamTypeVideo)
				_ = g.Connect(o1, i2)

				n3 := g.NewNode("n3", "format")
				i3 := n3.AddInput("p2", filtergraph.StreamTypeVideo)
				n3.AddOutput("out_v", filtergraph.StreamTypeVideo)
				_ = g.Connect(o2, i3)

				return g
			},
			wantCount:   3,
			firstNodeID: "n1",
			lastNodeID:  "n3",
		},
		{
			name: "diamond graph (n1 splits to n2 and n3, joined in n4)",
			buildGraph: func() *filtergraph.Graph {
				g := filtergraph.NewGraph()
				n1 := g.NewNode("n1", "color")
				o1a := n1.AddOutput("p1a", filtergraph.StreamTypeVideo)
				o1b := n1.AddOutput("p1b", filtergraph.StreamTypeVideo)

				n2 := g.NewNode("n2", "scale")
				i2 := n2.AddInput("p1a", filtergraph.StreamTypeVideo)
				o2 := n2.AddOutput("p2", filtergraph.StreamTypeVideo)
				_ = g.Connect(o1a, i2)

				n3 := g.NewNode("n3", "boxblur")
				i3 := n3.AddInput("p1b", filtergraph.StreamTypeVideo)
				o3 := n3.AddOutput("p3", filtergraph.StreamTypeVideo)
				_ = g.Connect(o1b, i3)

				n4 := g.NewNode("n4", "overlay")
				i4a := n4.AddInput("p2", filtergraph.StreamTypeVideo)
				i4b := n4.AddInput("p3", filtergraph.StreamTypeVideo)
				n4.AddOutput("out_v", filtergraph.StreamTypeVideo)
				_ = g.Connect(o2, i4a)
				_ = g.Connect(o3, i4b)

				return g
			},
			wantCount:   4,
			firstNodeID: "n1",
			lastNodeID:  "n4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := tt.buildGraph()
			sorted, err := g.TopologicalSort()
			if tt.shouldErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if !strings.Contains(err.Error(), tt.errSubstring) {
					t.Errorf("error %v does not contain %q", err, tt.errSubstring)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(sorted) != tt.wantCount {
				t.Fatalf("expected %d nodes, got %d", tt.wantCount, len(sorted))
			}
			if tt.firstNodeID != "" && sorted[0].ID != tt.firstNodeID {
				t.Errorf("first node = %s, want %s", sorted[0].ID, tt.firstNodeID)
			}
			if tt.lastNodeID != "" && sorted[len(sorted)-1].ID != tt.lastNodeID {
				t.Errorf("last node = %s, want %s", sorted[len(sorted)-1].ID, tt.lastNodeID)
			}
		})
	}
}

func TestGraph_ConnectErrorsTable(t *testing.T) {
	tests := []struct {
		name         string
		setupAndTest func() error
		errSubstring string
	}{
		{
			name: "stream type mismatch (video to audio)",
			setupAndTest: func() error {
				g := filtergraph.NewGraph()
				n1 := g.NewNode("n1", "scale")
				o1 := n1.AddOutput("v1", filtergraph.StreamTypeVideo)

				n2 := g.NewNode("n2", "volume")
				i2 := n2.AddInput("v1", filtergraph.StreamTypeAudio)
				return g.Connect(o1, i2)
			},
			errSubstring: "mismatch",
		},
		{
			name: "connect to already connected input pad",
			setupAndTest: func() error {
				g := filtergraph.NewGraph()
				n1 := g.NewNode("n1", "color")
				o1 := n1.AddOutput("c1", filtergraph.StreamTypeVideo)

				n2 := g.NewNode("n2", "color")
				o2 := n2.AddOutput("c2", filtergraph.StreamTypeVideo)

				n3 := g.NewNode("n3", "scale")
				i3 := n3.AddInput("in", filtergraph.StreamTypeVideo)

				_ = g.Connect(o1, i3)
				return g.Connect(o2, i3)
			},
			errSubstring: "already connected",
		},
		{
			name: "cycle detection (direct 2-node cycle)",
			setupAndTest: func() error {
				g := filtergraph.NewGraph()
				n1 := g.NewNode("n1", "scale")
				o1 := n1.AddOutput("o1", filtergraph.StreamTypeVideo)
				i1 := n1.AddInput("i1", filtergraph.StreamTypeVideo)

				n2 := g.NewNode("n2", "boxblur")
				i2 := n2.AddInput("i2", filtergraph.StreamTypeVideo)
				o2 := n2.AddOutput("o2", filtergraph.StreamTypeVideo)

				if err := g.Connect(o1, i2); err != nil {
					return err
				}
				return g.Connect(o2, i1)
			},
			errSubstring: "cycle",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.setupAndTest()
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.errSubstring) {
				t.Errorf("expected error to contain %q, got: %v", tt.errSubstring, err)
			}
		})
	}
}

func TestGraph_AutoSplitPassTable(t *testing.T) {
	tests := []struct {
		name           string
		streamType     filtergraph.StreamType
		consumerCount  int
		expectedFilter string
	}{
		{
			name:           "video 2 consumers split",
			streamType:     filtergraph.StreamTypeVideo,
			consumerCount:  2,
			expectedFilter: "split",
		},
		{
			name:           "video 4 consumers split",
			streamType:     filtergraph.StreamTypeVideo,
			consumerCount:  4,
			expectedFilter: "split",
		},
		{
			name:           "audio 3 consumers asplit",
			streamType:     filtergraph.StreamTypeAudio,
			consumerCount:  3,
			expectedFilter: "asplit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := filtergraph.NewGraph()
			srcNode := g.NewNode("src", "source")
			srcOut := srcNode.AddOutput("s_out", tt.streamType)

			for i := 0; i < tt.consumerCount; i++ {
				cons := g.NewNode(g.NextPadID("cons"), "filter")
				inPad := cons.AddInput("s_out", tt.streamType)
				cons.AddOutput("c_out", tt.streamType)
				if err := g.Connect(srcOut, inPad); err != nil {
					t.Fatalf("failed connecting consumer %d: %v", i, err)
				}
			}

			splitPass := &filtergraph.AutoSplitPass{}
			if err := splitPass.Run(g); err != nil {
				t.Fatalf("AutoSplitPass failed: %v", err)
			}

			// Graph should now have 1 (src) + 1 (split) + N (consumers)
			expectedNodes := 2 + tt.consumerCount
			if g.NodeCount() != expectedNodes {
				t.Fatalf("NodeCount = %d, want %d", g.NodeCount(), expectedNodes)
			}

			filterStr, err := g.FormattedFilterComplex()
			if err != nil {
				t.Fatalf("FormattedFilterComplex failed: %v", err)
			}
			if !strings.Contains(filterStr, tt.expectedFilter) {
				t.Errorf("FormattedFilterComplex should contain %q:\n%s", tt.expectedFilter, filterStr)
			}
		})
	}
}

func TestGraph_DeadCodeEliminationPassTable(t *testing.T) {
	tests := []struct {
		name          string
		sinkLabels    []string
		expectedAlive int
	}{
		{
			name:          "custom sink label out_v prunes orphan",
			sinkLabels:    []string{"out_v"},
			expectedAlive: 2, // src1 + scale
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := filtergraph.NewGraph()

			// Live branch: src1 -> scale -> [out_v]
			src1 := g.NewNode("src1", "color")
			out1 := src1.AddOutput("c1", filtergraph.StreamTypeVideo)
			scale := g.NewNode("scale", "scale")
			inScale := scale.AddInput("c1", filtergraph.StreamTypeVideo)
			scale.AddOutput("out_v", filtergraph.StreamTypeVideo)
			_ = g.Connect(out1, inScale)

			// Dead branch: orphan -> dead_blur
			orphan := g.NewNode("orphan", "color")
			outOrphan := orphan.AddOutput("dead_c", filtergraph.StreamTypeVideo)
			deadBlur := g.NewNode("dead_blur", "boxblur")
			inDeadBlur := deadBlur.AddInput("dead_c", filtergraph.StreamTypeVideo)
			deadBlur.AddOutput("unused_output", filtergraph.StreamTypeVideo)
			_ = g.Connect(outOrphan, inDeadBlur)

			dce := &filtergraph.DeadCodeEliminationPass{
				SinkPadLabels: tt.sinkLabels,
			}
			if err := dce.Run(g); err != nil {
				t.Fatalf("DCE failed: %v", err)
			}

			if g.NodeCount() != tt.expectedAlive {
				t.Fatalf("NodeCount after DCE = %d, want %d", g.NodeCount(), tt.expectedAlive)
			}
		})
	}
}
