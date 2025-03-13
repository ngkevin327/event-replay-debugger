package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	gwmw "github.com/replay/platform/apps/api-gateway/internal/middleware"
)

// GraphHandler serves workflow graph artifacts.
type GraphHandler struct {
	Loader GraphLoader
}

// GraphLoader loads graph JSON for an incident.
type GraphLoader interface {
	LoadGraph(ctx context.Context, incidentID string) ([]byte, error)
}

// GetGraph handles GET /v1/incidents/{incidentId}/graph.
func (h *GraphHandler) GetGraph(w http.ResponseWriter, r *http.Request) {
	incidentID := chi.URLParam(r, "incidentId")
	_, ok := gwmw.OrgIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", "session required")
		return
	}
	body, err := h.Loader.LoadGraph(r.Context(), incidentID)
	if err != nil {
		writeError(w, http.StatusConflict, "not_ready", "graph not ready")
		return
	}
	var doc any
	_ = json.Unmarshal(body, &doc)
	writeJSON(w, http.StatusOK, doc)
}
