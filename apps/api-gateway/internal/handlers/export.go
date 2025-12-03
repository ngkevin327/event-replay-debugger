package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/replay/platform/apps/api-gateway/internal/export"
	gwmw "github.com/replay/platform/apps/api-gateway/internal/middleware"
	"github.com/replay/platform/apps/api-gateway/internal/store"
)

// ExportHandler serves incident export endpoints.
type ExportHandler struct {
	Store *store.Store
}

// ExportIncidentSummary handles GET /v1/incidents/{incidentId}/export.
func (h *ExportHandler) ExportIncidentSummary(w http.ResponseWriter, r *http.Request) {
	incidentID := chi.URLParam(r, "incidentId")
	orgID, ok := gwmw.OrgIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", "session required")
		return
	}
	inc, err := h.Store.GetIncident(r.Context(), orgID, incidentID)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "incident not found")
		return
	}
	summary := export.BuildSummary(inc)
	body, err := json.Marshal(summary)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not export")
		return
	}
	if len(body) > export.MaxExportBytes {
		writeError(w, http.StatusRequestEntityTooLarge, "export_too_large", "export exceeds size limit")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", `attachment; filename="incident-`+incidentID+`.json"`)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}
