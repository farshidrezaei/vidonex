package subtitles_test

import (
	"strings"
	"testing"
	"time"

	"github.com/farshidrezaei/vidonyx/filtergraph"
	"github.com/farshidrezaei/vidonyx/subtitles"
	"github.com/farshidrezaei/vidonyx/types"
)

func TestParseSRTTable(t *testing.T) {
	sampleSRT := `1
00:00:01,000 --> 00:00:04,500
Welcome to Vidonyx!
Declarative Video Composition in Go.

2
00:00:05,000 --> 00:00:08,250
High Performance & Fast Rendering.
`

	tests := []struct {
		name      string
		inputSRT  string
		wantCues  int
		checkCues func(t *testing.T, track *subtitles.SubtitleTrack)
	}{
		{
			name:     "valid 2-cue srt",
			inputSRT: sampleSRT,
			wantCues: 2,
			checkCues: func(t *testing.T, track *subtitles.SubtitleTrack) {
				cue1 := track.Cues[0]
				if cue1.StartTime != 1*time.Second || cue1.EndTime != 4500*time.Millisecond {
					t.Errorf("Cue 1 timing = %v to %v, want 1s to 4.5s", cue1.StartTime, cue1.EndTime)
				}
				if !strings.Contains(cue1.Text, "Welcome to Vidonyx!") {
					t.Errorf("Cue 1 text does not contain expected string: %s", cue1.Text)
				}

				cue2 := track.Cues[1]
				if cue2.StartTime != 5*time.Second || cue2.EndTime != 8250*time.Millisecond {
					t.Errorf("Cue 2 timing = %v to %v, want 5s to 8.25s", cue2.StartTime, cue2.EndTime)
				}
			},
		},
		{
			name:      "empty srt",
			inputSRT:  "",
			wantCues:  0,
			checkCues: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			track, err := subtitles.ParseSRT(strings.NewReader(tt.inputSRT))
			if err != nil {
				t.Fatalf("ParseSRT failed: %v", err)
			}
			if len(track.Cues) != tt.wantCues {
				t.Fatalf("Cues count = %d, want %d", len(track.Cues), tt.wantCues)
			}
			if tt.checkCues != nil {
				tt.checkCues(t, track)
			}
		})
	}
}

func TestParseVTTTable(t *testing.T) {
	sampleVTT := `WEBVTT

00:01.000 --> 00:04.000
Caption line one

00:05.000 --> 00:09.500 position:50% line:80%
Caption line two with VTT settings
`

	track, err := subtitles.ParseVTT(strings.NewReader(sampleVTT))
	if err != nil {
		t.Fatalf("ParseVTT failed: %v", err)
	}

	if len(track.Cues) != 2 {
		t.Fatalf("Expected 2 cues, got %d", len(track.Cues))
	}

	if track.Cues[0].StartTime != 1*time.Second || track.Cues[0].EndTime != 4*time.Second {
		t.Errorf("Cue 0 timing mismatch: %v - %v", track.Cues[0].StartTime, track.Cues[0].EndTime)
	}
}

func TestAttachSubtitlesCompiler(t *testing.T) {
	graph := filtergraph.NewGraph()
	dummyIn := &filtergraph.Pad{
		ID:         "video_in",
		StreamType: filtergraph.StreamTypeVideo,
		IsInput:    false,
	}

	track := subtitles.NewSubtitleTrack()
	track.AddCue(subtitles.SubtitleCue{
		Index:     1,
		StartTime: 1 * time.Second,
		EndTime:   3 * time.Second,
		Text:      "Hello World!",
	})

	outPad, err := subtitles.AttachSubtitles(graph, dummyIn, track, types.Res1080p)
	if err != nil {
		t.Fatalf("AttachSubtitles failed: %v", err)
	}

	if outPad == nil {
		t.Fatalf("expected non-nil output pad")
	}

	filterComplex, err := graph.FormattedFilterComplex()
	if err != nil {
		t.Fatalf("FormattedFilterComplex failed: %v", err)
	}

	if !strings.Contains(filterComplex, "drawtext") || !strings.Contains(filterComplex, "Hello World!") {
		t.Errorf("filter complex does not contain expected drawtext: %s", filterComplex)
	}
}
