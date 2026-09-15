# Changelog

All notable changes to the **Vidonex** project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Planned for v1.4.0
- **Interactive OpenAPI & Scalar Documentation (`/docs`)**: Embedded interactive REST sandbox directly in `vidonex serve`.
- **Compiler Performance Benchmark Suite**: Memory allocation and throughput benchmark tests comparing Go DAG compiler against industry alternatives.
- **GitHub Action Workflow (`farshidrezaei/vidonex-action`)**: Reusable CI action for automated video composition and changelog renders.
- **DevContainer Environment**: Turnkey containerized workspace with Go, Node, Wails, and FFmpeg pre-configured.

---

## [1.3.0] - 2026-09-15

### Next-Gen Codecs & WebAssembly Browser Engine Release 🌐

Vidonex v1.3.0 introduces production export profiles for next-generation video formats (AV1, ProRes 4444 with Alpha, WebM) and compiles the core declarative composition engine into WebAssembly (`vidonex.wasm`) for pure client-side execution in web browsers.

### Added

#### 🎞️ Next-Generation Video Codec Profiles (`presets/codecs.go`)
- **`ProfileAV1()`**: Turnkey SVT-AV1 (`libsvtav1`) encoding configuration featuring 10-bit color (`yuv420p10le`), Opus audio, and tune flags delivering 35% higher compression than H.264/HEVC.
- **`ProfileProRes4444()`**: Apple ProRes 4444 configuration with 10-bit color and alpha channel preservation (`yuva444p10le`) for visual effects workflows and transparent overlays.
- **`ProfileProRes422HQ()`**: Apple ProRes 422 HQ broadcast mastering preset with uncompressed 24-bit PCM audio.
- **`ProfileWebMVP9()` & `ProfileWebMAlpha()`**: High-compatibility open web delivery presets supporting standard video and transparent overlay alphas.
- Table-driven unit test suite in [`presets/codecs_test.go`](file:///home/farshid/Desktop/projects/vidonyx/presets/codecs_test.go).

#### 🌐 WebAssembly (WASM) In-Browser Pipeline (`cmd/wasm/main.go`)
- **Standalone Go WebAssembly Engine**: Compiled with `GOOS=js GOARCH=wasm` into a 5.2MB self-contained WebAssembly binary (`ui/public/vidonex.wasm`).
- Exposes `globalThis.vidonex` to JavaScript:
  - `vidonex.validate(specString)`: Performs full AST validation and returns `{ valid, duration, tracks, error }`.
  - `vidonex.compile(specString, outputPath)`: Generates the full FFmpeg filter complex, CLI arguments, and execution commands with zero server roundtrips.
  - `vidonex.mermaid(specString)`: Generates real-time Mermaid.js DAG flowchart code directly in-browser.
  - `vidonex.version()`: Returns engine version string (`"1.3.0"`).
- **Nuxt Composable (`useWasmCompiler.ts`)**: Reactive WebAssembly loader and typed bridge for client-side timeline compilation.
- **Build Target**: Added `make wasm` to the root [`Makefile`](file:///home/farshid/Desktop/projects/vidonyx/Makefile).

---

## [1.2.0] - 2026-09-15

### Broadcast Audio & Cinematic Color Grading Release 🎨

Vidonex v1.2.0 elevates audio and video fidelity with industry-standard broadcast loudness normalization, 3D Look-Up Table (LUT) color grading, an expanded cinematic transition library, and turnkey social media templates.

### Added

#### 🎚️ Broadcast & Streaming Audio Normalization (`effects/loudnorm.go`)
- **EBU R128 `LoudnormFilter`**:
  - Implements full FFmpeg `loudnorm` filter integration for broadcast and streaming audio loudness compliance.
  - Dedicated industry presets:
    - `LoudnormSpotifyYouTube()`: -14.0 LUFS integrated loudness, 11.0 LU LRA, -1.0 dBFS true peak.
    - `LoudnormEBUR128()`: -23.0 LUFS European broadcast standard, 7.0 LU LRA, -1.0 dBFS true peak.
    - `LoudnormPodcast()`: -16.0 LUFS vocal speech optimization, 9.0 LU LRA, -1.5 dBFS true peak.
    - `LoudnormAppleMusic()`: -16.0 LUFS Apple ecosystem preset, 11.0 LU LRA, -1.0 dBFS true peak.
  - Full table-driven unit test coverage (`effects/loudnorm_test.go`).

#### 🎨 Cinematic Color Grading & 3D LUTs (`effects/color.go`)
- **`LUT3DFilter`**:
  - Support for industry-standard `.cube` and `.3dl` Look-Up Tables with configurable tetrahedral, trilinear, or nearest-neighbor interpolation.
- **`ColorBalanceFilter`**:
  - Independent shadow, midtone, and highlight RGB color balance adjustments (`rs`, `gs`, `bs`, `rm`, `gm`, `bm`, `rh`, `gh`, `bh`).
- **`ColorGradingFilter`**:
  - Unified color pipeline controlling contrast, brightness, saturation, gamma, temperature (warm/cool), and tint (green/magenta) with automatic filter graph node chaining.
  - Comprehensive table-driven unit tests (`effects/color_test.go`).

#### 🎬 50+ Extended XFade Cinematic Transitions Catalog (`timeline/transition.go`)
- Expanded `TransitionType` constants to support the complete FFmpeg `xfade` matrix:
  - Directional wipes: `wipeleft`, `wiperight`, `wipeup`, `wipedown`, `wipetl`, `wipetr`, `wipebl`, `wipebr`.
  - Directional slides & smooths: `slideleft`, `slideright`, `slideup`, `slidedown`, `smoothleft`, `smoothright`, `smoothup`, `smoothdown`.
  - Shapes & Iris: `circlecrop`, `circleopen`, `circleclose`, `rectcrop`, `radial`, `distance`, `horzopen`, `horzclose`, `vertopen`, `vertclose`, `hlslice`, `vuslice`, `zoomin`, `squeezev`, `squeezeh`.
  - Fades & textures: `fade`, `fadeblack`, `fadewhite`, `fadegrays`, `dissolve`, `pixelize`.
  - Catalog query helpers: `AllVideoTransitions()` and `IsValidTransitionType()`.

#### 📱 Turnkey Social Media Production Templates (`presets/templates.go`)
- `TemplateTikTokSplitScreen()`: 9:16 vertical 1080x1920 layout with top gameplay/reaction half and bottom creator facecam half.
- `TemplatePodcastAudiogram()`: 1:1 square layout combining background visuals, dynamic neon audio waveforms, and vocal tracks.
- `TemplateYouTubeEndScreen()`: 16:9 1080p outro sequence with dedicated video and subscribe card zones.

---

## [1.1.0] - 2026-09-15

### Studio Ergonomics & UX Powerhouse Release 🚀

Vidonex v1.1.0 focuses on workstation responsiveness, tactile editing controls, and professional NLE ergonomics.

### Added

#### ⚡ Studio Workstation & Editing Ergonomics
- **Studio Command Palette (`Ctrl+K` / `Cmd+K`)**:
  - Searchable command palette built with `@nuxt/ui` v4 `UCommandPalette` and `UModal`.
  - Instant access to timeline operations (Split Clip `S`, Duplicate `Ctrl+D`, Delete `Del`, Toggle Snapping `N`), track additions, canvas aspect ratios (16:9, 9:16, 1:1, 21:9), playback transport controls, and modals.
  - Dedicated Command Palette button and keyboard shortcut indicator in the workstation header.
- **Live Stereo VU / Audio Peak Meter (`AudioPeakMeter.vue`)**:
  - Real-time dual-channel (L/R) peak level indicator displaying dynamic green (-60 dB to -12 dB), yellow (-12 dB to -3 dB), and red clip warnings.
  - Responsive Peak Hold ticks and live numeric dBFS telemetry integrated directly into the playback transport bar.
- **Dynamic Audio Waveform Visualization (`ClipWaveform.vue`)**:
  - Multi-sample responsive audio peaks drawn via HTML5 canvas with symmetrical dual-polarity envelopes and track-themed gradients (`emerald` for audio, `purple` for waveforms).
  - Automatically updates when clips are moved, trimmed, or resized.
- **Professional Repeating Filmstrip on Video Clips**:
  - Tiled background filmstrip with periodic frame dividers and border vignette, matching industry-standard NLEs (DaVinci Resolve, Premiere Pro, Final Cut).
- **Interactive Bezier & Easing Curve Graph Editor (`CurveEditor.vue`)**:
  - Interactive curve visualizer in the Keyframe Inspector displaying animated trajectory acceleration curves (Linear, EaseInQuad, EaseOutQuad, EaseInOutQuad, EaseInOutCubic) with animated preview play button.
- **Bilingual Internationalization**: Full English and Persian (`fa`) translations for all Command Palette commands, actions, and curve visualization tools.

---

## [1.0.0] - 2026-09-14

### Initial Stable Release 🎉

Vidonex v1.0.0 is the foundational release of the declarative video composition engine, FFmpeg filtergraph intermediate representation compiler, and modern non-linear desktop workstation.

### Added

#### ⚡ Core Engine & Compiler (`types`, `timeline`, `filtergraph`, `compiler`)
- **Declarative AST**: High-level timeline model supporting multi-layered video, audio, overlay, subtitle, and waveform tracks.
- **Filtergraph Directed Acyclic Graph (DAG)**: Type-safe intermediate representation with strict `StreamTypeVideo` and `StreamTypeAudio` pad typing.
- **Compiler Passes**:
  - `AutoSplitPass`: Statically analyzes pad consumer references and transparently injects `split` and `asplit` nodes to eliminate FFmpeg pad reuse errors.
  - `DeadCodeEliminationPass`: Prunes unlinked branches, dormant filters, and unreachable tracks before argument emission.
  - **Alpha Preservation & Normalization**: Preserves `yuva420p` alpha channels across overlays and chroma-key layers, finalizing with `yuv420p` encoding for universal playback.
- **Zero-Drift Rational Timing**: Implemented `types.Rational` exact fraction arithmetic across all timestamps and framerates, preventing floating-point frame drift.
- **Visualizer**: Graph export to Mermaid.js (`.mmd`) and Graphviz (`.dot`) formats.

#### 🎛️ Audio & Video Processing Suite
- **Sidechain Audio Ducking (`ducking`)**: Automated music attenuation keyed to speech commentary with customizable threshold, attack, release, and ratio curves.
- **Animated Audiograms (`waveform`)**: Pulsing neon audio waveforms for podcasts and social reels with `ModePeakToPeak`, logarithmic scaling, and custom dimensions.
- **ChromaKey & Despill (`chromakey`)**: Studio-grade green/blue screen isolation with color pickers, similarity tuning, smoothness feathering, and spill suppression.
- **Subtitles & Captions (`subtitles`)**: Full SRT and WebVTT parser with customizable fonts, font size, border outlines, vertical margins, and highlight boxes.
- **Transitions & Effects (`effects`)**: Native `xfade` video transitions (dissolve, wipe, slide), `acrossfade` audio transitions, scale, and text drawing.
- **Keyframe Motion & Animation (`animation`)**: 2D coordinate translation, scale, and opacity keyframe interpolation with Linear, Quad, Cubic, and Sine easing.

#### 🎬 Desktop Workstation & Web Studio (`desktop`, `ui`, `server`)
- **Modern Non-Linear Editor**: Engineered with Wails v2, Nuxt 4, Vue 3, Nuxt UI, and Tailwind CSS (<40MB baseline RAM).
- **Multi-Track Timeline**: Precise drag-and-drop, playhead scrubbing, split clips (`S`), trim handles, magnetic snapping, and crossfade ribbons.
- **Interactive Viewport**: On-canvas 8-point Transform Gizmo, rotational pivot, canvas zoom/pan, and arrow nudging (`1px` fine / `10px` coarse).
- **Categorized Media Library**: High-density TreeView asset browser for Videos, Audio, Images, and Subtitles with live search filtering.
- **Zero-Copy File Import & Native "Save As"**: Native OS file dialog integration for multi-gigabyte project exports.
- **Headless Server**: Embedded SQLite persistence, REST endpoints, and WebSocket telemetry server (`vidonex serve`).

#### 💻 Automation CLI (`cmd/vidonex`)
- `vidonex render`: Batch render YAML/JSON declarative projects with hardware acceleration flags.
- `vidonex serve`: Launch standalone web studio on any host/port.
- `vidonex validate`: Validate timeline AST constraints and media file existence.
- `vidonex graph`: Export DAG diagrams for CI/CD documentation.
- `vidonex probe`: Extract container streams, codecs, and durations as JSON.

#### 🚀 Hardware Acceleration (`presets`)
- Automatic hardware encoder detection and flag injection:
  - NVIDIA NVENC (`h264_nvenc`, `hevc_nvenc`)
  - Apple Silicon / macOS VideoToolbox (`h264_videotoolbox`, `hevc_videotoolbox`)
  - Intel QuickSync (`h264_qsv`, `hevc_qsv`)
  - Linux VAAPI (`h264_vaapi`, `hevc_vaapi`)

#### 🧪 Testing & CI
- Table-driven unit tests across all packages with race detector enabled (`-race`).
- Synthetic media end-to-end integration test suite (`tests/e2e`) verifying real FFmpeg and FFprobe execution.
- Automated multi-platform release pipeline building desktop and CLI binaries for Linux, macOS, and Windows.
