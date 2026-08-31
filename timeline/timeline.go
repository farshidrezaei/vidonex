package timeline

import (
	"iter"
	"time"

	"github.com/farshidrezaei/vidonyx/types"
)

// Timeline is the top-level declarative description of a composition.
type Timeline struct {
	Canvas          types.Size
	FPS             types.Rational
	BackgroundColor types.Color
	ExplicitDuration time.Duration
	Tracks          []*Track
	Metadata        map[string]any
}

// Option configures a Timeline during construction.
type Option func(*Timeline)

// WithCanvas sets the target composition dimensions.
func WithCanvas(size types.Size) Option {
	return func(t *Timeline) {
		t.Canvas = size
	}
}

// WithFPS sets the target composition frame rate.
func WithFPS(fps types.Rational) Option {
	return func(t *Timeline) {
		t.FPS = fps
	}
}

// WithBackgroundColor sets the canvas background color.
func WithBackgroundColor(color types.Color) Option {
	return func(t *Timeline) {
		t.BackgroundColor = color
	}
}

// WithDuration explicitly sets the timeline duration.
func WithDuration(d time.Duration) Option {
	return func(t *Timeline) {
		t.ExplicitDuration = d
	}
}

// New creates a new Timeline configured with fluent options and sensible defaults.
func New(opts ...Option) *Timeline {
	tl := &Timeline{
		Canvas:          types.Res1080p,
		FPS:             types.FPS30,
		BackgroundColor: types.ColorBlack,
		Tracks:          make([]*Track, 0),
		Metadata:        make(map[string]any),
	}
	for _, opt := range opts {
		opt(tl)
	}
	return tl
}

// AddTrack appends one or more tracks to the timeline.
func (t *Timeline) AddTrack(tracks ...*Track) *Timeline {
	t.Tracks = append(t.Tracks, tracks...)
	return t
}

// Duration computes the effective duration of the timeline based on its longest track,
// or returns ExplicitDuration if set.
func (t *Timeline) Duration() time.Duration {
	if t.ExplicitDuration > 0 {
		return t.ExplicitDuration
	}
	var maxDur time.Duration
	for _, tr := range t.Tracks {
		if d := tr.Duration(); d > maxDur {
			maxDur = d
		}
	}
	return maxDur
}

// TotalFrames calculates the total number of frames for the composition.
func (t *Timeline) TotalFrames() int64 {
	ratDur := types.RationalFromDuration(t.Duration())
	return ratDur.ToFrames(t.FPS)
}

// AllTracks returns an iterator over tracks in the timeline.
func (t *Timeline) AllTracks() iter.Seq[*Track] {
	return func(yield func(*Track) bool) {
		for _, tr := range t.Tracks {
			if !yield(tr) {
				return
			}
		}
	}
}
