package types_test

import (
	"testing"
	"time"

	"github.com/farshidrezaei/vidonyx/types"
)

func TestRational_Arithmetic(t *testing.T) {
	r1 := types.NewRational(1, 2)
	r2 := types.NewRational(1, 3)

	// Add
	add := r1.Add(r2)
	if expected := types.NewRational(5, 6); !add.Equal(expected) {
		t.Fatalf("Add: expected %v, got %v", expected, add)
	}

	// Sub
	sub := r1.Sub(r2)
	if expected := types.NewRational(1, 6); !sub.Equal(expected) {
		t.Fatalf("Sub: expected %v, got %v", expected, sub)
	}

	// Mul
	mul := r1.Mul(r2)
	if expected := types.NewRational(1, 6); !mul.Equal(expected) {
		t.Fatalf("Mul: expected %v, got %v", expected, mul)
	}

	// Div
	div := r1.Div(r2)
	if expected := types.NewRational(3, 2); !div.Equal(expected) {
		t.Fatalf("Div: expected %v, got %v", expected, div)
	}
}

func TestRational_FrameConversions(t *testing.T) {
	fps := types.FPS30
	duration := types.NewRational(5, 1) // 5 seconds

	frames := duration.ToFrames(fps)
	if frames != 150 {
		t.Fatalf("ToFrames: expected 150, got %d", frames)
	}

	durBack := types.FramesToRational(150, fps)
	if !durBack.Equal(duration) {
		t.Fatalf("FramesToRational: expected %v, got %v", duration, durBack)
	}
}

func TestRational_DurationAndSeconds(t *testing.T) {
	d := 2500 * time.Millisecond
	r := types.RationalFromDuration(d)

	if r.Seconds() != 2.5 {
		t.Fatalf("Seconds: expected 2.5, got %f", r.Seconds())
	}

	if r.Duration() != d {
		t.Fatalf("Duration: expected %v, got %v", d, r.Duration())
	}
}

func TestRational_StandardFPS(t *testing.T) {
	tests := []struct {
		fps      types.Rational
		expected string
	}{
		{types.FPS24, "24"},
		{types.FPS30, "30"},
		{types.FPS60, "60"},
		{types.FPS29_97, "30000/1001"},
		{types.FPS23_976, "24000/1001"},
	}

	for _, tt := range tests {
		if tt.fps.FFmpegString() != tt.expected {
			t.Errorf("FFmpegString for %v: expected %s, got %s", tt.fps, tt.expected, tt.fps.FFmpegString())
		}
	}
}
