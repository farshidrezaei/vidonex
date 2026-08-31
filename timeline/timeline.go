package timeline

import (
	"iter"
	"time"

	"github.com/farshidrezaei/vidonyx/types"
)

// Timeline is the top-level declarative description of a composition.
type Timeline struct {
	Canvas           types.Size
	FPS              types.Rational
	BackgroundColor  types.Color
	ExplicitDuration time.Duration
	Tracks           []*Track
	Metadata         map[string]any
}

// Option configures a Timeline during construction.
type Option func(*Timeline)

// WithCanvas sets the target composition dimensions.
func WithCanvas(size types.Size) Option {
	return func(timeline *Timeline) {
		timeline.Canvas = size
	}
}

// WithFPS sets the target composition frame rate.
func WithFPS(frameRate types.Rational) Option {
	return func(timeline *Timeline) {
		timeline.FPS = frameRate
	}
}

// WithBackgroundColor sets the canvas background color.
func WithBackgroundColor(color types.Color) Option {
	return func(timeline *Timeline) {
		timeline.BackgroundColor = color
	}
}

// WithDuration explicitly sets the timeline duration.
func WithDuration(duration time.Duration) Option {
	return func(timeline *Timeline) {
		timeline.ExplicitDuration = duration
	}
}

// New creates a new Timeline configured with fluent options and sensible defaults.
func New(options ...Option) *Timeline {
	timeline := &Timeline{
		Canvas:          types.Res1080p,
		FPS:             types.FPS30,
		BackgroundColor: types.ColorBlack,
		Tracks:          make([]*Track, 0),
		Metadata:        make(map[string]any),
	}
	for _, option := range options {
		option(timeline)
	}
	return timeline
}

// AddTrack appends one or more tracks to the timeline.
func (timeline *Timeline) AddTrack(tracks ...*Track) *Timeline {
	timeline.Tracks = append(timeline.Tracks, tracks...)
	return timeline
}

// Duration computes the effective duration of the timeline based on its longest track,
// or returns ExplicitDuration if set.
func (timeline *Timeline) Duration() time.Duration {
	if timeline.ExplicitDuration > 0 {
		return timeline.ExplicitDuration
	}
	var maximumDuration time.Duration
	for _, track := range timeline.Tracks {
		if trackDuration := track.Duration(); trackDuration > maximumDuration {
			maximumDuration = trackDuration
		}
	}
	return maximumDuration
}

// TotalFrames calculates the total number of frames for the composition.
func (timeline *Timeline) TotalFrames() int64 {
	rationalDuration := types.RationalFromDuration(timeline.Duration())
	return rationalDuration.ToFrames(timeline.FPS)
}

// AllTracks returns an iterator over tracks in the timeline.
func (timeline *Timeline) AllTracks() iter.Seq[*Track] {
	return func(yield func(*Track) bool) {
		for _, track := range timeline.Tracks {
			if !yield(track) {
				return
			}
		}
	}
}
