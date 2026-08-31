package types_test

import (
	"testing"
	"time"

	"github.com/farshidrezaei/vidonyx/types"
)

func TestRational_NewAndReduce(t *testing.T) {
	tests := []struct {
		name        string
		num         int64
		den         int64
		expectedNum int64
		expectedDen int64
		shouldPanic bool
	}{
		{name: "simple fraction", num: 1, den: 2, expectedNum: 1, expectedDen: 2},
		{name: "reducible fraction", num: 4, den: 8, expectedNum: 1, expectedDen: 2},
		{name: "negative numerator", num: -3, den: 9, expectedNum: -1, expectedDen: 3},
		{name: "negative denominator", num: 3, den: -9, expectedNum: -1, expectedDen: 3},
		{name: "both negative", num: -4, den: -16, expectedNum: 1, expectedDen: 4},
		{name: "zero numerator", num: 0, den: 10, expectedNum: 0, expectedDen: 1},
		{name: "division by zero panics", num: 5, den: 0, shouldPanic: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic for den=0, but did not panic")
					}
				}()
				_ = types.NewRational(tt.num, tt.den)
				return
			}

			r := types.NewRational(tt.num, tt.den)
			if r.Num != tt.expectedNum || r.Den != tt.expectedDen {
				t.Errorf("NewRational(%d, %d) = %d/%d, want %d/%d", tt.num, tt.den, r.Num, r.Den, tt.expectedNum, tt.expectedDen)
			}
		})
	}
}

func TestRational_FromFloat(t *testing.T) {
	tests := []struct {
		name     string
		val      float64
		expected types.Rational
	}{
		{name: "half", val: 0.5, expected: types.NewRational(1, 2)},
		{name: "quarter", val: 0.25, expected: types.NewRational(1, 4)},
		{name: "one third approx", val: 1.0 / 3.0, expected: types.NewRational(333333, 1000000)},
		{name: "zero", val: 0.0, expected: types.NewRational(0, 1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := types.RationalFromFloat(tt.val)
			if !got.Equal(tt.expected) {
				t.Errorf("RationalFromFloat(%f) = %v, want %v", tt.val, got, tt.expected)
			}
		})
	}
}

func TestRational_ArithmeticTable(t *testing.T) {
	tests := []struct {
		name     string
		op       string // "add", "sub", "mul", "div"
		a        types.Rational
		b        types.Rational
		expected types.Rational
	}{
		{name: "add same denom", op: "add", a: types.NewRational(1, 4), b: types.NewRational(2, 4), expected: types.NewRational(3, 4)},
		{name: "add diff denom", op: "add", a: types.NewRational(1, 2), b: types.NewRational(1, 3), expected: types.NewRational(5, 6)},
		{name: "sub same denom", op: "sub", a: types.NewRational(3, 4), b: types.NewRational(1, 4), expected: types.NewRational(1, 2)},
		{name: "sub diff denom", op: "sub", a: types.NewRational(1, 2), b: types.NewRational(1, 3), expected: types.NewRational(1, 6)},
		{name: "mul simple", op: "mul", a: types.NewRational(2, 3), b: types.NewRational(3, 4), expected: types.NewRational(1, 2)},
		{name: "mul by zero", op: "mul", a: types.NewRational(2, 3), b: types.NewRational(0, 1), expected: types.NewRational(0, 1)},
		{name: "div simple", op: "div", a: types.NewRational(1, 2), b: types.NewRational(1, 4), expected: types.NewRational(2, 1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got types.Rational
			switch tt.op {
			case "add":
				got = tt.a.Add(tt.b)
			case "sub":
				got = tt.a.Sub(tt.b)
			case "mul":
				got = tt.a.Mul(tt.b)
			case "div":
				got = tt.a.Div(tt.b)
			}
			if !got.Equal(tt.expected) {
				t.Errorf("%s(%v, %v) = %v, want %v", tt.op, tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

func TestRational_InvertAndNegate(t *testing.T) {
	r := types.NewRational(3, 4)

	inv := r.Invert()
	if !inv.Equal(types.NewRational(4, 3)) {
		t.Errorf("Invert(%v) = %v, want 4/3", r, inv)
	}

	neg := r.Negate()
	if !neg.Equal(types.NewRational(-3, 4)) {
		t.Errorf("Negate(%v) = %v, want -3/4", r, neg)
	}

	// Invert zero should panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic when inverting 0, but did not panic")
		}
	}()
	_ = types.NewRational(0, 1).Invert()
}

func TestRational_Comparisons(t *testing.T) {
	tests := []struct {
		name       string
		a          types.Rational
		b          types.Rational
		wantEqual  bool
		wantGt     bool
		wantLt     bool
		wantIsZero bool
	}{
		{
			name:       "equal fractions",
			a:          types.NewRational(2, 4),
			b:          types.NewRational(1, 2),
			wantEqual:  true,
			wantGt:     false,
			wantLt:     false,
			wantIsZero: false,
		},
		{
			name:       "a greater than b",
			a:          types.NewRational(3, 4),
			b:          types.NewRational(1, 2),
			wantEqual:  false,
			wantGt:     true,
			wantLt:     false,
			wantIsZero: false,
		},
		{
			name:       "a less than b",
			a:          types.NewRational(1, 3),
			b:          types.NewRational(1, 2),
			wantEqual:  false,
			wantGt:     false,
			wantLt:     true,
			wantIsZero: false,
		},
		{
			name:       "zero value",
			a:          types.NewRational(0, 5),
			b:          types.NewRational(1, 1),
			wantEqual:  false,
			wantGt:     false,
			wantLt:     true,
			wantIsZero: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Equal(tt.b); got != tt.wantEqual {
				t.Errorf("%v.Equal(%v) = %v, want %v", tt.a, tt.b, got, tt.wantEqual)
			}
			if got := tt.a.GreaterThan(tt.b); got != tt.wantGt {
				t.Errorf("%v.GreaterThan(%v) = %v, want %v", tt.a, tt.b, got, tt.wantGt)
			}
			if got := tt.a.LessThan(tt.b); got != tt.wantLt {
				t.Errorf("%v.LessThan(%v) = %v, want %v", tt.a, tt.b, got, tt.wantLt)
			}
			if got := tt.a.IsZero(); got != tt.wantIsZero {
				t.Errorf("%v.IsZero() = %v, want %v", tt.a, got, tt.wantIsZero)
			}
		})
	}
}

func TestRational_FrameConversionsTable(t *testing.T) {
	tests := []struct {
		name       string
		duration   time.Duration
		fps        types.Rational
		wantFrames int64
	}{
		{name: "1s at 30fps", duration: 1 * time.Second, fps: types.FPS30, wantFrames: 30},
		{name: "5s at 60fps", duration: 5 * time.Second, fps: types.FPS60, wantFrames: 300},
		{name: "2.5s at 24fps", duration: 2500 * time.Millisecond, fps: types.FPS24, wantFrames: 60},
		{name: "1s at 29.97fps", duration: 1 * time.Second, fps: types.FPS2997, wantFrames: 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := types.RationalFromDuration(tt.duration)
			frames := r.ToFrames(tt.fps)
			if frames != tt.wantFrames {
				t.Errorf("ToFrames(%v at %v) = %d, want %d", tt.duration, tt.fps, frames, tt.wantFrames)
			}

			durBack := types.FramesToRational(frames, tt.fps)
			if durBack.Duration() != tt.duration {
				t.Logf("FramesToRational roundtrip: %v -> %d frames -> %v", tt.duration, frames, durBack.Duration())
			}
		})
	}
}

func TestRational_StringAndFFmpeg(t *testing.T) {
	tests := []struct {
		name           string
		r              types.Rational
		expectedStr    string
		expectedFFmpeg string
	}{
		{name: "integer fraction", r: types.NewRational(30, 1), expectedStr: "30", expectedFFmpeg: "30"},
		{name: "standard fractional", r: types.FPS2997, expectedStr: "30000/1001", expectedFFmpeg: "30000/1001"},
		{name: "standard 23.976", r: types.FPS23976, expectedStr: "24000/1001", expectedFFmpeg: "24000/1001"},
		{name: "reduced fraction", r: types.NewRational(16, 9), expectedStr: "16/9", expectedFFmpeg: "16/9"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.String(); got != tt.expectedStr {
				t.Errorf("String() = %s, want %s", got, tt.expectedStr)
			}
			if got := tt.r.FFmpegString(); got != tt.expectedFFmpeg {
				t.Errorf("FFmpegString() = %s, want %s", got, tt.expectedFFmpeg)
			}
		})
	}
}
