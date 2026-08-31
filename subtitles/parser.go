package subtitles

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

var (
	// ErrInvalidSubtitleFormat indicates malformed subtitle structure.
	ErrInvalidSubtitleFormat = errors.New("subtitles: invalid subtitle file format")
)

// ParseSRT parses an SRT subtitle stream into a SubtitleTrack.
func ParseSRT(reader io.Reader) (*SubtitleTrack, error) {
	scanner := bufio.NewScanner(reader)
	track := NewSubtitleTrack()

	var currentCue *SubtitleCue
	var textLines []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = strings.TrimPrefix(line, "\uFEFF") // Strip UTF-8 BOM

		if line == "" {
			if currentCue != nil {
				currentCue.Text = strings.Join(textLines, "\n")
				track.AddCue(*currentCue)
				currentCue = nil
				textLines = nil
			}
			continue
		}

		if currentCue == nil {
			// Either cue index number or timestamp line
			if strings.Contains(line, "-->") {
				start, end, err := parseTimestampRange(line)
				if err != nil {
					return nil, err
				}
				currentCue = &SubtitleCue{
					Index:     len(track.Cues) + 1,
					StartTime: start,
					EndTime:   end,
				}
			} else {
				// Index number
				index, err := strconv.Atoi(line)
				if err == nil {
					currentCue = &SubtitleCue{
						Index: index,
					}
				}
			}
			continue
		}

		if currentCue.StartTime == 0 && currentCue.EndTime == 0 && strings.Contains(line, "-->") {
			start, end, err := parseTimestampRange(line)
			if err != nil {
				return nil, err
			}
			currentCue.StartTime = start
			currentCue.EndTime = end
			continue
		}

		textLines = append(textLines, line)
	}

	if currentCue != nil {
		currentCue.Text = strings.Join(textLines, "\n")
		track.AddCue(*currentCue)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("subtitles: error reading srt: %w", err)
	}

	return track, nil
}

// ParseVTT parses a WebVTT subtitle stream into a SubtitleTrack.
func ParseVTT(reader io.Reader) (*SubtitleTrack, error) {
	scanner := bufio.NewScanner(reader)
	track := NewSubtitleTrack()

	// Verify WEBVTT header
	if !scanner.Scan() {
		return nil, ErrInvalidSubtitleFormat
	}
	header := strings.TrimSpace(strings.TrimPrefix(scanner.Text(), "\uFEFF"))
	if !strings.HasPrefix(header, "WEBVTT") {
		return nil, fmt.Errorf("%w: missing WEBVTT header", ErrInvalidSubtitleFormat)
	}

	var currentCue *SubtitleCue
	var textLines []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			if currentCue != nil {
				currentCue.Text = strings.Join(textLines, "\n")
				track.AddCue(*currentCue)
				currentCue = nil
				textLines = nil
			}
			continue
		}

		if strings.Contains(line, "-->") {
			start, end, err := parseTimestampRange(line)
			if err != nil {
				return nil, err
			}
			currentCue = &SubtitleCue{
				Index:     len(track.Cues) + 1,
				StartTime: start,
				EndTime:   end,
			}
			continue
		}

		if currentCue != nil {
			textLines = append(textLines, line)
		}
	}

	if currentCue != nil {
		currentCue.Text = strings.Join(textLines, "\n")
		track.AddCue(*currentCue)
	}

	return track, nil
}

func parseTimestampRange(line string) (time.Duration, time.Duration, error) {
	parts := strings.Split(line, "-->")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("%w: malformed timestamp range %q", ErrInvalidSubtitleFormat, line)
	}

	startStr := strings.TrimSpace(parts[0])
	endStr := strings.TrimSpace(strings.Fields(parts[1])[0]) // Discard any trailing VTT cue settings

	start, err := ParseTimestamp(startStr)
	if err != nil {
		return 0, 0, err
	}

	end, err := ParseTimestamp(endStr)
	if err != nil {
		return 0, 0, err
	}

	return start, end, nil
}

// ParseTimestamp parses SRT ("00:01:23,456") or VTT ("01:23.456" / "00:01:23.456") timestamp strings.
func ParseTimestamp(timestampString string) (time.Duration, error) {
	timestampString = strings.ReplaceAll(strings.TrimSpace(timestampString), ",", ".")
	parts := strings.Split(timestampString, ":")

	var hours, minutes int
	var secondsFloat float64
	var err error

	switch len(parts) {
	case 2: // MM:SS.mmm
		minutes, err = strconv.Atoi(parts[0])
		if err != nil {
			return 0, fmt.Errorf("invalid minutes in %q: %w", timestampString, err)
		}
		secondsFloat, err = strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid seconds in %q: %w", timestampString, err)
		}
	case 3: // HH:MM:SS.mmm
		hours, err = strconv.Atoi(parts[0])
		if err != nil {
			return 0, fmt.Errorf("invalid hours in %q: %w", timestampString, err)
		}
		minutes, err = strconv.Atoi(parts[1])
		if err != nil {
			return 0, fmt.Errorf("invalid minutes in %q: %w", timestampString, err)
		}
		secondsFloat, err = strconv.ParseFloat(parts[2], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid seconds in %q: %w", timestampString, err)
		}
	default:
		return 0, fmt.Errorf("%w: unrecognized timestamp format %q", ErrInvalidSubtitleFormat, timestampString)
	}

	totalSeconds := float64(hours*3600+minutes*60) + secondsFloat
	return time.Duration(totalSeconds * float64(time.Second)), nil
}
