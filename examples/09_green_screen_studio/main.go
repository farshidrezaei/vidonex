// Package main demonstrates studio green screen chroma key removal with despill suppression using Vidonex.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/farshidrezaei/vidonex/chromakey"
	"github.com/farshidrezaei/vidonex/composer"
	"github.com/farshidrezaei/vidonex/executor"
	"github.com/farshidrezaei/vidonex/timeline"
	"github.com/farshidrezaei/vidonex/types"
)

func main() {
	// 1. Setup 1080p 30fps Timeline
	compositionTimeline := timeline.New(
		timeline.WithCanvas(types.Res1080p),
		timeline.WithFPS(types.FPS30),
	)

	// 2. Futuristic Studio Background Layer (Layer 0)
	bgTrack := timeline.NewTrack("background_track", timeline.TrackKindVideo).SetZIndex(0)
	bgClip := timeline.NewClip("cyber_studio", "assets/sci_fi_hq.mp4", 0, 10*time.Second)
	bgTrack.AddClip(bgClip)

	// 3. Studio Presenter with Green Screen Keying & Despill (Layer 1)
	presenterTrack := timeline.NewTrack("presenter_track", timeline.TrackKindOverlay).SetZIndex(1)
	presenterClip := timeline.NewClip("host", "assets/presenter_greenscreen.mp4", 0, 10*time.Second).
		WithChromaKey(chromakey.Options{
			KeyColor:      chromakey.StudioGreenScreen, // 0x00FF00
			Similarity:    0.32,                        // Fine-tuned green similarity
			Blend:         0.12,                        // Smooth edge feathering
			Despill:       true,                        // Suppress green reflection on skin/clothes
			DespillType:   chromakey.DespillGreen,
			DespillExpand: 0.0,
		}).
		WithScale(0.85).
		WithPosition(types.Point{X: 150, Y: 160}) // Align lower-left
	presenterTrack.AddClip(presenterClip)

	compositionTimeline.AddTrack(bgTrack, presenterTrack)

	// 4. Initialize Composer & Compile
	composerInstance := composer.New(
		composer.WithExecutor(executor.NewMockExecutor()),
	)

	compilation, err := composerInstance.Compile(compositionTimeline, "output_studio_greenscreen.mp4")
	if err != nil {
		log.Fatalf("Compilation failed: %v", err)
	}

	fmt.Println("=== Generated Chroma Key Filtergraph ===")
	fmt.Println(compilation.FilterComplex)
	fmt.Println()

	mermaidDiagram, _ := compilation.Mermaid()
	fmt.Println("=== Mermaid Studio Compositor Flowchart ===")
	fmt.Println(mermaidDiagram)

	// 5. Render
	ctx := context.Background()
	_, err = composerInstance.Render(ctx, compositionTimeline, "output_studio_greenscreen.mp4", func(event executor.ProgressEvent) {
		fmt.Printf("[Chroma Render Progress] Speed: %.2fx | %.1f%%\n", event.Speed, event.Percentage)
	})
	if err != nil {
		log.Fatalf("Render failed: %v", err)
	}
}
