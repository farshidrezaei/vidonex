// Package types defines exact mathematical, geometric, and temporal primitives for Vidonex.
package types

import (
	"fmt"
	"strconv"
	"strings"
)

// Color represents an RGBA color with 8-bit channels.
type Color struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

// Common color constants.
var (
	// ColorTransparent represents fully transparent black (Alpha = 0).
	ColorTransparent = Color{R: 0, G: 0, B: 0, A: 0}
	// ColorBlack represents solid opaque black.
	ColorBlack = Color{R: 0, G: 0, B: 0, A: 255}
	// ColorWhite represents solid opaque white.
	ColorWhite = Color{R: 255, G: 255, B: 255, A: 255}
	// ColorRed represents solid opaque red.
	ColorRed = Color{R: 255, G: 0, B: 0, A: 255}
	// ColorGreen represents solid opaque green.
	ColorGreen = Color{R: 0, G: 255, B: 0, A: 255}
	// ColorBlue represents solid opaque blue.
	ColorBlue = Color{R: 0, G: 0, B: 255, A: 255}
	// ColorYellow represents solid opaque yellow.
	ColorYellow = Color{R: 255, G: 255, B: 0, A: 255}
)

// RGB constructs an opaque Color (Alpha=255).
func RGB(r, g, b uint8) Color {
	return Color{R: r, G: g, B: b, A: 255}
}

// RGBA constructs a Color with custom Alpha.
func RGBA(r, g, b, a uint8) Color {
	return Color{R: r, G: g, B: b, A: a}
}

// Hex parses a hex color string (e.g., "#RGB", "#RGBA", "#RRGGBB", "#RRGGBBAA", or without "#").
func Hex(s string) (Color, error) {
	s = strings.TrimPrefix(strings.TrimSpace(s), "#")
	switch len(s) {
	case 3: // RGB
		r, err := strconv.ParseUint(string(s[0])+string(s[0]), 16, 8)
		if err != nil {
			return Color{}, fmt.Errorf("types.Hex: invalid hex %q: %w", s, err)
		}
		g, err := strconv.ParseUint(string(s[1])+string(s[1]), 16, 8)
		if err != nil {
			return Color{}, fmt.Errorf("types.Hex: invalid hex %q: %w", s, err)
		}
		b, err := strconv.ParseUint(string(s[2])+string(s[2]), 16, 8)
		if err != nil {
			return Color{}, fmt.Errorf("types.Hex: invalid hex %q: %w", s, err)
		}
		return RGB(uint8(r), uint8(g), uint8(b)), nil

	case 6: // RRGGBB
		val, err := strconv.ParseUint(s, 16, 32)
		if err != nil {
			return Color{}, fmt.Errorf("types.Hex: invalid hex %q: %w", s, err)
		}
		return RGB(uint8(val>>16), uint8((val>>8)&0xFF), uint8(val&0xFF)), nil

	case 8: // RRGGBBAA
		val, err := strconv.ParseUint(s, 16, 32)
		if err != nil {
			return Color{}, fmt.Errorf("types.Hex: invalid hex %q: %w", s, err)
		}
		return RGBA(uint8(val>>24), uint8((val>>16)&0xFF), uint8((val>>8)&0xFF), uint8(val&0xFF)), nil

	default:
		return Color{}, fmt.Errorf("types.Hex: unsupported hex format %q", s)
	}
}

// MustHex parses hex string and panics on error.
func MustHex(s string) Color {
	c, err := Hex(s)
	if err != nil {
		panic(err)
	}
	return c
}

// FFmpegColor formats the color for FFmpeg filters (such as color=c=..., drawbox, drawtext).
// For opaque colors: "0xRRGGBB"
// For transparent colors: "0xRRGGBB@A.AA" (where Alpha is 0.0 to 1.0)
func (c Color) FFmpegColor() string {
	if c.A == 0 {
		return "black@0.0"
	}
	if c.A == 255 {
		return fmt.Sprintf("0x%02X%02X%02X", c.R, c.G, c.B)
	}
	alphaFloat := float64(c.A) / 255.0
	return fmt.Sprintf("0x%02X%02X%02X@%.2f", c.R, c.G, c.B, alphaFloat)
}

// HexRGB returns "#RRGGBB".
func (c Color) HexRGB() string {
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
}

// HexRGBA returns "#RRGGBBAA".
func (c Color) HexRGBA() string {
	return fmt.Sprintf("#%02X%02X%02X%02X", c.R, c.G, c.B, c.A)
}
