package filtergraph

import (
	"fmt"
)

// StreamType defines whether a pad carries video or audio frames.
type StreamType uint8

const (
	// StreamTypeVideo denotes a video frame stream.
	StreamTypeVideo StreamType = iota
	// StreamTypeAudio denotes an audio frame stream.
	StreamTypeAudio
)

// String returns a human-readable stream type name.
func (st StreamType) String() string {
	switch st {
	case StreamTypeVideo:
		return "video"
	case StreamTypeAudio:
		return "audio"
	default:
		return "unknown"
	}
}

// Pad represents an input or output terminal on a filter node.
type Pad struct {
	ID         string
	StreamType StreamType
	Node       *Node
	IsInput    bool
}

// Label returns the bracketed pad name used in FFmpeg filtergraphs (e.g. "[v0]" or "[0:v]").
func (p *Pad) Label() string {
	if p == nil || p.ID == "" {
		return ""
	}
	return fmt.Sprintf("[%s]", p.ID)
}

// String returns the string representation of the Pad.
func (p *Pad) String() string {
	dir := "out"
	if p.IsInput {
		dir = "in"
	}
	nodeName := "none"
	if p.Node != nil {
		nodeName = p.Node.ID
	}
	return fmt.Sprintf("Pad(%s, %s, %s, node=%s)", p.ID, p.StreamType, dir, nodeName)
}
