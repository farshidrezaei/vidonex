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
	// TransitionFadeBlack fades to black then into clip B.
	TransitionFadeBlack TransitionType = "fadeblack"
	// TransitionFadeWhite fades to white then into clip B.
	TransitionFadeWhite TransitionType = "fadewhite"
	// TransitionFadeGrays fades to grayscale then into clip B.
	TransitionFadeGrays TransitionType = "fadegrays"
	// TransitionDissolve dissolves smoothly between clip A and clip B.
	TransitionDissolve TransitionType = "dissolve"
	// TransitionPixelize pixelates into clip B.
	TransitionPixelize TransitionType = "pixelize"

	// Directional Wipes
	TransitionWipeLeft        TransitionType = "wipeleft"
	TransitionWipeRight       TransitionType = "wiperight"
	TransitionWipeUp          TransitionType = "wipeup"
	TransitionWipeDown        TransitionType = "wipedown"
	TransitionWipeTopLeft     TransitionType = "wipetl"
	TransitionWipeTopRight    TransitionType = "wipetr"
	TransitionWipeBottomLeft  TransitionType = "wipebl"
	TransitionWipeBottomRight TransitionType = "wipebr"

	// Directional Slides
	TransitionSlideLeft  TransitionType = "slideleft"
	TransitionSlideRight TransitionType = "slideright"
	TransitionSlideUp    TransitionType = "slideup"
	TransitionSlideDown  TransitionType = "slidedown"

	// Smooth Directional Motion
	TransitionSmoothLeft  TransitionType = "smoothleft"
	TransitionSmoothRight TransitionType = "smoothright"
	TransitionSmoothUp    TransitionType = "smoothup"
	TransitionSmoothDown  TransitionType = "smoothdown"

	// Shapes & Iris
	TransitionCircleCrop   TransitionType = "circlecrop"
	TransitionCircleOpen   TransitionType = "circleopen"
	TransitionCircleClose  TransitionType = "circleclose"
	TransitionRectCrop     TransitionType = "rectcrop"
	TransitionRadial       TransitionType = "radial"
	TransitionDistance     TransitionType = "distance"
	TransitionHorzOpen     TransitionType = "horzopen"
	TransitionHorzClose    TransitionType = "horzclose"
	TransitionVertOpen     TransitionType = "vertopen"
	TransitionVertClose    TransitionType = "vertclose"
	TransitionHorzSlice    TransitionType = "hlslice"
	TransitionVertSlice    TransitionType = "vuslice"
	TransitionZoomIn       TransitionType = "zoomin"
	TransitionSqueezeVert  TransitionType = "squeezev"
	TransitionSqueezeHorz  TransitionType = "squeezeh"

	// Audio Transitions
	TransitionAcrossFade TransitionType = "acrossfade"
)

// AllVideoTransitions returns a list of all supported video transition types.
func AllVideoTransitions() []TransitionType {
	return []TransitionType{
		TransitionFade,
		TransitionFadeBlack,
		TransitionFadeWhite,
		TransitionFadeGrays,
		TransitionDissolve,
		TransitionPixelize,
		TransitionWipeLeft,
		TransitionWipeRight,
		TransitionWipeUp,
		TransitionWipeDown,
		TransitionWipeTopLeft,
		TransitionWipeTopRight,
		TransitionWipeBottomLeft,
		TransitionWipeBottomRight,
		TransitionSlideLeft,
		TransitionSlideRight,
		TransitionSlideUp,
		TransitionSlideDown,
		TransitionSmoothLeft,
		TransitionSmoothRight,
		TransitionSmoothUp,
		TransitionSmoothDown,
		TransitionCircleCrop,
		TransitionCircleOpen,
		TransitionCircleClose,
		TransitionRectCrop,
		TransitionRadial,
		TransitionDistance,
		TransitionHorzOpen,
		TransitionHorzClose,
		TransitionVertOpen,
		TransitionVertClose,
		TransitionHorzSlice,
		TransitionVertSlice,
		TransitionZoomIn,
		TransitionSqueezeVert,
		TransitionSqueezeHorz,
	}
}

// IsValidTransitionType reports whether the specified transition type is recognized.
func IsValidTransitionType(transitionType TransitionType) bool {
	if transitionType == TransitionAcrossFade {
		return true
	}
	for _, supported := range AllVideoTransitions() {
		if supported == transitionType {
			return true
		}
	}
	return false
}

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
