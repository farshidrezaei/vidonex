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

---

## Declarative YAML & CLI Automation

You can also define audiograms declaratively and render directly via the `vidonex` CLI:

```yaml
version: "1.0"
canvas:
  preset: "square_1080"
tracks:
  - id: "voice_track"
    kind: "audio"
    clips:
      - id: "podcast_audio"
        source: "assets/podcast.mp3"
        duration: 30s

  - id: "audiogram_track"
    kind: "waveform"
    z_index: 2
    waveform:
      source_audio: "voice_track"
      mode: "p2p" # p2p, line, cline, dot, bars, spectrum, wave, circular
      color: "#00FFCC"
      secondary_color: "#6366F1"
      scale: "log"
      density: 50
      glow: true
```

Render with CLI:
```bash
vidonex render audiogram.yaml -o audiogram.mp4
```
