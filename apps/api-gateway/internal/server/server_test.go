package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/replay/platform/apps/api-gateway/internal/config"
	"github.com/replay/platform/apps/api-gateway/internal/server"
)

func TestHealthOK(t *testing.T) {
	srv := server.New(config.Config{HTTPAddr: ":0"}, server.RouteDeps{})
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	srv.Router().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status: %d", rec.Code)
	}
}

func TestReadyOK(t *testing.T) {
	srv := server.New(config.Config{HTTPAddr: ":0"}, server.RouteDeps{})
	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	rec := httptest.NewRecorder()
	srv.Router().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status: %d", rec.Code)
	}
}
