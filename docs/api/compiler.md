# Compiler Package (`compiler`)

The `compiler` package turns a `timeline.Timeline` AST into an intermediate representation Directed Acyclic Graph (`filtergraph.Graph`), applies optimization passes, and emits final FFmpeg arguments.

## Compiler Usage

```go
import "github.com/farshidrezaei/vidonex/compiler"

comp := compiler.New()

// Compile AST into CompilationResult
result, err := comp.Compile(tl, "output.mp4")
if err != nil {
    log.Fatalf("Compilation error: %v", err)
}
```

---

## `CompilationResult`

The result object contains all compiled metadata:

```go
// Generate the full FFmpeg command slice
args := result.BuildFFmpegArgs()
// e.g. ["-y", "-i", "input.mp4", "-filter_complex", "...", "-map", "[outv]", "output.mp4"]

// Generate interactive Mermaid flowchart diagram
mermaidString, err := result.Mermaid()

// Inspect the underlying Graph IR
graph := result.Graph()
```

---

## Hardware Acceleration Flags

The compiler accepts acceleration options to inject NVENC, VideoToolbox, QSV, or VAAPI encoding arguments:

```go
comp := compiler.New(
    compiler.WithHardwareAcceleration(presets.AcceleratorNVENC),
)
```
