# CLI Automation

The `vidonex` CLI is a single static binary designed for headless batch rendering, containerized cloud workers, and automated video generation.

## Installation

```bash
go install github.com/farshidrezaei/vidonex/cmd/vidonex@latest
```

---

## Commands

### `vidonex render`

Renders a declarative project file (YAML or JSON) to an output video container.

```bash
vidonex render <project-file> -o <output-path> [flags]
```

**Key Flags:**
- `-o, --output <string>`: Destination output path (e.g. `output.mp4`).
- `--gpu <string>`: Hardware acceleration engine: `auto`, `nvenc` (NVIDIA), `videotoolbox` (Apple), `qsv` (Intel), `vaapi` (Linux), `none`. Default: `auto`.
- `--dry-run`: Validates the project and prints the generated FFmpeg CLI command without executing the render.
- `--overwrite`: Overwrites destination output file if it already exists.

```bash
# Example: Render with NVIDIA NVENC
vidonex render project.yaml -o final.mp4 --gpu nvenc

# Example: Inspect FFmpeg command line
vidonex render project.yaml --dry-run
```

---

### `vidonex serve`

Starts the embedded SQLite persistence layer, REST API, and WebSocket server for the Web Studio workstation.

```bash
vidonex serve --port 8080 --host 0.0.0.0
```

- **Web UI**: Open `http://localhost:8080` in your browser.
- **WebSocket Endpoint**: `ws://localhost:8080/ws` for live telemetry and project sync.
- **REST Endpoints**: `/api/projects`, `/api/media/probe`, `/api/render`.

---

### `vidonex graph`

Compiles the project and exports the filtergraph intermediate representation (DAG) to standard visual formats.

```bash
# Export Mermaid diagram (for Markdown, Notion, GitHub)
vidonex graph project.yaml --format mermaid -o filtergraph.mmd

# Export Graphviz DOT format
vidonex graph project.yaml --format dot -o filtergraph.dot
```

---

### `vidonex validate`

Checks project syntax, schema rules, track constraints, and verifies that all referenced media files exist on the filesystem.

```bash
vidonex validate project.yaml
```

---

### `vidonex probe`

Inspects audio and video streams inside a container using FFprobe, with optional structured JSON output.

```bash
vidonex probe assets/video.mp4 --json
```

---

### `vidonex about`

Displays comprehensive system environment diagnostics, installed engine version, detected FFmpeg/FFprobe binaries and versions, active hardware accelerators (e.g. NVIDIA NVENC, Apple VideoToolbox, VAAPI, Intel QSV), documentation resources, and an instant release check against GitHub.

```bash
vidonex about
```

---

### `vidonex upgrade`

Checks and automatically updates the local `vidonex` binary directly from GitHub Releases.

```bash
# Check if a newer version is available without installing
vidonex upgrade --check

# Download and perform atomic in-place binary upgrade
vidonex upgrade

# Force reinstall or downgrade to the latest release
vidonex upgrade --force
```
