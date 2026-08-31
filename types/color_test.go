package types_test

import (
	"testing"

	"github.com/farshidrezaei/vidonyx/types"
)

func TestColor_HexParsingTable(t *testing.T) {
	tests := []struct {
		name      string
		hex       string
		wantR     uint8
		wantG     uint8
		wantB     uint8
		wantA     uint8
		shouldErr bool
	}{
		{name: "3-char shorthand", hex: "#FFF", wantR: 255, wantG: 255, wantB: 255, wantA: 255},
		{name: "3-char no prefix", hex: "F00", wantR: 255, wantG: 0, wantB: 0, wantA: 255},
		{name: "6-char full", hex: "#1A2B3C", wantR: 0x1A, wantG: 0x2B, wantB: 0x3C, wantA: 255},
		{name: "6-char no prefix", hex: "00FF00", wantR: 0, wantG: 255, wantB: 0, wantA: 255},
		{name: "8-char with alpha", hex: "#11223380", wantR: 0x11, wantG: 0x22, wantB: 0x33, wantA: 0x80},
		{name: "invalid length", hex: "#1234", shouldErr: true},
		{name: "invalid characters", hex: "#ZZZZZZ", shouldErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := types.Hex(tt.hex)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("Hex(%q) expected error, got nil", tt.hex)
				}
				return
			}
			if err != nil {
				t.Fatalf("Hex(%q) unexpected error: %v", tt.hex, err)
			}
			if c.R != tt.wantR || c.G != tt.wantG || c.B != tt.wantB || c.A != tt.wantA {
				t.Errorf("Hex(%q) = RGBA(%d,%d,%d,%d), want RGBA(%d,%d,%d,%d)",
					tt.hex, c.R, c.G, c.B, c.A, tt.wantR, tt.wantG, tt.wantB, tt.wantA)
			}
		})
	}
}

func TestColor_FFmpegColorTable(t *testing.T) {
	tests := []struct {
		name     string
		color    types.Color
		expected string
	}{
		{name: "fully transparent", color: types.ColorTransparent, expected: "black@0.0"},
		{name: "opaque black", color: types.ColorBlack, expected: "0x000000"},
		{name: "opaque white", color: types.ColorWhite, expected: "0xFFFFFF"},
		{name: "opaque red", color: types.ColorRed, expected: "0xFF0000"},
		{name: "semi-transparent white", color: types.RGBA(255, 255, 255, 128), expected: "0xFFFFFF@0.50"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.color.FFmpegColor(); got != tt.expected {
				t.Errorf("FFmpegColor() = %q, want %q", got, tt.expected)
			}
		})
	}
}
