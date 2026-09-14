# Go SDK API Reference

Vidonex provides a type-safe Go SDK designed with zero external third-party dependencies in its core packages.

## Package Architecture

| Package | Import Path | Responsibility |
| :--- | :--- | :--- |
| **`types`** | `github.com/farshidrezaei/vidonex/types` | Exact math (`Rational`), geometry (`Point`, `Rect`, `Size`), and colors |
| **`timeline`** | `github.com/farshidrezaei/vidonex/timeline` | Declarative timeline AST: Timeline, Tracks, Clips, Transitions |
| **`filtergraph`** | `github.com/farshidrezaei/vidonex/filtergraph` | Intermediate representation: DAG, Nodes, Typed Pads, AutoSplitPass |
| **`compiler`** | `github.com/farshidrezaei/vidonex/compiler` | Multi-stage compiler: AST $\to$ DAG $\to$ FFmpeg CLI arguments |
| **`composer`** | `github.com/farshidrezaei/vidonex/composer` | High-level facade for compilation and real-time execution |
| **`ducking`** | `github.com/farshidrezaei/vidonex/ducking` | Sidechain audio compressor ducking options and compiler hooks |
| **`waveform`** | `github.com/farshidrezaei/vidonex/waveform` | Audiogram waveform visualizer modes and parameters |
| **`chromakey`** | `github.com/farshidrezaei/vidonex/chromakey` | ChromaKey green/blue screen isolation with despill |
| **`effects`** | `github.com/farshidrezaei/vidonex/effects` | Pre-built filter nodes (Scale, DrawText, XFade, Volume) |
| **`presets`** | `github.com/farshidrezaei/vidonex/presets` | Social media platform presets (TikTok 9:16, YouTube 4K) & GPU encoders |
| **`executor`** | `github.com/farshidrezaei/vidonex/executor` | Subprocess runner with live progress and event parsing |

---

## Engineering Contracts

1. **Exact Rational Math**: Never use `float64` for timestamp math or frame calculations. Always use `types.Rational`.
2. **Type Safety**: All filter pads are typed as `StreamTypeVideo` or `StreamTypeAudio`. Connecting mismatching types returns compile-time errors.
3. **No Concatenation**: Filtergraph strings are emitted only after structural validation and optimization passes.
