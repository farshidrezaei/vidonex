// Package main demonstrates sequential video cutting and timeline trimming with Vidonyx.
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
	// 1. Build a 1080p 30fps timeline with two sequential clips
	tl := timeline.New(
		timeline.WithCanvas(types.Res1080p),
		timeline.WithFPS(types.FPS30),
		timeline.WithBackgroundColor(types.ColorBlack),
	)

	videoTrack := timeline.NewTrack("main_video", timeline.TrackKindVideo).SetZIndex(0)

	clip1 := timeline.NewClip("clip_1", "assets/intro.mp4", 0, 5*time.Second).
		WithTrim(1 * time.Second) // Trim first 1s from source

	clip2 := timeline.NewClip("clip_2", "assets/feature.mp4", 5*time.Second, 10*time.Second)

	videoTrack.AddClip(clip1, clip2)
	tl.AddTrack(videoTrack)

	// 2. Initialize Composer with MockExecutor for demonstration
	c := composer.New(
		composer.WithExecutor(executor.NewMockExecutor()),
	)

	// 3. Compile & Inspect Graph and Diagram
	compilation, err := c.Compile(tl, "output_simple_cut.mp4")
	if err != nil {
		log.Fatalf("Compilation failed: %v", err)
	}

	fmt.Println("=== Generated FFmpeg CLI Arguments ===")
	for _, arg := range compilation.Args {
		fmt.Printf("%s ", arg)
	}
	fmt.Printf("\n\n")

	mermaid, _ := compilation.Mermaid()
	fmt.Println("=== Mermaid Diagram ===")
	fmt.Println(mermaid)

	// 4. Render
	ctx := context.Background()
	_, err = c.Render(ctx, tl, "output_simple_cut.mp4", func(event executor.ProgressEvent) {
		fmt.Printf("Rendering Progress: %.1f%%\n", event.Percentage)
	})
	if err != nil {
		log.Fatalf("Render failed: %v", err)
	}
}
