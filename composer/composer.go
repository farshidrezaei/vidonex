package composer

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/farshidrezaei/vidonyx/compiler"
	"github.com/farshidrezaei/vidonyx/executor"
	"github.com/farshidrezaei/vidonyx/timeline"
)

// RenderResult contains the metadata and artifacts of a rendered composition.
type RenderResult struct {
	Compilation *compiler.CompilationResult
	OutputPath  string
}

// ProgressHandler is a callback invoked with live rendering telemetry.
type ProgressHandler func(executor.ProgressEvent)

// ComposerOption configures the Composer runtime.
type ComposerOption func(*Composer)

// WithExecutor configures a custom CommandExecutor (e.g. MockExecutor for testing).
func WithExecutor(exec executor.CommandExecutor) ComposerOption {
	return func(c *Composer) {
		c.executor = exec
	}
}

// WithLogger sets the structured logger.
func WithLogger(logger *slog.Logger) ComposerOption {
	return func(c *Composer) {
		c.logger = logger
	}
}

// WithEncoding sets custom video/audio encoding parameters.
func WithEncoding(opts compiler.EncodingOptions) ComposerOption {
	return func(c *Composer) {
		c.encoding = opts
	}
}

// WithBinaryPath sets custom binary name or absolute path for ffmpeg.
func WithBinaryPath(binPath string) ComposerOption {
	return func(c *Composer) {
		c.binPath = binPath
	}
}

// Composer is the top-level facade and orchestrator for declarative video composition and rendering.
type Composer struct {
	executor executor.CommandExecutor
	logger   *slog.Logger
	encoding compiler.EncodingOptions
	binPath  string
}

// New creates an initialized Composer with standard defaults.
func New(opts ...ComposerOption) *Composer {
	c := &Composer{
		executor: executor.NewOSExecutor(),
		logger:   slog.Default(),
		encoding: compiler.DefaultEncodingOptions(),
		binPath:  "ffmpeg",
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Inspect / Compile translates a timeline into a compilation result without rendering.
func (c *Composer) Compile(tl *timeline.Timeline, outputPath string) (*compiler.CompilationResult, error) {
	comp := compiler.New(c.logger).SetEncodingOptions(c.encoding)
	return comp.Compile(tl, outputPath)
}

// Render compiles the timeline into an FFmpeg command and executes the render pipeline.
func (c *Composer) Render(ctx context.Context, tl *timeline.Timeline, outputPath string, onProgress ProgressHandler) (*RenderResult, error) {
	c.logger.Info("starting composition compilation", slog.String("output", outputPath))

	compilation, err := c.Compile(tl, outputPath)
	if err != nil {
		return nil, fmt.Errorf("composer: compilation failed: %w", err)
	}

	c.logger.Info("executing render process",
		slog.Int("input_count", len(compilation.Inputs)),
		slog.String("output", outputPath))

	var progressCb func(executor.ProgressEvent)
	if onProgress != nil {
		progressCb = func(ev executor.ProgressEvent) {
			onProgress(ev)
		}
	}

	err = c.executor.Run(ctx, c.binPath, compilation.Args, tl.Duration(), progressCb)
	if err != nil {
		return nil, fmt.Errorf("composer: render execution failed: %w", err)
	}

	c.logger.Info("render completed successfully", slog.String("output", outputPath))

	return &RenderResult{
		Compilation: compilation,
		OutputPath:  outputPath,
	}, nil
}
