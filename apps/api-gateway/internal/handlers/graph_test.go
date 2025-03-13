package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

type stubGraph struct{}

func (stubGraph) LoadGraph(ctx context.Context, incidentID string) ([]byte, error) {
	return []byte(`{"nodes":[],"edges":[]}`), nil
}

func TestGetGraph(t *testing.T) {
	h := &GraphHandler{Loader: stubGraph{}}
	req := httptest.NewRequest(http.MethodGet, "/v1/incidents/inc-1/graph", nil)
	rec := httptest.NewRecorder()
	h.GetGraph(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
}
