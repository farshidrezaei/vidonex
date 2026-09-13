package presets_test

import (
	"testing"

	"github.com/farshidrezaei/vidonex/compiler"
	"github.com/farshidrezaei/vidonex/presets"
	"github.com/farshidrezaei/vidonex/timeline"
	"github.com/farshidrezaei/vidonex/types"
)

func TestPlatformPresetsTable(t *testing.T) {
	tests := []struct {
		name         string
		preset       presets.PlatformPreset
		wantWidth    int
		wantHeight   int
		wantFPS      types.Rational
		wantCRF      int
		wantAudioBit string
	}{
		{
			name:         "TikTok 1080p 60fps",
			preset:       presets.TikTokVertical1080p60(),
			wantWidth:    1080,
			wantHeight:   1920,
			wantFPS:      types.FPS60,
			wantCRF:      21,
			wantAudioBit: "192k",
		},
		{
			name:         "Instagram Square 1080p",
			preset:       presets.InstagramSquare1080p(),
			wantWidth:    1080,
			wantHeight:   1080,
			wantFPS:      types.FPS30,
			wantCRF:      22,
			wantAudioBit: "192k",
		},
		{
			name:         "YouTube 4K 60fps",
			preset:       presets.YouTube4K60(),
			wantWidth:    3840,
			wantHeight:   2160,
			wantFPS:      types.FPS60,
			wantCRF:      18,
			wantAudioBit: "320k",
		},
		{
			name:         "Web Fast 720p 30fps",
			preset:       presets.WebFast720p30(),
			wantWidth:    1280,
			wantHeight:   720,
			wantFPS:      types.FPS30,
			wantCRF:      24,
			wantAudioBit: "128k",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.preset.Canvas.Width != tt.wantWidth || tt.preset.Canvas.Height != tt.wantHeight {
				t.Errorf("Canvas = %dx%d, want %dx%d", tt.preset.Canvas.Width, tt.preset.Canvas.Height, tt.wantWidth, tt.wantHeight)
			}
			if !tt.preset.FPS.Equal(tt.wantFPS) {
				t.Errorf("FPS = %v, want %v", tt.preset.FPS, tt.wantFPS)
			}
			if tt.preset.Encoding.CRF != tt.wantCRF {
				t.Errorf("CRF = %d, want %d", tt.preset.Encoding.CRF, tt.wantCRF)
			}
			if tt.preset.Encoding.AudioBitrate != tt.wantAudioBit {
				t.Errorf("AudioBitrate = %s, want %s", tt.preset.Encoding.AudioBitrate, tt.wantAudioBit)
			}

			// Verify ApplyToTimeline
			tl := timeline.New()
			tt.preset.ApplyToTimeline(tl)
			if tl.Canvas != tt.preset.Canvas {
				t.Errorf("ApplyToTimeline canvas mismatch")
			}
			if !tl.FPS.Equal(tt.preset.FPS) {
				t.Errorf("ApplyToTimeline FPS mismatch")
			}
		})
	}
}

func TestConfigureHardwareEncodingTable(t *testing.T) {
	tests := []struct {
		name          string
		accelerator   presets.HardwareAccelerator
		useHEVC       bool
		wantCodec     string
		wantPresetSet bool
	}{
		{
			name:          "NVENC h264",
			accelerator:   presets.AcceleratorNVENC,
			useHEVC:       false,
			wantCodec:     "h264_nvenc",
			wantPresetSet: true,
		},
		{
			name:          "NVENC hevc",
			accelerator:   presets.AcceleratorNVENC,
			useHEVC:       true,
			wantCodec:     "hevc_nvenc",
			wantPresetSet: true,
		},
		{
			name:          "VideoToolbox h264",
			accelerator:   presets.AcceleratorVideoToolbox,
			useHEVC:       false,
			wantCodec:     "h264_videotoolbox",
			wantPresetSet: false,
		},
		{
			name:          "VideoToolbox hevc",
			accelerator:   presets.AcceleratorVideoToolbox,
			useHEVC:       true,
			wantCodec:     "hevc_videotoolbox",
			wantPresetSet: false,
		},
		{
			name:          "QSV h264",
			accelerator:   presets.AcceleratorQSV,
			useHEVC:       false,
			wantCodec:     "h264_qsv",
			wantPresetSet: false,
		},
		{
			name:          "VAAPI hevc",
			accelerator:   presets.AcceleratorVAAPI,
			useHEVC:       true,
			wantCodec:     "hevc_vaapi",
			wantPresetSet: false,
		},
		{
			name:          "CPU fallback",
			accelerator:   presets.AcceleratorNone,
			useHEVC:       false,
			wantCodec:     "libx264",
			wantPresetSet: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := compiler.DefaultEncodingOptions()
			presets.ConfigureHardwareEncoding(&opts, tt.accelerator, tt.useHEVC)

			if opts.VideoCodec != tt.wantCodec {
				t.Errorf("VideoCodec = %s, want %s", opts.VideoCodec, tt.wantCodec)
			}
		})
	}
}
