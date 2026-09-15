package timeline_test

import (
	"testing"

	"github.com/farshidrezaei/vidonex/timeline"
)

func TestTransitionsCatalog_Table(t *testing.T) {
	testCases := []struct {
		name           string
		transitionType timeline.TransitionType
		expectedValid  bool
	}{
		{
			name:           "Standard dissolve",
			transitionType: timeline.TransitionDissolve,
			expectedValid:  true,
		},
		{
			name:           "Wipe left",
			transitionType: timeline.TransitionWipeLeft,
			expectedValid:  true,
		},
		{
			name:           "Circle crop iris",
			transitionType: timeline.TransitionCircleCrop,
			expectedValid:  true,
		},
		{
			name:           "Zoom in",
			transitionType: timeline.TransitionZoomIn,
			expectedValid:  true,
		},
		{
			name:           "Audio acrossfade",
			transitionType: timeline.TransitionAcrossFade,
			expectedValid:  true,
		},
		{
			name:           "Unknown invalid transition",
			transitionType: timeline.TransitionType("invalid_warp_jump"),
			expectedValid:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			valid := timeline.IsValidTransitionType(tc.transitionType)
			if valid != tc.expectedValid {
				t.Errorf("expected IsValidTransitionType(%q)=%v, got %v", tc.transitionType, tc.expectedValid, valid)
			}
		})
	}

	allTransitions := timeline.AllVideoTransitions()
	if len(allTransitions) < 30 {
		t.Errorf("expected comprehensive transition catalog with 30+ transitions, got %d", len(allTransitions))
	}
}
