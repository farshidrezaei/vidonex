package compiler_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/farshidrezaei/vidonex/compiler"
	"github.com/farshidrezaei/vidonex/filtergraph"
	"github.com/farshidrezaei/vidonex/timeline"
	"github.com/farshidrezaei/vidonex/types"
)

func TestChainTrackTransitions_Table(t *testing.T) {
	t.Parallel()

	testScenarios := []struct {
		name                     string
		clips                    []*timeline.Clip
		transitions              func(clips []*timeline.Clip) []*timeline.Transition
		expectedChainedDuration  time.Duration
		expectedNodeSubstrings   []string
		unexpectedNodeSubstrings []string
	}{
		{
			name: "Two clips with dissolve transition",
			clips: []*timeline.Clip{
				timeline.NewClip("clip_1", "assets/clip1.mp4", 0, 5*time.Second),
				timeline.NewClip("clip_2", "assets/clip2.mp4", 5*time.Second, 5*time.Second),
			},
			transitions: func(clips []*timeline.Clip) []*timeline.Transition {
				return []*timeline.Transition{
					timeline.NewTransition("trans_1_2", timeline.TransitionDissolve, 1*time.Second, clips[0], clips[1]),
				}
			},
			expectedChainedDuration: 9500 * time.Millisecond, // 5 + 5 - 0.5 = 9.5s
			expectedNodeSubstrings: []string{
				"xfade",
				"transition=dissolve",
			},
		},
		{
			name: "Three clips with two sequential transitions",
			clips: []*timeline.Clip{
				timeline.NewClip("scene_a", "assets/a.mp4", 0, 5*time.Second),
				timeline.NewClip("scene_b", "assets/b.mp4", 4*time.Second, 6*time.Second),
				timeline.NewClip("scene_c", "assets/c.mp4", 9*time.Second, 4*time.Second),
			},
			transitions: func(clips []*timeline.Clip) []*timeline.Transition {
				return []*timeline.Transition{
					timeline.NewTransition("trans_ab", timeline.TransitionFade, 1*time.Second, clips[0], clips[1]),
					timeline.NewTransition("trans_bc", timeline.TransitionWipeLeft, 1*time.Second, clips[1], clips[2]),
				}
			},
			expectedChainedDuration: 14 * time.Second, // 5 + 6 - 0.5 + 4 - 0.5 = 14s
			expectedNodeSubstrings: []string{
				"transition=fade",
				"transition=wipeleft",
			},
		},
		{
			name: "Three clips where first pair is a cut and second pair has a transition",
			clips: []*timeline.Clip{
				timeline.NewClip("part_1", "assets/p1.mp4", 0, 4*time.Second),
				timeline.NewClip("part_2", "assets/p2.mp4", 4*time.Second, 5*time.Second),
				timeline.NewClip("part_3", "assets/p3.mp4", 9*time.Second, 6*time.Second),
			},
			transitions: func(clips []*timeline.Clip) []*timeline.Transition {
				// Only transition between part_2 and part_3
				return []*timeline.Transition{
					timeline.NewTransition("trans_2_3", timeline.TransitionCircleCrop, 2*time.Second, clips[1], clips[2]),
				}
			},
			expectedChainedDuration: 14 * time.Second, // 4 + 5 + 6 - 1 = 14s
			expectedNodeSubstrings: []string{
				"concat=n=2:v=1:a=0",
				"transition=circlecrop",
			},
		},
		{
			name: "Four clips with alternating transition and cut",
			clips: []*timeline.Clip{
				timeline.NewClip("clip_a", "assets/a.mp4", 0, 3*time.Second),
				timeline.NewClip("clip_b", "assets/b.mp4", 3*time.Second, 4*time.Second),
				timeline.NewClip("clip_c", "assets/c.mp4", 7*time.Second, 5*time.Second),
				timeline.NewClip("clip_d", "assets/d.mp4", 12*time.Second, 4*time.Second),
			},
			transitions: func(clips []*timeline.Clip) []*timeline.Transition {
				return []*timeline.Transition{
					timeline.NewTransition("trans_ab", timeline.TransitionSlideLeft, 1*time.Second, clips[0], clips[1]),
					// cut between b and c
					timeline.NewTransition("trans_cd", timeline.TransitionSlideRight, 1*time.Second, clips[2], clips[3]),
				}
			},
			expectedChainedDuration: 15 * time.Second, // 3 + 4 - 0.5 + 5 + 4 - 0.5 = 15s
			expectedNodeSubstrings: []string{
				"transition=slideleft",
				"concat=n=2:v=1:a=0",
				"transition=slideright",
			},
		},
	}

	for _, scenario := range testScenarios {
		scenario := scenario
		t.Run(scenario.name, func(t *testing.T) {
			t.Parallel()

			graph := filtergraph.NewGraph()
			pads := make([]*filtergraph.Pad, len(scenario.clips))
			for index, clip := range scenario.clips {
				node := graph.NewNode(fmt.Sprintf("source_node_%s", clip.ID), "testsrc")
				pads[index] = node.AddOutput(graph.NextPadID("test_pad"), filtergraph.StreamTypeVideo)
			}

			transitions := scenario.transitions(scenario.clips)
			chainedPad, chainedDuration, err := compiler.ChainTrackTransitionsVideo(graph, pads, transitions, scenario.clips)
			if err != nil {
				t.Fatalf("unexpected error chaining transitions: %v", err)
			}

			if chainedPad == nil {
				t.Fatal("expected non-nil chained pad")
			}

			if chainedDuration != scenario.expectedChainedDuration {
				t.Errorf("expected chained duration %v, got %v", scenario.expectedChainedDuration, chainedDuration)
			}

			filterComplexString, err := graph.FormattedFilterComplex()
			if err != nil {
				t.Fatalf("failed to format filter complex: %v", err)
			}
			for _, expectedSubstr := range scenario.expectedNodeSubstrings {
				if !strings.Contains(filterComplexString, expectedSubstr) {
					t.Errorf("expected filter string to contain %q, but got:\n%s", expectedSubstr, filterComplexString)
				}
			}
			for _, unexpectedSubstr := range scenario.unexpectedNodeSubstrings {
				if strings.Contains(filterComplexString, unexpectedSubstr) {
					t.Errorf("expected filter string NOT to contain %q, but got:\n%s", unexpectedSubstr, filterComplexString)
				}
			}
		})
	}
}

func TestFullTimelineCompilationWithSparseTransitions(t *testing.T) {
	t.Parallel()

	compositionTimeline := timeline.New(
		timeline.WithCanvas(types.Res1080p),
		timeline.WithFPS(types.FPS30),
	)

	track := timeline.NewTrack("main_track", timeline.TrackKindVideo).SetZIndex(0)
	clipA := timeline.NewClip("clip_a", "assets/a.mp4", 0, 5*time.Second)
	clipB := timeline.NewClip("clip_b", "assets/b.mp4", 5*time.Second, 5*time.Second)
	clipC := timeline.NewClip("clip_c", "assets/c.mp4", 10*time.Second, 5*time.Second)

	// Sparse transition: Only between B and C
	transitionBC := timeline.NewTransition("trans_bc", timeline.TransitionDissolve, 1*time.Second, clipB, clipC)

	track.AddClip(clipA, clipB, clipC)
	track.AddTransition(transitionBC)
	compositionTimeline.AddTrack(track)

	compilerInstance := compiler.New(nil)
	result, err := compilerInstance.Compile(compositionTimeline, "output_sparse.mp4")
	if err != nil {
		t.Fatalf("unexpected compilation failure: %v", err)
	}

	if result.FilterComplex == "" {
		t.Fatal("expected non-empty filter complex")
	}

	// Verify that internal node IDs exist in the DAG
	if _, exists := result.Graph.GetNode("vconcat_clip_a_clip_b"); !exists {
		t.Error("expected DAG node 'vconcat_clip_a_clip_b' to exist")
	}

	if _, exists := result.Graph.GetNode("xfade_trans_clip_b_clip_c"); !exists {
		t.Error("expected DAG node 'xfade_trans_clip_b_clip_c' to exist")
	}

	if !strings.Contains(result.FilterComplex, "concat=n=2:v=1:a=0") {
		t.Errorf("expected filter complex to contain cut concat node 'concat=n=2:v=1:a=0', got:\n%s", result.FilterComplex)
	}

	if !strings.Contains(result.FilterComplex, "xfade=transition=dissolve") {
		t.Errorf("expected filter complex to contain xfade node 'xfade=transition=dissolve', got:\n%s", result.FilterComplex)
	}
}
