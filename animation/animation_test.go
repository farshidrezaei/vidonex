package animation_test

import (
	"strings"
	"testing"
	"time"

	"github.com/farshidrezaei/vidonyx/animation"
	"github.com/farshidrezaei/vidonyx/types"
)

func TestCalculateProgressTable(t *testing.T) {
	tests := []struct {
		name     string
		progress float64
		easing   animation.EasingFunction
		wantMin  float64
		wantMax  float64
	}{
		{
			name:     "linear at 0.5",
			progress: 0.5,
			easing:   animation.EasingLinear,
			wantMin:  0.5,
			wantMax:  0.5,
		},
		{
			name:     "ease in quad starts slower",
			progress: 0.5,
			easing:   animation.EasingEaseInQuad,
			wantMin:  0.25,
			wantMax:  0.25,
		},
		{
			name:     "ease out quad starts faster",
			progress: 0.5,
			easing:   animation.EasingEaseOutQuad,
			wantMin:  0.75,
			wantMax:  0.75,
		},
		{
			name:     "boundary below zero clamps",
			progress: -0.5,
			easing:   animation.EasingLinear,
			wantMin:  0.0,
			wantMax:  0.0,
		},
		{
			name:     "boundary above one clamps",
			progress: 1.5,
			easing:   animation.EasingLinear,
			wantMin:  1.0,
			wantMax:  1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := animation.CalculateProgress(tt.progress, tt.easing)
			if got < tt.wantMin || got > tt.wantMax {
				t.Errorf("CalculateProgress(%f, %s) = %f, want [%f, %f]", tt.progress, tt.easing, got, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestFloatKeyframeTrackTable(t *testing.T) {
	track := animation.NewFloatKeyframeTrack()
	track.AddKeyframe(0*time.Second, 0.0, animation.EasingLinear)
	track.AddKeyframe(2*time.Second, 100.0, animation.EasingLinear)
	track.AddKeyframe(4*time.Second, 200.0, animation.EasingEaseInQuad)

	tests := []struct {
		name      string
		timestamp time.Duration
		wantValue float64
	}{
		{
			name:      "at start 0s",
			timestamp: 0 * time.Second,
			wantValue: 0.0,
		},
		{
			name:      "midway 1s linear",
			timestamp: 1 * time.Second,
			wantValue: 50.0,
		},
		{
			name:      "exact keyframe 2s",
			timestamp: 2 * time.Second,
			wantValue: 100.0,
		},
		{
			name:      "beyond last keyframe 5s clamps",
			timestamp: 5 * time.Second,
			wantValue: 200.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := track.Evaluate(tt.timestamp)
			if got != tt.wantValue {
				t.Errorf("Evaluate(%v) = %f, want %f", tt.timestamp, got, tt.wantValue)
			}
		})
	}

	// Verify FFmpeg expression generation
	expr := track.ToFFmpegExpression()
	if !strings.Contains(expr, "if(lt(t,") {
		t.Errorf("expected generated expression to contain nested IFs, got: %s", expr)
	}
}

func TestPositionTrackTable(t *testing.T) {
	track := animation.NewPositionTrack()
	track.AddKeyframe(0*time.Second, types.Point{X: 100, Y: 200}, animation.EasingLinear)
	track.AddKeyframe(2*time.Second, types.Point{X: 500, Y: 800}, animation.EasingEaseOutQuad)

	xExpr, yExpr := track.ToFFmpegExpressions()
	if !strings.Contains(xExpr, "100.0000") || !strings.Contains(xExpr, "500.0000") {
		t.Errorf("xExpr does not contain expected coordinates: %s", xExpr)
	}
	if !strings.Contains(yExpr, "200.0000") || !strings.Contains(yExpr, "800.0000") {
		t.Errorf("yExpr does not contain expected coordinates: %s", yExpr)
	}
}
