package composer_test

import (
	"context"
	"testing"
	"time"

	"github.com/farshidrezaei/vidonex/composer"
	"github.com/farshidrezaei/vidonex/executor"
	"github.com/farshidrezaei/vidonex/timeline"
	"github.com/farshidrezaei/vidonex/types"
)

func TestComposer_EndToEndPipelineWithMock(t *testing.T) {
	mockExec := executor.NewMockExecutor()

	c := composer.New(
		composer.WithExecutor(mockExec),
	)

	tl := timeline.New(
		timeline.WithCanvas(types.Res1080p),
		timeline.WithFPS(types.FPS30),
	)

	vTrack := timeline.NewTrack("v0", timeline.TrackKindVideo)
	clip := timeline.NewClip("c0", "intro.mp4", 0, 5*time.Second)
	vTrack.AddClip(clip)
	tl.AddTrack(vTrack)

	ctx := context.Background()
	var progressTicks int

	result, err := c.Render(ctx, tl, "final.mp4", func(_ executor.ProgressEvent) {
		progressTicks++
	})

	if err != nil {
		t.Fatalf("Composer Render failed: %v", err)
	}

	if result.OutputPath != "final.mp4" {
		t.Fatalf("Expected output path final.mp4, got %s", result.OutputPath)
	}

	if len(mockExec.Commands) != 1 {
		t.Fatalf("Expected 1 executed command, got %d", len(mockExec.Commands))
	}

	if progressTicks != 4 {
		t.Fatalf("Expected 4 progress updates, got %d", progressTicks)
	}
}
