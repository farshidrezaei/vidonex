package executor_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/farshidrezaei/vidonyx/executor"
)

func TestParseProgressStreamTable(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		totalDuration time.Duration
		wantEvents    int
		checkLast     func(t *testing.T, ev executor.ProgressEvent)
	}{
		{
			name: "full 2-step progress stream",
			input: `frame=100
fps=60.0
total_size=1000000
out_time_us=5000000
speed=2.0x
progress=continue
frame=200
fps=60.0
total_size=2000000
out_time_us=10000000
speed=2.0x
progress=end
`,
			totalDuration: 10 * time.Second,
			wantEvents:    2,
			checkLast: func(t *testing.T, ev executor.ProgressEvent) {
				if ev.Frame != 200 {
					t.Errorf("Frame = %d, want 200", ev.Frame)
				}
				if ev.FPS != 60.0 {
					t.Errorf("FPS = %f, want 60.0", ev.FPS)
				}
				if ev.Percentage != 100.0 {
					t.Errorf("Percentage = %f, want 100.0", ev.Percentage)
				}
				if ev.Progress != "end" {
					t.Errorf("Progress = %s, want end", ev.Progress)
				}
			},
		},
		{
			name: "empty input stream",
			input: `
`,
			totalDuration: 10 * time.Second,
			wantEvents:    0,
			checkLast:     nil,
		},
		{
			name: "partial malformed lines ignored gracefully",
			input: `garbage line without equals
invalid_key_only
frame=50
out_time_us=2500000
progress=end
`,
			totalDuration: 5 * time.Second,
			wantEvents:    1,
			checkLast: func(t *testing.T, ev executor.ProgressEvent) {
				if ev.Frame != 50 {
					t.Errorf("Frame = %d, want 50", ev.Frame)
				}
				if ev.Percentage != 50.0 {
					t.Errorf("Percentage = %f, want 50.0", ev.Percentage)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var events []executor.ProgressEvent
			err := executor.ParseProgressStream(strings.NewReader(tt.input), tt.totalDuration, func(ev executor.ProgressEvent) {
				events = append(events, ev)
			})

			if err != nil {
				t.Fatalf("ParseProgressStream failed: %v", err)
			}

			if len(events) != tt.wantEvents {
				t.Fatalf("Events count = %d, want %d", len(events), tt.wantEvents)
			}

			if tt.checkLast != nil && len(events) > 0 {
				tt.checkLast(t, events[len(events)-1])
			}
		})
	}
}

func TestMockExecutorTable(t *testing.T) {
	tests := []struct {
		name         string
		steps        int
		simulateErr  error
		expectedRuns int
		shouldErr    bool
	}{
		{
			name:         "successful 4 step run",
			steps:        4,
			simulateErr:  nil,
			expectedRuns: 1,
			shouldErr:    false,
		},
		{
			name:         "simulated error",
			steps:        4,
			simulateErr:  context.DeadlineExceeded,
			expectedRuns: 1,
			shouldErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := executor.NewMockExecutor()
			mock.ProgressStep = tt.steps
			mock.SimulateErr = tt.simulateErr

			ctx := context.Background()
			var receivedEvents int

			err := mock.Run(ctx, "ffmpeg", []string{"-i", "a.mp4", "b.mp4"}, 10*time.Second, func(_ executor.ProgressEvent) {
				receivedEvents++
			})

			if tt.shouldErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(mock.Commands) != tt.expectedRuns {
				t.Errorf("Commands recorded = %d, want %d", len(mock.Commands), tt.expectedRuns)
			}
			if receivedEvents != tt.steps {
				t.Errorf("Received events = %d, want %d", receivedEvents, tt.steps)
			}
		})
	}
}
