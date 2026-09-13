// Package api provides HTTP REST API route handlers and request/response payloads for the Vidonex Web Studio.
package api

import (
	"github.com/farshidrezaei/vidonex/presets"
	"github.com/farshidrezaei/vidonex/spec"
)

// CreateProjectRequest defines payload for initializing a new project.
type CreateProjectRequest struct {
	Name            string          `json:"name"`
	Description     string          `json:"description,omitempty"`
	Width           int             `json:"width,omitempty"`
	Height          int             `json:"height,omitempty"`
	FrameRate       float64         `json:"frame_rate,omitempty"`
	BackgroundColor string          `json:"background_color,omitempty"`
	Specification   *spec.VideoSpec `json:"specification,omitempty"`
}

// UpdateProjectRequest defines payload for auto-saving and updating a project.
type UpdateProjectRequest struct {
	Name            string          `json:"name"`
	Description     string          `json:"description,omitempty"`
	Width           int             `json:"width"`
	Height          int             `json:"height"`
	FrameRate       float64         `json:"frame_rate"`
	BackgroundColor string          `json:"background_color"`
	Specification   *spec.VideoSpec `json:"specification"`
}

// ValidateSpecRequest contains a specification payload to validate.
type ValidateSpecRequest struct {
	Specification *spec.VideoSpec `json:"specification"`
}

// ValidateSpecResponse contains validation result and any diagnostics.
type ValidateSpecResponse struct {
	Valid   bool     `json:"valid"`
	Errors  []string `json:"errors,omitempty"`
	Message string   `json:"message"`
}

// GenerateGraphRequest specifies graph visualization options.
type GenerateGraphRequest struct {
	Specification *spec.VideoSpec `json:"specification"`
	Format        string          `json:"format"` // "mermaid" or "dot"
}

// GenerateGraphResponse contains the generated graph definition.
type GenerateGraphResponse struct {
	Format  string `json:"format"`
	Content string `json:"content"`
}

// StartRenderRequest configures and initiates a background video render job.
type StartRenderRequest struct {
	ProjectID            string                     `json:"project_id"`
	Specification        *spec.VideoSpec            `json:"specification"`
	OutputFormat         string                     `json:"output_format,omitempty"` // "mp4", "mkv", "mov", "webm", "gif"
	Preset               string                     `json:"preset,omitempty"`
	HardwareAcceleration presets.HardwareAccelerator `json:"hardware_acceleration,omitempty"`
}

// StartRenderResponse returns the dispatched job identifier.
type StartRenderResponse struct {
	JobID      string `json:"job_id"`
	Status     string `json:"status"`
	OutputPath string `json:"output_path"`
}

// StandardResponse standardizes JSON envelope for error and success responses.
type StandardResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}
