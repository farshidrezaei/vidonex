# Changelog

All notable changes, releases, and roadmap milestones for **Vidonex** are documented here.

This project strictly adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html) and [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

---

## [Unreleased] <Badge text="In Development" type="warning" />

### Roadmap for v1.1.0
- **Direct Viewport Acceleration**: Embedded MPV / OpenGL surface rendering for real-time 4K60 video playback.
- **WebAssembly (Wasm) In-Browser Pipeline**: Pure client-side timeline composition without needing a local server.
- **Multicam Switcher**: Synchronized multi-angle camera switching for podcasts and studio interviews.
- **LUT & 3D Color Grading Engine**: Import `.cube` 3D LUTs with real-time lift/gamma/gain wheels.

---

## [v1.0.0] - 2026-09-14 <Badge text="Latest Stable" type="tip" />

### Initial Production Release 🎉

Vidonex v1.0.0 is the official initial release of the declarative video composition engine, FFmpeg filtergraph intermediate representation compiler, and modern non-linear desktop workstation.

### Key Highlights

::: tip ⚡ Core Engine & Compiler
- **Declarative AST**: High-level timeline model supporting multi-layered video, audio, overlay, subtitle, and waveform tracks.
- **Filtergraph Directed Acyclic Graph (DAG)**: Type-safe intermediate representation with strict `StreamTypeVideo` and `StreamTypeAudio` pad typing.
- **`AutoSplitPass`**: Automatically detects output pads with multiple downstream consumers and transparently injects `split`/`asplit` nodes to eliminate FFmpeg pad reuse crashes.
- **`DeadCodeEliminationPass`**: Prunes dangling nodes, unlinked pads, and dead branches before emitting final CLI arguments.
- **Zero-Drift Rational Timing**: Implements `types.Rational` exact fractions across all timestamps, avoiding floating-point frame drift.
- **Format & Alpha Normalization**: Preserves transparency channels (`yuva420p`) across overlays, finalizing with universally compatible `yuv420p` video.
:::

::: tip 🎬 Desktop Workstation & Web Studio
- **Lightweight Wails v2 App**: Engineered with Vue 3, Nuxt 4, and Nuxt UI (<40MB baseline RAM).
- **Multi-Track Timeline**: Precise drag-and-drop, playhead scrubbing, split clips (`S`), trim handles, magnetic snapping, and crossfade ribbons.
- **Interactive Viewport**: 8-point Transform Gizmo, rotational pivot, canvas zoom/pan, and arrow nudging.
- **High-Density Media Library**: TreeView asset browser for Videos, Audio, Images, and Subtitles.
- **Zero-Copy File Import & Native "Save As"**: Native OS file dialog integration for multi-gigabyte project exports.
:::

::: tip 🎛️ Audio & Video Processing Suite
- **Sidechain Audio Ducking**: Automatically attenuates background music when voiceover is detected.
- **Animated Audiograms**: Generates dynamic neon waveforms from voice recordings.
- **ChromaKey & Despill**: Studio-grade green/blue screen removal with color pickers and despill filtering.
- **Styled Subtitles**: Burns in SRT/VTT captions with custom typography, outlines, and highlight cards.
- **Keyframe Motion**: 2D coordinate translation, scale, and opacity keyframes with mathematical easing curves.
:::

::: tip 💻 CLI Automation & GPU Acceleration
- **Single Static Binary**: `vidonex render`, `vidonex serve`, `vidonex validate`, `vidonex graph`, `vidonex probe`.
- **GPU Acceleration**: Auto-detection for NVIDIA NVENC, Apple VideoToolbox, Intel QSV, and VAAPI.
:::

---

## Release Links

- [GitHub Release Assets & Checksums](https://github.com/farshidrezaei/vidonex/releases/tag/v1.0.0)
- [Source Code at v1.0.0](https://github.com/farshidrezaei/vidonex/tree/v1.0.0)
