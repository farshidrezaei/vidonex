package timeline

import (
	"iter"
	"time"

	"github.com/farshidrezaei/vidonyx/subtitles"
)

// TrackKind specifies the primary media layer type of a track.
type TrackKind uint8

const (
	// TrackKindVideo represents video and image media layers.
	TrackKindVideo TrackKind = iota
	// TrackKindAudio represents sound effects, voiceover, and music layers.
	TrackKindAudio
	// TrackKindOverlay represents graphic overlays, stickers, and Picture-in-Picture layers.
	TrackKindOverlay
	// TrackKindText represents title cards, subtitles, and captions.
	TrackKindText
	// TrackKindSubtitle represents dedicated subtitle and caption streams.
	TrackKindSubtitle
)

// String returns the human-readable name of the track kind.
func (kind TrackKind) String() string {
	switch kind {
	case TrackKindVideo:
		return "video"
	case TrackKindAudio:
		return "audio"
	case TrackKindOverlay:
		return "overlay"
	case TrackKindText:
		return "text"
	case TrackKindSubtitle:
		return "subtitle"
	default:
		return "unknown"
	}
}

// Track defines a sequence or layer of clips within a timeline.
type Track struct {
	ID            string
	Kind          TrackKind
	ZIndex        int
	Muted         bool
	Volume        float64
	Clips         []*Clip
	Transitions   []*Transition
	SubtitleTrack *subtitles.SubtitleTrack
}

// NewTrack creates a new Track with default volume.
func NewTrack(id string, kind TrackKind) *Track {
	return &Track{
		ID:          id,
		Kind:        kind,
		Volume:      1.0,
		Clips:       make([]*Clip, 0),
		Transitions: make([]*Transition, 0),
	}
}

// NewSubtitleTrack creates a track dedicated to rendering subtitles.
func NewSubtitleTrack(id string, subtitleTrack *subtitles.SubtitleTrack) *Track {
	return &Track{
		ID:            id,
		Kind:          TrackKindSubtitle,
		Volume:        1.0,
		Clips:         make([]*Clip, 0),
		Transitions:   make([]*Transition, 0),
		SubtitleTrack: subtitleTrack,
	}
}

// SetSubtitleTrack attaches a subtitle track.
func (track *Track) SetSubtitleTrack(subtitleTrack *subtitles.SubtitleTrack) *Track {
	track.SubtitleTrack = subtitleTrack
	return track
}

// SetZIndex sets the visual layer priority (higher Z-index renders on top).
func (track *Track) SetZIndex(zIndex int) *Track {
	track.ZIndex = zIndex
	return track
}

// SetVolume sets the track audio multiplier (1.0 = 100%).
func (track *Track) SetVolume(volume float64) *Track {
	track.Volume = volume
	return track
}

// SetMuted mutes or unmutes the track.
func (track *Track) SetMuted(muted bool) *Track {
	track.Muted = muted
	return track
}

// AddClip appends one or more clips to the track.
func (track *Track) AddClip(clips ...*Clip) *Track {
	track.Clips = append(track.Clips, clips...)
	return track
}

// AddTransition adds one or more transitions between clips on this track.
func (track *Track) AddTransition(transitions ...*Transition) *Track {
	track.Transitions = append(track.Transitions, transitions...)
	return track
}

// Duration returns the timestamp where the last clip ends on this track.
func (track *Track) Duration() time.Duration {
	var maximumEndTimestamp time.Duration
	for _, clip := range track.Clips {
		if endTimestamp := clip.TimelineEnd(); endTimestamp > maximumEndTimestamp {
			maximumEndTimestamp = endTimestamp
		}
	}
	if track.SubtitleTrack != nil && len(track.SubtitleTrack.Cues) > 0 {
		lastCue := track.SubtitleTrack.Cues[len(track.SubtitleTrack.Cues)-1]
		if lastCue.EndTime > maximumEndTimestamp {
			maximumEndTimestamp = lastCue.EndTime
		}
	}
	return maximumEndTimestamp
}

// AllClips returns an iterator over all clips in the track.
func (track *Track) AllClips() iter.Seq[*Clip] {
	return func(yield func(*Clip) bool) {
		for _, clip := range track.Clips {
			if !yield(clip) {
				return
			}
		}
	}
}
