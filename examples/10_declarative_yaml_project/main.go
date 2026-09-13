// Recipe 10: Declarative YAML/JSON Video Composition Project.
//
// This example demonstrates how to load, resolve relative media paths, compile,
// and render a complete multi-track video composition from a declarative vidonex.yaml file.
package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/farshidrezaei/vidonex/composer"
	"github.com/farshidrezaei/vidonex/executor"
	"github.com/farshidrezaei/vidonex/spec"
)

func main() {
	specPath := "examples/10_declarative_yaml_project/vidonex.yaml"
	fmt.Printf("=== Recipe 10: Declarative YAML Video Composition ===\n")
	fmt.Printf("Loading project specification: %s\n\n", specPath)

	// 1. Parse YAML Specification
	specification, err := spec.ParseFile(specPath)
	if err != nil {
		log.Fatalf("Failed parsing specification: %v", err)
	}

	// 2. Resolve relative media paths against directory of spec file
	specDirectory := filepath.Dir(specPath)
	spec.ResolveAssetPaths(specification, specDirectory)

	// 3. Convert to Timeline AST
	compositionTimeline, encodingOptions, outputPath, err := spec.ToTimeline(specification)
	if err != nil {
		log.Fatalf("Failed converting spec to timeline: %v", err)
	}

	fmt.Printf("Composition Summary:\n")
	fmt.Printf("  • Canvas   : %dx%d (%s fps)\n", compositionTimeline.Canvas.Width, compositionTimeline.Canvas.Height, compositionTimeline.FPS.String())
	fmt.Printf("  • Duration : %s\n", compositionTimeline.Duration().Round(10*time.Millisecond))
	fmt.Printf("  • Tracks   : %d layers\n", len(compositionTimeline.Tracks))
	fmt.Printf("  • Output   : %s\n\n", outputPath)

	// 4. Compile Filtergraph DAG
	composerOptions := []composer.Option{
		composer.WithExecutor(executor.NewMockExecutor()),
	}
	if encodingOptions != nil {
		composerOptions = append(composerOptions, composer.WithEncoding(*encodingOptions))
	}

	composerInstance := composer.New(composerOptions...)

	compilationResult, err := composerInstance.Compile(compositionTimeline, outputPath)
	if err != nil {
		log.Fatalf("Compilation failed: %v", err)
	}

	fmt.Println("=== Generated FFmpeg Filtergraph Command ===")
	fmt.Printf("ffmpeg %s\n\n", compilationResult.FilterComplex)

	// 5. Render with Mock Executor (or use OSExecutor for real renders)
	ctx := context.Background()
	_, err = composerInstance.Render(ctx, compositionTimeline, outputPath, func(event executor.ProgressEvent) {
		fmt.Printf("Rendering: %.1f%% (Speed: %.2fx, FPS: %.1f)\n", event.Percentage, event.Speed, event.FPS)
	})
	if err != nil {
		log.Fatalf("Render failed: %v", err)
	}

	fmt.Println("Render completed successfully!")
}
