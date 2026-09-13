# Recipe 10: Declarative YAML Project & Vidonex CLI

This directory demonstrates a complete, portable, standalone video editing project defined purely in declarative YAML (`vidonex.yaml`).

---

## 📁 Project Structure

```
.
├── vidonex.yaml      # Declarative project specification
├── main.go           # Go programmatic execution runner
└── README.md         # Documentation and CLI usage instructions
```

---

## 🚀 Running via Vidonex CLI

You can validate, graph, and render this project directly from your terminal using the `vidonex` command-line tool:

### 1. Validate Specification
```bash
vidonex validate vidonex.yaml
```

### 2. Export Filtergraph Flowchart
```bash
vidonex graph vidonex.yaml --format mermaid
```

### 3. Dry-Run Compilation
```bash
vidonex render vidonex.yaml --dry-run
```

### 4. Render Video with Hardware Acceleration & Live Progress
```bash
vidonex render vidonex.yaml -o final.mp4 --log-level info
```

---

## 💻 Running Programmatically in Go

```bash
go run ./examples/10_declarative_yaml_project/main.go
```
