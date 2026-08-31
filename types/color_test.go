package types_test

import (
	"testing"

	"github.com/farshidrezaei/vidonyx/types"
)

func TestColor_HexParsing(t *testing.T) {
	c1, err := types.Hex("#FFAA00")
	if err != nil {
		t.Fatalf("Hex #FFAA00 failed: %v", err)
	}
	if c1.R != 255 || c1.G != 170 || c1.B != 0 || c1.A != 255 {
		t.Fatalf("Hex parsed incorrectly: %+v", c1)
	}

	c2, err := types.Hex("00FF0080") // RGBA with 50% alpha
	if err != nil {
		t.Fatalf("Hex 00FF0080 failed: %v", err)
	}
	if c2.R != 0 || c2.G != 255 || c2.B != 0 || c2.A != 128 {
		t.Fatalf("Hex parsed incorrectly: %+v", c2)
	}
}

func TestColor_FFmpegColor(t *testing.T) {
	black := types.ColorBlack
	if s := black.FFmpegColor(); s != "0x000000" {
		t.Fatalf("Black FFmpegColor: expected 0x000000, got %s", s)
	}

	transparent := types.ColorTransparent
	if s := transparent.FFmpegColor(); s != "black@0.0" {
		t.Fatalf("Transparent FFmpegColor: expected black@0.0, got %s", s)
	}

	whiteHalf := types.RGBA(255, 255, 255, 128)
	if s := whiteHalf.FFmpegColor(); s != "0xFFFFFF@0.50" {
		t.Fatalf("Half white FFmpegColor: expected 0xFFFFFF@0.50, got %s", s)
	}
}
