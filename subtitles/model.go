// Package subtitles provides subtitle parsing (SRT, WebVTT), caption rendering, and animated word highlights.
package subtitles

import (
	"time"

	"github.com/farshidrezaei/vidonyx/types"
)

// WordCue represents a single word timestamped within a subtitle sentence (Karaoke/TikTok captions).
type WordCue struct {
	Word      string
	StartTime time.Duration
	EndTime   time.Duration
}

// SubtitleCue represents a single subtitle entry with start/end timestamps and text content.
type SubtitleCue struct {
	Index     int
	StartTime time.Duration
	EndTime   time.Duration
	Text      string
	Words     []WordCue
}

// Duration returns the display duration of the subtitle cue.
func (cue SubtitleCue) Duration() time.Duration {
	return cue.EndTime - cue.StartTime
}

// SubtitleStyle defines visual formatting, colors, borders, and placement for rendered subtitles.
type SubtitleStyle struct {
	FontName       string
	FontFile       string
	FontSize       int
	PrimaryColor   types.Color
	HighlightColor types.Color
	OutlineColor   types.Color
	OutlineWidth   int
	Box            bool
	BoxColor       types.Color
	BoxBorderWidth int
	Alignment      types.Alignment
	MarginBottom   int
	MarginTop      int
}

// DefaultSubtitleStyle returns modern high-legibility subtitle styles (white text, subtle dark box).
func DefaultSubtitleStyle() SubtitleStyle {
	return SubtitleStyle{
		FontSize:       36,
		PrimaryColor:   types.ColorWhite,
		HighlightColor: types.ColorRed,
		OutlineColor:   types.ColorBlack,
		OutlineWidth:   2,
		Box:            true,
		BoxColor:       types.RGBA(0, 0, 0, 180),
		BoxBorderWidth: 8,
		Alignment:      types.AlignBottomCenter,
		MarginBottom:   60,
		MarginTop:      60,
	}
}

// SubtitleTrack encapsulates a collection of subtitle cues and their shared rendering style.
type SubtitleTrack struct {
	Cues  []SubtitleCue
	Style SubtitleStyle
}

// NewSubtitleTrack creates an empty SubtitleTrack with standard styling.
func NewSubtitleTrack() *SubtitleTrack {
	return &SubtitleTrack{
		Cues:  make([]SubtitleCue, 0),
		Style: DefaultSubtitleStyle(),
	}
}

// SetStyle configures custom subtitle styling.
func (track *SubtitleTrack) SetStyle(style SubtitleStyle) *SubtitleTrack {
	track.Style = style
	return track
}

// AddCue appends a subtitle cue to the track.
func (track *SubtitleTrack) AddCue(cue SubtitleCue) *SubtitleTrack {
	track.Cues = append(track.Cues, cue)
	return track
}
