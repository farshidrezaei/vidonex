// Package main demonstrates platform-specific presets (TikTok, YouTube 4K) and GPU hardware acceleration with Vidonex.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/farshidrezaei/vidonex/composer"
	"github.com/farshidrezaei/vidonex/executor"
	"github.com/farshidrezaei/vidonex/presets"
	"github.com/farshidrezaei/vidonex/timeline"
)

func main() {
	// 1. Choose a platform preset: TikTok Vertical 1080p 60fps (9:16)
	tiktokPreset := presets.TikTokVertical1080p60()
	fmt.Printf("Using Preset: %s (%s)\n\n", tiktokPreset.Name, tiktokPreset.Description)

	// 2. Initialize Timeline and configure with the preset
	compositionTimeline := timeline.New()
	tiktokPreset.ApplyToTimeline(compositionTimeline)

	videoTrack := timeline.NewTrack("vertical_main", timeline.TrackKindVideo)
	videoClip := timeline.NewClip("portrait_video", "assets/mobile_recording.mp4", 0, 10*time.Second)
	videoTrack.AddClip(videoClip)

	compositionTimeline.AddTrack(videoTrack)

	// 3. Initialize Composer configured with Preset and Hardware Acceleration (NVENC)
	composerInstance := composer.New(
		composer.WithExecutor(executor.NewMockExecutor()),
		composer.WithPreset(tiktokPreset),
		composer.WithHardwareAcceleration(presets.AcceleratorNVENC, false),
	)

	// 4. Compile
	compilation, err := composerInstance.Compile(compositionTimeline, "output_tiktok_nvenc.mp4")
	if err != nil {
		log.Fatalf("Compilation failed: %v", err)
	}

	fmt.Println("=== Generated FFmpeg CLI Command with GPU Acceleration ===")
	for _, arg := range compilation.Args {
		fmt.Printf("%s ", arg)
	}
	fmt.Printf("\n\n")

	// 5. Render
	ctx := context.Background()
	_, err = composerInstance.Render(ctx, compositionTimeline, "output_tiktok_nvenc.mp4", func(event executor.ProgressEvent) {
		fmt.Printf("[Hardware Encode Progress] %.1f%%\n", event.Percentage)
	})
	if err != nil {
		log.Fatalf("Render failed: %v", err)
	}
}
