package timeline

import (
	"time"

	"github.com/farshidrezaei/vidonex/animation"
	"github.com/farshidrezaei/vidonex/chromakey"
	"github.com/farshidrezaei/vidonex/types"
)

// Clip represents a media element placed at a specific time range in a track.
type Clip struct {
	ID                string
	Source            string // Filepath or URL
	TimelineStart     time.Duration
	Duration          time.Duration
	SourceStart       time.Duration // Trim in offset within the source media
	Speed             float64       // Playback rate multiplier (default 1.0)
	Volume            float64       // Audio volume multiplier (default 1.0)
	Opacity           float64       // Visual opacity 0.0 (transparent) to 1.0 (opaque)
	Position          types.Point   // Absolute (X, Y) coordinate on canvas
	Alignment         types.Alignment
	Scale             float64  // Scale factor (1.0 = original size)
	Rotation          float64  // Rotation in degrees (e.g. 90.0, 180.0)
	BlendMode         string   // Blending mode for compositing (e.g. "normal", "multiply", "screen", "overlay")
	HasAudioStream    bool     // Indicates whether the source file contains audio (default true for non-image sources)
	FadeInDuration    time.Duration
	FadeOutDuration   time.Duration
	PositionTrack     *animation.PositionTrack
	ScaleTrack        *animation.FloatKeyframeTrack
	OpacityTrack      *animation.FloatKeyframeTrack
	ChromaKeyOptions  *chromakey.Options
	Effects           []Effect // Attached filters and visual transformations
}

// NewClip creates a Clip with sensible defaults.
func NewClip(id, source string, start, duration time.Duration) *Clip {
	return &Clip{
		ID:             id,
		Source:         source,
		TimelineStart:  start,
		Duration:       duration,
		SourceStart:    0,
		Speed:          1.0,
		Volume:         1.0,
		Opacity:        1.0,
		Scale:          1.0,
		HasAudioStream: true,
		Alignment:      types.AlignCenter,
		Effects:        make([]Effect, 0),
	}
}

// TimelineEnd returns the absolute end timestamp of the clip on the timeline.
func (clip *Clip) TimelineEnd() time.Duration {
	return clip.TimelineStart + clip.Duration
}

// WithTrim sets the in-point offset within the source media.
func (clip *Clip) WithTrim(sourceStart time.Duration) *Clip {
	clip.SourceStart = sourceStart
	return clip
}

// WithSpeed sets playback speed rate (e.g. 2.0 = double speed, 0.5 = slow motion).
func (clip *Clip) WithSpeed(speed float64) *Clip {
	clip.Speed = speed
	return clip
}

// WithVolume sets the audio volume level (1.0 = 100%).
func (clip *Clip) WithVolume(volume float64) *Clip {
	clip.Volume = volume
	return clip
}

// WithOpacity sets the visual transparency level (0.0 to 1.0).
func (clip *Clip) WithOpacity(opacity float64) *Clip {
	clip.Opacity = opacity
	return clip
}

// WithChromaKey attaches green/blue screen removal configuration to the clip.
func (clip *Clip) WithChromaKey(customOptions ...chromakey.Options) *Clip {
	if len(customOptions) > 0 {
		clip.ChromaKeyOptions = &customOptions[0]
	} else {
		opts := chromakey.DefaultOptions()
		clip.ChromaKeyOptions = &opts
	}
	return clip
}

// WithFadeIn configures an opacity and audio fade-in transition at the start of the clip.
func (clip *Clip) WithFadeIn(duration time.Duration) *Clip {
	clip.FadeInDuration = duration
	return clip
}

// WithFadeOut configures an opacity and audio fade-out transition at the end of the clip.
func (clip *Clip) WithFadeOut(duration time.Duration) *Clip {
	clip.FadeOutDuration = duration
	return clip
}

// WithPosition sets the exact (X, Y) canvas coordinates.
func (clip *Clip) WithPosition(position types.Point) *Clip {
	clip.Position = position
	return clip
}

// WithPositionTrack attaches an animated position keyframe track.
func (clip *Clip) WithPositionTrack(track *animation.PositionTrack) *Clip {
	clip.PositionTrack = track
	return clip
}

// WithScale sets the size scale multiplier.
func (clip *Clip) WithScale(scale float64) *Clip {
	clip.Scale = scale
	return clip
}

// WithScaleTrack attaches an animated scale keyframe track.
func (clip *Clip) WithScaleTrack(track *animation.FloatKeyframeTrack) *Clip {
	clip.ScaleTrack = track
	return clip
}

// WithOpacityTrack attaches an animated opacity keyframe track.
func (clip *Clip) WithOpacityTrack(track *animation.FloatKeyframeTrack) *Clip {
	clip.OpacityTrack = track
	return clip
}

// WithRotation sets the rotation in degrees.
func (clip *Clip) WithRotation(degrees float64) *Clip {
	clip.Rotation = degrees
	return clip
}

// WithBlendMode sets the compositing blend mode.
func (clip *Clip) WithBlendMode(mode string) *Clip {
	clip.BlendMode = mode
	return clip
}

// WithHasAudio sets whether the clip source contains an audio stream.
func (clip *Clip) WithHasAudio(hasAudio bool) *Clip {
	clip.HasAudioStream = hasAudio
	return clip
}

// AddEffect attaches a filter effect to the clip.
func (clip *Clip) AddEffect(effect Effect) *Clip {
	clip.Effects = append(clip.Effects, effect)
	return clip
}
