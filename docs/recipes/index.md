# Runnable Recipes Catalog

Explore tested, production-ready recipes demonstrating Vidonex capabilities.

| Recipe | Description | Source Code |
| :--- | :--- | :--- |
| **[Picture-in-Picture](/recipes/picture-in-picture)** | Overlay facecam video with border and opacity | `examples/02_picture_in_picture` |
| **[Sidechain Audio Ducking](/recipes/audio-ducking)** | Automatically lower music volume during commentary | `examples/06_audio_ducking` |
| **[Animated Waveforms](/recipes/animated-waveforms)** | Neon pulsing audiograms for podcasts & reels | `examples/08_podcast_audio_waveform` |
| **[Green Screen Studio](/recipes/green-screen)** | ChromaKey removal with despill onto virtual studio | `examples/09_green_screen_studio` |
| **[Burn-in Subtitles](/recipes/subtitles)** | Styled SRT captions with highlight boxes | `examples/05_animated_subtitles` |

---

## Running Recipes Locally

All recipes are directly runnable from the repository root:

```bash
go run ./examples/02_picture_in_picture/main.go
go run ./examples/06_audio_ducking/main.go
go run ./examples/08_podcast_audio_waveform/main.go
go run ./examples/09_green_screen_studio/main.go
```
