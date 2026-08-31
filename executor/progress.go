package executor

import (
	"bufio"
	"io"
	"strconv"
	"strings"
	"time"
)

// ProgressEvent contains real-time render telemetry reported by FFmpeg.
type ProgressEvent struct {
	Frame       int64
	FPS         float64
	Bitrate     string
	TotalSize   int64
	Time        time.Duration
	Speed       float64
	Progress    string // "continue" or "end"
	Percentage  float64
}

// ParseProgressStream reads key-value progress lines from r and fires onEvent when an event completes.
func ParseProgressStream(r io.Reader, totalDuration time.Duration, onEvent func(ProgressEvent)) error {
	scanner := bufio.NewScanner(r)
	current := ProgressEvent{}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		switch key {
		case "frame":
			if f, err := strconv.ParseInt(val, 10, 64); err == nil {
				current.Frame = f
			}
		case "fps":
			if fps, err := strconv.ParseFloat(val, 64); err == nil {
				current.FPS = fps
			}
		case "bitrate":
			current.Bitrate = val
		case "total_size":
			if size, err := strconv.ParseInt(val, 10, 64); err == nil {
				current.TotalSize = size
			}
		case "out_time_us":
			if us, err := strconv.ParseInt(val, 10, 64); err == nil {
				current.Time = time.Duration(us) * time.Microsecond
				if totalDuration > 0 {
					current.Percentage = (current.Time.Seconds() / totalDuration.Seconds()) * 100.0
					if current.Percentage > 100.0 {
						current.Percentage = 100.0
					}
				}
			}
		case "speed":
			val = strings.TrimSuffix(val, "x")
			if sp, err := strconv.ParseFloat(val, 64); err == nil {
				current.Speed = sp
			}
		case "progress":
			current.Progress = val
			if onEvent != nil {
				onEvent(current)
			}
			// Reset for next block
			current = ProgressEvent{}
		}
	}

	return scanner.Err()
}
