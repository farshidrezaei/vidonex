package compiler

import (
	"fmt"

	"github.com/farshidrezaei/vidonyx/filtergraph"
	"github.com/farshidrezaei/vidonyx/types"
)

// InjectVideoNormalizer creates a normalization chain ensuring uniform resolution, FPS, SAR and pixel format:
// [in] -> format=yuva420p -> scale -> pad(color=black@0.0) -> setsar -> fps -> format=yuva420p -> [out]
func InjectVideoNormalizer(graph *filtergraph.Graph, inputPad *filtergraph.Pad, canvas types.Size, frameRate types.Rational) (*filtergraph.Pad, error) {
	normalizerID := graph.NextPadID("video_norm")

	// Ensure alpha-capable pixel format early so transparent PNGs and padding retain 0 alpha
	preFormatNode := graph.NewNode(fmt.Sprintf("%s_preformat", normalizerID), "format")
	preFormatNode.SetParam("pix_fmts", "yuva420p")
	preFormatInput := preFormatNode.AddInput(inputPad.ID, filtergraph.StreamTypeVideo)
	preFormatOutput := preFormatNode.AddOutput(graph.NextPadID("preformat_out"), filtergraph.StreamTypeVideo)
	if err := graph.Connect(inputPad, preFormatInput); err != nil {
		return nil, err
	}

	// Scale & Pad to canvas with letterboxing/pillarboxing
	scaleNode := graph.NewNode(fmt.Sprintf("%s_scale", normalizerID), "scale")
	scaleNode.SetParam("w", canvas.Width)
	scaleNode.SetParam("h", canvas.Height)
	scaleNode.SetParam("force_original_aspect_ratio", "decrease")

	scaleInput := scaleNode.AddInput(preFormatOutput.ID, filtergraph.StreamTypeVideo)
	scaleOutput := scaleNode.AddOutput(graph.NextPadID("scale_out"), filtergraph.StreamTypeVideo)
	if err := graph.Connect(preFormatOutput, scaleInput); err != nil {
		return nil, err
	}

	// Pad node with 100% transparent padding (black@0.0)
	padNode := graph.NewNode(fmt.Sprintf("%s_pad", normalizerID), "pad")
	padNode.SetParam("w", canvas.Width)
	padNode.SetParam("h", canvas.Height)
	padNode.SetParam("x", "(ow-iw)/2")
	padNode.SetParam("y", "(oh-ih)/2")
	padNode.SetParam("color", "black@0.0")

	padInput := padNode.AddInput(scaleOutput.ID, filtergraph.StreamTypeVideo)
	padOutput := padNode.AddOutput(graph.NextPadID("pad_out"), filtergraph.StreamTypeVideo)
	if err := graph.Connect(scaleOutput, padInput); err != nil {
		return nil, err
	}

	// SetSAR node
	sarNode := graph.NewNode(fmt.Sprintf("%s_sar", normalizerID), "setsar")
	sarNode.SetParam("sar", "1")
	sarInput := sarNode.AddInput(padOutput.ID, filtergraph.StreamTypeVideo)
	sarOutput := sarNode.AddOutput(graph.NextPadID("sar_out"), filtergraph.StreamTypeVideo)
	if err := graph.Connect(padOutput, sarInput); err != nil {
		return nil, err
	}

	// FPS node
	fpsNode := graph.NewNode(fmt.Sprintf("%s_fps", normalizerID), "fps")
	fpsNode.SetParam("fps", frameRate.FFmpegString())
	fpsInput := fpsNode.AddInput(sarOutput.ID, filtergraph.StreamTypeVideo)
	fpsOutput := fpsNode.AddOutput(graph.NextPadID("fps_out"), filtergraph.StreamTypeVideo)
	if err := graph.Connect(sarOutput, fpsInput); err != nil {
		return nil, err
	}

	// Format node: preserve alpha transparency (yuva420p) for compositing
	formatNode := graph.NewNode(fmt.Sprintf("%s_format", normalizerID), "format")
	formatNode.SetParam("pix_fmts", "yuva420p")
	formatInput := formatNode.AddInput(fpsOutput.ID, filtergraph.StreamTypeVideo)
	finalOutput := formatNode.AddOutput(normalizerID, filtergraph.StreamTypeVideo)
	if err := graph.Connect(fpsOutput, formatInput); err != nil {
		return nil, err
	}

	return finalOutput, nil
}

// InjectAudioNormalizer creates a normalization chain ensuring 48kHz stereo FLTP audio:
// [in] -> aresample=48000 -> aformat=sample_fmts=fltp:channel_layouts=stereo -> [out]
func InjectAudioNormalizer(graph *filtergraph.Graph, inputPad *filtergraph.Pad) (*filtergraph.Pad, error) {
	normalizerID := graph.NextPadID("audio_norm")

	resampleNode := graph.NewNode(fmt.Sprintf("%s_resample", normalizerID), "aresample")
	resampleNode.SetParam("sample_rate", 48000)
	resampleInput := resampleNode.AddInput(inputPad.ID, filtergraph.StreamTypeAudio)
	resampleOutput := resampleNode.AddOutput(graph.NextPadID("resample_out"), filtergraph.StreamTypeAudio)
	if err := graph.Connect(inputPad, resampleInput); err != nil {
		return nil, err
	}

	formatNode := graph.NewNode(fmt.Sprintf("%s_aformat", normalizerID), "aformat")
	formatNode.SetParam("sample_fmts", "fltp")
	formatNode.SetParam("channel_layouts", "stereo")
	formatInput := formatNode.AddInput(resampleOutput.ID, filtergraph.StreamTypeAudio)
	finalOutput := formatNode.AddOutput(normalizerID, filtergraph.StreamTypeAudio)
	if err := graph.Connect(resampleOutput, formatInput); err != nil {
		return nil, err
	}

	return finalOutput, nil
}
