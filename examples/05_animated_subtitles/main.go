// Package main demonstrates subtitle track burn-in and styled social media captions with Vidonex.
package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/farshidrezaei/vidonex/composer"
	"github.com/farshidrezaei/vidonex/executor"
	"github.com/farshidrezaei/vidonex/subtitles"
	"github.com/farshidrezaei/vidonex/timeline"
	"github.com/farshidrezaei/vidonex/types"
)

func main() {
	// 1. Setup 1080p 30fps Timeline
	compositionTimeline := timeline.New(
		timeline.WithCanvas(types.Res1080p),
		timeline.WithFPS(types.FPS30),
	)

	// 2. Base Video Layer
	videoTrack := timeline.NewTrack("main_video", timeline.TrackKindVideo).SetZIndex(0)
	videoClip := timeline.NewClip("vlog", "assets/vlog.mp4", 0, 12*time.Second)
	videoTrack.AddClip(videoClip)

	// 3. Parse SRT Subtitles and Configure Style
	srtContent := `1
00:00:01,000 --> 00:00:04,500
Hey everyone! Welcome to Vidonex.

2
00:00:05,000 --> 00:00:08,000
Declarative, type-safe video editing in pure Go.

3
00:00:08,500 --> 00:00:11,500
Subtitles, transitions, and animations with zero effort!
`

	subTrack, err := subtitles.ParseSRT(strings.NewReader(srtContent))
	if err != nil {
		log.Fatalf("Failed to parse SRT: %v", err)
	}

	// Modern TikTok / Reels style: Yellow text with black bounding box
	subTrack.SetStyle(subtitles.SubtitleStyle{
		FontSize:       48,
		PrimaryColor:   types.RGB(255, 220, 0), // Vibrant yellow
		OutlineColor:   types.ColorBlack,
		OutlineWidth:   3,
		Box:            true,
		BoxColor:       types.RGBA(0, 0, 0, 200), // Semi-transparent black box
		BoxBorderWidth: 10,
		Alignment:      types.AlignBottomCenter,
		MarginBottom:   80,
	})

	subtitleTimelineTrack := timeline.NewSubtitleTrack("captions", subTrack)

	compositionTimeline.AddTrack(videoTrack, subtitleTimelineTrack)

	// 4. Compile & Inspect
	composerInstance := composer.New(
		composer.WithExecutor(executor.NewMockExecutor()),
	)

	compilation, err := composerInstance.Compile(compositionTimeline, "output_subtitled.mp4")
	if err != nil {
		log.Fatalf("Compilation failed: %v", err)
	}

	fmt.Println("=== Generated Filtergraph String ===")
	fmt.Println(compilation.FilterComplex)
	fmt.Println()

	mermaidDiagram, _ := compilation.Mermaid()
	fmt.Println("=== Mermaid Subtitle Flowchart ===")
	fmt.Println(mermaidDiagram)

	// 5. Render
	ctx := context.Background()
	_, err = composerInstance.Render(ctx, compositionTimeline, "output_subtitled.mp4", func(event executor.ProgressEvent) {
		fmt.Printf("[Render Progress] Speed: %.2fx | %.1f%%\n", event.Speed, event.Percentage)
	})
	if err != nil {
		log.Fatalf("Render failed: %v", err)
	}
}
