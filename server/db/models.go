// Package db provides lightweight, pure-Go SQLite persistence for Vidonyx projects, media assets, and render jobs.
package db

import (
	"time"
)

// ProjectRecord represents a persisted video project in the database.
type ProjectRecord struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Width           int       `json:"width"`
	Height          int       `json:"height"`
	FrameRate       float64   `json:"frame_rate"`
	BackgroundColor string    `json:"background_color"`
	Specification   string    `json:"specification"` // JSON string of spec.VideoSpec
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// MediaAssetRecord represents an uploaded audio, video, or image file in the media library.
type MediaAssetRecord struct {
	ID            string    `json:"id"`
	ProjectID     string    `json:"project_id"`
	FileName      string    `json:"file_name"`
	FilePath      string    `json:"file_path"`
	FileType      string    `json:"file_type"` // "video", "audio", "image"
	FileSizeBytes int64     `json:"file_size_bytes"`
	Duration      float64   `json:"duration_seconds"`
	Width         int       `json:"width"`
	Height        int       `json:"height"`
	FrameRate     float64   `json:"frame_rate"`
	SampleRate    int       `json:"sample_rate"`
	Channels      int       `json:"channels"`
	ThumbnailPath string    `json:"thumbnail_path"`
	WaveformData  string    `json:"waveform_data"` // JSON array of normalized peak levels [0.0 - 1.0]
	CreatedAt     time.Time `json:"created_at"`
}

// RenderJobRecord represents an asynchronous video rendering task.
type RenderJobRecord struct {
	ID                 string     `json:"id"`
	ProjectID          string     `json:"project_id"`
	Status             string     `json:"status"` // "pending", "rendering", "completed", "failed", "cancelled"
	ProgressPercentage float64    `json:"progress_percentage"`
	CurrentFrame       int64      `json:"current_frame"`
	CurrentFPS         float64    `json:"current_fps"`
	CurrentTimeSeconds float64    `json:"current_time_seconds"`
	RenderSpeed        float64    `json:"render_speed"`
	OutputPath         string     `json:"output_path"`
	ErrorMessage       string     `json:"error_message,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	FinishedAt         *time.Time `json:"finished_at,omitempty"`
}
