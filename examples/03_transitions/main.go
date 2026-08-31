// Package main demonstrates sequential video transitions (xfade) and audio cross-fading (acrossfade) with Vidonyx.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/farshidrezaei/vidonyx/composer"
	"github.com/farshidrezaei/vidonyx/executor"
	"github.com/farshidrezaei/vidonyx/timeline"
	"github.com/farshidrezaei/vidonyx/types"
)

func main() {
	// 1. Setup 1080p 30fps Timeline
	compositionTimeline := timeline.New(
		timeline.WithCanvas(types.Res1080p),
		timeline.WithFPS(types.FPS30),
	)

	// 2. Main Video Track with 3 sequential scenes and transitions
	mainVideoTrack := timeline.NewTrack("main_sequence", timeline.TrackKindVideo).SetZIndex(0)

	scene1 := timeline.NewClip("scene_1", "assets/scene1.mp4", 0, 5*time.Second)
	scene2 := timeline.NewClip("scene_2", "assets/scene2.mp4", 4*time.Second, 6*time.Second)
	scene3 := timeline.NewClip("scene_3", "assets/scene3.mp4", 9*time.Second, 5*time.Second)

	// Transition 1: Dissolve from Scene 1 to Scene 2 (1.0 second duration)
	transition1 := timeline.NewTransition("trans_1_2", timeline.TransitionDissolve, 1*time.Second, scene1, scene2)

	// Transition 2: WipeLeft from Scene 2 to Scene 3 (1.0 second duration)
	transition2 := timeline.NewTransition("trans_2_3", timeline.TransitionWipeLeft, 1*time.Second, scene2, scene3)

	mainVideoTrack.AddClip(scene1, scene2, scene3)
	mainVideoTrack.AddTransition(transition1, transition2)

	compositionTimeline.AddTrack(mainVideoTrack)

	// 3. Initialize Composer with MockExecutor
	composerInstance := composer.New(
		composer.WithExecutor(executor.NewMockExecutor()),
	)

	// 4. Compile & Output
	compilation, err := composerInstance.Compile(compositionTimeline, "output_transitions.mp4")
	if err != nil {
		log.Fatalf("Compilation failed: %v", err)
	}

	fmt.Println("=== Generated Filtergraph String ===")
	fmt.Println(compilation.FilterComplex)
	fmt.Println()

	mermaidDiagram, _ := compilation.Mermaid()
	fmt.Println("=== Mermaid Transition Flowchart ===")
	fmt.Println(mermaidDiagram)

	// 5. Render
	ctx := context.Background()
	_, err = composerInstance.Render(ctx, compositionTimeline, "output_transitions.mp4", func(event executor.ProgressEvent) {
		fmt.Printf("[Render Telemetry] Progress: %.1f%% (Speed: %.2fx)\n", event.Percentage, event.Speed)
	})
	if err != nil {
		log.Fatalf("Render failed: %v", err)
	}
}
