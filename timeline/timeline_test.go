package timeline_test

import (
	"testing"
	"time"

	"github.com/farshidrezaei/vidonyx/timeline"
	"github.com/farshidrezaei/vidonyx/types"
)

func TestTimeline_BuilderAndDuration(t *testing.T) {
	tl := timeline.New(
		timeline.WithCanvas(types.Res1080p),
		timeline.WithFPS(types.FPS60),
	)

	videoTrack := timeline.NewTrack("v_main", timeline.TrackKindVideo)
	clip1 := timeline.NewClip("c1", "video1.mp4", 0, 5*time.Second)
	clip2 := timeline.NewClip("c2", "video2.mp4", 5*time.Second, 10*time.Second)
	videoTrack.AddClip(clip1, clip2)

	audioTrack := timeline.NewTrack("a_bgm", timeline.TrackKindAudio)
	bgm := timeline.NewClip("bgm1", "music.mp3", 0, 20*time.Second)
	audioTrack.AddClip(bgm)

	tl.AddTrack(videoTrack, audioTrack)

	if tl.Duration() != 20*time.Second {
		t.Fatalf("Expected duration 20s, got %v", tl.Duration())
	}

	frames := tl.TotalFrames()
	if frames != 1200 { // 20s * 60 fps
		t.Fatalf("Expected 1200 frames, got %d", frames)
	}

	if err := timeline.Validate(tl); err != nil {
		t.Fatalf("Validation failed unexpectedly: %v", err)
	}
}

func TestTimeline_ValidationErrors(t *testing.T) {
	// Invalid canvas dimensions (odd numbers), negative clip duration, speed <= 0
	tl := timeline.New(
		timeline.WithCanvas(types.NewSize(1921, 1080)),
		timeline.WithFPS(types.NewRational(0, 1)), // invalid fps
	)

	badTrack := timeline.NewTrack("t1", timeline.TrackKindVideo)
	badClip := timeline.NewClip("c_bad", "", 0, -1*time.Second).WithSpeed(-2.0).WithOpacity(2.5)
	badTrack.AddClip(badClip)
	tl.AddTrack(badTrack)

	err := timeline.Validate(tl)
	if err == nil {
		t.Fatal("Expected validation to fail, but got nil")
	}

	t.Logf("Successfully captured aggregated validation errors:\n%v", err)
}
