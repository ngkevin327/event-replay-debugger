package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/replay/platform/apps/ingestion/internal/server"
)

func TestServerHealth(t *testing.T) {
	srv := server.New(":0", server.Deps{})
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	srv.Router().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status: %d", rec.Code)
	}
}
