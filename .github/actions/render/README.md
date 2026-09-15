# Vidonex Video Render Action

Official GitHub Action for rendering declarative video timelines using the [Vidonex](https://github.com/farshidrezaei/vidonex) video composition engine and FFmpeg.

## Features

- **Automated Video CI/CD**: Compile and render dynamic programmatic videos directly on PR merges or scheduled cron workflows.
- **Hardware Acceleration Ready**: Supports NVENC, VAAPI, and Apple VideoToolbox runners.
- **Zero-Setup FFmpeg**: Automatically installs FFmpeg and FFprobe on Linux and macOS runners if not present.
- **Artifact Pipeline**: Exposes the rendered file path as a GitHub Actions step output for direct uploading with `actions/upload-artifact`.

## Usage Example

```yaml
name: Render Daily Reel

on:
  push:
    branches: [main]
  schedule:
    - cron: '0 0 * * *' # Midnight daily

jobs:
  render:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Repository
        uses: actions/checkout@v4

      - name: Render Video with Vidonex
        id: vidonex-render
        uses: farshidrezaei/vidonex/.github/actions/render@v1.4.0
        with:
          project: "spec/daily_reel.yaml"
          output: "exports/daily_reel.mp4"
          hwaccel: "none"

      - name: Upload Video Artifact
        uses: actions/upload-artifact@v4
        with:
          name: daily-reel-video
          path: ${{ steps.vidonex-render.outputs.rendered-file }}
```

## Inputs

| Input | Description | Required | Default |
| :--- | :--- | :--- | :--- |
| `project` | Path to Vidonex timeline specification file (`.yaml` or `.json`) | **Yes** | — |
| `output` | Destination file path for the rendered MP4/WebM video | **Yes** | — |
| `version` | Vidonex engine version | No | `1.4.0` |
| `hwaccel` | Hardware acceleration (`none`, `nvenc`, `vaapi`, `videotoolbox`) | No | `none` |
| `install-ffmpeg` | Automatically install FFmpeg if missing | No | `true` |

## Outputs

| Output | Description |
| :--- | :--- |
| `rendered-file` | Absolute path to the rendered video file |
