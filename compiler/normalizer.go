package compiler

import (
	"fmt"

	"github.com/farshidrezaei/vidonyx/filtergraph"
	"github.com/farshidrezaei/vidonyx/types"
)

// InjectVideoNormalizer creates a normalization chain ensuring uniform resolution, FPS, SAR and pixel format:
// [in] -> scale -> pad -> setsar -> fps -> format=yuv420p -> [out]
func InjectVideoNormalizer(g *filtergraph.Graph, inPad *filtergraph.Pad, canvas types.Size, fps types.Rational) (*filtergraph.Pad, error) {
	normID := g.NextPadID("v_norm")

	// Scale & Pad to canvas with letterboxing/pillarboxing
	scaleNode := g.NewNode(fmt.Sprintf("%s_scale", normID), "scale")
	scaleNode.SetParam("w", canvas.Width)
	scaleNode.SetParam("h", canvas.Height)
	scaleNode.SetParam("force_original_aspect_ratio", "decrease")

	scaleIn := scaleNode.AddInput(inPad.ID, filtergraph.StreamTypeVideo)
	scaleOut := scaleNode.AddOutput(g.NextPadID("sc_out"), filtergraph.StreamTypeVideo)
	if err := g.Connect(inPad, scaleIn); err != nil {
		return nil, err
	}

	// Pad node
	padNode := g.NewNode(fmt.Sprintf("%s_pad", normID), "pad")
	padNode.SetParam("w", canvas.Width)
	padNode.SetParam("h", canvas.Height)
	padNode.SetParam("x", "(ow-iw)/2")
	padNode.SetParam("y", "(oh-ih)/2")
	padNode.SetParam("color", "black@0.0")

	padIn := padNode.AddInput(scaleOut.ID, filtergraph.StreamTypeVideo)
	padOut := padNode.AddOutput(g.NextPadID("pad_out"), filtergraph.StreamTypeVideo)
	if err := g.Connect(scaleOut, padIn); err != nil {
		return nil, err
	}

	// SetSAR node
	sarNode := g.NewNode(fmt.Sprintf("%s_sar", normID), "setsar")
	sarNode.SetParam("sar", "1")
	sarIn := sarNode.AddInput(padOut.ID, filtergraph.StreamTypeVideo)
	sarOut := sarNode.AddOutput(g.NextPadID("sar_out"), filtergraph.StreamTypeVideo)
	if err := g.Connect(padOut, sarIn); err != nil {
		return nil, err
	}

	// FPS node
	fpsNode := g.NewNode(fmt.Sprintf("%s_fps", normID), "fps")
	fpsNode.SetParam("fps", fps.FFmpegString())
	fpsIn := fpsNode.AddInput(sarOut.ID, filtergraph.StreamTypeVideo)
	fpsOut := fpsNode.AddOutput(g.NextPadID("fps_out"), filtergraph.StreamTypeVideo)
	if err := g.Connect(sarOut, fpsIn); err != nil {
		return nil, err
	}

	// Format node
	formatNode := g.NewNode(fmt.Sprintf("%s_fmt", normID), "format")
	formatNode.SetParam("pix_fmts", "yuv420p")
	fmtIn := formatNode.AddInput(fpsOut.ID, filtergraph.StreamTypeVideo)
	finalOut := formatNode.AddOutput(normID, filtergraph.StreamTypeVideo)
	if err := g.Connect(fpsOut, fmtIn); err != nil {
		return nil, err
	}

	return finalOut, nil
}

// InjectAudioNormalizer creates a normalization chain ensuring 48kHz stereo FLTP audio:
// [in] -> aresample=48000 -> aformat=sample_fmts=fltp:channel_layouts=stereo -> [out]
func InjectAudioNormalizer(g *filtergraph.Graph, inPad *filtergraph.Pad) (*filtergraph.Pad, error) {
	normID := g.NextPadID("a_norm")

	resampleNode := g.NewNode(fmt.Sprintf("%s_resample", normID), "aresample")
	resampleNode.SetParam("sample_rate", 48000)
	resampleIn := resampleNode.AddInput(inPad.ID, filtergraph.StreamTypeAudio)
	resampleOut := resampleNode.AddOutput(g.NextPadID("res_out"), filtergraph.StreamTypeAudio)
	if err := g.Connect(inPad, resampleIn); err != nil {
		return nil, err
	}

	formatNode := g.NewNode(fmt.Sprintf("%s_aformat", normID), "aformat")
	formatNode.SetParam("sample_fmts", "fltp")
	formatNode.SetParam("channel_layouts", "stereo")
	fmtIn := formatNode.AddInput(resampleOut.ID, filtergraph.StreamTypeAudio)
	finalOut := formatNode.AddOutput(normID, filtergraph.StreamTypeAudio)
	if err := g.Connect(resampleOut, fmtIn); err != nil {
		return nil, err
	}

	return finalOut, nil
}
