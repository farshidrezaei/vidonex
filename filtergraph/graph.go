// Package filtergraph implements the intermediate representation Directed Acyclic Graph (DAG) for FFmpeg filter chains.
package filtergraph

import (
	"errors"
	"fmt"
	"iter"
	"strings"
)

var (
	// ErrCycleDetected indicates that adding a connection would create a cycle in the DAG.
	ErrCycleDetected = errors.New("filtergraph: cycle detected in graph")
	// ErrStreamTypeMismatch indicates an attempt to connect mismatched stream types (e.g. video to audio).
	ErrStreamTypeMismatch = errors.New("filtergraph: stream type mismatch")
	// ErrPadAlreadyConnected indicates that the target input pad is already connected to another source.
	ErrPadAlreadyConnected = errors.New("filtergraph: input pad is already connected")
	// ErrNodeNotFound indicates that the requested filter node does not exist in the graph.
	ErrNodeNotFound = errors.New("filtergraph: node not found")
)

// Graph represents a Directed Acyclic Graph (DAG) of FFmpeg filter nodes and connected pads.
type Graph struct {
	nodes        map[string]*Node
	orderedNodes []*Node
	edges        map[*Pad]*Pad   // Key: OutputPad, Value: InputPad (or primary link)
	inToOut      map[*Pad]*Pad   // Key: InputPad, Value: OutputPad (reverse mapping)
	padConsumers map[*Pad][]*Pad // Key: OutputPad, Value: all consuming InputPads
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
func (graph *Graph) AddNode(node *Node) error {
	if node == nil {
		return errors.New("filtergraph: cannot add nil node")
	}
	if _, exists := graph.nodes[node.ID]; exists {
		return fmt.Errorf("filtergraph: node with ID %q already exists", node.ID)
	}
	graph.nodes[node.ID] = node
	graph.orderedNodes = append(graph.orderedNodes, node)
	return nil
}

// NewNode creates and registers a new filter node in the graph.
func (graph *Graph) NewNode(id, filterName string) *Node {
	node := NewNode(id, filterName)
	_ = graph.AddNode(node)
	return node
}

// RemoveNode removes a node and all associated connections from the graph.
func (graph *Graph) RemoveNode(id string) error {
	node, exists := graph.nodes[id]
	if !exists {
		return ErrNodeNotFound
	}

	// Disconnect all inputs
	for _, inputPad := range node.Inputs {
		_ = graph.Disconnect(inputPad)
	}

	// Disconnect all outputs
	for _, outputPad := range node.Outputs {
		consumers := append([]*Pad(nil), graph.padConsumers[outputPad]...)
		for _, inputPad := range consumers {
			_ = graph.Disconnect(inputPad)
		}
		delete(graph.edges, outputPad)
		delete(graph.padConsumers, outputPad)
	}

	delete(graph.nodes, id)
	for index, n := range graph.orderedNodes {
		if n.ID == id {
			graph.orderedNodes = append(graph.orderedNodes[:index], graph.orderedNodes[index+1:]...)
			break
		}
	}
	return nil
}

// GetNode retrieves a node by ID.
func (graph *Graph) GetNode(id string) (*Node, bool) {
	node, ok := graph.nodes[id]
	return node, ok
}

// NextPadID generates a deterministic unique pad label (e.g. "p1", "p2").
func (graph *Graph) NextPadID(prefix string) string {
	graph.padCounter++
	if prefix == "" {
		prefix = "pad"
	}
	return fmt.Sprintf("%s_%d", prefix, graph.padCounter)
}

// Connect establishes a directed connection from an output pad to an input pad.
func (graph *Graph) Connect(sourceOutputPad, destinationInputPad *Pad) error {
	if sourceOutputPad == nil || destinationInputPad == nil {
		return errors.New("filtergraph: cannot connect nil pads")
	}
	if sourceOutputPad.IsInput || !destinationInputPad.IsInput {
		return errors.New("filtergraph: connect requires (sourceOutputPad: IsInput=false, destinationInputPad: IsInput=true)")
	}
	if sourceOutputPad.StreamType != destinationInputPad.StreamType {
		return fmt.Errorf("%w: cannot connect %s output to %s input", ErrStreamTypeMismatch, sourceOutputPad.StreamType, destinationInputPad.StreamType)
	}
	if existing, connected := graph.inToOut[destinationInputPad]; connected && existing != sourceOutputPad {
		return fmt.Errorf("%w: input %s is already connected to %s", ErrPadAlreadyConnected, destinationInputPad.ID, existing.ID)
	}

	// Check for cycles if both nodes belong to the graph
	if sourceOutputPad.Node != nil && destinationInputPad.Node != nil {
		if graph.wouldCreateCycle(sourceOutputPad.Node, destinationInputPad.Node) {
			return ErrCycleDetected
		}
	}

	graph.inToOut[destinationInputPad] = sourceOutputPad
	graph.edges[sourceOutputPad] = destinationInputPad
	graph.padConsumers[sourceOutputPad] = append(graph.padConsumers[sourceOutputPad], destinationInputPad)

	// Ensure pads share the same label ID for valid FFmpeg filtergraph linking
	if destinationInputPad.ID == "" || destinationInputPad.ID != sourceOutputPad.ID {
		destinationInputPad.ID = sourceOutputPad.ID
	}

	return nil
}

// Disconnect removes the connection feeding into destinationInputPad.
func (graph *Graph) Disconnect(destinationInputPad *Pad) error {
	sourceOutputPad, connected := graph.inToOut[destinationInputPad]
	if !connected {
		return nil
	}

	delete(graph.inToOut, destinationInputPad)
	consumers := graph.padConsumers[sourceOutputPad]
	for index, inputPad := range consumers {
		if inputPad == destinationInputPad {
			graph.padConsumers[sourceOutputPad] = append(consumers[:index], consumers[index+1:]...)
			break
		}
	}
	if len(graph.padConsumers[sourceOutputPad]) == 0 {
		delete(graph.edges, sourceOutputPad)
		delete(graph.padConsumers, sourceOutputPad)
	}
	return nil
}

// GetSourcePad returns the output pad feeding into destinationInputPad.
func (graph *Graph) GetSourcePad(destinationInputPad *Pad) (*Pad, bool) {
	sourceOutputPad, ok := graph.inToOut[destinationInputPad]
	return sourceOutputPad, ok
}

// GetConsumerPads returns all input pads consuming sourceOutputPad.
func (graph *Graph) GetConsumerPads(sourceOutputPad *Pad) []*Pad {
	return graph.padConsumers[sourceOutputPad]
}

// Nodes returns an iterator over all nodes in registration order.
func (graph *Graph) Nodes() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		for _, node := range graph.orderedNodes {
			if !yield(node) {
				return
			}
		}
	}
}

// NodeCount returns the number of nodes in the graph.
func (graph *Graph) NodeCount() int {
	return len(graph.nodes)
}

// TopologicalSort computes an ordered sequence of nodes respecting dependencies using Kahn's algorithm.
func (graph *Graph) TopologicalSort() ([]*Node, error) {
	inDegree := make(map[string]int)
	for _, node := range graph.nodes {
		inDegree[node.ID] = 0
	}

	// Calculate in-degrees (number of predecessor nodes)
	for _, node := range graph.nodes {
		for _, inputPad := range node.Inputs {
			if sourceOutputPad, ok := graph.inToOut[inputPad]; ok && sourceOutputPad.Node != nil {
				inDegree[node.ID]++
			}
		}
	}

	// Queue nodes with zero in-degree
	queue := make([]*Node, 0)
	for _, node := range graph.orderedNodes {
		if inDegree[node.ID] == 0 {
			queue = append(queue, node)
		}
	}

	result := make([]*Node, 0, len(graph.nodes))

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)

		// For each successor node connected to current's outputs
		for _, outputPad := range current.Outputs {
			for _, inputPad := range graph.padConsumers[outputPad] {
				if inputPad.Node == nil {
					continue
				}
				successorID := inputPad.Node.ID
				inDegree[successorID]--
				if inDegree[successorID] == 0 {
					queue = append(queue, inputPad.Node)
				}
			}
		}
	}

	if len(result) != len(graph.nodes) {
		return nil, ErrCycleDetected
	}

	return result, nil
}

// FormattedFilterComplex renders the full filtergraph as a semicolon-separated FFmpeg -filter_complex string.
func (graph *Graph) FormattedFilterComplex() (string, error) {
	sortedNodes, err := graph.TopologicalSort()
	if err != nil {
		return "", err
	}

	var stringBuilder strings.Builder
	for index, node := range sortedNodes {
		if index > 0 {
			stringBuilder.WriteString(";\n")
		}
		stringBuilder.WriteString(node.String())
	}
	return stringBuilder.String(), nil
}

// wouldCreateCycle checks if adding an edge from source to destination would introduce a cycle.
func (graph *Graph) wouldCreateCycle(sourceNode, destinationNode *Node) bool {
	if sourceNode.ID == destinationNode.ID {
		return true
	}
	visited := make(map[string]bool)
	var depthFirstSearch func(current *Node) bool
	depthFirstSearch = func(current *Node) bool {
		if current.ID == sourceNode.ID {
			return true
		}
		visited[current.ID] = true
		for _, outputPad := range current.Outputs {
			for _, inputPad := range graph.padConsumers[outputPad] {
				if inputPad.Node != nil && !visited[inputPad.Node.ID] {
					if depthFirstSearch(inputPad.Node) {
						return true
					}
				}
			}
		}
		return false
	}
	return depthFirstSearch(destinationNode)
}
