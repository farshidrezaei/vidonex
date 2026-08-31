package visualizer

import (
	"fmt"
	"strings"

	"github.com/farshidrezaei/vidonyx/filtergraph"
)

// ToMermaid generates a Mermaid.js diagram (flowchart LR) visualizing the filtergraph DAG.
func ToMermaid(g *filtergraph.Graph) (string, error) {
	var sb strings.Builder
	sb.WriteString("graph LR\n")

	// 1. Declare Nodes
	for node := range g.Nodes() {
		filterLabel := escapeMermaid(node.FormattedFilter())
		sb.WriteString(fmt.Sprintf("    %s[\"%s<br/><i>%s</i>\"]\n", node.ID, node.ID, filterLabel))
	}

	sb.WriteString("\n")

	// 2. Declare Links / Edges
	for node := range g.Nodes() {
		for _, outPad := range node.Outputs {
			consumers := g.GetConsumerPads(outPad)
			for _, inPad := range consumers {
				if inPad.Node != nil {
					padLabel := escapeMermaid(outPad.ID)
					sb.WriteString(fmt.Sprintf("    %s -->|%s| %s\n", node.ID, padLabel, inPad.Node.ID))
				}
			}
		}
	}

	return sb.String(), nil
}

func escapeMermaid(s string) string {
	s = strings.ReplaceAll(s, "\"", "#quot;")
	s = strings.ReplaceAll(s, "[", "#91;")
	s = strings.ReplaceAll(s, "]", "#93;")
	return s
}
