package types_test

import (
	"testing"

	"github.com/farshidrezaei/vidonyx/types"
)

func TestSize_AspectAndScale(t *testing.T) {
	s := types.Res1080p // 1920x1080
	ar := s.AspectRatio()
	if !ar.Equal(types.NewRational(16, 9)) {
		t.Fatalf("AspectRatio: expected 16/9, got %s", ar)
	}

	scaled := s.Scale(0.5)
	if scaled.Width != 960 || scaled.Height != 540 {
		t.Fatalf("Scale: expected 960x540, got %dx%d", scaled.Width, scaled.Height)
	}
}

func TestSize_FitAndFill(t *testing.T) {
	source := types.NewSize(1000, 500) // 2:1
	target := types.Res1080p           // 1920x1080 (16:9)

	fit := source.FitWithin(target)
	if fit.Width != 1920 || fit.Height != 960 {
		t.Fatalf("FitWithin: expected 1920x960, got %dx%d", fit.Width, fit.Height)
	}

	fill := source.FillWithin(target)
	if fill.Width != 2160 || fill.Height != 1080 {
		t.Fatalf("FillWithin: expected 2160x1080, got %dx%d", fill.Width, fill.Height)
	}
}

func TestCalculateOffset(t *testing.T) {
	container := types.NewSize(1920, 1080)
	item := types.NewSize(400, 300)

	center := types.CalculateOffset(container, item, types.AlignCenter)
	if center.X != 760 || center.Y != 390 {
		t.Fatalf("AlignCenter: expected (760, 390), got (%d, %d)", center.X, center.Y)
	}

	bottomRight := types.CalculateOffset(container, item, types.AlignBottomRight)
	if bottomRight.X != 1520 || bottomRight.Y != 780 {
		t.Fatalf("AlignBottomRight: expected (1520, 780), got (%d, %d)", bottomRight.X, bottomRight.Y)
	}
}
