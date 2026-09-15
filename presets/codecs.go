// Package presets provides platform templates, hardware acceleration flags, and turnkey encoding profiles.
package presets

import (
	"fmt"

	"github.com/farshidrezaei/vidonex/compiler"
)

// ProfileAV1 returns an optimized SVT-AV1 encoding configuration.
// AV1 offers 30-40% higher compression efficiency than H.264/HEVC with 10-bit color depth.
func ProfileAV1(qualityPreset, constantRateFactor int) compiler.EncodingOptions {
	crf := constantRateFactor
	if crf <= 0 {
		crf = 28
	}

	presetNum := qualityPreset
	if presetNum <= 0 {
		presetNum = 7
	}

	return compiler.EncodingOptions{
		VideoCodec:   "libsvtav1",
		AudioCodec:   "libopus",
		PixelFormat:  "yuv420p10le",
		AudioBitrate: "128k",
		CRF:          crf,
		Preset:       fmt.Sprintf("%d", presetNum),
		CustomArgs: []string{
			"-svtav1-params", "tune=0",
		},
	}
}

// ProfileProRes4444 returns an Apple ProRes 4444 configuration with 10-bit color and alpha channel preservation.
// Essential for professional master exports and visual effects compositing with transparency.
func ProfileProRes4444() compiler.EncodingOptions {
	return compiler.EncodingOptions{
		VideoCodec:   "prores_ks",
		AudioCodec:   "pcm_s24le",
		PixelFormat:  "yuva444p10le",
		AudioBitrate: "1536k",
		CRF:          0,
		Preset:       "",
		CustomArgs: []string{
			"-profile:v", "4",
			"-bits_per_mb", "8000",
		},
	}
}

// ProfileProRes422HQ returns an Apple ProRes 422 HQ broadcast mastering configuration.
func ProfileProRes422HQ() compiler.EncodingOptions {
	return compiler.EncodingOptions{
		VideoCodec:   "prores_ks",
		AudioCodec:   "pcm_s24le",
		PixelFormat:  "yuv422p10le",
		AudioBitrate: "1536k",
		CRF:          0,
		Preset:       "",
		CustomArgs: []string{
			"-profile:v", "3",
		},
	}
}

// ProfileWebMVP9 returns a WebM VP9 + Opus encoding configuration for web and open streaming.
func ProfileWebMVP9(constantRateFactor int) compiler.EncodingOptions {
	crf := constantRateFactor
	if crf <= 0 {
		crf = 30
	}

	return compiler.EncodingOptions{
		VideoCodec:   "libvpx-vp9",
		AudioCodec:   "libopus",
		PixelFormat:  "yuv420p",
		AudioBitrate: "128k",
		CRF:          crf,
		Preset:       "medium",
		CustomArgs: []string{
			"-b:v", "0",
			"-row-mt", "1",
		},
	}
}

// ProfileWebMAlpha returns a WebM VP9 configuration with alpha transparency channel preservation.
func ProfileWebMAlpha(constantRateFactor int) compiler.EncodingOptions {
	crf := constantRateFactor
	if crf <= 0 {
		crf = 28
	}

	return compiler.EncodingOptions{
		VideoCodec:   "libvpx-vp9",
		AudioCodec:   "libopus",
		PixelFormat:  "yuva420p",
		AudioBitrate: "128k",
		CRF:          crf,
		Preset:       "good",
		CustomArgs: []string{
			"-b:v", "0",
			"-auto-alt-ref", "0",
			"-row-mt", "1",
		},
	}
}
