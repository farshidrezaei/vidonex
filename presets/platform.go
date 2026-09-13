package presets

import (
	"github.com/farshidrezaei/vidonex/compiler"
	"github.com/farshidrezaei/vidonex/timeline"
	"github.com/farshidrezaei/vidonex/types"
)

// PlatformPreset defines turnkey canvas dimensions, frame rate, and encoding configurations for major platforms.
type PlatformPreset struct {
	Name        string
	Description string
	Canvas      types.Size
	FPS         types.Rational
	Encoding    compiler.EncodingOptions
}

// TikTokVertical1080p60 returns optimized 9:16 vertical 60fps preset for TikTok.
func TikTokVertical1080p60() PlatformPreset {
	return PlatformPreset{
		Name:        "TikTok 1080p 60fps (9:16)",
		Description: "Optimized for TikTok vertical video feed and smooth 60fps playback",
		Canvas:      types.ResPortrait1080p,
		FPS:         types.FPS60,
		Encoding: compiler.EncodingOptions{
			VideoCodec:   "libx264",
			AudioCodec:   "aac",
			PixelFormat:  "yuv420p",
			AudioBitrate: "192k",
			CRF:          21,
			Preset:       "fast",
		},
	}
}

// InstagramReels1080p30 returns 9:16 vertical 30fps preset for Instagram Reels and Stories.
func InstagramReels1080p30() PlatformPreset {
	return PlatformPreset{
		Name:        "Instagram Reels 1080p 30fps (9:16)",
		Description: "Standard format for Instagram Reels, Stories, and Facebook Reels",
		Canvas:      types.ResPortrait1080p,
		FPS:         types.FPS30,
		Encoding: compiler.EncodingOptions{
			VideoCodec:   "libx264",
			AudioCodec:   "aac",
			PixelFormat:  "yuv420p",
			AudioBitrate: "192k",
			CRF:          22,
			Preset:       "medium",
		},
	}
}

// InstagramSquare1080p returns 1:1 square preset for Instagram Feed posts.
func InstagramSquare1080p() PlatformPreset {
	return PlatformPreset{
		Name:        "Instagram Square 1080p (1:1)",
		Description: "Square aspect ratio for Instagram timeline posts and carousels",
		Canvas:      types.ResSquare1080,
		FPS:         types.FPS30,
		Encoding: compiler.EncodingOptions{
			VideoCodec:   "libx264",
			AudioCodec:   "aac",
			PixelFormat:  "yuv420p",
			AudioBitrate: "192k",
			CRF:          22,
			Preset:       "medium",
		},
	}
}

// YouTube4K60 returns high-fidelity 4K UHD 60fps preset for YouTube long-form content.
func YouTube4K60() PlatformPreset {
	return PlatformPreset{
		Name:        "YouTube 4K UHD 60fps (16:9)",
		Description: "High-bitrate master preset for 4K cinema and gaming on YouTube",
		Canvas:      types.Res4K,
		FPS:         types.FPS60,
		Encoding: compiler.EncodingOptions{
			VideoCodec:   "libx264",
			AudioCodec:   "aac",
			PixelFormat:  "yuv420p",
			AudioBitrate: "320k",
			CRF:          18,
			Preset:       "slow",
		},
	}
}

// YouTube1080p60 returns standard 1080p 60fps preset for YouTube.
func YouTube1080p60() PlatformPreset {
	return PlatformPreset{
		Name:        "YouTube 1080p 60fps (16:9)",
		Description: "Standard high-definition widescreen video for YouTube",
		Canvas:      types.Res1080p,
		FPS:         types.FPS60,
		Encoding: compiler.EncodingOptions{
			VideoCodec:   "libx264",
			AudioCodec:   "aac",
			PixelFormat:  "yuv420p",
			AudioBitrate: "256k",
			CRF:          20,
			Preset:       "medium",
		},
	}
}

// WebFast720p30 returns a lightweight 720p preset for fast web streaming and low bandwidth.
func WebFast720p30() PlatformPreset {
	return PlatformPreset{
		Name:        "Web Fast 720p 30fps (16:9)",
		Description: "Compact file size optimized for fast loading on websites and mobile browsers",
		Canvas:      types.Res720p,
		FPS:         types.FPS30,
		Encoding: compiler.EncodingOptions{
			VideoCodec:   "libx264",
			AudioCodec:   "aac",
			PixelFormat:  "yuv420p",
			AudioBitrate: "128k",
			CRF:          24,
			Preset:       "veryfast",
		},
	}
}

// ApplyToTimeline sets the canvas resolution and frame rate on the timeline to match the preset.
func (preset PlatformPreset) ApplyToTimeline(compositionTimeline *timeline.Timeline) {
	compositionTimeline.Canvas = preset.Canvas
	compositionTimeline.FPS = preset.FPS
}
