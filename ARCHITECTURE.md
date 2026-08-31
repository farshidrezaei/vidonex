# Vidonyx Architecture & Deep Design Guide

This document provides a comprehensive technical overview of the internal architecture, design principles, intermediate representation (IR), and execution pipeline of the **Vidonyx** video composition engine.

---

## 1. System Overview & Architectural Pipeline

Vidonyx uses a **Multi-Stage Compiler Pipeline** to translate high-level declarative video compositions into valid, highly-optimized FFmpeg filtergraphs:

```
┌─────────────────────────────────────────────────────────────┐
│ 1. Declarative AST Layer (Timeline, Tracks, Clips, Effects) │
└──────────────────────────────┬──────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│ 2. Semantic Analysis & Normalization Pass                   │
└──────────────────────────────┬──────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│ 3. Filtergraph IR (Directed Acyclic Graph - DAG)            │
└──────────────────────────────┬──────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│ 4. Graph Optimization Passes (AutoSplit, DCE, Type Checking)│
└──────────────────────────────┬──────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│ 5. Target Code Generation (FFmpeg CLI / Mermaid / DOT)      │
└──────────────────────────────┬──────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│ 6. Process Execution & Real-time Progress Telemetry         │
└──────────────────────────────┬──────────────────────────────┘
```

---

## 2. Layer-by-Layer Breakdown

### 2.1. Fundamental Types (`types/`)

- **`types.Rational` (`Num/Den int64`)**:
  - Floating-point calculations drift when rendering thousands of frames across hours of video. Vidonyx strictly relies on fractional arithmetic for timecode, frame rate, and aspect ratio calculations.
  - Standard frame rate definitions (`FPS24`, `FPS30`, `FPS60`, `FPS23_976`, `FPS29_97`, `FPS59_94`) are exact integer fractions.
- **`types.Size` & `types.Rect`**:
  - Enforces even-dimension rounding (`makeEven`) on all scale operations, guaranteeing compatibility with `yuv420p` chroma subsampling and H.264/HEVC encoders.
- **`types.Color`**:
  - Parses Hex (`#RGB`, `#RRGGBB`, `#RRGGBBAA`) and serializes to FFmpeg expressions (`0xRRGGBB`, `0xRRGGBB@A.AA`, `black@0.0`).

---

### 2.2. Declarative Timeline AST (`timeline/`)

- **`Timeline`**: Defines canvas bounds (`Width`, `Height`), `FPS`, background color, duration, and track stacks.
- **`Track`**: Layer grouping categorized by `TrackKind` (`TrackKindVideo`, `TrackKindAudio`, `TrackKindOverlay`, `TrackKindText`) with explicit Z-index rendering priority.
- **`Clip`**: Media slices with in-point source trims (`SourceStart`), timeline placement (`TimelineStart`, `Duration`), playback rate multiplier (`Speed`), opacity (`Opacity`), and position transform (`Position`, `Scale`, `Alignment`).
- **Semantic Validator (`timeline.Validate`)**:
  - Uses `errors.Join` to collect all invalid states simultaneously (e.g. odd canvas dimensions, zero FPS, negative durations, out-of-bound opacity).

---

### 2.3. Filtergraph Intermediate Representation (`filtergraph/`)

The Filtergraph package is a standalone graph library modeling FFmpeg filter topologies as a Directed Acyclic Graph (DAG):

- **`Node`**: Represents a single FFmpeg filter (e.g. `scale`, `overlay`, `amix`, `xfade`, `trim`).
- **`Pad`**: Represents typed input and output pins (`StreamTypeVideo`, `StreamTypeAudio`).
- **`Graph`**:
  - Enforces stream type compatibility (e.g. audio pads cannot connect to video inputs).
  - Prevents graph cycles on connection using depth-first search (DFS).
  - Employs **Kahn's Algorithm** for topological sorting to serialize the filter execution order deterministically.
  - Provides zero-allocation Go iterators (`Nodes() iter.Seq[*Node]`).

---

### 2.4. Optimization Passes (`filtergraph/pass*.go`)

1. **`AutoSplitPass`**:
   - In FFmpeg, a filter output pad can only be consumed once. If an output pad is connected to $N > 1$ input pads, this pass automatically injects a `split` (for video) or `asplit` (for audio) node with $N$ outputs and rewires downstream links.
2. **`DeadCodeEliminationPass` (DCE)**:
   - Traverses backwards from designated sink pads (`[out_v]`, `[out_a]`) using DFS, pruning disconnected subgraphs and orphan nodes before CLI emission.
3. **`ValidatePass`**:
   - Verifies graph acyclicity, stream type matching, and ensures all internal input pads are connected.

---

### 2.5. Multi-Stage Compiler (`compiler/`)

1. **Input Indexing**: Discovers all media sources across clips and generates optimal `-i <path>` CLI flags with stable indices (`0:v`, `0:a`, `1:v`, etc.).
2. **Clip Video Pipeline**: Injects `trim`, `setpts`, `format=yuva420p` + `colorchannelmixer` (for opacity), and `scale` (for transform scaling).
3. **Stream Normalization**: Injects `scale=w:h:force_original_aspect_ratio=decrease`, `pad`, `setsar=1`, `fps`, and `format=yuv420p` for video; and `aresample=48000`, `aformat=sample_fmts=fltp:channel_layouts=stereo` for audio.
4. **Compositor**: Stacks visual layers on top of a synthesized background canvas (`color=c=...:s=...:r=...:d=...`) using `overlay` filters enabled during active clip time intervals (`enable='between(t, start, end)'`).
5. **Audio Mixer**: Aligns audio clips using `adelay` and combines all streams into a master stereo mix via `amix`.

---

### 2.6. Execution & Telemetry (`executor/`)

- **`CommandExecutor` Interface**: Abstracts binary process invocation.
- **`OSExecutor`**: Invokes FFmpeg via `os/exec.CommandContext`, attaching `-progress pipe:1` to stream real-time progress without blocking.
- **`ProgressParser`**: Line-by-line streaming state machine that parses FFmpeg progress key-value pairs (`frame=`, `fps=`, `out_time_us=`, `speed=`) into structured `ProgressEvent` structs.
- **`MockExecutor`**: Records commands and simulates progress ticks for deterministic sub-millisecond unit tests in CI/CD without requiring external binaries.

---

### 2.7. Visualizers (`visualizer/`)

- **Mermaid.js Generator (`ToMermaid`)**: Outputs standard GitHub/Markdown flowchart syntax.
- **Graphviz DOT Generator (`ToDOT`)**: Outputs Graphviz `.dot` files for high-resolution visual layout rendering.
