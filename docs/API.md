# Vidonex API Reference

Complete API documentation for the **Vidonex** video composition engine and compiler.

---

## Table of Contents

- [Package `types`](#package-types)
  - [`Rational`](#rational)
  - [`Size`, `Point`, `Rect`](#size-point-rect)
  - [`Color`](#color)
- [Package `timeline`](#package-timeline)
  - [`Timeline`](#timeline)
  - [`Track`](#track)
  - [`Clip`](#clip)
  - [`Transition`](#transition)
- [Package `filtergraph`](#package-filtergraph)
  - [`Graph`](#graph)
  - [`Node`](#node)
  - [`Pad`](#pad)
  - [`Pass` & `Pipeline`](#pass--pipeline)
- [Package `compiler`](#package-compiler)
  - [`Compiler`](#compiler)
  - [`CompilationResult`](#compilationresult)
  - [`EncodingOptions`](#encodingoptions)
- [Package `executor`](#package-executor)
  - [`CommandExecutor`](#commandexecutor)
  - [`OSExecutor`](#osexecutor)
  - [`MockExecutor`](#mockexecutor)
  - [`ProgressEvent`](#progressevent)
- [Package `composer`](#package-composer)
  - [`Composer`](#composer)

---

## Package `types`

```go
import "github.com/farshidrezaei/vidonex/types"
```

### `Rational`
Exact fractional number model (`Num/Den int64`) avoiding floating-point precision drift.

```go
r := types.NewRational(16, 9)
r2 := types.RationalFromDuration(5 * time.Second)
r3 := types.RationalFromFloat(29.97)

// Arithmetic
sum := r.Add(r2)
diff := r.Sub(r2)
prod := r.Mul(r2)
quot := r.Div(r2)

// Frame rate conversions
frames := r.ToFrames(types.FPS30)
dur := types.FramesToRational(frames, types.FPS30)
```

Standard Frame Rates:
- `types.FPS24`, `types.FPS25`, `types.FPS30`, `types.FPS50`, `types.FPS60`
- `types.FPS23_976`, `types.FPS29_97`, `types.FPS59_94`

### `Size`, `Point`, `Rect`
```go
sz := types.NewSize(1920, 1080)
scaled := sz.Scale(0.5)               // 960x540 (auto-rounded to even integers)
fit := sz.FitWithin(types.Res720p)    // Aspect ratio fitting
fill := sz.FillWithin(types.Res4K)    // Aspect ratio filling

pt := types.CalculateOffset(containerSize, itemSize, types.AlignCenter)
```

### `Color`
```go
black := types.ColorBlack
trans := types.ColorTransparent
custom := types.RGB(255, 128, 0)
alpha := types.RGBA(255, 255, 255, 128)
fromHex, err := types.Hex("#FFAA0080")

// Format for FFmpeg filters
fmtStr := custom.FFmpegColor() // "0xFFAA00"
```

---

## Package `timeline`

```go
import "github.com/farshidrezaei/vidonex/timeline"
```

### `Timeline`
```go
tl := timeline.New(
    timeline.WithCanvas(types.Res1080p),
    timeline.WithFPS(types.FPS60),
    timeline.WithBackgroundColor(types.ColorBlack),
    timeline.WithDuration(15 * time.Second), // Optional explicit override
)

tl.AddTrack(videoTrack, audioTrack)
err := timeline.Validate(tl)
```

### `Track`
```go
vTrack := timeline.NewTrack("main_video", timeline.TrackKindVideo).
    SetZIndex(0).
    SetVolume(1.0).
    SetMuted(false)

vTrack.AddClip(clip1, clip2)
```

### `Clip`
```go
clip := timeline.NewClip("clip_id", "path/to/media.mp4", 0, 10*time.Second).
    WithTrim(2 * time.Second).                   // Skip first 2s of source
    WithSpeed(1.5).                              // 1.5x playback speed
    WithVolume(0.8).                             // 80% audio volume
    WithOpacity(0.9).                            // 90% visual opacity
    WithPosition(types.Point{X: 100, Y: 100}).   // Canvas placement
    WithScale(0.5).                              // 50% scale
    WithRotation(332.0).                         // 332 degrees rotation
    WithBlendMode("multiply")                    // Compositing blend mode
```

### `Transition`
```go
trans := timeline.NewTransition("trans_1_2", timeline.TransitionDissolve, 1*time.Second, clip1, clip2)
vTrack.AddTransition(trans)
```
Supported transition types:
- `TransitionDissolve` (`"dissolve"`)
- `TransitionFade` (`"fade"`)
- `TransitionWipeLeft` (`"wipeleft"`)
- `TransitionWipeRight` (`"wiperight"`)
- `TransitionWipeUp` (`"wipeup"`)
- `TransitionWipeDown` (`"wipedown"`)
- `TransitionSlideLeft` (`"slideleft"`)
- `TransitionSlideRight` (`"slideright"`)
- `TransitionCircleCrop` (`"circlecrop"`)
- `TransitionCircleOpen` (`"circleopen"`)
- `TransitionZoomIn` (`"zoomin"`)
- `TransitionAcrossFade` (`"acrossfade"`)
```

---

## Package `filtergraph`

```go
import "github.com/farshidrezaei/vidonex/filtergraph"
```

### `Graph`
Directed Acyclic Graph modeling FFmpeg filter operations.

```go
g := filtergraph.NewGraph()

// Create nodes
scaleNode := g.NewNode("scale_1", "scale")
scaleNode.SetParam("w", 1920).SetParam("h", 1080)
inPad := scaleNode.AddInput("0:v", filtergraph.StreamTypeVideo)
outPad := scaleNode.AddOutput("v_scaled", filtergraph.StreamTypeVideo)

// Connect nodes
_ = g.Connect(outPad, nextInPad)

// Run optimization passes
pipeline := filtergraph.DefaultPipeline(nil)
_ = pipeline.Execute(g)

// Generate -filter_complex string
str, err := g.FormattedFilterComplex()
```

---

## Package `compiler`

```go
import "github.com/farshidrezaei/vidonex/compiler"
```

```go
c := compiler.New(logger).SetEncodingOptions(compiler.EncodingOptions{
    VideoCodec:   "libx264",
    AudioCodec:   "aac",
    PixelFormat:  "yuv420p",
    AudioBitrate: "192k",
    CRF:          23,
    Preset:       "medium",
})

result, err := c.Compile(tl, "output.mp4")
// result.Args          -> Full FFmpeg CLI flags
// result.FilterComplex -> -filter_complex string
// result.Mermaid()     -> Mermaid.js flowchart
// result.DOT()         -> Graphviz DOT string
```

---

## Package `composer`

```go
import "github.com/farshidrezaei/vidonex/composer"
```

```go
c := composer.New(
    composer.WithLogger(slog.Default()),
    composer.WithBinaryPath("ffmpeg"),
)

// Dry-run inspection
compilation, err := c.Compile(tl, "output.mp4")

// Render with progress callback
ctx := context.Background()
result, err := c.Render(ctx, tl, "output.mp4", func(ev executor.ProgressEvent) {
    fmt.Printf("Progress: %.1f%% (Speed: %.2fx, FPS: %.1f)\n", ev.Percentage, ev.Speed, ev.FPS)
})
```

---

## Package `server` & RESTful API

Vidonex includes an embedded pure-Go server with SQLite persistence, REST endpoints, and WebSocket telemetry.

### Starting the Server programmatically:

```go
import "github.com/farshidrezaei/vidonex/server"

srv, err := server.New(server.Config{
    Port:            8080,
    Host:            "0.0.0.0",
    DataDirectory:   "/path/to/.vidonex",
    StaticDirectory: "/path/to/ui/dist",
    Logger:          slog.Default(),
})
if err != nil {
    log.Fatal(err)
}

if err := srv.Start(); err != nil {
    log.Fatal(err)
}
```

### REST API Endpoints:

| Method | Route | Description |
| :--- | :--- | :--- |
| `GET` | `/api/projects` | List all saved projects |
| `POST` | `/api/projects` | Create a new project workspace |
| `GET` | `/api/projects/:id` | Get project details and JSON timeline specification |
| `PUT` | `/api/projects/:id` | Update project metadata and video timeline specification |
| `DELETE` | `/api/projects/:id` | Delete project workspace |
| `GET` | `/api/projects/:id/assets`| List all uploaded media assets for project |
| `POST` | `/api/media/upload` | Upload video, audio, or image asset (with auto ffprobe) |
| `DELETE` | `/api/media/:id` | Delete media asset from disk and database |
| `GET` | `/api/media/files/*` | Range-supported streaming media file server |
| `POST` | `/api/spec/validate` | Validate declarative video timeline specification |
| `POST` | `/api/spec/graph` | Generate Mermaid flowchart and Graphviz DOT from spec |
| `POST` | `/api/render/start` | Start asynchronous render job with progress telemetry |
| `GET` | `/api/render/:id` | Get render job status, progress percentage, and output |
| `POST` | `/api/render/:id/cancel` | Cancel active background render process |
| `GET` | `/api/exports/*` | Download completed rendered MP4 video exports |
| `GET` | `/api/version/check` | Check GitHub release updates, return version and host diagnostics |
| `POST` | `/api/version/upgrade` | Trigger in-app self-update with SSE progress streaming |

### WebSocket Telemetry (`/ws`):

Clients can connect to `ws://localhost:8080/ws` to receive real-time render telemetry:

```json
{
  "type": "render_progress",
  "job_id": "job_abc123",
  "data": {
    "frame": 120,
    "fps": 58.4,
    "current_time": 4.0,
    "percentage": 40.0,
    "speed": 1.95,
    "bitrate": "4500kbits/s",
    "total_size": 2400000
  }
}
```
