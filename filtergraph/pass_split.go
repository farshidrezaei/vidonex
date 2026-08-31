package filtergraph

import (
	"fmt"
)

// AutoSplitPass detects output pads connected to multiple consumers (fan-out > 1)
// and automatically injects 'split' or 'asplit' filter nodes.
type AutoSplitPass struct{}

// Name returns the unique identifier name of the pass.
func (pass *AutoSplitPass) Name() string {
	return "AutoSplitPass"
}

// Run executes the auto-split injection pass across the graph.
func (pass *AutoSplitPass) Run(graph *Graph) error {
	// Find all output pads that have more than 1 consumer
	type splitTarget struct {
		outputPad *Pad
		consumers []*Pad
	}

	targets := make([]splitTarget, 0)
	for node := range graph.Nodes() {
		for _, outputPad := range node.Outputs {
			consumers := graph.GetConsumerPads(outputPad)
			if len(consumers) > 1 {
				targets = append(targets, splitTarget{
					outputPad: outputPad,
					consumers: append([]*Pad(nil), consumers...),
				})
			}
		}
	}

	for _, target := range targets {
		consumerCount := len(target.consumers)
		filterName := "split"
		prefix := "video_split"
		if target.outputPad.StreamType == StreamTypeAudio {
			filterName = "asplit"
			prefix = "audio_split"
		}

		splitNodeID := graph.NextPadID(prefix)
		splitNode := NewNode(splitNodeID, filterName)
		if consumerCount > 2 {
			splitNode.SetParam("", consumerCount)
		}

		// Input pad on the split node
		splitInput := splitNode.AddInput(target.outputPad.ID, target.outputPad.StreamType)
		if err := graph.AddNode(splitNode); err != nil {
			return err
		}

		// Connect original outputPad to splitInput
		for _, consumerInput := range target.consumers {
			_ = graph.Disconnect(consumerInput)
		}
		if err := graph.Connect(target.outputPad, splitInput); err != nil {
			return fmt.Errorf("failed connecting source to split node: %w", err)
		}

		// Create N output pads on the split node and connect each to one consumer
		for index, consumerInput := range target.consumers {
			outputPadID := graph.NextPadID(fmt.Sprintf("%s_out_%d", prefix, index+1))
			splitOutput := splitNode.AddOutput(outputPadID, target.outputPad.StreamType)
			if err := graph.Connect(splitOutput, consumerInput); err != nil {
				return fmt.Errorf("failed connecting split out to consumer: %w", err)
			}
		}
	}

	return nil
}
