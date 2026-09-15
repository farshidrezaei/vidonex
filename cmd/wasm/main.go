//go:build js && wasm

// Package main implements the WebAssembly (WASM) entrypoint for the Vidonex video composition engine.
// It allows running the declarative AST parser, timeline validator, filtergraph compiler,
// and Mermaid.js DAG visualizer directly inside web browsers without requiring a backend server.
package main

import (
	"strings"
	"syscall/js"

	"github.com/farshidrezaei/vidonex/compiler"
	"github.com/farshidrezaei/vidonex/spec"
	"github.com/farshidrezaei/vidonex/timeline"
)

const WasmEngineVersion = "1.3.0"

func parseSpecFromString(content string) (*spec.VideoSpec, error) {
	trimmed := strings.TrimSpace(content)
	if strings.HasPrefix(trimmed, "{") {
		return spec.ParseJSON(strings.NewReader(trimmed))
	}
	return spec.ParseYAML(strings.NewReader(trimmed))
}

// validateTimeline validates a YAML or JSON specification string in the browser.
func validateTimeline(this js.Value, args []js.Value) any {
	if len(args) == 0 {
		return js.ValueOf(map[string]any{
			"valid": false,
			"error": "validateTimeline requires a YAML or JSON string argument",
		})
	}

	content := args[0].String()
	videoSpec, err := parseSpecFromString(content)
	if err != nil {
		return js.ValueOf(map[string]any{
			"valid": false,
			"error": "Failed to parse specification: " + err.Error(),
		})
	}

	compositionTimeline, _, _, err := spec.ToTimeline(videoSpec)
	if err != nil {
		return js.ValueOf(map[string]any{
			"valid": false,
			"error": "Failed to convert specification to timeline: " + err.Error(),
		})
	}

	if err := timeline.Validate(compositionTimeline); err != nil {
		return js.ValueOf(map[string]any{
			"valid": false,
			"error": "Timeline validation error: " + err.Error(),
		})
	}

	return js.ValueOf(map[string]any{
		"valid":    true,
		"duration": compositionTimeline.Duration().Seconds(),
		"tracks":   len(compositionTimeline.Tracks),
	})
}

// compileTimeline compiles a specification into a full FFmpeg filtergraph DAG and CLI command.
func compileTimeline(this js.Value, args []js.Value) any {
	if len(args) == 0 {
		return js.ValueOf(map[string]any{
			"success": false,
			"error":   "compileTimeline requires a specification string argument",
		})
	}

	content := args[0].String()
	outputPath := "output.mp4"
	if len(args) > 1 && args[1].String() != "" {
		outputPath = args[1].String()
	}

	videoSpec, err := parseSpecFromString(content)
	if err != nil {
		return js.ValueOf(map[string]any{
			"success": false,
			"error":   "Parser error: " + err.Error(),
		})
	}

	compositionTimeline, encodingOptions, resolvedOutput, err := spec.ToTimeline(videoSpec)
	if err != nil {
		return js.ValueOf(map[string]any{
			"success": false,
			"error":   "Timeline conversion error: " + err.Error(),
		})
	}

	if resolvedOutput != "" && (len(args) <= 1 || args[1].String() == "") {
		outputPath = resolvedOutput
	}

	comp := compiler.New(nil)
	if encodingOptions != nil {
		comp.SetEncodingOptions(*encodingOptions)
	}

	result, err := comp.Compile(compositionTimeline, outputPath)
	if err != nil {
		return js.ValueOf(map[string]any{
			"success": false,
			"error":   "Compilation error: " + err.Error(),
		})
	}

	mermaidCode, _ := result.Mermaid()
	var cliArgs []any
	for _, a := range result.Args {
		cliArgs = append(cliArgs, a)
	}

	var inputFiles []any
	for _, in := range result.Inputs {
		inputFiles = append(inputFiles, in)
	}

	fullCommand := "ffmpeg " + strings.Join(result.Args, " ")

	return js.ValueOf(map[string]any{
		"success":       true,
		"command":       fullCommand,
		"filterComplex": result.FilterComplex,
		"mermaid":       mermaidCode,
		"inputs":        cliArgs,
		"inputFiles":    inputFiles,
		"duration":      compositionTimeline.Duration().Seconds(),
	})
}

// renderMermaid returns just the Mermaid.js flowchart code for the compiled filtergraph.
func renderMermaid(this js.Value, args []js.Value) any {
	if len(args) == 0 {
		return js.ValueOf(map[string]any{
			"success": false,
			"error":   "renderMermaid requires a specification string",
		})
	}

	content := args[0].String()
	videoSpec, err := parseSpecFromString(content)
	if err != nil {
		return js.ValueOf(map[string]any{
			"success": false,
			"error":   err.Error(),
		})
	}

	compositionTimeline, _, _, err := spec.ToTimeline(videoSpec)
	if err != nil {
		return js.ValueOf(map[string]any{
			"success": false,
			"error":   err.Error(),
		})
	}

	comp := compiler.New(nil)
	result, err := comp.Compile(compositionTimeline, "output.mp4")
	if err != nil {
		return js.ValueOf(map[string]any{
			"success": false,
			"error":   err.Error(),
		})
	}

	code, err := result.Mermaid()
	if err != nil {
		return js.ValueOf(map[string]any{
			"success": false,
			"error":   err.Error(),
		})
	}

	return js.ValueOf(map[string]any{
		"success": true,
		"mermaid": code,
	})
}

func getVersion(this js.Value, args []js.Value) any {
	return js.ValueOf(WasmEngineVersion)
}

func main() {
	bridgeObject := js.Global().Get("Object").New()
	bridgeObject.Set("validate", js.FuncOf(validateTimeline))
	bridgeObject.Set("compile", js.FuncOf(compileTimeline))
	bridgeObject.Set("mermaid", js.FuncOf(renderMermaid))
	bridgeObject.Set("version", js.FuncOf(getVersion))

	js.Global().Set("vidonex", bridgeObject)

	// Keep WebAssembly event loop active
	select {}
}
