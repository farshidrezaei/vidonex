package effects_test

import (
	"testing"

	"github.com/farshidrezaei/vidonex/effects"
	"github.com/farshidrezaei/vidonex/filtergraph"
	"github.com/farshidrezaei/vidonex/timeline"
	"github.com/farshidrezaei/vidonex/types"
)

func TestEffects_VideoFilters(t *testing.T) {
	g := filtergraph.NewGraph()
	rawIn := &filtergraph.Pad{ID: "0:v", StreamType: filtergraph.StreamTypeVideo}

	// Test Scale
	scaleFilter := effects.ScaleFilter{Width: 1280, Height: 720, AspectRatio: "decrease"}
	outPad, err := scaleFilter.Apply(g, "scale_1", rawIn)
	if err != nil {
		t.Fatalf("ScaleFilter failed: %v", err)
	}

	// Test DrawText
	textFilter := effects.DrawTextFilter{
		Text:      "Vidonex Engine",
		FontSize:  32,
		FontColor: types.ColorWhite,
		Box:       true,
		BoxColor:  types.ColorBlack,
	}
	outPad, err = textFilter.Apply(g, "text_1", outPad)
	if err != nil {
		t.Fatalf("DrawTextFilter failed: %v", err)
	}

	// Test ChromaKey
	chroma := effects.ChromaKeyFilter{ColorHex: "0x00FF00", Similarity: 0.25}
	outPad, err = chroma.Apply(g, "chroma_1", outPad)
	if err != nil {
		t.Fatalf("ChromaKeyFilter failed: %v", err)
	}

	if outPad.StreamType != filtergraph.StreamTypeVideo {
		t.Fatalf("Expected video stream type, got %s", outPad.StreamType)
	}
}

func TestEffects_Transitions(t *testing.T) {
	g := filtergraph.NewGraph()
	inA := &filtergraph.Pad{ID: "v_a", StreamType: filtergraph.StreamTypeVideo}
	inB := &filtergraph.Pad{ID: "v_b", StreamType: filtergraph.StreamTypeVideo}

	xfade := effects.XFadeFilter{
		Transition: timeline.TransitionDissolve,
		Duration:   1.5,
		Offset:     4.0,
	}

	outPad, err := xfade.Apply(g, "xfade_1", inA, inB)
	if err != nil {
		t.Fatalf("XFadeFilter failed: %v", err)
	}

	node, _ := g.GetNode("xfade_1")
	if node.FilterName != "xfade" {
		t.Fatalf("Expected xfade filter name, got %s", node.FilterName)
	}
	_ = outPad
}
