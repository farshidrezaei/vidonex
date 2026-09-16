package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/farshidrezaei/vidonex/server/api"
)

func TestCompareSemanticVersionsTable(t *testing.T) {
	testCases := []struct {
		name     string
		versionA string
		versionB string
		expected int
	}{
		{
			name:     "identical versions",
			versionA: "1.6.0",
			versionB: "1.6.0",
			expected: 0,
		},
		{
			name:     "newer minor version",
			versionA: "1.7.0",
			versionB: "1.6.0",
			expected: 1,
		},
		{
			name:     "newer patch version",
			versionA: "1.6.1",
			versionB: "1.6.0",
			expected: 1,
		},
		{
			name:     "older patch version",
			versionA: "1.5.9",
			versionB: "1.6.0",
			expected: -1,
		},
		{
			name:     "with v prefix and suffix",
			versionA: "v2.0.0-beta.1",
			versionB: "1.6.0",
			expected: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := api.CompareSemanticVersions(tc.versionA, tc.versionB)
			if result != tc.expected {
				t.Fatalf("expected CompareSemanticVersions(%q, %q) = %d, got %d", tc.versionA, tc.versionB, tc.expected, result)
			}
		})
	}
}

func TestCollectSystemDiagnosticsTable(t *testing.T) {
	diagnostics := api.CollectSystemDiagnostics()
	if diagnostics.OperatingSystem == "" {
		t.Errorf("expected non-empty OS")
	}
	if diagnostics.Architecture == "" {
		t.Errorf("expected non-empty Architecture")
	}
	if diagnostics.CPUCount <= 0 {
		t.Errorf("expected CPUCount > 0, got %d", diagnostics.CPUCount)
	}
}

func TestHandleCheckUpdate_MethodTable(t *testing.T) {
	handlers := api.NewHandlers(nil, nil, t.TempDir(), t.TempDir(), nil)

	testCases := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{
			name:           "disallowed POST",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "allowed GET",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/api/version/check", nil)
			w := httptest.NewRecorder()
			handlers.HandleCheckUpdate(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, w.Code)
			}
		})
	}
}
