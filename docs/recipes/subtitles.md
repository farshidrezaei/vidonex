# Burn-in Subtitles & Captions Recipe

Parse standard SRT or WebVTT files and burn stylized captions directly into video frames with custom fonts and highlight boxes.

```go
package main

import (
	"context"
	"time"

	"github.com/farshidrezaei/vidonex/composer"
	"github.com/farshidrezaei/vidonex/subtitles"
	"github.com/farshidrezaei/vidonex/timeline"
	"github.com/farshidrezaei/vidonex/types"
)

func main() {
	tl := timeline.New(
		timeline.WithCanvas(types.Res1080p),
		timeline.WithFPS(types.FPS60),
	)

	videoTrack := timeline.NewTrack("video", timeline.TrackKindVideo)
	videoTrack.AddClip(timeline.NewClip("clip", "assets/interview.mp4", 0, 20*time.Second))

	// Subtitle Styling
	style := subtitles.Style{
		FontName:  "Inter Bold",
		FontSize:  32,
		TextColor: types.HexColor("#FFCC00"), // Yellow social captions
		BoxColor:  types.Black,
		BoxAlpha:  0.75,
		MarginV:   60,
	}

	subTrack := timeline.NewTrack("subtitles", timeline.TrackKindSubtitle)
	subTrack.AddClip(
		timeline.NewClip("captions", "assets/dialogue.srt", 0, 20*time.Second).
			WithSubtitleStyle(style),
	)

	tl.AddTrack(videoTrack, subTrack)

	c := composer.New()
	_, _ = c.Render(context.Background(), tl, "subtitled.mp4", nil)
}
```
