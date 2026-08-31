package subtitles

import (
	"fmt"
	"strings"

	"github.com/farshidrezaei/vidonyx/filtergraph"
	"github.com/farshidrezaei/vidonyx/types"
)

// AttachSubtitles converts all subtitle cues on a track into chained drawtext filter nodes.
func AttachSubtitles(graph *filtergraph.Graph, inputVideoPad *filtergraph.Pad, subtitleTrack *SubtitleTrack, canvasSize types.Size) (*filtergraph.Pad, error) {
	if subtitleTrack == nil || len(subtitleTrack.Cues) == 0 {
		return inputVideoPad, nil
	}

	currentPad := inputVideoPad
	style := subtitleTrack.Style

	for index, cue := range subtitleTrack.Cues {
		drawTextNode := graph.NewNode(fmt.Sprintf("sub_drawtext_%d", index), "drawtext")

		// Text content (escaping single quotes for FFmpeg)
		escapedText := escapeDrawText(cue.Text)
		drawTextNode.SetParam("text", fmt.Sprintf("'%s'", escapedText))

		// Font styling
		if style.FontFile != "" {
			drawTextNode.SetParam("fontfile", fmt.Sprintf("'%s'", style.FontFile))
		}
		if style.FontSize > 0 {
			drawTextNode.SetParam("fontsize", style.FontSize)
		} else {
			drawTextNode.SetParam("fontsize", 36)
		}

		drawTextNode.SetParam("fontcolor", style.PrimaryColor.FFmpegColor())

		// Outline / Border
		if style.OutlineWidth > 0 {
			drawTextNode.SetParam("borderw", style.OutlineWidth)
			drawTextNode.SetParam("bordercolor", style.OutlineColor.FFmpegColor())
		}

		// Background box
		if style.Box {
			drawTextNode.SetParam("box", "1")
			drawTextNode.SetParam("boxcolor", style.BoxColor.FFmpegColor())
			drawTextNode.SetParam("boxborderw", style.BoxBorderWidth)
		}

		// Positioning (X, Y)
		xExpr, yExpr := calculateSubtitleCoordinates(style, canvasSize)
		drawTextNode.SetParam("x", xExpr)
		drawTextNode.SetParam("y", yExpr)

		// Display time interval
		startSec := cue.StartTime.Seconds()
		endSec := cue.EndTime.Seconds()
		drawTextNode.SetParam("enable", fmt.Sprintf("between(t,%.4f,%.4f)", startSec, endSec))

		inputPad := drawTextNode.AddInput(currentPad.ID, filtergraph.StreamTypeVideo)
		outputPad := drawTextNode.AddOutput(graph.NextPadID("sub_out"), filtergraph.StreamTypeVideo)

		if err := graph.Connect(currentPad, inputPad); err != nil {
			return nil, fmt.Errorf("subtitles: failed connecting subtitle cue %d: %w", index, err)
		}

		currentPad = outputPad
	}

	return currentPad, nil
}

func calculateSubtitleCoordinates(style SubtitleStyle, _ types.Size) (string, string) {
	var xExpr, yExpr string

	switch style.Alignment {
	case types.AlignTopLeft, types.AlignCenterLeft, types.AlignBottomLeft:
		xExpr = "40"
	case types.AlignTopRight, types.AlignCenterRight, types.AlignBottomRight:
		xExpr = "w-tw-40"
	default: // Centered
		xExpr = "(w-tw)/2"
	}

	switch style.Alignment {
	case types.AlignTopLeft, types.AlignTopCenter, types.AlignTopRight:
		margin := style.MarginTop
		if margin <= 0 {
			margin = 60
		}
		yExpr = fmt.Sprintf("%d", margin)
	case types.AlignCenterLeft, types.AlignCenter, types.AlignCenterRight:
		yExpr = "(h-th)/2"
	default: // Bottom aligned
		margin := style.MarginBottom
		if margin <= 0 {
			margin = 60
		}
		yExpr = fmt.Sprintf("h-th-%d", margin)
	}

	return xExpr, yExpr
}

func escapeDrawText(text string) string {
	text = strings.ReplaceAll(text, "\\", "\\\\")
	text = strings.ReplaceAll(text, "'", "\\'")
	text = strings.ReplaceAll(text, ":", "\\:")
	text = strings.ReplaceAll(text, "%", "\\%")
	return text
}
