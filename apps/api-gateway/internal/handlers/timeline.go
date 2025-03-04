package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/replay/platform/apps/api-gateway/internal/cache"
	gwmw "github.com/replay/platform/apps/api-gateway/internal/middleware"
)

// TimelineHandler serves incident timeline artifacts.
type TimelineHandler struct {
	Loader TimelineLoader
	Cache  *cache.TimelineCache
}

// CachedLoader wraps a TimelineLoader with TimelineCache.
type CachedLoader struct {
	Inner TimelineLoader
	Cache *cache.TimelineCache
}

// LoadTimeline returns cached timeline when present.
func (c *CachedLoader) LoadTimeline(ctx context.Context, incidentID string) ([]byte, int, error) {
	if c.Cache != nil {
		if body, ver, ok := c.Cache.Get(incidentID); ok {
			return body, ver, nil
		}
	}
	body, ver, err := c.Inner.LoadTimeline(ctx, incidentID)
	if err == nil && c.Cache != nil {
		c.Cache.Set(incidentID, body, ver)
	}
	return body, ver, err
}

// TimelineLoader fetches timeline JSON for an incident.
type TimelineLoader interface {
	LoadTimeline(ctx context.Context, incidentID string) ([]byte, int, error)
}

// GetTimeline handles GET /v1/incidents/{incidentId}/timeline.
func (h *TimelineHandler) GetTimeline(w http.ResponseWriter, r *http.Request) {
	incidentID := chi.URLParam(r, "incidentId")
	_, ok := gwmw.OrgIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", "session required")
		return
	}
	body, version, err := h.Loader.LoadTimeline(r.Context(), incidentID)
	if err != nil {
		writeError(w, http.StatusConflict, "not_ready", "timeline not ready")
		return
	}
	var doc any
	_ = json.Unmarshal(body, &doc)
	writeJSON(w, http.StatusOK, map[string]any{
		"version": version,
		"timeline": doc,
	})
}
