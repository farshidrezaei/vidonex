# ChromaKey & Despill (`chromakey`)

The `chromakey` package isolates colored backdrops (green or blue screens) and suppresses green spill reflections on skin and clothing.

## Quick Example

```go
import "github.com/farshidrezaei/vidonex/chromakey"

options := chromakey.Options{
    KeyColor:   types.HexColor("#00FF00"),
    Similarity: 0.25, // Tolerance for color variation
    Smoothness: 0.10, // Edge feathering
    Despill:    true, // Remove green reflections from actor
}

// Injects colorkey and despill nodes into filtergraph DAG
outputPad, err := chromakey.ApplyChromaKey(graph, inputPad, options)
```

## Studio Presets

Vidonex includes tested studio presets:
- `chromakey.StudioGreenScreen`: Calibrated for well-lit studio backdrops.
- `chromakey.StudioBlueScreen`: Calibrated for blue cycloramas.
