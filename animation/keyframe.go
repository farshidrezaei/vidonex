package animation

import (
	"fmt"
	"sort"
	"time"

	"github.com/farshidrezaei/vidonex/types"
)

// Keyframe represents a discrete animated value at a specific timeline timestamp.
type Keyframe[T any] struct {
	Time   time.Duration
	Value  T
	Easing EasingFunction
}

// FloatKeyframeTrack manages a sorted timeline of float64 keyframes.
type FloatKeyframeTrack struct {
	Keyframes []Keyframe[float64]
}

// NewFloatKeyframeTrack creates an empty FloatKeyframeTrack.
func NewFloatKeyframeTrack() *FloatKeyframeTrack {
	return &FloatKeyframeTrack{
		Keyframes: make([]Keyframe[float64], 0),
	}
}

// AddKeyframe inserts a new keyframe in chronological order.
func (track *FloatKeyframeTrack) AddKeyframe(timestamp time.Duration, value float64, easing EasingFunction) *FloatKeyframeTrack {
	if easing == "" {
		easing = EasingLinear
	}
	track.Keyframes = append(track.Keyframes, Keyframe[float64]{
		Time:   timestamp,
		Value:  value,
		Easing: easing,
	})
	sort.SliceStable(track.Keyframes, func(i, j int) bool {
		return track.Keyframes[i].Time < track.Keyframes[j].Time
	})
	return track
}

// Evaluate computes the interpolated float value at the given timestamp.
func (track *FloatKeyframeTrack) Evaluate(timestamp time.Duration) float64 {
	if len(track.Keyframes) == 0 {
		return 0.0
	}
	if len(track.Keyframes) == 1 || timestamp <= track.Keyframes[0].Time {
		return track.Keyframes[0].Value
	}
	lastIndex := len(track.Keyframes) - 1
	if timestamp >= track.Keyframes[lastIndex].Time {
		return track.Keyframes[lastIndex].Value
	}

	for index := 0; index < lastIndex; index++ {
		current := track.Keyframes[index]
		next := track.Keyframes[index+1]

		if timestamp >= current.Time && timestamp <= next.Time {
			intervalDuration := next.Time - current.Time
			if intervalDuration <= 0 {
				return next.Value
			}
			progress := float64(timestamp-current.Time) / float64(intervalDuration)
			return InterpolateFloat(current.Value, next.Value, progress, next.Easing)
		}
	}

	return track.Keyframes[lastIndex].Value
}

// ToFFmpegExpression generates nested IF expressions for FFmpeg filters using variable "t".
func (track *FloatKeyframeTrack) ToFFmpegExpression() string {
	return track.ToFFmpegExpressionWithVariable("t")
}

// ToFFmpegExpressionWithVariable generates nested IF expressions for FFmpeg filters with a custom time variable (e.g. "t" or "in_time").
func (track *FloatKeyframeTrack) ToFFmpegExpressionWithVariable(timeVariable string) string {
	if len(track.Keyframes) == 0 {
		return "0.0"
	}
	if len(track.Keyframes) == 1 {
		return fmt.Sprintf("%.4f", track.Keyframes[0].Value)
	}

	expressions := make([]string, 0, len(track.Keyframes)-1)
	for index := 0; index < len(track.Keyframes)-1; index++ {
		k1 := track.Keyframes[index]
		k2 := track.Keyframes[index+1]
		expr := GenerateFFmpegExpressionWithVariable(k1.Value, k2.Value, k1.Time, k2.Time, k2.Easing, timeVariable)
		expressions = append(expressions, expr)
	}

	return expressions[0]
}

// PositionKeyframe represents an (X, Y) coordinate keyframe.
type PositionKeyframe struct {
	Time     time.Duration
	Position types.Point
	Easing   EasingFunction
}

// PositionTrack manages keyframed spatial animation across 2D pixel space.
type PositionTrack struct {
	Keyframes []PositionKeyframe
}

// NewPositionTrack creates an empty PositionTrack.
func NewPositionTrack() *PositionTrack {
	return &PositionTrack{
		Keyframes: make([]PositionKeyframe, 0),
	}
}

// AddKeyframe adds a position keyframe in chronological order.
func (track *PositionTrack) AddKeyframe(timestamp time.Duration, position types.Point, easing EasingFunction) *PositionTrack {
	if easing == "" {
		easing = EasingLinear
	}
	track.Keyframes = append(track.Keyframes, PositionKeyframe{
		Time:     timestamp,
		Position: position,
		Easing:   easing,
	})
	sort.SliceStable(track.Keyframes, func(i, j int) bool {
		return track.Keyframes[i].Time < track.Keyframes[j].Time
	})
	return track
}

// ToFFmpegExpressions generates (XExpression, YExpression) strings for FFmpeg overlay filters.
func (track *PositionTrack) ToFFmpegExpressions() (string, string) {
	if len(track.Keyframes) == 0 {
		return "0", "0"
	}
	if len(track.Keyframes) == 1 {
		return fmt.Sprintf("%d", track.Keyframes[0].Position.X), fmt.Sprintf("%d", track.Keyframes[0].Position.Y)
	}

	xTrack := NewFloatKeyframeTrack()
	yTrack := NewFloatKeyframeTrack()

	for _, kf := range track.Keyframes {
		xTrack.AddKeyframe(kf.Time, float64(kf.Position.X), kf.Easing)
		yTrack.AddKeyframe(kf.Time, float64(kf.Position.Y), kf.Easing)
	}

	return xTrack.ToFFmpegExpression(), yTrack.ToFFmpegExpression()
}

// KenBurnsAnimation builds a slow pan and zoom effect across a visual element.
type KenBurnsAnimation struct {
	StartScale float64
	EndScale   float64
	StartPos   types.Point
	EndPos     types.Point
	Duration   time.Duration
	Easing     EasingFunction
}

// NewKenBurns creates a new KenBurns pan-and-zoom animation definition.
func NewKenBurns(startScale, endScale float64, startPos, endPos types.Point, duration time.Duration, easing EasingFunction) *KenBurnsAnimation {
	if easing == "" {
		easing = EasingLinear
	}
	return &KenBurnsAnimation{
		StartScale: startScale,
		EndScale:   endScale,
		StartPos:   startPos,
		EndPos:     endPos,
		Duration:   duration,
		Easing:     easing,
	}
}

// GenerateScaleExpression creates dynamic FFmpeg scale expression for Ken Burns.
func (kenBurns *KenBurnsAnimation) GenerateScaleExpression() string {
	return GenerateFFmpegExpressionWithVariable(kenBurns.StartScale, kenBurns.EndScale, 0, kenBurns.Duration, kenBurns.Easing, "in_time")
}
