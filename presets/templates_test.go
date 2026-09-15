package presets_test

import (
	"testing"
	"time"

	"github.com/farshidrezaei/vidonex/presets"
	"github.com/farshidrezaei/vidonex/types"
)

func TestTemplates_Table(t *testing.T) {
	t.Run("TikTokSplitScreen", func(t *testing.T) {
		testCases := []struct {
			name          string
			topSource     string
			bottomSource  string
			duration      time.Duration
			expectedWidth int
			expectedError bool
		}{
			{
				name:          "Valid split screen setup",
				topSource:     "gameplay.mp4",
				bottomSource:  "facecam.mp4",
				duration:      15 * time.Second,
				expectedWidth: 1080,
				expectedError: false,
			},
			{
				name:          "Missing source path",
				topSource:     "",
				bottomSource:  "facecam.mp4",
				duration:      15 * time.Second,
				expectedWidth: 0,
				expectedError: true,
			},
			{
				name:          "Invalid zero duration",
				topSource:     "gameplay.mp4",
				bottomSource:  "facecam.mp4",
				duration:      0,
				expectedWidth: 0,
				expectedError: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tl, err := presets.TemplateTikTokSplitScreen(tc.topSource, tc.bottomSource, tc.duration)
				if tc.expectedError {
					if err == nil {
						t.Fatalf("expected error, got nil")
					}
					return
				}
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if tl.Canvas.Width != tc.expectedWidth {
					t.Errorf("expected width %d, got %d", tc.expectedWidth, tl.Canvas.Width)
				}
				if len(tl.Tracks) != 2 {
					t.Errorf("expected 2 tracks, got %d", len(tl.Tracks))
				}
			})
		}
	})

	t.Run("PodcastAudiogram", func(t *testing.T) {
		tl, err := presets.TemplatePodcastAudiogram("voice.mp3", "cover.png", 30*time.Second, "Episode 1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tl.Canvas != types.ResSquare1080 {
			t.Errorf("expected 1080x1080 canvas, got %v", tl.Canvas)
		}
		if len(tl.Tracks) != 3 {
			t.Errorf("expected 3 tracks (background, waveform, audio), got %d", len(tl.Tracks))
		}
	})

	t.Run("YouTubeEndScreen", func(t *testing.T) {
		tl, err := presets.TemplateYouTubeEndScreen("outro.mp4", 20*time.Second, "Vidonex")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tl.Canvas != types.Res1080p {
			t.Errorf("expected 1080p canvas, got %v", tl.Canvas)
		}
	})
}
