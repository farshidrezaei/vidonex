package e2e_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/farshidrezaei/vidonyx/probe"
)

func TestCLI_EndToEndSuiteTable(t *testing.T) {
	temporaryDirectory := t.TempDir()
	assets := generateSyntheticMedia(t, temporaryDirectory)
	mediaProber := probe.NewFFprobeProber("ffprobe")

	// Compile CLI binary into temporary directory
	cliBinaryPath := filepath.Join(temporaryDirectory, "vidonyx_cli")
	compileCmd := exec.Command("go", "build", "-o", cliBinaryPath, "../../cmd/vidonyx")
	if compileOut, err := compileCmd.CombinedOutput(); err != nil {
		t.Fatalf("failed compiling vidonyx CLI: %v\nOutput: %s", err, string(compileOut))
	}

	// 1. Create a sample YAML project file in a dedicated subdirectory to test relative path resolution
	projectDir := filepath.Join(temporaryDirectory, "vlog_project")
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("failed creating project directory: %v", err)
	}

	// Copy media into relative subfolders
	mediaDir := filepath.Join(projectDir, "media")
	audioDir := filepath.Join(projectDir, "audio")
	if err := os.MkdirAll(mediaDir, 0755); err != nil {
		t.Fatalf("failed creating media directory: %v", err)
	}
	if err := os.MkdirAll(audioDir, 0755); err != nil {
		t.Fatalf("failed creating audio directory: %v", err)
	}

	// Copy assets into relative folder structure
	copyFile(t, assets.videoAlpha, filepath.Join(mediaDir, "main.mp4"))
	copyFile(t, assets.audioVoice, filepath.Join(audioDir, "speech.mp3"))
	copyFile(t, assets.audioSoundtrack, filepath.Join(audioDir, "bgm.mp3"))

	yamlSpecPath := filepath.Join(projectDir, "vidonyx.yaml")
	yamlContent := `
version: "1.0"
canvas:
  preset: "720p"
  background_color: "black"
fps: "30"
output: "rendered_vlog.mp4"
tracks:
  - id: "video_layer"
    kind: "video"
    clips:
      - id: "clip_main"
        source: "./media/main.mp4"
        start: "0s"
        duration: "3s"
        fade_in: "500ms"
  - id: "voice_layer"
    kind: "audio"
    clips:
      - id: "speech_clip"
        source: "./audio/speech.mp3"
        start: "0.5s"
        duration: "2s"
  - id: "music_layer"
    kind: "audio"
    duck_under: "voice_layer"
    ducking:
      threshold: 0.1
      ratio: 4.0
    clips:
      - id: "bgm_clip"
        source: "./audio/bgm.mp3"
        start: "0s"
        duration: "3s"
        volume: 0.6
  - id: "subtitles_layer"
    kind: "subtitle"
    subtitles:
      srt_content: |
        1
        00:00:00,500 --> 00:00:02,500
        Vidonyx CLI YAML Automation!
      style:
        font_size: 32
        color: "yellow"
        box: true
        box_color: "black"
        alignment: "bottom-center"
        margin_bottom: 40
  - id: "wave_layer"
    kind: "waveform"
    waveform:
      source_audio: "voice_layer"
      mode: "p2p"
      color: "#00F0B4"
      scale: "sqrt"
      width: 600
      height: 100
      position:
        x: 100
        y: 550
`
	if err := os.WriteFile(yamlSpecPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed writing yaml spec: %v", err)
	}

	// 2. Create a sample JSON project file
	jsonSpecPath := filepath.Join(projectDir, "vidonyx.json")
	jsonContent := `{
  "version": "1.0",
  "preset": "tiktok_1080p60",
  "output": "rendered_tiktok.mp4",
  "tracks": [
    {
      "id": "main_video",
      "kind": "video",
      "clips": [
        {
          "id": "clip1",
          "source": "./media/main.mp4",
          "start": 0,
          "duration": 2.0
        }
      ]
    }
  ]
}`
	if err := os.WriteFile(jsonSpecPath, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("failed writing json spec: %v", err)
	}

	tests := []struct {
		name                 string
		arguments            []string
		expectedOutputStdout string
		verifyOutputFile     string
		expectedWidth        int
		expectedHeight       int
	}{
		{
			name:                 "cli_version",
			arguments:            []string{"version"},
			expectedOutputStdout: "vidonyx engine version 1.0.0",
		},
		{
			name:                 "cli_validate_yaml",
			arguments:            []string{"validate", yamlSpecPath},
			expectedOutputStdout: "is valid!",
		},
		{
			name:                 "cli_validate_json",
			arguments:            []string{"validate", jsonSpecPath},
			expectedOutputStdout: "is valid!",
		},
		{
			name:                 "cli_graph_mermaid",
			arguments:            []string{"graph", yamlSpecPath, "--format", "mermaid"},
			expectedOutputStdout: "graph LR",
		},
		{
			name:                 "cli_probe_json",
			arguments:            []string{"probe", assets.videoAlpha, "--json"},
			expectedOutputStdout: `"VideoStream":`,
		},
		{
			name:                 "cli_render_yaml_project_with_relative_paths",
			arguments:            []string{"render", yamlSpecPath, "--log-level", "debug"},
			expectedOutputStdout: "Composition Rendered Successfully",
			verifyOutputFile:     filepath.Join(projectDir, "rendered_vlog.mp4"),
			expectedWidth:        1280,
			expectedHeight:       720,
		},
		{
			name:                 "cli_render_json_project",
			arguments:            []string{"render", jsonSpecPath},
			expectedOutputStdout: "Composition Rendered Successfully",
			verifyOutputFile:     filepath.Join(projectDir, "rendered_tiktok.mp4"),
			expectedWidth:        1080,
			expectedHeight:       1920,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			cmd := exec.CommandContext(ctx, cliBinaryPath, tt.arguments...)
			outputBytes, err := cmd.CombinedOutput()
			outputStr := string(outputBytes)

			if err != nil {
				t.Fatalf("CLI command %v failed: %v\nOutput: %s", tt.arguments, err, outputStr)
			}

			if !strings.Contains(outputStr, tt.expectedOutputStdout) {
				t.Errorf("CLI output did not contain %q. Full output:\n%s", tt.expectedOutputStdout, outputStr)
			}

			// If rendering an output file, verify with ffprobe
			if tt.verifyOutputFile != "" {
				fileInfo, statErr := os.Stat(tt.verifyOutputFile)
				if statErr != nil {
					t.Fatalf("Expected output file does not exist: %v", statErr)
				}
				if fileInfo.Size() < 5000 {
					t.Errorf("Output file size %d bytes is suspiciously small", fileInfo.Size())
				}

				metadata, probeErr := mediaProber.Probe(ctx, tt.verifyOutputFile)
				if probeErr != nil {
					t.Fatalf("ffprobe failed on %q: %v", tt.verifyOutputFile, probeErr)
				}

				if metadata.VideoStream.Width != tt.expectedWidth || metadata.VideoStream.Height != tt.expectedHeight {
					t.Errorf("Rendered dimensions = %dx%d, want %dx%d",
						metadata.VideoStream.Width, metadata.VideoStream.Height,
						tt.expectedWidth, tt.expectedHeight,
					)
				}
			}
		})
	}
}

func copyFile(t *testing.T, sourcePath, destinationPath string) {
	t.Helper()
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		t.Fatalf("failed reading %q: %v", sourcePath, err)
	}
	if err := os.WriteFile(destinationPath, data, 0644); err != nil {
		t.Fatalf("failed writing %q: %v", destinationPath, err)
	}
}
