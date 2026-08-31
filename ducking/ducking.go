// Package ducking provides automated audio ducking and sidechain compression for lowering background music during dialogue.
package ducking

import (
	"fmt"

	"github.com/farshidrezaei/vidonyx/filtergraph"
)

// Options configures sidechain compression parameters.
type Options struct {
	Threshold           float64 // Detection threshold (0.001 to 1.0, default 0.1)
	Ratio               float64 // Compression ratio (default 4.0)
	AttackMilliseconds  int     // Attack time in ms (default 20ms)
	ReleaseMilliseconds int     // Release time in ms (default 300ms)
	MakeupGain          float64 // Post-compression makeup gain (default 1.0)
}

// DefaultOptions returns recommended natural-sounding voiceover ducking parameters.
func DefaultOptions() Options {
	return Options{
		Threshold:           0.1,
		Ratio:               4.0,
		AttackMilliseconds:  20,
		ReleaseMilliseconds: 300,
		MakeupGain:          1.0,
	}
}

// ApplySidechainDucking connects a music audio pad and a voice audio pad into a sidechaincompress filter node.
func ApplySidechainDucking(graph *filtergraph.Graph, nodeID string, musicPad, voicePad *filtergraph.Pad, options Options) (*filtergraph.Pad, error) {
	if options.Ratio <= 0 {
		options.Ratio = 4.0
	}
	if options.AttackMilliseconds <= 0 {
		options.AttackMilliseconds = 20
	}
	if options.ReleaseMilliseconds <= 0 {
		options.ReleaseMilliseconds = 300
	}
	if options.Threshold <= 0 {
		options.Threshold = 0.1
	}

	sidechainNode := graph.NewNode(nodeID, "sidechaincompress")
	sidechainNode.SetParam("threshold", fmt.Sprintf("%.3f", options.Threshold))
	sidechainNode.SetParam("ratio", fmt.Sprintf("%.1f", options.Ratio))
	sidechainNode.SetParam("attack", options.AttackMilliseconds)
	sidechainNode.SetParam("release", options.ReleaseMilliseconds)

	// Input 0: Main audio to be ducked (Music/BGM)
	// Input 1: Trigger/Key audio (Voice/Dialogue)
	inputMusic := sidechainNode.AddInput(musicPad.ID, filtergraph.StreamTypeAudio)
	inputVoice := sidechainNode.AddInput(voicePad.ID, filtergraph.StreamTypeAudio)
	outputDuckedMusic := sidechainNode.AddOutput(graph.NextPadID("ducked_music"), filtergraph.StreamTypeAudio)

	if err := graph.Connect(musicPad, inputMusic); err != nil {
		return nil, fmt.Errorf("ducking: failed connecting music pad: %w", err)
	}
	if err := graph.Connect(voicePad, inputVoice); err != nil {
		return nil, fmt.Errorf("ducking: failed connecting voice trigger pad: %w", err)
	}

	return outputDuckedMusic, nil
}
