package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/farshidrezaei/vidonyx/compiler"
	"github.com/farshidrezaei/vidonyx/composer"
	"github.com/farshidrezaei/vidonyx/executor"
	"github.com/farshidrezaei/vidonyx/presets"
	"github.com/farshidrezaei/vidonyx/probe"
	"github.com/farshidrezaei/vidonyx/server/db"
	"github.com/farshidrezaei/vidonyx/server/ws"
	"github.com/farshidrezaei/vidonyx/spec"
	"github.com/farshidrezaei/vidonyx/timeline"
)

// Handlers bundles all route dependencies.
type Handlers struct {
	database       *db.Database
	websocketHub   *ws.Hub
	mediaDirectory string
	exportDirectory string
	prober         probe.MediaProber
	logger         *slog.Logger

	activeRendersMutex sync.Mutex
	activeRenderCancels map[string]context.CancelFunc
}

// NewHandlers creates an initialized API Handlers instance.
func NewHandlers(database *db.Database, websocketHub *ws.Hub, mediaDirectory string, exportDirectory string, logger *slog.Logger) *Handlers {
	if logger == nil {
		logger = slog.Default()
	}
	_ = os.MkdirAll(mediaDirectory, 0755)
	_ = os.MkdirAll(exportDirectory, 0755)

	return &Handlers{
		database:            database,
		websocketHub:        websocketHub,
		mediaDirectory:      mediaDirectory,
		exportDirectory:     exportDirectory,
		prober:              probe.NewCachedProber(probe.NewFFprobeProber("ffprobe")),
		logger:              logger,
		activeRenderCancels: make(map[string]context.CancelFunc),
	}
}

func generateRandomID(prefix string) string {
	bytes := make([]byte, 8)
	_, _ = rand.Read(bytes)
	return fmt.Sprintf("%s_%s", prefix, hex.EncodeToString(bytes))
}

func respondJSON(responseWriter http.ResponseWriter, statusCode int, payload any) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)
	_ = json.NewEncoder(responseWriter).Encode(payload)
}

func respondError(responseWriter http.ResponseWriter, statusCode int, message string) {
	respondJSON(responseWriter, statusCode, StandardResponse{
		Success: false,
		Message: message,
	})
}

// HandleListProjects returns all saved projects.
func (h *Handlers) HandleListProjects(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	projects, err := h.database.ListProjects(httpRequest.Context())
	if err != nil {
		h.logger.Error("failed listing projects", "error", err)
		respondError(responseWriter, http.StatusInternalServerError, "failed listing projects")
		return
	}
	respondJSON(responseWriter, http.StatusOK, StandardResponse{
		Success: true,
		Data:    projects,
	})
}

// HandleCreateProject creates and saves a new project.
func (h *Handlers) HandleCreateProject(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	var request CreateProjectRequest
	if err := json.NewDecoder(httpRequest.Body).Decode(&request); err != nil {
		respondError(responseWriter, http.StatusBadRequest, "invalid request body")
		return
	}

	if request.Name == "" {
		request.Name = "Untitled Project"
	}
	if request.Width <= 0 {
		request.Width = 1920
	}
	if request.Height <= 0 {
		request.Height = 1080
	}
	if request.FrameRate <= 0 {
		request.FrameRate = 30.0
	}
	if request.BackgroundColor == "" {
		request.BackgroundColor = "#000000"
	}

	initialSpec := request.Specification
	if initialSpec == nil {
		initialSpec = &spec.VideoSpec{
			Version: "1.0",
			Canvas: spec.CanvasSpec{
				Width:           request.Width,
				Height:          request.Height,
				BackgroundColor: request.BackgroundColor,
			},
			FPS:    request.FrameRate,
			Tracks: []spec.TrackSpec{},
		}
	}

	specBytes, err := json.Marshal(initialSpec)
	if err != nil {
		respondError(responseWriter, http.StatusInternalServerError, "failed marshaling specification")
		return
	}

	now := time.Now().UTC()
	projectRecord := db.ProjectRecord{
		ID:              generateRandomID("proj"),
		Name:            request.Name,
		Description:     request.Description,
		Width:           request.Width,
		Height:          request.Height,
		FrameRate:       request.FrameRate,
		BackgroundColor: request.BackgroundColor,
		Specification:   string(specBytes),
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := h.database.CreateProject(httpRequest.Context(), projectRecord); err != nil {
		h.logger.Error("failed creating project", "error", err)
		respondError(responseWriter, http.StatusInternalServerError, "failed creating project")
		return
	}

	respondJSON(responseWriter, http.StatusCreated, StandardResponse{
		Success: true,
		Data:    projectRecord,
	})
}

// HandleGetProject retrieves a project by ID.
func (h *Handlers) HandleGetProject(responseWriter http.ResponseWriter, httpRequest *http.Request, projectID string) {
	project, err := h.database.GetProject(httpRequest.Context(), projectID)
	if err != nil {
		h.logger.Error("failed retrieving project", "error", err, "project_id", projectID)
		respondError(responseWriter, http.StatusInternalServerError, "failed retrieving project")
		return
	}
	if project == nil {
		respondError(responseWriter, http.StatusNotFound, "project not found")
		return
	}

	respondJSON(responseWriter, http.StatusOK, StandardResponse{
		Success: true,
		Data:    project,
	})
}

// HandleUpdateProject updates a project's metadata and specification.
func (h *Handlers) HandleUpdateProject(responseWriter http.ResponseWriter, httpRequest *http.Request, projectID string) {
	var request UpdateProjectRequest
	if err := json.NewDecoder(httpRequest.Body).Decode(&request); err != nil {
		respondError(responseWriter, http.StatusBadRequest, "invalid request body")
		return
	}

	specBytes, err := json.Marshal(request.Specification)
	if err != nil {
		respondError(responseWriter, http.StatusBadRequest, "invalid specification")
		return
	}

	projectRecord := db.ProjectRecord{
		ID:              projectID,
		Name:            request.Name,
		Description:     request.Description,
		Width:           request.Width,
		Height:          request.Height,
		FrameRate:       request.FrameRate,
		BackgroundColor: request.BackgroundColor,
		Specification:   string(specBytes),
		UpdatedAt:       time.Now().UTC(),
	}

	if err := h.database.UpdateProject(httpRequest.Context(), projectRecord); err != nil {
		h.logger.Error("failed updating project", "error", err, "project_id", projectID)
		respondError(responseWriter, http.StatusInternalServerError, "failed updating project")
		return
	}

	respondJSON(responseWriter, http.StatusOK, StandardResponse{
		Success: true,
		Message: "project updated successfully",
	})
}

// HandleDeleteProject removes a project.
func (h *Handlers) HandleDeleteProject(responseWriter http.ResponseWriter, httpRequest *http.Request, projectID string) {
	if err := h.database.DeleteProject(httpRequest.Context(), projectID); err != nil {
		h.logger.Error("failed deleting project", "error", err, "project_id", projectID)
		respondError(responseWriter, http.StatusInternalServerError, "failed deleting project")
		return
	}

	respondJSON(responseWriter, http.StatusOK, StandardResponse{
		Success: true,
		Message: "project deleted successfully",
	})
}

// HandleUploadMedia handles multi-part file uploads, probes metadata, and generates thumbnails.
func (h *Handlers) HandleUploadMedia(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	err := httpRequest.ParseMultipartForm(500 << 20) // 500 MB max memory/temp
	if err != nil {
		respondError(responseWriter, http.StatusBadRequest, "failed parsing upload form")
		return
	}

	projectID := httpRequest.FormValue("project_id")
	if projectID == "" {
		respondError(responseWriter, http.StatusBadRequest, "missing project_id")
		return
	}

	file, header, err := httpRequest.FormFile("file")
	if err != nil {
		respondError(responseWriter, http.StatusBadRequest, "missing file in form")
		return
	}
	defer func() { _ = file.Close() }()

	assetID := generateRandomID("asset")
	fileExtension := strings.ToLower(filepath.Ext(header.Filename))
	destinationFileName := fmt.Sprintf("%s%s", assetID, fileExtension)
	destinationFilePath := filepath.Join(h.mediaDirectory, destinationFileName)

	destinationFile, err := os.Create(destinationFilePath)
	if err != nil {
		h.logger.Error("failed creating asset file", "error", err)
		respondError(responseWriter, http.StatusInternalServerError, "failed saving file")
		return
	}
	defer func() { _ = destinationFile.Close() }()

	fileSize, err := io.Copy(destinationFile, file)
	if err != nil {
		h.logger.Error("failed copying uploaded file", "error", err)
		respondError(responseWriter, http.StatusInternalServerError, "failed writing file")
		return
	}

	fileType := "video"
	switch fileExtension {
	case ".mp3", ".wav", ".aac", ".flac", ".ogg", ".m4a":
		fileType = "audio"
	case ".png", ".jpg", ".jpeg", ".webp", ".svg", ".bmp", ".gif":
		fileType = "image"
	}

	// Probe stream properties
	var durationSeconds float64
	var width, height, sampleRate, channels int
	var frameRate float64

	probeCtx, cancel := context.WithTimeout(httpRequest.Context(), 15*time.Second)
	defer cancel()

	mediaMetadata, probeErr := h.prober.Probe(probeCtx, destinationFilePath)
	if probeErr == nil && mediaMetadata != nil {
		durationSeconds = mediaMetadata.Duration.Seconds()
		if mediaMetadata.VideoStream != nil {
			width = mediaMetadata.VideoStream.Width
			height = mediaMetadata.VideoStream.Height
			frameRate = mediaMetadata.VideoStream.FPS.Seconds()
		}
		if mediaMetadata.AudioStream != nil {
			sampleRate = mediaMetadata.AudioStream.SampleRate
			channels = mediaMetadata.AudioStream.Channels
		}
	}

	// Extract thumbnail for video/image
	thumbnailPath := ""
	switch fileType {
	case "video":
		thumbFileName := fmt.Sprintf("thumb_%s.jpg", assetID)
		thumbFullPath := filepath.Join(h.mediaDirectory, thumbFileName)
		cmd := exec.CommandContext(probeCtx, "ffmpeg", "-y", "-ss", "00:00:00.500", "-i", destinationFilePath, "-vframes", "1", "-vf", "scale=320:-1", thumbFullPath)
		if cmd.Run() == nil {
			thumbnailPath = thumbFileName
		}
	case "image":
		thumbnailPath = destinationFileName
	}

	assetRecord := db.MediaAssetRecord{
		ID:            assetID,
		ProjectID:     projectID,
		FileName:      header.Filename,
		FilePath:      destinationFileName,
		FileType:      fileType,
		FileSizeBytes: fileSize,
		Duration:      durationSeconds,
		Width:         width,
		Height:        height,
		FrameRate:     frameRate,
		SampleRate:    sampleRate,
		Channels:      channels,
		ThumbnailPath: thumbnailPath,
		WaveformData:  "[]",
		CreatedAt:     time.Now().UTC(),
	}

	if err := h.database.SaveMediaAsset(httpRequest.Context(), assetRecord); err != nil {
		h.logger.Error("failed saving media asset record", "error", err)
		respondError(responseWriter, http.StatusInternalServerError, "failed saving asset metadata")
		return
	}

	respondJSON(responseWriter, http.StatusCreated, StandardResponse{
		Success: true,
		Data:    assetRecord,
	})
}

// HandleListAssets returns all media assets for a project.
func (h *Handlers) HandleListAssets(responseWriter http.ResponseWriter, httpRequest *http.Request, projectID string) {
	assets, err := h.database.ListMediaAssets(httpRequest.Context(), projectID)
	if err != nil {
		h.logger.Error("failed listing assets", "error", err, "project_id", projectID)
		respondError(responseWriter, http.StatusInternalServerError, "failed listing assets")
		return
	}

	respondJSON(responseWriter, http.StatusOK, StandardResponse{
		Success: true,
		Data:    assets,
	})
}

// HandleDeleteAsset removes an asset from database and disk.
func (h *Handlers) HandleDeleteAsset(responseWriter http.ResponseWriter, httpRequest *http.Request, assetID string) {
	if err := h.database.DeleteMediaAsset(httpRequest.Context(), assetID); err != nil {
		h.logger.Error("failed deleting asset", "error", err, "asset_id", assetID)
		respondError(responseWriter, http.StatusInternalServerError, "failed deleting asset")
		return
	}

	respondJSON(responseWriter, http.StatusOK, StandardResponse{
		Success: true,
		Message: "asset deleted successfully",
	})
}

// HandleValidateSpec validates the declarative timeline specification.
func (h *Handlers) HandleValidateSpec(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	var request ValidateSpecRequest
	if err := json.NewDecoder(httpRequest.Body).Decode(&request); err != nil || request.Specification == nil {
		respondError(responseWriter, http.StatusBadRequest, "invalid specification payload")
		return
	}

	compositionTimeline, _, _, err := spec.ToTimeline(request.Specification)
	if err != nil {
		respondJSON(responseWriter, http.StatusOK, ValidateSpecResponse{
			Valid:   false,
			Errors:  []string{err.Error()},
			Message: "specification conversion error",
		})
		return
	}

	if err := timeline.Validate(compositionTimeline); err != nil {
		respondJSON(responseWriter, http.StatusOK, ValidateSpecResponse{
			Valid:   false,
			Errors:  []string{err.Error()},
			Message: "timeline validation failed",
		})
		return
	}

	respondJSON(responseWriter, http.StatusOK, ValidateSpecResponse{
		Valid:   true,
		Message: "specification is valid",
	})
}

// HandleGenerateGraph returns the Mermaid or DOT diagram of the compiled filtergraph DAG.
func (h *Handlers) HandleGenerateGraph(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	var request GenerateGraphRequest
	if err := json.NewDecoder(httpRequest.Body).Decode(&request); err != nil || request.Specification == nil {
		respondError(responseWriter, http.StatusBadRequest, "invalid request body")
		return
	}

	compositionTimeline, _, _, err := spec.ToTimeline(request.Specification)
	if err != nil {
		respondError(responseWriter, http.StatusBadRequest, fmt.Sprintf("failed converting spec to timeline: %v", err))
		return
	}

	graphCompiler := compiler.New(h.logger)
	compilationResult, err := graphCompiler.Compile(compositionTimeline, "output.mp4")
	if err != nil {
		respondError(responseWriter, http.StatusBadRequest, fmt.Sprintf("failed compiling filtergraph: %v", err))
		return
	}

	format := strings.ToLower(request.Format)
	var content string
	if format == "dot" {
		content, _ = compilationResult.DOT()
	} else {
		format = "mermaid"
		content, _ = compilationResult.Mermaid()
	}

	respondJSON(responseWriter, http.StatusOK, GenerateGraphResponse{
		Format:  format,
		Content: content,
	})
}

// HandleStartRender initiates background video rendering and broadcasts progress via WebSocket.
func (h *Handlers) HandleStartRender(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	var request StartRenderRequest
	if err := json.NewDecoder(httpRequest.Body).Decode(&request); err != nil || request.Specification == nil {
		respondError(responseWriter, http.StatusBadRequest, "invalid render request payload")
		return
	}

	jobID := generateRandomID("job")
	outputExtension := ".mp4"
	if request.OutputFormat != "" {
		outputExtension = "." + strings.TrimPrefix(request.OutputFormat, ".")
	}

	outputFileName := fmt.Sprintf("render_%s%s", jobID, outputExtension)
	outputFilePath := filepath.Join(h.exportDirectory, outputFileName)

	// Resolve media sources in spec to absolute paths
	for trackIndex := range request.Specification.Tracks {
		for clipIndex := range request.Specification.Tracks[trackIndex].Clips {
			clipSource := request.Specification.Tracks[trackIndex].Clips[clipIndex].Source
			if clipSource != "" && !filepath.IsAbs(clipSource) {
				request.Specification.Tracks[trackIndex].Clips[clipIndex].Source = filepath.Join(h.mediaDirectory, clipSource)
			}
		}
	}

	compositionTimeline, _, _, err := spec.ToTimeline(request.Specification)
	if err != nil {
		respondError(responseWriter, http.StatusBadRequest, fmt.Sprintf("failed converting spec to timeline: %v", err))
		return
	}

	renderJobRecord := db.RenderJobRecord{
		ID:                 jobID,
		ProjectID:          request.ProjectID,
		Status:             "rendering",
		ProgressPercentage: 0.0,
		OutputPath:         outputFileName,
		CreatedAt:          time.Now().UTC(),
	}

	if err := h.database.SaveRenderJob(httpRequest.Context(), renderJobRecord); err != nil {
		h.logger.Error("failed saving render job", "error", err)
		respondError(responseWriter, http.StatusInternalServerError, "failed creating render job")
		return
	}

	renderCtx, cancel := context.WithCancel(context.Background())
	h.activeRendersMutex.Lock()
	h.activeRenderCancels[jobID] = cancel
	h.activeRendersMutex.Unlock()

	// Launch async composition render
	go func() {
		defer func() {
			h.activeRendersMutex.Lock()
			delete(h.activeRenderCancels, jobID)
			h.activeRendersMutex.Unlock()
			cancel()
		}()

		var composerOptions []composer.Option
		if request.HardwareAcceleration != "" {
			composerOptions = append(composerOptions, composer.WithHardwareAcceleration(request.HardwareAcceleration, false))
		}
		if request.Preset != "" {
			switch strings.ToLower(request.Preset) {
			case "tiktok", "reels", "shorts":
				composerOptions = append(composerOptions, composer.WithPreset(presets.TikTokVertical1080p60()))
			case "youtube_4k", "4k":
				composerOptions = append(composerOptions, composer.WithPreset(presets.YouTube4K60()))
			case "youtube_1080p", "1080p":
				composerOptions = append(composerOptions, composer.WithPreset(presets.YouTube1080p60()))
			}
		}

		composerInstance := composer.New(composerOptions...)

		progressCallback := func(event executor.ProgressEvent) {
			_ = h.database.UpdateRenderJobProgress(
				context.Background(),
				jobID,
				event.Percentage,
				event.Frame,
				event.FPS,
				event.Time.Seconds(),
				event.Speed,
			)

			h.websocketHub.Broadcast("render_progress", map[string]any{
				"job_id":        jobID,
				"percentage":    event.Percentage,
				"current_frame": event.Frame,
				"current_fps":   event.FPS,
				"current_time":  event.Time.Seconds(),
				"render_speed":  event.Speed,
				"bitrate":       event.Bitrate,
				"total_size":    event.TotalSize,
			})
		}

		renderResult, renderErr := composerInstance.Render(renderCtx, compositionTimeline, outputFilePath, progressCallback)
		if renderErr != nil {
			h.logger.Error("render job failed", "job_id", jobID, "error", renderErr)
			_ = h.database.CompleteRenderJob(context.Background(), jobID, "failed", renderErr.Error())
			h.websocketHub.Broadcast("render_failed", map[string]any{
				"job_id": jobID,
				"error":  renderErr.Error(),
			})
			return
		}

		h.logger.Info("render job completed successfully", "job_id", jobID, "output_path", renderResult.OutputPath)
		_ = h.database.CompleteRenderJob(context.Background(), jobID, "completed", "")
		h.websocketHub.Broadcast("render_complete", map[string]any{
			"job_id":      jobID,
			"output_path": outputFileName,
			"download_url": fmt.Sprintf("/api/exports/%s", outputFileName),
		})
	}()

	respondJSON(responseWriter, http.StatusAccepted, StartRenderResponse{
		JobID:      jobID,
		Status:     "rendering",
		OutputPath: outputFileName,
	})
}

// HandleCancelRender cancels an ongoing render job.
func (h *Handlers) HandleCancelRender(responseWriter http.ResponseWriter, httpRequest *http.Request, jobID string) {
	h.activeRendersMutex.Lock()
	cancelFunc, exists := h.activeRenderCancels[jobID]
	h.activeRendersMutex.Unlock()

	if !exists {
		respondError(responseWriter, http.StatusNotFound, "render job not found or already finished")
		return
	}

	cancelFunc()
	_ = h.database.CompleteRenderJob(httpRequest.Context(), jobID, "cancelled", "Cancelled by user")
	h.websocketHub.Broadcast("render_cancelled", map[string]any{
		"job_id": jobID,
	})

	respondJSON(responseWriter, http.StatusOK, StandardResponse{
		Success: true,
		Message: "render job cancelled",
	})
}

// HandleGetRenderJob returns status of a render job.
func (h *Handlers) HandleGetRenderJob(responseWriter http.ResponseWriter, httpRequest *http.Request, jobID string) {
	job, err := h.database.GetRenderJob(httpRequest.Context(), jobID)
	if err != nil {
		h.logger.Error("failed retrieving render job", "error", err, "job_id", jobID)
		respondError(responseWriter, http.StatusInternalServerError, "failed retrieving render job")
		return
	}
	if job == nil {
		respondError(responseWriter, http.StatusNotFound, "render job not found")
		return
	}

	respondJSON(responseWriter, http.StatusOK, StandardResponse{
		Success: true,
		Data:    job,
	})
}
