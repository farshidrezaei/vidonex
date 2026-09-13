// Package spec provides declarative YAML and JSON specification parsing, relative path resolution, and timeline compilation.
package spec

// VideoSpec represents the top-level declarative project specification for a video composition.
type VideoSpec struct {
	Version              string            `json:"version" yaml:"version"`
	Canvas               CanvasSpec        `json:"canvas" yaml:"canvas"`
	FPS                  any               `json:"fps,omitempty" yaml:"fps,omitempty"`
	Output               string            `json:"output,omitempty" yaml:"output,omitempty"`
	Preset               string            `json:"preset,omitempty" yaml:"preset,omitempty"`
	HardwareAcceleration string            `json:"hardware_acceleration,omitempty" yaml:"hardware_acceleration,omitempty"`
	Encoding             *EncodingSpec     `json:"encoding,omitempty" yaml:"encoding,omitempty"`
	Tracks               []TrackSpec       `json:"tracks" yaml:"tracks"`
	Metadata             map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// CanvasSpec defines the output video resolution and background properties.
type CanvasSpec struct {
	Width           int    `json:"width,omitempty" yaml:"width,omitempty"`
	Height          int    `json:"height,omitempty" yaml:"height,omitempty"`
	Preset          string `json:"preset,omitempty" yaml:"preset,omitempty"` // e.g. "1080p", "720p", "4k", "square_1080", "portrait_1080"
	BackgroundColor string `json:"background_color,omitempty" yaml:"background_color,omitempty"`
}

// EncodingSpec customizes low-level video and audio encoder parameters.
type EncodingSpec struct {
	VideoCodec   string `json:"video_codec,omitempty" yaml:"video_codec,omitempty"`
	AudioCodec   string `json:"audio_codec,omitempty" yaml:"audio_codec,omitempty"`
	PixelFormat  string `json:"pixel_format,omitempty" yaml:"pixel_format,omitempty"`
	AudioBitrate string `json:"audio_bitrate,omitempty" yaml:"audio_bitrate,omitempty"`
	CRF          int    `json:"crf,omitempty" yaml:"crf,omitempty"`
	Preset       string `json:"preset,omitempty" yaml:"preset,omitempty"`
}

// TrackSpec defines a single media layer (video, audio, overlay, subtitle, or waveform).
type TrackSpec struct {
	ID          string           `json:"id" yaml:"id"`
	Kind        string           `json:"kind" yaml:"kind"` // "video", "audio", "overlay", "subtitle", "waveform"
	ZIndex      int              `json:"z_index,omitempty" yaml:"z_index,omitempty"`
	Volume      *float64         `json:"volume,omitempty" yaml:"volume,omitempty"`
	Muted       bool             `json:"muted,omitempty" yaml:"muted,omitempty"`
	DuckUnder   string           `json:"duck_under,omitempty" yaml:"duck_under,omitempty"`
	Ducking     *DuckingSpec     `json:"ducking,omitempty" yaml:"ducking,omitempty"`
	Subtitles   *SubtitlesSpec   `json:"subtitles,omitempty" yaml:"subtitles,omitempty"`
	Waveform    *WaveformSpec    `json:"waveform,omitempty" yaml:"waveform,omitempty"`
	Clips       []ClipSpec       `json:"clips,omitempty" yaml:"clips,omitempty"`
	Transitions []TransitionSpec `json:"transitions,omitempty" yaml:"transitions,omitempty"`
}

// ClipSpec defines an individual media asset placed at a specific time range.
type ClipSpec struct {
	ID          string          `json:"id" yaml:"id"`
	Source      string          `json:"source,omitempty" yaml:"source,omitempty"`
	Start       any             `json:"start,omitempty" yaml:"start,omitempty"` // number, "2.5s", "00:00:02.500"
	Duration    any             `json:"duration,omitempty" yaml:"duration,omitempty"`
	Trim        any             `json:"trim,omitempty" yaml:"trim,omitempty"`
	Speed       float64         `json:"speed,omitempty" yaml:"speed,omitempty"`
	Volume      *float64        `json:"volume,omitempty" yaml:"volume,omitempty"`
	Opacity     *float64        `json:"opacity,omitempty" yaml:"opacity,omitempty"`
	Scale       float64         `json:"scale,omitempty" yaml:"scale,omitempty"`
	Rotation    float64         `json:"rotation,omitempty" yaml:"rotation,omitempty"`
	BlendMode   string          `json:"blend_mode,omitempty" yaml:"blend_mode,omitempty"`
	FadeIn      any             `json:"fade_in,omitempty" yaml:"fade_in,omitempty"`
	FadeOut     any             `json:"fade_out,omitempty" yaml:"fade_out,omitempty"`
	Position    *PositionSpec   `json:"position,omitempty" yaml:"position,omitempty"`
	ChromaKey   *ChromaKeySpec  `json:"chroma_key,omitempty" yaml:"chroma_key,omitempty"`
	Keyframes   *KeyframesSpec  `json:"keyframes,omitempty" yaml:"keyframes,omitempty"`
	HasAudio    *bool           `json:"has_audio,omitempty" yaml:"has_audio,omitempty"`
}

// PositionSpec sets the visual placement on the canvas.
type PositionSpec struct {
	X         int    `json:"x,omitempty" yaml:"x,omitempty"`
	Y         int    `json:"y,omitempty" yaml:"y,omitempty"`
	Alignment string `json:"alignment,omitempty" yaml:"alignment,omitempty"` // "center", "top_left", "bottom_center", etc.
}

// TransitionSpec configures video/audio transitions between adjacent clips.
type TransitionSpec struct {
	ID       string `json:"id,omitempty" yaml:"id,omitempty"`
	Type     string `json:"type" yaml:"type"` // "dissolve", "fade", "wipeleft", "wiperight", "circleopen"
	Duration any    `json:"duration" yaml:"duration"`
	From     string `json:"from,omitempty" yaml:"from,omitempty"`
	To       string `json:"to,omitempty" yaml:"to,omitempty"`
}

// SubtitlesSpec defines subtitle burning options from SRT/VTT content or files.
type SubtitlesSpec struct {
	File       string             `json:"file,omitempty" yaml:"file,omitempty"`
	SRTContent string             `json:"srt_content,omitempty" yaml:"srt_content,omitempty"`
	VTTContent string             `json:"vtt_content,omitempty" yaml:"vtt_content,omitempty"`
	Style      *SubtitleStyleSpec `json:"style,omitempty" yaml:"style,omitempty"`
}

// SubtitleStyleSpec configures font size, colors, bounding boxes, and alignment.
type SubtitleStyleSpec struct {
	FontSize       int    `json:"font_size,omitempty" yaml:"font_size,omitempty"`
	Color          string `json:"color,omitempty" yaml:"color,omitempty"`
	OutlineColor   string `json:"outline_color,omitempty" yaml:"outline_color,omitempty"`
	OutlineWidth   int    `json:"outline_width,omitempty" yaml:"outline_width,omitempty"`
	Box            bool   `json:"box,omitempty" yaml:"box,omitempty"`
	BoxColor       string `json:"box_color,omitempty" yaml:"box_color,omitempty"`
	BoxBorderWidth int    `json:"box_border_width,omitempty" yaml:"box_border_width,omitempty"`
	Alignment      string `json:"alignment,omitempty" yaml:"alignment,omitempty"`
	MarginBottom   int    `json:"margin_bottom,omitempty" yaml:"margin_bottom,omitempty"`
}

// DuckingSpec configures sidechain compression parameters.
type DuckingSpec struct {
	Threshold           float64 `json:"threshold,omitempty" yaml:"threshold,omitempty"`
	Ratio               float64 `json:"ratio,omitempty" yaml:"ratio,omitempty"`
	AttackMilliseconds  int     `json:"attack_ms,omitempty" yaml:"attack_ms,omitempty"`
	ReleaseMilliseconds int     `json:"release_ms,omitempty" yaml:"release_ms,omitempty"`
}

// WaveformSpec configures animated audio waveform generation.
type WaveformSpec struct {
	SourceAudio string        `json:"source_audio" yaml:"source_audio"`
	Mode        string        `json:"mode,omitempty" yaml:"mode,omitempty"` // "p2p", "line", "cline", "dot"
	Color       string        `json:"color,omitempty" yaml:"color,omitempty"`
	Scale       string        `json:"scale,omitempty" yaml:"scale,omitempty"` // "sqrt", "log", "lin"
	Width       int           `json:"width,omitempty" yaml:"width,omitempty"`
	Height      int           `json:"height,omitempty" yaml:"height,omitempty"`
	Position    *PositionSpec `json:"position,omitempty" yaml:"position,omitempty"`
	Opacity     float64       `json:"opacity,omitempty" yaml:"opacity,omitempty"`
}

// ChromaKeySpec configures green/blue screen removal.
type ChromaKeySpec struct {
	KeyColor      string  `json:"key_color,omitempty" yaml:"key_color,omitempty"`
	Similarity    float64 `json:"similarity,omitempty" yaml:"similarity,omitempty"`
	Blend         float64 `json:"blend,omitempty" yaml:"blend,omitempty"`
	Despill       bool    `json:"despill,omitempty" yaml:"despill,omitempty"`
	DespillType   string  `json:"despill_type,omitempty" yaml:"despill_type,omitempty"` // "green", "blue"
	DespillExpand float64 `json:"despill_expand,omitempty" yaml:"despill_expand,omitempty"`
}

// KeyframesSpec configures dynamic position, scale, and opacity animation tracks.
type KeyframesSpec struct {
	Position []PositionKeyframeSpec `json:"position,omitempty" yaml:"position,omitempty"`
	Scale    []ScalarKeyframeSpec   `json:"scale,omitempty" yaml:"scale,omitempty"`
	Opacity  []ScalarKeyframeSpec   `json:"opacity,omitempty" yaml:"opacity,omitempty"`
}

// PositionKeyframeSpec represents a 2D spatial coordinate keyframe.
type PositionKeyframeSpec struct {
	Time   any    `json:"time" yaml:"time"`
	X      int    `json:"x" yaml:"x"`
	Y      int    `json:"y" yaml:"y"`
	Easing string `json:"easing,omitempty" yaml:"easing,omitempty"` // "linear", "ease_in_quad", "ease_out_quad", etc.
}

// ScalarKeyframeSpec represents a float keyframe (for scale or opacity).
type ScalarKeyframeSpec struct {
	Time   any     `json:"time" yaml:"time"`
	Value  float64 `json:"value" yaml:"value"`
	Easing string  `json:"easing,omitempty" yaml:"easing,omitempty"`
}
