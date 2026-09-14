# Green Screen & Despill Studio Recipe

Isolate green-screen presenter footage, remove green color reflections, and place them on an animated studio background.

```go
package main

import (
	"context"
	"time"

	"github.com/farshidrezaei/vidonex/chromakey"
	"github.com/farshidrezaei/vidonex/composer"
	"github.com/farshidrezaei/vidonex/timeline"
	"github.com/farshidrezaei/vidonex/types"
)

func main() {
	tl := timeline.New(
		timeline.WithCanvas(types.Res1080p),
		timeline.WithFPS(types.FPS60),
	)

	// Background Virtual Studio Track
	bgTrack := timeline.NewTrack("studio_bg", timeline.TrackKindVideo)
	bgTrack.AddClip(
		timeline.NewClip("backdrop", "assets/virtual_studio.mp4", 0, 15*time.Second),
	)

	// Green Screen Presenter Clip with ChromaKey Options
	keyOptions := chromakey.Options{
		KeyColor:   types.HexColor("#00FF00"),
		Similarity: 0.25,
		Smoothness: 0.10,
		Despill:    true,
	}

	presenterTrack := timeline.NewTrack("presenter", timeline.TrackKindOverlay).SetZIndex(1)
	presenterTrack.AddClip(
		timeline.NewClip("actor", "assets/greenscreen_actor.mp4", 0, 15*time.Second).
			WithChromaKey(keyOptions),
	)

	tl.AddTrack(bgTrack, presenterTrack)

	c := composer.New()
	_, _ = c.Render(context.Background(), tl, "studio_output.mp4", nil)
}
```
