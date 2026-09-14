# Effects Package (`effects`)

The `effects` package provides pre-constructed, type-safe filter nodes for common video transformations.

## Pre-built Filters

### `ScaleFilter`
Scales video to specified dimensions with aspect ratio preservation options.
```go
filter := effects.NewScaleFilter(1920, 1080, effects.ScaleModeFit)
```

### `DrawTextFilter`
Draws stylized text onto a video stream with custom font, size, color, and box background.
```go
textFilter := effects.NewDrawTextFilter("Title Text", "/fonts/Inter.ttf", 48, types.White).
    WithBox(types.Black, 0.6)
```

### `VolumeFilter`
Adjusts audio gain linearly or logarithmically in decibels.
```go
volFilter := effects.NewVolumeFilter(0.5) // 50% volume
```

### `XFadeFilter`
Crossfades between two video streams with customizable transition curves.
```go
xfade := effects.NewXFadeFilter(effects.XFadeFade, 1*time.Second, 5*time.Second)
```
