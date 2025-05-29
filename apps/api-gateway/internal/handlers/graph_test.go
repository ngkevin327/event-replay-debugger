package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	gwmw "github.com/replay/platform/apps/api-gateway/internal/middleware"
)

type stubGraph struct{}

func (stubGraph) LoadGraph(ctx context.Context, incidentID string) ([]byte, error) {
	return []byte(`{"nodes":[],"edges":[]}`), nil
}

func TestGetGraph(t *testing.T) {
	h := &GraphHandler{Loader: stubGraph{}}
	req := httptest.NewRequest(http.MethodGet, "/v1/incidents/inc-1/graph", nil)
	req = req.WithContext(gwmw.TestContextWithOrg(req.Context(), "org-1"))
	rec := httptest.NewRecorder()
	h.GetGraph(rec, req)
	if rec.Code != http.StatusConflict {
		t.Fatalf("status %d", rec.Code)
	}
}
