package visualizer

import (
	"fmt"
	"strings"

	"github.com/farshidrezaei/vidonyx/filtergraph"
)

// ToMermaid generates a Mermaid.js diagram (flowchart LR) visualizing the filtergraph DAG.
func ToMermaid(graph *filtergraph.Graph) (string, error) {
	var stringBuilder strings.Builder
	stringBuilder.WriteString("graph LR\n")

	// 1. Declare Nodes
	for node := range graph.Nodes() {
		filterLabel := escapeMermaid(node.FormattedFilter())
		fmt.Fprintf(&stringBuilder, "    %s[\"%s<br/><i>%s</i>\"]\n", node.ID, node.ID, filterLabel)
	}

	stringBuilder.WriteString("\n")

	// 2. Declare Links / Edges
	for node := range graph.Nodes() {
		for _, outputPad := range node.Outputs {
			consumers := graph.GetConsumerPads(outputPad)
			for _, inputPad := range consumers {
				if inputPad.Node != nil {
					padLabel := escapeMermaid(outputPad.ID)
					fmt.Fprintf(&stringBuilder, "    %s -->|%s| %s\n", node.ID, padLabel, inputPad.Node.ID)
				}
			}
		}
	}

	return stringBuilder.String(), nil
}

func escapeMermaid(s string) string {
	s = strings.ReplaceAll(s, "\"", "#quot;")
	s = strings.ReplaceAll(s, "[", "#91;")
	s = strings.ReplaceAll(s, "]", "#93;")
	return s
}
