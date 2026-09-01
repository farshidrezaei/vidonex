// Package main implements the Vidonyx standalone command-line video composition interface.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/farshidrezaei/vidonyx/composer"
	"github.com/farshidrezaei/vidonyx/executor"
	"github.com/farshidrezaei/vidonyx/internal/cli"
	"github.com/farshidrezaei/vidonyx/presets"
	"github.com/farshidrezaei/vidonyx/probe"
	"github.com/farshidrezaei/vidonyx/server"
	"github.com/farshidrezaei/vidonyx/spec"
)

const (
	EngineVersion = "1.0.0"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := strings.ToLower(os.Args[1])

	switch command {
	case "serve", "studio", "ui":
		executeServeCommand(os.Args[2:])
	case "render":
		executeRenderCommand(os.Args[2:])
	case "validate":
		executeValidateCommand(os.Args[2:])
	case "graph":
		executeGraphCommand(os.Args[2:])
	case "probe":
		executeProbeCommand(os.Args[2:])
	case "version", "--version", "-v":
		fmt.Printf("vidonyx engine version %s\n", EngineVersion)
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %q\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func executeRenderCommand(arguments []string) {
	fs := flag.NewFlagSet("render", flag.ExitOnError)
	outputFlag := fs.String("o", "", "Output file path (overrides spec)")
	outputLongFlag := fs.String("output", "", "Output file path (overrides spec)")
	gpuFlag := fs.String("gpu", "", "Hardware accelerator (nvenc, videotoolbox, qsv, vaapi)")
	logLevelFlag := fs.String("log-level", "info", "Log level (debug, info, warn, error)")
	logFormatFlag := fs.String("log-format", "text", "Log format (text, json)")
	dryRunFlag := fs.Bool("dry-run", false, "Compile without invoking FFmpeg render")

	if err := fs.Parse(reorderFlags(arguments)); err != nil {
		fmt.Fprintf(os.Stderr, "Failed parsing flags: %v\n", err)
		os.Exit(1)
	}

	positionalArgs := fs.Args()
	if len(positionalArgs) < 1 {
		fmt.Fprintln(os.Stderr, "Error: missing specification file path. Usage: vidonyx render <project.yaml> [flags]")
		os.Exit(1)
	}

	specPath := positionalArgs[0]
	outputPathOverride := *outputFlag
	if outputPathOverride == "" {
		outputPathOverride = *outputLongFlag
	}

	ui := cli.NewUI(os.Stdout)
	ui.PrintBanner()

	logger := cli.SetupLogger(*logLevelFlag, *logFormatFlag, os.Stderr)
	logger.Debug("starting render command", "spec_file", specPath)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Step 1: Parse and Resolve
	ui.PrintStep(1, 4, "🔍", fmt.Sprintf("Parsing Specification: %s", filepath.Base(specPath)))
	specification, err := spec.ParseFile(specPath)
	if err != nil {
		ui.PrintError(err)
		os.Exit(1)
	}

	// Resolve relative media paths against directory of spec file
	specDirectory := filepath.Dir(specPath)
	spec.ResolveAssetPaths(specification, specDirectory)

	// Apply CLI overrides
	if outputPathOverride != "" {
		specification.Output = outputPathOverride
	}
	if *gpuFlag != "" {
		specification.HardwareAcceleration = *gpuFlag
	}

	// Step 2: Convert to Timeline AST
	ui.PrintStep(2, 4, "📐", "Constructing Timeline AST & Semantic Constraints")
	compositionTimeline, encodingOptions, finalOutputPath, err := spec.ToTimeline(specification)
	if err != nil {
		ui.PrintError(err)
		os.Exit(1)
	}

	ui.PrintSuccess(fmt.Sprintf("Canvas: %dx%d @ %s fps | Duration: %s | Tracks: %d",
		compositionTimeline.Canvas.Width,
		compositionTimeline.Canvas.Height,
		compositionTimeline.FPS.String(),
		compositionTimeline.Duration().Round(10*time.Millisecond),
		len(compositionTimeline.Tracks),
	))

	// Step 3: Compile Filtergraph DAG
	ui.PrintStep(3, 4, "⚙️ ", "Compiling Filtergraph DAG (Optimization Passes: AutoSplit, DCE)")

	composerOptions := []composer.Option{
		composer.WithLogger(logger),
	}
	if encodingOptions != nil {
		composerOptions = append(composerOptions, composer.WithEncoding(*encodingOptions))
	}
	if specification.Preset != "" {
		resolvedPreset, presetErr := resolvePlatformPreset(specification.Preset)
		if presetErr == nil {
			composerOptions = append(composerOptions, composer.WithPreset(resolvedPreset))
		}
	}
	if specification.HardwareAcceleration != "" && strings.ToLower(specification.HardwareAcceleration) != "none" {
		accelerator := presets.HardwareAccelerator(strings.ToLower(specification.HardwareAcceleration))
		composerOptions = append(composerOptions, composer.WithHardwareAcceleration(accelerator, false))
	}

	composerInstance := composer.New(composerOptions...)

	compilationResult, err := composerInstance.Compile(compositionTimeline, finalOutputPath)
	if err != nil {
		ui.PrintError(err)
		os.Exit(1)
	}

	nodeCount := 0
	for range compilationResult.Graph.Nodes() {
		nodeCount++
	}

	ui.PrintSuccess(fmt.Sprintf("Filtergraph compiled with %d DAG nodes", nodeCount))

	if *dryRunFlag {
		ui.PrintSuccess("Dry-run complete. Skipping FFmpeg render process.")
		fmt.Printf("\nGenerated FFmpeg Arguments:\nffmpeg %s\n", strings.Join(compilationResult.Args, " "))
		return
	}

	// Step 4: Render
	ui.PrintStep(4, 4, "🎬", fmt.Sprintf("Rendering Composition -> %s", finalOutputPath))

	progressBar := cli.NewProgressBar(os.Stdout, compositionTimeline.Duration())
	renderStartTime := time.Now()

	renderResult, err := composerInstance.Render(ctx, compositionTimeline, finalOutputPath, func(event executor.ProgressEvent) {
		progressBar.Update(event)
	})
	progressBar.Finish()

	if err != nil {
		ui.PrintError(err)
		os.Exit(1)
	}

	renderDuration := time.Since(renderStartTime)

	fileInfo, statErr := os.Stat(renderResult.OutputPath)
	var fileSizeBytes int64
	if statErr == nil {
		fileSizeBytes = fileInfo.Size()
	}

	ui.PrintSummaryCard(renderResult.OutputPath, fileSizeBytes, renderDuration, compositionTimeline.Duration())
}

func executeValidateCommand(arguments []string) {
	fs := flag.NewFlagSet("validate", flag.ExitOnError)
	logLevelFlag := fs.String("log-level", "info", "Log level")

	if err := fs.Parse(reorderFlags(arguments)); err != nil {
		os.Exit(1)
	}

	if len(fs.Args()) < 1 {
		fmt.Fprintln(os.Stderr, "Error: missing spec file. Usage: vidonyx validate <project.yaml>")
		os.Exit(1)
	}

	specPath := fs.Args()[0]
	ui := cli.NewUI(os.Stdout)
	cli.SetupLogger(*logLevelFlag, "text", os.Stderr)

	specification, err := spec.ParseFile(specPath)
	if err != nil {
		ui.PrintError(fmt.Errorf("syntax error: %w", err))
		os.Exit(1)
	}

	spec.ResolveAssetPaths(specification, filepath.Dir(specPath))

	compositionTimeline, _, _, err := spec.ToTimeline(specification)
	if err != nil {
		ui.PrintError(fmt.Errorf("validation error: %w", err))
		os.Exit(1)
	}

	ui.PrintSuccess(fmt.Sprintf("Specification %q is valid! (%dx%d @ %s fps, %d tracks)",
		specPath,
		compositionTimeline.Canvas.Width,
		compositionTimeline.Canvas.Height,
		compositionTimeline.FPS.String(),
		len(compositionTimeline.Tracks),
	))
}

func executeGraphCommand(arguments []string) {
	fs := flag.NewFlagSet("graph", flag.ExitOnError)
	formatFlag := fs.String("format", "mermaid", "Graph format (mermaid, dot)")
	outputFlag := fs.String("o", "", "Write graph to output file")

	if err := fs.Parse(reorderFlags(arguments)); err != nil {
		os.Exit(1)
	}

	if len(fs.Args()) < 1 {
		fmt.Fprintln(os.Stderr, "Error: missing spec file. Usage: vidonyx graph <project.yaml> [--format mermaid|dot]")
		os.Exit(1)
	}

	specPath := fs.Args()[0]
	specification, err := spec.ParseFile(specPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	spec.ResolveAssetPaths(specification, filepath.Dir(specPath))

	compositionTimeline, _, _, err := spec.ToTimeline(specification)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	composerInstance := composer.New()
	compilation, err := composerInstance.Compile(compositionTimeline, "temp.mp4")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Compilation error: %v\n", err)
		os.Exit(1)
	}

	var graphContent string
	if strings.ToLower(*formatFlag) == "dot" {
		graphContent, err = compilation.DOT()
	} else {
		graphContent, err = compilation.Mermaid()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Graph generation error: %v\n", err)
		os.Exit(1)
	}

	if *outputFlag != "" {
		if writeErr := os.WriteFile(*outputFlag, []byte(graphContent), 0644); writeErr != nil {
			fmt.Fprintf(os.Stderr, "Failed writing graph file: %v\n", writeErr)
			os.Exit(1)
		}
		fmt.Printf("Graph saved to %s\n", *outputFlag)
	} else {
		fmt.Println(graphContent)
	}
}

func executeProbeCommand(arguments []string) {
	fs := flag.NewFlagSet("probe", flag.ExitOnError)
	jsonFlag := fs.Bool("json", false, "Output in JSON format")

	if err := fs.Parse(reorderFlags(arguments)); err != nil {
		os.Exit(1)
	}

	if len(fs.Args()) < 1 {
		fmt.Fprintln(os.Stderr, "Error: missing media file. Usage: vidonyx probe <media.mp4> [--json]")
		os.Exit(1)
	}

	mediaPath := fs.Args()[0]
	mediaProber := probe.NewFFprobeProber("ffprobe")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	metadata, err := mediaProber.Probe(ctx, mediaPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Probe error: %v\n", err)
		os.Exit(1)
	}

	if *jsonFlag {
		jsonBytes, jsonErr := json.MarshalIndent(metadata, "", "  ")
		if jsonErr != nil {
			fmt.Fprintf(os.Stderr, "JSON error: %v\n", jsonErr)
			os.Exit(1)
		}
		fmt.Println(string(jsonBytes))
	} else {
		fmt.Printf("Media Probe Summary: %s\n", mediaPath)
		fmt.Printf("  Container : %s (Duration: %s, Bitrate: %d bps)\n", metadata.FormatName, metadata.Duration, metadata.BitRate)
		if metadata.VideoStream != nil {
			fmt.Printf("  Video     : %s (%dx%d @ %s fps, PixelFormat: %s)\n",
				metadata.VideoStream.Codec,
				metadata.VideoStream.Width,
				metadata.VideoStream.Height,
				metadata.VideoStream.FPS.String(),
				metadata.VideoStream.PixelFormat,
			)
		}
		if metadata.AudioStream != nil {
			fmt.Printf("  Audio     : %s (%d Hz, %d Channels)\n",
				metadata.AudioStream.Codec,
				metadata.AudioStream.SampleRate,
				metadata.AudioStream.Channels,
			)
		}
	}
}

func resolvePlatformPreset(presetName string) (presets.PlatformPreset, error) {
	switch strings.ToLower(strings.TrimSpace(presetName)) {
	case "tiktok", "tiktok_1080p60", "tiktok_vertical":
		return presets.TikTokVertical1080p60(), nil
	case "instagram_reels", "reels":
		return presets.InstagramReels1080p30(), nil
	case "instagram_square", "square":
		return presets.InstagramSquare1080p(), nil
	case "youtube_4k", "youtube_4k60":
		return presets.YouTube4K60(), nil
	case "youtube_1080", "youtube_1080p60":
		return presets.YouTube1080p60(), nil
	case "web_fast", "web_720p30":
		return presets.WebFast720p30(), nil
	default:
		return presets.PlatformPreset{}, fmt.Errorf("unknown platform preset %q", presetName)
	}
}

func reorderFlags(arguments []string) []string {
	flags := make([]string, 0, len(arguments))
	positionals := make([]string, 0, len(arguments))

	skipNext := false
	for i, arg := range arguments {
		if skipNext {
			skipNext = false
			continue
		}
		if strings.HasPrefix(arg, "-") {
			flags = append(flags, arg)
			if !strings.Contains(arg, "=") && i+1 < len(arguments) && !strings.HasPrefix(arguments[i+1], "-") {
				if arg != "--json" && arg != "-json" && arg != "--dry-run" && arg != "-dry-run" {
					flags = append(flags, arguments[i+1])
					skipNext = true
				}
			}
		} else {
			positionals = append(positionals, arg)
		}
	}
	return append(flags, positionals...)
}

func executeServeCommand(arguments []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	portFlag := fs.Int("port", 8080, "HTTP server listening port")
	hostFlag := fs.String("host", "0.0.0.0", "HTTP server listening host")
	dataDirFlag := fs.String("data-dir", "", "Data directory for database, media, and exports")
	staticDirFlag := fs.String("static-dir", "ui/dist", "Directory containing built Web Studio UI static files")
	logLevelFlag := fs.String("log-level", "info", "Log verbosity level (debug, info, warn, error)")
	logFormatFlag := fs.String("log-format", "text", "Log format (text, json)")

	if err := fs.Parse(reorderFlags(arguments)); err != nil {
		fmt.Fprintf(os.Stderr, "Failed parsing flags: %v\n", err)
		os.Exit(1)
	}

	ui := cli.NewUI(os.Stdout)
	ui.PrintBanner()

	logger := cli.SetupLogger(*logLevelFlag, *logFormatFlag, os.Stderr)

	staticDirectory := *staticDirFlag
	if staticDirectory == "ui/dist" {
		if _, err := os.Stat("ui/.output/public"); err == nil {
			staticDirectory = "ui/.output/public"
		}
	}

	serverInstance, err := server.New(server.Config{
		Port:            *portFlag,
		Host:            *hostFlag,
		DataDirectory:   *dataDirFlag,
		StaticDirectory: staticDirectory,
		Logger:          logger,
	})
	if err != nil {
		ui.PrintError(fmt.Errorf("failed creating server: %w", err))
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func() {
		<-ctx.Done()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		_ = serverInstance.Stop(shutdownCtx)
	}()

	fmt.Printf("\n✨ Vidonyx Web Studio running on http://localhost:%d\n\n", *portFlag)
	if err := serverInstance.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		ui.PrintError(fmt.Errorf("server error: %w", err))
		os.Exit(1)
	}
}

func printUsage() {
	usage := `
Vidonyx - Declarative Video Composition & FFmpeg Filtergraph Engine

USAGE:
  vidonyx <command> [arguments] [flags]

COMMANDS:
  serve    [flags]              Launch the Vidonyx Web Studio GUI server
  render   <project.yaml|json>  Compile and render video composition
  validate <project.yaml|json>  Validate specification syntax and constraints
  graph    <project.yaml|json>  Export Mermaid.js or Graphviz DOT filtergraph diagram
  probe    <media.mp4>          Inspect media container and stream properties
  version                       Show Vidonyx engine version

SERVE FLAGS:
  --port <number>               HTTP port (default: 8080)
  --host <string>               Listening host (default: 0.0.0.0)
  --data-dir <path>             Custom data directory (default: ~/.vidonyx)
  --static-dir <path>           Static assets directory (default: ui/dist)
  --log-level <level>           Log verbosity (debug, info, warn, error)

RENDER FLAGS:
  -o, --output <path>           Output file path (overrides project spec)
  --gpu <accelerator>           Hardware acceleration (nvenc, videotoolbox, qsv, vaapi)
  --log-level <level>           Log verbosity (debug, info, warn, error)
  --log-format <format>         Log format (text, json)
  --dry-run                     Compile filtergraph without executing FFmpeg

EXAMPLES:
  vidonyx serve --port 8080
  vidonyx render project.yaml -o final.mp4
  vidonyx render project.yaml --gpu nvenc --log-level debug
  vidonyx validate project.yaml
  vidonyx graph project.yaml --format mermaid
  vidonyx probe gameplay.mp4 --json
`
	fmt.Println(strings.TrimSpace(usage))
}
