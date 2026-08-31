package compiler

import (
	"fmt"
	"sort"

	"github.com/farshidrezaei/vidonyx/filtergraph"
	"github.com/farshidrezaei/vidonyx/timeline"
)

// ProcessClipVideo processes a single clip's video pipeline:
// Trim -> Speed (setpts) -> Opacity -> Custom Effects -> Normalized Stream
func ProcessClipVideo(g *filtergraph.Graph, rawInPad *filtergraph.Pad, clip *timeline.Clip, tl *timeline.Timeline) (*filtergraph.Pad, error) {
	currPad := rawInPad

	// 1. Trim & SetPTS
	trimNode := g.NewNode(fmt.Sprintf("trim_v_%s", clip.ID), "trim")
	startSec := clip.SourceStart.Seconds()
	durSec := clip.Duration.Seconds()
	trimNode.SetParam("start", fmt.Sprintf("%.4f", startSec))
	trimNode.SetParam("duration", fmt.Sprintf("%.4f", durSec))

	trimIn := trimNode.AddInput(currPad.ID, filtergraph.StreamTypeVideo)
	trimOut := trimNode.AddOutput(g.NextPadID("trim_out_v"), filtergraph.StreamTypeVideo)
	if err := g.Connect(currPad, trimIn); err != nil {
		return nil, err
	}
	currPad = trimOut

	// SetPTS based on speed and start offset
	ptsNode := g.NewNode(fmt.Sprintf("pts_v_%s", clip.ID), "setpts")
	if clip.Speed > 0 && clip.Speed != 1.0 {
		ptsNode.SetParam("expr", fmt.Sprintf("(PTS-STARTPTS)/%.4f", clip.Speed))
	} else {
		ptsNode.SetParam("expr", "PTS-STARTPTS")
	}
	ptsIn := ptsNode.AddInput(currPad.ID, filtergraph.StreamTypeVideo)
	ptsOut := ptsNode.AddOutput(g.NextPadID("pts_out_v"), filtergraph.StreamTypeVideo)
	if err := g.Connect(currPad, ptsIn); err != nil {
		return nil, err
	}
	currPad = ptsOut

	// 2. Opacity
	if clip.Opacity < 1.0 && clip.Opacity >= 0.0 {
		fmtYuva := g.NewNode(fmt.Sprintf("fmt_yuva_%s", clip.ID), "format")
		fmtYuva.SetParam("pix_fmts", "yuva420p")
		yuvaIn := fmtYuva.AddInput(currPad.ID, filtergraph.StreamTypeVideo)
		yuvaOut := fmtYuva.AddOutput(g.NextPadID("yuva_out"), filtergraph.StreamTypeVideo)
		if err := g.Connect(currPad, yuvaIn); err != nil {
			return nil, err
		}

		mixNode := g.NewNode(fmt.Sprintf("opacity_%s", clip.ID), "colorchannelmixer")
		mixNode.SetParam("aa", fmt.Sprintf("%.2f", clip.Opacity))
		mixIn := mixNode.AddInput(yuvaOut.ID, filtergraph.StreamTypeVideo)
		mixOut := mixNode.AddOutput(g.NextPadID("opac_out"), filtergraph.StreamTypeVideo)
		if err := g.Connect(yuvaOut, mixIn); err != nil {
			return nil, err
		}
		currPad = mixOut
	}

	// 3. Clip scale & position handling
	if clip.Scale != 1.0 && clip.Scale > 0 {
		scaleNode := g.NewNode(fmt.Sprintf("scale_%s", clip.ID), "scale")
		scaleNode.SetParam("w", fmt.Sprintf("iw*%.4f", clip.Scale))
		scaleNode.SetParam("h", fmt.Sprintf("ih*%.4f", clip.Scale))
		sIn := scaleNode.AddInput(currPad.ID, filtergraph.StreamTypeVideo)
		sOut := scaleNode.AddOutput(g.NextPadID("sc_out"), filtergraph.StreamTypeVideo)
		if err := g.Connect(currPad, sIn); err != nil {
			return nil, err
		}
		currPad = sOut
	}

	return currPad, nil
}

// ProcessClipAudio processes a single clip's audio pipeline:
// Trim -> atempo (speed) -> Volume -> Delay (adelay to TimelineStart)
func ProcessClipAudio(g *filtergraph.Graph, rawInPad *filtergraph.Pad, clip *timeline.Clip) (*filtergraph.Pad, error) {
	currPad := rawInPad

	// 1. atrim & asetpts
	atrimNode := g.NewNode(fmt.Sprintf("atrim_%s", clip.ID), "atrim")
	startSec := clip.SourceStart.Seconds()
	durSec := clip.Duration.Seconds()
	atrimNode.SetParam("start", fmt.Sprintf("%.4f", startSec))
	atrimNode.SetParam("duration", fmt.Sprintf("%.4f", durSec))

	atrimIn := atrimNode.AddInput(currPad.ID, filtergraph.StreamTypeAudio)
	atrimOut := atrimNode.AddOutput(g.NextPadID("atrim_out"), filtergraph.StreamTypeAudio)
	if err := g.Connect(currPad, atrimIn); err != nil {
		return nil, err
	}
	currPad = atrimOut

	// asetpts
	aptsNode := g.NewNode(fmt.Sprintf("apts_%s", clip.ID), "asetpts")
	aptsNode.SetParam("expr", "PTS-STARTPTS")
	aptsIn := aptsNode.AddInput(currPad.ID, filtergraph.StreamTypeAudio)
	aptsOut := aptsNode.AddOutput(g.NextPadID("apts_out"), filtergraph.StreamTypeAudio)
	if err := g.Connect(currPad, aptsIn); err != nil {
		return nil, err
	}
	currPad = aptsOut

	// 2. Volume
	if clip.Volume != 1.0 && clip.Volume >= 0.0 {
		volNode := g.NewNode(fmt.Sprintf("vol_%s", clip.ID), "volume")
		volNode.SetParam("volume", fmt.Sprintf("%.2f", clip.Volume))
		volIn := volNode.AddInput(currPad.ID, filtergraph.StreamTypeAudio)
		volOut := volNode.AddOutput(g.NextPadID("vol_out"), filtergraph.StreamTypeAudio)
		if err := g.Connect(currPad, volIn); err != nil {
			return nil, err
		}
		currPad = volOut
	}

	// 3. adelay (to position audio on the timeline)
	if clip.TimelineStart > 0 {
		delayMs := clip.TimelineStart.Milliseconds()
		delayNode := g.NewNode(fmt.Sprintf("delay_%s", clip.ID), "adelay")
		delayNode.SetParam("delays", fmt.Sprintf("%d|%d", delayMs, delayMs))
		dIn := delayNode.AddInput(currPad.ID, filtergraph.StreamTypeAudio)
		dOut := delayNode.AddOutput(g.NextPadID("delay_out"), filtergraph.StreamTypeAudio)
		if err := g.Connect(currPad, dIn); err != nil {
			return nil, err
		}
		currPad = dOut
	}

	return currPad, nil
}

// BuildVideoCompositor builds the background canvas and overlays all video tracks by Z-index.
func BuildVideoCompositor(g *filtergraph.Graph, tl *timeline.Timeline, processedVideoPads []*ClipVideoPad) (*filtergraph.Pad, error) {
	// 1. Generate base color background canvas
	bgNode := g.NewNode("bg_canvas", "color")
	bgNode.SetParam("c", tl.BackgroundColor.FFmpegColor())
	bgNode.SetParam("s", tl.Canvas.String())
	bgNode.SetParam("r", tl.FPS.FFmpegString())
	bgNode.SetParam("d", fmt.Sprintf("%.4f", tl.Duration().Seconds()))
	bgOut := bgNode.AddOutput(g.NextPadID("base_canvas"), filtergraph.StreamTypeVideo)

	currCanvas := bgOut

	// Sort clips by Track Z-Index and TimelineStart
	sort.SliceStable(processedVideoPads, func(i, j int) bool {
		if processedVideoPads[i].ZIndex != processedVideoPads[j].ZIndex {
			return processedVideoPads[i].ZIndex < processedVideoPads[j].ZIndex
		}
		return processedVideoPads[i].Clip.TimelineStart < processedVideoPads[j].Clip.TimelineStart
	})

	// Overlay each clip on top of the canvas
	for i, item := range processedVideoPads {
		overlayNode := g.NewNode(fmt.Sprintf("overlay_%d_%s", i, item.Clip.ID), "overlay")
		
		// Position coordinates
		xExpr := fmt.Sprintf("%d", item.Clip.Position.X)
		yExpr := fmt.Sprintf("%d", item.Clip.Position.Y)
		overlayNode.SetParam("x", xExpr)
		overlayNode.SetParam("y", yExpr)
		
		// Time interval enable expression
		startSec := item.Clip.TimelineStart.Seconds()
		endSec := item.Clip.TimelineEnd().Seconds()
		overlayNode.SetParam("enable", fmt.Sprintf("between(t,%.4f,%.4f)", startSec, endSec))
		overlayNode.SetParam("eof_action", "pass")

		inBase := overlayNode.AddInput(currCanvas.ID, filtergraph.StreamTypeVideo)
		inOverlay := overlayNode.AddInput(item.Pad.ID, filtergraph.StreamTypeVideo)
		outCanvas := overlayNode.AddOutput(g.NextPadID("comp_v"), filtergraph.StreamTypeVideo)

		if err := g.Connect(currCanvas, inBase); err != nil {
			return nil, err
		}
		if err := g.Connect(item.Pad, inOverlay); err != nil {
			return nil, err
		}

		currCanvas = outCanvas
	}

	return currCanvas, nil
}

// BuildAudioMixer mixes all processed audio streams into a single stereo output.
func BuildAudioMixer(g *filtergraph.Graph, tl *timeline.Timeline, audioPads []*filtergraph.Pad) (*filtergraph.Pad, error) {
	if len(audioPads) == 0 {
		// Generate silent audio stream matching timeline duration
		silentNode := g.NewNode("silent_audio", "anullsrc")
		silentNode.SetParam("r", "48000")
		silentNode.SetParam("cl", "stereo")
		silentNode.SetParam("d", fmt.Sprintf("%.4f", tl.Duration().Seconds()))
		return silentNode.AddOutput("out_a", filtergraph.StreamTypeAudio), nil
	}

	if len(audioPads) == 1 {
		return audioPads[0], nil
	}

	// Mix multiple audio streams with amix
	amixNode := g.NewNode("audio_mixer", "amix")
	amixNode.SetParam("inputs", len(audioPads))
	amixNode.SetParam("duration", "longest")
	amixNode.SetParam("dropout_transition", "0")

	for _, pad := range audioPads {
		inPad := amixNode.AddInput(pad.ID, filtergraph.StreamTypeAudio)
		if err := g.Connect(pad, inPad); err != nil {
			return nil, err
		}
	}

	amixOut := amixNode.AddOutput("mixed_a", filtergraph.StreamTypeAudio)
	return amixOut, nil
}

// ClipVideoPad pairs a processed video pad with its Clip metadata and Track Z-Index.
type ClipVideoPad struct {
	Clip   *timeline.Clip
	ZIndex int
	Pad    *filtergraph.Pad
}
