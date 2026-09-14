# Animated Audiogram Waveforms Recipe

Create dynamic, pulsing neon audiograms from voice recordings for podcasts, Instagram Reels, and YouTube Shorts.

```go
package main

import (
	"context"
	"time"

	"github.com/farshidrezaei/vidonex/composer"
	"github.com/farshidrezaei/vidonex/timeline"
	"github.com/farshidrezaei/vidonex/types"
	"github.com/farshidrezaei/vidonex/waveform"
)

func main() {
	// Square 1:1 format for social feeds
	tl := timeline.New(
		timeline.WithCanvas(types.Size{Width: 1080, Height: 1080}),
		timeline.WithFPS(types.FPS60),
	)

	// Waveform configuration
	options := waveform.Options{
		Mode:       waveform.ModePeakToPeak,
		Color:      types.HexColor("#00FFCC"), // Vibrant Cyan Neon
		Scale:      waveform.ScaleLogarithmic,
		Width:      900,
		Height:     240,
	}

	waveTrack := timeline.NewTrack("waveform_fx", timeline.TrackKindWaveform)
	waveTrack.AddClip(
		timeline.NewClip("voice", "assets/podcast.mp3", 0, 30*time.Second).
			WithPosition(types.Point{X: 90, Y: 420}),
	)

	tl.AddTrack(waveTrack)

	c := composer.New()
	_, _ = c.Render(context.Background(), tl, "audiogram.mp4", nil)
}
```
