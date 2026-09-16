package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/farshidrezaei/vidonex/internal/cli"
	"github.com/farshidrezaei/vidonex/presets"
	"github.com/farshidrezaei/vidonex/server/api"
)

func executeAboutCommand(_ []string) {
	ui := cli.NewUI(os.Stdout)
	ui.PrintBanner()

	fmt.Printf("\n%s\n", ui.Colorize(cli.ColorWhiteBold, "Vidonex Video Composition & Filtergraph Engine"))
	fmt.Printf("Version:         %s\n", ui.Colorize(cli.ColorCyan+cli.ColorBold, "v"+EngineVersion))
	fmt.Printf("Architecture:    %s / %s (%d CPU cores)\n", runtime.GOOS, runtime.GOARCH, runtime.NumCPU())
	fmt.Printf("Go Runtime:      %s\n", runtime.Version())

	// FFmpeg & Media Tools
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		ffmpegPath = "not found in PATH"
	}
	ffmpegVersion := getToolVersion("ffmpeg")
	fmt.Printf("FFmpeg Binary:   %s\n", ffmpegPath)
	if ffmpegVersion != "" {
		fmt.Printf("FFmpeg Version:  %s\n", ui.Colorize(cli.ColorDim, ffmpegVersion))
	}

	ffprobePath, err := exec.LookPath("ffprobe")
	if err != nil {
		ffprobePath = "not found in PATH"
	}
	ffprobeVersion := getToolVersion("ffprobe")
	fmt.Printf("FFprobe Binary:  %s\n", ffprobePath)
	if ffprobeVersion != "" {
		fmt.Printf("FFprobe Version: %s\n", ui.Colorize(cli.ColorDim, ffprobeVersion))
	}

	// GPU Hardware Acceleration
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	detectedAccelerator := presets.DetectHardware(ctx, "ffmpeg")
	cancel()

	acceleratorName := "None (CPU Software Encoding)"
	if detectedAccelerator != presets.AcceleratorNone {
		acceleratorName = fmt.Sprintf("%s (Hardware Accelerated)", strings.ToUpper(string(detectedAccelerator)))
	}
	fmt.Printf("GPU Accelerator: %s\n", ui.Colorize(cli.ColorGreen, acceleratorName))

	// Repository & Resources
	fmt.Printf("\n%s\n", ui.Colorize(cli.ColorWhiteBold, "Documentation & Community"))
	fmt.Println("GitHub:          https://github.com/farshidrezaei/vidonex")
	fmt.Println("Docs:            https://farshidrezaei.github.io/vidonex/")
	fmt.Println("API Sandbox:     http://localhost:8080/docs (when running `vidonex serve`)")
	fmt.Println("License:         MIT License")

	// Quick Update Check (non-blocking, timeout 2s)
	fmt.Printf("\nChecking for updates... ")
	checkCtx, checkCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer checkCancel()

	release, err := api.FetchLatestRelease(checkCtx, false)
	if err != nil {
		fmt.Printf("%s\n", ui.Colorize(cli.ColorDim, "(Could not connect to GitHub to verify latest version)"))
	} else {
		cleanTag := strings.TrimPrefix(release.TagName, "v")
		if api.CompareSemanticVersions(cleanTag, EngineVersion) > 0 {
			fmt.Printf("%s\n", ui.Colorize(cli.ColorYellow+cli.ColorBold, fmt.Sprintf("A new release v%s is available! (Run `vidonex upgrade` to update)", cleanTag)))
		} else {
			fmt.Printf("%s\n", ui.Colorize(cli.ColorGreen, fmt.Sprintf("You are on the latest version (v%s).", EngineVersion)))
		}
	}
	fmt.Println()
}

func getToolVersion(tool string) string {
	cmd := exec.Command(tool, "-version")
	outputBytes, err := cmd.Output()
	if err != nil {
		return ""
	}
	lines := strings.Split(string(outputBytes), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	}
	return ""
}
