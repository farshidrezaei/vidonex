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
	SourceStart   time.Duration // Trim in
	Speed         float64       // Playback rate multiplier (default 1.0)
	Volume        float64       // Audio volume multiplier (default 1.0)
	Opacity       float64       // Visual opacity 0.0 (transparent) to 1.0 (opaque)
	Position      types.Point   // Absolute (X, Y) coordinate on canvas
	Alignment     types.Alignment
	Scale         float64  // Scale factor (1.0 = original)
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
func (c *Clip) TimelineEnd() time.Duration {
	return c.TimelineStart + c.Duration
}

// WithTrim sets the in-point offset within the source media.
func (c *Clip) WithTrim(sourceStart time.Duration) *Clip {
	c.SourceStart = sourceStart
	return c
}

// WithSpeed sets playback speed rate (e.g. 2.0 = double speed, 0.5 = slow motion).
func (c *Clip) WithSpeed(speed float64) *Clip {
	c.Speed = speed
	return c
}

// WithVolume sets the audio volume level (1.0 = 100%).
func (c *Clip) WithVolume(volume float64) *Clip {
	c.Volume = volume
	return c
}

// WithOpacity sets the visual transparency level (0.0 to 1.0).
func (c *Clip) WithOpacity(opacity float64) *Clip {
	c.Opacity = opacity
	return c
}

// WithPosition sets the exact (X, Y) canvas coordinates.
func (c *Clip) WithPosition(p types.Point) *Clip {
	c.Position = p
	return c
}

// WithScale sets the size scale multiplier.
func (c *Clip) WithScale(scale float64) *Clip {
	c.Scale = scale
	return c
}

// WithRotation sets the rotation in degrees.
func (c *Clip) WithRotation(degrees float64) *Clip {
	c.Rotation = degrees
	return c
}

// AddEffect attaches a filter effect to the clip.
func (c *Clip) AddEffect(e Effect) *Clip {
	c.Effects = append(c.Effects, e)
	return c
}
