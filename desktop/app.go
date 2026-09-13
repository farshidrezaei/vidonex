// Package desktop implements the Wails v2 native desktop application shell for Vidonex.
package desktop

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/farshidrezaei/vidonex/presets"
	"github.com/farshidrezaei/vidonex/probe"
	"github.com/farshidrezaei/vidonex/server"
	"github.com/farshidrezaei/vidonex/server/db"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// ApplicationVersion defines the semantic version of the desktop application.
const ApplicationVersion = "1.0.0"

// DefaultServerPort specifies the default port (0 indicates dynamic ephemeral allocation).
const DefaultServerPort = 0

// ImportedAssetResult carries metadata of an asset imported directly from the native filesystem.
type ImportedAssetResult struct {
	Asset db.MediaAssetRecord `json:"asset"`
	Error string              `json:"error,omitempty"`
}

// SystemCapabilities provides information on the host OS and available hardware encoders.
type SystemCapabilities struct {
	OS                    string                        `json:"os"`
	Architecture          string                        `json:"architecture"`
	CPUCount              int                           `json:"cpu_count"`
	AvailableAccelerators []presets.HardwareAccelerator `json:"available_accelerators"`
	FFmpegVersion         string                        `json:"ffmpeg_version"`
	FFprobeVersion        string                        `json:"ffprobe_version"`
	ServerAddress         string                        `json:"server_address"`
	ServerPort            int                           `json:"server_port"`
}

// App manages the desktop workstation lifecycle, background server, and native OS APIs.
type App struct {
	context       context.Context
	server        *server.Server
	listener      net.Listener
	serverPort    int
	serverURL     string
	mediaProber   probe.MediaProber
	dataDirectory string
	logger        *slog.Logger
}

// NewApp initializes a new desktop App instance.
func NewApp(dataDirectory string, logger *slog.Logger) (*App, error) {
	if logger == nil {
		logger = slog.Default()
	}

	if dataDirectory == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			homeDir = "."
		}
		dataDirectory = filepath.Join(homeDir, ".vidonex")
	}

	if err := os.MkdirAll(dataDirectory, 0755); err != nil {
		return nil, fmt.Errorf("failed creating app data directory: %w", err)
	}

	// Bind ephemeral listener on loopback
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("failed allocating server listener: %w", err)
	}

	serverPort := listener.Addr().(*net.TCPAddr).Port
	serverURL := fmt.Sprintf("http://127.0.0.1:%d", serverPort)

	serverInstance, err := server.New(server.Config{
		Host:          "127.0.0.1",
		Port:          serverPort,
		DataDirectory: dataDirectory,
		Logger:        logger,
	})
	if err != nil {
		_ = listener.Close()
		return nil, fmt.Errorf("failed initializing embedded server: %w", err)
	}

	appInstance := &App{
		server:        serverInstance,
		listener:      listener,
		serverPort:    serverPort,
		serverURL:     serverURL,
		mediaProber:   probe.NewCachedProber(probe.NewFFprobeProber("ffprobe")),
		dataDirectory: dataDirectory,
		logger:        logger,
	}

	return appInstance, nil
}

// Startup is called by Wails when the application begins.
func (a *App) Startup(ctx context.Context) {
	a.context = ctx
	a.logger.Info("vidonex desktop starting up", "server_port", a.serverPort)

	// Launch embedded HTTP server in background
	go func() {
		if err := a.server.StartListener(a.listener); err != nil {
			a.logger.Error("embedded server stopped", "error", err)
		}
	}()
}

// Shutdown is called by Wails when the application terminates.
func (a *App) Shutdown(ctx context.Context) {
	a.logger.Info("vidonex desktop shutting down")
	shutdownContext, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_ = a.server.Stop(shutdownContext)
	if a.listener != nil {
		_ = a.listener.Close()
	}
}

// Handler returns the HTTP Handler for proxying REST API and WebSocket requests within Wails AssetServer.
func (a *App) Handler() http.Handler {
	return a.server.Handler()
}

// GetServerInfo returns the internal local server address and port.
func (a *App) GetServerInfo() map[string]any {
	return map[string]any{
		"port": a.serverPort,
		"url":  a.serverURL,
		"ws":   fmt.Sprintf("ws://127.0.0.1:%d/ws", a.serverPort),
	}
}

// GetAppVersion returns the current engine and application release version.
func (a *App) GetAppVersion() string {
	return ApplicationVersion
}

// GetSystemCapabilities inspects the system for GPU acceleration and FFmpeg tools.
func (a *App) GetSystemCapabilities() SystemCapabilities {
	accelerators := make([]presets.HardwareAccelerator, 0)

	// Test NVENC
	if checkCommandSuccess("ffmpeg", "-hide_banner", "-encoders") {
		if checkCommandGrep("ffmpeg", []string{"-hide_banner", "-encoders"}, "h264_nvenc") {
			accelerators = append(accelerators, presets.AcceleratorNVENC)
		}
		if runtime.GOOS == "darwin" && checkCommandGrep("ffmpeg", []string{"-hide_banner", "-encoders"}, "h264_videotoolbox") {
			accelerators = append(accelerators, presets.AcceleratorVideoToolbox)
		}
		if runtime.GOOS == "linux" && checkCommandGrep("ffmpeg", []string{"-hide_banner", "-encoders"}, "h264_vaapi") {
			accelerators = append(accelerators, presets.AcceleratorVAAPI)
		}
		if checkCommandGrep("ffmpeg", []string{"-hide_banner", "-encoders"}, "h264_qsv") {
			accelerators = append(accelerators, presets.AcceleratorQSV)
		}
	}

	ffmpegVersion := getFirstLineOfOutput("ffmpeg", "-version")
	ffprobeVersion := getFirstLineOfOutput("ffprobe", "-version")

	return SystemCapabilities{
		OS:                    runtime.GOOS,
		Architecture:          runtime.GOARCH,
		CPUCount:              runtime.NumCPU(),
		AvailableAccelerators: accelerators,
		FFmpegVersion:         ffmpegVersion,
		FFprobeVersion:        ffprobeVersion,
		ServerAddress:         a.serverURL,
		ServerPort:            a.serverPort,
	}
}

// SelectMediaFiles opens a native OS multi-file open dialog for videos, audios, and images.
func (a *App) SelectMediaFiles() ([]string, error) {
	selectedFiles, err := wailsRuntime.OpenMultipleFilesDialog(a.context, wailsRuntime.OpenDialogOptions{
		Title: "Import Media Files",
		Filters: []wailsRuntime.FileFilter{
			{
				DisplayName: "Media Files (Video, Audio, Image)",
				Pattern:     "*.mp4;*.mkv;*.mov;*.webm;*.avi;*.flv;*.mp3;*.wav;*.aac;*.flac;*.ogg;*.m4a;*.png;*.jpg;*.jpeg;*.webp;*.svg",
			},
			{
				DisplayName: "Video Files (*.mp4, *.mov, *.mkv, *.webm)",
				Pattern:     "*.mp4;*.mkv;*.mov;*.webm;*.avi;*.flv",
			},
			{
				DisplayName: "Audio Files (*.mp3, *.wav, *.aac, *.flac, *.ogg)",
				Pattern:     "*.mp3;*.wav;*.aac;*.flac;*.ogg;*.m4a",
			},
			{
				DisplayName: "Image Files (*.png, *.jpg, *.jpeg, *.webp, *.svg)",
				Pattern:     "*.png;*.jpg;*.jpeg;*.webp;*.svg",
			},
			{
				DisplayName: "All Files (*.*)",
				Pattern:     "*.*",
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed opening file dialog: %w", err)
	}

	return selectedFiles, nil
}

// SaveProjectFileDialog opens a native OS file save dialog for Vidonex project files.
func (a *App) SaveProjectFileDialog(defaultName string) (string, error) {
	if defaultName == "" {
		defaultName = "project.vidonex.json"
	}
	savePath, err := wailsRuntime.SaveFileDialog(a.context, wailsRuntime.SaveDialogOptions{
		Title:           "Save Vidonex Project",
		DefaultFilename: defaultName,
		Filters: []wailsRuntime.FileFilter{
			{
				DisplayName: "Vidonex Project File (*.vidonex.json, *.vdx)",
				Pattern:     "*.vidonex.json;*.vdx;*.json;*.yaml;*.yml",
			},
			{
				DisplayName: "All Files (*.*)",
				Pattern:     "*.*",
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed opening save dialog: %w", err)
	}
	return savePath, nil
}

// OpenProjectFileDialog opens a native OS file open dialog for project files.
func (a *App) OpenProjectFileDialog() (string, error) {
	filePath, err := wailsRuntime.OpenFileDialog(a.context, wailsRuntime.OpenDialogOptions{
		Title: "Open Vidonex Project File",
		Filters: []wailsRuntime.FileFilter{
			{
				DisplayName: "Vidonex Project Files (*.vidonex.json, *.vdx, *.yaml, *.yml)",
				Pattern:     "*.vidonex.json;*.vdx;*.json;*.yaml;*.yml",
			},
			{
				DisplayName: "All Files (*.*)",
				Pattern:     "*.*",
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed opening project dialog: %w", err)
	}
	return filePath, nil
}

// SelectExportDirectory opens a native OS folder picker to choose an export destination.
func (a *App) SelectExportDirectory() (string, error) {
	directoryPath, err := wailsRuntime.OpenDirectoryDialog(a.context, wailsRuntime.OpenDialogOptions{
		Title: "Select Video Export Destination Directory",
	})
	if err != nil {
		return "", fmt.Errorf("failed opening directory dialog: %w", err)
	}
	return directoryPath, nil
}

// ImportLocalMedia probes and imports local media files into a project without network copying.
func (a *App) ImportLocalMedia(projectID string, filePaths []string) []ImportedAssetResult {
	results := make([]ImportedAssetResult, 0, len(filePaths))
	mediaDirectory := a.server.MediaDirectory()
	databaseInstance := a.server.Database()

	for _, originalPath := range filePaths {
		fileInfo, statErr := os.Stat(originalPath)
		if statErr != nil {
			results = append(results, ImportedAssetResult{
				Error: fmt.Sprintf("file %q not found: %v", originalPath, statErr),
			})
			continue
		}

		assetID := generateRandomID("asset")
		fileName := filepath.Base(originalPath)
		fileExtension := strings.ToLower(filepath.Ext(originalPath))

		fileType := "video"
		switch fileExtension {
		case ".mp3", ".wav", ".aac", ".flac", ".ogg", ".m4a":
			fileType = "audio"
		case ".png", ".jpg", ".jpeg", ".webp", ".svg", ".bmp", ".gif":
			fileType = "image"
		}

		probeContext, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		mediaMetadata, probeErr := a.mediaProber.Probe(probeContext, originalPath)
		cancel()

		var durationSeconds float64
		var width, height, sampleRate, channels int
		var frameRate float64

		if probeErr == nil && mediaMetadata != nil {
			durationSeconds = mediaMetadata.Duration.Seconds()
			if mediaMetadata.VideoStream != nil {
				width = mediaMetadata.VideoStream.Width
				height = mediaMetadata.VideoStream.Height
				frameRate = mediaMetadata.VideoStream.FPS.Seconds()
			}
			if mediaMetadata.AudioStream != nil {
				sampleRate = mediaMetadata.AudioStream.SampleRate
				channels = mediaMetadata.AudioStream.Channels
			}
		}

		thumbnailPath := ""
		switch fileType {
		case "video":
			thumbFileName := fmt.Sprintf("thumb_%s.jpg", assetID)
			thumbFullPath := filepath.Join(mediaDirectory, thumbFileName)
			thumbCmd := exec.Command("ffmpeg", "-y", "-ss", "00:00:00.500", "-i", originalPath, "-vframes", "1", "-vf", "scale=320:-1", thumbFullPath)
			if thumbCmd.Run() == nil {
				thumbnailPath = thumbFileName
			}
		case "image":
			thumbFileName := fmt.Sprintf("thumb_%s.jpg", assetID)
			thumbFullPath := filepath.Join(mediaDirectory, thumbFileName)
			thumbCmd := exec.Command("ffmpeg", "-y", "-i", originalPath, "-vframes", "1", "-vf", "scale=320:-1", thumbFullPath)
			if thumbCmd.Run() == nil {
				thumbnailPath = thumbFileName
			} else {
				thumbnailPath = originalPath
			}
		}

		assetRecord := db.MediaAssetRecord{
			ID:            assetID,
			ProjectID:     projectID,
			FileName:      fileName,
			FilePath:      originalPath, // Native Absolute Path on disk!
			FileType:      fileType,
			FileSizeBytes: fileInfo.Size(),
			Duration:      durationSeconds,
			Width:         width,
			Height:        height,
			FrameRate:     frameRate,
			SampleRate:    sampleRate,
			Channels:      channels,
			ThumbnailPath: thumbnailPath,
			WaveformData:  "[]",
			CreatedAt:     time.Now().UTC(),
		}

		if saveErr := databaseInstance.SaveMediaAsset(context.Background(), assetRecord); saveErr != nil {
			results = append(results, ImportedAssetResult{
				Error: fmt.Sprintf("failed persisting asset %q: %v", fileName, saveErr),
			})
			continue
		}

		results = append(results, ImportedAssetResult{
			Asset: assetRecord,
		})
	}

	return results
}

func generateRandomID(prefix string) string {
	bytes := make([]byte, 8)
	_, _ = rand.Read(bytes)
	return fmt.Sprintf("%s_%s", prefix, hex.EncodeToString(bytes))
}

func checkCommandSuccess(name string, args ...string) bool {
	cmd := exec.Command(name, args...)
	return cmd.Run() == nil
}

func checkCommandGrep(name string, args []string, pattern string) bool {
	cmd := exec.Command(name, args...)
	outputBytes, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(outputBytes), pattern)
}

func getFirstLineOfOutput(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	outputBytes, err := cmd.Output()
	if err != nil {
		return ""
	}
	lines := strings.Split(string(outputBytes), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	}
	return ""
}
