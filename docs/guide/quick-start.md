# Quick Start

Get started with Vidonex via the Desktop Studio, CLI tool, or Go SDK.

## Prerequisites

- **FFmpeg & FFprobe**: Version 5.0 or later installed in your `$PATH`.
  ```bash
  # macOS
  brew install ffmpeg

  # Ubuntu / Debian
  sudo apt install ffmpeg

  # Windows (via Chocolatey or Scoop)
  choco install ffmpeg
  ```

---

## 1. Using the Desktop Studio

Pre-built standalone executables are available for Linux, macOS, and Windows:

1. Visit [GitHub Releases](https://github.com/farshidrezaei/vidonex/releases/latest).
2. Download the package for your operating system:
   - **Linux**: `vidonex-linux-amd64`
   - **macOS**: `vidonex-darwin-universal` (Apple Silicon & Intel)
   - **Windows**: `vidonex-windows-amd64.exe`
3. Launch the binary. No external runtime or browser installation required!

---

## 2. Using the CLI Tool

Install the standalone command-line interface:

```bash
go install github.com/farshidrezaei/vidonex/cmd/vidonex@latest
```

Verify installation:
```bash
vidonex version
```

### Render a Project
```bash
vidonex render project.yaml -o output.mp4 --gpu nvenc
```

### Start the Web Studio
```bash
vidonex serve --port 8080
```
Open your browser at `http://localhost:8080` to access the full-featured Web Studio.

---

## 3. Using the Go Engine SDK

Add the Vidonex Go module to your project:

```bash
go get github.com/farshidrezaei/vidonex
```

### Your First Programmatic Video

Create `main.go`:

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/farshidrezaei/vidonex/composer"
	"github.com/farshidrezaei/vidonex/executor"
	"github.com/farshidrezaei/vidonex/timeline"
	"github.com/farshidrezaei/vidonex/types"
)

func main() {
	// Create canvas at 1080p 60fps
	tl := timeline.New(
		timeline.WithCanvas(types.Res1080p),
		timeline.WithFPS(types.FPS60),
	)

	// Add video track
	videoTrack := timeline.NewTrack("main", timeline.TrackKindVideo)
	videoTrack.AddClip(
		timeline.NewClip("clip1", "assets/sample.mp4", 0, 5*time.Second),
	)
	tl.AddTrack(videoTrack)

	// Render using OS FFmpeg
	c := composer.New()
	_, err := c.Render(context.Background(), tl, "out.mp4", func(e executor.ProgressEvent) {
		fmt.Printf("\rProgress: %.1f%%", e.Percentage)
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("\nDone!")
}
```
