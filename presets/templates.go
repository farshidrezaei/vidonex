// Package presets provides platform templates, hardware acceleration flags, and turnkey timeline blueprints.
package presets

import (
	"errors"
	"time"

	"github.com/farshidrezaei/vidonex/timeline"
	"github.com/farshidrezaei/vidonex/types"
)

// TemplateTikTokSplitScreen constructs a 9:16 vertical timeline (1080x1920) with a split-screen layout.
// The top half (50% height) hosts the reaction or gameplay video, and the bottom half hosts the creator facecam.
func TemplateTikTokSplitScreen(
	topVideoSourcePath, bottomVideoSourcePath string,
	duration time.Duration,
) (*timeline.Timeline, error) {
	if topVideoSourcePath == "" || bottomVideoSourcePath == "" {
		return nil, errors.New("presets: TemplateTikTokSplitScreen requires non-empty top and bottom video paths")
	}
	if duration <= 0 {
		return nil, errors.New("presets: TemplateTikTokSplitScreen requires a positive duration")
	}

	compositionTimeline := timeline.New(
		timeline.WithCanvas(types.ResPortrait1080p),
		timeline.WithFPS(types.FPS60),
	)

	// Base Track (Top Video - Top Half)
	topTrack := timeline.NewTrack("track_top_video", timeline.TrackKindVideo).SetZIndex(0)
	topClip := timeline.NewClip("clip_top_video", topVideoSourcePath, 0, duration).
		WithScale(1.0).
		WithPosition(types.Point{X: 0, Y: -480})
	topTrack.AddClip(topClip)

	// Bottom Track (Bottom Video - Bottom Half)
	bottomTrack := timeline.NewTrack("track_bottom_video", timeline.TrackKindOverlay).SetZIndex(1)
	bottomClip := timeline.NewClip("clip_bottom_video", bottomVideoSourcePath, 0, duration).
		WithScale(1.0).
		WithPosition(types.Point{X: 0, Y: 480})
	bottomTrack.AddClip(bottomClip)

	compositionTimeline.AddTrack(topTrack)
	compositionTimeline.AddTrack(bottomTrack)

	return compositionTimeline, nil
}

// TemplatePodcastAudiogram constructs an audiogram timeline with an audio track, dynamic audio waveform,
// background graphic, and title overlay.
func TemplatePodcastAudiogram(
	audioSourcePath, backgroundSourcePath string,
	duration time.Duration,
	podcastTitleText string,
) (*timeline.Timeline, error) {
	if audioSourcePath == "" {
		return nil, errors.New("presets: TemplatePodcastAudiogram requires a non-empty audio source path")
	}
	if duration <= 0 {
		return nil, errors.New("presets: TemplatePodcastAudiogram requires a positive duration")
	}

	compositionTimeline := timeline.New(
		timeline.WithCanvas(types.ResSquare1080),
		timeline.WithFPS(types.FPS30),
	)

	// 1. Background Image / Video Layer
	if backgroundSourcePath != "" {
		backgroundTrack := timeline.NewTrack("track_background", timeline.TrackKindVideo).SetZIndex(0)
		bgClip := timeline.NewClip("clip_background", backgroundSourcePath, 0, duration)
		if podcastTitleText != "" {
			bgClip.AddEffect(&timeline.DrawTextEffect{
				Text:      podcastTitleText,
				FontSize:  48,
				FontColor: "white",
				X:         "(w-text_w)/2",
				Y:         "100",
			})
		}
		backgroundTrack.AddClip(bgClip)
		compositionTimeline.AddTrack(backgroundTrack)
	}

	// 2. Dynamic Audio Waveform Overlay Layer
	waveformTrack := timeline.NewTrack("track_waveform", timeline.TrackKindWaveform).SetZIndex(1)
	waveformClip := timeline.NewClip("clip_audiogram", audioSourcePath, 0, duration).
		WithPosition(types.Point{X: 0, Y: 200})
	waveformTrack.AddClip(waveformClip)
	compositionTimeline.AddTrack(waveformTrack)

	// 3. Audio Track
	audioTrack := timeline.NewTrack("track_audio_vocal", timeline.TrackKindAudio)
	audioClip := timeline.NewClip("clip_vocal", audioSourcePath, 0, duration)
	audioTrack.AddClip(audioClip)
	compositionTimeline.AddTrack(audioTrack)

	return compositionTimeline, nil
}

// TemplateYouTubeEndScreen constructs a 16:9 landscape timeline (1920x1080) with outro video and placeholder slots.
func TemplateYouTubeEndScreen(
	outroVideoSourcePath string,
	duration time.Duration,
	channelTitle string,
) (*timeline.Timeline, error) {
	if outroVideoSourcePath == "" {
		return nil, errors.New("presets: TemplateYouTubeEndScreen requires a non-empty video source path")
	}
	if duration <= 0 {
		return nil, errors.New("presets: TemplateYouTubeEndScreen requires a positive duration")
	}

	compositionTimeline := timeline.New(
		timeline.WithCanvas(types.Res1080p),
		timeline.WithFPS(types.FPS60),
	)

	// Main Outro Video
	baseTrack := timeline.NewTrack("track_outro_video", timeline.TrackKindVideo).SetZIndex(0)
	videoClip := timeline.NewClip("clip_outro", outroVideoSourcePath, 0, duration)
	if channelTitle != "" {
		videoClip.AddEffect(&timeline.DrawTextEffect{
			Text:      channelTitle,
			FontSize:  36,
			FontColor: "white",
			X:         "(w-text_w)/2",
			Y:         "h-150",
		})
	}
	baseTrack.AddClip(videoClip)
	compositionTimeline.AddTrack(baseTrack)

	return compositionTimeline, nil
}
