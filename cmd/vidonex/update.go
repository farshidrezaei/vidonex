package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/farshidrezaei/vidonex/internal/cli"
	"github.com/farshidrezaei/vidonex/server/api"
)

func executeUpgradeCommand(arguments []string) {
	fs := flag.NewFlagSet("upgrade", flag.ExitOnError)
	checkOnlyFlag := fs.Bool("check", false, "Only check if an update is available without installing")
	forceFlag := fs.Bool("force", false, "Force reinstall even if already on the latest version")

	if err := fs.Parse(reorderFlags(arguments)); err != nil {
		fmt.Fprintf(os.Stderr, "Failed parsing flags: %v\n", err)
		os.Exit(1)
	}

	ui := cli.NewUI(os.Stdout)
	ui.PrintBanner()

	fmt.Printf("\n%s\n", ui.Colorize(cli.ColorCyan+cli.ColorBold, "Checking for latest Vidonex release..."))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	release, err := api.FetchLatestRelease(ctx, true)
	if err != nil {
		ui.PrintError(fmt.Errorf("failed fetching latest release: %w", err))
		os.Exit(1)
	}

	latestVersion := strings.TrimPrefix(release.TagName, "v")
	hasUpdate := api.CompareSemanticVersions(latestVersion, EngineVersion) > 0

	fmt.Printf("Current Installed Version: %s\n", ui.Colorize(cli.ColorDim, "v"+EngineVersion))
	fmt.Printf("Latest Available Version:  %s\n\n", ui.Colorize(cli.ColorBold, "v"+latestVersion))

	if !hasUpdate && !*forceFlag {
		ui.PrintSuccess(fmt.Sprintf("Vidonex is already up to date (v%s).", EngineVersion))
		return
	}

	if *checkOnlyFlag {
		if hasUpdate {
			fmt.Printf("%s\n", ui.Colorize(cli.ColorYellow+cli.ColorBold, fmt.Sprintf("✨ Update available: v%s -> v%s", EngineVersion, latestVersion)))
			fmt.Printf("Release URL: %s\n", release.HTMLURL)
			fmt.Println("Run `vidonex upgrade` to install the update.")
		} else {
			fmt.Printf("No update available.\n")
		}
		return
	}

	// Upgrade Execution
	fmt.Printf("%s\n", ui.Colorize(cli.ColorGreen+cli.ColorBold, fmt.Sprintf("Upgrading Vidonex to v%s...", latestVersion)))

	// Find suitable asset for current OS and architecture
	expectedOS := runtime.GOOS
	expectedArch := runtime.GOARCH

	var selectedAsset *api.ReleaseAsset
	for _, asset := range release.Assets {
		assetNameLower := strings.ToLower(asset.Name)
		if strings.Contains(assetNameLower, expectedOS) && strings.Contains(assetNameLower, expectedArch) {
			selectedAsset = &asset
			break
		}
	}

	currentExecutablePath, err := os.Executable()
	if err != nil {
		ui.PrintError(fmt.Errorf("could not determine running executable path: %w", err))
		os.Exit(1)
	}
	currentExecutablePath, err = filepath.EvalSymlinks(currentExecutablePath)
	if err != nil {
		ui.PrintError(fmt.Errorf("could not resolve symlink: %w", err))
		os.Exit(1)
	}

	if selectedAsset == nil {
		fmt.Printf("%s\n", ui.Colorize(cli.ColorYellow, fmt.Sprintf("No precompiled binary asset matching %s/%s was found in release v%s.", expectedOS, expectedArch, latestVersion)))
		fmt.Printf("You can update Vidonex via Go:\n\n")
		fmt.Printf("  %s\n\n", ui.Colorize(cli.ColorWhiteBold, "go install github.com/farshidrezaei/vidonex/cmd/vidonex@latest"))
		fmt.Printf("Or visit the GitHub release page directly:\n  %s\n", release.HTMLURL)
		return
	}

	fmt.Printf("Downloading %s (%s)...\n", selectedAsset.Name, cli.FormatFileSize(selectedAsset.Size))

	downloadReq, err := http.NewRequestWithContext(ctx, http.MethodGet, selectedAsset.BrowserDownloadURL, nil)
	if err != nil {
		ui.PrintError(fmt.Errorf("failed preparing download: %w", err))
		os.Exit(1)
	}
	downloadReq.Header.Set("User-Agent", "Vidonex-Updater/"+EngineVersion)

	resp, err := http.DefaultClient.Do(downloadReq)
	if err != nil {
		ui.PrintError(fmt.Errorf("failed downloading release: %w", err))
		os.Exit(1)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		ui.PrintError(fmt.Errorf("download failed with HTTP status %d", resp.StatusCode))
		os.Exit(1)
	}

	tempDir := filepath.Dir(currentExecutablePath)
	tempFile, err := os.CreateTemp(tempDir, "vidonex_update_*")
	if err != nil {
		// Fallback to system temp directory if binary dir is not writable
		tempFile, err = os.CreateTemp("", "vidonex_update_*")
		if err != nil {
			ui.PrintError(fmt.Errorf("failed creating temporary file: %w", err))
			os.Exit(1)
		}
	}
	tempFilePath := tempFile.Name()
	defer func() { _ = os.Remove(tempFilePath) }()

	// If asset is a tar.gz archive, extract the binary
	if strings.HasSuffix(selectedAsset.Name, ".tar.gz") || strings.HasSuffix(selectedAsset.Name, ".tgz") {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			ui.PrintError(fmt.Errorf("failed decompressing gzip: %w", err))
			os.Exit(1)
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
				ui.PrintError(fmt.Errorf("failed reading tar archive: %w", err))
				os.Exit(1)
			}
			if filepath.Base(header.Name) == "vidonex" || filepath.Base(header.Name) == "vidonex.exe" {
				if _, err := io.Copy(tempFile, tarReader); err != nil {
					ui.PrintError(fmt.Errorf("failed writing extracted binary: %w", err))
					os.Exit(1)
				}
				found = true
				break
			}
		}
		if !found {
			ui.PrintError(fmt.Errorf("archive did not contain vidonex binary"))
			os.Exit(1)
		}
	} else {
		if _, err := io.Copy(tempFile, resp.Body); err != nil {
			ui.PrintError(fmt.Errorf("failed writing binary: %w", err))
			os.Exit(1)
		}
	}
	_ = tempFile.Close()

	// Ensure binary has executable permissions
	if err := os.Chmod(tempFilePath, 0755); err != nil {
		ui.PrintError(fmt.Errorf("failed setting executable permissions: %w", err))
		os.Exit(1)
	}

	// Atomic binary replacement
	backupPath := currentExecutablePath + ".old"
	_ = os.Remove(backupPath)
	if err := os.Rename(currentExecutablePath, backupPath); err != nil {
		ui.PrintError(fmt.Errorf("failed backing up old binary (insufficient permissions? try with sudo): %w", err))
		os.Exit(1)
	}

	if err := os.Rename(tempFilePath, currentExecutablePath); err != nil {
		// Rollback
		_ = os.Rename(backupPath, currentExecutablePath)
		ui.PrintError(fmt.Errorf("failed replacing binary: %w", err))
		os.Exit(1)
	}
	_ = os.Remove(backupPath)

	ui.PrintSuccess(fmt.Sprintf("Successfully upgraded Vidonex to v%s!", latestVersion))
	fmt.Printf("Binary location: %s\n\n", currentExecutablePath)
}
