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

// GetReplay handles GET /v1/replays/{replayId}.
func (h *ReplaysHandler) GetReplay(w http.ResponseWriter, r *http.Request) {
	replayID := chi.URLParam(r, "replayId")
	orgID, ok := gwmw.OrgIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", "session required")
		return
	}
	run, err := h.Store.GetReplayRun(r.Context(), orgID, replayID)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "replay not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"replay":              run,
		"divergence_summary": DivergenceSummary(run),
	})
}

// DivergenceSummary builds a compact divergence payload for API clients.
func DivergenceSummary(run store.ReplayRun) map[string]any {
	out := map[string]any{"status": run.Status}
	if run.DivergenceIndex != nil {
		out["first_mismatch_index"] = *run.DivergenceIndex
	}
	if run.ReportURI != nil {
		out["report_uri"] = *run.ReportURI
	}
	return out
}

// CancelReplay handles DELETE /v1/replays/{replayId}.
func (h *ReplaysHandler) CancelReplay(w http.ResponseWriter, r *http.Request) {
	replayID := chi.URLParam(r, "replayId")
	orgID, ok := gwmw.OrgIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", "session required")
		return
	}
	if _, err := h.Store.GetReplayRun(r.Context(), orgID, replayID); err != nil {
		writeError(w, http.StatusNotFound, "not_found", "replay not found")
		return
	}
	if err := h.Store.CancelReplayRun(r.Context(), replayID); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not cancel replay")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "cancelled"})
}
