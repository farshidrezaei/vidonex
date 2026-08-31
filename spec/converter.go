package spec

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/farshidrezaei/vidonyx/animation"
	"github.com/farshidrezaei/vidonyx/chromakey"
	"github.com/farshidrezaei/vidonyx/compiler"
	"github.com/farshidrezaei/vidonyx/ducking"
	"github.com/farshidrezaei/vidonyx/presets"
	"github.com/farshidrezaei/vidonyx/subtitles"
	"github.com/farshidrezaei/vidonyx/timeline"
	"github.com/farshidrezaei/vidonyx/types"
	"github.com/farshidrezaei/vidonyx/waveform"
)

// ToTimeline compiles the VideoSpec into a timeline.Timeline AST and compiler.EncodingOptions.
func ToTimeline(specification *VideoSpec) (*timeline.Timeline, *compiler.EncodingOptions, string, error) {
	if specification == nil {
		return nil, nil, "", fmt.Errorf("spec.ToTimeline: specification cannot be nil")
	}

	// 1. Resolve Platform Preset (if specified at root level)
	var platformPreset *presets.PlatformPreset
	if specification.Preset != "" {
		resolvedPreset, err := resolvePlatformPreset(specification.Preset)
		if err != nil {
			return nil, nil, "", err
		}
		platformPreset = &resolvedPreset
	}

	// 2. Resolve Canvas Size
	canvasSize, err := resolveCanvasSize(specification.Canvas, platformPreset)
	if err != nil {
		return nil, nil, "", err
	}

	// 3. Resolve Frame Rate
	frameRate, err := resolveFrameRate(specification.FPS, platformPreset)
	if err != nil {
		return nil, nil, "", err
	}

	// 4. Resolve Background Color
	backgroundColor := types.ColorBlack
	if specification.Canvas.BackgroundColor != "" {
		parsedColor, colorErr := ParseColorValue(specification.Canvas.BackgroundColor)
		if colorErr != nil {
			return nil, nil, "", fmt.Errorf("spec.ToTimeline: invalid background color: %w", colorErr)
		}
		backgroundColor = parsedColor
	}

	compositionTimeline := timeline.New(
		timeline.WithCanvas(canvasSize),
		timeline.WithFPS(frameRate),
		timeline.WithBackgroundColor(backgroundColor),
	)

	// 5. Build Tracks
	for _, trackSpec := range specification.Tracks {
		builtTrack, trackErr := buildTrack(trackSpec)
		if trackErr != nil {
			return nil, nil, "", trackErr
		}
		compositionTimeline.AddTrack(builtTrack)
	}

	// Validate the reconstructed timeline
	if err := timeline.Validate(compositionTimeline); err != nil {
		return nil, nil, "", fmt.Errorf("spec.ToTimeline: timeline validation failed: %w", err)
	}

	// 6. Build Encoding Options
	encodingOptions := resolveEncodingOptions(specification, platformPreset)

	outputPath := specification.Output
	if outputPath == "" {
		outputPath = "output.mp4"
	}

	return compositionTimeline, encodingOptions, outputPath, nil
}

func resolvePlatformPreset(presetName string) (presets.PlatformPreset, error) {
	cleanName := strings.ToLower(strings.TrimSpace(presetName))
	switch cleanName {
	case "tiktok", "tiktok_1080p60", "tiktok_vertical", "reels_60":
		return presets.TikTokVertical1080p60(), nil
	case "instagram_reels", "reels_30", "reels":
		return presets.InstagramReels1080p30(), nil
	case "instagram_square", "square", "square_1080":
		return presets.InstagramSquare1080p(), nil
	case "youtube_4k", "youtube_4k60", "4k60":
		return presets.YouTube4K60(), nil
	case "youtube_1080", "youtube_1080p60", "1080p60":
		return presets.YouTube1080p60(), nil
	case "web_fast", "web_720p30", "fast":
		return presets.WebFast720p30(), nil
	default:
		return presets.PlatformPreset{}, fmt.Errorf("spec.resolvePlatformPreset: unknown preset %q", presetName)
	}
}

func resolveCanvasSize(canvas CanvasSpec, preset *presets.PlatformPreset) (types.Size, error) {
	if preset != nil {
		return preset.Canvas, nil
	}

	if canvas.Width > 0 && canvas.Height > 0 {
		return types.NewSize(canvas.Width, canvas.Height), nil
	}

	switch strings.ToLower(strings.TrimSpace(canvas.Preset)) {
	case "1080p", "fhd":
		return types.Res1080p, nil
	case "720p", "hd":
		return types.Res720p, nil
	case "4k", "uhd":
		return types.Res4K, nil
	case "square_1080", "square":
		return types.ResSquare1080, nil
	case "portrait_1080", "vertical_1080", "tiktok":
		return types.ResPortrait1080p, nil
	case "":
		return types.Res1080p, nil
	default:
		return types.Size{}, fmt.Errorf("spec.resolveCanvasSize: unknown canvas preset %q", canvas.Preset)
	}
}

func resolveFrameRate(rawFPS any, preset *presets.PlatformPreset) (types.Rational, error) {
	if preset != nil {
		return preset.FPS, nil
	}

	if rawFPS == nil {
		return types.FPS30, nil
	}

	switch value := rawFPS.(type) {
	case int:
		return types.NewRational(int64(value), 1), nil
	case float64:
		if value == 29.97 {
			return types.FPS2997, nil
		}
		if value == 59.94 {
			return types.FPS5994, nil
		}
		return types.NewRational(int64(value*1000), 1000), nil
	case string:
		clean := strings.TrimSpace(value)
		switch clean {
		case "30", "30fps":
			return types.FPS30, nil
		case "60", "60fps":
			return types.FPS60, nil
		case "24", "24fps":
			return types.FPS24, nil
		case "25", "25fps":
			return types.FPS25, nil
		case "29.97":
			return types.FPS2997, nil
		case "59.94":
			return types.FPS5994, nil
		default:
			parsedNum, err := strconv.Atoi(clean)
			if err != nil {
				return types.FPS30, fmt.Errorf("spec.resolveFrameRate: invalid fps string %q: %w", clean, err)
			}
			return types.NewRational(int64(parsedNum), 1), nil
		}
	default:
		return types.FPS30, fmt.Errorf("spec.resolveFrameRate: invalid fps format %T", rawFPS)
	}
}

func buildTrack(trackSpec TrackSpec) (*timeline.Track, error) {
	// Handle Subtitle Track
	if trackSpec.Kind == "subtitle" || trackSpec.Subtitles != nil {
		return buildSubtitleTrack(trackSpec)
	}

	// Handle Waveform Track
	if trackSpec.Kind == "waveform" || trackSpec.Waveform != nil {
		return buildWaveformTrack(trackSpec)
	}

	var trackKind timeline.TrackKind
	switch strings.ToLower(trackSpec.Kind) {
	case "audio":
		trackKind = timeline.TrackKindAudio
	case "overlay":
		trackKind = timeline.TrackKindOverlay
	case "video", "":
		trackKind = timeline.TrackKindVideo
	default:
		return nil, fmt.Errorf("spec.buildTrack: unknown track kind %q", trackSpec.Kind)
	}

	track := timeline.NewTrack(trackSpec.ID, trackKind).
		SetZIndex(trackSpec.ZIndex).
		SetMuted(trackSpec.Muted)

	if trackSpec.Volume != nil {
		track.SetVolume(*trackSpec.Volume)
	}

	// Audio Ducking configuration
	if trackSpec.DuckUnder != "" {
		duckingOptions := ducking.DefaultOptions()
		if trackSpec.Ducking != nil {
			if trackSpec.Ducking.Threshold > 0 {
				duckingOptions.Threshold = trackSpec.Ducking.Threshold
			}
			if trackSpec.Ducking.Ratio > 0 {
				duckingOptions.Ratio = trackSpec.Ducking.Ratio
			}
			if trackSpec.Ducking.AttackMilliseconds > 0 {
				duckingOptions.AttackMilliseconds = trackSpec.Ducking.AttackMilliseconds
			}
			if trackSpec.Ducking.ReleaseMilliseconds > 0 {
				duckingOptions.ReleaseMilliseconds = trackSpec.Ducking.ReleaseMilliseconds
			}
		}
		track.WithDucking(trackSpec.DuckUnder, duckingOptions)
	}

	// Build Clips
	clipMap := make(map[string]*timeline.Clip)
	for _, clipSpec := range trackSpec.Clips {
		builtClip, clipErr := buildClip(clipSpec)
		if clipErr != nil {
			return nil, clipErr
		}
		track.AddClip(builtClip)
		clipMap[builtClip.ID] = builtClip
	}

	// Build Transitions
	for _, transitionSpec := range trackSpec.Transitions {
		builtTransition, transErr := buildTransition(transitionSpec, clipMap)
		if transErr != nil {
			return nil, transErr
		}
		track.AddTransition(builtTransition)
	}

	return track, nil
}

func buildSubtitleTrack(trackSpec TrackSpec) (*timeline.Track, error) {
	subtitlesSpec := trackSpec.Subtitles
	if subtitlesSpec == nil {
		return nil, fmt.Errorf("spec.buildSubtitleTrack: track %q is kind subtitle but missing subtitles configuration", trackSpec.ID)
	}

	var parsedTrack *subtitles.SubtitleTrack
	var err error

	if subtitlesSpec.SRTContent != "" {
		parsedTrack, err = subtitles.ParseSRT(strings.NewReader(subtitlesSpec.SRTContent))
		if err != nil {
			return nil, fmt.Errorf("spec.buildSubtitleTrack: failed parsing srt_content: %w", err)
		}
	} else if subtitlesSpec.VTTContent != "" {
		parsedTrack, err = subtitles.ParseVTT(strings.NewReader(subtitlesSpec.VTTContent))
		if err != nil {
			return nil, fmt.Errorf("spec.buildSubtitleTrack: failed parsing vtt_content: %w", err)
		}
	} else if subtitlesSpec.File != "" {
		fileBytes, readErr := os.ReadFile(subtitlesSpec.File)
		if readErr != nil {
			return nil, fmt.Errorf("spec.buildSubtitleTrack: failed reading subtitle file %q: %w", subtitlesSpec.File, readErr)
		}
		if strings.HasSuffix(strings.ToLower(subtitlesSpec.File), ".vtt") {
			parsedTrack, err = subtitles.ParseVTT(strings.NewReader(string(fileBytes)))
		} else {
			parsedTrack, err = subtitles.ParseSRT(strings.NewReader(string(fileBytes)))
		}
		if err != nil {
			return nil, fmt.Errorf("spec.buildSubtitleTrack: failed parsing subtitle file %q: %w", subtitlesSpec.File, err)
		}
	} else {
		return nil, fmt.Errorf("spec.buildSubtitleTrack: subtitle track %q requires srt_content, vtt_content, or file", trackSpec.ID)
	}

	// Apply Style if defined
	if subtitlesSpec.Style != nil {
		styleSpec := subtitlesSpec.Style
		style := subtitles.SubtitleStyle{
			FontSize:       styleSpec.FontSize,
			OutlineWidth:   styleSpec.OutlineWidth,
			Box:            styleSpec.Box,
			BoxBorderWidth: styleSpec.BoxBorderWidth,
			Alignment:      ParseAlignmentValue(styleSpec.Alignment),
			MarginBottom:   styleSpec.MarginBottom,
		}
		if styleSpec.Color != "" {
			parsedColor, colorErr := ParseColorValue(styleSpec.Color)
			if colorErr != nil {
				return nil, fmt.Errorf("spec.buildSubtitleTrack: invalid subtitle color: %w", colorErr)
			}
			style.PrimaryColor = parsedColor
		}
		if styleSpec.OutlineColor != "" {
			parsedColor, colorErr := ParseColorValue(styleSpec.OutlineColor)
			if colorErr != nil {
				return nil, fmt.Errorf("spec.buildSubtitleTrack: invalid subtitle outline_color: %w", colorErr)
			}
			style.OutlineColor = parsedColor
		}
		if styleSpec.BoxColor != "" {
			parsedColor, colorErr := ParseColorValue(styleSpec.BoxColor)
			if colorErr != nil {
				return nil, fmt.Errorf("spec.buildSubtitleTrack: invalid subtitle box_color: %w", colorErr)
			}
			style.BoxColor = parsedColor
		}
		parsedTrack.SetStyle(style)
	}

	return timeline.NewSubtitleTrack(trackSpec.ID, parsedTrack), nil
}

func buildWaveformTrack(trackSpec TrackSpec) (*timeline.Track, error) {
	waveSpec := trackSpec.Waveform
	if waveSpec == nil {
		return nil, fmt.Errorf("spec.buildWaveformTrack: track %q is missing waveform configuration", trackSpec.ID)
	}

	waveformOptions := waveform.DefaultOptions()
	if waveSpec.Mode != "" {
		waveformOptions.Mode = waveform.Mode(waveSpec.Mode)
	}
	if waveSpec.Scale != "" {
		waveformOptions.Scale = waveSpec.Scale
	}
	if waveSpec.Width > 0 && waveSpec.Height > 0 {
		waveformOptions.Size = types.NewSize(waveSpec.Width, waveSpec.Height)
	}
	if waveSpec.Color != "" {
		parsedColor, colorErr := ParseColorValue(waveSpec.Color)
		if colorErr != nil {
			return nil, fmt.Errorf("spec.buildWaveformTrack: invalid waveform color: %w", colorErr)
		}
		waveformOptions.Color = parsedColor
	}
	if waveSpec.Position != nil {
		waveformOptions.Position = types.Point{X: waveSpec.Position.X, Y: waveSpec.Position.Y}
	}
	if waveSpec.Opacity > 0 {
		waveformOptions.Opacity = waveSpec.Opacity
	}

	return timeline.NewWaveformTrack(trackSpec.ID, waveSpec.SourceAudio, waveformOptions), nil
}

func buildClip(clipSpec ClipSpec) (*timeline.Clip, error) {
	startTime, err := ParseDurationValue(clipSpec.Start)
	if err != nil {
		return nil, fmt.Errorf("spec.buildClip: invalid start time for clip %q: %w", clipSpec.ID, err)
	}

	durationTime, err := ParseDurationValue(clipSpec.Duration)
	if err != nil {
		return nil, fmt.Errorf("spec.buildClip: invalid duration for clip %q: %w", clipSpec.ID, err)
	}

	clip := timeline.NewClip(clipSpec.ID, clipSpec.Source, startTime, durationTime)

	if clipSpec.Trim != nil {
		trimTime, trimErr := ParseDurationValue(clipSpec.Trim)
		if trimErr != nil {
			return nil, fmt.Errorf("spec.buildClip: invalid trim for clip %q: %w", clipSpec.ID, trimErr)
		}
		clip.WithTrim(trimTime)
	}

	if clipSpec.Speed > 0 {
		clip.WithSpeed(clipSpec.Speed)
	}
	if clipSpec.Volume != nil {
		clip.WithVolume(*clipSpec.Volume)
	}
	if clipSpec.Opacity != nil {
		clip.WithOpacity(*clipSpec.Opacity)
	}
	if clipSpec.Scale > 0 {
		clip.WithScale(clipSpec.Scale)
	}
	if clipSpec.Rotation != 0 {
		clip.WithRotation(clipSpec.Rotation)
	}
	if clipSpec.FadeIn != nil {
		fadeInTime, fadeErr := ParseDurationValue(clipSpec.FadeIn)
		if fadeErr != nil {
			return nil, fmt.Errorf("spec.buildClip: invalid fade_in for clip %q: %w", clipSpec.ID, fadeErr)
		}
		clip.WithFadeIn(fadeInTime)
	}
	if clipSpec.FadeOut != nil {
		fadeOutTime, fadeErr := ParseDurationValue(clipSpec.FadeOut)
		if fadeErr != nil {
			return nil, fmt.Errorf("spec.buildClip: invalid fade_out for clip %q: %w", clipSpec.ID, fadeErr)
		}
		clip.WithFadeOut(fadeOutTime)
	}
	if clipSpec.Position != nil {
		clip.WithPosition(types.Point{X: clipSpec.Position.X, Y: clipSpec.Position.Y})
	}

	// Chroma Key
	if clipSpec.ChromaKey != nil {
		chromaOptions := chromakey.DefaultOptions()
		if clipSpec.ChromaKey.KeyColor != "" {
			parsedColor, colorErr := ParseColorValue(clipSpec.ChromaKey.KeyColor)
			if colorErr != nil {
				return nil, fmt.Errorf("spec.buildClip: invalid chroma key color: %w", colorErr)
			}
			chromaOptions.KeyColor = parsedColor
		}
		if clipSpec.ChromaKey.Similarity > 0 {
			chromaOptions.Similarity = clipSpec.ChromaKey.Similarity
		}
		if clipSpec.ChromaKey.Blend > 0 {
			chromaOptions.Blend = clipSpec.ChromaKey.Blend
		}
		chromaOptions.Despill = clipSpec.ChromaKey.Despill
		if clipSpec.ChromaKey.DespillType != "" {
			chromaOptions.DespillType = chromakey.DespillMode(clipSpec.ChromaKey.DespillType)
		}
		if clipSpec.ChromaKey.DespillExpand > 0 {
			chromaOptions.DespillExpand = clipSpec.ChromaKey.DespillExpand
		}
		clip.WithChromaKey(chromaOptions)
	}

	// Keyframe Animation Tracks
	if clipSpec.Keyframes != nil {
		if len(clipSpec.Keyframes.Position) > 0 {
			posTrack := animation.NewPositionTrack()
			for _, kf := range clipSpec.Keyframes.Position {
				kfTime, kfErr := ParseDurationValue(kf.Time)
				if kfErr != nil {
					return nil, fmt.Errorf("spec.buildClip: invalid position keyframe time: %w", kfErr)
				}
				posTrack.AddKeyframe(kfTime, types.Point{X: kf.X, Y: kf.Y}, ParseEasingValue(kf.Easing))
			}
			clip.WithPositionTrack(posTrack)
		}

		if len(clipSpec.Keyframes.Scale) > 0 {
			scaleTrack := animation.NewFloatKeyframeTrack()
			for _, kf := range clipSpec.Keyframes.Scale {
				kfTime, kfErr := ParseDurationValue(kf.Time)
				if kfErr != nil {
					return nil, fmt.Errorf("spec.buildClip: invalid scale keyframe time: %w", kfErr)
				}
				scaleTrack.AddKeyframe(kfTime, kf.Value, ParseEasingValue(kf.Easing))
			}
			clip.WithScaleTrack(scaleTrack)
		}

		if len(clipSpec.Keyframes.Opacity) > 0 {
			opacityTrack := animation.NewFloatKeyframeTrack()
			for _, kf := range clipSpec.Keyframes.Opacity {
				kfTime, kfErr := ParseDurationValue(kf.Time)
				if kfErr != nil {
					return nil, fmt.Errorf("spec.buildClip: invalid opacity keyframe time: %w", kfErr)
				}
				opacityTrack.AddKeyframe(kfTime, kf.Value, ParseEasingValue(kf.Easing))
			}
			clip.WithOpacityTrack(opacityTrack)
		}
	}

	return clip, nil
}

func buildTransition(transitionSpec TransitionSpec, clipMap map[string]*timeline.Clip) (*timeline.Transition, error) {
	durationTime, err := ParseDurationValue(transitionSpec.Duration)
	if err != nil {
		return nil, fmt.Errorf("spec.buildTransition: invalid transition duration: %w", err)
	}

	var clipFrom, clipTo *timeline.Clip
	if transitionSpec.From != "" {
		clipFrom = clipMap[transitionSpec.From]
	}
	if transitionSpec.To != "" {
		clipTo = clipMap[transitionSpec.To]
	}

	transitionKind := timeline.TransitionType(strings.ToLower(transitionSpec.Type))
	if transitionKind == "" {
		transitionKind = timeline.TransitionDissolve
	}

	transitionID := transitionSpec.ID
	if transitionID == "" {
		transitionID = fmt.Sprintf("trans_%s", transitionKind)
	}

	return timeline.NewTransition(transitionID, transitionKind, durationTime, clipFrom, clipTo), nil
}

func resolveEncodingOptions(specification *VideoSpec, platformPreset *presets.PlatformPreset) *compiler.EncodingOptions {
	encodingOptions := compiler.DefaultEncodingOptions()

	if platformPreset != nil {
		encodingOptions.VideoCodec = platformPreset.Encoding.VideoCodec
		encodingOptions.AudioCodec = platformPreset.Encoding.AudioCodec
		encodingOptions.PixelFormat = platformPreset.Encoding.PixelFormat
		encodingOptions.AudioBitrate = platformPreset.Encoding.AudioBitrate
		encodingOptions.CRF = platformPreset.Encoding.CRF
		encodingOptions.Preset = platformPreset.Encoding.Preset
	}

	// Apply Hardware Accelerator
	if specification.HardwareAcceleration != "" && strings.ToLower(specification.HardwareAcceleration) != "none" {
		accelerator := presets.HardwareAccelerator(strings.ToLower(specification.HardwareAcceleration))
		presets.ConfigureHardwareEncoding(&encodingOptions, accelerator, false)
	}

	// Custom Encoding Overrides
	if specification.Encoding != nil {
		if specification.Encoding.VideoCodec != "" {
			encodingOptions.VideoCodec = specification.Encoding.VideoCodec
		}
		if specification.Encoding.AudioCodec != "" {
			encodingOptions.AudioCodec = specification.Encoding.AudioCodec
		}
		if specification.Encoding.PixelFormat != "" {
			encodingOptions.PixelFormat = specification.Encoding.PixelFormat
		}
		if specification.Encoding.AudioBitrate != "" {
			encodingOptions.AudioBitrate = specification.Encoding.AudioBitrate
		}
		if specification.Encoding.CRF > 0 {
			encodingOptions.CRF = specification.Encoding.CRF
		}
		if specification.Encoding.Preset != "" {
			encodingOptions.Preset = specification.Encoding.Preset
		}
	}

	return &encodingOptions
}
