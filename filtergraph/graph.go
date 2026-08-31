package filtergraph

import (
	"errors"
	"fmt"
	"iter"
	"strings"
)

var (
	ErrCycleDetected       = errors.New("filtergraph: cycle detected in graph")
	ErrStreamTypeMismatch  = errors.New("filtergraph: stream type mismatch")
	ErrPadAlreadyConnected = errors.New("filtergraph: input pad is already connected")
	ErrNodeNotFound        = errors.New("filtergraph: node not found")
)

// Graph represents a Directed Acyclic Graph (DAG) of FFmpeg filter nodes and connected pads.
type Graph struct {
	nodes        map[string]*Node
	orderedNodes []*Node
	edges        map[*Pad]*Pad // Key: OutPad, Value: InPad (or primary link)
	inToOut      map[*Pad]*Pad // Key: InPad, Value: OutPad (reverse mapping)
	padConsumers map[*Pad][]*Pad // Key: OutPad, Value: all consuming InPads
	padCounter   int
}

// NewGraph initializes an empty filtergraph DAG.
func NewGraph() *Graph {
	return &Graph{
		nodes:        make(map[string]*Node),
		orderedNodes: make([]*Node, 0),
		edges:        make(map[*Pad]*Pad),
		inToOut:      make(map[*Pad]*Pad),
		padConsumers: make(map[*Pad][]*Pad),
	}
}

// AddNode adds an existing node to the graph.
func (g *Graph) AddNode(node *Node) error {
	if node == nil {
		return errors.New("filtergraph: cannot add nil node")
	}
	if _, exists := g.nodes[node.ID]; exists {
		return fmt.Errorf("filtergraph: node with ID %q already exists", node.ID)
	}
	g.nodes[node.ID] = node
	g.orderedNodes = append(g.orderedNodes, node)
	return nil
}

// NewNode creates and registers a new filter node in the graph.
func (g *Graph) NewNode(id, filterName string) *Node {
	node := NewNode(id, filterName)
	_ = g.AddNode(node)
	return node
}

// RemoveNode removes a node and all associated connections from the graph.
func (g *Graph) RemoveNode(id string) error {
	node, exists := g.nodes[id]
	if !exists {
		return ErrNodeNotFound
	}

	// Disconnect all inputs
	for _, inPad := range node.Inputs {
		_ = g.Disconnect(inPad)
	}

	// Disconnect all outputs
	for _, outPad := range node.Outputs {
		consumers := append([]*Pad(nil), g.padConsumers[outPad]...)
		for _, inPad := range consumers {
			_ = g.Disconnect(inPad)
		}
		delete(g.edges, outPad)
		delete(g.padConsumers, outPad)
	}

	delete(g.nodes, id)
	for i, n := range g.orderedNodes {
		if n.ID == id {
			g.orderedNodes = append(g.orderedNodes[:i], g.orderedNodes[i+1:]...)
			break
		}
	}
	return nil
}

// GetNode retrieves a node by ID.
func (g *Graph) GetNode(id string) (*Node, bool) {
	node, ok := g.nodes[id]
	return node, ok
}

// NextPadID generates a deterministic unique pad label (e.g. "p1", "p2").
func (g *Graph) NextPadID(prefix string) string {
	g.padCounter++
	if prefix == "" {
		prefix = "p"
	}
	return fmt.Sprintf("%s_%d", prefix, g.padCounter)
}

// Connect establishes a directed connection from an output pad to an input pad.
func (g *Graph) Connect(srcOut, dstIn *Pad) error {
	if srcOut == nil || dstIn == nil {
		return errors.New("filtergraph: cannot connect nil pads")
	}
	if srcOut.IsInput || !dstIn.IsInput {
		return errors.New("filtergraph: connect requires (srcOut: IsInput=false, dstIn: IsInput=true)")
	}
	if srcOut.StreamType != dstIn.StreamType {
		return fmt.Errorf("%w: cannot connect %s output to %s input", ErrStreamTypeMismatch, srcOut.StreamType, dstIn.StreamType)
	}
	if existing, connected := g.inToOut[dstIn]; connected && existing != srcOut {
		return fmt.Errorf("%w: input %s is already connected to %s", ErrPadAlreadyConnected, dstIn.ID, existing.ID)
	}

	// Check for cycles if both nodes belong to the graph
	if srcOut.Node != nil && dstIn.Node != nil {
		if g.wouldCreateCycle(srcOut.Node, dstIn.Node) {
			return ErrCycleDetected
		}
	}

	g.inToOut[dstIn] = srcOut
	g.edges[srcOut] = dstIn
	g.padConsumers[srcOut] = append(g.padConsumers[srcOut], dstIn)

	// Ensure pads share the same label ID for valid FFmpeg filtergraph linking
	if dstIn.ID == "" || dstIn.ID != srcOut.ID {
		dstIn.ID = srcOut.ID
	}

	return nil
}

// Disconnect removes the connection feeding into dstIn.
func (g *Graph) Disconnect(dstIn *Pad) error {
	srcOut, connected := g.inToOut[dstIn]
	if !connected {
		return nil
	}

	delete(g.inToOut, dstIn)
	consumers := g.padConsumers[srcOut]
	for i, inPad := range consumers {
		if inPad == dstIn {
			g.padConsumers[srcOut] = append(consumers[:i], consumers[i+1:]...)
			break
		}
	}
	if len(g.padConsumers[srcOut]) == 0 {
		delete(g.edges, srcOut)
		delete(g.padConsumers, srcOut)
	}
	return nil
}

// GetSourcePad returns the output pad feeding into dstIn.
func (g *Graph) GetSourcePad(dstIn *Pad) (*Pad, bool) {
	srcOut, ok := g.inToOut[dstIn]
	return srcOut, ok
}

// GetConsumerPads returns all input pads consuming srcOut.
func (g *Graph) GetConsumerPads(srcOut *Pad) []*Pad {
	return g.padConsumers[srcOut]
}

// Nodes returns an iterator over all nodes in registration order.
func (g *Graph) Nodes() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		for _, n := range g.orderedNodes {
			if !yield(n) {
				return
			}
		}
	}
}

// NodeCount returns the number of nodes in the graph.
func (g *Graph) NodeCount() int {
	return len(g.nodes)
}

// TopologicalSort computes an ordered sequence of nodes respecting dependencies using Kahn's algorithm.
func (g *Graph) TopologicalSort() ([]*Node, error) {
	inDegree := make(map[string]int)
	for _, node := range g.nodes {
		inDegree[node.ID] = 0
	}

	// Calculate in-degrees (number of predecessor nodes)
	for _, node := range g.nodes {
		for _, inPad := range node.Inputs {
			if srcOut, ok := g.inToOut[inPad]; ok && srcOut.Node != nil {
				inDegree[node.ID]++
			}
		}
	}

	// Queue nodes with zero in-degree
	queue := make([]*Node, 0)
	for _, node := range g.orderedNodes {
		if inDegree[node.ID] == 0 {
			queue = append(queue, node)
		}
	}

	result := make([]*Node, 0, len(g.nodes))

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		result = append(result, curr)

		// For each successor node connected to curr's outputs
		for _, outPad := range curr.Outputs {
			for _, inPad := range g.padConsumers[outPad] {
				if inPad.Node == nil {
					continue
				}
				succID := inPad.Node.ID
				inDegree[succID]--
				if inDegree[succID] == 0 {
					queue = append(queue, inPad.Node)
				}
			}
		}
	}

	if len(result) != len(g.nodes) {
		return nil, ErrCycleDetected
	}

	return result, nil
}

// FormattedFilterComplex renders the full filtergraph as a semicolon-separated FFmpeg -filter_complex string.
func (g *Graph) FormattedFilterComplex() (string, error) {
	sortedNodes, err := g.TopologicalSort()
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	for i, node := range sortedNodes {
		if i > 0 {
			sb.WriteString(";\n")
		}
		sb.WriteString(node.String())
	}
	return sb.String(), nil
}

// wouldCreateCycle checks if adding an edge from src to dst would introduce a cycle.
func (g *Graph) wouldCreateCycle(src, dst *Node) bool {
	if src.ID == dst.ID {
		return true
	}
	visited := make(map[string]bool)
	var dfs func(curr *Node) bool
	dfs = func(curr *Node) bool {
		if curr.ID == src.ID {
			return true
		}
		visited[curr.ID] = true
		for _, outPad := range curr.Outputs {
			for _, inPad := range g.padConsumers[outPad] {
				if inPad.Node != nil && !visited[inPad.Node.ID] {
					if dfs(inPad.Node) {
						return true
					}
				}
			}
		}
		return false
	}
	return dfs(dst)
}
