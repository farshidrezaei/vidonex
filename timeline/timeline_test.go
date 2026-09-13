package timeline_test

import (
	"strings"
	"testing"
	"time"

	"github.com/farshidrezaei/vidonex/timeline"
	"github.com/farshidrezaei/vidonex/types"
)

func TestTimeline_ValidationTable(t *testing.T) {
	validClip := func() *timeline.Clip {
		return timeline.NewClip("c1", "source.mp4", 0, 5*time.Second)
	}

	tests := []struct {
		name          string
		buildTimeline func() *timeline.Timeline
		shouldPass    bool
		errorContains []string
	}{
		{
			name: "valid basic timeline",
			buildTimeline: func() *timeline.Timeline {
				tl := timeline.New(timeline.WithCanvas(types.Res1080p), timeline.WithFPS(types.FPS30))
				tr := timeline.NewTrack("v0", timeline.TrackKindVideo)
				tr.AddClip(validClip())
				tl.AddTrack(tr)
				return tl
			},
			shouldPass: true,
		},
		{
			name: "invalid canvas odd dimensions",
			buildTimeline: func() *timeline.Timeline {
				tl := timeline.New(timeline.WithCanvas(types.NewSize(1921, 1080)))
				tr := timeline.NewTrack("v0", timeline.TrackKindVideo)
				tr.AddClip(validClip())
				tl.AddTrack(tr)
				return tl
			},
			shouldPass:    false,
			errorContains: []string{"must be even numbers"},
		},
		{
			name: "zero fps",
			buildTimeline: func() *timeline.Timeline {
				tl := timeline.New(timeline.WithFPS(types.NewRational(0, 1)))
				tr := timeline.NewTrack("v0", timeline.TrackKindVideo)
				tr.AddClip(validClip())
				tl.AddTrack(tr)
				return tl
			},
			shouldPass:    false,
			errorContains: []string{"positive fraction"},
		},
		{
			name: "empty source in clip",
			buildTimeline: func() *timeline.Timeline {
				tl := timeline.New()
				tr := timeline.NewTrack("v0", timeline.TrackKindVideo)
				tr.AddClip(timeline.NewClip("c_empty", "", 0, 5*time.Second))
				tl.AddTrack(tr)
				return tl
			},
			shouldPass:    false,
			errorContains: []string{"empty source"},
		},
		{
			name: "negative clip duration",
			buildTimeline: func() *timeline.Timeline {
				tl := timeline.New()
				tr := timeline.NewTrack("v0", timeline.TrackKindVideo)
				tr.AddClip(timeline.NewClip("c_neg", "src.mp4", 0, -5*time.Second))
				tl.AddTrack(tr)
				return tl
			},
			shouldPass:    false,
			errorContains: []string{"duration must be positive"},
		},
		{
			name: "invalid opacity outside range",
			buildTimeline: func() *timeline.Timeline {
				tl := timeline.New()
				tr := timeline.NewTrack("v0", timeline.TrackKindVideo)
				tr.AddClip(timeline.NewClip("c_opac", "src.mp4", 0, 5*time.Second).WithOpacity(1.5))
				tl.AddTrack(tr)
				return tl
			},
			shouldPass:    false,
			errorContains: []string{"opacity must be in [0.0, 1.0]"},
		},
		{
			name: "invalid transition with missing clips",
			buildTimeline: func() *timeline.Timeline {
				tl := timeline.New()
				tr := timeline.NewTrack("v0", timeline.TrackKindVideo)
				tr.AddClip(validClip())
				tr.AddTransition(timeline.NewTransition("tr1", timeline.TransitionDissolve, 1*time.Second, nil, nil))
				tl.AddTrack(tr)
				return tl
			},
			shouldPass:    false,
			errorContains: []string{"must have both ClipA and ClipB set"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tl := tt.buildTimeline()
			err := timeline.Validate(tl)
			if tt.shouldPass {
				if err != nil {
					t.Fatalf("unexpected validation error: %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("expected validation error, but got nil")
			}

			errStr := err.Error()
			for _, expectedSnippet := range tt.errorContains {
				if !strings.Contains(errStr, expectedSnippet) {
					t.Errorf("expected error to contain %q, but got: %v", expectedSnippet, errStr)
				}
			}
		})
	}
}

func TestTimeline_DurationCalculationTable(t *testing.T) {
	tests := []struct {
		name          string
		buildTimeline func() *timeline.Timeline
		wantDuration  time.Duration
	}{
		{
			name: "single clip track",
			buildTimeline: func() *timeline.Timeline {
				tl := timeline.New()
				tr := timeline.NewTrack("t1", timeline.TrackKindVideo)
				tr.AddClip(timeline.NewClip("c1", "s.mp4", 0, 10*time.Second))
				tl.AddTrack(tr)
				return tl
			},
			wantDuration: 10 * time.Second,
		},
		{
			name: "multiple tracks longest wins",
			buildTimeline: func() *timeline.Timeline {
				tl := timeline.New()
				tr1 := timeline.NewTrack("t1", timeline.TrackKindVideo)
				tr1.AddClip(timeline.NewClip("c1", "s1.mp4", 0, 10*time.Second))

				tr2 := timeline.NewTrack("t2", timeline.TrackKindAudio)
				tr2.AddClip(timeline.NewClip("c2", "s2.mp3", 5*time.Second, 15*time.Second)) // ends at 20s

				tl.AddTrack(tr1, tr2)
				return tl
			},
			wantDuration: 20 * time.Second,
		},
		{
			name: "explicit duration overrides tracks",
			buildTimeline: func() *timeline.Timeline {
				tl := timeline.New(timeline.WithDuration(30 * time.Second))
				tr1 := timeline.NewTrack("t1", timeline.TrackKindVideo)
				tr1.AddClip(timeline.NewClip("c1", "s1.mp4", 0, 10*time.Second))
				tl.AddTrack(tr1)
				return tl
			},
			wantDuration: 30 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tl := tt.buildTimeline()
			if got := tl.Duration(); got != tt.wantDuration {
				t.Errorf("Duration() = %v, want %v", got, tt.wantDuration)
			}
		})
	}
}
