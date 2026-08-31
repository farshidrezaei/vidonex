package types

import (
	"fmt"
	"math"
	"time"
)

// Rational represents an exact fractional number as Num/Den.
// It avoids floating-point precision drift across long video timelines.
type Rational struct {
	Num int64
	Den int64
}

// Common video standard frame rates.
var (
	FPS23_976 = Rational{Num: 24000, Den: 1001}
	FPS24     = Rational{Num: 24, Den: 1}
	FPS25     = Rational{Num: 25, Den: 1}
	FPS29_97  = Rational{Num: 30000, Den: 1001}
	FPS30     = Rational{Num: 30, Den: 1}
	FPS50     = Rational{Num: 50, Den: 1}
	FPS59_94  = Rational{Num: 60000, Den: 1001}
	FPS60     = Rational{Num: 60, Den: 1}
)

// NewRational creates a new Rational fraction and reduces it to simplest terms.
func NewRational(num, den int64) Rational {
	if den == 0 {
		panic("types.NewRational: division by zero")
	}
	r := Rational{Num: num, Den: den}
	return r.Reduce()
}

// RationalFromFloat converts a float64 to a rational approximation with precision.
func RationalFromFloat(f float64) Rational {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		panic("types.RationalFromFloat: NaN or Inf is not supported")
	}

	const precision = 1000000
	den := int64(precision)
	num := int64(math.Round(f * precision))
	return NewRational(num, den)
}

// RationalFromDuration converts a time.Duration to a Rational in seconds.
func RationalFromDuration(d time.Duration) Rational {
	return NewRational(d.Nanoseconds(), int64(time.Second))
}

// Reduce simplifies the fraction using the greatest common divisor (GCD).
func (r Rational) Reduce() Rational {
	if r.Den == 0 {
		return r
	}
	if r.Num == 0 {
		return Rational{Num: 0, Den: 1}
	}
	g := gcd(abs(r.Num), abs(r.Den))
	num := r.Num / g
	den := r.Den / g
	if den < 0 {
		num = -num
		den = -den
	}
	return Rational{Num: num, Den: den}
}

// IsZero reports whether the rational value is zero.
func (r Rational) IsZero() bool {
	return r.Num == 0
}

// Seconds returns the rational value as a float64 number of seconds.
func (r Rational) Seconds() float64 {
	if r.Den == 0 {
		return 0
	}
	return float64(r.Num) / float64(r.Den)
}

// Duration returns the rational value as a time.Duration.
func (r Rational) Duration() time.Duration {
	if r.Den == 0 {
		return 0
	}
	// Calculate nanoseconds with high precision
	ns := (r.Num * int64(time.Second)) / r.Den
	return time.Duration(ns)
}

// Add computes r + other.
func (r Rational) Add(other Rational) Rational {
	if r.Den == other.Den {
		return NewRational(r.Num+other.Num, r.Den)
	}
	num := r.Num*other.Den + other.Num*r.Den
	den := r.Den * other.Den
	return NewRational(num, den)
}

// Sub computes r - other.
func (r Rational) Sub(other Rational) Rational {
	if r.Den == other.Den {
		return NewRational(r.Num-other.Num, r.Den)
	}
	num := r.Num*other.Den - other.Num*r.Den
	den := r.Den * other.Den
	return NewRational(num, den)
}

// Mul computes r * other.
func (r Rational) Mul(other Rational) Rational {
	return NewRational(r.Num*other.Num, r.Den*other.Den)
}

// Div computes r / other.
func (r Rational) Div(other Rational) Rational {
	if other.Num == 0 {
		panic("types.Rational.Div: division by zero")
	}
	return NewRational(r.Num*other.Den, r.Den*other.Num)
}

// Invert returns 1 / r.
func (r Rational) Invert() Rational {
	if r.Num == 0 {
		panic("types.Rational.Invert: cannot invert zero")
	}
	return NewRational(r.Den, r.Num)
}

// Negate returns -r.
func (r Rational) Negate() Rational {
	return Rational{Num: -r.Num, Den: r.Den}
}

// Equal reports whether r == other.
func (r Rational) Equal(other Rational) bool {
	r1 := r.Reduce()
	r2 := other.Reduce()
	return r1.Num == r2.Num && r1.Den == r2.Den
}

// GreaterThan reports whether r > other.
func (r Rational) GreaterThan(other Rational) bool {
	return (r.Num * other.Den) > (other.Num * r.Den)
}

// LessThan reports whether r < other.
func (r Rational) LessThan(other Rational) bool {
	return (r.Num * other.Den) < (other.Num * r.Den)
}

// ToFrames calculates the number of frames for this duration at a given FPS.
func (r Rational) ToFrames(fps Rational) int64 {
	framesRat := r.Mul(fps)
	return int64(math.Round(framesRat.Seconds()))
}

// FramesToRational converts a frame count at a given FPS to a duration Rational.
func FramesToRational(frames int64, fps Rational) Rational {
	return NewRational(frames*fps.Den, fps.Num)
}

// String returns a human-readable representation like "30/1" or "30000/1001".
func (r Rational) String() string {
	if r.Den == 1 {
		return fmt.Sprintf("%d", r.Num)
	}
	return fmt.Sprintf("%d/%d", r.Num, r.Den)
}

// FFmpegString returns the formatted string suitable for FFmpeg expressions.
func (r Rational) FFmpegString() string {
	if r.Den == 1 {
		return fmt.Sprintf("%d", r.Num)
	}
	return fmt.Sprintf("%d/%d", r.Num, r.Den)
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
