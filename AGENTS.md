# AGENTS.md: Developer & AI Agent Context Guide

This guide is optimized for autonomous AI coding agents (and human engineers) working in the **Vidonyx** repository.

---

## 1. Project Mission & Identity
- **Repository**: `github.com/farshidrezaei/vidonyx`
- **Purpose**: Declarative Video Composition & FFmpeg Filtergraph Compiler Engine in Go (conceptually similar to Remotion / Editframe / CapCut backend, but written in native Go).
- **Core Strategy**: Abstract video composition into a high-level AST (`timeline/`), compile to an Intermediate Representation DAG (`filtergraph/`), run graph optimization passes (split-injection, dead-code elimination, normalization), and emit valid FFmpeg CLI arguments (`-filter_complex`).

---

## 2. Invariants & Rules for AI Agents

1. **No String Concatenation for Filtergraphs**:
   - NEVER construct raw FFmpeg filter strings manually with `fmt.Sprintf` or string concatenation in domain code.
   - Always build a `filtergraph.Graph`, add `filtergraph.Node`s with `filtergraph.Pad`s, and connect them with `g.Connect(srcOut, dstIn)`.
2. **Never Use Floating-Point `float64` for Core Time Calculations**:
   - Always use `types.Rational` or `time.Duration` to prevent floating-point precision drift across video frames.
3. **Stream Type Safety**:
   - Video pads (`filtergraph.StreamTypeVideo`) can ONLY connect to video pads. Audio pads (`filtergraph.StreamTypeAudio`) can ONLY connect to audio pads.
4. **Standard Library First**:
   - The core packages (`types`, `timeline`, `filtergraph`, `compiler`) must remain 100% standard library only.
5. **Deterministic Testing with Mocks**:
   - Unit tests MUST NOT depend on a locally installed `ffmpeg` binary. Use `executor.NewMockExecutor()` and graph structure assertions.
6. **Even Dimension Constraint**:
   - Video resolutions must always have even width and height for H.264 / `yuv420p` compatibility (`types.makeEven`).

---

## 3. Package Layout Quick Reference

| Package | Purpose | Key Types / Entrypoints |
| :--- | :--- | :--- |
| `types` | Exact Math & Geometry | `Rational`, `Size`, `Point`, `Rect`, `Color`, `Alignment` |
| `timeline` | Declarative Timeline AST | `Timeline`, `Track`, `Clip`, `Transition`, `Effect`, `Validate()` |
| `filtergraph` | Graph Intermediate Representation | `Graph`, `Node`, `Pad`, `AutoSplitPass`, `DeadCodeEliminationPass` |
| `compiler` | AST $\rightarrow$ DAG $\rightarrow$ CLI Compiler | `Compiler`, `CompilationResult`, `BuildFFmpegArgs()` |
| `effects` | Pre-built Type-Safe Filters | `ScaleFilter`, `DrawTextFilter`, `ChromaKeyFilter`, `VolumeFilter`, `XFadeFilter` |
| `executor` | Subprocess & Telemetry | `CommandExecutor`, `OSExecutor`, `MockExecutor`, `ParseProgressStream()` |
| `visualizer` | Graph Visualization | `ToMermaid()`, `ToDOT()` |
| `composer` | High-Level Facade | `Composer`, `New()`, `Render()`, `Compile()` |
| `examples` | Concrete Recipes | `01_simple_cut`, `02_picture_in_picture` |

---

## 4. Common Developer Workflows

### Running All Tests with Race Detector:
```bash
go test -race -v ./...
```

### Running Table-Driven Tests for Specific Packages:
```bash
go test -race -v ./timeline/...
go test -race -v ./filtergraph/...
go test -race -v ./compiler/...
```

### Running Example Recipes:
```bash
go run ./examples/01_simple_cut/main.go
go run ./examples/02_picture_in_picture/main.go
```
