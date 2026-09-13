// Package main provides the root entrypoint and Wails v2 desktop integration.
package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/farshidrezaei/vidonex/desktop"
	"github.com/farshidrezaei/vidonex/ui"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

func main() {
	appLogger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	app, err := desktop.NewApp("", appLogger)
	if err != nil {
		appLogger.Error("failed initializing vidonex desktop application", "error", err)
		os.Exit(1)
	}

	subFS, err := ui.GetFS()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: running without embedded assets: %v\n", err)
	}

	appOptions := &options.App{
		Title:             "Vidonex Studio - Video Composition Workstation",
		Width:             1440,
		Height:            900,
		MinWidth:          1024,
		MinHeight:         700,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		BackgroundColour:  &options.RGBA{R: 11, G: 15, B: 25, A: 255},
		AssetServer: &assetserver.Options{
			Assets:  subFS,
			Handler: app.Handler(),
		},
		LogLevel:   logger.INFO,
		OnStartup:  app.Startup,
		OnShutdown: app.Shutdown,
		Bind: []any{
			app,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			BackdropType:         windows.Mica,
			DisableWindowIcon:    false,
		},
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHidden(),
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  false,
		},
		Linux: &linux.Options{
			WindowIsTranslucent: false,
		},
	}

	if err := wails.Run(appOptions); err != nil {
		appLogger.Error("application execution error", "error", err)
		os.Exit(1)
	}
}
