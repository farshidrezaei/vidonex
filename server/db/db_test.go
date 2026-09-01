package db_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/farshidrezaei/vidonyx/server/db"
)

func TestDatabase_Table(t *testing.T) {
	tempDirectory := t.TempDir()
	databasePath := filepath.Join(tempDirectory, "test_vidonyx.db")

	databaseInstance, err := db.Open(databasePath)
	if err != nil {
		t.Fatalf("failed opening test database: %v", err)
	}
	defer func() { _ = databaseInstance.Close() }()

	ctx := context.Background()

	t.Run("Projects CRUD Table", func(t *testing.T) {
		testCases := []struct {
			name        string
			project     db.ProjectRecord
			updateName  string
			shouldError bool
		}{
			{
				name: "Standard 1080p Project",
				project: db.ProjectRecord{
					ID:              "proj_test_1",
					Name:            "YouTube Video 1080p",
					Description:     "Test composition",
					Width:           1920,
					Height:          1080,
					FrameRate:       30.0,
					BackgroundColor: "#000000",
					Specification:   `{"version":"1.0"}`,
					CreatedAt:       time.Now().UTC(),
					UpdatedAt:       time.Now().UTC(),
				},
				updateName:  "YouTube Video 1080p Updated",
				shouldError: false,
			},
			{
				name: "TikTok 9:16 Vertical Project",
				project: db.ProjectRecord{
					ID:              "proj_test_2",
					Name:            "TikTok Clip",
					Description:     "Vertical short",
					Width:           1080,
					Height:          1920,
					FrameRate:       60.0,
					BackgroundColor: "#111827",
					Specification:   `{"version":"1.0"}`,
					CreatedAt:       time.Now().UTC(),
					UpdatedAt:       time.Now().UTC(),
				},
				updateName:  "TikTok Clip Final",
				shouldError: false,
			},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				// Create
				err := databaseInstance.CreateProject(ctx, testCase.project)
				if (err != nil) != testCase.shouldError {
					t.Fatalf("CreateProject() error = %v, wantErr %v", err, testCase.shouldError)
				}

				// Get
				retrieved, err := databaseInstance.GetProject(ctx, testCase.project.ID)
				if err != nil {
					t.Fatalf("GetProject() error = %v", err)
				}
				if retrieved == nil || retrieved.Name != testCase.project.Name {
					t.Fatalf("GetProject() returned unexpected project: %+v", retrieved)
				}

				// Update
				testCase.project.Name = testCase.updateName
				if err := databaseInstance.UpdateProject(ctx, testCase.project); err != nil {
					t.Fatalf("UpdateProject() error = %v", err)
				}

				// Verify Update
				updated, err := databaseInstance.GetProject(ctx, testCase.project.ID)
				if err != nil {
					t.Fatalf("GetProject() after update error = %v", err)
				}
				if updated.Name != testCase.updateName {
					t.Fatalf("expected updated name %q, got %q", testCase.updateName, updated.Name)
				}
			})
		}

		// List projects
		allProjects, err := databaseInstance.ListProjects(ctx)
		if err != nil {
			t.Fatalf("ListProjects() error = %v", err)
		}
		if len(allProjects) != 2 {
			t.Fatalf("expected 2 projects, got %d", len(allProjects))
		}

		// Delete project
		if err := databaseInstance.DeleteProject(ctx, "proj_test_1"); err != nil {
			t.Fatalf("DeleteProject() error = %v", err)
		}

		deleted, err := databaseInstance.GetProject(ctx, "proj_test_1")
		if err != nil {
			t.Fatalf("GetProject() after deletion error = %v", err)
		}
		if deleted != nil {
			t.Fatalf("expected nil for deleted project, got %+v", deleted)
		}
	})

	t.Run("Media Assets Table", func(t *testing.T) {
		asset := db.MediaAssetRecord{
			ID:            "asset_test_1",
			ProjectID:     "proj_test_2",
			FileName:      "intro.mp4",
			FilePath:      "asset_test_1.mp4",
			FileType:      "video",
			FileSizeBytes: 1048576,
			Duration:      12.5,
			Width:         1080,
			Height:        1920,
			FrameRate:     30.0,
			CreatedAt:     time.Now().UTC(),
		}

		if err := databaseInstance.SaveMediaAsset(ctx, asset); err != nil {
			t.Fatalf("SaveMediaAsset() error = %v", err)
		}

		assets, err := databaseInstance.ListMediaAssets(ctx, "proj_test_2")
		if err != nil {
			t.Fatalf("ListMediaAssets() error = %v", err)
		}
		if len(assets) != 1 || assets[0].FileName != "intro.mp4" {
			t.Fatalf("unexpected assets list: %+v", assets)
		}

		if err := databaseInstance.DeleteMediaAsset(ctx, "asset_test_1"); err != nil {
			t.Fatalf("DeleteMediaAsset() error = %v", err)
		}

		remaining, err := databaseInstance.ListMediaAssets(ctx, "proj_test_2")
		if err != nil {
			t.Fatalf("ListMediaAssets() after delete error = %v", err)
		}
		if len(remaining) != 0 {
			t.Fatalf("expected 0 assets after delete, got %d", len(remaining))
		}
	})

	t.Run("Render Jobs Table", func(t *testing.T) {
		job := db.RenderJobRecord{
			ID:                 "job_test_1",
			ProjectID:          "proj_test_2",
			Status:             "rendering",
			ProgressPercentage: 45.0,
			CurrentFrame:       450,
			CurrentFPS:         59.8,
			CurrentTimeSeconds: 7.5,
			RenderSpeed:        1.99,
			OutputPath:         "render_job_test_1.mp4",
			CreatedAt:          time.Now().UTC(),
		}

		if err := databaseInstance.SaveRenderJob(ctx, job); err != nil {
			t.Fatalf("SaveRenderJob() error = %v", err)
		}

		retrievedJob, err := databaseInstance.GetRenderJob(ctx, "job_test_1")
		if err != nil {
			t.Fatalf("GetRenderJob() error = %v", err)
		}
		if retrievedJob == nil || retrievedJob.Status != "rendering" {
			t.Fatalf("unexpected render job: %+v", retrievedJob)
		}

		if err := databaseInstance.CompleteRenderJob(ctx, "job_test_1", "completed", ""); err != nil {
			t.Fatalf("CompleteRenderJob() error = %v", err)
		}

		completedJob, err := databaseInstance.GetRenderJob(ctx, "job_test_1")
		if err != nil {
			t.Fatalf("GetRenderJob() after complete error = %v", err)
		}
		if completedJob.Status != "completed" || completedJob.FinishedAt == nil {
			t.Fatalf("expected completed status with finished timestamp, got %+v", completedJob)
		}
	})
}
