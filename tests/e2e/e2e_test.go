package e2e_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/farshidrezaei/vidonyx/animation"
	"github.com/farshidrezaei/vidonyx/chromakey"
	"github.com/farshidrezaei/vidonyx/composer"
	"github.com/farshidrezaei/vidonyx/ducking"
	"github.com/farshidrezaei/vidonyx/presets"
	"github.com/farshidrezaei/vidonyx/probe"
	"github.com/farshidrezaei/vidonyx/subtitles"
	"github.com/farshidrezaei/vidonyx/timeline"
	"github.com/farshidrezaei/vidonyx/types"
	"github.com/farshidrezaei/vidonyx/waveform"
)

type syntheticMediaAssets struct {
	videoAlpha        string
	videoBeta         string
	audioVoice        string
	audioSoundtrack   string
	greenScreenSource string
}

func generateSyntheticMedia(t *testing.T, directoryPath string) syntheticMediaAssets {
	t.Helper()

	assets := syntheticMediaAssets{
		videoAlpha:        filepath.Join(directoryPath, "synth_video_alpha.mp4"),
		videoBeta:         filepath.Join(directoryPath, "synth_video_beta.mp4"),
		audioVoice:        filepath.Join(directoryPath, "synth_audio_voice.mp3"),
		audioSoundtrack:   filepath.Join(directoryPath, "synth_audio_music.mp3"),
		greenScreenSource: filepath.Join(directoryPath, "synth_greenscreen.mp4"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	commands := [][]string{
		// 1. Synthetic Video Alpha (testsrc2 pattern + 440Hz sine audio, 3 seconds, 640x360)
		{"-y", "-f", "lavfi", "-i", "testsrc2=size=640x360:rate=30", "-f", "lavfi", "-i", "sine=frequency=440:duration=3", "-t", "3", "-c:v", "libx264", "-c:a", "aac", "-pix_fmt", "yuv420p", assets.videoAlpha},
		// 2. Synthetic Video Beta (smptebars color pattern + 660Hz sine audio, 3 seconds, 640x360)
		{"-y", "-f", "lavfi", "-i", "smptebars=size=640x360:rate=30", "-f", "lavfi", "-i", "sine=frequency=660:duration=3", "-t", "3", "-c:v", "libx264", "-c:a", "aac", "-pix_fmt", "yuv420p", assets.videoBeta},
		// 3. Synthetic Audio Voice (800Hz sine beep, 2 seconds)
		{"-y", "-f", "lavfi", "-i", "sine=frequency=800:duration=2", "-c:a", "libmp3lame", assets.audioVoice},
		// 4. Synthetic Audio Music (300Hz sine tone, 4 seconds)
		{"-y", "-f", "lavfi", "-i", "sine=frequency=300:duration=4", "-c:a", "libmp3lame", assets.audioSoundtrack},
		// 5. Synthetic Green Screen Video (Solid green canvas, 3 seconds)
		{"-y", "-f", "lavfi", "-i", "color=c=0x00FF00:s=640x360:r=30:d=3", "-c:v", "libx264", "-pix_fmt", "yuv420p", assets.greenScreenSource},
	}

	for _, commandArguments := range commands {
		cmd := exec.CommandContext(ctx, "ffmpeg", commandArguments...)
		if outputBytes, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("failed generating synthetic media asset %v: %v\nOutput: %s", commandArguments, err, string(outputBytes))
		}
	}

	return assets
}

func TestEndToEnd_FullCompositionSuiteTable(t *testing.T) {
	temporaryDirectory := t.TempDir()
	assets := generateSyntheticMedia(t, temporaryDirectory)
	mediaProber := probe.NewFFprobeProber("ffprobe")

	tests := []struct {
		name                 string
		buildTimeline        func(assets syntheticMediaAssets) *timeline.Timeline
		configureComposer    func() []composer.Option
		expectedCanvas       types.Size
		expectedMinDuration  time.Duration
		expectedMaxDuration  time.Duration
		verifyAudioPresence  bool
	}{
		{
			name: "01_simple_cut_and_canvas_normalization",
			buildTimeline: func(assets syntheticMediaAssets) *timeline.Timeline {
				compositionTimeline := timeline.New(
					timeline.WithCanvas(types.Res720p),
					timeline.WithFPS(types.FPS30),
				)
				videoTrack := timeline.NewTrack("main_track", timeline.TrackKindVideo)
				videoTrack.AddClip(timeline.NewClip("clip_1", assets.videoAlpha, 0, 3*time.Second))
				compositionTimeline.AddTrack(videoTrack)
				return compositionTimeline
			},
			configureComposer: func() []composer.Option {
				return nil
			},
			expectedCanvas:      types.Res720p,
			expectedMinDuration: 2800 * time.Millisecond,
			expectedMaxDuration: 3200 * time.Millisecond,
			verifyAudioPresence: true,
		},
		{
			name: "02_picture_in_picture_with_opacity_and_positioning",
			buildTimeline: func(assets syntheticMediaAssets) *timeline.Timeline {
				compositionTimeline := timeline.New(
					timeline.WithCanvas(types.Res720p),
					timeline.WithFPS(types.FPS30),
				)
				backgroundTrack := timeline.NewTrack("bg", timeline.TrackKindVideo).SetZIndex(0)
				backgroundTrack.AddClip(timeline.NewClip("bg_clip", assets.videoAlpha, 0, 3*time.Second))

				pipTrack := timeline.NewTrack("pip", timeline.TrackKindOverlay).SetZIndex(1)
				pipClip := timeline.NewClip("pip_clip", assets.videoBeta, 500*time.Millisecond, 2*time.Second).
					WithScale(0.4).
					WithPosition(types.Point{X: 50, Y: 50}).
					WithOpacity(0.85)
				pipTrack.AddClip(pipClip)

				compositionTimeline.AddTrack(backgroundTrack, pipTrack)
				return compositionTimeline
			},
			configureComposer: func() []composer.Option {
				return nil
			},
			expectedCanvas:      types.Res720p,
			expectedMinDuration: 2800 * time.Millisecond,
			expectedMaxDuration: 3200 * time.Millisecond,
			verifyAudioPresence: true,
		},
		{
			name: "03_sequential_transitions_xfade_and_acrossfade",
			buildTimeline: func(assets syntheticMediaAssets) *timeline.Timeline {
				compositionTimeline := timeline.New(
					timeline.WithCanvas(types.Res720p),
					timeline.WithFPS(types.FPS30),
				)
				track := timeline.NewTrack("main_sequence", timeline.TrackKindVideo)
				clipA := timeline.NewClip("clip_a", assets.videoAlpha, 0, 2*time.Second)
				clipB := timeline.NewClip("clip_b", assets.videoBeta, 1500*time.Millisecond, 2*time.Second)
				transition := timeline.NewTransition("dissolve_t", timeline.TransitionDissolve, 500*time.Millisecond, clipA, clipB)

				track.AddClip(clipA, clipB)
				track.AddTransition(transition)
				compositionTimeline.AddTrack(track)
				return compositionTimeline
			},
			configureComposer: func() []composer.Option {
				return nil
			},
			expectedCanvas:      types.Res720p,
			expectedMinDuration: 3300 * time.Millisecond,
			expectedMaxDuration: 3700 * time.Millisecond,
			verifyAudioPresence: true,
		},
		{
			name: "04_animated_motion_ken_burns_and_keyframe_tracks",
			buildTimeline: func(assets syntheticMediaAssets) *timeline.Timeline {
				compositionTimeline := timeline.New(
					timeline.WithCanvas(types.Res720p),
					timeline.WithFPS(types.FPS30),
				)
				videoTrack := timeline.NewTrack("anim_track", timeline.TrackKindVideo)

				positionTrack := animation.NewPositionTrack().
					AddKeyframe(0, types.Point{X: 0, Y: 0}, animation.EasingLinear).
					AddKeyframe(2*time.Second, types.Point{X: 100, Y: 80}, animation.EasingEaseInOutQuad)

				scaleTrack := animation.NewFloatKeyframeTrack().
					AddKeyframe(0, 1.0, animation.EasingLinear).
					AddKeyframe(2*time.Second, 1.2, animation.EasingEaseInQuad)

				clip := timeline.NewClip("motion_clip", assets.videoAlpha, 0, 2*time.Second).
					WithPositionTrack(positionTrack).
					WithScaleTrack(scaleTrack).
					WithFadeIn(500 * time.Millisecond).
					WithFadeOut(500 * time.Millisecond)

				videoTrack.AddClip(clip)
				compositionTimeline.AddTrack(videoTrack)
				return compositionTimeline
			},
			configureComposer: func() []composer.Option {
				return nil
			},
			expectedCanvas:      types.Res720p,
			expectedMinDuration: 1800 * time.Millisecond,
			expectedMaxDuration: 2200 * time.Millisecond,
			verifyAudioPresence: true,
		},
		{
			name: "05_styled_subtitles_burn_in",
			buildTimeline: func(assets syntheticMediaAssets) *timeline.Timeline {
				compositionTimeline := timeline.New(
					timeline.WithCanvas(types.Res720p),
					timeline.WithFPS(types.FPS30),
				)
				videoTrack := timeline.NewTrack("v", timeline.TrackKindVideo)
				videoTrack.AddClip(timeline.NewClip("v1", assets.videoAlpha, 0, 3*time.Second))

				rawSRT := `1
00:00:00,500 --> 00:00:02,500
Vidonyx Real End-To-End Subtitles
`
				subTrack, err := subtitles.ParseSRT(strings.NewReader(rawSRT))
				if err != nil {
					t.Fatalf("failed parsing SRT: %v", err)
				}
				subTrack.SetStyle(subtitles.SubtitleStyle{
					FontSize:     36,
					PrimaryColor: types.ColorYellow,
					Box:          true,
					BoxColor:     types.RGBA(0, 0, 0, 180),
					Alignment:    types.AlignBottomCenter,
					MarginBottom: 40,
				})

				subTimelineTrack := timeline.NewSubtitleTrack("subs", subTrack)
				compositionTimeline.AddTrack(videoTrack, subTimelineTrack)
				return compositionTimeline
			},
			configureComposer: func() []composer.Option {
				return nil
			},
			expectedCanvas:      types.Res720p,
			expectedMinDuration: 2800 * time.Millisecond,
			expectedMaxDuration: 3200 * time.Millisecond,
			verifyAudioPresence: true,
		},
		{
			name: "06_smart_audio_ducking_with_sidechain_compression",
			buildTimeline: func(assets syntheticMediaAssets) *timeline.Timeline {
				compositionTimeline := timeline.New(
					timeline.WithCanvas(types.Res720p),
					timeline.WithFPS(types.FPS30),
				)
				videoTrack := timeline.NewTrack("v", timeline.TrackKindVideo)
				videoTrack.AddClip(timeline.NewClip("v1", assets.videoAlpha, 0, 3*time.Second))

				voiceTrack := timeline.NewTrack("voice_channel", timeline.TrackKindAudio)
				voiceTrack.AddClip(timeline.NewClip("voice_clip", assets.audioVoice, 500*time.Millisecond, 2*time.Second))

				musicTrack := timeline.NewTrack("music_channel", timeline.TrackKindAudio).
					WithDucking("voice_channel", ducking.Options{
						Threshold:           0.1,
						Ratio:               4.0,
						AttackMilliseconds:  20,
						ReleaseMilliseconds: 300,
					})
				musicTrack.AddClip(timeline.NewClip("music_clip", assets.audioSoundtrack, 0, 3*time.Second))

				compositionTimeline.AddTrack(videoTrack, voiceTrack, musicTrack)
				return compositionTimeline
			},
			configureComposer: func() []composer.Option {
				return nil
			},
			expectedCanvas:      types.Res720p,
			expectedMinDuration: 2800 * time.Millisecond,
			expectedMaxDuration: 3200 * time.Millisecond,
			verifyAudioPresence: true,
		},
		{
			name: "07_podcast_audiogram_with_animated_waveform",
			buildTimeline: func(assets syntheticMediaAssets) *timeline.Timeline {
				compositionTimeline := timeline.New(
					timeline.WithCanvas(types.ResSquare1080),
					timeline.WithFPS(types.FPS30),
				)
				bgTrack := timeline.NewTrack("bg", timeline.TrackKindVideo)
				bgTrack.AddClip(timeline.NewClip("bg_clip", assets.videoAlpha, 0, 2*time.Second))

				audioTrack := timeline.NewTrack("speech", timeline.TrackKindAudio)
				audioTrack.AddClip(timeline.NewClip("speech_clip", assets.audioVoice, 0, 2*time.Second))

				waveTrack := timeline.NewWaveformTrack("wave", "speech", waveform.Options{
					Size:     types.NewSize(800, 150),
					Mode:     waveform.ModePeakToPeak,
					Color:    types.RGB(0, 240, 180),
					Scale:    "sqrt",
					Position: types.Point{X: 140, Y: 800},
				})

				compositionTimeline.AddTrack(bgTrack, audioTrack, waveTrack)
				return compositionTimeline
			},
			configureComposer: func() []composer.Option {
				return nil
			},
			expectedCanvas:      types.ResSquare1080,
			expectedMinDuration: 1800 * time.Millisecond,
			expectedMaxDuration: 2200 * time.Millisecond,
			verifyAudioPresence: true,
		},
		{
			name: "08_studio_greenscreen_chromakey_and_despill",
			buildTimeline: func(assets syntheticMediaAssets) *timeline.Timeline {
				compositionTimeline := timeline.New(
					timeline.WithCanvas(types.Res720p),
					timeline.WithFPS(types.FPS30),
				)
				bgTrack := timeline.NewTrack("bg", timeline.TrackKindVideo).SetZIndex(0)
				bgTrack.AddClip(timeline.NewClip("bg_clip", assets.videoAlpha, 0, 2*time.Second))

				greenTrack := timeline.NewTrack("green_presenter", timeline.TrackKindOverlay).SetZIndex(1)
				greenClip := timeline.NewClip("host", assets.greenScreenSource, 0, 2*time.Second).
					WithChromaKey(chromakey.Options{
						KeyColor:    chromakey.StudioGreenScreen,
						Similarity:  0.30,
						Blend:       0.10,
						Despill:     true,
						DespillType: chromakey.DespillGreen,
					})
				greenTrack.AddClip(greenClip)

				compositionTimeline.AddTrack(bgTrack, greenTrack)
				return compositionTimeline
			},
			configureComposer: func() []composer.Option {
				return nil
			},
			expectedCanvas:      types.Res720p,
			expectedMinDuration: 1800 * time.Millisecond,
			expectedMaxDuration: 2200 * time.Millisecond,
			verifyAudioPresence: true,
		},
		{
			name: "09_tiktok_vertical_preset_rendering",
			buildTimeline: func(assets syntheticMediaAssets) *timeline.Timeline {
				preset := presets.TikTokVertical1080p60()
				compositionTimeline := timeline.New()
				preset.ApplyToTimeline(compositionTimeline)

				videoTrack := timeline.NewTrack("v", timeline.TrackKindVideo)
				videoTrack.AddClip(timeline.NewClip("v1", assets.videoAlpha, 0, 2*time.Second))
				compositionTimeline.AddTrack(videoTrack)
				return compositionTimeline
			},
			configureComposer: func() []composer.Option {
				return []composer.Option{
					composer.WithPreset(presets.TikTokVertical1080p60()),
				}
			},
			expectedCanvas:      types.ResPortrait1080p,
			expectedMinDuration: 1800 * time.Millisecond,
			expectedMaxDuration: 2200 * time.Millisecond,
			verifyAudioPresence: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputFilePath := filepath.Join(temporaryDirectory, tt.name+".mp4")
			compositionTimeline := tt.buildTimeline(assets)

			composerOptions := tt.configureComposer()
			composerInstance := composer.New(composerOptions...)

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			renderResult, err := composerInstance.Render(ctx, compositionTimeline, outputFilePath, nil)
			if err != nil {
				t.Fatalf("Render failed for scenario %q: %v", tt.name, err)
			}

			if renderResult.OutputPath != outputFilePath {
				t.Errorf("OutputPath = %q, want %q", renderResult.OutputPath, outputFilePath)
			}

			// Verify file existence and non-empty byte size
			fileInfo, err := os.Stat(outputFilePath)
			if err != nil {
				t.Fatalf("Output file does not exist: %v", err)
			}
			if fileInfo.Size() < 5000 {
				t.Errorf("Rendered file size %d bytes is suspiciously small", fileInfo.Size())
			}

			// Probe the rendered output video with ffprobe
			metadata, err := mediaProber.Probe(ctx, outputFilePath)
			if err != nil {
				t.Fatalf("ffprobe failed on rendered output %q: %v", outputFilePath, err)
			}

			if metadata.VideoStream == nil {
				t.Fatalf("Rendered file is missing video stream")
			}

			if metadata.VideoStream.Width != tt.expectedCanvas.Width || metadata.VideoStream.Height != tt.expectedCanvas.Height {
				t.Errorf("Rendered Canvas = %dx%d, want %dx%d", metadata.VideoStream.Width, metadata.VideoStream.Height, tt.expectedCanvas.Width, tt.expectedCanvas.Height)
			}

			if metadata.Duration < tt.expectedMinDuration || metadata.Duration > tt.expectedMaxDuration {
				t.Errorf("Rendered Duration = %v, want within [%v, %v]", metadata.Duration, tt.expectedMinDuration, tt.expectedMaxDuration)
			}

			if tt.verifyAudioPresence && metadata.AudioStream == nil {
				t.Errorf("Rendered file was expected to contain an audio stream, but none was found")
			}
		})
	}
}
