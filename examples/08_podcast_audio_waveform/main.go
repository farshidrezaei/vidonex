// Package main demonstrates creating a modern podcast audiogram with animated waveform visualization using Vidonyx.
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
	"github.com/farshidrezaei/vidonyx/waveform"
)

func main() {
	// 1. Setup Square 1080x1080 Timeline (Standard for Instagram / Twitter / LinkedIn Audiograms)
	compositionTimeline := timeline.New(
		timeline.WithCanvas(types.ResSquare1080),
		timeline.WithFPS(types.FPS30),
	)

	// 2. Podcast Cover Art / Video Background Layer
	bgTrack := timeline.NewTrack("cover_art", timeline.TrackKindVideo).SetZIndex(0)
	bgClip := timeline.NewClip("artwork", "assets/podcast_cover.jpg", 0, 12*time.Second)
	bgTrack.AddClip(bgClip)

	// 3. Audio Voice Episode Track
	audioTrack := timeline.NewTrack("episode_audio", timeline.TrackKindAudio)
	audioClip := timeline.NewClip("voice", "assets/podcast_clip.mp3", 0, 12*time.Second)
	audioTrack.AddClip(audioClip)

	// 4. Animated Waveform Visualizer Layer (driven by the voice audio track)
	waveformTrack := timeline.NewWaveformTrack(
		"waveform_overlay",
		"episode_audio",
		waveform.Options{
			Size:     types.NewSize(900, 160),
			Mode:     waveform.ModePeakToPeak,
			Color:    types.RGB(0, 240, 180), // Vibrant neon mint
			Scale:    "sqrt",
			FPS:      types.FPS30,
			Position: types.Point{X: 90, Y: 840}, // Bottom area of the square canvas
		},
	).SetZIndex(1)

	compositionTimeline.AddTrack(bgTrack, audioTrack, waveformTrack)

	// 5. Initialize Composer & Compile
	composerInstance := composer.New(
		composer.WithExecutor(executor.NewMockExecutor()),
	)

	compilation, err := composerInstance.Compile(compositionTimeline, "output_audiogram.mp4")
	if err != nil {
		log.Fatalf("Compilation failed: %v", err)
	}

	fmt.Println("=== Generated Filtergraph String for Audiogram ===")
	fmt.Println(compilation.FilterComplex)
	fmt.Println()

	mermaidDiagram, _ := compilation.Mermaid()
	fmt.Println("=== Mermaid Waveform Flowchart ===")
	fmt.Println(mermaidDiagram)

	// 6. Render
	ctx := context.Background()
	_, err = composerInstance.Render(ctx, compositionTimeline, "output_audiogram.mp4", func(event executor.ProgressEvent) {
		fmt.Printf("[Audiogram Render] Speed: %.2fx | %.1f%%\n", event.Speed, event.Percentage)
	})
	if err != nil {
		log.Fatalf("Render failed: %v", err)
	}
}
