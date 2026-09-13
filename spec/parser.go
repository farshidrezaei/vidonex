package spec

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/farshidrezaei/vidonex/animation"
	"github.com/farshidrezaei/vidonex/types"
	"gopkg.in/yaml.v3"
)

// ParseJSON parses a JSON stream into a VideoSpec.
func ParseJSON(reader io.Reader) (*VideoSpec, error) {
	var specification VideoSpec
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&specification); err != nil {
		return nil, fmt.Errorf("spec.ParseJSON: failed decoding JSON: %w", err)
	}
	return &specification, nil
}

// ParseYAML parses a YAML stream into a VideoSpec.
func ParseYAML(reader io.Reader) (*VideoSpec, error) {
	var specification VideoSpec
	decoder := yaml.NewDecoder(reader)
	if err := decoder.Decode(&specification); err != nil {
		return nil, fmt.Errorf("spec.ParseYAML: failed decoding YAML: %w", err)
	}
	return &specification, nil
}

// ParseFile loads and parses a spec file (.json, .yaml, or .yml) from disk.
func ParseFile(filePath string) (*VideoSpec, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("spec.ParseFile: failed reading %q: %w", filePath, err)
	}

	extension := strings.ToLower(filepath.Ext(filePath))
	switch extension {
	case ".json":
		return ParseJSON(strings.NewReader(string(fileBytes)))
	case ".yaml", ".yml":
		return ParseYAML(strings.NewReader(string(fileBytes)))
	default:
		// Attempt YAML first (since YAML is a superset of JSON), fallback to JSON
		if spec, yamlErr := ParseYAML(strings.NewReader(string(fileBytes))); yamlErr == nil {
			return spec, nil
		}
		return ParseJSON(strings.NewReader(string(fileBytes)))
	}
}

// ParseDurationValue converts numbers, duration strings ("5s", "500ms"), or timestamps ("00:00:05.500") into time.Duration.
func ParseDurationValue(raw any) (time.Duration, error) {
	if raw == nil {
		return 0, nil
	}

	switch value := raw.(type) {
	case int:
		return time.Duration(value) * time.Second, nil
	case int64:
		return time.Duration(value) * time.Second, nil
	case float64:
		return time.Duration(value * float64(time.Second)), nil
	case string:
		clean := strings.TrimSpace(value)
		if clean == "" || clean == "0" {
			return 0, nil
		}

		// Check for standard Go duration string e.g. "5s", "500ms", "2.5s"
		if parsedDuration, err := time.ParseDuration(clean); err == nil {
			return parsedDuration, nil
		}

		// Check for plain numeric string e.g. "5.5"
		if parsedSeconds, err := strconv.ParseFloat(clean, 64); err == nil {
			return time.Duration(parsedSeconds * float64(time.Second)), nil
		}

		// Check for timestamp string: "HH:MM:SS" or "HH:MM:SS.mmm" or "MM:SS.mmm"
		parts := strings.Split(clean, ":")
		if len(parts) == 3 {
			hours, errH := strconv.Atoi(parts[0])
			minutes, errM := strconv.Atoi(parts[1])
			seconds, errS := strconv.ParseFloat(parts[2], 64)
			if errH == nil && errM == nil && errS == nil {
				totalSeconds := float64(hours*3600+minutes*60) + seconds
				return time.Duration(totalSeconds * float64(time.Second)), nil
			}
		} else if len(parts) == 2 {
			minutes, errM := strconv.Atoi(parts[0])
			seconds, errS := strconv.ParseFloat(parts[1], 64)
			if errM == nil && errS == nil {
				totalSeconds := float64(minutes*60) + seconds
				return time.Duration(totalSeconds * float64(time.Second)), nil
			}
		}

		return 0, fmt.Errorf("spec.ParseDurationValue: unsupported duration format %q", clean)
	default:
		return 0, fmt.Errorf("spec.ParseDurationValue: invalid duration type %T", raw)
	}
}

// ParseColorValue converts color strings (#RRGGBB, #RRGGBBAA, or named colors) into types.Color.
func ParseColorValue(raw string) (types.Color, error) {
	trimmed := strings.TrimSpace(strings.ToLower(raw))
	if trimmed == "" {
		return types.ColorTransparent, nil
	}

	switch trimmed {
	case "black":
		return types.ColorBlack, nil
	case "white":
		return types.ColorWhite, nil
	case "red":
		return types.ColorRed, nil
	case "green":
		return types.ColorGreen, nil
	case "blue":
		return types.ColorBlue, nil
	case "yellow":
		return types.ColorYellow, nil
	case "transparent":
		return types.ColorTransparent, nil
	default:
		return types.Hex(raw)
	}
}

// ParseAlignmentValue converts alignment strings into types.Alignment.
func ParseAlignmentValue(raw string) types.Alignment {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "top_left", "top-left", "topleft":
		return types.AlignTopLeft
	case "top_center", "top-center", "topcenter", "top":
		return types.AlignTopCenter
	case "top_right", "top-right", "topright":
		return types.AlignTopRight
	case "center", "middle":
		return types.AlignCenter
	case "bottom_left", "bottom-left", "bottomleft":
		return types.AlignBottomLeft
	case "bottom_right", "bottom-right", "bottomright":
		return types.AlignBottomRight
	case "bottom_center", "bottom-center", "bottomcenter", "bottom":
		fallthrough
	default:
		return types.AlignBottomCenter
	}
}

// ParseEasingValue converts easing strings into animation.EasingFunction.
func ParseEasingValue(raw string) animation.EasingFunction {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "ease_in_quad", "easeinquad":
		return animation.EasingEaseInQuad
	case "ease_out_quad", "easeoutquad":
		return animation.EasingEaseOutQuad
	case "ease_in_out_quad", "easeinoutquad":
		return animation.EasingEaseInOutQuad
	case "ease_in_cubic", "easeincubic":
		return animation.EasingEaseInCubic
	case "ease_out_cubic", "easeoutcubic":
		return animation.EasingEaseOutCubic
	case "ease_in_out_cubic", "easeinoutcubic":
		return animation.EasingEaseInOutCubic
	case "ease_in_sine", "easeinsine":
		return animation.EasingEaseInSine
	case "ease_out_sine", "easeoutsine":
		return animation.EasingEaseOutSine
	case "ease_in_out_sine", "easeinoutsine":
		return animation.EasingEaseInOutSine
	case "linear":
		fallthrough
	default:
		return animation.EasingLinear
	}
}
