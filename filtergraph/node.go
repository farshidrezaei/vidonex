package filtergraph

import (
	"fmt"
	"sort"
	"strings"
)

// FilterParam represents a key-value parameter passed to an FFmpeg filter.
type FilterParam struct {
	Key   string
	Value any
}

// Node represents a discrete filter operation within the FFmpeg filtergraph DAG.
type Node struct {
	ID         string
	FilterName string
	Params     []FilterParam
	Inputs     []*Pad
	Outputs    []*Pad
	Comment    string
}

// NewNode creates an empty Node.
func NewNode(id, filterName string) *Node {
	return &Node{
		ID:         id,
		FilterName: filterName,
		Params:     make([]FilterParam, 0),
		Inputs:     make([]*Pad, 0),
		Outputs:    make([]*Pad, 0),
	}
}

// SetParam sets or updates a key-value parameter on the node.
func (n *Node) SetParam(key string, value any) *Node {
	for i, p := range n.Params {
		if p.Key == key {
			n.Params[i].Value = value
			return n
		}
	}
	n.Params = append(n.Params, FilterParam{Key: key, Value: value})
	return n
}

// AddInput appends a new input pad to the node.
func (n *Node) AddInput(id string, st StreamType) *Pad {
	pad := &Pad{
		ID:         id,
		StreamType: st,
		Node:       n,
		IsInput:    true,
	}
	n.Inputs = append(n.Inputs, pad)
	return pad
}

// AddOutput appends a new output pad to the node.
func (n *Node) AddOutput(id string, st StreamType) *Pad {
	pad := &Pad{
		ID:         id,
		StreamType: st,
		Node:       n,
		IsInput:    false,
	}
	n.Outputs = append(n.Outputs, pad)
	return pad
}

// IsSource reports whether this node produces streams without inputs (e.g. color, anullsrc, or input media).
func (n *Node) IsSource() bool {
	return len(n.Inputs) == 0
}

// IsSink reports whether this node is a terminal consumer.
func (n *Node) IsSink() bool {
	return len(n.Outputs) == 0
}

// FormattedFilter formats the filter name and arguments (e.g., "scale=1920:1080:flags=lanczos").
func (n *Node) FormattedFilter() string {
	if len(n.Params) == 0 {
		return n.FilterName
	}

	var sb strings.Builder
	sb.WriteString(n.FilterName)
	sb.WriteString("=")

	for i, p := range n.Params {
		if i > 0 {
			sb.WriteString(":")
		}
		if p.Key == "" {
			sb.WriteString(fmt.Sprintf("%v", p.Value))
		} else {
			sb.WriteString(fmt.Sprintf("%s=%v", p.Key, p.Value))
		}
	}
	return sb.String()
}

// String returns a single filter expression with input and output pad labels
// (e.g. "[0:v][1:v]overlay=x=100:y=200[v_out]").
func (n *Node) String() string {
	var sb strings.Builder

	// Inputs
	for _, in := range n.Inputs {
		sb.WriteString(in.Label())
	}

	// Filter
	sb.WriteString(n.FormattedFilter())

	// Outputs
	for _, out := range n.Outputs {
		sb.WriteString(out.Label())
	}

	return sb.String()
}

// ParamsMap returns a copy of parameters as a map.
func (n *Node) ParamsMap() map[string]any {
	m := make(map[string]any, len(n.Params))
	for _, p := range n.Params {
		m[p.Key] = p.Value
	}
	return m
}

// SortParams ensures deterministic parameter ordering for unit testing and reproducibility.
func (n *Node) SortParams() {
	sort.Slice(n.Params, func(i, j int) bool {
		return n.Params[i].Key < n.Params[j].Key
	})
}
