# Timeline Package (`timeline`)

The `timeline` package provides the declarative Abstract Syntax Tree (AST) representing video compositions.

## Timeline

```go
import "github.com/farshidrezaei/vidonex/timeline"

tl := timeline.New(
    timeline.WithCanvas(types.Size{Width: 1920, Height: 1080}),
    timeline.WithFPS(types.FPS60),
    timeline.WithBackgroundColor(types.Black),
)
```

### Methods
- `tl.AddTrack(tracks ...*Track) *Timeline`: Appends tracks to the timeline.
- `tl.Tracks() iter.Seq[*Track]`: Modern Go iterator over all registered tracks.
- `tl.Validate() error`: Verifies timeline consistency, media duration constraints, and stream alignment.

---

## Tracks

```go
track := timeline.NewTrack("main_video", timeline.TrackKindVideo).
    SetZIndex(0)
```

### Track Kinds
- `TrackKindVideo`: Primary background or normalized video stream.
- `TrackKindOverlay`: Picture-in-picture, alpha graphics, or animated layers.
- `TrackKindAudio`: Audio streams, music, or dialog.
- `TrackKindSubtitle`: Burned-in caption layers.
- `TrackKindWaveform`: Synthesized audio waveform visualizations.

---

## Clips

```go
clip := timeline.NewClip("clip_id", "path/to/media.mp4", 0, 10*time.Second).
    WithTrim(2*time.Second, 12*time.Second).
    WithScale(0.5).
    WithPosition(types.Point{X: 100, Y: 100}).
    WithOpacity(0.9).
    WithVolume(0.8)
```

### Fluent Chaining Methods
- `WithTrim(start, end time.Duration) *Clip`: Sets source media in/out points.
- `WithPosition(point types.Point) *Clip`: Coordinates on canvas.
- `WithScale(factor float64) *Clip`: Scale multiplier.
- `WithOpacity(alpha float64) *Clip`: Transparency `0.0` (invisible) to `1.0` (solid).
- `WithVolume(level float64) *Clip`: Audio volume factor.
