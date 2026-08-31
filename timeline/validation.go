package timeline

import (
	"errors"
	"fmt"
)

// Validate checks the semantic consistency of a Timeline:
// - Non-zero canvas dimensions
// - Positive frame rate (FPS)
// - Valid clip parameters (Duration > 0, Speed > 0, Opacity in [0, 1])
// - Transition overlap checks
func Validate(t *Timeline) error {
	if t == nil {
		return errors.New("timeline is nil")
	}

	var errs []error

	// Canvas dimensions check
	if t.Canvas.IsZero() {
		errs = append(errs, errors.New("timeline canvas width and height must be positive"))
	}
	if t.Canvas.Width%2 != 0 || t.Canvas.Height%2 != 0 {
		errs = append(errs, fmt.Errorf("timeline canvas dimensions (%dx%d) must be even numbers", t.Canvas.Width, t.Canvas.Height))
	}

	// FPS check
	if t.FPS.IsZero() || t.FPS.Num <= 0 || t.FPS.Den <= 0 {
		errs = append(errs, errors.New("timeline FPS must be a positive fraction"))
	}

	// Validate Tracks & Clips
	for trackIdx, track := range t.Tracks {
		if track == nil {
			errs = append(errs, fmt.Errorf("track index %d is nil", trackIdx))
			continue
		}

		for clipIdx, clip := range track.Clips {
			if clip == nil {
				errs = append(errs, fmt.Errorf("track %q clip index %d is nil", track.ID, clipIdx))
				continue
			}

			if clip.Source == "" {
				errs = append(errs, fmt.Errorf("track %q clip %q has empty source", track.ID, clip.ID))
			}
			if clip.Duration <= 0 {
				errs = append(errs, fmt.Errorf("track %q clip %q duration must be positive", track.ID, clip.ID))
			}
			if clip.Speed <= 0 {
				errs = append(errs, fmt.Errorf("track %q clip %q speed must be greater than 0", track.ID, clip.ID))
			}
			if clip.Opacity < 0 || clip.Opacity > 1.0 {
				errs = append(errs, fmt.Errorf("track %q clip %q opacity must be in [0.0, 1.0], got %f", track.ID, clip.ID, clip.Opacity))
			}
			if clip.Volume < 0 {
				errs = append(errs, fmt.Errorf("track %q clip %q volume must be non-negative, got %f", track.ID, clip.ID, clip.Volume))
			}

			for effIdx, eff := range clip.Effects {
				if eff == nil {
					errs = append(errs, fmt.Errorf("track %q clip %q effect index %d is nil", track.ID, clip.ID, effIdx))
					continue
				}
				if err := eff.Validate(); err != nil {
					errs = append(errs, fmt.Errorf("track %q clip %q effect %s invalid: %w", track.ID, clip.ID, eff.Type(), err))
				}
			}
		}

		for transIdx, trans := range track.Transitions {
			if trans == nil {
				errs = append(errs, fmt.Errorf("track %q transition index %d is nil", track.ID, transIdx))
				continue
			}
			if trans.Duration <= 0 {
				errs = append(errs, fmt.Errorf("track %q transition %q duration must be positive", track.ID, trans.ID))
			}
			if trans.ClipA == nil || trans.ClipB == nil {
				errs = append(errs, fmt.Errorf("track %q transition %q must have both ClipA and ClipB set", track.ID, trans.ID))
			}
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
