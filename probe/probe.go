// Package probe provides automated media inspection, stream metadata extraction, and caching using ffprobe.
package probe

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/farshidrezaei/vidonex/types"
)

var (
	// ErrProbeFailed indicates ffprobe failed to execute or returned an error code.
	ErrProbeFailed = errors.New("probe: media probe failed")
	// ErrNoMediaStreams indicates the probed media file contains no audio or video streams.
	ErrNoMediaStreams = errors.New("probe: no valid media streams found")
)

// VideoStreamMetadata describes the technical parameters of a video stream.
type VideoStreamMetadata struct {
	Width       int
	Height      int
	FPS         types.Rational
	Codec       string
	PixelFormat string
	SAR         types.Rational
	DAR         types.Rational
	Duration    time.Duration
	FrameCount  int64
}

// AspectRatio returns the video display aspect ratio (Width:Height).
func (videoStream *VideoStreamMetadata) AspectRatio() types.Rational {
	if videoStream.Width <= 0 || videoStream.Height <= 0 {
		return types.NewRational(16, 9)
	}
	return types.NewRational(int64(videoStream.Width), int64(videoStream.Height))
}

// Size returns the video frame dimensions as a types.Size.
func (videoStream *VideoStreamMetadata) Size() types.Size {
	return types.NewSize(videoStream.Width, videoStream.Height)
}

// AudioStreamMetadata describes the technical parameters of an audio stream.
type AudioStreamMetadata struct {
	Codec         string
	SampleRate    int
	Channels      int
	ChannelLayout string
	BitRate       int64
	Duration      time.Duration
}

// MediaMetadata encapsulates the combined technical metadata of a media container.
type MediaMetadata struct {
	Path        string
	FormatName  string
	Duration    time.Duration
	BitRate     int64
	Size        int64
	VideoStream *VideoStreamMetadata
	AudioStream *AudioStreamMetadata
}

// HasVideo reports whether the media has a video stream.
func (metadata *MediaMetadata) HasVideo() bool {
	return metadata.VideoStream != nil
}

// HasAudio reports whether the media has an audio stream.
func (metadata *MediaMetadata) HasAudio() bool {
	return metadata.AudioStream != nil
}

// MediaProber abstracts media inspection operations.
type MediaProber interface {
	Probe(ctx context.Context, mediaPath string) (*MediaMetadata, error)
}

// FFprobeProber inspects media files using the ffprobe CLI binary.
type FFprobeProber struct {
	binaryPath string
}

// NewFFprobeProber creates a new FFprobeProber.
func NewFFprobeProber(binaryPath string) *FFprobeProber {
	if binaryPath == "" {
		binaryPath = "ffprobe"
	}
	return &FFprobeProber{
		binaryPath: binaryPath,
	}
}

// Probe invokes ffprobe with JSON output formatting and parses stream metadata.
func (prober *FFprobeProber) Probe(ctx context.Context, mediaPath string) (*MediaMetadata, error) {
	arguments := []string{
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		mediaPath,
	}

	command := exec.CommandContext(ctx, prober.binaryPath, arguments...)
	outputBytes, err := command.Output()
	if err != nil {
		return nil, fmt.Errorf("%w: failed executing %s on %q: %v", ErrProbeFailed, prober.binaryPath, mediaPath, err)
	}

	return ParseFFprobeJSON(outputBytes, mediaPath)
}

// CachedProber wraps a MediaProber with a thread-safe in-memory cache.
type CachedProber struct {
	underlyingProber MediaProber
	cacheMutex       sync.RWMutex
	cache            map[string]*MediaMetadata
}

// NewCachedProber creates a thread-safe caching wrapper around a MediaProber.
func NewCachedProber(underlyingProber MediaProber) *CachedProber {
	return &CachedProber{
		underlyingProber: underlyingProber,
		cache:            make(map[string]*MediaMetadata),
	}
}

// Probe checks the in-memory cache before delegating to the underlying prober.
func (cachedProber *CachedProber) Probe(ctx context.Context, mediaPath string) (*MediaMetadata, error) {
	cachedProber.cacheMutex.RLock()
	if cachedMetadata, exists := cachedProber.cache[mediaPath]; exists {
		cachedProber.cacheMutex.RUnlock()
		return cachedMetadata, nil
	}
	cachedProber.cacheMutex.RUnlock()

	metadata, err := cachedProber.underlyingProber.Probe(ctx, mediaPath)
	if err != nil {
		return nil, err
	}

	cachedProber.cacheMutex.Lock()
	cachedProber.cache[mediaPath] = metadata
	cachedProber.cacheMutex.Unlock()

	return metadata, nil
}

// MockProber provides deterministic mock metadata for unit testing without external binaries.
type MockProber struct {
	mutex        sync.RWMutex
	metadataMap  map[string]*MediaMetadata
	defaultMeta  *MediaMetadata
	simulateErr  error
	probedPaths  []string
}

// NewMockProber creates an initialized MockProber.
func NewMockProber() *MockProber {
	return &MockProber{
		metadataMap: make(map[string]*MediaMetadata),
		defaultMeta: &MediaMetadata{
			FormatName: "mp4",
			Duration:   10 * time.Second,
			VideoStream: &VideoStreamMetadata{
				Width:       1920,
				Height:      1080,
				FPS:         types.FPS30,
				Codec:       "h264",
				PixelFormat: "yuv420p",
				Duration:    10 * time.Second,
			},
			AudioStream: &AudioStreamMetadata{
				Codec:         "aac",
				SampleRate:    48000,
				Channels:      2,
				ChannelLayout: "stereo",
				Duration:      10 * time.Second,
			},
		},
		probedPaths: make([]string, 0),
	}
}

// SetMetadata registers specific metadata for a given file path.
func (mockProber *MockProber) SetMetadata(mediaPath string, metadata *MediaMetadata) *MockProber {
	mockProber.mutex.Lock()
	defer mockProber.mutex.Unlock()
	mockProber.metadataMap[mediaPath] = metadata
	return mockProber
}

// SetSimulatedError sets an error to be returned by Probe.
func (mockProber *MockProber) SetSimulatedError(err error) *MockProber {
	mockProber.mutex.Lock()
	defer mockProber.mutex.Unlock()
	mockProber.simulateErr = err
	return mockProber
}

// Probe returns registered mock metadata or the default mock metadata.
func (mockProber *MockProber) Probe(_ context.Context, mediaPath string) (*MediaMetadata, error) {
	mockProber.mutex.Lock()
	defer mockProber.mutex.Unlock()

	mockProber.probedPaths = append(mockProber.probedPaths, mediaPath)

	if mockProber.simulateErr != nil {
		return nil, mockProber.simulateErr
	}

	if metadata, exists := mockProber.metadataMap[mediaPath]; exists {
		return metadata, nil
	}

	// Return a copy of default metadata customized with the requested path
	copied := *mockProber.defaultMeta
	copied.Path = mediaPath
	return &copied, nil
}

// ProbedPaths returns a slice of all probed paths recorded by the mock.
func (mockProber *MockProber) ProbedPaths() []string {
	mockProber.mutex.RLock()
	defer mockProber.mutex.RUnlock()
	return append([]string(nil), mockProber.probedPaths...)
}

// rawFFprobeOutput represents ffprobe's JSON output structure.
type rawFFprobeOutput struct {
	Streams []rawStream `json:"streams"`
	Format  rawFormat   `json:"format"`
}

type rawStream struct {
	CodecType      string `json:"codec_type"`
	CodecName      string `json:"codec_name"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
	PixFmt         string `json:"pix_fmt"`
	RFrameRate     string `json:"r_frame_rate"`
	AvgFrameRate   string `json:"avg_frame_rate"`
	SampleAspectRatio string `json:"sample_aspect_ratio"`
	DisplayAspectRatio string `json:"display_aspect_ratio"`
	Duration       string `json:"duration"`
	NbFrames       string `json:"nb_frames"`
	SampleRate     string `json:"sample_rate"`
	Channels       int    `json:"channels"`
	ChannelLayout  string `json:"channel_layout"`
	BitRate        string `json:"bit_rate"`
}

type rawFormat struct {
	Filename   string `json:"filename"`
	FormatName string `json:"format_name"`
	Duration   string `json:"duration"`
	Size       string `json:"size"`
	BitRate    string `json:"bit_rate"`
}

// ParseFFprobeJSON parses JSON output from ffprobe.
func ParseFFprobeJSON(jsonData []byte, mediaPath string) (*MediaMetadata, error) {
	var rawOutput rawFFprobeOutput
	if err := json.Unmarshal(jsonData, &rawOutput); err != nil {
		return nil, fmt.Errorf("probe: failed to parse ffprobe json: %w", err)
	}

	containerDuration := parseDurationSeconds(rawOutput.Format.Duration)
	containerSize, _ := strconv.ParseInt(rawOutput.Format.Size, 10, 64)
	containerBitRate, _ := strconv.ParseInt(rawOutput.Format.BitRate, 10, 64)

	result := &MediaMetadata{
		Path:       mediaPath,
		FormatName: rawOutput.Format.FormatName,
		Duration:   containerDuration,
		Size:       containerSize,
		BitRate:    containerBitRate,
	}

	for _, stream := range rawOutput.Streams {
		switch stream.CodecType {
		case "video":
			if result.VideoStream == nil {
				frameRate := parseRationalString(stream.AvgFrameRate)
				if frameRate.IsZero() {
					frameRate = parseRationalString(stream.RFrameRate)
				}
				if frameRate.IsZero() {
					frameRate = types.FPS30
				}

				streamDuration := parseDurationSeconds(stream.Duration)
				if streamDuration == 0 {
					streamDuration = containerDuration
				}

				frameCount, _ := strconv.ParseInt(stream.NbFrames, 10, 64)
				sar := parseRationalString(stream.SampleAspectRatio)
				dar := parseRationalString(stream.DisplayAspectRatio)

				result.VideoStream = &VideoStreamMetadata{
					Width:       stream.Width,
					Height:      stream.Height,
					FPS:         frameRate,
					Codec:       stream.CodecName,
					PixelFormat: stream.PixFmt,
					SAR:         sar,
					DAR:         dar,
					Duration:    streamDuration,
					FrameCount:  frameCount,
				}
			}

		case "audio":
			if result.AudioStream == nil {
				sampleRate, _ := strconv.Atoi(stream.SampleRate)
				bitRate, _ := strconv.ParseInt(stream.BitRate, 10, 64)
				streamDuration := parseDurationSeconds(stream.Duration)
				if streamDuration == 0 {
					streamDuration = containerDuration
				}

				result.AudioStream = &AudioStreamMetadata{
					Codec:         stream.CodecName,
					SampleRate:    sampleRate,
					Channels:      stream.Channels,
					ChannelLayout: stream.ChannelLayout,
					BitRate:       bitRate,
					Duration:      streamDuration,
				}
			}
		}
	}

	if result.VideoStream == nil && result.AudioStream == nil {
		return nil, ErrNoMediaStreams
	}

	return result, nil
}

func parseDurationSeconds(secondsString string) time.Duration {
	if secondsString == "" || secondsString == "N/A" {
		return 0
	}
	secondsFloat, err := strconv.ParseFloat(secondsString, 64)
	if err != nil {
		return 0
	}
	return time.Duration(secondsFloat * float64(time.Second))
}

func parseRationalString(fractionString string) types.Rational {
	if fractionString == "" || fractionString == "0/0" || fractionString == "N/A" {
		return types.Rational{}
	}
	parts := strings.Split(fractionString, "/")
	if len(parts) == 1 {
		num, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return types.Rational{}
		}
		return types.NewRational(num, 1)
	}
	if len(parts) == 2 {
		num, errNum := strconv.ParseInt(parts[0], 10, 64)
		den, errDen := strconv.ParseInt(parts[1], 10, 64)
		if errNum != nil || errDen != nil || den == 0 {
			return types.Rational{}
		}
		return types.NewRational(num, den)
	}
	return types.Rational{}
}
