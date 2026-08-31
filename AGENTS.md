# AGENTS.md: Developer & AI Agent Context Guide

This guide is optimized for autonomous AI coding agents (and human engineers) working in the **Vidonyx** repository.

---

## 1. Project Mission & Identity
- **Repository**: `github.com/farshidrezaei/vidonyx`
- **Purpose**: Declarative Video Composition & FFmpeg Filtergraph Compiler Engine in Go (conceptually similar to Remotion / Editframe / CapCut backend, but written natively in Go).
- **Core Strategy**: Abstract video composition into a high-level AST (`timeline/`), compile to an Intermediate Representation DAG (`filtergraph/`), run graph optimization passes (split-injection, dead-code elimination, normalization), and emit valid FFmpeg CLI arguments (`-filter_complex`).

---

## 2. Core Mandatory Project Contracts & Invariants

1. **Descriptive Naming with No Arbitrary Abbreviations**:
   - Function and method names MUST clearly express their exact operation (e.g., `InjectVideoNormalizer`, `BuildVideoCompositor`, `CalculateOffset`, `ApplySidechainDucking`, `ApplyWaveformVisualizer`, `ApplyChromaKey`).
   - Variable and parameter names MUST be descriptive and meaningful (e.g., `compositionTimeline`, `durationSeconds`, `sourceOutputPad`, `destinationInputPad`, `consumerCount`, `frameRate`, `waveformOptions`).
   - Never use cryptic abbreviations (`tl`, `tr`, `dur`, `res`, `pct`, `sz`, `sc`, `pts`, `inPad`, `outPad`).
2. **No String Concatenation for Filtergraphs**:
   - NEVER construct raw FFmpeg filter strings manually with `fmt.Sprintf` or string concatenation in domain code.
   - Always build a `filtergraph.Graph`, add `filtergraph.Node`s with `filtergraph.Pad`s, and connect them with `graph.Connect(sourceOutputPad, destinationInputPad)`.
3. **Never Use Floating-Point `float64` for Core Time Calculations**:
   - Always use `types.Rational` or `time.Duration` to prevent floating-point precision drift across video frames.
4. **Stream Type Safety**:
   - Video pads (`filtergraph.StreamTypeVideo`) can ONLY connect to video pads. Audio pads (`filtergraph.StreamTypeAudio`) can ONLY connect to audio pads.
5. **Table-Driven Tests Covering All Scenarios & Real E2E Tests**:
   - All unit and feature tests MUST be Table-Driven (`Test...Table`), covering happy paths, edge cases, zero-values, and error conditions.
   - The end-to-end suite (`tests/e2e/`) generates synthetic media and tests all 9 composition scenarios against real FFmpeg and FFprobe.
6. **Mandatory Quality Checks (`golangci-lint` & Race Detector)**:
   - Always run and ensure 100% clean passes for:
     ```bash
     golangci-lint run ./...
     go test -race -v ./...
     ```
7. **Standard Library First & Modern Go Standards**:
   - Core packages (`types`, `timeline`, `filtergraph`, `compiler`) must remain 100% standard library only.
   - Use modern Go idioms (Go iterators `iter.Seq`, `log/slog`, `errors.Join`, generics).
8. **Mandatory Continuous Documentation Synchronization**:
   - With every feature addition, refactor, or code modification, ALL relevant documentation (`.md` files, diagrams, contracts, and example recipes) MUST be updated immediately to keep code and documentation 100% in sync.

---

## 3. Package Layout Quick Reference

| Package | Purpose | Key Types / Entrypoints |
| :--- | :--- | :--- |
| `types` | Exact Math & Geometry | `Rational`, `Size`, `Point`, `Rect`, `Color`, `Alignment` |
| `timeline` | Declarative Timeline AST | `Timeline`, `Track`, `Clip`, `Transition`, `Effect`, `Validate()` |
| `filtergraph` | Graph Intermediate Representation | `Graph`, `Node`, `Pad`, `AutoSplitPass`, `DeadCodeEliminationPass` |
| `compiler` | AST $\rightarrow$ DAG $\rightarrow$ CLI Compiler | `Compiler`, `CompilationResult`, `BuildFFmpegArgs()` |
| `animation` | Keyframing & Math Easing | `PositionTrack`, `FloatKeyframeTrack`, `KenBurnsAnimation`, `Easing` |
| `subtitles` | SRT/VTT Parsing & Caption Burn-in | `ParseSRT()`, `ParseVTT()`, `SubtitleTrack`, `AttachSubtitles()` |
| `ducking` | Sidechain Audio Ducking | `ApplySidechainDucking()`, `Options`, `DefaultOptions()` |
| `waveform` | Animated Audio Waveforms | `ApplyWaveformVisualizer()`, `ModePeakToPeak`, `Options` |
| `chromakey` | Green/Blue Screen & Despill | `ApplyChromaKey()`, `Options`, `StudioGreenScreen` |
| `presets` | Social Presets & GPU Acceleration | `TikTokVertical1080p60()`, `YouTube4K60()`, `AcceleratorNVENC` |
| `effects` | Pre-built Type-Safe Filters | `ScaleFilter`, `DrawTextFilter`, `XFadeFilter`, `VolumeFilter` |
| `probe` | Media Probing & Stream Caching | `MediaProber`, `FFprobeProber`, `CachedProber`, `MockProber` |
| `executor` | Subprocess & Live Telemetry | `CommandExecutor`, `OSExecutor`, `MockExecutor`, `ParseProgressStream()` |
| `visualizer` | Graph Visualization | `ToMermaid()`, `ToDOT()` |
| `spec` | Declarative YAML/JSON Parser & Path Resolver | `ParseFile()`, `ParseYAML()`, `ParseJSON()`, `ToTimeline()` |
| `cmd/vidonyx` | Standalone CLI Tool | `render`, `validate`, `graph`, `probe`, `version` |
| `composer` | High-Level Facade | `Composer`, `New()`, `Render()`, `Compile()` |
| `tests/e2e` | Real-FFmpeg E2E Test Suite | `TestEndToEnd_FullCompositionSuiteTable`, `TestCLI_EndToEndSuiteTable` |
| `examples` | Concrete Recipes | `01` through `10` |

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

### Running Real E2E Tests:
```bash
go test -race -v ./tests/e2e/...
```

### Running Example Recipes:
```bash
go run ./examples/01_simple_cut/main.go
go run ./examples/02_picture_in_picture/main.go
go run ./examples/03_transitions/main.go
go run ./examples/04_animated_motion/main.go
go run ./examples/05_animated_subtitles/main.go
go run ./examples/06_audio_ducking/main.go
go run ./examples/07_platform_presets/main.go
go run ./examples/08_podcast_audio_waveform/main.go
go run ./examples/09_green_screen_studio/main.go
```
