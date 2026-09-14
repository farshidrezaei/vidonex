---
layout: home

hero:
  name: "Vidonex"
  text: "Next-Gen Video Engine in Go"
  tagline: "Declarative video composition, intermediate representation DAG compiler, and modern non-linear desktop workstation."
  image:
    src: /logo.png
    alt: Vidonex Logo
  actions:
    - theme: brand
      text: Get Started
      link: /guide/quick-start
    - theme: alt
      text: Go SDK Reference
      link: /api/
    - theme: alt
      text: View on GitHub
      link: https://github.com/farshidrezaei/vidonex

features:
  - icon: ⚡
    title: Pure Go Engine
    details: Zero-dependency core compiler that turns declarative timeline trees into optimized FFmpeg filtergraphs in sub-millisecond speeds.
  - icon: 📐
    title: Graph Optimization & IR
    details: Auto-split stream isolation, dead-code elimination, and format normalization prevent FFmpeg runtime pad collisions.
  - icon: 🎬
    title: Native Desktop Workstation
    details: Full-fledged non-linear editing studio engineered with Wails v2, Vue 3, and Nuxt UI. Zero Electron bloat.
  - icon: 🎛️
    title: Studio Effects Built-in
    details: Sidechain audio ducking, reactive podcast waveforms, green screen despill, styled SRT burn-in, and keyframe motion.
  - icon: ⏱️
    title: Zero-Drift Math
    details: Exact fractional rational timing prevents cumulative millisecond drift across thousands of frames.
  - icon: 🚀
    title: Hardware Acceleration
    details: Native one-flag auto-detection for NVIDIA NVENC, Apple VideoToolbox, Intel QSV, and VAAPI.
---

## 🚀 Quick Example

Compose and render a multi-layered video with picture-in-picture in pure Go:

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

	// Layer 0: Background Video
	baseTrack := timeline.NewTrack("video", timeline.TrackKindVideo)
	baseTrack.AddClip(timeline.NewClip("intro", "gameplay.mp4", 0, 10*time.Second))

	// Layer 1: Facecam PiP Overlay
	pipTrack := timeline.NewTrack("pip", timeline.TrackKindOverlay).SetZIndex(1)
	pipTrack.AddClip(
		timeline.NewClip("webcam", "facecam.mp4", 0, 10*time.Second).
			WithScale(0.3).
			WithPosition(types.Point{X: 1300, Y: 40}),
	)

	tl.AddTrack(baseTrack, pipTrack)

	// Compile & Render via OS FFmpeg
	c := composer.New()
	_, _ = c.Render(context.Background(), tl, "output.mp4", nil)
}
```
