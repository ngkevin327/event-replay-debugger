package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/replay/platform/apps/ingestion/internal/metrics"
)

// Server hosts ingestion HTTP endpoints.
type Server struct {
	router chi.Router
	http   *http.Server
}

// New constructs the ingestion server.
func New(addr string) *Server {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	s := &Server{router: r}
	s.RegisterRoutes()
	s.http = &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return s
}

// RegisterRoutes mounts ingestion API routes.
func (s *Server) RegisterRoutes() {
	s.router.Get("/health", s.handleHealth)
	s.router.Handle("/metrics", metrics.Handler())
	s.router.Post("/v1/ingest/batch", s.handleIngestBatch)
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleIngestBatch(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer metrics.ObserveBatch(start)
	if r.Method != http.MethodPost {
		metrics.IncBatchError()
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method_not_allowed"})
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]string{
		"status": "accepted",
		"note":   "batch ingest stub",
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// Router exposes the router for tests.
func (s *Server) Router() chi.Router {
	return s.router
}

// Run listens until shutdown.
func (s *Server) Run() error {
	return s.http.ListenAndServe()
}

// Shutdown stops the server gracefully.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
