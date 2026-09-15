package presets_test

import (
	"testing"

	"github.com/farshidrezaei/vidonex/presets"
)

func TestCodecsProfiles_Table(t *testing.T) {
	testCases := []struct {
		name              string
		profileName       string
		videoCodec        string
		audioCodec        string
		pixelFormat       string
		expectedCRF       int
		expectAlphaFormat bool
	}{
		{
			name:              "AV1 10-bit profile",
			profileName:       "AV1",
			videoCodec:        "libsvtav1",
			audioCodec:        "libopus",
			pixelFormat:       "yuv420p10le",
			expectedCRF:       28,
			expectAlphaFormat: false,
		},
		{
			name:              "ProRes 4444 with Alpha Channel",
			profileName:       "ProRes4444",
			videoCodec:        "prores_ks",
			audioCodec:        "pcm_s24le",
			pixelFormat:       "yuva444p10le",
			expectedCRF:       0,
			expectAlphaFormat: true,
		},
		{
			name:              "ProRes 422 HQ Broadcast Mastering",
			profileName:       "ProRes422HQ",
			videoCodec:        "prores_ks",
			audioCodec:        "pcm_s24le",
			pixelFormat:       "yuv422p10le",
			expectedCRF:       0,
			expectAlphaFormat: false,
		},
		{
			name:              "WebM VP9 Standard",
			profileName:       "WebMVP9",
			videoCodec:        "libvpx-vp9",
			audioCodec:        "libopus",
			pixelFormat:       "yuv420p",
			expectedCRF:       30,
			expectAlphaFormat: false,
		},
		{
			name:              "WebM VP9 Transparent Alpha",
			profileName:       "WebMAlpha",
			videoCodec:        "libvpx-vp9",
			audioCodec:        "libopus",
			pixelFormat:       "yuva420p",
			expectedCRF:       28,
			expectAlphaFormat: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var opts struct {
				VideoCodec  string
				AudioCodec  string
				PixelFormat string
				CRF         int
			}

			switch tc.profileName {
			case "AV1":
				p := presets.ProfileAV1(0, 0)
				opts.VideoCodec = p.VideoCodec
				opts.AudioCodec = p.AudioCodec
				opts.PixelFormat = p.PixelFormat
				opts.CRF = p.CRF
			case "ProRes4444":
				p := presets.ProfileProRes4444()
				opts.VideoCodec = p.VideoCodec
				opts.AudioCodec = p.AudioCodec
				opts.PixelFormat = p.PixelFormat
				opts.CRF = p.CRF
			case "ProRes422HQ":
				p := presets.ProfileProRes422HQ()
				opts.VideoCodec = p.VideoCodec
				opts.AudioCodec = p.AudioCodec
				opts.PixelFormat = p.PixelFormat
				opts.CRF = p.CRF
			case "WebMVP9":
				p := presets.ProfileWebMVP9(0)
				opts.VideoCodec = p.VideoCodec
				opts.AudioCodec = p.AudioCodec
				opts.PixelFormat = p.PixelFormat
				opts.CRF = p.CRF
			case "WebMAlpha":
				p := presets.ProfileWebMAlpha(0)
				opts.VideoCodec = p.VideoCodec
				opts.AudioCodec = p.AudioCodec
				opts.PixelFormat = p.PixelFormat
				opts.CRF = p.CRF
			}

			if opts.VideoCodec != tc.videoCodec {
				t.Errorf("expected VideoCodec %s, got %s", tc.videoCodec, opts.VideoCodec)
			}
			if opts.AudioCodec != tc.audioCodec {
				t.Errorf("expected AudioCodec %s, got %s", tc.audioCodec, opts.AudioCodec)
			}
			if opts.PixelFormat != tc.pixelFormat {
				t.Errorf("expected PixelFormat %s, got %s", tc.pixelFormat, opts.PixelFormat)
			}
			if opts.CRF != tc.expectedCRF {
				t.Errorf("expected CRF %d, got %d", tc.expectedCRF, opts.CRF)
			}
		})
	}
}
