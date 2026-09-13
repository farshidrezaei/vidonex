// Package main demonstrates Picture-in-Picture overlay composition with Vidonex.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/farshidrezaei/vidonex/composer"
	"github.com/farshidrezaei/vidonex/executor"
	"github.com/farshidrezaei/vidonex/timeline"
	"github.com/farshidrezaei/vidonex/types"
)

func main() {
	// 1. Setup Canvas
	tl := timeline.New(
		timeline.WithCanvas(types.Res1080p),
		timeline.WithFPS(types.FPS60),
	)

	// 2. Base Video Layer (Layer 0)
	baseTrack := timeline.NewTrack("base_track", timeline.TrackKindVideo).SetZIndex(0)
	baseClip := timeline.NewClip("gameplay", "assets/gameplay.mp4", 0, 15*time.Second)
	baseTrack.AddClip(baseClip)

	// 3. Picture-in-Picture Webcam Overlay (Layer 1)
	pipTrack := timeline.NewTrack("pip_track", timeline.TrackKindOverlay).SetZIndex(1)
	webcamClip := timeline.NewClip("webcam", "assets/camera.mp4", 2*time.Second, 10*time.Second).
		WithScale(0.25).
		WithPosition(types.Point{X: 1920 - 480 - 40, Y: 40}). // Top-Right corner with margin
		WithOpacity(0.95)
	pipTrack.AddClip(webcamClip)

	// 4. Background Music Layer
	bgmTrack := timeline.NewTrack("bgm_track", timeline.TrackKindAudio)
	bgmClip := timeline.NewClip("bgm", "assets/ambient.mp3", 0, 15*time.Second).
		WithVolume(0.4)
	bgmTrack.AddClip(bgmClip)

	tl.AddTrack(baseTrack, pipTrack, bgmTrack)

	// 5. Initialize Composer & Render
	c := composer.New(
		composer.WithExecutor(executor.NewMockExecutor()),
	)

	comp, err := c.Compile(tl, "output_pip.mp4")
	if err != nil {
		log.Fatalf("Compilation failed: %v", err)
	}

	dotGraph, _ := comp.DOT()
	fmt.Println("=== Graphviz DOT Graph ===")
	fmt.Println(dotGraph)

	ctx := context.Background()
	_, err = c.Render(ctx, tl, "output_pip.mp4", func(event executor.ProgressEvent) {
		fmt.Printf("[Render Progress] Speed: %.2fx | %.1f%%\n", event.Speed, event.Percentage)
	})
	if err != nil {
		log.Fatalf("Render failed: %v", err)
	}
}
