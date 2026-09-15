package waveform_test

import (
	"strings"
	"testing"

	"github.com/farshidrezaei/vidonex/filtergraph"
	"github.com/farshidrezaei/vidonex/types"
	"github.com/farshidrezaei/vidonex/waveform"
)

func TestApplyWaveformVisualizerTable(t *testing.T) {
	tests := []struct {
		name         string
		options      waveform.Options
		wantSnippets []string
	}{
		{
			name:    "default waveform options with dual neon colors",
			options: waveform.DefaultOptions(),
			wantSnippets: []string{
				"showwaves=s=1200x200",
				"mode=p2p",
				"scale=sqrt",
				"r=30",
				"format=pix_fmts=yuva420p",
			},
		},
		{
			name: "custom 60fps line mode",
			options: waveform.Options{
				Size:     types.NewSize(800, 150),
				Mode:     waveform.ModeLine,
				Color:    types.ColorRed,
				Scale:    "log",
				FPS:      types.FPS60,
				Position: types.Point{X: 100, Y: 500},
			},
			wantSnippets: []string{
				"showwaves=s=800x150",
				"mode=line",
				"scale=log",
				"r=60",
			},
		},
		{
			name: "oscilloscope wave mode maps to cline",
			options: waveform.Options{
				Size:           types.NewSize(1000, 180),
				Mode:           waveform.ModeWave,
				Color:          types.RGB(0, 230, 180),
				SecondaryColor: types.RGB(99, 102, 241),
				Scale:          "lin",
				FPS:            types.FPS30,
				Position:       types.Point{X: 0, Y: 400},
			},
			wantSnippets: []string{
				"showwaves=s=1000x180",
				"mode=cline",
				"scale=lin",
			},
		},
		{
			name: "spectrum bars mode maps to p2p",
			options: waveform.Options{
				Size:      types.NewSize(900, 220),
				Mode:      waveform.ModeSpectrum,
				Color:     types.RGB(255, 51, 102),
				Density:   64,
				Roundness: 8,
				Glow:      true,
			},
			wantSnippets: []string{
				"showwaves=s=900x220",
				"mode=p2p",
			},
		},
		{
			name: "frequency dot mode maps to dot",
			options: waveform.Options{
				Size:    types.NewSize(600, 120),
				Mode:    waveform.ModeDot,
				Color:   types.RGB(0, 255, 255),
				Density: 32,
			},
			wantSnippets: []string{
				"showwaves=s=600x120",
				"mode=dot",
			},
		},
		{
			name: "circular audiogram maps to p2p",
			options: waveform.Options{
				Size:  types.NewSize(500, 500),
				Mode:  waveform.ModeCircular,
				Color: types.RGB(168, 85, 247),
			},
			wantSnippets: []string{
				"showwaves=s=500x500",
				"mode=p2p",
			},
		},
		{
			name:    "zero-value options fallback to sensible defaults",
			options: waveform.Options{},
			wantSnippets: []string{
				"showwaves=s=1200x200",
				"mode=p2p",
				"scale=sqrt",
				"r=30",
				"format=pix_fmts=yuva420p",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graph := filtergraph.NewGraph()
			dummyAudioIn := &filtergraph.Pad{
				ID:         "speech_audio_in",
				StreamType: filtergraph.StreamTypeAudio,
				IsInput:    false,
			}

			outVideoPad, err := waveform.ApplyWaveformVisualizer(graph, "waves_test", dummyAudioIn, tt.options)
			if err != nil {
				t.Fatalf("ApplyWaveformVisualizer failed: %v", err)
			}

			if outVideoPad == nil {
				t.Fatalf("expected non-nil output video pad")
			}

			filterComplex, err := graph.FormattedFilterComplex()
			if err != nil {
				t.Fatalf("FormattedFilterComplex failed: %v", err)
			}

			for _, snippet := range tt.wantSnippets {
				if !strings.Contains(filterComplex, snippet) {
					t.Errorf("filterComplex missing %q\nFull string: %s", snippet, filterComplex)
				}
			}
		})
	}
}
