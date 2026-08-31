package filtergraph

import (
	"errors"
	"fmt"
)

// ValidatePass checks graph integrity:
// 1. Detects cycles (acyclicity).
// 2. Verifies stream type compatibility on all links.
// 3. Verifies no dangling unassigned internal inputs.
type ValidatePass struct{}

// Name returns the unique pass name.
func (pass *ValidatePass) Name() string {
	return "ValidatePass"
}

// Run verifies graph integrity and type correctness.
func (pass *ValidatePass) Run(graph *Graph) error {
	var collectedErrors []error

	// 1. Check for cycles via Topological Sort
	if _, err := graph.TopologicalSort(); err != nil {
		collectedErrors = append(collectedErrors, err)
	}

	// 2. Check stream types and connectivity
	for node := range graph.Nodes() {
		for _, inputPad := range node.Inputs {
			sourceOutputPad, connected := graph.GetSourcePad(inputPad)
			if !connected {
				// Non-connected inputs (unless external input stream like [0:v])
				if !isExternalPad(inputPad.ID) {
					collectedErrors = append(collectedErrors, fmt.Errorf("node %q input pad %q is disconnected", node.ID, inputPad.ID))
				}
				continue
			}
			if sourceOutputPad.StreamType != inputPad.StreamType {
				collectedErrors = append(collectedErrors, fmt.Errorf("node %q input %q (%s) connected to mismatched output %q (%s)",
					node.ID, inputPad.ID, inputPad.StreamType, sourceOutputPad.ID, sourceOutputPad.StreamType))
			}
		}
	}

	if len(collectedErrors) > 0 {
		return errors.Join(collectedErrors...)
	}
	return nil
}

func isExternalPad(id string) bool {
	// Matches FFmpeg input stream syntax like "0:v", "1:a", "0:v:0"
	if len(id) >= 3 && id[1] == ':' {
		return true
	}
	return false
}
