// Package main demonstrates automatic background music audio ducking during voiceover with Vidonex.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/farshidrezaei/vidonex/composer"
	"github.com/farshidrezaei/vidonex/ducking"
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

	// 2. Base Video Layer (e.g. gameplay or cinematic scene)
	videoTrack := timeline.NewTrack("cinematic_video", timeline.TrackKindVideo).SetZIndex(0)
	videoClip := timeline.NewClip("scene", "assets/gameplay.mp4", 0, 15*time.Second)
	videoTrack.AddClip(videoClip)

	// 3. Voiceover Track (Host speaking between 3s and 10s)
	voiceTrack := timeline.NewTrack("voiceover_track", timeline.TrackKindAudio)
	voiceClip := timeline.NewClip("commentary", "assets/speech.mp3", 3*time.Second, 7*time.Second).
		WithVolume(1.0)
	voiceTrack.AddClip(voiceClip)

	// 4. Background Music Track (Configured to automatically duck under voiceover)
	bgmTrack := timeline.NewTrack("bgm_track", timeline.TrackKindAudio).
		WithDucking("voiceover_track", ducking.Options{
			Threshold:           0.08,
			Ratio:               5.0,
			AttackMilliseconds:  25,
			ReleaseMilliseconds: 400,
		})

	bgmClip := timeline.NewClip("soundtrack", "assets/music.mp3", 0, 15*time.Second).
		WithVolume(0.8)
	bgmTrack.AddClip(bgmClip)

	compositionTimeline.AddTrack(videoTrack, voiceTrack, bgmTrack)

	// 5. Initialize Composer & Compile
	composerInstance := composer.New(
		composer.WithExecutor(executor.NewMockExecutor()),
	)

	compilation, err := composerInstance.Compile(compositionTimeline, "output_ducked.mp4")
	if err != nil {
		log.Fatalf("Compilation failed: %v", err)
	}

	fmt.Println("=== Generated Filtergraph String ===")
	fmt.Println(compilation.FilterComplex)
	fmt.Println()

	mermaidDiagram, _ := compilation.Mermaid()
	fmt.Println("=== Mermaid Ducking Flowchart ===")
	fmt.Println(mermaidDiagram)

	// 6. Render
	ctx := context.Background()
	_, err = composerInstance.Render(ctx, compositionTimeline, "output_ducked.mp4", func(event executor.ProgressEvent) {
		fmt.Printf("[Render Progress] Speed: %.2fx | %.1f%%\n", event.Speed, event.Percentage)
	})
	if err != nil {
		log.Fatalf("Render failed: %v", err)
	}
}
