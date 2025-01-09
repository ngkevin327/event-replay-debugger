package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	gwmw "github.com/replay/platform/apps/api-gateway/internal/middleware"
	"github.com/replay/platform/apps/api-gateway/internal/store"
)

// IncidentsHandler serves incident CRUD.
type IncidentsHandler struct {
	Store *store.Store
}

type createIncidentRequest struct {
	WindowStart  time.Time `json:"window_start"`
	WindowEnd    time.Time `json:"window_end"`
	TopicFilters []string  `json:"topic_filters"`
}

// CreateIncident handles POST /v1/projects/{projectId}/incidents.
func (h *IncidentsHandler) CreateIncident(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectId")
	orgID, ok := gwmw.OrgIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", "session required")
		return
	}
	project, err := h.Store.GetProject(r.Context(), projectID)
	if err != nil || project.OrgID != orgID {
		writeError(w, http.StatusNotFound, "not_found", "project not found")
		return
	}
	var req createIncidentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid body")
		return
	}
	if req.WindowEnd.Before(req.WindowStart) {
		writeError(w, http.StatusBadRequest, "validation_error", "window_end must be after window_start")
		return
	}
	inc, err := h.Store.InsertIncident(r.Context(), projectID, orgID, store.WindowBounds{
		WindowStart: req.WindowStart, WindowEnd: req.WindowEnd, TopicFilters: req.TopicFilters,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not create incident")
		return
	}
	writeJSON(w, http.StatusCreated, inc)
}
