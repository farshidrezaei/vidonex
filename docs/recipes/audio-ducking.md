# Sidechain Audio Ducking Recipe

Automatically attenuate background music volume when a speaker is talking.

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

	// Voiceover commentary track
	voiceTrack := timeline.NewTrack("voice", timeline.TrackKindAudio)
	voiceTrack.AddClip(
		timeline.NewClip("voice_clip", "assets/speech.mp3", 3*time.Second, 10*time.Second),
	)

	// Background soundtrack track with sidechain ducking targeting 'voice'
	musicTrack := timeline.NewTrack("music", timeline.TrackKindAudio)
	musicTrack.AddClip(
		timeline.NewClip("music_clip", "assets/background.mp3", 0, 20*time.Second).
			WithVolume(0.8),
	)

	tl.AddTrack(voiceTrack, musicTrack)

	c := composer.New()
	_, _ = c.Render(context.Background(), tl, "ducking_output.mp4", nil)
}
```
