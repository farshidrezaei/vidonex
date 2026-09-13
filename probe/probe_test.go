package probe_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/farshidrezaei/vidonex/probe"
	"github.com/farshidrezaei/vidonex/types"
)

func TestParseFFprobeJSONTable(t *testing.T) {
	sampleFullMediaJSON := `{
  "streams": [
    {
      "codec_type": "video",
      "codec_name": "h264",
      "width": 1920,
      "height": 1080,
      "pix_fmt": "yuv420p",
      "avg_frame_rate": "30/1",
      "duration": "12.500000",
      "nb_frames": "375"
    },
    {
      "codec_type": "audio",
      "codec_name": "aac",
      "sample_rate": "48000",
      "channels": 2,
      "channel_layout": "stereo",
      "bit_rate": "192000",
      "duration": "12.500000"
    }
  ],
  "format": {
    "filename": "sample.mp4",
    "format_name": "mov,mp4,m4a,3gp,3g2,mj2",
    "duration": "12.500000",
    "size": "5242880",
    "bit_rate": "3355443"
  }
}`

	sampleVideoOnlyJSON := `{
  "streams": [
    {
      "codec_type": "video",
      "codec_name": "vp9",
      "width": 1280,
      "height": 720,
      "pix_fmt": "yuv420p",
      "avg_frame_rate": "60/1",
      "duration": "5.000000"
    }
  ],
  "format": {
    "format_name": "matroska,webm",
    "duration": "5.000000"
  }
}`

	sampleAudioOnlyJSON := `{
  "streams": [
    {
      "codec_type": "audio",
      "codec_name": "mp3",
      "sample_rate": "44100",
      "channels": 2,
      "channel_layout": "stereo",
      "duration": "30.000000"
    }
  ],
  "format": {
    "format_name": "mp3",
    "duration": "30.000000"
  }
}`

	sampleNoStreamsJSON := `{
  "streams": [],
  "format": {
    "format_name": "data"
  }
}`

	tests := []struct {
		name          string
		jsonData      string
		mediaPath     string
		shouldErr     bool
		errContains   string
		checkMetadata func(t *testing.T, meta *probe.MediaMetadata)
	}{
		{
			name:      "full video and audio container",
			jsonData:  sampleFullMediaJSON,
			mediaPath: "sample.mp4",
			shouldErr: false,
			checkMetadata: func(t *testing.T, meta *probe.MediaMetadata) {
				if !meta.HasVideo() || !meta.HasAudio() {
					t.Fatalf("expected both video and audio streams")
				}
				if meta.VideoStream.Width != 1920 || meta.VideoStream.Height != 1080 {
					t.Errorf("Video dimensions = %dx%d, want 1920x1080", meta.VideoStream.Width, meta.VideoStream.Height)
				}
				if !meta.VideoStream.FPS.Equal(types.FPS30) {
					t.Errorf("Video FPS = %v, want 30", meta.VideoStream.FPS)
				}
				if meta.AudioStream.SampleRate != 48000 || meta.AudioStream.Channels != 2 {
					t.Errorf("Audio = %dHz %d channels, want 48000Hz 2 channels", meta.AudioStream.SampleRate, meta.AudioStream.Channels)
				}
				if meta.Duration != 12500*time.Millisecond {
					t.Errorf("Duration = %v, want 12.5s", meta.Duration)
				}
			},
		},
		{
			name:      "video only container",
			jsonData:  sampleVideoOnlyJSON,
			mediaPath: "video.webm",
			shouldErr: false,
			checkMetadata: func(t *testing.T, meta *probe.MediaMetadata) {
				if !meta.HasVideo() || meta.HasAudio() {
					t.Fatalf("expected video-only stream")
				}
				if meta.VideoStream.Codec != "vp9" {
					t.Errorf("Video Codec = %s, want vp9", meta.VideoStream.Codec)
				}
			},
		},
		{
			name:      "audio only container",
			jsonData:  sampleAudioOnlyJSON,
			mediaPath: "music.mp3",
			shouldErr: false,
			checkMetadata: func(t *testing.T, meta *probe.MediaMetadata) {
				if meta.HasVideo() || !meta.HasAudio() {
					t.Fatalf("expected audio-only stream")
				}
				if meta.AudioStream.Codec != "mp3" {
					t.Errorf("Audio Codec = %s, want mp3", meta.AudioStream.Codec)
				}
			},
		},
		{
			name:        "no media streams error",
			jsonData:    sampleNoStreamsJSON,
			mediaPath:   "empty.dat",
			shouldErr:   true,
			errContains: "no valid media streams",
		},
		{
			name:        "malformed json error",
			jsonData:    "{invalid_json}",
			mediaPath:   "corrupt.mp4",
			shouldErr:   true,
			errContains: "failed to parse ffprobe json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata, err := probe.ParseFFprobeJSON([]byte(tt.jsonData), tt.mediaPath)
			if tt.shouldErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("expected error containing %q, got: %v", tt.errContains, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.checkMetadata != nil {
				tt.checkMetadata(t, metadata)
			}
		})
	}
}

func TestCachedProberTable(t *testing.T) {
	mockProber := probe.NewMockProber()
	cachedProber := probe.NewCachedProber(mockProber)
	ctx := context.Background()

	// 1. First call triggers underlying prober
	metadataFirst, err := cachedProber.Probe(ctx, "video.mp4")
	if err != nil {
		t.Fatalf("first probe failed: %v", err)
	}

	// 2. Second call should return cached instance
	metadataSecond, err := cachedProber.Probe(ctx, "video.mp4")
	if err != nil {
		t.Fatalf("second probe failed: %v", err)
	}

	if metadataFirst != metadataSecond {
		t.Errorf("expected cached pointer instance to match")
	}

	// Underlying mock should only have recorded 1 invocation for "video.mp4"
	if len(mockProber.ProbedPaths()) != 1 {
		t.Errorf("underlying prober calls = %d, want 1", len(mockProber.ProbedPaths()))
	}
}

func TestMockProberTable(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func() *probe.MockProber
		mediaPath     string
		shouldErr     bool
		checkMetadata func(t *testing.T, meta *probe.MediaMetadata)
	}{
		{
			name: "default fallback metadata",
			setupMock: func() *probe.MockProber {
				return probe.NewMockProber()
			},
			mediaPath: "clip.mp4",
			shouldErr: false,
			checkMetadata: func(t *testing.T, meta *probe.MediaMetadata) {
				if meta.Path != "clip.mp4" {
					t.Errorf("Path = %s, want clip.mp4", meta.Path)
				}
				if meta.VideoStream.Width != 1920 {
					t.Errorf("Width = %d, want 1920", meta.VideoStream.Width)
				}
			},
		},
		{
			name: "custom registered metadata",
			setupMock: func() *probe.MockProber {
				m := probe.NewMockProber()
				m.SetMetadata("custom_4k.mp4", &probe.MediaMetadata{
					Path: "custom_4k.mp4",
					VideoStream: &probe.VideoStreamMetadata{
						Width:  3840,
						Height: 2160,
						FPS:    types.FPS60,
					},
				})
				return m
			},
			mediaPath: "custom_4k.mp4",
			shouldErr: false,
			checkMetadata: func(t *testing.T, meta *probe.MediaMetadata) {
				if meta.VideoStream.Width != 3840 || meta.VideoStream.Height != 2160 {
					t.Errorf("Width x Height = %dx%d, want 3840x2160", meta.VideoStream.Width, meta.VideoStream.Height)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.setupMock()
			metadata, err := mock.Probe(context.Background(), tt.mediaPath)
			if tt.shouldErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.checkMetadata != nil {
				tt.checkMetadata(t, metadata)
			}
		})
	}
}
