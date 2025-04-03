package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	gwmw "github.com/replay/platform/apps/api-gateway/internal/middleware"
	"github.com/replay/platform/apps/api-gateway/internal/store"
)

// ReplaysHandler serves replay control APIs.
type ReplaysHandler struct {
	Store *store.Store
}

type createReplayRequest struct {
	TimingMode string `json:"timing_mode"`
}

// CreateReplay handles POST /v1/incidents/{incidentId}/replays.
func (h *ReplaysHandler) CreateReplay(w http.ResponseWriter, r *http.Request) {
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
	if inc.Status != store.IncidentReady {
		writeError(w, http.StatusConflict, "not_ready", "incident timeline not ready")
		return
	}
	var req createReplayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid body")
		return
	}
	mode := req.TimingMode
	if mode == "" {
		mode = "strict"
	}
	if mode != "strict" && mode != "compressed" {
		writeError(w, http.StatusBadRequest, "validation_error", "timing_mode must be strict or compressed")
		return
	}
	run, err := h.Store.CreateReplayRun(r.Context(), incidentID, orgID, mode)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not create replay")
		return
	}
	writeJSON(w, http.StatusCreated, run)
}
