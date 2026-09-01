package compiler

import (
	"fmt"
)

// EncodingOptions configures video and audio codecs and bitrates.
type EncodingOptions struct {
	VideoCodec   string // e.g. "libx264"
	AudioCodec   string // e.g. "aac"
	PixelFormat  string // e.g. "yuv420p"
	AudioBitrate string // e.g. "192k"
	CRF          int    // Constant Rate Factor (default 23)
	Preset       string // e.g. "medium", "fast"
	CustomArgs   []string
}

// DefaultEncodingOptions returns standard web-compatible H.264/AAC encoding options.
func DefaultEncodingOptions() EncodingOptions {
	return EncodingOptions{
		VideoCodec:   "libx264",
		AudioCodec:   "aac",
		PixelFormat:  "yuv420p",
		AudioBitrate: "192k",
		CRF:          23,
		Preset:       "medium",
	}
}

// BuildFFmpegArgs assembles the complete CLI argument slice for FFmpeg.
func BuildFFmpegArgs(inputs []string, filterComplex string, outVLabel, outALabel, outputPath string, opts EncodingOptions) []string {
	args := make([]string, 0, 32)

	// Overwrite output files without asking
	args = append(args, "-y")

	// Hide FFmpeg banner
	args = append(args, "-hide_banner")

	// Inputs
	for _, inPath := range inputs {
		if isImageSource(inPath) {
			args = append(args, "-loop", "1", "-i", inPath)
		} else {
			args = append(args, "-i", inPath)
		}
	}

	// Filter complex
	if filterComplex != "" {
		args = append(args, "-filter_complex", filterComplex)
	}

	// Stream mapping
	if outVLabel != "" {
		args = append(args, "-map", fmt.Sprintf("[%s]", outVLabel))
	}
	if outALabel != "" {
		args = append(args, "-map", fmt.Sprintf("[%s]", outALabel))
	}

	// Video encoding settings
	if opts.VideoCodec != "" {
		args = append(args, "-c:v", opts.VideoCodec)
	}
	if opts.PixelFormat != "" {
		args = append(args, "-pix_fmt", opts.PixelFormat)
	}
	if opts.CRF > 0 {
		args = append(args, "-crf", fmt.Sprintf("%d", opts.CRF))
	}
	if opts.Preset != "" {
		args = append(args, "-preset", opts.Preset)
	}

	// Audio encoding settings
	if opts.AudioCodec != "" {
		args = append(args, "-c:a", opts.AudioCodec)
	}
	if opts.AudioBitrate != "" {
		args = append(args, "-b:a", opts.AudioBitrate)
	}

	// Custom additional arguments
	if len(opts.CustomArgs) > 0 {
		args = append(args, opts.CustomArgs...)
	}

	// Output file path
	args = append(args, outputPath)

	return args
}
