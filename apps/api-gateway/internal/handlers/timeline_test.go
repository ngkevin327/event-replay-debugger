package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	gwmw "github.com/replay/platform/apps/api-gateway/internal/middleware"
)

type stubTimeline struct{}

func (stubTimeline) LoadTimeline(ctx context.Context, incidentID string) ([]byte, int, error) {
	_ = ctx
	_ = incidentID
	return []byte(`{"events":[]}`), 1, nil
}

func TestGetTimeline(t *testing.T) {
	h := &TimelineHandler{Loader: stubTimeline{}}
	req := httptest.NewRequest(http.MethodGet, "/v1/incidents/inc-1/timeline", nil)
	req = req.WithContext(gwmw.TestContextWithOrg(req.Context(), "org-1"))
	rec := httptest.NewRecorder()
	h.GetTimeline(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
}
