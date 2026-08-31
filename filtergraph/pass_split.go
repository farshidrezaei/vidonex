package filtergraph

import (
	"fmt"
)

// AutoSplitPass detects output pads connected to multiple consumers (fan-out > 1)
// and automatically injects 'split' or 'asplit' filter nodes.
type AutoSplitPass struct{}

func (p *AutoSplitPass) Name() string {
	return "AutoSplitPass"
}

func (p *AutoSplitPass) Run(g *Graph) error {
	// Find all output pads that have more than 1 consumer
	type splitTarget struct {
		outPad    *Pad
		consumers []*Pad
	}

	targets := make([]splitTarget, 0)
	for node := range g.Nodes() {
		for _, outPad := range node.Outputs {
			consumers := g.GetConsumerPads(outPad)
			if len(consumers) > 1 {
				targets = append(targets, splitTarget{
					outPad:    outPad,
					consumers: append([]*Pad(nil), consumers...),
				})
			}
		}
	}

	for _, target := range targets {
		nConsumers := len(target.consumers)
		filterName := "split"
		prefix := "v_split"
		if target.outPad.StreamType == StreamTypeAudio {
			filterName = "asplit"
			prefix = "a_split"
		}

		splitNodeID := g.NextPadID(prefix)
		splitNode := NewNode(splitNodeID, filterName)
		if nConsumers > 2 {
			splitNode.SetParam("", nConsumers)
		}

		// Input pad on the split node
		splitIn := splitNode.AddInput(target.outPad.ID, target.outPad.StreamType)
		if err := g.AddNode(splitNode); err != nil {
			return err
		}

		// Connect original outPad to splitIn
		// First disconnect all original consumers
		for _, consumerIn := range target.consumers {
			_ = g.Disconnect(consumerIn)
		}
		if err := g.Connect(target.outPad, splitIn); err != nil {
			return fmt.Errorf("failed connecting source to split node: %w", err)
		}

		// Create N output pads on the split node and connect each to one consumer
		for i, consumerIn := range target.consumers {
			outPadID := g.NextPadID(fmt.Sprintf("%s_out_%d", prefix, i+1))
			splitOut := splitNode.AddOutput(outPadID, target.outPad.StreamType)
			if err := g.Connect(splitOut, consumerIn); err != nil {
				return fmt.Errorf("failed connecting split out to consumer: %w", err)
			}
		}
	}

	return nil
}
