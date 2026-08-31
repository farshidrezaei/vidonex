package spec_test

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/farshidrezaei/vidonyx/animation"
	"github.com/farshidrezaei/vidonyx/spec"
	"github.com/farshidrezaei/vidonyx/timeline"
	"github.com/farshidrezaei/vidonyx/types"
)

func TestSpec_ParserAndConverterTable(t *testing.T) {
	tests := []struct {
		name                 string
		rawSpec              string
		isYAML               bool
		expectedCanvas       types.Size
		expectedFPS          types.Rational
		expectedTrackCount   int
		expectedOutputPath   string
		expectedVideoCodec   string
		verifyTimelineValues func(t *testing.T, tl *timeline.Timeline)
	}{
		{
			name: "yaml_vlog_with_ducking_subtitles_and_waveform",
			isYAML: true,
			rawSpec: `
version: "1.0"
canvas:
  preset: "1080p"
  background_color: "black"
fps: "30"
output: "vlog_final.mp4"
tracks:
  - id: "video_track"
    kind: "video"
    clips:
      - id: "clip_intro"
        source: "./media/intro.mp4"
        start: "0s"
        duration: "5s"
        fade_in: "500ms"
        scale: 1.0
  - id: "voice_track"
    kind: "audio"
    clips:
      - id: "voice_clip"
        source: "./audio/voice.mp3"
        start: "1s"
        duration: "4s"
  - id: "bgm_track"
    kind: "audio"
    duck_under: "voice_track"
    ducking:
      threshold: 0.12
      ratio: 4.5
      attack_ms: 25
      release_ms: 350
    clips:
      - id: "music_clip"
        source: "./audio/music.mp3"
        start: "0s"
        duration: "5s"
        volume: 0.8
  - id: "subtitles_track"
    kind: "subtitle"
    subtitles:
      srt_content: |
        1
        00:00:01,000 --> 00:00:03,500
        Welcome to Vidonyx Declarative Engine!
      style:
        font_size: 42
        color: "yellow"
        box: true
        box_color: "black"
        alignment: "bottom-center"
        margin_bottom: 50
  - id: "wave_track"
    kind: "waveform"
    waveform:
      source_audio: "voice_track"
      mode: "p2p"
      color: "#00F0B4"
      scale: "sqrt"
      width: 800
      height: 120
      position:
        x: 140
        y: 850
`,
			expectedCanvas:     types.Res1080p,
			expectedFPS:        types.FPS30,
			expectedTrackCount: 5,
			expectedOutputPath: "vlog_final.mp4",
			expectedVideoCodec: "libx264",
			verifyTimelineValues: func(t *testing.T, tl *timeline.Timeline) {
				t.Helper()
				if tl.Tracks[2].DucksUnderTrackID != "voice_track" {
					t.Errorf("DucksUnderTrackID = %q, want %q", tl.Tracks[2].DucksUnderTrackID, "voice_track")
				}
				if tl.Tracks[2].DuckingOptions.Ratio != 4.5 {
					t.Errorf("Ducking Ratio = %v, want 4.5", tl.Tracks[2].DuckingOptions.Ratio)
				}
			},
		},
		{
			name: "json_tiktok_preset_with_chromakey_and_keyframes",
			isYAML: false,
			rawSpec: `{
  "version": "1.0",
  "preset": "tiktok_1080p60",
  "output": "tiktok_out.mp4",
  "hardware_acceleration": "none",
  "tracks": [
    {
      "id": "background",
      "kind": "video",
      "clips": [
        {
          "id": "bg_clip",
          "source": "assets/bg.mp4",
          "start": 0,
          "duration": 3.0
        }
      ]
    },
    {
      "id": "host_overlay",
      "kind": "overlay",
      "z_index": 1,
      "clips": [
        {
          "id": "green_host",
          "source": "assets/green.mp4",
          "start": "0.5s",
          "duration": 2.5,
          "chroma_key": {
            "key_color": "#00FF00",
            "similarity": 0.32,
            "blend": 0.12,
            "despill": true,
            "despill_type": "green"
          },
          "keyframes": {
            "position": [
              { "time": "0s", "x": 100, "y": 200, "easing": "linear" },
              { "time": "2s", "x": 150, "y": 350, "easing": "ease_in_out_quad" }
            ],
            "scale": [
              { "time": 0, "value": 1.0, "easing": "linear" },
              { "time": 2.0, "value": 1.2, "easing": "ease_in_quad" }
            ]
          }
        }
      ]
    }
  ]
}`,
			expectedCanvas:     types.ResPortrait1080p,
			expectedFPS:        types.FPS60,
			expectedTrackCount: 2,
			expectedOutputPath: "tiktok_out.mp4",
			expectedVideoCodec: "libx264",
			verifyTimelineValues: func(t *testing.T, tl *timeline.Timeline) {
				t.Helper()
				hostClip := tl.Tracks[1].Clips[0]
				if hostClip.ChromaKeyOptions == nil {
					t.Fatalf("hostClip.ChromaKeyOptions is nil")
				}
				if hostClip.ChromaKeyOptions.Similarity != 0.32 {
					t.Errorf("ChromaKey similarity = %v, want 0.32", hostClip.ChromaKeyOptions.Similarity)
				}
				if hostClip.PositionTrack == nil {
					t.Fatalf("hostClip.PositionTrack is nil")
				}
				if hostClip.ScaleTrack == nil {
					t.Fatalf("hostClip.ScaleTrack is nil")
				}
			},
		},
		{
			name: "yaml_transitions_between_clips",
			isYAML: true,
			rawSpec: `
version: "1.0"
canvas:
  width: 1280
  height: 720
fps: 24
output: "transition_test.mp4"
tracks:
  - id: "main"
    kind: "video"
    clips:
      - id: "clip_a"
        source: "a.mp4"
        start: 0
        duration: 3.0
      - id: "clip_b"
        source: "b.mp4"
        start: 2.5
        duration: 3.0
    transitions:
      - id: "trans_ab"
        type: "dissolve"
        duration: "0.5s"
        from: "clip_a"
        to: "clip_b"
`,
			expectedCanvas:     types.Res720p,
			expectedFPS:        types.FPS24,
			expectedTrackCount: 1,
			expectedOutputPath: "transition_test.mp4",
			expectedVideoCodec: "libx264",
			verifyTimelineValues: func(t *testing.T, tl *timeline.Timeline) {
				t.Helper()
				if len(tl.Tracks[0].Transitions) != 1 {
					t.Fatalf("Transitions count = %d, want 1", len(tl.Tracks[0].Transitions))
				}
				trans := tl.Tracks[0].Transitions[0]
				if trans.Type != timeline.TransitionDissolve {
					t.Errorf("Transition type = %v, want dissolve", trans.Type)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var parsedSpec *spec.VideoSpec
			var err error

			if tt.isYAML {
				parsedSpec, err = spec.ParseYAML(strings.NewReader(tt.rawSpec))
			} else {
				parsedSpec, err = spec.ParseJSON(strings.NewReader(tt.rawSpec))
			}

			if err != nil {
				t.Fatalf("Failed parsing spec: %v", err)
			}

			timelineAST, encodingOptions, outputPath, err := spec.ToTimeline(parsedSpec)
			if err != nil {
				t.Fatalf("Failed converting spec to timeline: %v", err)
			}

			if timelineAST.Canvas.Width != tt.expectedCanvas.Width || timelineAST.Canvas.Height != tt.expectedCanvas.Height {
				t.Errorf("Canvas = %dx%d, want %dx%d", timelineAST.Canvas.Width, timelineAST.Canvas.Height, tt.expectedCanvas.Width, tt.expectedCanvas.Height)
			}

			if timelineAST.FPS.Num != tt.expectedFPS.Num || timelineAST.FPS.Den != tt.expectedFPS.Den {
				t.Errorf("FPS = %v, want %v", timelineAST.FPS, tt.expectedFPS)
			}

			if len(timelineAST.Tracks) != tt.expectedTrackCount {
				t.Errorf("Tracks count = %d, want %d", len(timelineAST.Tracks), tt.expectedTrackCount)
			}

			if outputPath != tt.expectedOutputPath {
				t.Errorf("OutputPath = %q, want %q", outputPath, tt.expectedOutputPath)
			}

			if encodingOptions.VideoCodec != tt.expectedVideoCodec {
				t.Errorf("VideoCodec = %q, want %q", encodingOptions.VideoCodec, tt.expectedVideoCodec)
			}

			if tt.verifyTimelineValues != nil {
				tt.verifyTimelineValues(t, timelineAST)
			}
		})
	}
}

func TestSpec_RelativePathResolution(t *testing.T) {
	rawYAML := `
version: "1.0"
output: "rendered/output.mp4"
tracks:
  - id: "v"
    kind: "video"
    clips:
      - id: "c1"
        source: "./media/intro.mp4"
        duration: "3s"
      - id: "c2"
        source: "https://example.com/stream.m3u8"
        start: "3s"
        duration: "3s"
  - id: "s"
    kind: "subtitle"
    subtitles:
      file: "captions/english.srt"
`
	specification, err := spec.ParseYAML(strings.NewReader(rawYAML))
	if err != nil {
		t.Fatalf("Failed parsing YAML: %v", err)
	}

	baseDir := "/home/user/project"
	spec.ResolveAssetPaths(specification, baseDir)

	if specification.Output != filepath.Join(baseDir, "rendered/output.mp4") {
		t.Errorf("Output path = %q, want %q", specification.Output, filepath.Join(baseDir, "rendered/output.mp4"))
	}

	clip1Source := specification.Tracks[0].Clips[0].Source
	expectedClip1 := filepath.Join(baseDir, "media/intro.mp4")
	if clip1Source != expectedClip1 {
		t.Errorf("Clip 1 Source = %q, want %q", clip1Source, expectedClip1)
	}

	// URL should remain unchanged
	clip2Source := specification.Tracks[0].Clips[1].Source
	if clip2Source != "https://example.com/stream.m3u8" {
		t.Errorf("Clip 2 URL = %q, want unchanged URL", clip2Source)
	}

	subtitleFile := specification.Tracks[1].Subtitles.File
	expectedSub := filepath.Join(baseDir, "captions/english.srt")
	if subtitleFile != expectedSub {
		t.Errorf("Subtitle File = %q, want %q", subtitleFile, expectedSub)
	}
}

func TestSpec_ParseHelpersTable(t *testing.T) {
	t.Run("durations", func(t *testing.T) {
		tests := []struct {
			input    any
			expected time.Duration
		}{
			{input: 5, expected: 5 * time.Second},
			{input: 2.5, expected: 2500 * time.Millisecond},
			{input: "3s", expected: 3 * time.Second},
			{input: "500ms", expected: 500 * time.Millisecond},
			{input: "00:01:30.500", expected: 90500 * time.Millisecond},
			{input: "02:15", expected: 135 * time.Second},
		}

		for _, tt := range tests {
			dur, err := spec.ParseDurationValue(tt.input)
			if err != nil {
				t.Fatalf("ParseDurationValue(%v) error: %v", tt.input, err)
			}
			if dur != tt.expected {
				t.Errorf("ParseDurationValue(%v) = %v, want %v", tt.input, dur, tt.expected)
			}
		}
	})

	t.Run("colors", func(t *testing.T) {
		tests := []struct {
			input    string
			expected types.Color
		}{
			{input: "red", expected: types.ColorRed},
			{input: "green", expected: types.ColorGreen},
			{input: "blue", expected: types.ColorBlue},
			{input: "yellow", expected: types.ColorYellow},
			{input: "black", expected: types.ColorBlack},
			{input: "white", expected: types.ColorWhite},
			{input: "#FF0000", expected: types.RGB(255, 0, 0)},
		}

		for _, tt := range tests {
			color, err := spec.ParseColorValue(tt.input)
			if err != nil {
				t.Fatalf("ParseColorValue(%q) error: %v", tt.input, err)
			}
			if color != tt.expected {
				t.Errorf("ParseColorValue(%q) = %v, want %v", tt.input, color, tt.expected)
			}
		}
	})

	t.Run("easings", func(t *testing.T) {
		if spec.ParseEasingValue("linear") != animation.EasingLinear {
			t.Errorf("linear failed")
		}
		if spec.ParseEasingValue("ease_in_quad") != animation.EasingEaseInQuad {
			t.Errorf("ease_in_quad failed")
		}
		if spec.ParseEasingValue("ease_in_out_sine") != animation.EasingEaseInOutSine {
			t.Errorf("ease_in_out_sine failed")
		}
	})

	t.Run("alignments", func(t *testing.T) {
		if spec.ParseAlignmentValue("center") != types.AlignCenter {
			t.Errorf("center failed")
		}
		if spec.ParseAlignmentValue("bottom_center") != types.AlignBottomCenter {
			t.Errorf("bottom_center failed")
		}
		if spec.ParseAlignmentValue("top_left") != types.AlignTopLeft {
			t.Errorf("top_left failed")
		}
	})
}
