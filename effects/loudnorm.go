// Package effects provides pre-built, type-safe video and audio filters for FFmpeg filtergraphs.
package effects

import (
	"fmt"

	"github.com/farshidrezaei/vidonex/filtergraph"
)

// LoudnormFilter implements the industry-standard EBU R128 audio loudness normalization filter.
// It normalizes audio to streaming and broadcast loudness targets (e.g., YouTube, Spotify, Apple Podcasts, EBU).
type LoudnormFilter struct {
	// IntegratedLoudness is the target integrated loudness in LUFS (e.g., -14.0 for YouTube, -23.0 for EBU R128).
	IntegratedLoudness float64
	// LoudnessRange is the target loudness range in LU (LRA, typically between 7.0 and 11.0).
	LoudnessRange float64
	// TruePeak is the maximum allowed true peak level in dBFS (typically -1.0 or -1.5).
	TruePeak float64
	// Linear enables linear normalization when set to true.
	Linear bool
	// DualMono treats mono audio streams as dual mono when enabled.
	DualMono bool
}

// LoudnormSpotifyYouTube returns a preset configured for YouTube and Spotify streaming (-14.0 LUFS, -1.0 dBFS TP).
func LoudnormSpotifyYouTube() LoudnormFilter {
	return LoudnormFilter{
		IntegratedLoudness: -14.0,
		LoudnessRange:      11.0,
		TruePeak:           -1.0,
		Linear:             true,
		DualMono:           false,
	}
}

// LoudnormEBUR128 returns a preset configured for European broadcast standards (-23.0 LUFS, -1.0 dBFS TP).
func LoudnormEBUR128() LoudnormFilter {
	return LoudnormFilter{
		IntegratedLoudness: -23.0,
		LoudnessRange:      7.0,
		TruePeak:           -1.0,
		Linear:             true,
		DualMono:           false,
	}
}

// LoudnormPodcast returns a preset optimized for voice and spoken-word podcasts (-16.0 LUFS, -1.5 dBFS TP).
func LoudnormPodcast() LoudnormFilter {
	return LoudnormFilter{
		IntegratedLoudness: -16.0,
		LoudnessRange:      9.0,
		TruePeak:           -1.5,
		Linear:             true,
		DualMono:           false,
	}
}

// LoudnormAppleMusic returns a preset configured for Apple Music and Apple Podcasts (-16.0 LUFS, -1.0 dBFS TP).
func LoudnormAppleMusic() LoudnormFilter {
	return LoudnormFilter{
		IntegratedLoudness: -16.0,
		LoudnessRange:      11.0,
		TruePeak:           -1.0,
		Linear:             true,
		DualMono:           false,
	}
}

// Apply attaches an EBU R128 loudnorm filter node to the given graph.
func (filter LoudnormFilter) Apply(
	graph *filtergraph.Graph,
	nodeIdentifier string,
	sourceOutputPad *filtergraph.Pad,
) (*filtergraph.Pad, error) {
	node := graph.NewNode(nodeIdentifier, "loudnorm")

	integratedLoudnessTarget := filter.IntegratedLoudness
	if integratedLoudnessTarget == 0 {
		integratedLoudnessTarget = -24.0
	}

	loudnessRangeTarget := filter.LoudnessRange
	if loudnessRangeTarget == 0 {
		loudnessRangeTarget = 7.0
	}

	truePeakTarget := filter.TruePeak
	if truePeakTarget == 0 {
		truePeakTarget = -2.0
	}

	node.SetParam("i", fmt.Sprintf("%.1f", integratedLoudnessTarget))
	node.SetParam("lra", fmt.Sprintf("%.1f", loudnessRangeTarget))
	node.SetParam("tp", fmt.Sprintf("%.1f", truePeakTarget))

	if filter.Linear {
		node.SetParam("linear", "true")
	} else {
		node.SetParam("linear", "false")
	}

	if filter.DualMono {
		node.SetParam("dual_mono", "true")
	}

	destinationInputPad := node.AddInput(sourceOutputPad.ID, filtergraph.StreamTypeAudio)
	outputPad := node.AddOutput(graph.NextPadID("loudnorm_out"), filtergraph.StreamTypeAudio)

	if err := graph.Connect(sourceOutputPad, destinationInputPad); err != nil {
		return nil, err
	}

	return outputPad, nil
}
