// Package server provides the HTTP & WebSocket server for the Vidonex Web Studio.
package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/farshidrezaei/vidonex/server/api"
	"github.com/farshidrezaei/vidonex/server/db"
	"github.com/farshidrezaei/vidonex/server/ws"
)

// Config configures the server runtime.
type Config struct {
	Port            int
	Host            string
	DataDirectory   string
	StaticDirectory string
	Logger          *slog.Logger
}

// Server encapsulates the HTTP server, database, WebSocket hub, and routing.
type Server struct {
	config       Config
	httpServer   *http.Server
	database     *db.Database
	websocketHub *ws.Hub
	apiHandlers  *api.Handlers
	logger       *slog.Logger
}

// New creates and initializes a new Server.
func New(config Config) (*Server, error) {
	if config.Port <= 0 {
		config.Port = 8080
	}
	if config.DataDirectory == "" {
		homeDir, _ := os.UserHomeDir()
		config.DataDirectory = filepath.Join(homeDir, ".vidonex")
	}
	if config.Logger == nil {
		config.Logger = slog.Default()
	}

	databasePath := filepath.Join(config.DataDirectory, "vidonex.db")
	mediaDirectory := filepath.Join(config.DataDirectory, "media")
	exportDirectory := filepath.Join(config.DataDirectory, "exports")

	databaseInstance, err := db.Open(databasePath)
	if err != nil {
		return nil, fmt.Errorf("failed initializing database: %w", err)
	}

	websocketHub := ws.NewHub(config.Logger)
	apiHandlers := api.NewHandlers(databaseInstance, websocketHub, mediaDirectory, exportDirectory, config.Logger)

	serverInstance := &Server{
		config:       config,
		database:     databaseInstance,
		websocketHub: websocketHub,
		apiHandlers:  apiHandlers,
		logger:       config.Logger,
	}

	return serverInstance, nil
}

// Handler constructs and returns the HTTP handler (mux) including all REST, WebSocket, and file endpoints.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	// WebSocket endpoint
	mux.HandleFunc("/ws", s.websocketHub.HandleWebSocket)

	// API Endpoints
	mux.HandleFunc("/api/projects", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.apiHandlers.HandleListProjects(w, r)
		case http.MethodPost:
			s.apiHandlers.HandleCreateProject(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/projects/", func(w http.ResponseWriter, r *http.Request) {
		trimmedPath := strings.TrimPrefix(r.URL.Path, "/api/projects/")
		segments := strings.Split(trimmedPath, "/")
		projectID := segments[0]

		if len(segments) == 1 {
			switch r.Method {
			case http.MethodGet:
				s.apiHandlers.HandleGetProject(w, r, projectID)
			case http.MethodPut, http.MethodPost:
				s.apiHandlers.HandleUpdateProject(w, r, projectID)
			case http.MethodDelete:
				s.apiHandlers.HandleDeleteProject(w, r, projectID)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		if len(segments) == 2 && segments[1] == "assets" {
			if r.Method == http.MethodGet {
				s.apiHandlers.HandleListAssets(w, r, projectID)
				return
			}
		}

		http.NotFound(w, r)
	})

	mux.HandleFunc("/api/media/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			s.apiHandlers.HandleUploadMedia(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/api/media/", func(w http.ResponseWriter, r *http.Request) {
		trimmedPath := strings.TrimPrefix(r.URL.Path, "/api/media/")
		if r.Method == http.MethodDelete && trimmedPath != "" {
			s.apiHandlers.HandleDeleteAsset(w, r, trimmedPath)
			return
		}
		http.NotFound(w, r)
	})

	mux.HandleFunc("/api/spec/validate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			s.apiHandlers.HandleValidateSpec(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/api/spec/graph", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			s.apiHandlers.HandleGenerateGraph(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/api/render/start", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			s.apiHandlers.HandleStartRender(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/api/render/", func(w http.ResponseWriter, r *http.Request) {
		trimmedPath := strings.TrimPrefix(r.URL.Path, "/api/render/")
		segments := strings.Split(trimmedPath, "/")
		jobID := segments[0]

		if len(segments) == 1 && r.Method == http.MethodGet {
			s.apiHandlers.HandleGetRenderJob(w, r, jobID)
			return
		}
		if len(segments) == 2 && segments[1] == "cancel" && r.Method == http.MethodPost {
			s.apiHandlers.HandleCancelRender(w, r, jobID)
			return
		}
		http.NotFound(w, r)
	})

	// Static Media & Exports Serving (with Range request support)
	mediaDirectory := filepath.Join(s.config.DataDirectory, "media")
	exportDirectory := filepath.Join(s.config.DataDirectory, "exports")

	mux.Handle("/api/media/files/", http.StripPrefix("/api/media/files/", http.FileServer(http.Dir(mediaDirectory))))
	mux.Handle("/api/exports/", http.StripPrefix("/api/exports/", http.FileServer(http.Dir(exportDirectory))))

	// Static Web UI Hosting (SPA Fallback)
	if s.config.StaticDirectory != "" {
		fileServer := http.FileServer(http.Dir(s.config.StaticDirectory))
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/ws") {
				http.NotFound(w, r)
				return
			}
			indexPath := filepath.Join(s.config.StaticDirectory, "index.html")
			if r.URL.Path == "/" {
				http.ServeFile(w, r, indexPath)
				return
			}
			filePath := filepath.Join(s.config.StaticDirectory, filepath.Clean(r.URL.Path))
			fileInfo, err := os.Stat(filePath)
			if os.IsNotExist(err) || (err == nil && fileInfo.IsDir()) {
				http.ServeFile(w, r, indexPath)
				return
			}
			fileServer.ServeHTTP(w, r)
		})
	}

	return withCORS(mux)
}

// Start launches the WebSocket hub and starts listening on the configured HTTP port.
func (s *Server) Start() error {
	go s.websocketHub.Run()

	address := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	s.httpServer = &http.Server{
		Addr:              address,
		Handler:           s.Handler(),
		ReadHeaderTimeout: 15 * time.Second,
	}

	s.logger.Info("vidonex web studio server listening", "address", fmt.Sprintf("http://localhost:%d", s.config.Port))
	return s.httpServer.ListenAndServe()
}

// StartListener starts listening on an existing net.Listener (useful for random available port selection in desktop).
func (s *Server) StartListener(listener net.Listener) error {
	go s.websocketHub.Run()

	s.httpServer = &http.Server{
		Handler:           s.Handler(),
		ReadHeaderTimeout: 15 * time.Second,
	}

	return s.httpServer.Serve(listener)
}

// Database returns the underlying persistence database instance.
func (s *Server) Database() *db.Database {
	return s.database
}

// DataDirectory returns the root data storage directory path.
func (s *Server) DataDirectory() string {
	return s.config.DataDirectory
}

// MediaDirectory returns the media assets directory path.
func (s *Server) MediaDirectory() string {
	return filepath.Join(s.config.DataDirectory, "media")
}

// ExportDirectory returns the export videos directory path.
func (s *Server) ExportDirectory() string {
	return filepath.Join(s.config.DataDirectory, "exports")
}

// Config returns the configuration of the server.
func (s *Server) Config() Config {
	return s.config
}

// Stop gracefully shuts down the HTTP server and database connection.
func (s *Server) Stop(ctx context.Context) error {
	var errs []error
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			errs = append(errs, err)
		}
	}
	if s.database != nil {
		if err := s.database.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("errors during shutdown: %v", errs)
	}
	return nil
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
