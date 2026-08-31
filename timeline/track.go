package timeline

import (
	"iter"
	"time"
)

// TrackKind specifies the primary media layer type of a track.
type TrackKind uint8

const (
	// TrackKindVideo represents video and image media layers.
	TrackKindVideo TrackKind = iota
	// TrackKindAudio represents sound effects, voiceover, and music layers.
	TrackKindAudio
	// TrackKindOverlay represents graphic overlays, stickers, and PiP layers.
	TrackKindOverlay
	// TrackKindText represents title cards, subtitles, and captions.
	TrackKindText
)

// String returns the name of the track kind.
func (k TrackKind) String() string {
	switch k {
	case TrackKindVideo:
		return "video"
	case TrackKindAudio:
		return "audio"
	case TrackKindOverlay:
		return "overlay"
	case TrackKindText:
		return "text"
	default:
		return "unknown"
	}
}

// Track defines a sequence or layer of clips within a timeline.
type Track struct {
	ID          string
	Kind        TrackKind
	ZIndex      int
	Muted       bool
	Volume      float64
	Clips       []*Clip
	Transitions []*Transition
}

// NewTrack creates a new Track.
func NewTrack(id string, kind TrackKind) *Track {
	return &Track{
		ID:          id,
		Kind:        kind,
		Volume:      1.0,
		Clips:       make([]*Clip, 0),
		Transitions: make([]*Transition, 0),
	}
}

// SetZIndex sets the visual layer priority (higher Z-index renders on top).
func (t *Track) SetZIndex(z int) *Track {
	t.ZIndex = z
	return t
}

// SetVolume sets the track audio multiplier (1.0 = 100%).
func (t *Track) SetVolume(vol float64) *Track {
	t.Volume = vol
	return t
}

// SetMuted mutes or unmutes the track.
func (t *Track) SetMuted(muted bool) *Track {
	t.Muted = muted
	return t
}

// AddClip appends one or more clips to the track.
func (t *Track) AddClip(clips ...*Clip) *Track {
	t.Clips = append(t.Clips, clips...)
	return t
}

// AddTransition adds a transition between clips on this track.
func (t *Track) AddTransition(tr *Transition) *Track {
	t.Transitions = append(t.Transitions, tr)
	return t
}

// Duration returns the timestamp where the last clip ends on this track.
func (t *Track) Duration() time.Duration {
	var maxEnd time.Duration
	for _, c := range t.Clips {
		if end := c.TimelineEnd(); end > maxEnd {
			maxEnd = end
		}
	}
	return maxEnd
}

// AllClips returns an iterator over all clips in the track.
func (t *Track) AllClips() iter.Seq[*Clip] {
	return func(yield func(*Clip) bool) {
		for _, c := range t.Clips {
			if !yield(c) {
				return
			}
		}
	}
}
