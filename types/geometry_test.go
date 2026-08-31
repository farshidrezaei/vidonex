package types_test

import (
	"testing"

	"github.com/farshidrezaei/vidonyx/types"
)

func TestSize_AspectAndScaleTable(t *testing.T) {
	tests := []struct {
		name       string
		size       types.Size
		scale      float64
		expectedAR types.Rational
		expectedSz types.Size
	}{
		{
			name:       "1080p 16:9 scale half",
			size:       types.Res1080p,
			scale:      0.5,
			expectedAR: types.NewRational(16, 9),
			expectedSz: types.NewSize(960, 540),
		},
		{
			name:       "720p 16:9 scale double",
			size:       types.Res720p,
			scale:      2.0,
			expectedAR: types.NewRational(16, 9),
			expectedSz: types.NewSize(2560, 1440),
		},
		{
			name:       "vertical 9:16",
			size:       types.ResPortrait1080p,
			scale:      0.5,
			expectedAR: types.NewRational(9, 16),
			expectedSz: types.NewSize(540, 960),
		},
		{
			name:       "square 1:1",
			size:       types.ResSquare1080,
			scale:      0.5,
			expectedAR: types.NewRational(1, 1),
			expectedSz: types.NewSize(540, 540),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ar := tt.size.AspectRatio()
			if !ar.Equal(tt.expectedAR) {
				t.Errorf("AspectRatio(%v) = %v, want %v", tt.size, ar, tt.expectedAR)
			}

			scaled := tt.size.Scale(tt.scale)
			if scaled != tt.expectedSz {
				t.Errorf("Scale(%v, %f) = %v, want %v", tt.size, tt.scale, scaled, tt.expectedSz)
			}
		})
	}
}

func TestSize_FitAndFillTable(t *testing.T) {
	tests := []struct {
		name     string
		source   types.Size
		target   types.Size
		wantFit  types.Size
		wantFill types.Size
	}{
		{
			name:     "wider into narrower 16:9",
			source:   types.NewSize(1000, 500), // 2:1
			target:   types.Res1080p,          // 1920x1080 (16:9)
			wantFit:  types.NewSize(1920, 960),
			wantFill: types.NewSize(2160, 1080),
		},
		{
			name:     "taller into wider 16:9",
			source:   types.NewSize(1080, 1920), // 9:16 vertical
			target:   types.Res1080p,           // 1920x1080 (16:9)
			wantFit:  types.NewSize(608, 1080),
			wantFill: types.NewSize(1920, 3414), // makeEven will produce even bounds
		},
		{
			name:     "exact same aspect ratio",
			source:   types.Res720p,
			target:   types.Res1080p,
			wantFit:  types.Res1080p,
			wantFill: types.Res1080p,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fit := tt.source.FitWithin(tt.target)
			if fit.Width != tt.wantFit.Width || fit.Height != tt.wantFit.Height {
				t.Errorf("FitWithin(%v into %v) = %v, want %v", tt.source, tt.target, fit, tt.wantFit)
			}

			fill := tt.source.FillWithin(tt.target)
			if fill.Width != tt.wantFill.Width || fill.Height != tt.wantFill.Height {
				t.Errorf("FillWithin(%v into %v) = %v, want %v", tt.source, tt.target, fill, tt.wantFill)
			}
		})
	}
}

func TestCalculateOffsetTable(t *testing.T) {
	container := types.NewSize(1920, 1080)
	item := types.NewSize(400, 200)

	tests := []struct {
		name      string
		alignment types.Alignment
		wantX     int
		wantY     int
	}{
		{name: "TopLeft", alignment: types.AlignTopLeft, wantX: 0, wantY: 0},
		{name: "TopCenter", alignment: types.AlignTopCenter, wantX: 760, wantY: 0},
		{name: "TopRight", alignment: types.AlignTopRight, wantX: 1520, wantY: 0},
		{name: "CenterLeft", alignment: types.AlignCenterLeft, wantX: 0, wantY: 440},
		{name: "Center", alignment: types.AlignCenter, wantX: 760, wantY: 440},
		{name: "CenterRight", alignment: types.AlignCenterRight, wantX: 1520, wantY: 440},
		{name: "BottomLeft", alignment: types.AlignBottomLeft, wantX: 0, wantY: 880},
		{name: "BottomCenter", alignment: types.AlignBottomCenter, wantX: 760, wantY: 880},
		{name: "BottomRight", alignment: types.AlignBottomRight, wantX: 1520, wantY: 880},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pt := types.CalculateOffset(container, item, tt.alignment)
			if pt.X != tt.wantX || pt.Y != tt.wantY {
				t.Errorf("CalculateOffset(%v, %v, %v) = (%d, %d), want (%d, %d)",
					container, item, tt.alignment, pt.X, pt.Y, tt.wantX, tt.wantY)
			}
		})
	}
}

func TestRect_OriginAndSize(t *testing.T) {
	r := types.Rect{X: 10, Y: 20, Width: 300, Height: 200}
	if r.Origin() != (types.Point{X: 10, Y: 20}) {
		t.Errorf("Rect.Origin() = %v, want (10, 20)", r.Origin())
	}
	if r.Size() != (types.Size{Width: 300, Height: 200}) {
		t.Errorf("Rect.Size() = %v, want 300x200", r.Size())
	}
}
