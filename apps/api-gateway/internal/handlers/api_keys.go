package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/replay/platform/apps/api-gateway/internal/audit"
	"github.com/replay/platform/apps/api-gateway/internal/auth"
	gwmw "github.com/replay/platform/apps/api-gateway/internal/middleware"
	"github.com/replay/platform/apps/api-gateway/internal/store"
)

// APIKeysHandler manages project API keys.
type APIKeysHandler struct {
	Store *store.Store
	Audit *audit.Logger
}

type createKeyRequest struct {
	Name   string   `json:"name"`
	Scopes []string `json:"scopes"`
}

type createKeyResponse struct {
	ID        string   `json:"id"`
	ProjectID string   `json:"project_id"`
	Name      string   `json:"name"`
	Key       string   `json:"key"`
	Scopes    []string `json:"scopes"`
}

// CreateAPIKey handles POST /v1/projects/{id}/api-keys.
func (h *APIKeysHandler) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	if projectID == "" {
		writeError(w, http.StatusBadRequest, "validation_error", "project id required")
		return
	}
	var req createKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	if req.Name == "" || len(req.Scopes) == 0 {
		writeError(w, http.StatusBadRequest, "validation_error", "name and scopes required")
		return
	}
	for _, s := range req.Scopes {
		if !auth.ValidateScope([]string{"ingest", "read", "replay"}, s) {
			writeError(w, http.StatusBadRequest, "validation_error", "invalid scope")
			return
		}
	}

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

	plain, prefix, hash, err := auth.GenerateKey()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not generate key")
		return
	}
	row, err := h.Store.CreateAPIKey(r.Context(), projectID, req.Name, prefix, hash, req.Scopes)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not store key")
		return
	}
	if actor, ok := gwmw.UserIDFromContext(r.Context()); ok && h.Audit != nil {
		h.Audit.LogAction(actor, "api_key.create", "project", projectID, r.RemoteAddr)
	}
	writeJSON(w, http.StatusCreated, createKeyResponse{
		ID:        row.ID,
		ProjectID: row.ProjectID,
		Name:      row.Name,
		Key:       plain,
		Scopes:    row.Scopes,
	})
}

// RotateAPIKey handles POST /v1/projects/{id}/api-keys/{keyId}/rotate.
func (h *APIKeysHandler) RotateAPIKey(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	keyID := chi.URLParam(r, "keyId")
	if projectID == "" || keyID == "" {
		writeError(w, http.StatusBadRequest, "validation_error", "project and key id required")
		return
	}
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
	if err := h.Store.RevokeAPIKey(r.Context(), keyID); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not revoke key")
		return
	}
	plain, prefix, hash, err := auth.GenerateKey()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not generate key")
		return
	}
	row, err := h.Store.CreateAPIKey(r.Context(), projectID, "rotated", prefix, hash, []string{"ingest", "read", "replay"})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not store key")
		return
	}
	if actor, ok := gwmw.UserIDFromContext(r.Context()); ok && h.Audit != nil {
		h.Audit.LogAction(actor, "api_key.rotate", "api_key", keyID, r.RemoteAddr)
	}
	writeJSON(w, http.StatusCreated, createKeyResponse{
		ID:        row.ID,
		ProjectID: row.ProjectID,
		Name:      row.Name,
		Key:       plain,
		Scopes:    row.Scopes,
	})
}
