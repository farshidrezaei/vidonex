# Picture-in-Picture (PiP) Recipe

Learn how to overlay a secondary camera feed (e.g. facecam) onto a primary video stream.

```go
package main

import (
	"context"
	"time"

	"github.com/farshidrezaei/vidonex/composer"
	"github.com/farshidrezaei/vidonex/timeline"
	"github.com/farshidrezaei/vidonex/types"
)

func main() {
	tl := timeline.New(
		timeline.WithCanvas(types.Res1080p),
		timeline.WithFPS(types.FPS60),
	)

	// Main Gameplay Track (Layer 0)
	mainTrack := timeline.NewTrack("gameplay", timeline.TrackKindVideo).SetZIndex(0)
	mainTrack.AddClip(
		timeline.NewClip("gameplay_clip", "assets/gameplay.mp4", 0, 15*time.Second),
	)

	// Facecam Overlay Track (Layer 1)
	pipTrack := timeline.NewTrack("facecam", timeline.TrackKindOverlay).SetZIndex(1)
	pipTrack.AddClip(
		timeline.NewClip("facecam_clip", "assets/facecam.mp4", 2*time.Second, 10*time.Second).
			WithScale(0.25).
			WithPosition(types.Point{X: 1400, Y: 40}).
			WithOpacity(0.95),
	)

	tl.AddTrack(mainTrack, pipTrack)

	c := composer.New()
	_, _ = c.Render(context.Background(), tl, "pip_output.mp4", nil)
}
```
