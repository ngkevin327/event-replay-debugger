package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/replay/platform/apps/api-gateway/internal/config"
	gwmw "github.com/replay/platform/apps/api-gateway/internal/middleware"
)

// Server wraps the HTTP listener and router.
type Server struct {
	cfg    config.Config
	router chi.Router
	http   *http.Server
}

// New builds a Server with routes registered.
func New(cfg config.Config, deps RouteDeps) *Server {
	r := chi.NewRouter()
	r.Use(gwmw.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(gwmw.CORS)

	s := &Server{cfg: cfg, router: r}
	s.registerRoutes(deps)

	s.http = &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	return s
}

func (s *Server) registerRoutes(deps RouteDeps) {
	s.router.Get("/health", s.handleHealth)
	s.router.Get("/ready", s.handleReady)
	if deps.Store != nil {
		RegisterV1Routes(s.router, deps)
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleReady(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// Router exposes the chi router for tests.
func (s *Server) Router() chi.Router {
	return s.router
}

// Run starts listening until error or shutdown.
func (s *Server) Run() error {
	return s.http.ListenAndServe()
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
