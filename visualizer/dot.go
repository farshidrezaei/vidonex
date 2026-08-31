package visualizer

import (
	"fmt"
	"strings"

	"github.com/farshidrezaei/vidonyx/filtergraph"
)

// ToDOT generates a Graphviz DOT representation of the filtergraph.
func ToDOT(g *filtergraph.Graph) (string, error) {
	var sb strings.Builder
	sb.WriteString("digraph Filtergraph {\n")
	sb.WriteString("    rankdir=LR;\n")
	sb.WriteString("    node [shape=box, style=\"rounded,filled\", fillcolor=\"#F0F4F8\", fontname=\"Helvetica\"];\n")
	sb.WriteString("    edge [fontname=\"Helvetica\", fontsize=10];\n\n")

	// Nodes
	for node := range g.Nodes() {
		filterLabel := strings.ReplaceAll(node.FormattedFilter(), "\"", "\\\"")
		sb.WriteString(fmt.Sprintf("    \"%s\" [label=\"%s\\n%s\"];\n", node.ID, node.ID, filterLabel))
	}

	sb.WriteString("\n")

	// Edges
	for node := range g.Nodes() {
		for _, outPad := range node.Outputs {
			consumers := g.GetConsumerPads(outPad)
			for _, inPad := range consumers {
				if inPad.Node != nil {
					sb.WriteString(fmt.Sprintf("    \"%s\" -> \"%s\" [label=\"[%s]\"];\n", node.ID, inPad.Node.ID, outPad.ID))
				}
			}
		}
	}

	sb.WriteString("}\n")
	return sb.String(), nil
}
