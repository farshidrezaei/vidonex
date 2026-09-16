package api

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// UpdateProgress represents the real-time state and completion percentage of an update.
type UpdateProgress struct {
	Stage      string  `json:"stage"`      // "checking", "downloading", "extracting", "replacing", "completed", "failed"
	Percentage float64 `json:"percentage"` // 0 to 100
	Message    string  `json:"message"`
	Error      string  `json:"error,omitempty"`
}

// ProgressCallback is called whenever the update progress changes.
type ProgressCallback func(progress UpdateProgress)

// MatchAssetForCurrentPlatform finds the release asset matching host OS and CPU architecture.
func MatchAssetForCurrentPlatform(assets []ReleaseAsset, targetOS, targetArch string) *ReleaseAsset {
	expectedOS := strings.ToLower(targetOS)
	expectedArch := strings.ToLower(targetArch)

	for _, asset := range assets {
		assetNameLower := strings.ToLower(asset.Name)
		if strings.Contains(assetNameLower, expectedOS) && strings.Contains(assetNameLower, expectedArch) {
			return &asset
		}
	}
	return nil
}

// CountingReader tracks read bytes for progress calculation.
type CountingReader struct {
	Reader     io.Reader
	Total      int64
	Current    int64
	OnProgress func(percentage float64)
}

func (cr *CountingReader) Read(p []byte) (int, error) {
	n, err := cr.Reader.Read(p)
	if n > 0 {
		cr.Current += int64(n)
		if cr.Total > 0 && cr.OnProgress != nil {
			percentage := float64(cr.Current) / float64(cr.Total) * 100
			if percentage > 100 {
				percentage = 100
			}
			cr.OnProgress(percentage)
		}
	}
	return n, err
}

// ExecuteSelfUpdate performs the end-to-end download, extraction, and atomic binary replacement.
func ExecuteSelfUpdate(ctx context.Context, release *GitHubReleaseResponse, onProgress ProgressCallback) error {
	if release == nil {
		return fmt.Errorf("release metadata cannot be nil")
	}

	report := func(stage string, percentage float64, message string, err error) {
		if onProgress != nil {
			errStr := ""
			if err != nil {
				errStr = err.Error()
			}
			onProgress(UpdateProgress{
				Stage:      stage,
				Percentage: percentage,
				Message:    message,
				Error:      errStr,
			})
		}
	}

	report("checking", 5, "Resolving target asset for platform...", nil)

	selectedAsset := MatchAssetForCurrentPlatform(release.Assets, runtime.GOOS, runtime.GOARCH)
	if selectedAsset == nil {
		err := fmt.Errorf("no matching precompiled binary asset found for %s/%s in release %s", runtime.GOOS, runtime.GOARCH, release.TagName)
		report("failed", 0, err.Error(), err)
		return err
	}

	currentExecutablePath, err := os.Executable()
	if err != nil {
		err = fmt.Errorf("could not determine running executable path: %w", err)
		report("failed", 0, err.Error(), err)
		return err
	}
	currentExecutablePath, err = filepath.EvalSymlinks(currentExecutablePath)
	if err != nil {
		err = fmt.Errorf("could not resolve symlink: %w", err)
		report("failed", 0, err.Error(), err)
		return err
	}

	// Prepare temporary file
	tempDir := filepath.Dir(currentExecutablePath)
	tempFile, err := os.CreateTemp(tempDir, "vidonex_update_*")
	if err != nil {
		// Fallback to system temp directory if binary directory is not writable
		tempFile, err = os.CreateTemp("", "vidonex_update_*")
		if err != nil {
			err = fmt.Errorf("failed creating temporary download file: %w", err)
			report("failed", 0, err.Error(), err)
			return err
		}
	}
	tempFilePath := tempFile.Name()
	defer func() { _ = os.Remove(tempFilePath) }()

	report("downloading", 10, fmt.Sprintf("Downloading %s...", selectedAsset.Name), nil)

	downloadReq, err := http.NewRequestWithContext(ctx, http.MethodGet, selectedAsset.BrowserDownloadURL, nil)
	if err != nil {
		err = fmt.Errorf("failed preparing download request: %w", err)
		report("failed", 0, err.Error(), err)
		return err
	}
	downloadReq.Header.Set("User-Agent", "Vidonex-SelfUpdater/"+CurrentEngineVersion)

	httpClient := &http.Client{Timeout: 0}
	resp, err := httpClient.Do(downloadReq)
	if err != nil {
		err = fmt.Errorf("failed downloading release: %w", err)
		report("failed", 0, err.Error(), err)
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("download failed with HTTP status %d", resp.StatusCode)
		report("failed", 0, err.Error(), err)
		return err
	}

	countingReader := &CountingReader{
		Reader: resp.Body,
		Total:  selectedAsset.Size,
		OnProgress: func(percentage float64) {
			// Map download to 10% - 80% range
			mappedPercentage := 10 + (percentage * 0.7)
			report("downloading", mappedPercentage, fmt.Sprintf("Downloading: %.0f%%", percentage), nil)
		},
	}

	report("extracting", 80, "Extracting binary...", nil)

	// If asset is a tar.gz / tgz archive, extract the binary
	if strings.HasSuffix(selectedAsset.Name, ".tar.gz") || strings.HasSuffix(selectedAsset.Name, ".tgz") {
		gzipReader, err := gzip.NewReader(countingReader)
		if err != nil {
			err = fmt.Errorf("failed decompressing gzip: %w", err)
			report("failed", 0, err.Error(), err)
			return err
		}
		defer func() { _ = gzipReader.Close() }()

		tarReader := tar.NewReader(gzipReader)
		found := false
		for {
			header, err := tarReader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				err = fmt.Errorf("failed reading tar archive: %w", err)
				report("failed", 0, err.Error(), err)
				return err
			}
			baseName := filepath.Base(header.Name)
			if baseName == "vidonex" || baseName == "vidonex.exe" {
				if _, err := io.Copy(tempFile, tarReader); err != nil {
					err = fmt.Errorf("failed writing extracted binary: %w", err)
					report("failed", 0, err.Error(), err)
					return err
				}
				found = true
				break
			}
		}
		if !found {
			err = fmt.Errorf("archive %s did not contain vidonex binary", selectedAsset.Name)
			report("failed", 0, err.Error(), err)
			return err
		}
	} else {
		// Raw binary
		if _, err := io.Copy(tempFile, countingReader); err != nil {
			err = fmt.Errorf("failed writing binary: %w", err)
			report("failed", 0, err.Error(), err)
			return err
		}
	}
	_ = tempFile.Close()

	// Ensure binary has executable permissions
	if err := os.Chmod(tempFilePath, 0755); err != nil {
		err = fmt.Errorf("failed setting executable permissions: %w", err)
		report("failed", 0, err.Error(), err)
		return err
	}

	report("replacing", 90, "Applying atomic binary replacement...", nil)

	// Atomic binary replacement with rollback
	backupPath := currentExecutablePath + ".old"
	_ = os.Remove(backupPath)
	if err := os.Rename(currentExecutablePath, backupPath); err != nil {
		err = fmt.Errorf("insufficient file permissions to overwrite binary at %s (try running with higher privileges): %w", currentExecutablePath, err)
		report("failed", 0, err.Error(), err)
		return err
	}

	if err := os.Rename(tempFilePath, currentExecutablePath); err != nil {
		// Rollback old binary
		_ = os.Rename(backupPath, currentExecutablePath)
		err = fmt.Errorf("failed replacing binary: %w", err)
		report("failed", 0, err.Error(), err)
		return err
	}
	_ = os.Remove(backupPath)

	report("completed", 100, fmt.Sprintf("Successfully upgraded Vidonex to %s!", release.TagName), nil)
	return nil
}
