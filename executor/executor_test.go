package executor_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/farshidrezaei/vidonyx/executor"
)

func TestParseProgressStream(t *testing.T) {
	sampleOutput := `
frame=15
fps=30.00
stream_0_0_q=28.0
bitrate=   1250.4kbits/s
total_size=1245000
out_time_us=5000000
out_time_ms=5000000
out_time=00:00:05.000000
dup_frames=0
drop_frames=0
speed=1.98x
progress=continue
frame=30
fps=30.00
stream_0_0_q=28.0
bitrate=   1250.4kbits/s
total_size=2490000
out_time_us=10000000
out_time_ms=10000000
out_time=00:00:10.000000
dup_frames=0
drop_frames=0
speed=2.01x
progress=end
`

	var events []executor.ProgressEvent
	err := executor.ParseProgressStream(strings.NewReader(sampleOutput), 10*time.Second, func(ev executor.ProgressEvent) {
		events = append(events, ev)
	})

	if err != nil {
		t.Fatalf("ParseProgressStream failed: %v", err)
	}

	if len(events) != 2 {
		t.Fatalf("Expected 2 events, got %d", len(events))
	}

	if events[0].Percentage != 50.0 {
		t.Fatalf("Event 0 Percentage: expected 50.0, got %f", events[0].Percentage)
	}
	if events[1].Percentage != 100.0 || events[1].Progress != "end" {
		t.Fatalf("Event 1: expected 100%% and end, got %f and %s", events[1].Percentage, events[1].Progress)
	}
}

func TestMockExecutor(t *testing.T) {
	mock := executor.NewMockExecutor()
	ctx := context.Background()

	var pcts []float64
	err := mock.Run(ctx, "ffmpeg", []string{"-i", "input.mp4", "output.mp4"}, 10*time.Second, func(ev executor.ProgressEvent) {
		pcts = append(pcts, ev.Percentage)
	})

	if err != nil {
		t.Fatalf("Mock Run failed: %v", err)
	}

	if len(mock.Commands) != 1 {
		t.Fatalf("Expected 1 captured command, got %d", len(mock.Commands))
	}

	if len(pcts) != 4 {
		t.Fatalf("Expected 4 progress ticks, got %d: %v", len(pcts), pcts)
	}
}
