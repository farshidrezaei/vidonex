package effects_test

import (
	"testing"

	"github.com/farshidrezaei/vidonex/effects"
	"github.com/farshidrezaei/vidonex/filtergraph"
)

func TestLoudnormFilter_Table(t *testing.T) {
	testCases := []struct {
		name                   string
		filter                 effects.LoudnormFilter
		sourceStreamType       filtergraph.StreamType
		expectedIntegratedLUFS string
		expectedTruePeak       string
		expectedError          bool
	}{
		{
			name:                   "Spotify and YouTube preset",
			filter:                 effects.LoudnormSpotifyYouTube(),
			sourceStreamType:       filtergraph.StreamTypeAudio,
			expectedIntegratedLUFS: "-14.0",
			expectedTruePeak:       "-1.0",
			expectedError:          false,
		},
		{
			name:                   "EBU R128 European broadcast preset",
			filter:                 effects.LoudnormEBUR128(),
			sourceStreamType:       filtergraph.StreamTypeAudio,
			expectedIntegratedLUFS: "-23.0",
			expectedTruePeak:       "-1.0",
			expectedError:          false,
		},
		{
			name:                   "Podcast vocal speech preset",
			filter:                 effects.LoudnormPodcast(),
			sourceStreamType:       filtergraph.StreamTypeAudio,
			expectedIntegratedLUFS: "-16.0",
			expectedTruePeak:       "-1.5",
			expectedError:          false,
		},
		{
			name:                   "Apple Music preset",
			filter:                 effects.LoudnormAppleMusic(),
			sourceStreamType:       filtergraph.StreamTypeAudio,
			expectedIntegratedLUFS: "-16.0",
			expectedTruePeak:       "-1.0",
			expectedError:          false,
		},
		{
			name: "Custom loudnorm settings with dual mono",
			filter: effects.LoudnormFilter{
				IntegratedLoudness: -18.0,
				LoudnessRange:      8.5,
				TruePeak:           -2.0,
				Linear:             true,
				DualMono:           true,
			},
			sourceStreamType:       filtergraph.StreamTypeAudio,
			expectedIntegratedLUFS: "-18.0",
			expectedTruePeak:       "-2.0",
			expectedError:          false,
		},
		{
			name:                   "Stream type mismatch error (video pad to audio filter)",
			filter:                 effects.LoudnormSpotifyYouTube(),
			sourceStreamType:       filtergraph.StreamTypeVideo,
			expectedIntegratedLUFS: "",
			expectedTruePeak:       "",
			expectedError:          true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			graph := filtergraph.NewGraph()
			sourcePad := &filtergraph.Pad{
				ID:         "source_audio_pad",
				StreamType: tc.sourceStreamType,
			}

			outputPad, err := tc.filter.Apply(graph, "loudnorm_node", sourcePad)

			if tc.expectedError {
				if err == nil {
					t.Fatalf("expected error for scenario %q, got nil", tc.name)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error for scenario %q: %v", tc.name, err)
			}

			if outputPad == nil {
				t.Fatalf("expected valid output pad, got nil")
			}

			if outputPad.StreamType != filtergraph.StreamTypeAudio {
				t.Errorf("expected output pad stream type Audio, got %s", outputPad.StreamType)
			}

			node, exists := graph.GetNode("loudnorm_node")
			if !exists || node == nil {
				t.Fatalf("expected node loudnorm_node to be present in graph")
			}

			if node.FilterName != "loudnorm" {
				t.Errorf("expected filter name loudnorm, got %s", node.FilterName)
			}

			if val, ok := node.GetParam("i"); !ok || val != tc.expectedIntegratedLUFS {
				t.Errorf("expected param i=%s, got %v", tc.expectedIntegratedLUFS, val)
			}

			if val, ok := node.GetParam("tp"); !ok || val != tc.expectedTruePeak {
				t.Errorf("expected param tp=%s, got %v", tc.expectedTruePeak, val)
			}
		})
	}
}
