# Declarative Project Specification

Vidonex allows you to declare video compositions entirely using YAML or JSON. Relative asset paths are automatically resolved relative to the location of the specification file.

## Schema Example

```yaml
version: "1.0"

canvas:
  width: 1920
  height: 1080
  background_color: "#000000"

fps: 60

tracks:
  # Base Video Track
  - id: "track_main"
    kind: "video"
    z_index: 0
    clips:
      - id: "clip_intro"
        source: "assets/intro.mp4"
        offset: 0s
        duration: 8s
        trim_start: 2s
        volume: 1.0

      - id: "clip_main"
        source: "assets/main_footage.mp4"
        offset: 8s
        duration: 20s
        transition:
          type: "crossfade"
          duration: 1s

  # Picture-in-Picture Facecam Track
  - id: "track_overlay"
    kind: "overlay"
    z_index: 1
    clips:
      - id: "clip_webcam"
        source: "assets/webcam.mp4"
        offset: 2s
        duration: 10s
        scale: 0.28
        position:
          x: 1350
          y: 60
        opacity: 0.95
        chromakey:
          color: "#00FF00"
          similarity: 0.25
          smoothness: 0.1
          despill: true

  # Audio & Soundtrack Track
  - id: "track_audio"
    kind: "audio"
    clips:
      - id: "clip_bgm"
        source: "assets/music.mp3"
        offset: 0s
        duration: 28s
        volume: 0.35
        ducking:
          sidechain_track: "track_main"
          threshold_db: -20
          ratio: 4
```

---

## Field Reference

### Canvas
- `width` *(integer)*: Canvas width in pixels (e.g. `1920` or `1080`).
- `height` *(integer)*: Canvas height in pixels (e.g. `1080` or `1920`).
- `background_color` *(string)*: Hex color code (default: `#000000`).

### Track Kinds
- `video`: Primary video layer (normalized to canvas size).
- `overlay`: Picture-in-picture or transparent graphics layer.
- `audio`: Background music, sound effects, or voiceover.
- `subtitle`: Burned-in caption tracks.

### Clip Properties
- `source` *(string)*: Relative or absolute path to media asset.
- `offset` *(duration)*: Placement time on timeline (e.g. `2s`, `500ms`).
- `duration` *(duration)*: Duration to render on timeline.
- `trim_start` *(duration)*: Internal source trim offset.
- `scale` *(float)*: Size scaling multiplier (default: `1.0`).
- `position` *({ x, y })*: Canvas pixel coordinates.
- `opacity` *(float)*: Alpha transparency `0.0` to `1.0`.
- `volume` *(float)*: Audio volume multiplier (`1.0` = 100%).
