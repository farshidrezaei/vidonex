package desktop_test

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/farshidrezaei/vidonex/desktop"
)

func TestDesktopAppInitializationTable(t *testing.T) {
	temporaryDirectory, err := os.MkdirTemp("", "vidonex_desktop_test_*")
	if err != nil {
		t.Fatalf("failed creating temp directory: %v", err)
	}
	defer func() { _ = os.RemoveAll(temporaryDirectory) }()

	tests := []struct {
		name          string
		dataDirectory string
		expectError   bool
	}{
		{
			name:          "ValidDataDirectory",
			dataDirectory: filepath.Join(temporaryDirectory, "app_data"),
			expectError:   false,
		},
		{
			name:          "DefaultEmptyDataDirectory",
			dataDirectory: "",
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
				Level: slog.LevelError,
			}))

			appInstance, err := desktop.NewApp(tt.dataDirectory, logger)
			if tt.expectError && err == nil {
				t.Fatalf("expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Fatalf("expected no error but got: %v", err)
			}
			if appInstance == nil {
				t.Fatalf("expected non-nil appInstance")
			}

			serverInfo := appInstance.GetServerInfo()
			if serverInfo["port"] == nil || serverInfo["port"].(int) <= 0 {
				t.Errorf("expected positive server port, got %v", serverInfo["port"])
			}
			if appInstance.GetAppVersion() != desktop.ApplicationVersion {
				t.Errorf("expected version %s, got %s", desktop.ApplicationVersion, appInstance.GetAppVersion())
			}

			capabilities := appInstance.GetSystemCapabilities()
			if capabilities.CPUCount <= 0 {
				t.Errorf("expected positive CPU count, got %d", capabilities.CPUCount)
			}
			if capabilities.ServerPort <= 0 {
				t.Errorf("expected positive ServerPort in capabilities, got %d", capabilities.ServerPort)
			}
		})
	}
}
