package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	gwmw "github.com/replay/platform/apps/api-gateway/internal/middleware"
	"github.com/replay/platform/apps/api-gateway/internal/store"
)

// SnapshotsHandler ingests HTTP offset snapshots.
type SnapshotsHandler struct {
	Store *store.Store
}

type ingestSnapshotRequest struct {
	ConsumerGroup string `json:"consumer_group"`
	Topic         string `json:"topic"`
	Partition     int    `json:"partition"`
	Offset        int64  `json:"offset"`
}

// IngestSnapshot handles POST /v1/incidents/{incidentId}/snapshots.
func (h *SnapshotsHandler) IngestSnapshot(w http.ResponseWriter, r *http.Request) {
	incidentID := chi.URLParam(r, "incidentId")
	orgID, ok := gwmw.OrgIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", "session required")
		return
	}
	if _, err := h.Store.GetIncident(r.Context(), orgID, incidentID); err != nil {
		writeError(w, http.StatusNotFound, "not_found", "incident not found")
		return
	}
	var req ingestSnapshotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid body")
		return
	}
	if err := h.Store.InsertOffsetSnapshot(r.Context(), incidentID, req.ConsumerGroup, req.Topic, req.Partition, req.Offset); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not store snapshot")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"status": "ok"})
}
