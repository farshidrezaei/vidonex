package filtergraph

import (
	"strings"
)

// DeadCodeEliminationPass removes nodes and subgraphs whose outputs are not connected
// to designated output sinks (e.g. "[out_v]", "[out_a]", or terminal sink nodes).
type DeadCodeEliminationPass struct {
	// SinkPadLabels is an optional list of designated sink pad labels. If empty, defaults to any node with unconnected output pads named "out" or sink nodes.
	SinkPadLabels []string
}

// Name returns the unique pass identifier.
func (pass *DeadCodeEliminationPass) Name() string {
	return "DeadCodeEliminationPass"
}

// Run prunes dead nodes and disconnected subgraphs from the graph.
func (pass *DeadCodeEliminationPass) Run(graph *Graph) error {
	aliveNodes := make(map[string]bool)

	// Step 1: Identify terminal nodes / sink nodes
	for node := range graph.Nodes() {
		// If the node has no outputs, it's inherently a sink
		if node.IsSink() {
			aliveNodes[node.ID] = true
			continue
		}

		// Check if any output pad matches designated sink names (like "out_v", "out_a", "out", etc.)
		// or is an unconsumed leaf output with "out" in its label
		for _, outputPad := range node.Outputs {
			if len(pass.SinkPadLabels) > 0 {
				for _, label := range pass.SinkPadLabels {
					if outputPad.ID == label {
						aliveNodes[node.ID] = true
						break
					}
				}
			} else {
				if strings.HasPrefix(outputPad.ID, "out") || strings.HasSuffix(outputPad.ID, "_out") {
					aliveNodes[node.ID] = true
					break
				}
			}
		}
	}

	// If no designated sinks were matched, treat any node with unconsumed output as alive
	if len(aliveNodes) == 0 {
		for node := range graph.Nodes() {
			for _, outputPad := range node.Outputs {
				if len(graph.GetConsumerPads(outputPad)) == 0 {
					aliveNodes[node.ID] = true
					break
				}
			}
		}
	}

	// Step 2: Backward traversal (DFS) to mark all predecessor nodes as alive
	visited := make(map[string]bool)
	var markAlive func(node *Node)
	markAlive = func(current *Node) {
		if current == nil || visited[current.ID] {
			return
		}
		visited[current.ID] = true
		aliveNodes[current.ID] = true

		for _, inputPad := range current.Inputs {
			if sourceOutputPad, ok := graph.GetSourcePad(inputPad); ok && sourceOutputPad.Node != nil {
				markAlive(sourceOutputPad.Node)
			}
		}
	}

	for node := range graph.Nodes() {
		if aliveNodes[node.ID] {
			markAlive(node)
		}
	}

	// Step 3: Remove any node not marked as alive
	deadNodes := make([]string, 0)
	for node := range graph.Nodes() {
		if !aliveNodes[node.ID] {
			deadNodes = append(deadNodes, node.ID)
		}
	}

	for _, deadID := range deadNodes {
		_ = graph.RemoveNode(deadID)
	}

	return nil
}
