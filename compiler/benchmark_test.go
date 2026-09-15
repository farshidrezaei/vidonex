package compiler_test

import (
	"fmt"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/farshidrezaei/vidonex/chromakey"
	"github.com/farshidrezaei/vidonex/compiler"
	"github.com/farshidrezaei/vidonex/ducking"
	"github.com/farshidrezaei/vidonex/filtergraph"
	"github.com/farshidrezaei/vidonex/timeline"
	"github.com/farshidrezaei/vidonex/types"
	"github.com/farshidrezaei/vidonex/waveform"
)

// BenchmarkCompiler_SimpleTimeline benchmarks compiling a single-track cut composition.
func BenchmarkCompiler_SimpleTimeline(b *testing.B) {
	timelineInstance := timeline.New(
		timeline.WithCanvas(types.Res1080p),
		timeline.WithFPS(types.FPS60),
	)
	track := timeline.NewTrack("video_track_0", timeline.TrackKindVideo)
	track.AddClip(timeline.NewClip("clip_1", "intro.mp4", 0, 10*time.Second))
	track.AddClip(timeline.NewClip("clip_2", "main.mp4", 10*time.Second, 30*time.Second))
	timelineInstance.AddTrack(track)

	silentLogger := slog.New(slog.NewTextHandler(io.Discard, nil))
	compilerInstance := compiler.New(silentLogger)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		compilationResult, err := compilerInstance.Compile(timelineInstance, "output.mp4")
		if err != nil {
			b.Fatalf("failed compiling simple timeline: %v", err)
		}
		if compilationResult == nil {
			b.Fatal("nil compilation result")
		}
	}
}

// BenchmarkCompiler_ComplexMultiTrack benchmarks compiling a heavy multi-track composition with PIP, ChromaKey, Ducking, and Waveform.
func BenchmarkCompiler_ComplexMultiTrack(b *testing.B) {
	timelineInstance := timeline.New(
		timeline.WithCanvas(types.Res1080p),
		timeline.WithFPS(types.FPS60),
	)

	// Track 0: Background Video
	backgroundTrack := timeline.NewTrack("bg_track", timeline.TrackKindVideo).SetZIndex(0)
	backgroundTrack.AddClip(timeline.NewClip("bg_clip", "nature_background.mp4", 0, 30*time.Second))

	// Track 1: Presenter Green Screen Overlay
	overlayTrack := timeline.NewTrack("presenter_track", timeline.TrackKindOverlay).SetZIndex(1)
	presenterClip := timeline.NewClip("presenter_clip", "presenter_green.mp4", 0, 30*time.Second).
		WithChromaKey(chromakey.Options{
			KeyColor:    chromakey.StudioGreenScreen,
			Similarity:  0.30,
			Blend:       0.08,
			Despill:     true,
			DespillType: chromakey.DespillGreen,
		})
	overlayTrack.AddClip(presenterClip)

	// Track 2: Voice Track (Triggers Ducking)
	voiceTrack := timeline.NewTrack("voice_track", timeline.TrackKindAudio)
	voiceTrack.AddClip(timeline.NewClip("voice_clip", "narration.wav", 0, 30*time.Second))

	// Track 3: Background Music (Ducked)
	musicTrack := timeline.NewTrack("music_track", timeline.TrackKindAudio).
		WithDucking("voice_track", ducking.Options{
			Threshold:           0.1,
			Ratio:               4.0,
			AttackMilliseconds:  20,
			ReleaseMilliseconds: 300,
		})
	musicTrack.AddClip(timeline.NewClip("music_clip", "cinematic_bgm.mp3", 0, 30*time.Second))

	// Track 4: Audio Waveform Visualizer
	waveformTrack := timeline.NewWaveformTrack("waveform_layer", "voice_track", waveform.Options{
		Size:     types.NewSize(1280, 200),
		Mode:     waveform.ModePeakToPeak,
		Scale:    "sqrt",
		Position: types.Point{X: 320, Y: 800},
	})

	timelineInstance.AddTrack(backgroundTrack, overlayTrack, voiceTrack, musicTrack, waveformTrack)

	silentLogger := slog.New(slog.NewTextHandler(io.Discard, nil))
	compilerInstance := compiler.New(silentLogger)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		compilationResult, err := compilerInstance.Compile(timelineInstance, "complex_output.mp4")
		if err != nil {
			b.Fatalf("failed compiling complex multi-track timeline: %v", err)
		}
		if compilationResult == nil {
			b.Fatal("nil compilation result")
		}
	}
}

// BenchmarkFiltergraph_TopologicalSort_500Nodes benchmarks Kahn's topological sort on a 500-node DAG.
func BenchmarkFiltergraph_TopologicalSort_500Nodes(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		graph := filtergraph.NewGraph()
		previousOutput := graph.NewNode("source_0", "color").AddOutput("out", filtergraph.StreamTypeVideo)

		for nodeIndex := 1; nodeIndex < 500; nodeIndex++ {
			currentNode := graph.NewNode(fmt.Sprintf("node_%d", nodeIndex), "scale")
			currentInput := currentNode.AddInput("in", filtergraph.StreamTypeVideo)
			currentOutput := currentNode.AddOutput("out", filtergraph.StreamTypeVideo)
			_ = graph.Connect(previousOutput, currentInput)
			previousOutput = currentOutput
		}
		b.StartTimer()

		sortedNodes, err := graph.TopologicalSort()
		if err != nil {
			b.Fatalf("failed topological sort: %v", err)
		}
		if len(sortedNodes) != 500 {
			b.Fatalf("expected 500 sorted nodes, got %d", len(sortedNodes))
		}
	}
}
