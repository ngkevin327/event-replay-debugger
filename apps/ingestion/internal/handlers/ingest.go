package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/replay/platform/apps/ingestion/internal/auth"
	"github.com/replay/platform/apps/ingestion/internal/batch"
	"github.com/replay/platform/apps/ingestion/internal/index"
	"github.com/replay/platform/apps/ingestion/internal/metrics"
	"github.com/replay/platform/packages/shared-go/event"
)

// IngestHandler accepts agent batch uploads.
type IngestHandler struct {
	Validator *auth.Validator
	Uploader  *batch.Uploader
	Writer    *index.Writer
	OrgID     string
}

type ingestRequest struct {
	Events []event.CapturedEvent `json:"events"`
}

type ingestResponse struct {
	Accepted          int `json:"accepted"`
	DuplicatesIgnored int `json:"duplicates_ignored"`
	Rejected          int `json:"rejected"`
}

// ServeHTTP handles POST /v1/ingest/batch.
func (h *IngestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer metrics.ObserveBatch(start)

	key := r.Header.Get("X-Replay-Key")
	projectID, err := h.Validator.ValidateAPIKey(r.Context(), key)
	if err != nil {
		metrics.IncBatchError()
		writeError(w, http.StatusUnauthorized, "unauthorized", "invalid api key")
		return
	}

	var req ingestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		metrics.IncBatchError()
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid body")
		return
	}

	filtered := make([]event.CapturedEvent, 0, len(req.Events))
	rejected := 0
	for _, ev := range req.Events {
		if ev.ProjectID != projectID {
			rejected++
			continue
		}
		filtered = append(filtered, ev)
	}

	if h.Uploader != nil && len(filtered) > 0 {
		h.Uploader.OrgID = h.OrgID
		uploaded, err := h.Uploader.UploadBatch(r.Context(), filtered)
		if err != nil {
			metrics.IncBatchError()
			writeError(w, http.StatusInternalServerError, "upload_failed", "payload upload failed")
			return
		}
		filtered = uploaded
	}

	if h.Writer != nil && len(filtered) > 0 {
		if err := h.Writer.WriteBatch(r.Context(), h.OrgID, filtered); err != nil {
			metrics.IncBatchError()
			writeError(w, http.StatusInternalServerError, "index_failed", "clickhouse insert failed")
			return
		}
	}

	writeJSON(w, http.StatusAccepted, ingestResponse{
		Accepted:          len(filtered),
		DuplicatesIgnored: 0,
		Rejected:          rejected,
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, msg string) {
	writeJSON(w, status, map[string]string{"error": code, "message": msg})
}
