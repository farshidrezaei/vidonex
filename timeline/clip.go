package timeline

import (
	"time"

	"github.com/farshidrezaei/vidonyx/types"
)

// Clip represents a media element placed at a specific time range in a track.
type Clip struct {
	ID            string
	Source        string // Filepath or URL
	TimelineStart time.Duration
	Duration      time.Duration
	SourceStart   time.Duration // Trim in offset within the source media
	Speed         float64       // Playback rate multiplier (default 1.0)
	Volume        float64       // Audio volume multiplier (default 1.0)
	Opacity       float64       // Visual opacity 0.0 (transparent) to 1.0 (opaque)
	Position      types.Point   // Absolute (X, Y) coordinate on canvas
	Alignment     types.Alignment
	Scale         float64  // Scale factor (1.0 = original size)
	Rotation      float64  // Rotation in degrees (e.g. 90.0, 180.0)
	Effects       []Effect // Attached filters and visual transformations
}

// NewClip creates a Clip with sensible defaults.
func NewClip(id, source string, start, duration time.Duration) *Clip {
	return &Clip{
		ID:            id,
		Source:        source,
		TimelineStart: start,
		Duration:      duration,
		SourceStart:   0,
		Speed:         1.0,
		Volume:        1.0,
		Opacity:       1.0,
		Scale:         1.0,
		Alignment:     types.AlignCenter,
		Effects:       make([]Effect, 0),
	}
}

// TimelineEnd returns the absolute end timestamp of the clip on the timeline.
func (clip *Clip) TimelineEnd() time.Duration {
	return clip.TimelineStart + clip.Duration
}

// WithTrim sets the in-point offset within the source media.
func (clip *Clip) WithTrim(sourceStart time.Duration) *Clip {
	clip.SourceStart = sourceStart
	return clip
}

// WithSpeed sets playback speed rate (e.g. 2.0 = double speed, 0.5 = slow motion).
func (clip *Clip) WithSpeed(speed float64) *Clip {
	clip.Speed = speed
	return clip
}

// WithVolume sets the audio volume level (1.0 = 100%).
func (clip *Clip) WithVolume(volume float64) *Clip {
	clip.Volume = volume
	return clip
}

// WithOpacity sets the visual transparency level (0.0 to 1.0).
func (clip *Clip) WithOpacity(opacity float64) *Clip {
	clip.Opacity = opacity
	return clip
}

// WithPosition sets the exact (X, Y) canvas coordinates.
func (clip *Clip) WithPosition(position types.Point) *Clip {
	clip.Position = position
	return clip
}

// WithScale sets the size scale multiplier.
func (clip *Clip) WithScale(scale float64) *Clip {
	clip.Scale = scale
	return clip
}

// WithRotation sets the rotation in degrees.
func (clip *Clip) WithRotation(degrees float64) *Clip {
	clip.Rotation = degrees
	return clip
}

// AddEffect attaches a filter effect to the clip.
func (clip *Clip) AddEffect(effect Effect) *Clip {
	clip.Effects = append(clip.Effects, effect)
	return clip
}
