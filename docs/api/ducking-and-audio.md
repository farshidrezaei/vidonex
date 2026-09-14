# Audio & Sidechain Ducking (`ducking`)

The `ducking` package allows you to automatically attenuate background music or sound effects whenever dialogue or voiceover is present.

## Quick Example

```go
import "github.com/farshidrezaei/vidonex/ducking"

options := ducking.Options{
    ThresholdDB: -25.0, // Trigger ducking when voice reaches -25dB
    Ratio:       4.0,   // Compression ratio
    AttackMS:    20.0,  // Fast attack to instantly dip music
    ReleaseMS:   400.0, // Smooth recovery fade back
}

// Connect background audio pad and speech lead pad
err := ducking.ApplySidechainDucking(graph, musicPad, voicePad, options)
```

## Options Reference

| Parameter | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `ThresholdDB` | `float64` | `-20.0` | Decibel volume at which ducking activates |
| `Ratio` | `float64` | `4.0` | Compression factor applied to background audio |
| `AttackMS` | `float64` | `20.0` | Time in ms to lower volume once voice begins |
| `ReleaseMS` | `float64` | `300.0` | Time in ms to restore volume once voice stops |
| `MakeupGain` | `float64` | `1.0` | Post-compression gain multiplier |
