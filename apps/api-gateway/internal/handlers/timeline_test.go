package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

type stubTimeline struct{}

func (stubTimeline) LoadTimeline(ctx context.Context, incidentID string) ([]byte, int, error) {
	return []byte(`{"events":[]}`), 1, nil
}

func TestGetTimeline(t *testing.T) {
	h := &TimelineHandler{Loader: stubTimeline{}}
	req := httptest.NewRequest(http.MethodGet, "/v1/incidents/inc-1/timeline", nil)
	rec := httptest.NewRecorder()
	h.GetTimeline(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
}
