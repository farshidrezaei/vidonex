// Package presets provides platform-specific render configurations and GPU hardware acceleration profiles.
package presets

import (
	"context"
	"os/exec"
	"strings"

	"github.com/farshidrezaei/vidonyx/compiler"
)

// HardwareAccelerator enumerates supported GPU and dedicated hardware encoders.
type HardwareAccelerator string

// Hardware accelerator constants.
const (
	// AcceleratorNone uses pure software CPU encoding (libx264, libx265).
	AcceleratorNone HardwareAccelerator = "none"
	// AcceleratorNVENC uses NVIDIA GPU hardware encoding (h264_nvenc, hevc_nvenc).
	AcceleratorNVENC HardwareAccelerator = "nvenc"
	// AcceleratorVideoToolbox uses Apple Silicon / macOS hardware encoding (h264_videotoolbox, hevc_videotoolbox).
	AcceleratorVideoToolbox HardwareAccelerator = "videotoolbox"
	// AcceleratorQSV uses Intel QuickSync hardware encoding (h264_qsv, hevc_qsv).
	AcceleratorQSV HardwareAccelerator = "qsv"
	// AcceleratorVAAPI uses Linux VAAPI hardware encoding (h264_vaapi, hevc_vaapi).
	AcceleratorVAAPI HardwareAccelerator = "vaapi"
)

// DetectHardware inspects the local ffmpeg binary to identify available GPU acceleration encoders.
func DetectHardware(ctx context.Context, ffmpegBinary string) HardwareAccelerator {
	if ffmpegBinary == "" {
		ffmpegBinary = "ffmpeg"
	}

	command := exec.CommandContext(ctx, ffmpegBinary, "-encoders")
	outputBytes, err := command.Output()
	if err != nil {
		return AcceleratorNone
	}

	outputString := string(outputBytes)

	if strings.Contains(outputString, "h264_nvenc") {
		return AcceleratorNVENC
	}
	if strings.Contains(outputString, "h264_videotoolbox") {
		return AcceleratorVideoToolbox
	}
	if strings.Contains(outputString, "h264_qsv") {
		return AcceleratorQSV
	}
	if strings.Contains(outputString, "h264_vaapi") {
		return AcceleratorVAAPI
	}

	return AcceleratorNone
}

// ConfigureHardwareEncoding updates compiler.EncodingOptions to utilize the specified hardware accelerator.
func ConfigureHardwareEncoding(options *compiler.EncodingOptions, accelerator HardwareAccelerator, useHEVC bool) {
	switch accelerator {
	case AcceleratorNVENC:
		if useHEVC {
			options.VideoCodec = "hevc_nvenc"
		} else {
			options.VideoCodec = "h264_nvenc"
		}
		options.Preset = "p4" // Balanced NVENC preset

	case AcceleratorVideoToolbox:
		if useHEVC {
			options.VideoCodec = "hevc_videotoolbox"
		} else {
			options.VideoCodec = "h264_videotoolbox"
		}
		options.Preset = "" // VideoToolbox handles bitrate automatically

	case AcceleratorQSV:
		if useHEVC {
			options.VideoCodec = "hevc_qsv"
		} else {
			options.VideoCodec = "h264_qsv"
		}

	case AcceleratorVAAPI:
		if useHEVC {
			options.VideoCodec = "hevc_vaapi"
		} else {
			options.VideoCodec = "h264_vaapi"
		}

	default:
		if useHEVC {
			options.VideoCodec = "libx265"
		} else {
			options.VideoCodec = "libx264"
		}
	}
}
