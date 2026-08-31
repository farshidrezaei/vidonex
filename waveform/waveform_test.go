package waveform_test

import (
	"strings"
	"testing"

	"github.com/farshidrezaei/vidonyx/filtergraph"
	"github.com/farshidrezaei/vidonyx/types"
	"github.com/farshidrezaei/vidonyx/waveform"
)

func TestApplyWaveformVisualizerTable(t *testing.T) {
	tests := []struct {
		name         string
		options      waveform.Options
		wantSnippets []string
	}{
		{
			name:    "default waveform options",
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
