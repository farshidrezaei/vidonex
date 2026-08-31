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

func (p *ValidatePass) Name() string {
	return "ValidatePass"
}

func (p *ValidatePass) Run(g *Graph) error {
	var errs []error

	// 1. Check for cycles via Topological Sort
	if _, err := g.TopologicalSort(); err != nil {
		errs = append(errs, err)
	}

	// 2. Check stream types and connectivity
	for node := range g.Nodes() {
		for _, inPad := range node.Inputs {
			srcOut, connected := g.GetSourcePad(inPad)
			if !connected {
				// Non-connected inputs (unless external input stream like [0:v])
				if !isExternalPad(inPad.ID) {
					errs = append(errs, fmt.Errorf("node %q input pad %q is disconnected", node.ID, inPad.ID))
				}
				continue
			}
			if srcOut.StreamType != inPad.StreamType {
				errs = append(errs, fmt.Errorf("node %q input %q (%s) connected to mismatched output %q (%s)",
					node.ID, inPad.ID, inPad.StreamType, srcOut.ID, srcOut.StreamType))
			}
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
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
