# AGENTS.md: Developer & AI Agent Context Guide

This guide is optimized for autonomous AI coding agents (and human engineers) working in the **Vidonyx** repository.

---

## 1. Project Mission & Identity
- **Repository**: `github.com/farshidrezaei/vidonyx`
- **Purpose**: Declarative Video Composition & FFmpeg Filtergraph Compiler Engine in Go (conceptually similar to Remotion / Editframe / CapCut backend, but written in native Go).
- **Core Strategy**: Abstract video composition into a high-level AST (`timeline/`), compile to an Intermediate Representation DAG (`filtergraph/`), run graph optimization passes (split-injection, dead-code elimination, normalization), and emit valid FFmpeg CLI arguments (`-filter_complex`).

---

## 2. Core Mandatory Project Contracts & Invariants

1. **Descriptive Naming with No Arbitrary Abbreviations**:
   - Function and method names MUST clearly express their exact operation (e.g. `InjectVideoNormalizer`, `BuildVideoCompositor`, `CalculateOffset`).
   - Variable and parameter names MUST be descriptive and meaningful (e.g. `compositionTimeline`, `durationSeconds`, `sourceOutputPad`, `destinationInputPad`, `consumerCount`).
   - Never use cryptic abbreviations (`tl`, `tr`, `dur`, `res`, `pct`, `sz`, `sc`, `pts`, `inPad`, `outPad`).
2. **No String Concatenation for Filtergraphs**:
   - NEVER construct raw FFmpeg filter strings manually with `fmt.Sprintf` or string concatenation in domain code.
   - Always build a `filtergraph.Graph`, add `filtergraph.Node`s with `filtergraph.Pad`s, and connect them with `graph.Connect(sourceOutputPad, destinationInputPad)`.
3. **Never Use Floating-Point `float64` for Core Time Calculations**:
   - Always use `types.Rational` or `time.Duration` to prevent floating-point precision drift across video frames.
4. **Stream Type Safety**:
   - Video pads (`filtergraph.StreamTypeVideo`) can ONLY connect to video pads. Audio pads (`filtergraph.StreamTypeAudio`) can ONLY connect to audio pads.
5. **Table-Driven Tests Covering All Scenarios**:
   - All unit and feature tests MUST be Table-Driven (`Test...Table`), covering happy paths, edge cases, zero-values, and error conditions.
6. **Mandatory Quality Checks (`golangci-lint` & Race Detector)**:
   - Always run and ensure 100% clean passes for:
     ```bash
     golangci-lint run ./...
     go test -race -v ./...
     ```
7. **Standard Library First & Modern Go Standards**:
   - Core packages (`types`, `timeline`, `filtergraph`, `compiler`) must remain 100% standard library only.
   - Use modern Go idioms (Go iterators `iter.Seq`, `log/slog`, `errors.Join`, generics).

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

### Running Linter:
```bash
golangci-lint run ./...
```

### Running All Tests with Race Detector:
```bash
go test -race -v ./...
```

### Running Example Recipes:
```bash
go run ./examples/01_simple_cut/main.go
go run ./examples/02_picture_in_picture/main.go
```
