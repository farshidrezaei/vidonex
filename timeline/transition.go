package timeline

import (
	"time"
)

// TransitionType enumerates visual and audio transition types.
type TransitionType string

const (
	// Visual Transitions
	TransitionFade       TransitionType = "fade"
	TransitionDissolve   TransitionType = "dissolve"
	TransitionWipeLeft   TransitionType = "wipeleft"
	TransitionWipeRight  TransitionType = "wiperight"
	TransitionWipeUp     TransitionType = "wipeup"
	TransitionWipeDown   TransitionType = "wipedown"
	TransitionSlideLeft  TransitionType = "slideleft"
	TransitionSlideRight TransitionType = "slideright"
	TransitionSlideUp    TransitionType = "slideup"
	TransitionSlideDown  TransitionType = "slidedown"
	TransitionCircleCrop TransitionType = "circlecrop"
	TransitionZoomIn     TransitionType = "zoomin"

	// Audio Transitions
	TransitionAcrossFade TransitionType = "acrossfade"
)

// Transition specifies an overlap transition between two adjacent clips on a track.
type Transition struct {
	ID       string
	Type     TransitionType
	Duration time.Duration
	ClipA    *Clip // Preceding outgoing clip
	ClipB    *Clip // Succeeding incoming clip
}

// NewTransition creates a new Transition between clipA and clipB.
func NewTransition(id string, transType TransitionType, duration time.Duration, clipA, clipB *Clip) *Transition {
	return &Transition{
		ID:       id,
		Type:     transType,
		Duration: duration,
		ClipA:    clipA,
		ClipB:    clipB,
	}
}
