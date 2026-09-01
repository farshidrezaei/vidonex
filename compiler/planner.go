package compiler

import (
	"fmt"
	"sort"

	"github.com/farshidrezaei/vidonyx/chromakey"
	"github.com/farshidrezaei/vidonyx/effects"
	"github.com/farshidrezaei/vidonyx/filtergraph"
	"github.com/farshidrezaei/vidonyx/timeline"
)

// ProcessClipVideo processes a single clip's video pipeline:
// Trim -> Speed (setpts) -> ChromaKey -> Fade In/Out -> Opacity / Animated Opacity -> Animated Scale (zoompan) -> Normalized Stream
func ProcessClipVideo(graph *filtergraph.Graph, rawInputPad *filtergraph.Pad, clip *timeline.Clip, compositionTimeline *timeline.Timeline) (*filtergraph.Pad, error) {
	currentPad := rawInputPad

	// 1. Trim & SetPTS
	trimNode := graph.NewNode(fmt.Sprintf("trim_video_%s", clip.ID), "trim")
	startSeconds := clip.SourceStart.Seconds()
	durationSeconds := clip.Duration.Seconds()
	trimNode.SetParam("start", fmt.Sprintf("%.4f", startSeconds))
	trimNode.SetParam("duration", fmt.Sprintf("%.4f", durationSeconds))

	trimInput := trimNode.AddInput(currentPad.ID, filtergraph.StreamTypeVideo)
	trimOutput := trimNode.AddOutput(graph.NextPadID("trim_out_video"), filtergraph.StreamTypeVideo)
	if err := graph.Connect(currentPad, trimInput); err != nil {
		return nil, err
	}
	currentPad = trimOutput

	// SetPTS based on speed and start offset (accurately shifts frames to TimelineStart to handle gaps and start offsets)
	ptsNode := graph.NewNode(fmt.Sprintf("pts_video_%s", clip.ID), "setpts")
	timelineStartSeconds := clip.TimelineStart.Seconds()

	if clip.Speed > 0 && clip.Speed != 1.0 {
		if timelineStartSeconds > 0 {
			ptsNode.SetParam("expr", fmt.Sprintf("((PTS-STARTPTS)/%.4f)+%.4f/TB", clip.Speed, timelineStartSeconds))
		} else {
			ptsNode.SetParam("expr", fmt.Sprintf("(PTS-STARTPTS)/%.4f", clip.Speed))
		}
	} else {
		if timelineStartSeconds > 0 {
			ptsNode.SetParam("expr", fmt.Sprintf("PTS-STARTPTS+%.4f/TB", timelineStartSeconds))
		} else {
			ptsNode.SetParam("expr", "PTS-STARTPTS")
		}
	}
	ptsInput := ptsNode.AddInput(currentPad.ID, filtergraph.StreamTypeVideo)
	ptsOutput := ptsNode.AddOutput(graph.NextPadID("pts_out_video"), filtergraph.StreamTypeVideo)
	if err := graph.Connect(currentPad, ptsInput); err != nil {
		return nil, err
	}
	currentPad = ptsOutput

	// 2. Chroma Keying (Green/Blue screen removal + Despill)
	if clip.ChromaKeyOptions != nil {
		chromaPad, err := chromakey.ApplyChromaKey(graph, fmt.Sprintf("chroma_%s", clip.ID), currentPad, *clip.ChromaKeyOptions)
		if err != nil {
			return nil, fmt.Errorf("failed applying chromakey to clip %q: %w", clip.ID, err)
		}
		currentPad = chromaPad
	}

	// 3. Fade In & Fade Out
	if clip.FadeInDuration > 0 {
		fadeInNode := graph.NewNode(fmt.Sprintf("fade_in_%s", clip.ID), "fade")
		fadeInNode.SetParam("t", "in")
		fadeInNode.SetParam("st", "0")
		fadeInNode.SetParam("d", fmt.Sprintf("%.4f", clip.FadeInDuration.Seconds()))
		fadeInInput := fadeInNode.AddInput(currentPad.ID, filtergraph.StreamTypeVideo)
		fadeInOutput := fadeInNode.AddOutput(graph.NextPadID("fade_in_out"), filtergraph.StreamTypeVideo)
		if err := graph.Connect(currentPad, fadeInInput); err != nil {
			return nil, err
		}
		currentPad = fadeInOutput
	}

	if clip.FadeOutDuration > 0 {
		fadeOutStart := clip.Duration.Seconds() - clip.FadeOutDuration.Seconds()
		if fadeOutStart < 0 {
			fadeOutStart = 0
		}
		fadeOutNode := graph.NewNode(fmt.Sprintf("fade_out_%s", clip.ID), "fade")
		fadeOutNode.SetParam("t", "out")
		fadeOutNode.SetParam("st", fmt.Sprintf("%.4f", fadeOutStart))
		fadeOutNode.SetParam("d", fmt.Sprintf("%.4f", clip.FadeOutDuration.Seconds()))
		fadeOutInput := fadeOutNode.AddInput(currentPad.ID, filtergraph.StreamTypeVideo)
		fadeOutOutput := fadeOutNode.AddOutput(graph.NextPadID("fade_out_out"), filtergraph.StreamTypeVideo)
		if err := graph.Connect(currentPad, fadeOutInput); err != nil {
			return nil, err
		}
		currentPad = fadeOutOutput
	}

	// 4. Opacity (Static or Animated Keyframe Track)
	if clip.OpacityTrack != nil || (clip.Opacity < 1.0 && clip.Opacity >= 0.0) {
		formatYuva := graph.NewNode(fmt.Sprintf("format_yuva_%s", clip.ID), "format")
		formatYuva.SetParam("pix_fmts", "yuva420p")
		yuvaInput := formatYuva.AddInput(currentPad.ID, filtergraph.StreamTypeVideo)
		yuvaOutput := formatYuva.AddOutput(graph.NextPadID("yuva_out"), filtergraph.StreamTypeVideo)
		if err := graph.Connect(currentPad, yuvaInput); err != nil {
			return nil, err
		}

		mixerNode := graph.NewNode(fmt.Sprintf("opacity_%s", clip.ID), "colorchannelmixer")
		if clip.OpacityTrack != nil {
			mixerNode.SetParam("aa", fmt.Sprintf("'%s'", clip.OpacityTrack.ToFFmpegExpression()))
		} else {
			mixerNode.SetParam("aa", fmt.Sprintf("%.2f", clip.Opacity))
		}

		mixerInput := mixerNode.AddInput(yuvaOutput.ID, filtergraph.StreamTypeVideo)
		mixerOutput := mixerNode.AddOutput(graph.NextPadID("opacity_out"), filtergraph.StreamTypeVideo)
		if err := graph.Connect(yuvaOutput, mixerInput); err != nil {
			return nil, err
		}
		currentPad = mixerOutput
	}

	// 5. Clip scale handling (Static or Animated Keyframe Track)
	// Uses zoompan filter instead of scale with eval=frame because dynamically changing
	// output dimensions per frame causes FFmpeg to hang or deadlock when combined with
	// downstream overlay filters. zoompan maintains constant output dimensions.
	if clip.ScaleTrack != nil {
		zoompanNode := graph.NewNode(fmt.Sprintf("zoompan_%s", clip.ID), "zoompan")
		scaleExpression := clip.ScaleTrack.ToFFmpegExpressionWithVariable("in_time")
		zoompanNode.SetParam("z", fmt.Sprintf("'%s'", scaleExpression))
		zoompanNode.SetParam("x", "'iw/2-(iw/zoom/2)'")
		zoompanNode.SetParam("y", "'ih/2-(ih/zoom/2)'")
		zoompanNode.SetParam("d", "1")
		zoompanNode.SetParam("s", fmt.Sprintf("%dx%d", compositionTimeline.Canvas.Width, compositionTimeline.Canvas.Height))
		zoompanNode.SetParam("fps", compositionTimeline.FPS.FFmpegString())
		zoompanInput := zoompanNode.AddInput(currentPad.ID, filtergraph.StreamTypeVideo)
		zoompanOutput := zoompanNode.AddOutput(graph.NextPadID("zoompan_out"), filtergraph.StreamTypeVideo)
		if err := graph.Connect(currentPad, zoompanInput); err != nil {
			return nil, err
		}

		// Ensure monotonic PTS after zoompan so overlay filter syncs instantaneously without blocking
		ptsZoomNode := graph.NewNode(fmt.Sprintf("pts_zoom_%s", clip.ID), "setpts")
		ptsZoomNode.SetParam("expr", "N/FRAME_RATE/TB")
		ptsZoomInput := ptsZoomNode.AddInput(zoompanOutput.ID, filtergraph.StreamTypeVideo)
		ptsZoomOutput := ptsZoomNode.AddOutput(graph.NextPadID("pts_zoom_out"), filtergraph.StreamTypeVideo)
		if err := graph.Connect(zoompanOutput, ptsZoomInput); err != nil {
			return nil, err
		}

		currentPad = ptsZoomOutput
	} else if clip.Scale != 1.0 && clip.Scale > 0 {
		scaleNode := graph.NewNode(fmt.Sprintf("scale_%s", clip.ID), "scale")
		scaleNode.SetParam("w", fmt.Sprintf("iw*%.4f", clip.Scale))
		scaleNode.SetParam("h", fmt.Sprintf("ih*%.4f", clip.Scale))
		scaleInput := scaleNode.AddInput(currentPad.ID, filtergraph.StreamTypeVideo)
		scaleOutput := scaleNode.AddOutput(graph.NextPadID("scale_out"), filtergraph.StreamTypeVideo)
		if err := graph.Connect(currentPad, scaleInput); err != nil {
			return nil, err
		}
		currentPad = scaleOutput
	}

	return currentPad, nil
}

// ProcessClipAudio processes a single clip's audio pipeline:
// Trim -> atempo (speed) -> afade (in/out) -> Volume -> Delay (adelay to TimelineStart)
func ProcessClipAudio(graph *filtergraph.Graph, rawInputPad *filtergraph.Pad, clip *timeline.Clip) (*filtergraph.Pad, error) {
	currentPad := rawInputPad

	// 1. atrim & asetpts
	atrimNode := graph.NewNode(fmt.Sprintf("atrim_%s", clip.ID), "atrim")
	startSeconds := clip.SourceStart.Seconds()
	durationSeconds := clip.Duration.Seconds()
	atrimNode.SetParam("start", fmt.Sprintf("%.4f", startSeconds))
	atrimNode.SetParam("duration", fmt.Sprintf("%.4f", durationSeconds))

	atrimInput := atrimNode.AddInput(currentPad.ID, filtergraph.StreamTypeAudio)
	atrimOutput := atrimNode.AddOutput(graph.NextPadID("atrim_out"), filtergraph.StreamTypeAudio)
	if err := graph.Connect(currentPad, atrimInput); err != nil {
		return nil, err
	}
	currentPad = atrimOutput

	// asetpts
	aptsNode := graph.NewNode(fmt.Sprintf("apts_%s", clip.ID), "asetpts")
	aptsNode.SetParam("expr", "PTS-STARTPTS")
	aptsInput := aptsNode.AddInput(currentPad.ID, filtergraph.StreamTypeAudio)
	aptsOutput := aptsNode.AddOutput(graph.NextPadID("apts_out"), filtergraph.StreamTypeAudio)
	if err := graph.Connect(currentPad, aptsInput); err != nil {
		return nil, err
	}
	currentPad = aptsOutput

	// 2. Audio Fade In & Fade Out
	if clip.FadeInDuration > 0 {
		afadeInNode := graph.NewNode(fmt.Sprintf("afade_in_%s", clip.ID), "afade")
		afadeInNode.SetParam("t", "in")
		afadeInNode.SetParam("st", "0")
		afadeInNode.SetParam("d", fmt.Sprintf("%.4f", clip.FadeInDuration.Seconds()))
		afadeInInput := afadeInNode.AddInput(currentPad.ID, filtergraph.StreamTypeAudio)
		afadeInOutput := afadeInNode.AddOutput(graph.NextPadID("afade_in_out"), filtergraph.StreamTypeAudio)
		if err := graph.Connect(currentPad, afadeInInput); err != nil {
			return nil, err
		}
		currentPad = afadeInOutput
	}

	if clip.FadeOutDuration > 0 {
		fadeOutStart := clip.Duration.Seconds() - clip.FadeOutDuration.Seconds()
		if fadeOutStart < 0 {
			fadeOutStart = 0
		}
		afadeOutNode := graph.NewNode(fmt.Sprintf("afade_out_%s", clip.ID), "afade")
		afadeOutNode.SetParam("t", "out")
		afadeOutNode.SetParam("st", fmt.Sprintf("%.4f", fadeOutStart))
		afadeOutNode.SetParam("d", fmt.Sprintf("%.4f", clip.FadeOutDuration.Seconds()))
		afadeOutInput := afadeOutNode.AddInput(currentPad.ID, filtergraph.StreamTypeAudio)
		afadeOutOutput := afadeOutNode.AddOutput(graph.NextPadID("afade_out_out"), filtergraph.StreamTypeAudio)
		if err := graph.Connect(currentPad, afadeOutInput); err != nil {
			return nil, err
		}
		currentPad = afadeOutOutput
	}

	// 3. Volume
	if clip.Volume != 1.0 && clip.Volume >= 0.0 {
		volumeNode := graph.NewNode(fmt.Sprintf("volume_%s", clip.ID), "volume")
		volumeNode.SetParam("volume", fmt.Sprintf("%.2f", clip.Volume))
		volumeInput := volumeNode.AddInput(currentPad.ID, filtergraph.StreamTypeAudio)
		volumeOutput := volumeNode.AddOutput(graph.NextPadID("volume_out"), filtergraph.StreamTypeAudio)
		if err := graph.Connect(currentPad, volumeInput); err != nil {
			return nil, err
		}
		currentPad = volumeOutput
	}

	// 4. adelay (to position audio on the timeline)
	if clip.TimelineStart > 0 {
		delayMilliseconds := clip.TimelineStart.Milliseconds()
		delayNode := graph.NewNode(fmt.Sprintf("delay_%s", clip.ID), "adelay")
		delayNode.SetParam("delays", fmt.Sprintf("%d|%d", delayMilliseconds, delayMilliseconds))
		delayInput := delayNode.AddInput(currentPad.ID, filtergraph.StreamTypeAudio)
		delayOutput := delayNode.AddOutput(graph.NextPadID("delay_out"), filtergraph.StreamTypeAudio)
		if err := graph.Connect(currentPad, delayInput); err != nil {
			return nil, err
		}
		currentPad = delayOutput
	}

	return currentPad, nil
}

// ChainTrackTransitionsVideo connects adjacent video clips on a track using xfade transitions.
func ChainTrackTransitionsVideo(graph *filtergraph.Graph, videoPads []*filtergraph.Pad, transitions []*timeline.Transition, clips []*timeline.Clip) (*filtergraph.Pad, error) {
	if len(videoPads) == 0 {
		return nil, nil
	}
	if len(videoPads) == 1 || len(transitions) == 0 {
		return videoPads[0], nil
	}

	currentStream := videoPads[0]
	accumulatedDuration := clips[0].Duration.Seconds()

	for index := 0; index < len(transitions) && index+1 < len(videoPads); index++ {
		transition := transitions[index]
		nextPad := videoPads[index+1]

		transitionDuration := transition.Duration.Seconds()
		transitionOffset := accumulatedDuration - transitionDuration
		if transitionOffset < 0 {
			transitionOffset = 0
		}

		xfadeFilter := effects.XFadeFilter{
			Transition: transition.Type,
			Duration:   transitionDuration,
			Offset:     transitionOffset,
		}

		transitionNodeID := graph.NextPadID(fmt.Sprintf("xfade_trans_%d", index))
		transitionedOutput, err := xfadeFilter.Apply(graph, transitionNodeID, currentStream, nextPad)
		if err != nil {
			return nil, fmt.Errorf("compiler: failed to apply xfade transition %q: %w", transition.ID, err)
		}

		currentStream = transitionedOutput
		accumulatedDuration += clips[index+1].Duration.Seconds() - transitionDuration
	}

	return currentStream, nil
}

// ChainTrackTransitionsAudio connects adjacent audio clips on a track using acrossfade transitions.
func ChainTrackTransitionsAudio(graph *filtergraph.Graph, audioPads []*filtergraph.Pad, transitions []*timeline.Transition) (*filtergraph.Pad, error) {
	if len(audioPads) == 0 {
		return nil, nil
	}
	if len(audioPads) == 1 || len(transitions) == 0 {
		return audioPads[0], nil
	}

	currentStream := audioPads[0]

	for index := 0; index < len(transitions) && index+1 < len(audioPads); index++ {
		transition := transitions[index]
		nextPad := audioPads[index+1]

		acrossFadeFilter := effects.AcrossFadeFilter{
			Duration: transition.Duration.Seconds(),
			Curve1:   "tri",
			Curve2:   "tri",
		}

		transitionNodeID := graph.NextPadID(fmt.Sprintf("acrossfade_trans_%d", index))
		transitionedOutput, err := acrossFadeFilter.Apply(graph, transitionNodeID, currentStream, nextPad)
		if err != nil {
			return nil, fmt.Errorf("compiler: failed to apply acrossfade transition %q: %w", transition.ID, err)
		}

		currentStream = transitionedOutput
	}

	return currentStream, nil
}

// BuildVideoCompositor builds the background canvas and overlays all video tracks by Z-index.
func BuildVideoCompositor(graph *filtergraph.Graph, compositionTimeline *timeline.Timeline, processedVideoPads []*ClipVideoPad) (*filtergraph.Pad, error) {
	// 1. Generate base color background canvas
	backgroundNode := graph.NewNode("bg_canvas", "color")
	backgroundNode.SetParam("c", compositionTimeline.BackgroundColor.FFmpegColor())
	backgroundNode.SetParam("s", compositionTimeline.Canvas.String())
	backgroundNode.SetParam("r", compositionTimeline.FPS.FFmpegString())
	backgroundNode.SetParam("d", fmt.Sprintf("%.4f", compositionTimeline.Duration().Seconds()))
	backgroundOutput := backgroundNode.AddOutput(graph.NextPadID("base_canvas"), filtergraph.StreamTypeVideo)

	currentCanvas := backgroundOutput

	// Sort clips by Track Z-Index and TimelineStart
	sort.SliceStable(processedVideoPads, func(i, j int) bool {
		if processedVideoPads[i].ZIndex != processedVideoPads[j].ZIndex {
			return processedVideoPads[i].ZIndex < processedVideoPads[j].ZIndex
		}
		return processedVideoPads[i].Clip.TimelineStart < processedVideoPads[j].Clip.TimelineStart
	})

	// Overlay each clip on top of the canvas
	for index, item := range processedVideoPads {
		overlayNode := graph.NewNode(fmt.Sprintf("overlay_%d_%s", index, item.Clip.ID), "overlay")

		var xExpression, yExpression string
		if item.Clip.PositionTrack != nil {
			posX, posY := item.Clip.PositionTrack.ToFFmpegExpressions()
			xExpression = fmt.Sprintf("'%s'", posX)
			yExpression = fmt.Sprintf("'%s'", posY)
			overlayNode.SetParam("eval", "frame")
		} else {
			xExpression = fmt.Sprintf("%d", item.Clip.Position.X)
			yExpression = fmt.Sprintf("%d", item.Clip.Position.Y)
		}

		overlayNode.SetParam("x", xExpression)
		overlayNode.SetParam("y", yExpression)

		// Time interval enable expression
		startSeconds := item.Clip.TimelineStart.Seconds()
		endSeconds := item.Clip.TimelineEnd().Seconds()
		overlayNode.SetParam("enable", fmt.Sprintf("'between(t,%.4f,%.4f)'", startSeconds, endSeconds))
		overlayNode.SetParam("eof_action", "pass")

		inputBase := overlayNode.AddInput(currentCanvas.ID, filtergraph.StreamTypeVideo)
		inputOverlay := overlayNode.AddInput(item.Pad.ID, filtergraph.StreamTypeVideo)
		outputCanvas := overlayNode.AddOutput(graph.NextPadID("composite_video"), filtergraph.StreamTypeVideo)

		if err := graph.Connect(currentCanvas, inputBase); err != nil {
			return nil, err
		}
		if err := graph.Connect(item.Pad, inputOverlay); err != nil {
			return nil, err
		}

		currentCanvas = outputCanvas
	}

	// Format final output to yuv420p for standard MP4 encoding compatibility
	finalFormatNode := graph.NewNode("final_format_v", "format")
	finalFormatNode.SetParam("pix_fmts", "yuv420p")
	finalFormatInput := finalFormatNode.AddInput(currentCanvas.ID, filtergraph.StreamTypeVideo)
	finalFormatOutput := finalFormatNode.AddOutput(graph.NextPadID("final_canvas_yuv"), filtergraph.StreamTypeVideo)
	if err := graph.Connect(currentCanvas, finalFormatInput); err != nil {
		return nil, err
	}

	return finalFormatOutput, nil
}

// BuildAudioMixer mixes all processed audio streams into a single stereo output.
func BuildAudioMixer(graph *filtergraph.Graph, compositionTimeline *timeline.Timeline, audioPads []*filtergraph.Pad) (*filtergraph.Pad, error) {
	if len(audioPads) == 0 {
		// Generate silent audio stream matching timeline duration
		silentNode := graph.NewNode("silent_audio", "anullsrc")
		silentNode.SetParam("r", "48000")
		silentNode.SetParam("cl", "stereo")
		silentNode.SetParam("d", fmt.Sprintf("%.4f", compositionTimeline.Duration().Seconds()))
		return silentNode.AddOutput("out_a", filtergraph.StreamTypeAudio), nil
	}

	if len(audioPads) == 1 {
		return audioPads[0], nil
	}

	// Mix multiple audio streams with amix
	amixNode := graph.NewNode("audio_mixer", "amix")
	amixNode.SetParam("inputs", len(audioPads))
	amixNode.SetParam("duration", "longest")
	amixNode.SetParam("dropout_transition", "0")

	for _, pad := range audioPads {
		inputPad := amixNode.AddInput(pad.ID, filtergraph.StreamTypeAudio)
		if err := graph.Connect(pad, inputPad); err != nil {
			return nil, err
		}
	}

	amixOutput := amixNode.AddOutput("mixed_audio", filtergraph.StreamTypeAudio)
	return amixOutput, nil
}

// ClipVideoPad pairs a processed video pad with its Clip metadata and Track Z-Index.
type ClipVideoPad struct {
	Clip   *timeline.Clip
	ZIndex int
	Pad    *filtergraph.Pad
}
