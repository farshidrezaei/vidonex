package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/farshidrezaei/vidonex/server/api"
	"github.com/farshidrezaei/vidonex/server/db"
	"github.com/farshidrezaei/vidonex/server/ws"
	"github.com/farshidrezaei/vidonex/spec"
)

func TestAPIEndpoints_Table(t *testing.T) {
	tempDirectory := t.TempDir()
	databasePath := filepath.Join(tempDirectory, "test_api.db")
	mediaDirectory := filepath.Join(tempDirectory, "media")
	exportDirectory := filepath.Join(tempDirectory, "exports")

	databaseInstance, err := db.Open(databasePath)
	if err != nil {
		t.Fatalf("failed opening test database: %v", err)
	}
	defer func() { _ = databaseInstance.Close() }()

	websocketHub := ws.NewHub(nil)
	apiHandlers := api.NewHandlers(databaseInstance, websocketHub, mediaDirectory, exportDirectory, nil)

	t.Run("Create and List Projects API", func(t *testing.T) {
		createPayload := api.CreateProjectRequest{
			Name:            "Test API Project",
			Description:     "Integration test",
			Width:           1920,
			Height:          1080,
			FrameRate:       30.0,
			BackgroundColor: "#000000",
		}
		bodyBytes, _ := json.Marshal(createPayload)

		req := httptest.NewRequest(http.MethodPost, "/api/projects", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		apiHandlers.HandleCreateProject(rr, req)

		if rr.Code != http.StatusCreated {
			t.Fatalf("expected status 201 Created, got %d: %s", rr.Code, rr.Body.String())
		}

		var createResp api.StandardResponse
		if err := json.NewDecoder(rr.Body).Decode(&createResp); err != nil || !createResp.Success {
			t.Fatalf("failed decoding response: %v", err)
		}

		// List
		listReq := httptest.NewRequest(http.MethodGet, "/api/projects", nil)
		listRR := httptest.NewRecorder()

		apiHandlers.HandleListProjects(listRR, listReq)

		if listRR.Code != http.StatusOK {
			t.Fatalf("expected status 200 OK, got %d", listRR.Code)
		}
	})

	t.Run("Validate Spec API Table", func(t *testing.T) {
		testCases := []struct {
			name      string
			spec      spec.VideoSpec
			wantValid bool
		}{
			{
				name: "Valid Specification",
				spec: spec.VideoSpec{
					Version: "1.0",
					Canvas: spec.CanvasSpec{
						Width:  1920,
						Height: 1080,
					},
					Tracks: []spec.TrackSpec{
						{
							ID:   "video_track_1",
							Kind: "video",
							Clips: []spec.ClipSpec{
								{
									ID:       "clip_1",
									Source:   "test.mp4",
									Start:    0,
									Duration: 5.0,
								},
							},
						},
					},
				},
				wantValid: true,
			},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				reqPayload := api.ValidateSpecRequest{
					Specification: &testCase.spec,
				}
				bodyBytes, _ := json.Marshal(reqPayload)

				req := httptest.NewRequest(http.MethodPost, "/api/spec/validate", bytes.NewReader(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()

				apiHandlers.HandleValidateSpec(rr, req)

				var resp api.ValidateSpecResponse
				if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
					t.Fatalf("failed decoding validate response: %v", err)
				}

				if resp.Valid != testCase.wantValid {
					t.Fatalf("expected Valid=%v, got %v (errors: %v)", testCase.wantValid, resp.Valid, resp.Errors)
				}
			})
		}
	})

	t.Run("Generate Graph API", func(t *testing.T) {
		reqPayload := api.GenerateGraphRequest{
			Specification: &spec.VideoSpec{
				Version: "1.0",
				Canvas: spec.CanvasSpec{
					Width:  1920,
					Height: 1080,
				},
				Tracks: []spec.TrackSpec{
					{
						ID:   "video_track_1",
						Kind: "video",
						Clips: []spec.ClipSpec{
							{
								ID:       "clip_1",
								Source:   "input.mp4",
								Start:    0,
								Duration: 5.0,
							},
						},
					},
				},
			},
			Format: "mermaid",
		}
		bodyBytes, _ := json.Marshal(reqPayload)

		req := httptest.NewRequest(http.MethodPost, "/api/spec/graph", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		apiHandlers.HandleGenerateGraph(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d: %s", rr.Code, rr.Body.String())
		}

		var resp api.GenerateGraphResponse
		if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
			t.Fatalf("failed decoding graph response: %v", err)
		}

		if resp.Format != "mermaid" || len(resp.Content) == 0 {
			t.Fatalf("unexpected graph response: %+v", resp)
		}
	})

	t.Run("Start Render API with Custom OutputPath", func(t *testing.T) {
		_ = databaseInstance.CreateProject(context.Background(), db.ProjectRecord{
			ID:              "proj_test_123",
			Name:            "Test Project",
			Width:           1280,
			Height:          720,
			FrameRate:       30,
			BackgroundColor: "#000000",
		})
		customOutput := filepath.Join(tempDirectory, "custom_export", "my_custom_video.mp4")
		renderPayload := api.StartRenderRequest{
			ProjectID: "proj_test_123",
			Specification: &spec.VideoSpec{
				Version: "1.0",
				Canvas: spec.CanvasSpec{
					Width:  1280,
					Height: 720,
				},
				Tracks: []spec.TrackSpec{
					{
						ID:    "video_track_1",
						Kind:  "video",
						Clips: []spec.ClipSpec{},
					},
				},
			},
			OutputFormat: "mp4",
			OutputPath:   customOutput,
		}
		bodyBytes, _ := json.Marshal(renderPayload)

		req := httptest.NewRequest(http.MethodPost, "/api/render/start", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		apiHandlers.HandleStartRender(rr, req)

		if rr.Code != http.StatusAccepted {
			t.Fatalf("expected status 202 Accepted, got %d: %s", rr.Code, rr.Body.String())
		}

		var resp api.StartRenderResponse
		if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
			t.Fatalf("failed decoding render response: %v", err)
		}

		if resp.OutputPath != "my_custom_video.mp4" {
			t.Fatalf("expected OutputPath 'my_custom_video.mp4', got %q", resp.OutputPath)
		}
	})
}
