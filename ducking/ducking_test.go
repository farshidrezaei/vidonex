package ducking_test

import (
	"strings"
	"testing"

	"github.com/farshidrezaei/vidonex/ducking"
	"github.com/farshidrezaei/vidonex/filtergraph"
)

func TestApplySidechainDuckingTable(t *testing.T) {
	tests := []struct {
		name         string
		options      ducking.Options
		wantSnippets []string
	}{
		{
			name:    "default options",
			options: ducking.DefaultOptions(),
			wantSnippets: []string{
				"sidechaincompress=",
				"threshold=0.100",
				"ratio=4.0",
				"attack=20",
				"release=300",
			},
		},
		{
			name: "custom aggressive ducking",
			options: ducking.Options{
				Threshold:           0.05,
				Ratio:               8.0,
				AttackMilliseconds:  10,
				ReleaseMilliseconds: 500,
			},
			wantSnippets: []string{
				"threshold=0.050",
				"ratio=8.0",
				"attack=10",
				"release=500",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graph := filtergraph.NewGraph()
			musicPad := &filtergraph.Pad{
				ID:         "music_in",
				StreamType: filtergraph.StreamTypeAudio,
				IsInput:    false,
			}
			voicePad := &filtergraph.Pad{
				ID:         "voice_in",
				StreamType: filtergraph.StreamTypeAudio,
				IsInput:    false,
			}

			outPad, err := ducking.ApplySidechainDucking(graph, "duck_node", musicPad, voicePad, tt.options)
			if err != nil {
				t.Fatalf("ApplySidechainDucking failed: %v", err)
			}

			if outPad == nil {
				t.Fatalf("expected non-nil output pad")
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
