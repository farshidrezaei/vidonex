# Vidonex Architecture & Deep Design Guide

This document provides a comprehensive technical overview of the internal architecture, design principles, intermediate representation (IR), and execution pipeline of the **Vidonex** video composition engine.

---

## 1. System Overview & Architectural Pipeline

Vidonex uses a **Multi-Stage Compiler Pipeline** to translate high-level declarative video compositions into valid, highly-optimized FFmpeg filtergraphs:

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
  - Floating-point calculations drift when rendering thousands of frames across hours of video. Vidonex strictly relies on fractional arithmetic for timecode, frame rate, and aspect ratio calculations.
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

1. **Input Indexing**: Discovers all media sources across clips and generates optimal `-i <path>` CLI flags with stable indices (`0:v`, `0:a`, `1:v`, etc.), applying `-loop 1` for static images.
2. **Clip Video Pipeline**: Injects `trim`, `setpts` (with timeline start offset `PTS-STARTPTS+<TimelineStart>/TB` to cleanly preserve timeline gaps), `format=yuva420p` + `colorchannelmixer` (for opacity), and `scale` (for transform scaling).
3. **Stream Normalization**: Injects `format=yuva420p`, `scale=w:h:force_original_aspect_ratio=decrease`, `pad=color=black@0.0`, `setsar=1`, `fps`, and `format=yuva420p` to maintain 100% alpha transparency for overlays and transparent PNGs.
4. **Compositor**: Stacks visual layers on top of a synthesized background canvas (`color=c=...:s=...:r=...:d=...`) using `overlay` filters enabled during active clip time intervals (`enable='between(t, start, end)'`), finishing with a final `format=yuv420p` node for standard MP4 player compatibility.
5. **Audio Mixer**: Aligns audio clips using `adelay` and combines all streams into a master stereo mix via `amix`.

---

### 2.6. Embedded Server & Web Studio Workstation (`server/` & `ui/`)

1. **Embedded Persistence (`server/db/`)**: High-performance SQLite persistence (`mattn/go-sqlite3`) managing projects, timeline specifications, assets, and render job queue states.
2. **RESTful API & WebSocket Hub (`server/api/` & `server/ws/`)**: Exposes project management, live asset probing (`ffprobe`), file uploading, and real-time WebSocket progress broadcasting.
3. **Web Studio Workstation (`ui/`)**:
   - **Nuxt 4 + Nuxt UI v4 (Vue 3, Pinia, Tailwind CSS)**: Modern non-linear video editing workstation.
   - **Customizable Splitter Layout**: Dynamic resizable panes powered by Nuxt UI `USplitter` separating the Media Library, Inspector, Canvas Viewport, and Timeline.
   - **Pro Multi-Track Timeline**: Features clip thumbnail filmstrips, audio waveform SVG textures, sticky time ruler, scrubbing needle line, split clip (`S`), duplicate (`Ctrl+D`), and adjacent clip crossfade (`XFade`) visual ribbons.
   - **Interactive Viewport with Hand Tool & Zoom**: Direct on-canvas 8-point resize, position dragging, rotation, dynamic aspect-ratio fitting, smooth wheel zoom, Spacebar hand pan tool, and fullscreen mode.
   - **Atomic Snapshot History & Exact Duration**: Bulletproof Undo/Redo stack with automatic synchronization of timeline and playback duration matching exact content bounds.
   - **Graph Visualizer**: Interactive Mermaid.js DAG inspection modal with SVG export and refresh triggers.

---

### 2.7. Execution & Telemetry (`executor/`)

- **`CommandExecutor` Interface**: Abstracts binary process invocation.
- **`OSExecutor`**: Invokes FFmpeg via `os/exec.CommandContext`, attaching `-progress pipe:1` to stream real-time progress without blocking.
- **`ProgressParser`**: Line-by-line streaming state machine that parses FFmpeg progress key-value pairs (`frame=`, `fps=`, `out_time_us=`, `speed=`) into structured `ProgressEvent` structs.
- **`MockExecutor`**: Records commands and simulates progress ticks for deterministic sub-millisecond unit tests in CI/CD without requiring external binaries.

---

### 2.8. Visualizers (`visualizer/`)

- **Mermaid.js Generator (`ToMermaid`)**: Outputs standard GitHub/Markdown flowchart syntax.
- **Graphviz DOT Generator (`ToDOT`)**: Outputs Graphviz `.dot` files for high-resolution visual layout rendering.

---

### 2.9. Nuxt 4 Web Studio Workstation Details (`ui/`)

- **Framework & Tooling**: Nuxt 4, Vue 3 Composition API, Nuxt UI v4, Nuxt i18n (LTR English & RTL Persian), Pinia state stores, and Tailwind CSS.
- **Pro Multi-Track Timeline**: Dynamic video/audio/overlay/subtitle/waveform tracks, clip dragging, In/Out trimming, splitting (`S`), adjacent clip crossfades (`XFade`), sticky time ruler, and magnetic snapping.
- **Interactive Viewport & Transform Gizmo**: Direct on-canvas 8-point resize, translation, rotation, smart guidelines, Spacebar hand tool pan, wheel zoom, zoom presets, and instant 0ms latency rendering.
- **Domain Inspectors**: Dedicated controllers for Chroma Key, Sidechain Ducking, Podcast Waveforms, SRT/VTT Subtitles, Ken Burns Keyframes, and Mermaid DAG visualization.
- **Persistence & Synchronization**: Auto-saving project spec to embedded SQLite, aspect-ratio presets (16:9, 9:16, 1:1, 21:9), and exact media asset uploads.
