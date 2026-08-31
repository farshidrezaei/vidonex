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
func (node *Node) SetParam(key string, value any) *Node {
	for index, param := range node.Params {
		if param.Key == key {
			node.Params[index].Value = value
			return node
		}
	}
	node.Params = append(node.Params, FilterParam{Key: key, Value: value})
	return node
}

// AddInput appends a new input pad to the node.
func (node *Node) AddInput(id string, streamType StreamType) *Pad {
	pad := &Pad{
		ID:         id,
		StreamType: streamType,
		Node:       node,
		IsInput:    true,
	}
	node.Inputs = append(node.Inputs, pad)
	return pad
}

// AddOutput appends a new output pad to the node.
func (node *Node) AddOutput(id string, streamType StreamType) *Pad {
	pad := &Pad{
		ID:         id,
		StreamType: streamType,
		Node:       node,
		IsInput:    false,
	}
	node.Outputs = append(node.Outputs, pad)
	return pad
}

// IsSource reports whether this node produces streams without inputs (e.g. color, anullsrc, or input media).
func (node *Node) IsSource() bool {
	return len(node.Inputs) == 0
}

// IsSink reports whether this node is a terminal consumer.
func (node *Node) IsSink() bool {
	return len(node.Outputs) == 0
}

// FormattedFilter formats the filter name and arguments (e.g., "scale=1920:1080:flags=lanczos").
func (node *Node) FormattedFilter() string {
	if len(node.Params) == 0 {
		return node.FilterName
	}

	var stringBuilder strings.Builder
	stringBuilder.WriteString(node.FilterName)
	stringBuilder.WriteString("=")

	for index, param := range node.Params {
		if index > 0 {
			stringBuilder.WriteString(":")
		}
		if param.Key == "" {
			fmt.Fprintf(&stringBuilder, "%v", param.Value)
		} else {
			fmt.Fprintf(&stringBuilder, "%s=%v", param.Key, param.Value)
		}
	}
	return stringBuilder.String()
}

// String returns a single filter expression with input and output pad labels
// (e.g. "[0:v][1:v]overlay=x=100:y=200[v_out]").
func (node *Node) String() string {
	var stringBuilder strings.Builder

	// Inputs
	for _, inputPad := range node.Inputs {
		stringBuilder.WriteString(inputPad.Label())
	}

	// Filter
	stringBuilder.WriteString(node.FormattedFilter())

	// Outputs
	for _, outputPad := range node.Outputs {
		stringBuilder.WriteString(outputPad.Label())
	}

	return stringBuilder.String()
}

// ParamsMap returns a copy of parameters as a map.
func (node *Node) ParamsMap() map[string]any {
	paramMap := make(map[string]any, len(node.Params))
	for _, param := range node.Params {
		paramMap[param.Key] = param.Value
	}
	return paramMap
}

// SortParams ensures deterministic parameter ordering for unit testing and reproducibility.
func (node *Node) SortParams() {
	sort.Slice(node.Params, func(i, j int) bool {
		return node.Params[i].Key < node.Params[j].Key
	})
}
