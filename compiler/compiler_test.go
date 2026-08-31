package compiler_test

import (
	"strings"
	"testing"
	"time"

	"github.com/farshidrezaei/vidonyx/animation"
	"github.com/farshidrezaei/vidonyx/chromakey"
	"github.com/farshidrezaei/vidonyx/compiler"
	"github.com/farshidrezaei/vidonyx/ducking"
	"github.com/farshidrezaei/vidonyx/subtitles"
	"github.com/farshidrezaei/vidonyx/timeline"
	"github.com/farshidrezaei/vidonyx/types"
	"github.com/farshidrezaei/vidonyx/waveform"
)

func TestCompiler_CompositionsTable(t *testing.T) {
	tests := []struct {
		name                string
		buildTimeline       func() *timeline.Timeline
		encodingOpts        compiler.EncodingOptions
		outputPath          string
		expectedInputCount  int
		expectedArgSnippets []string
	}{
		{
			name: "single clip basic composition",
			buildTimeline: func() *timeline.Timeline {
				tl := timeline.New(timeline.WithCanvas(types.Res1080p), timeline.WithFPS(types.FPS30))
				tr := timeline.NewTrack("v0", timeline.TrackKindVideo)
				tr.AddClip(timeline.NewClip("c1", "input.mp4", 0, 5*time.Second))
				tl.AddTrack(tr)
				return tl
			},
			encodingOpts:       compiler.DefaultEncodingOptions(),
			outputPath:         "out_single.mp4",
			expectedInputCount: 1,
			expectedArgSnippets: []string{
				"-i input.mp4",
				"-filter_complex",
				"-map [out_v]",
				"-map [out_a]",
				"-c:v libx264",
				"out_single.mp4",
			},
		},
		{
			name: "green screen studio clip with despill",
			buildTimeline: func() *timeline.Timeline {
				tl := timeline.New(timeline.WithCanvas(types.Res1080p), timeline.WithFPS(types.FPS30))

				bgTrack := timeline.NewTrack("bg", timeline.TrackKindVideo).SetZIndex(0)
				bgTrack.AddClip(timeline.NewClip("bg_clip", "futuristic_room.mp4", 0, 8*time.Second))

				presenterTrack := timeline.NewTrack("presenter", timeline.TrackKindOverlay).SetZIndex(1)
				presenterClip := timeline.NewClip("host", "studio_green.mp4", 0, 8*time.Second).
					WithChromaKey(chromakey.Options{
						KeyColor:    chromakey.StudioGreenScreen,
						Similarity:  0.30,
						Blend:       0.10,
						Despill:     true,
						DespillType: chromakey.DespillGreen,
					})
				presenterTrack.AddClip(presenterClip)

				tl.AddTrack(bgTrack, presenterTrack)
				return tl
			},
			encodingOpts:       compiler.DefaultEncodingOptions(),
			outputPath:         "out_greenscreen.mp4",
			expectedInputCount: 2,
			expectedArgSnippets: []string{
				"chromakey=color=0x00FF00:similarity=0.30:blend=0.10",
				"despill=type=green:expand=0.00",
				"out_greenscreen.mp4",
			},
		},
		{
			name: "podcast audiogram with audio waveform visualizer",
			buildTimeline: func() *timeline.Timeline {
				tl := timeline.New(timeline.WithCanvas(types.Res1080p), timeline.WithFPS(types.FPS30))

				bgTrack := timeline.NewTrack("bg", timeline.TrackKindVideo)
				bgTrack.AddClip(timeline.NewClip("poster", "cover.jpg", 0, 10*time.Second))

				audioTrack := timeline.NewTrack("voice", timeline.TrackKindAudio)
				audioTrack.AddClip(timeline.NewClip("speech", "podcast.mp3", 0, 10*time.Second))

				waveTrack := timeline.NewWaveformTrack("wave_layer", "voice", waveform.Options{
					Size:     types.NewSize(1000, 200),
					Mode:     waveform.ModePeakToPeak,
					Scale:    "sqrt",
					Position: types.Point{X: 460, Y: 800},
				})

				tl.AddTrack(bgTrack, audioTrack, waveTrack)
				return tl
			},
			encodingOpts:       compiler.DefaultEncodingOptions(),
			outputPath:         "out_audiogram.mp4",
			expectedInputCount: 2,
			expectedArgSnippets: []string{
				"showwaves=s=1000x200",
				"mode=p2p",
				"scale=sqrt",
				"format=pix_fmts=yuva420p",
				"overlay=x=460:y=800",
				"out_audiogram.mp4",
			},
		},
		{
			name: "music track auto ducking under voiceover track",
			buildTimeline: func() *timeline.Timeline {
				tl := timeline.New(timeline.WithCanvas(types.Res1080p), timeline.WithFPS(types.FPS30))

				videoTrack := timeline.NewTrack("v0", timeline.TrackKindVideo)
				videoTrack.AddClip(timeline.NewClip("v1", "gameplay.mp4", 0, 10*time.Second))

				voiceTrack := timeline.NewTrack("voice_track", timeline.TrackKindAudio)
				voiceTrack.AddClip(timeline.NewClip("speech", "dialogue.mp3", 2*time.Second, 5*time.Second))

				musicTrack := timeline.NewTrack("music_track", timeline.TrackKindAudio).
					WithDucking("voice_track", ducking.Options{
						Threshold:           0.1,
						Ratio:               4.0,
						AttackMilliseconds:  20,
						ReleaseMilliseconds: 300,
					})
				musicTrack.AddClip(timeline.NewClip("bgm", "ambient.mp3", 0, 10*time.Second))

				tl.AddTrack(videoTrack, voiceTrack, musicTrack)
				return tl
			},
			encodingOpts:       compiler.DefaultEncodingOptions(),
			outputPath:         "out_ducked.mp4",
			expectedInputCount: 3,
			expectedArgSnippets: []string{
				"sidechaincompress=threshold=0.100:ratio=4.0:attack=20:release=300",
				"asplit",
				"amix=inputs=3",
				"out_ducked.mp4",
			},
		},
		{
			name: "clip with subtitle track burn in",
			buildTimeline: func() *timeline.Timeline {
				tl := timeline.New(timeline.WithCanvas(types.Res1080p), timeline.WithFPS(types.FPS30))
				videoTrack := timeline.NewTrack("v0", timeline.TrackKindVideo)
				videoTrack.AddClip(timeline.NewClip("c1", "input.mp4", 0, 5*time.Second))

				subTrack := subtitles.NewSubtitleTrack().
					AddCue(subtitles.SubtitleCue{
						Index:     1,
						StartTime: 1 * time.Second,
						EndTime:   4 * time.Second,
						Text:      "Vidonyx Subtitles In Action",
					})

				subTimelineTrack := timeline.NewSubtitleTrack("sub_track", subTrack)

				tl.AddTrack(videoTrack, subTimelineTrack)
				return tl
			},
			encodingOpts:       compiler.DefaultEncodingOptions(),
			outputPath:         "out_sub.mp4",
			expectedInputCount: 1,
			expectedArgSnippets: []string{
				"drawtext=",
				"Vidonyx Subtitles In Action",
				"between(t,1.0000,4.0000)",
				"out_sub.mp4",
			},
		},
		{
			name: "clip with fade in and fade out",
			buildTimeline: func() *timeline.Timeline {
				tl := timeline.New(timeline.WithCanvas(types.Res1080p), timeline.WithFPS(types.FPS30))
				tr := timeline.NewTrack("v0", timeline.TrackKindVideo)
				clip := timeline.NewClip("c1", "input.mp4", 0, 10*time.Second).
					WithFadeIn(1 * time.Second).
					WithFadeOut(2 * time.Second)
				tr.AddClip(clip)
				tl.AddTrack(tr)
				return tl
			},
			encodingOpts:       compiler.DefaultEncodingOptions(),
			outputPath:         "out_fade.mp4",
			expectedInputCount: 1,
			expectedArgSnippets: []string{
				"fade=t=in:st=0:d=1.0000",
				"fade=t=out:st=8.0000:d=2.0000",
				"afade=t=in:st=0:d=1.0000",
				"afade=t=out:st=8.0000:d=2.0000",
			},
		},
		{
			name: "clip with animated position and scale tracks",
			buildTimeline: func() *timeline.Timeline {
				tl := timeline.New(timeline.WithCanvas(types.Res1080p), timeline.WithFPS(types.FPS30))
				tr := timeline.NewTrack("v0", timeline.TrackKindVideo)

				posTrack := animation.NewPositionTrack().
					AddKeyframe(0, types.Point{X: 100, Y: 100}, animation.EasingLinear).
					AddKeyframe(5*time.Second, types.Point{X: 500, Y: 400}, animation.EasingEaseInOutQuad)

				scaleTrack := animation.NewFloatKeyframeTrack().
					AddKeyframe(0, 1.0, animation.EasingLinear).
					AddKeyframe(5*time.Second, 1.5, animation.EasingEaseInQuad)

				clip := timeline.NewClip("c1", "input.mp4", 0, 5*time.Second).
					WithPositionTrack(posTrack).
					WithScaleTrack(scaleTrack)

				tr.AddClip(clip)
				tl.AddTrack(tr)
				return tl
			},
			encodingOpts:       compiler.DefaultEncodingOptions(),
			outputPath:         "out_anim.mp4",
			expectedInputCount: 1,
			expectedArgSnippets: []string{
				"zoompan=z='if(lt(in_time,",
				"x='iw/2-(iw/zoom/2)'",
				"overlay=eval=frame:x='if(lt(t,",
			},
		},
		{
			name: "sequential clips with xfade transition and acrossfade audio",
			buildTimeline: func() *timeline.Timeline {
				tl := timeline.New(timeline.WithCanvas(types.Res1080p), timeline.WithFPS(types.FPS30))
				tr := timeline.NewTrack("v_main", timeline.TrackKindVideo)

				clipA := timeline.NewClip("clip_a", "intro.mp4", 0, 5*time.Second)
				clipB := timeline.NewClip("clip_b", "outro.mp4", 4*time.Second, 5*time.Second)
				transition := timeline.NewTransition("trans1", timeline.TransitionDissolve, 1*time.Second, clipA, clipB)

				tr.AddClip(clipA, clipB)
				tr.AddTransition(transition)
				tl.AddTrack(tr)
				return tl
			},
			encodingOpts:       compiler.DefaultEncodingOptions(),
			outputPath:         "out_trans.mp4",
			expectedInputCount: 2,
			expectedArgSnippets: []string{
				"xfade=transition=dissolve:duration=1.00:offset=4.00",
				"acrossfade=d=1.00:c1=tri:c2=tri",
				"out_trans.mp4",
			},
		},
		{
			name: "picture-in-picture with overlay coordinates and opacity",
			buildTimeline: func() *timeline.Timeline {
				tl := timeline.New(timeline.WithCanvas(types.Res1080p), timeline.WithFPS(types.FPS30))

				mainTrack := timeline.NewTrack("main", timeline.TrackKindVideo).SetZIndex(0)
				mainTrack.AddClip(timeline.NewClip("bg", "bg.mp4", 0, 10*time.Second))

				pipTrack := timeline.NewTrack("pip", timeline.TrackKindOverlay).SetZIndex(1)
				pipTrack.AddClip(timeline.NewClip("fg", "cam.mp4", 1*time.Second, 5*time.Second).
					WithScale(0.3).
					WithPosition(types.Point{X: 100, Y: 100}).
					WithOpacity(0.8))

				tl.AddTrack(mainTrack, pipTrack)
				return tl
			},
			encodingOpts:       compiler.DefaultEncodingOptions(),
			outputPath:         "out_pip.mp4",
			expectedInputCount: 2,
			expectedArgSnippets: []string{
				"-i bg.mp4",
				"-i cam.mp4",
				"overlay=x=100:y=100",
				"between(t,1.0000,6.0000)",
				"out_pip.mp4",
			},
		},
		{
			name: "multi-track audio mixing with delay",
			buildTimeline: func() *timeline.Timeline {
				tl := timeline.New(timeline.WithCanvas(types.Res720p), timeline.WithFPS(types.FPS25))

				vTrack := timeline.NewTrack("v", timeline.TrackKindVideo)
				vTrack.AddClip(timeline.NewClip("v1", "clip.mp4", 0, 10*time.Second))

				aTrack1 := timeline.NewTrack("voice", timeline.TrackKindAudio)
				aTrack1.AddClip(timeline.NewClip("voice_clip", "voice.mp3", 0, 5*time.Second).WithVolume(1.2))

				aTrack2 := timeline.NewTrack("music", timeline.TrackKindAudio)
				aTrack2.AddClip(timeline.NewClip("music_clip", "music.mp3", 2*time.Second, 8*time.Second).WithVolume(0.4))

				tl.AddTrack(vTrack, aTrack1, aTrack2)
				return tl
			},
			encodingOpts:       compiler.DefaultEncodingOptions(),
			outputPath:         "out_audio_mix.mp4",
			expectedInputCount: 3,
			expectedArgSnippets: []string{
				"amix=inputs=3",
				"adelay=delays=2000|2000",
				"out_audio_mix.mp4",
			},
		},
		{
			name: "custom encoding options",
			buildTimeline: func() *timeline.Timeline {
				tl := timeline.New(timeline.WithCanvas(types.Res1080p), timeline.WithFPS(types.FPS60))
				tr := timeline.NewTrack("v0", timeline.TrackKindVideo)
				tr.AddClip(timeline.NewClip("c1", "src.mp4", 0, 3*time.Second))
				tl.AddTrack(tr)
				return tl
			},
			encodingOpts: compiler.EncodingOptions{
				VideoCodec:   "libx265",
				AudioCodec:   "libopus",
				PixelFormat:  "yuv420p10le",
				AudioBitrate: "128k",
				CRF:          18,
				Preset:       "slow",
			},
			outputPath:         "out_hevc.mkv",
			expectedInputCount: 1,
			expectedArgSnippets: []string{
				"-c:v libx265",
				"-c:a libopus",
				"-pix_fmt yuv420p10le",
				"-crf 18",
				"-preset slow",
				"-b:a 128k",
				"out_hevc.mkv",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tl := tt.buildTimeline()
			c := compiler.New(nil).SetEncodingOptions(tt.encodingOpts)

			res, err := c.Compile(tl, tt.outputPath)
			if err != nil {
				t.Fatalf("Compile failed: %v", err)
			}

			if len(res.Inputs) != tt.expectedInputCount {
				t.Errorf("Inputs count = %d, want %d: %v", len(res.Inputs), tt.expectedInputCount, res.Inputs)
			}

			fullCmd := strings.Join(res.Args, " ")
			for _, snippet := range tt.expectedArgSnippets {
				if !strings.Contains(fullCmd, snippet) {
					t.Errorf("expected command to contain %q\nFull command:\n%s", snippet, fullCmd)
				}
			}

			// Verify visualizers
			mermaid, err := res.Mermaid()
			if err != nil || !strings.Contains(mermaid, "graph LR") {
				t.Errorf("Mermaid export failed: %v", err)
			}

			dot, err := res.DOT()
			if err != nil || !strings.Contains(dot, "digraph Filtergraph") {
				t.Errorf("DOT export failed: %v", err)
			}
		})
	}
}
