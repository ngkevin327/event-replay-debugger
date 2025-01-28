package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	gwmw "github.com/replay/platform/apps/api-gateway/internal/middleware"
	"github.com/replay/platform/apps/api-gateway/internal/quota"
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

// ListIncidents handles GET /v1/projects/{projectId}/incidents.
func (h *IncidentsHandler) ListIncidents(w http.ResponseWriter, r *http.Request) {
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
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	incidents, err := h.Store.ListIncidents(r.Context(), orgID, projectID, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not list incidents")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"incidents": incidents})
}

// GetIncident handles GET /v1/incidents/{incidentId}.
func (h *IncidentsHandler) GetIncident(w http.ResponseWriter, r *http.Request) {
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
	org, err := h.Store.GetOrganization(r.Context(), orgID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not load org")
		return
	}
	pct := RetentionCoveragePercent(org.PlanTier, inc.WindowStart, inc.WindowEnd, time.Now())
	inc.CoveragePercent = &pct
	writeJSON(w, http.StatusOK, inc)
}

// RetentionCoveragePercent estimates how much of the requested window is within plan retention.
func RetentionCoveragePercent(planTier string, windowStart, windowEnd, now time.Time) float64 {
	clampedStart, clampedEnd := quota.ClampQueryWindow(planTier, windowStart, windowEnd, now)
	requested := windowEnd.Sub(windowStart).Seconds()
	if requested <= 0 {
		return 0
	}
	available := clampedEnd.Sub(clampedStart).Seconds()
	if available <= 0 {
		return 0
	}
	pct := (available / requested) * 100
	if pct > 100 {
		return 100
	}
	return pct
}
