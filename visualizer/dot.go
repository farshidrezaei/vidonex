// Package visualizer generates diagrammatic representations of filtergraph DAGs in Mermaid.js and Graphviz DOT formats.
package visualizer

import (
	"fmt"
	"strings"

	"github.com/farshidrezaei/vidonex/filtergraph"
)

// ToDOT generates a Graphviz DOT representation of the filtergraph.
func ToDOT(graph *filtergraph.Graph) (string, error) {
	var stringBuilder strings.Builder
	stringBuilder.WriteString("digraph Filtergraph {\n")
	stringBuilder.WriteString("    rankdir=LR;\n")
	stringBuilder.WriteString("    node [shape=box, style=\"rounded,filled\", fillcolor=\"#F0F4F8\", fontname=\"Helvetica\"];\n")
	stringBuilder.WriteString("    edge [fontname=\"Helvetica\", fontsize=10];\n\n")

	// Nodes
	for node := range graph.Nodes() {
		filterLabel := strings.ReplaceAll(node.FormattedFilter(), "\"", "\\\"")
		fmt.Fprintf(&stringBuilder, "    \"%s\" [label=\"%s\\n%s\"];\n", node.ID, node.ID, filterLabel)
	}

	stringBuilder.WriteString("\n")

	// Edges
	for node := range graph.Nodes() {
		for _, outputPad := range node.Outputs {
			consumers := graph.GetConsumerPads(outputPad)
			for _, inputPad := range consumers {
				if inputPad.Node != nil {
					fmt.Fprintf(&stringBuilder, "    \"%s\" -> \"%s\" [label=\"[%s]\"];\n", node.ID, inputPad.Node.ID, outputPad.ID)
				}
			}
		}
	}

	stringBuilder.WriteString("}\n")
	return stringBuilder.String(), nil
}
