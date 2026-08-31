package cli_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/farshidrezaei/vidonyx/executor"
	"github.com/farshidrezaei/vidonyx/internal/cli"
)

func TestCLI_UITable(t *testing.T) {
	tests := []struct {
		name           string
		action         func(ui *cli.UI)
		expectedOutput string
	}{
		{
			name: "print_banner",
			action: func(ui *cli.UI) {
				ui.PrintBanner()
			},
			expectedOutput: "Declarative Video Composition & FFmpeg Filtergraph Engine in Go",
		},
		{
			name: "print_step",
			action: func(ui *cli.UI) {
				ui.PrintStep(1, 4, "🔍", "Parsing Project Specification")
			},
			expectedOutput: "[1/4] 🔍 Parsing Project Specification",
		},
		{
			name: "print_success",
			action: func(ui *cli.UI) {
				ui.PrintSuccess("Timeline validated successfully")
			},
			expectedOutput: "Timeline validated successfully",
		},
		{
			name: "print_error",
			action: func(ui *cli.UI) {
				ui.PrintError(errors.New("file not found"))
			},
			expectedOutput: "file not found",
		},
		{
			name: "print_summary_card",
			action: func(ui *cli.UI) {
				ui.PrintSummaryCard("output.mp4", 10485760, 2*time.Second, 10*time.Second)
			},
			expectedOutput: "Composition Rendered Successfully!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var outputBuffer bytes.Buffer
			ui := cli.NewUI(&outputBuffer)

			tt.action(ui)

			result := outputBuffer.String()
			if !strings.Contains(result, tt.expectedOutput) {
				t.Errorf("UI output did not contain %q. Got:\n%s", tt.expectedOutput, result)
			}
		})
	}
}

func TestCLI_ProgressBarTable(t *testing.T) {
	var buffer bytes.Buffer
	progressBar := cli.NewProgressBar(&buffer, 10*time.Second)

	// In non-TTY buffer, milestone percentages are written
	progressBar.Update(executor.ProgressEvent{
		Percentage: 15.0,
		Speed:      2.5,
		FPS:        60.0,
	})
	progressBar.Update(executor.ProgressEvent{
		Percentage: 35.0,
		Speed:      2.4,
		FPS:        59.5,
	})
	progressBar.Update(executor.ProgressEvent{
		Percentage: 100.0,
		Speed:      2.6,
		FPS:        60.0,
	})
	progressBar.Finish()

	output := buffer.String()
	if !strings.Contains(output, "Rendering:  10%") {
		t.Errorf("Progress output missing 10%% milestone: %s", output)
	}
	if !strings.Contains(output, "Rendering:  30%") {
		t.Errorf("Progress output missing 30%% milestone: %s", output)
	}
	if !strings.Contains(output, "Rendering: 100%") {
		t.Errorf("Progress output missing 100%% milestone: %s", output)
	}
}

func TestCLI_FormatHelpersTable(t *testing.T) {
	t.Run("format_file_size", func(t *testing.T) {
		tests := []struct {
			bytes    int64
			expected string
		}{
			{bytes: 500, expected: "500 B"},
			{bytes: 1024, expected: "1.00 KB"},
			{bytes: 1048576, expected: "1.00 MB"},
			{bytes: 1073741824, expected: "1.00 GB"},
		}

		for _, tt := range tests {
			res := cli.FormatFileSize(tt.bytes)
			if res != tt.expected {
				t.Errorf("FormatFileSize(%d) = %q, want %q", tt.bytes, res, tt.expected)
			}
		}
	})

	t.Run("format_duration", func(t *testing.T) {
		tests := []struct {
			dur      time.Duration
			expected string
		}{
			{dur: 5 * time.Second, expected: "00:05"},
			{dur: 65 * time.Second, expected: "01:05"},
			{dur: 3665 * time.Second, expected: "61:05"},
		}

		for _, tt := range tests {
			res := cli.FormatDuration(tt.dur)
			if res != tt.expected {
				t.Errorf("FormatDuration(%v) = %q, want %q", tt.dur, res, tt.expected)
			}
		}
	})
}

func TestCLI_LoggerSetupTable(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error"}
	formats := []string{"text", "json"}

	for _, level := range levels {
		for _, format := range formats {
			var logBuffer bytes.Buffer
			logger := cli.SetupLogger(level, format, &logBuffer)
			if logger == nil {
				t.Fatalf("SetupLogger(%q, %q) returned nil", level, format)
			}
			logger.Error("test error message", "key", "value")
			if logBuffer.Len() == 0 {
				t.Errorf("No logs written for level=%s format=%s", level, format)
			}
		}
	}
}
