// Package main demonstrates dynamic keyframing, easing animations, Ken Burns effects, and fade transitions with Vidonyx.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/farshidrezaei/vidonyx/animation"
	"github.com/farshidrezaei/vidonyx/composer"
	"github.com/farshidrezaei/vidonyx/executor"
	"github.com/farshidrezaei/vidonyx/timeline"
	"github.com/farshidrezaei/vidonyx/types"
)

func main() {
	// 1. Setup 1080p 60fps Timeline
	compositionTimeline := timeline.New(
		timeline.WithCanvas(types.Res1080p),
		timeline.WithFPS(types.FPS60),
	)

	// 2. Base Video Track with Ken Burns Zoom and Fade In/Out
	baseTrack := timeline.NewTrack("base_track", timeline.TrackKindVideo).SetZIndex(0)

	scaleTrack := animation.NewFloatKeyframeTrack().
		AddKeyframe(0*time.Second, 1.0, animation.EasingLinear).
		AddKeyframe(10*time.Second, 1.25, animation.EasingEaseInOutCubic)

	baseClip := timeline.NewClip("main_scenery", "assets/scenery.mp4", 0, 10*time.Second).
		WithScaleTrack(scaleTrack).
		WithFadeIn(1500 * time.Millisecond).
		WithFadeOut(2000 * time.Millisecond)

	baseTrack.AddClip(baseClip)

	// 3. Overlay Layer with Animated Motion Path (Moving Picture-in-Picture)
	overlayTrack := timeline.NewTrack("moving_overlay", timeline.TrackKindOverlay).SetZIndex(1)

	posTrack := animation.NewPositionTrack().
		AddKeyframe(0*time.Second, types.Point{X: 100, Y: 100}, animation.EasingLinear).
		AddKeyframe(4*time.Second, types.Point{X: 1400, Y: 100}, animation.EasingEaseInOutQuad).
		AddKeyframe(8*time.Second, types.Point{X: 1400, Y: 700}, animation.EasingEaseOutQuad)

	overlayClip := timeline.NewClip("floating_badge", "assets/badge.png", 0, 10*time.Second).
		WithScale(0.2).
		WithPositionTrack(posTrack).
		WithFadeIn(1 * time.Second)

	overlayTrack.AddClip(overlayClip)

	compositionTimeline.AddTrack(baseTrack, overlayTrack)

	// 4. Initialize Composer & Compile
	composerInstance := composer.New(
		composer.WithExecutor(executor.NewMockExecutor()),
	)

	compilation, err := composerInstance.Compile(compositionTimeline, "output_animated.mp4")
	if err != nil {
		log.Fatalf("Compilation failed: %v", err)
	}

	fmt.Println("=== Generated Filtergraph String ===")
	fmt.Println(compilation.FilterComplex)
	fmt.Println()

	mermaidDiagram, _ := compilation.Mermaid()
	fmt.Println("=== Mermaid Animation Flowchart ===")
	fmt.Println(mermaidDiagram)

	// 5. Render
	ctx := context.Background()
	_, err = composerInstance.Render(ctx, compositionTimeline, "output_animated.mp4", func(event executor.ProgressEvent) {
		fmt.Printf("[Render Progress] Speed: %.2fx | %.1f%%\n", event.Speed, event.Percentage)
	})
	if err != nil {
		log.Fatalf("Render failed: %v", err)
	}
}
