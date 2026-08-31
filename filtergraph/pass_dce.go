package filtergraph

import (
	"strings"
)

// DeadCodeEliminationPass removes nodes and subgraphs whose outputs are not connected
// to designated output sinks (e.g. "[out_v]", "[out_a]", or terminal sink nodes).
type DeadCodeEliminationPass struct {
	// Sinks is an optional list of designated sink pad labels. If empty, defaults to any node with unconnected output pads named "out" or sink nodes.
	SinkPadLabels []string
}

func (p *DeadCodeEliminationPass) Name() string {
	return "DeadCodeEliminationPass"
}

func (p *DeadCodeEliminationPass) Run(g *Graph) error {
	aliveNodes := make(map[string]bool)

	// Step 1: Identify terminal nodes / sink nodes
	for node := range g.Nodes() {
		// If the node has no outputs, it's inherently a sink
		if node.IsSink() {
			aliveNodes[node.ID] = true
			continue
		}

		// Check if any output pad matches designated sink names (like "out_v", "out_a", "out", etc.)
		// or is an unconsumed leaf output with "out" in its label
		for _, outPad := range node.Outputs {
			if len(p.SinkPadLabels) > 0 {
				for _, label := range p.SinkPadLabels {
					if outPad.ID == label {
						aliveNodes[node.ID] = true
						break
					}
				}
			} else {
				if strings.HasPrefix(outPad.ID, "out") || strings.HasSuffix(outPad.ID, "_out") {
					aliveNodes[node.ID] = true
					break
				}
			}
		}
	}

	// If no designated sinks were matched, treat any node with unconsumed output as alive
	if len(aliveNodes) == 0 {
		for node := range g.Nodes() {
			for _, outPad := range node.Outputs {
				if len(g.GetConsumerPads(outPad)) == 0 {
					aliveNodes[node.ID] = true
					break
				}
			}
		}
	}

	// Step 2: Backward traversal (DFS) to mark all predecessor nodes as alive
	visited := make(map[string]bool)
	var markAlive func(node *Node)
	markAlive = func(curr *Node) {
		if curr == nil || visited[curr.ID] {
			return
		}
		visited[curr.ID] = true
		aliveNodes[curr.ID] = true

		for _, inPad := range curr.Inputs {
			if srcOut, ok := g.GetSourcePad(inPad); ok && srcOut.Node != nil {
				markAlive(srcOut.Node)
			}
		}
	}

	for node := range g.Nodes() {
		if aliveNodes[node.ID] {
			markAlive(node)
		}
	}

	// Step 3: Remove any node not marked as alive
	deadNodes := make([]string, 0)
	for node := range g.Nodes() {
		if !aliveNodes[node.ID] {
			deadNodes = append(deadNodes, node.ID)
		}
	}

	for _, deadID := range deadNodes {
		_ = g.RemoveNode(deadID)
	}

	return nil
}
