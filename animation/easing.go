// Package animation provides keyframing, easing mathematical models, and dynamic parameter interpolation for video elements.
package animation

import (
	"fmt"
	"math"
	"time"
)

// EasingFunction defines the acceleration curve used for interpolation.
type EasingFunction string

// Standard easing curve constants.
const (
	// EasingLinear provides constant velocity with no acceleration.
	EasingLinear EasingFunction = "linear"
	// EasingEaseInQuad accelerates from zero velocity using quadratic curve.
	EasingEaseInQuad EasingFunction = "easeInQuad"
	// EasingEaseOutQuad decelerates to zero velocity using quadratic curve.
	EasingEaseOutQuad EasingFunction = "easeOutQuad"
	// EasingEaseInOutQuad accelerates halfway then decelerates using quadratic curve.
	EasingEaseInOutQuad EasingFunction = "easeInOutQuad"
	// EasingEaseInCubic accelerates using cubic curve.
	EasingEaseInCubic EasingFunction = "easeInCubic"
	// EasingEaseOutCubic decelerates using cubic curve.
	EasingEaseOutCubic EasingFunction = "easeOutCubic"
	// EasingEaseInOutCubic accelerates then decelerates using cubic curve.
	EasingEaseInOutCubic EasingFunction = "easeInOutCubic"
	// EasingEaseInSine accelerates using sinusoidal curve.
	EasingEaseInSine EasingFunction = "easeInSine"
	// EasingEaseOutSine decelerates using sinusoidal curve.
	EasingEaseOutSine EasingFunction = "easeOutSine"
	// EasingEaseInOutSine accelerates then decelerates using sinusoidal curve.
	EasingEaseInOutSine EasingFunction = "easeInOutSine"
)

// CalculateProgress transforms linear progress (0.0 to 1.0) into eased progress.
func CalculateProgress(progress float64, easing EasingFunction) float64 {
	progress = math.Max(0.0, math.Min(1.0, progress))

	switch easing {
	case EasingLinear:
		return progress
	case EasingEaseInQuad:
		return progress * progress
	case EasingEaseOutQuad:
		return progress * (2.0 - progress)
	case EasingEaseInOutQuad:
		if progress < 0.5 {
			return 2.0 * progress * progress
		}
		return -1.0 + (4.0-2.0*progress)*progress
	case EasingEaseInCubic:
		return progress * progress * progress
	case EasingEaseOutCubic:
		p := progress - 1.0
		return p*p*p + 1.0
	case EasingEaseInOutCubic:
		if progress < 0.5 {
			return 4.0 * progress * progress * progress
		}
		p := 2.0*progress - 2.0
		return 0.5*p*p*p + 1.0
	case EasingEaseInSine:
		return 1.0 - math.Cos((progress*math.Pi)/2.0)
	case EasingEaseOutSine:
		return math.Sin((progress * math.Pi) / 2.0)
	case EasingEaseInOutSine:
		return -(math.Cos(math.Pi*progress) - 1.0) / 2.0
	default:
		return progress
	}
}

// InterpolateFloat computes the interpolated value between start and end.
func InterpolateFloat(start, end, progress float64, easing EasingFunction) float64 {
	easedProgress := CalculateProgress(progress, easing)
	return start + (end-start)*easedProgress
}

// GenerateFFmpegExpression creates an FFmpeg expression computing dynamic value across time interval using default variable "t".
func GenerateFFmpegExpression(startValue, endValue float64, startTime, endTime time.Duration, easing EasingFunction) string {
	return GenerateFFmpegExpressionWithVariable(startValue, endValue, startTime, endTime, easing, "t")
}

// GenerateFFmpegExpressionWithVariable creates an FFmpeg expression computing dynamic value across time interval with a configurable time variable name (e.g. "t" or "in_time").
func GenerateFFmpegExpressionWithVariable(startValue, endValue float64, startTime, endTime time.Duration, easing EasingFunction, timeVariable string) string {
	if timeVariable == "" {
		timeVariable = "t"
	}
	startSec := startTime.Seconds()
	endSec := endTime.Seconds()
	durationSec := endSec - startSec

	if durationSec <= 0 {
		return fmt.Sprintf("%.4f", endValue)
	}

	normalizedT := fmt.Sprintf("((%s-%.4f)/%.4f)", timeVariable, startSec, durationSec)

	var progressExpr string
	switch easing {
	case EasingLinear:
		progressExpr = normalizedT
	case EasingEaseInQuad:
		progressExpr = fmt.Sprintf("pow(%s,2)", normalizedT)
	case EasingEaseOutQuad:
		progressExpr = fmt.Sprintf("(%s*(2-%s))", normalizedT, normalizedT)
	case EasingEaseInCubic:
		progressExpr = fmt.Sprintf("pow(%s,3)", normalizedT)
	default:
		progressExpr = normalizedT
	}

	delta := endValue - startValue
	interpolatedExpr := fmt.Sprintf("(%.4f+(%.4f*%s))", startValue, delta, progressExpr)

	return fmt.Sprintf("if(lt(%s,%.4f),%.4f,if(gt(%s,%.4f),%.4f,%s))",
		timeVariable, startSec, startValue, timeVariable, endSec, endValue, interpolatedExpr)
}
