package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
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

	err = api.ExecuteSelfUpdate(ctx, release, func(progress api.UpdateProgress) {
		switch progress.Stage {
		case "checking":
			fmt.Printf("🔍 %s\n", progress.Message)
		case "downloading":
			fmt.Printf("\r⬇️  %s (%.0f%%)", progress.Message, progress.Percentage)
		case "extracting":
			fmt.Printf("\n📦 %s\n", progress.Message)
		case "replacing":
			fmt.Printf("⚙️  %s\n", progress.Message)
		case "completed":
			fmt.Printf("\n")
			ui.PrintSuccess(fmt.Sprintf("Successfully upgraded Vidonex to v%s!", latestVersion))
		case "failed":
			fmt.Printf("\n")
			ui.PrintError(fmt.Errorf("%s: %s", progress.Message, progress.Error))
		}
	})

	if err != nil {
		os.Exit(1)
	}

	currentExecutablePath, _ := os.Executable()
	currentExecutablePath, _ = filepath.EvalSymlinks(currentExecutablePath)
	fmt.Printf("Binary location: %s\n\n", currentExecutablePath)
}
