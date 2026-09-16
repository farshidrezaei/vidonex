package api_test

import (
	"context"
	"testing"

	"github.com/farshidrezaei/vidonex/server/api"
)

func TestMatchAssetForCurrentPlatformTable(t *testing.T) {
	testCases := []struct {
		name         string
		assets       []api.ReleaseAsset
		targetOS     string
		targetArch   string
		expectedName string
	}{
		{
			name: "matches linux amd64 tar.gz",
			assets: []api.ReleaseAsset{
				{Name: "vidonex-darwin-arm64.tar.gz"},
				{Name: "vidonex-linux-amd64.tar.gz"},
				{Name: "vidonex-windows-amd64.zip"},
			},
			targetOS:     "linux",
			targetArch:   "amd64",
			expectedName: "vidonex-linux-amd64.tar.gz",
		},
		{
			name: "matches darwin arm64",
			assets: []api.ReleaseAsset{
				{Name: "vidonex-darwin-arm64.tar.gz"},
				{Name: "vidonex-linux-amd64.tar.gz"},
			},
			targetOS:     "darwin",
			targetArch:   "arm64",
			expectedName: "vidonex-darwin-arm64.tar.gz",
		},
		{
			name: "returns nil when no matching asset exists",
			assets: []api.ReleaseAsset{
				{Name: "vidonex-windows-amd64.zip"},
			},
			targetOS:     "freebsd",
			targetArch:   "riscv64",
			expectedName: "",
		},
		{
			name:         "handles empty assets list",
			assets:       []api.ReleaseAsset{},
			targetOS:     "linux",
			targetArch:   "amd64",
			expectedName: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			matched := api.MatchAssetForCurrentPlatform(tc.assets, tc.targetOS, tc.targetArch)
			if tc.expectedName == "" {
				if matched != nil {
					t.Fatalf("expected nil match, got %v", matched.Name)
				}
			} else {
				if matched == nil {
					t.Fatalf("expected match for %s, got nil", tc.expectedName)
				}
				if matched.Name != tc.expectedName {
					t.Fatalf("expected asset %s, got %s", tc.expectedName, matched.Name)
				}
			}
		})
	}
}

func TestExecuteSelfUpdate_NilRelease(t *testing.T) {
	err := api.ExecuteSelfUpdate(context.Background(), nil, nil)
	if err == nil {
		t.Fatalf("expected error when passing nil release, got nil")
	}
}
