package chromakey_test

import (
	"strings"
	"testing"

	"github.com/farshidrezaei/vidonex/chromakey"
	"github.com/farshidrezaei/vidonex/filtergraph"
	"github.com/farshidrezaei/vidonex/types"
)

func TestApplyChromaKeyTable(t *testing.T) {
	tests := []struct {
		name         string
		options      chromakey.Options
		wantSnippets []string
	}{
		{
			name:    "default green screen with despill",
			options: chromakey.DefaultOptions(),
			wantSnippets: []string{
				"chromakey=color=0x00FF00:similarity=0.30:blend=0.10",
				"despill=type=green:expand=0.00",
			},
		},
		{
			name: "blue screen without despill",
			options: chromakey.Options{
				KeyColor:    types.RGB(0, 0, 255),
				Similarity:  0.45,
				Blend:       0.05,
				Despill:     false,
				DespillType: chromakey.DespillBlue,
			},
			wantSnippets: []string{
				"chromakey=color=0x0000FF:similarity=0.45:blend=0.05",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graph := filtergraph.NewGraph()
			dummyVideoIn := &filtergraph.Pad{
				ID:         "camera_green_screen_in",
				StreamType: filtergraph.StreamTypeVideo,
				IsInput:    false,
			}

			outVideoPad, err := chromakey.ApplyChromaKey(graph, "chroma_test", dummyVideoIn, tt.options)
			if err != nil {
				t.Fatalf("ApplyChromaKey failed: %v", err)
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
