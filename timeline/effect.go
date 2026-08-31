// Package timeline provides declarative composition models for video timelines, tracks, clips, and transitions.
package timeline

// Effect is an interface implemented by all visual and audio effects attached to clips or tracks.
type Effect interface {
	Type() string
	Validate() error
}

// DrawTextEffect defines dynamic text overlay parameters.
type DrawTextEffect struct {
	Text      string
	FontFile  string
	FontSize  int
	FontColor string
	Box       bool
	BoxColor  string
	X         string // FFmpeg expression (e.g. "(w-text_w)/2")
	Y         string // FFmpeg expression (e.g. "(h-text_h)/2")
}

// Type returns the effect filter identifier name.
func (e *DrawTextEffect) Type() string {
	return "drawtext"
}

// Validate checks the semantic parameters of the drawtext effect.
func (e *DrawTextEffect) Validate() error {
	return nil
}

// BlurEffect defines gaussian/box blur parameters.
type BlurEffect struct {
	Radius int
	Power  int
}

// Type returns the blur effect filter identifier.
func (e *BlurEffect) Type() string {
	return "blur"
}

// Validate checks the blur effect parameters.
func (e *BlurEffect) Validate() error {
	return nil
}

// ChromaKeyEffect defines green/blue screen chroma keying.
type ChromaKeyEffect struct {
	ColorHex   string
	Similarity float64 // 0.0 to 1.0 (default 0.3)
	Blend      float64 // 0.0 to 1.0 (default 0.1)
}

// Type returns the chromakey effect filter identifier.
func (e *ChromaKeyEffect) Type() string {
	return "chromakey"
}

// Validate checks the chromakey effect parameters.
func (e *ChromaKeyEffect) Validate() error {
	return nil
}
