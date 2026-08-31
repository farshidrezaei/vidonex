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
func (mock *MockExecutor) Run(ctx context.Context, cmdName string, args []string, _ time.Duration, onProgress func(ProgressEvent)) error {
	mock.Commands = append(mock.Commands, ExecutedCommand{
		Command: cmdName,
		Args:    append([]string(nil), args...),
	})

	if mock.SimulateErr != nil {
		return mock.SimulateErr
	}

	// Simulate progress ticks
	if onProgress != nil && mock.ProgressStep > 0 {
		for i := 1; i <= mock.ProgressStep; i++ {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				percentage := (float64(i) / float64(mock.ProgressStep)) * 100.0
				progress := "continue"
				if i == mock.ProgressStep {
					progress = "end"
				}
				onProgress(ProgressEvent{
					Percentage: percentage,
					Progress:   progress,
				})
			}
		}
	}

	return nil
}
