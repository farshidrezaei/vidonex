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
func (streamType StreamType) String() string {
	switch streamType {
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
func (pad *Pad) Label() string {
	if pad == nil || pad.ID == "" {
		return ""
	}
	return fmt.Sprintf("[%s]", pad.ID)
}

// String returns the string representation of the Pad.
func (pad *Pad) String() string {
	direction := "out"
	if pad.IsInput {
		direction = "in"
	}
	nodeName := "none"
	if pad.Node != nil {
		nodeName = pad.Node.ID
	}
	return fmt.Sprintf("Pad(%s, %s, %s, node=%s)", pad.ID, pad.StreamType, direction, nodeName)
}
