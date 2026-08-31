package executor

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

// CommandExecutor abstracts the invocation of FFmpeg/FFprobe binaries.
type CommandExecutor interface {
	Run(ctx context.Context, cmdName string, args []string, totalDuration time.Duration, onProgress func(ProgressEvent)) error
}

// OSExecutor executes real OS processes with progress streaming and context cancellation.
type OSExecutor struct{}

// NewOSExecutor creates a new OSExecutor.
func NewOSExecutor() *OSExecutor {
	return &OSExecutor{}
}

// Run executes the command, pipes progress telemetry, and respects context cancellation.
func (e *OSExecutor) Run(ctx context.Context, cmdName string, args []string, totalDuration time.Duration, onProgress func(ProgressEvent)) error {
	// Prepend -progress pipe:1 to capture telemetry on stdout
	execArgs := make([]string, 0, len(args)+2)
	execArgs = append(execArgs, "-progress", "pipe:1")
	execArgs = append(execArgs, args...)

	cmd := exec.CommandContext(ctx, cmdName, execArgs...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("executor: failed to acquire stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("executor: failed to start %s: %w", cmdName, err)
	}

	// Parse progress asynchronously
	parseErrChan := make(chan error, 1)
	go func() {
		parseErrChan <- ParseProgressStream(stdout, totalDuration, onProgress)
	}()

	cmdErr := cmd.Wait()
	parseErr := <-parseErrChan

	if cmdErr != nil {
		return fmt.Errorf("executor: command %s failed: %w", cmdName, cmdErr)
	}
	if parseErr != nil {
		return fmt.Errorf("executor: progress parsing failed: %w", parseErr)
	}

	return nil
}
