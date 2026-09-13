package visualizer_test

import (
	"strings"
	"testing"

	"github.com/farshidrezaei/vidonex/filtergraph"
	"github.com/farshidrezaei/vidonex/visualizer"
)

func TestVisualizer_MermaidAndDOT(t *testing.T) {
	g := filtergraph.NewGraph()

	n1 := g.NewNode("n1", "scale")
	n1.SetParam("w", 1920).SetParam("h", 1080)
	n1.AddInput("0:v", filtergraph.StreamTypeVideo)
	out1 := n1.AddOutput("v_scale", filtergraph.StreamTypeVideo)

	n2 := g.NewNode("n2", "boxblur")
	n2.SetParam("lr", 5)
	in2 := n2.AddInput("v_scale", filtergraph.StreamTypeVideo)
	n2.AddOutput("out_v", filtergraph.StreamTypeVideo)

	_ = g.Connect(out1, in2)

	// Mermaid test
	mermaid, err := visualizer.ToMermaid(g)
	if err != nil {
		t.Fatalf("ToMermaid failed: %v", err)
	}
	if !strings.Contains(mermaid, "graph LR") || !strings.Contains(mermaid, "n1 -->") {
		t.Fatalf("Unexpected Mermaid output:\n%s", mermaid)
	}

	// DOT test
	dot, err := visualizer.ToDOT(g)
	if err != nil {
		t.Fatalf("ToDOT failed: %v", err)
	}
	if !strings.Contains(dot, "digraph Filtergraph") || !strings.Contains(dot, "\"n1\" -> \"n2\"") {
		t.Fatalf("Unexpected DOT output:\n%s", dot)
	}
}
