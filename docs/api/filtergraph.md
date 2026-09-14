# Filtergraph IR Package (`filtergraph`)

The `filtergraph` package is the core Directed Acyclic Graph (DAG) intermediate representation.

## Core Concepts

- **`Graph`**: Contains all filter nodes and tracks output pads.
- **`Node`**: Represents a single FFmpeg filter (e.g. `scale`, `overlay`, `split`, `volume`).
- **`Pad`**: Represents an input or output terminal. Pads are strictly typed as `StreamTypeVideo` or `StreamTypeAudio`.

```go
import "github.com/farshidrezaei/vidonex/filtergraph"

graph := filtergraph.New()

nodeA := graph.AddNode("scale", map[string]string{"w": "1920", "h": "1080"})
nodeB := graph.AddNode("format", map[string]string{"pix_fmts": "yuv420p"})

// Type-safe connection
err := graph.Connect(nodeA.Outputs()[0], nodeB.Inputs()[0])
```

---

## Compiler Passes

### `AutoSplitPass`
Traverses the graph, counts consumers of each output pad, and automatically injects `split` or `asplit` filters if any output pad is connected to more than one downstream input pad.

```go
pass := filtergraph.NewAutoSplitPass()
err := pass.Run(graph)
```

### `DeadCodeEliminationPass`
Removes dangling nodes, orphan pads, and subgraphs that do not lead to the final output pads.

```go
pass := filtergraph.NewDeadCodeEliminationPass()
err := pass.Run(graph)
```
