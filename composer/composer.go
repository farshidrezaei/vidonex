// Package composer provides the high-level orchestration facade for declarative video composition and rendering.
package composer

import (
	"context"
	"fmt"
	"log/slog"

	"os"

	"github.com/farshidrezaei/vidonex/compiler"
	"github.com/farshidrezaei/vidonex/executor"
	"github.com/farshidrezaei/vidonex/presets"
	"github.com/farshidrezaei/vidonex/probe"
	"github.com/farshidrezaei/vidonex/timeline"
)

// RenderResult contains the metadata and artifacts of a rendered composition.
type RenderResult struct {
	Compilation *compiler.CompilationResult
	OutputPath  string
}

// ProgressHandler is a callback invoked with live rendering telemetry.
type ProgressHandler func(executor.ProgressEvent)

// Option configures the Composer runtime.
type Option func(*Composer)

// WithProber configures a custom MediaProber for stream property detection.
func WithProber(customProber probe.MediaProber) Option {
	return func(composerInstance *Composer) {
		composerInstance.prober = customProber
	}
}

// WithExecutor configures a custom CommandExecutor (e.g. MockExecutor for testing).
func WithExecutor(customExecutor executor.CommandExecutor) Option {
	return func(composerInstance *Composer) {
		composerInstance.executor = customExecutor
	}
}

// WithLogger sets the structured logger.
func WithLogger(logger *slog.Logger) Option {
	return func(composerInstance *Composer) {
		composerInstance.logger = logger
	}
}

// WithEncoding sets custom video/audio encoding parameters.
func WithEncoding(options compiler.EncodingOptions) Option {
	return func(composerInstance *Composer) {
		composerInstance.encoding = options
	}
}

// WithPreset applies platform preset encoding parameters.
func WithPreset(preset presets.PlatformPreset) Option {
	return func(composerInstance *Composer) {
		composerInstance.encoding = preset.Encoding
	}
}

// WithHardwareAcceleration configures GPU encoder acceleration.
func WithHardwareAcceleration(accelerator presets.HardwareAccelerator, useHEVC bool) Option {
	return func(composerInstance *Composer) {
		presets.ConfigureHardwareEncoding(&composerInstance.encoding, accelerator, useHEVC)
	}
}

// WithBinaryPath sets custom binary name or absolute path for ffmpeg.
func WithBinaryPath(binaryPath string) Option {
	return func(composerInstance *Composer) {
		composerInstance.binaryPath = binaryPath
	}
}

// Composer is the top-level facade and orchestrator for declarative video composition and rendering.
type Composer struct {
	executor   executor.CommandExecutor
	logger     *slog.Logger
	encoding   compiler.EncodingOptions
	binaryPath string
	prober     probe.MediaProber
}

// New creates an initialized Composer with standard defaults.
func New(options ...Option) *Composer {
	composerInstance := &Composer{
		executor:   executor.NewOSExecutor(),
		logger:     slog.Default(),
		encoding:   compiler.DefaultEncodingOptions(),
		binaryPath: "ffmpeg",
		prober:     probe.NewCachedProber(probe.NewFFprobeProber("ffprobe")),
	}
	for _, option := range options {
		option(composerInstance)
	}
	return composerInstance
}

// Compile translates a timeline into a compilation result without rendering.
func (composerInstance *Composer) Compile(compositionTimeline *timeline.Timeline, outputPath string) (*compiler.CompilationResult, error) {
	if composerInstance.prober != nil {
		for _, track := range compositionTimeline.Tracks {
			for _, clip := range track.Clips {
				// If source file exists on disk, inspect presence of an audio stream
				if clip.Source != "" {
					if fileInfo, statErr := os.Stat(clip.Source); statErr == nil && !fileInfo.IsDir() {
						if metadata, probeErr := composerInstance.prober.Probe(context.Background(), clip.Source); probeErr == nil && metadata != nil {
							clip.HasAudioStream = (metadata.AudioStream != nil)
						}
					}
				}
			}
		}
	}

	comp := compiler.New(composerInstance.logger).SetEncodingOptions(composerInstance.encoding)
	return comp.Compile(compositionTimeline, outputPath)
}

// Render compiles the timeline into an FFmpeg command and executes the render pipeline.
func (composerInstance *Composer) Render(ctx context.Context, compositionTimeline *timeline.Timeline, outputPath string, onProgress ProgressHandler) (*RenderResult, error) {
	composerInstance.logger.Info("starting composition compilation", slog.String("output", outputPath))

	compilation, err := composerInstance.Compile(compositionTimeline, outputPath)
	if err != nil {
		return nil, fmt.Errorf("composer: compilation failed: %w", err)
	}

	composerInstance.logger.Info("executing render process",
		slog.Int("input_count", len(compilation.Inputs)),
		slog.String("output", outputPath))

	var progressCallback func(executor.ProgressEvent)
	if onProgress != nil {
		progressCallback = func(event executor.ProgressEvent) {
			onProgress(event)
		}
	}

	err = composerInstance.executor.Run(ctx, composerInstance.binaryPath, compilation.Args, compositionTimeline.Duration(), progressCallback)
	if err != nil {
		return nil, fmt.Errorf("composer: render execution failed: %w", err)
	}

	composerInstance.logger.Info("render completed successfully", slog.String("output", outputPath))

	return &RenderResult{
		Compilation: compilation,
		OutputPath:  outputPath,
	}, nil
}
