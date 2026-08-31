package filtergraph_test

import (
	"strings"
	"testing"

	"github.com/farshidrezaei/vidonyx/filtergraph"
)

func TestGraph_BasicLinearChain(t *testing.T) {
	g := filtergraph.NewGraph()

	// [0:v] -> scale -> [v1] -> rotate -> [out_v]
	scaleNode := g.NewNode("scale_1", "scale")
	scaleNode.SetParam("w", 1920).SetParam("h", 1080)
	in0 := scaleNode.AddInput("0:v", filtergraph.StreamTypeVideo)
	out0 := scaleNode.AddOutput("v1", filtergraph.StreamTypeVideo)

	rotateNode := g.NewNode("rotate_1", "rotate")
	rotateNode.SetParam("a", "PI/2")
	in1 := rotateNode.AddInput("v1", filtergraph.StreamTypeVideo)
	out1 := rotateNode.AddOutput("out_v", filtergraph.StreamTypeVideo)

	if err := g.Connect(out0, in1); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	sorted, err := g.TopologicalSort()
	if err != nil {
		t.Fatalf("TopologicalSort failed: %v", err)
	}

	if len(sorted) != 2 {
		t.Fatalf("Expected 2 nodes, got %d", len(sorted))
	}
	if sorted[0].ID != "scale_1" || sorted[1].ID != "rotate_1" {
		t.Fatalf("Unexpected topological order: %v, %v", sorted[0].ID, sorted[1].ID)
	}

	filterStr, err := g.FormattedFilterComplex()
	if err != nil {
		t.Fatalf("FormattedFilterComplex failed: %v", err)
	}

	expectedPrefix := "[0:v]scale=w=1920:h=1080[v1];\n[v1]rotate=a=PI/2[out_v]"
	if filterStr != expectedPrefix {
		t.Fatalf("Unexpected filter complex string:\nGot:\n%s\nExpected:\n%s", filterStr, expectedPrefix)
	}

	_ = in0
	_ = out1
}

func TestGraph_CycleDetection(t *testing.T) {
	g := filtergraph.NewGraph()

	nodeA := g.NewNode("node_A", "scale")
	inA := nodeA.AddInput("in_a", filtergraph.StreamTypeVideo)
	outA := nodeA.AddOutput("out_a", filtergraph.StreamTypeVideo)

	nodeB := g.NewNode("node_B", "rotate")
	inB := nodeB.AddInput("in_b", filtergraph.StreamTypeVideo)
	outB := nodeB.AddOutput("out_b", filtergraph.StreamTypeVideo)

	if err := g.Connect(outA, inB); err != nil {
		t.Fatalf("First connection failed: %v", err)
	}

	// Connecting B -> A creates a cycle!
	err := g.Connect(outB, inA)
	if err == nil || !strings.Contains(err.Error(), "cycle") {
		t.Fatalf("Expected cycle detection error, got: %v", err)
	}
}

func TestGraph_AutoSplitPass(t *testing.T) {
	g := filtergraph.NewGraph()

	// Single source node feeding two consumers
	srcNode := g.NewNode("src", "color")
	srcOut := srcNode.AddOutput("c1", filtergraph.StreamTypeVideo)

	consumerA := g.NewNode("cons_a", "scale")
	inA := consumerA.AddInput("c1", filtergraph.StreamTypeVideo)
	consumerA.AddOutput("out_a", filtergraph.StreamTypeVideo)

	consumerB := g.NewNode("cons_b", "boxblur")
	inB := consumerB.AddInput("c1", filtergraph.StreamTypeVideo)
	consumerB.AddOutput("out_b", filtergraph.StreamTypeVideo)

	if err := g.Connect(srcOut, inA); err != nil {
		t.Fatalf("Connect A failed: %v", err)
	}
	if err := g.Connect(srcOut, inB); err != nil {
		t.Fatalf("Connect B failed: %v", err)
	}

	// Run AutoSplitPass
	splitPass := &filtergraph.AutoSplitPass{}
	if err := splitPass.Run(g); err != nil {
		t.Fatalf("AutoSplitPass failed: %v", err)
	}

	// Graph should now have 4 nodes: src, split, cons_a, cons_b
	if g.NodeCount() != 4 {
		t.Fatalf("Expected 4 nodes after AutoSplit, got %d", g.NodeCount())
	}

	filterStr, err := g.FormattedFilterComplex()
	if err != nil {
		t.Fatalf("FormattedFilterComplex failed: %v", err)
	}

	if !strings.Contains(filterStr, "split") {
		t.Fatalf("Expected split filter in output string:\n%s", filterStr)
	}
}

func TestGraph_DeadCodeEliminationPass(t *testing.T) {
	g := filtergraph.NewGraph()

	// Live pipeline: src1 -> scale -> [out_v]
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

	if g.NodeCount() != 4 {
		t.Fatalf("Expected 4 nodes initially, got %d", g.NodeCount())
	}

	dce := &filtergraph.DeadCodeEliminationPass{
		SinkPadLabels: []string{"out_v"},
	}
	if err := dce.Run(g); err != nil {
		t.Fatalf("DCE pass failed: %v", err)
	}

	if g.NodeCount() != 2 {
		t.Fatalf("Expected 2 nodes after DCE, got %d", g.NodeCount())
	}
	if _, exists := g.GetNode("orphan"); exists {
		t.Fatal("Expected orphan node to be removed")
	}
	if _, exists := g.GetNode("dead_blur"); exists {
		t.Fatal("Expected dead_blur node to be removed")
	}
}
