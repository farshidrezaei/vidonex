package compiler_test

import (
	"strings"
	"testing"
	"time"

	"github.com/farshidrezaei/vidonyx/compiler"
	"github.com/farshidrezaei/vidonyx/timeline"
	"github.com/farshidrezaei/vidonyx/types"
)

func TestCompiler_PictureInPictureComposition(t *testing.T) {
	// Build a 1080p 30fps timeline with a base background video and an overlay PiP clip
	tl := timeline.New(
		timeline.WithCanvas(types.Res1080p),
		timeline.WithFPS(types.FPS30),
	)

	// Track 0: Main background video
	mainTrack := timeline.NewTrack("main_track", timeline.TrackKindVideo).SetZIndex(0)
	mainClip := timeline.NewClip("main_clip", "background.mp4", 0, 10*time.Second)
	mainTrack.AddClip(mainClip)

	// Track 1: PiP Overlay (Z-index 1)
	pipTrack := timeline.NewTrack("pip_track", timeline.TrackKindOverlay).SetZIndex(1)
	pipClip := timeline.NewClip("pip_clip", "facecam.mp4", 2*time.Second, 6*time.Second).
		WithScale(0.3).
		WithPosition(types.Point{X: 50, Y: 50}).
		WithOpacity(0.9)
	pipTrack.AddClip(pipClip)

	// Track 2: Background Music
	bgmTrack := timeline.NewTrack("bgm_track", timeline.TrackKindAudio)
	bgmClip := timeline.NewClip("bgm_clip", "music.mp3", 0, 10*time.Second).WithVolume(0.5)
	bgmTrack.AddClip(bgmClip)

	tl.AddTrack(mainTrack, pipTrack, bgmTrack)

	c := compiler.New(nil)
	res, err := c.Compile(tl, "output.mp4")
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	if len(res.Inputs) != 3 {
		t.Fatalf("Expected 3 unique inputs, got %d: %v", len(res.Inputs), res.Inputs)
	}

	// Verify CLI arguments
	fullCommand := strings.Join(res.Args, " ")
	if !strings.Contains(fullCommand, "-filter_complex") {
		t.Fatal("Expected -filter_complex flag in CLI arguments")
	}
	if !strings.Contains(fullCommand, "-map [out_v]") || !strings.Contains(fullCommand, "-map [out_a]") {
		t.Fatal("Expected -map [out_v] and -map [out_a] in CLI arguments")
	}

	// Verify Mermaid generation
	mermaidDiagram, err := res.Mermaid()
	if err != nil {
		t.Fatalf("Mermaid generation failed: %v", err)
	}
	if !strings.Contains(mermaidDiagram, "graph LR") {
		t.Fatalf("Expected Mermaid diagram to start with graph LR:\n%s", mermaidDiagram)
	}

	// Verify DOT generation
	dotGraph, err := res.DOT()
	if err != nil {
		t.Fatalf("DOT generation failed: %v", err)
	}
	if !strings.Contains(dotGraph, "digraph Filtergraph") {
		t.Fatalf("Expected DOT graph to contain digraph Filtergraph:\n%s", dotGraph)
	}
}
