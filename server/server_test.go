package server_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/farshidrezaei/vidonex/server"
)

func TestServer_MediaFileServingTable(t *testing.T) {
	tempDir := t.TempDir()
	dataDir := filepath.Join(tempDir, "data")
	mediaDir := filepath.Join(dataDir, "media")
	if err := os.MkdirAll(mediaDir, 0755); err != nil {
		t.Fatalf("failed creating test media dir: %v", err)
	}

	// 1. Create a file inside media directory
	internalFileName := "thumb_test.jpg"
	internalFilePath := filepath.Join(mediaDir, internalFileName)
	internalContent := []byte("fake-jpeg-data")
	if err := os.WriteFile(internalFilePath, internalContent, 0644); err != nil {
		t.Fatalf("failed creating internal media file: %v", err)
	}

	// 2. Create an external file on host filesystem (outside mediaDir)
	externalFilePath := filepath.Join(tempDir, "external_countdown.mp4")
	externalContent := []byte("fake-mp4-video-stream")
	if err := os.WriteFile(externalFilePath, externalContent, 0644); err != nil {
		t.Fatalf("failed creating external media file: %v", err)
	}

	serverInstance, err := server.New(server.Config{
		DataDirectory: dataDir,
	})
	if err != nil {
		t.Fatalf("failed creating server: %v", err)
	}

	handler := serverInstance.Handler()

	testScenarios := []struct {
		name           string
		requestURL     string
		expectedStatus int
		expectedBody   []byte
	}{
		{
			name:           "internal media file served successfully",
			requestURL:     "/api/media/files/" + internalFileName,
			expectedStatus: http.StatusOK,
			expectedBody:   internalContent,
		},
		{
			name:           "external host filesystem file served successfully",
			requestURL:     "/api/media/files/" + strings.TrimPrefix(externalFilePath, "/"),
			expectedStatus: http.StatusOK,
			expectedBody:   externalContent,
		},
		{
			name:           "nonexistent media file returns 404",
			requestURL:     "/api/media/files/does_not_exist.mp4",
			expectedStatus: http.StatusNotFound,
			expectedBody:   nil,
		},
		{
			name:           "empty media file path returns 404",
			requestURL:     "/api/media/files/",
			expectedStatus: http.StatusNotFound,
			expectedBody:   nil,
		},
	}

	for _, scenario := range testScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, scenario.requestURL, nil)

			handler.ServeHTTP(recorder, request)

			if recorder.Code != scenario.expectedStatus {
				t.Errorf("status = %d, want %d", recorder.Code, scenario.expectedStatus)
			}

			if scenario.expectedBody != nil {
				bodyBytes, _ := io.ReadAll(recorder.Body)
				if string(bodyBytes) != string(scenario.expectedBody) {
					t.Errorf("body = %q, want %q", string(bodyBytes), string(scenario.expectedBody))
				}
			}
		})
	}
}

func TestServer_DocumentationEndpointsTable(t *testing.T) {
	tempDir := t.TempDir()
	dataDir := filepath.Join(tempDir, "data")

	serverInstance, err := server.New(server.Config{
		DataDirectory: dataDir,
	})
	if err != nil {
		t.Fatalf("failed creating server: %v", err)
	}

	handler := serverInstance.Handler()

	testScenarios := []struct {
		name                string
		requestURL          string
		expectedStatus      int
		expectedContentType string
		expectedContains    string
	}{
		{
			name:                "openapi spec json endpoint returns 200 with valid spec",
			requestURL:          "/docs/openapi.json",
			expectedStatus:      http.StatusOK,
			expectedContentType: "application/json",
			expectedContains:    `"openapi": "3.0.3"`,
		},
		{
			name:                "scalar docs page returns 200 with interactive ui",
			requestURL:          "/docs",
			expectedStatus:      http.StatusOK,
			expectedContentType: "text/html",
			expectedContains:    "@scalar/api-reference",
		},
		{
			name:                "scalar docs page with trailing slash returns 200",
			requestURL:          "/docs/",
			expectedStatus:      http.StatusOK,
			expectedContentType: "text/html",
			expectedContains:    `data-url="/docs/openapi.json"`,
		},
	}

	for _, scenario := range testScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, scenario.requestURL, nil)

			handler.ServeHTTP(recorder, request)

			if recorder.Code != scenario.expectedStatus {
				t.Errorf("status = %d, want %d", recorder.Code, scenario.expectedStatus)
			}

			contentType := recorder.Header().Get("Content-Type")
			if !strings.Contains(contentType, scenario.expectedContentType) {
				t.Errorf("contentType = %q, want containing %q", contentType, scenario.expectedContentType)
			}

			bodyBytes, _ := io.ReadAll(recorder.Body)
			if !strings.Contains(string(bodyBytes), scenario.expectedContains) {
				t.Errorf("body missing %q", scenario.expectedContains)
			}
		})
	}
}
