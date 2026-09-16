package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/farshidrezaei/vidonex/presets"
)

// CurrentEngineVersion defines the current semantic version of the engine.
const CurrentEngineVersion = "1.6.0"

// GitHubRepository defines the upstream repository slug for update checks.
const GitHubRepository = "farshidrezaei/vidonex"

// ReleaseAsset represents an artifact asset attached to a GitHub release.
type ReleaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
	ContentType        string `json:"content_type"`
}

// GitHubReleaseResponse maps the GitHub API release payload.
type GitHubReleaseResponse struct {
	TagName     string         `json:"tag_name"`
	Name        string         `json:"name"`
	Body        string         `json:"body"`
	HTMLURL     string         `json:"html_url"`
	PublishedAt string         `json:"published_at"`
	Draft       bool           `json:"draft"`
	Prerelease  bool           `json:"prerelease"`
	Assets      []ReleaseAsset `json:"assets"`
}

// SystemDiagnostics represents host capabilities, runtime specs, and media tool availability.
type SystemDiagnostics struct {
	OperatingSystem       string                        `json:"os"`
	Architecture          string                        `json:"architecture"`
	CPUCount              int                           `json:"cpu_count"`
	AvailableAccelerators []presets.HardwareAccelerator `json:"available_accelerators"`
	FFmpegVersion         string                        `json:"ffmpeg_version"`
	FFprobeVersion        string                        `json:"ffprobe_version"`
	GoVersion             string                        `json:"go_version"`
}

// VersionCheckResponse is the JSON response payload for GET /api/version/check.
type VersionCheckResponse struct {
	CurrentVersion string            `json:"current_version"`
	LatestVersion  string            `json:"latest_version"`
	HasUpdate      bool              `json:"has_update"`
	ReleaseName    string            `json:"release_name,omitempty"`
	ReleaseNotes   string            `json:"release_notes,omitempty"`
	ReleaseURL     string            `json:"release_url,omitempty"`
	PublishedAt    string            `json:"published_at,omitempty"`
	Assets         []ReleaseAsset    `json:"assets,omitempty"`
	Diagnostics    SystemDiagnostics `json:"diagnostics"`
	CheckedAt      string            `json:"checked_at"`
}

var (
	releaseCacheMutex sync.Mutex
	cachedRelease     *GitHubReleaseResponse
	cachedReleaseTime time.Time
)

const releaseCacheDuration = 10 * time.Minute

// FetchLatestRelease retrieves the latest release information from GitHub, using an in-memory cache.
func FetchLatestRelease(ctx context.Context, forceRefresh bool) (*GitHubReleaseResponse, error) {
	releaseCacheMutex.Lock()
	if !forceRefresh && cachedRelease != nil && time.Since(cachedReleaseTime) < releaseCacheDuration {
		cachedCopy := *cachedRelease
		releaseCacheMutex.Unlock()
		return &cachedCopy, nil
	}
	releaseCacheMutex.Unlock()

	requestURL := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", GitHubRepository)
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating release check request: %w", err)
	}

	httpRequest.Header.Set("User-Agent", "Vidonex-UpdateChecker/"+CurrentEngineVersion)
	httpRequest.Header.Set("Accept", "application/vnd.github.v3+json")

	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("failed requesting latest release: %w", err)
	}
	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api returned status %d", httpResponse.StatusCode)
	}

	var release GitHubReleaseResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed decoding github release payload: %w", err)
	}

	releaseCacheMutex.Lock()
	cachedRelease = &release
	cachedReleaseTime = time.Now()
	releaseCacheMutex.Unlock()

	return &release, nil
}

// CompareSemanticVersions returns 1 if versionA > versionB, -1 if versionA < versionB, and 0 if equal.
func CompareSemanticVersions(versionA, versionB string) int {
	cleanA := strings.TrimPrefix(strings.TrimSpace(versionA), "v")
	cleanB := strings.TrimPrefix(strings.TrimSpace(versionB), "v")

	partsA := strings.Split(cleanA, ".")
	partsB := strings.Split(cleanB, ".")

	maxLen := len(partsA)
	if len(partsB) > maxLen {
		maxLen = len(partsB)
	}

	for i := 0; i < maxLen; i++ {
		var numA, numB int
		if i < len(partsA) {
			numA, _ = strconv.Atoi(strings.Split(partsA[i], "-")[0])
		}
		if i < len(partsB) {
			numB, _ = strconv.Atoi(strings.Split(partsB[i], "-")[0])
		}

		if numA > numB {
			return 1
		}
		if numA < numB {
			return -1
		}
	}

	return 0
}

// CollectSystemDiagnostics inspects the host environment for FFmpeg, GPUs, and hardware encoders.
func CollectSystemDiagnostics() SystemDiagnostics {
	accelerators := make([]presets.HardwareAccelerator, 0)

	encoderOutputBytes, err := exec.Command("ffmpeg", "-hide_banner", "-encoders").Output()
	if err == nil {
		outputString := string(encoderOutputBytes)
		if strings.Contains(outputString, "h264_nvenc") {
			accelerators = append(accelerators, presets.AcceleratorNVENC)
		}
		if runtime.GOOS == "darwin" && strings.Contains(outputString, "h264_videotoolbox") {
			accelerators = append(accelerators, presets.AcceleratorVideoToolbox)
		}
		if runtime.GOOS == "linux" && strings.Contains(outputString, "h264_vaapi") {
			accelerators = append(accelerators, presets.AcceleratorVAAPI)
		}
		if strings.Contains(outputString, "h264_qsv") {
			accelerators = append(accelerators, presets.AcceleratorQSV)
		}
	}

	ffmpegVersion := getFirstOutputLine("ffmpeg", "-version")
	ffprobeVersion := getFirstOutputLine("ffprobe", "-version")

	return SystemDiagnostics{
		OperatingSystem:       runtime.GOOS,
		Architecture:          runtime.GOARCH,
		CPUCount:              runtime.NumCPU(),
		AvailableAccelerators: accelerators,
		FFmpegVersion:         ffmpegVersion,
		FFprobeVersion:        ffprobeVersion,
		GoVersion:             runtime.Version(),
	}
}

func getFirstOutputLine(name string, args ...string) string {
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

// HandleCheckUpdate handles GET /api/version/check.
func (h *Handlers) HandleCheckUpdate(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	if httpRequest.Method != http.MethodGet {
		http.Error(responseWriter, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	forceRefresh := httpRequest.URL.Query().Get("force") == "true"
	diagnostics := CollectSystemDiagnostics()

	checkContext, cancel := context.WithTimeout(httpRequest.Context(), 5*time.Second)
	defer cancel()

	release, err := FetchLatestRelease(checkContext, forceRefresh)

	response := VersionCheckResponse{
		CurrentVersion: CurrentEngineVersion,
		LatestVersion:  CurrentEngineVersion,
		HasUpdate:      false,
		Diagnostics:    diagnostics,
		CheckedAt:      time.Now().UTC().Format(time.RFC3339),
	}

	if err != nil {
		h.logger.Warn("failed checking for updates", "error", err)
		// Return 200 with current version and diagnostics so client UI doesn't fail catastrophically
		respondJSON(responseWriter, http.StatusOK, response)
		return
	}

	latestCleanTag := strings.TrimPrefix(release.TagName, "v")
	response.LatestVersion = latestCleanTag
	response.HasUpdate = CompareSemanticVersions(latestCleanTag, CurrentEngineVersion) > 0
	response.ReleaseName = release.Name
	response.ReleaseNotes = release.Body
	response.ReleaseURL = release.HTMLURL
	response.PublishedAt = release.PublishedAt
	response.Assets = release.Assets

	respondJSON(responseWriter, http.StatusOK, response)
}
