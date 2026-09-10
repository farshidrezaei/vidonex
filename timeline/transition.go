// Package timeline provides declarative composition models for video timelines, tracks, clips, and transitions.
package timeline

import (
	"time"
)

// TransitionType enumerates visual and audio transition types.
type TransitionType string

// Visual and audio transition constants.
const (
	// TransitionFade fades from clip A to clip B.
	TransitionFade TransitionType = "fade"
	// TransitionDissolve dissolves smoothly between clip A and clip B.
	TransitionDissolve TransitionType = "dissolve"
	// TransitionWipeLeft wipes from right to left.
	TransitionWipeLeft TransitionType = "wipeleft"
	// TransitionWipeRight wipes from left to right.
	TransitionWipeRight TransitionType = "wiperight"
	// TransitionWipeUp wipes from bottom to top.
	TransitionWipeUp TransitionType = "wipeup"
	// TransitionWipeDown wipes from top to bottom.
	TransitionWipeDown TransitionType = "wipedown"
	// TransitionSlideLeft slides incoming clip leftwards.
	TransitionSlideLeft TransitionType = "slideleft"
	// TransitionSlideRight slides incoming clip rightwards.
	TransitionSlideRight TransitionType = "slideright"
	// TransitionSlideUp slides incoming clip upwards.
	TransitionSlideUp TransitionType = "slideup"
	// TransitionSlideDown slides incoming clip downwards.
	TransitionSlideDown TransitionType = "slidedown"
	// TransitionCircleCrop expands a circular crop into the next clip.
	TransitionCircleCrop TransitionType = "circlecrop"
	// TransitionCircleOpen opens a circle from the center into the next clip.
	TransitionCircleOpen TransitionType = "circleopen"
	// TransitionZoomIn zooms smoothly into the next clip.
	TransitionZoomIn TransitionType = "zoomin"

	// TransitionAcrossFade cross-fades audio streams.
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
func NewTransition(id string, transitionType TransitionType, duration time.Duration, clipA, clipB *Clip) *Transition {
	return &Transition{
		ID:       id,
		Type:     transitionType,
		Duration: duration,
		ClipA:    clipA,
		ClipB:    clipB,
	}
}
