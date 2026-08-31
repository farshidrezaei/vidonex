# Recipe 10: Declarative YAML Project & Vidonyx CLI

This directory demonstrates a complete, portable, standalone video editing project defined purely in declarative YAML (`vidonyx.yaml`).

---

## 📁 Project Structure

```
.
├── vidonyx.yaml      # Declarative project specification
├── main.go           # Go programmatic execution runner
└── README.md         # Documentation and CLI usage instructions
```

---

## 🚀 Running via Vidonyx CLI

You can validate, graph, and render this project directly from your terminal using the `vidonyx` command-line tool:

### 1. Validate Specification
```bash
vidonyx validate vidonyx.yaml
```

### 2. Export Filtergraph Flowchart
```bash
vidonyx graph vidonyx.yaml --format mermaid
```

### 3. Dry-Run Compilation
```bash
vidonyx render vidonyx.yaml --dry-run
```

### 4. Render Video with Hardware Acceleration & Live Progress
```bash
vidonyx render vidonyx.yaml -o final.mp4 --log-level info
```

---

## 💻 Running Programmatically in Go

```bash
go run ./examples/10_declarative_yaml_project/main.go
```
