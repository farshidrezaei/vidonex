package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	// Register pure-Go SQLite driver for standard database/sql driver interface
	_ "modernc.org/sqlite"
)

// Database encapsulates the SQLite database connection and operations for Vidonex.
type Database struct {
	connection *sql.DB
}

// Open initializes and opens a SQLite database connection at the specified file path,
// creating parent directories and applying migrations automatically.
func Open(databasePath string) (*Database, error) {
	if databasePath == "" {
		return nil, errors.New("database path cannot be empty")
	}

	databaseDirectory := filepath.Dir(databasePath)
	if err := os.MkdirAll(databaseDirectory, 0755); err != nil {
		return nil, fmt.Errorf("failed creating database directory: %w", err)
	}

	// SQLite connection string with WAL mode and busy timeout for concurrent safety
	dataStoreName := fmt.Sprintf("%s?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)", databasePath)
	databaseConnection, err := sql.Open("sqlite", dataStoreName)
	if err != nil {
		return nil, fmt.Errorf("failed opening sqlite database: %w", err)
	}

	databaseConnection.SetMaxOpenConns(1) // SQLite single-writer safe pool

	if err := databaseConnection.Ping(); err != nil {
		_ = databaseConnection.Close()
		return nil, fmt.Errorf("failed connecting to sqlite database: %w", err)
	}

	databaseInstance := &Database{
		connection: databaseConnection,
	}

	if err := databaseInstance.migrate(); err != nil {
		_ = databaseConnection.Close()
		return nil, fmt.Errorf("failed running database migrations: %w", err)
	}

	return databaseInstance, nil
}

// Close gracefully closes the database connection.
func (db *Database) Close() error {
	if db.connection != nil {
		return db.connection.Close()
	}
	return nil
}

// migrate creates the necessary database tables and indexes if they do not already exist.
func (db *Database) migrate() error {
	schemaQueries := `
	CREATE TABLE IF NOT EXISTS projects (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT NOT NULL DEFAULT '',
		width INTEGER NOT NULL DEFAULT 1920,
		height INTEGER NOT NULL DEFAULT 1080,
		frame_rate REAL NOT NULL DEFAULT 30.0,
		background_color TEXT NOT NULL DEFAULT '#000000',
		specification TEXT NOT NULL DEFAULT '{}',
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);

	CREATE TABLE IF NOT EXISTS media_assets (
		id TEXT PRIMARY KEY,
		project_id TEXT NOT NULL,
		file_name TEXT NOT NULL,
		file_path TEXT NOT NULL,
		file_type TEXT NOT NULL,
		file_size_bytes INTEGER NOT NULL DEFAULT 0,
		duration_seconds REAL NOT NULL DEFAULT 0.0,
		width INTEGER NOT NULL DEFAULT 0,
		height INTEGER NOT NULL DEFAULT 0,
		frame_rate REAL NOT NULL DEFAULT 0.0,
		sample_rate INTEGER NOT NULL DEFAULT 0,
		channels INTEGER NOT NULL DEFAULT 0,
		thumbnail_path TEXT NOT NULL DEFAULT '',
		waveform_data TEXT NOT NULL DEFAULT '[]',
		created_at TIMESTAMP NOT NULL,
		FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_media_assets_project_id ON media_assets(project_id);

	CREATE TABLE IF NOT EXISTS render_jobs (
		id TEXT PRIMARY KEY,
		project_id TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'pending',
		progress_percentage REAL NOT NULL DEFAULT 0.0,
		current_frame INTEGER NOT NULL DEFAULT 0,
		current_fps REAL NOT NULL DEFAULT 0.0,
		current_time_seconds REAL NOT NULL DEFAULT 0.0,
		render_speed REAL NOT NULL DEFAULT 0.0,
		output_path TEXT NOT NULL,
		error_message TEXT NOT NULL DEFAULT '',
		created_at TIMESTAMP NOT NULL,
		finished_at TIMESTAMP,
		FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_render_jobs_project_id ON render_jobs(project_id);
	`

	_, err := db.connection.Exec(schemaQueries)
	return err
}

// CreateProject inserts a new project record into the database.
func (db *Database) CreateProject(ctx context.Context, project ProjectRecord) error {
	query := `
	INSERT INTO projects (id, name, description, width, height, frame_rate, background_color, specification, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.connection.ExecContext(
		ctx,
		query,
		project.ID,
		project.Name,
		project.Description,
		project.Width,
		project.Height,
		project.FrameRate,
		project.BackgroundColor,
		project.Specification,
		project.CreatedAt,
		project.UpdatedAt,
	)
	return err
}

// GetProject retrieves a single project record by its unique identifier.
func (db *Database) GetProject(ctx context.Context, projectIdentifier string) (*ProjectRecord, error) {
	query := `
	SELECT id, name, description, width, height, frame_rate, background_color, specification, created_at, updated_at
	FROM projects
	WHERE id = ?
	`
	row := db.connection.QueryRowContext(ctx, query, projectIdentifier)

	var project ProjectRecord
	err := row.Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.Width,
		&project.Height,
		&project.FrameRate,
		&project.BackgroundColor,
		&project.Specification,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}

// ListProjects retrieves all project records ordered by last updated timestamp descending.
func (db *Database) ListProjects(ctx context.Context) ([]ProjectRecord, error) {
	query := `
	SELECT id, name, description, width, height, frame_rate, background_color, specification, created_at, updated_at
	FROM projects
	ORDER BY updated_at DESC
	`
	rows, err := db.connection.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var projects []ProjectRecord
	for rows.Next() {
		var project ProjectRecord
		if err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.Width,
			&project.Height,
			&project.FrameRate,
			&project.BackgroundColor,
			&project.Specification,
			&project.CreatedAt,
			&project.UpdatedAt,
		); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return projects, nil
}

// UpdateProject updates an existing project's metadata and specification JSON.
func (db *Database) UpdateProject(ctx context.Context, project ProjectRecord) error {
	query := `
	UPDATE projects
	SET name = ?, description = ?, width = ?, height = ?, frame_rate = ?, background_color = ?, specification = ?, updated_at = ?
	WHERE id = ?
	`
	result, err := db.connection.ExecContext(
		ctx,
		query,
		project.Name,
		project.Description,
		project.Width,
		project.Height,
		project.FrameRate,
		project.BackgroundColor,
		project.Specification,
		time.Now().UTC(),
		project.ID,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("project not found")
	}
	return nil
}

// DeleteProject removes a project record and all cascading associated assets and jobs.
func (db *Database) DeleteProject(ctx context.Context, projectIdentifier string) error {
	query := `DELETE FROM projects WHERE id = ?`
	_, err := db.connection.ExecContext(ctx, query, projectIdentifier)
	return err
}

// SaveMediaAsset inserts or updates a media asset record.
func (db *Database) SaveMediaAsset(ctx context.Context, asset MediaAssetRecord) error {
	query := `
	INSERT INTO media_assets (id, project_id, file_name, file_path, file_type, file_size_bytes, duration_seconds, width, height, frame_rate, sample_rate, channels, thumbnail_path, waveform_data, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		file_name = excluded.file_name,
		file_path = excluded.file_path,
		file_size_bytes = excluded.file_size_bytes,
		duration_seconds = excluded.duration_seconds,
		width = excluded.width,
		height = excluded.height,
		frame_rate = excluded.frame_rate,
		sample_rate = excluded.sample_rate,
		channels = excluded.channels,
		thumbnail_path = excluded.thumbnail_path,
		waveform_data = excluded.waveform_data
	`
	_, err := db.connection.ExecContext(
		ctx,
		query,
		asset.ID,
		asset.ProjectID,
		asset.FileName,
		asset.FilePath,
		asset.FileType,
		asset.FileSizeBytes,
		asset.Duration,
		asset.Width,
		asset.Height,
		asset.FrameRate,
		asset.SampleRate,
		asset.Channels,
		asset.ThumbnailPath,
		asset.WaveformData,
		asset.CreatedAt,
	)
	return err
}

// ListMediaAssets retrieves all media assets associated with a given project.
func (db *Database) ListMediaAssets(ctx context.Context, projectIdentifier string) ([]MediaAssetRecord, error) {
	query := `
	SELECT id, project_id, file_name, file_path, file_type, file_size_bytes, duration_seconds, width, height, frame_rate, sample_rate, channels, thumbnail_path, waveform_data, created_at
	FROM media_assets
	WHERE project_id = ?
	ORDER BY created_at DESC
	`
	rows, err := db.connection.QueryContext(ctx, query, projectIdentifier)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var assets []MediaAssetRecord
	for rows.Next() {
		var asset MediaAssetRecord
		if err := rows.Scan(
			&asset.ID,
			&asset.ProjectID,
			&asset.FileName,
			&asset.FilePath,
			&asset.FileType,
			&asset.FileSizeBytes,
			&asset.Duration,
			&asset.Width,
			&asset.Height,
			&asset.FrameRate,
			&asset.SampleRate,
			&asset.Channels,
			&asset.ThumbnailPath,
			&asset.WaveformData,
			&asset.CreatedAt,
		); err != nil {
			return nil, err
		}
		assets = append(assets, asset)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return assets, nil
}

// DeleteMediaAsset deletes an asset by its identifier.
func (db *Database) DeleteMediaAsset(ctx context.Context, assetIdentifier string) error {
	query := `DELETE FROM media_assets WHERE id = ?`
	_, err := db.connection.ExecContext(ctx, query, assetIdentifier)
	return err
}

// SaveRenderJob inserts a new render job into the database.
func (db *Database) SaveRenderJob(ctx context.Context, job RenderJobRecord) error {
	query := `
	INSERT INTO render_jobs (id, project_id, status, progress_percentage, current_frame, current_fps, current_time_seconds, render_speed, output_path, error_message, created_at, finished_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.connection.ExecContext(
		ctx,
		query,
		job.ID,
		job.ProjectID,
		job.Status,
		job.ProgressPercentage,
		job.CurrentFrame,
		job.CurrentFPS,
		job.CurrentTimeSeconds,
		job.RenderSpeed,
		job.OutputPath,
		job.ErrorMessage,
		job.CreatedAt,
		job.FinishedAt,
	)
	return err
}

// UpdateRenderJobProgress updates live telemetry data for an active render job.
func (db *Database) UpdateRenderJobProgress(ctx context.Context, jobIdentifier string, percentage float64, currentFrame int64, currentFPS float64, currentTime float64, speed float64) error {
	query := `
	UPDATE render_jobs
	SET progress_percentage = ?, current_frame = ?, current_fps = ?, current_time_seconds = ?, render_speed = ?
	WHERE id = ?
	`
	_, err := db.connection.ExecContext(ctx, query, percentage, currentFrame, currentFPS, currentTime, speed, jobIdentifier)
	return err
}

// CompleteRenderJob marks a render job as completed or failed.
func (db *Database) CompleteRenderJob(ctx context.Context, jobIdentifier string, status string, errorMessage string) error {
	now := time.Now().UTC()
	query := `
	UPDATE render_jobs
	SET status = ?, error_message = ?, finished_at = ?
	WHERE id = ?
	`
	_, err := db.connection.ExecContext(ctx, query, status, errorMessage, now, jobIdentifier)
	return err
}

// GetRenderJob retrieves a render job by its unique identifier.
func (db *Database) GetRenderJob(ctx context.Context, jobIdentifier string) (*RenderJobRecord, error) {
	query := `
	SELECT id, project_id, status, progress_percentage, current_frame, current_fps, current_time_seconds, render_speed, output_path, error_message, created_at, finished_at
	FROM render_jobs
	WHERE id = ?
	`
	row := db.connection.QueryRowContext(ctx, query, jobIdentifier)

	var job RenderJobRecord
	err := row.Scan(
		&job.ID,
		&job.ProjectID,
		&job.Status,
		&job.ProgressPercentage,
		&job.CurrentFrame,
		&job.CurrentFPS,
		&job.CurrentTimeSeconds,
		&job.RenderSpeed,
		&job.OutputPath,
		&job.ErrorMessage,
		&job.CreatedAt,
		&job.FinishedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &job, nil
}
