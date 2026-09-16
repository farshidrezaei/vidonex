# Desktop Studio Workstation

**Vidonex Studio** is a modern, lightweight non-linear video editing workstation engineered with **Wails v2**, **Vue 3**, **Nuxt UI**, and **Tailwind CSS**.

It provides the fluid tactile experience of Premiere Pro and CapCut without the memory overhead and lag of Electron-based apps.

<p align="center" style="margin-top: 1.5rem; margin-bottom: 2rem;">
  <img src="/vidonex-studio-demo.gif" alt="Vidonex Studio Live Demo" style="border-radius: 8px; box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.5); max-width: 100%;" />
</p>

---

## Workspace Layout

```
┌────────────────────────────────────────────────────────────────────────┐
│ Top Bar: Project Title • Canvas Dimensions • Undo/Redo • Export Button │
├───────────────┬───────────────────────────────────┬────────────────────┤
│ Media Library │ Interactive Viewport              │ Inspector Panel    │
│               │ - Transform Gizmo (8-point)       │ - Transform        │
│ - Videos      │ - Canvas Zoom & Pan               │ - ChromaKey        │
│ - Audios      │ - Alignment Guides                │ - Audio Ducking    │
│ - Subtitles   │ - Live Frame Preview              │ - Waveform FX      │
│               │                                   │ - Subtitle Styler  │
├───────────────┴───────────────────────────────────┴────────────────────┤
│ Multi-Track Timeline                                                   │
│ [Play/Pause] [Split S] [Snapping] [Zoom Slider]                        │
│ ── Track 1: Overlays / PIP ─────────────────────────────────────────── │
│ ── Track 2: Main Video ─────────────────────────────────────────────── │
│ ── Track 3: Voiceover / Ducking Lead ───────────────────────────────── │
│ ── Track 4: Background Music ───────────────────────────────────────── │
└────────────────────────────────────────────────────────────────────────┘
```

---

## Core Features

### 🎞️ Multi-Track Timeline
- **Track Stacking**: Video, audio, subtitle, overlay, and waveform tracks.
- **Precision Trimming**: Drag clip edges with ripple or roll edit semantics.
- **Clip Splitting**: Place the playhead and hit `S` to instantly split a clip.
- **Snapping**: Magnetic alignment to playhead and adjacent clips.
- **Crossfades**: Visually drag transition ribbons between adjacent video or audio clips.

### 📐 Transform Gizmo
- On-canvas 8-point interactive resize handles.
- Center rotational pivot.
- Pixel-precision arrow nudging (`1px` with arrows, `10px` with `Shift + Arrow`).

### 🎛️ Built-in FX Suite
- **ChromaKey**: Visual color picker to isolate green screens with real-time despill suppression.
- **Audio Ducking**: Graphically select dialogue tracks to automatically lower background music volume.
- **Animated Audiograms**: Real-time waveform rendering with customizable neon gradients, bars, and symmetry modes.
- **SRT / VTT Burn-in**: Custom fonts, shadow offsets, background cards, and yellow karaoke highlights.

### 🔔 Updates & System Diagnostics
- **One-Click In-App Updates**: When a new release is available on GitHub, the Studio alerts the user via a top banner and header notification dot. Users can upgrade automatically with live download progress streaming and atomic replacement.
- **About Modal**: Inspect engine version, operating system architecture, CPU core count, active GPU hardware acceleration (NVENC, VideoToolbox, VA-API, QSV), and detected FFmpeg binaries directly within the application.

