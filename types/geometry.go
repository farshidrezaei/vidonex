package types

import (
	"fmt"
	"math"
)

// Size defines an integer width and height in pixels.
type Size struct {
	Width  int
	Height int
}

// Common video standard resolutions.
var (
	Res4K            = Size{Width: 3840, Height: 2160}
	Res1080p         = Size{Width: 1920, Height: 1080}
	Res720p          = Size{Width: 1280, Height: 720}
	ResSquare1080    = Size{Width: 1080, Height: 1080}
	ResPortrait1080p = Size{Width: 1080, Height: 1920} // 9:16 vertical (Shorts / Reels / TikTok)
)

// NewSize creates a new Size struct.
func NewSize(width, height int) Size {
	return Size{Width: width, Height: height}
}

// IsZero reports whether size has zero area or dimensions.
func (s Size) IsZero() bool {
	return s.Width <= 0 || s.Height <= 0
}

// AspectRatio returns the aspect ratio as a simplified Rational fraction (e.g. 16/9).
func (s Size) AspectRatio() Rational {
	if s.IsZero() {
		return Rational{Num: 0, Den: 1}
	}
	return NewRational(int64(s.Width), int64(s.Height))
}

// Scale returns a new Size scaled by a factor, rounding to even numbers (crucial for H.264/yuv420p).
func (s Size) Scale(factor float64) Size {
	w := makeEven(int(math.Round(float64(s.Width) * factor)))
	h := makeEven(int(math.Round(float64(s.Height) * factor)))
	return Size{Width: w, Height: h}
}

// FitWithin calculates dimensions that fit entirely inside target while preserving aspect ratio.
func (s Size) FitWithin(target Size) Size {
	if s.IsZero() || target.IsZero() {
		return target
	}
	scaleW := float64(target.Width) / float64(s.Width)
	scaleH := float64(target.Height) / float64(s.Height)
	scale := math.Min(scaleW, scaleH)
	return s.Scale(scale)
}

// FillWithin calculates dimensions that cover target completely while preserving aspect ratio.
func (s Size) FillWithin(target Size) Size {
	if s.IsZero() || target.IsZero() {
		return target
	}
	scaleW := float64(target.Width) / float64(s.Width)
	scaleH := float64(target.Height) / float64(s.Height)
	scale := math.Max(scaleW, scaleH)
	return s.Scale(scale)
}

// String returns "WidthxHeight" representation.
func (s Size) String() string {
	return fmt.Sprintf("%dx%d", s.Width, s.Height)
}

// Point represents a 2D coordinate in integer pixel space.
type Point struct {
	X int
	Y int
}

// PointF represents a 2D coordinate in floating-point normalized space (0.0 to 1.0) or exact pixels.
type PointF struct {
	X float64
	Y float64
}

// Rect defines a rectangle with origin (X, Y) and dimensions (Width, Height).
type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

// Size returns the size component of the rectangle.
func (r Rect) Size() Size {
	return Size{Width: r.Width, Height: r.Height}
}

// Origin returns the top-left origin point of the rectangle.
func (r Rect) Origin() Point {
	return Point{X: r.X, Y: r.Y}
}

// Alignment represents spatial positioning anchors on the canvas.
type Alignment uint8

const (
	AlignTopLeft Alignment = iota
	AlignTopCenter
	AlignTopRight
	AlignCenterLeft
	AlignCenter
	AlignCenterRight
	AlignBottomLeft
	AlignBottomCenter
	AlignBottomRight
)

// CalculateOffset computes the (X, Y) offset to place an item of itemSize inside containerSize with the given alignment.
func CalculateOffset(containerSize, itemSize Size, align Alignment) Point {
	var x, y int

	switch align {
	case AlignTopLeft, AlignCenterLeft, AlignBottomLeft:
		x = 0
	case AlignTopCenter, AlignCenter, AlignBottomCenter:
		x = (containerSize.Width - itemSize.Width) / 2
	case AlignTopRight, AlignCenterRight, AlignBottomRight:
		x = containerSize.Width - itemSize.Width
	}

	switch align {
	case AlignTopLeft, AlignTopCenter, AlignTopRight:
		y = 0
	case AlignCenterLeft, AlignCenter, AlignCenterRight:
		y = (containerSize.Height - itemSize.Height) / 2
	case AlignBottomLeft, AlignBottomCenter, AlignBottomRight:
		y = containerSize.Height - itemSize.Height
	}

	return Point{X: x, Y: y}
}

func makeEven(n int) int {
	if n%2 != 0 {
		return n + 1
	}
	return n
}
