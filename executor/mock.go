package executor

import (
	"context"
	"time"
)

// ExecutedCommand records an invocation captured by MockExecutor.
type ExecutedCommand struct {
	Command string
	Args    []string
}

// MockExecutor records executions and simulates progress events for testing without binaries.
type MockExecutor struct {
	Commands     []ExecutedCommand
	SimulateErr  error
	ProgressStep int // Number of simulated progress events (e.g. 4 for 25%, 50%, 75%, 100%)
}

// NewMockExecutor creates an initialized MockExecutor.
func NewMockExecutor() *MockExecutor {
	return &MockExecutor{
		Commands:     make([]ExecutedCommand, 0),
		ProgressStep: 4,
	}
}

// Run records the command and simulates progress ticks.
func (m *MockExecutor) Run(ctx context.Context, cmdName string, args []string, totalDuration time.Duration, onProgress func(ProgressEvent)) error {
	m.Commands = append(m.Commands, ExecutedCommand{
		Command: cmdName,
		Args:    append([]string(nil), args...),
	})

	if m.SimulateErr != nil {
		return m.SimulateErr
	}

	// Simulate progress ticks
	if onProgress != nil && m.ProgressStep > 0 {
		for i := 1; i <= m.ProgressStep; i++ {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				pct := (float64(i) / float64(m.ProgressStep)) * 100.0
				prog := "continue"
				if i == m.ProgressStep {
					prog = "end"
				}
				onProgress(ProgressEvent{
					Percentage: pct,
					Progress:   prog,
				})
			}
		}
	}

	return nil
}
