package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	gwmw "github.com/replay/platform/apps/api-gateway/internal/middleware"
	"github.com/replay/platform/apps/api-gateway/internal/store"
)

// ProjectsHandler serves project CRUD.
type ProjectsHandler struct {
	Store *store.Store
}

type projectRequest struct {
	Name string `json:"name"`
}

// ListProjects handles GET /v1/projects.
func (h *ProjectsHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	orgID, ok := gwmw.OrgIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", "session required")
		return
	}
	projects, err := h.Store.ListProjectsByOrg(r.Context(), orgID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not list projects")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"projects": projects})
}

// CreateProject handles POST /v1/projects.
func (h *ProjectsHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	orgID, ok := gwmw.OrgIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", "session required")
		return
	}
	var req projectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "validation_error", "name required")
		return
	}
	project, err := h.Store.CreateProject(r.Context(), orgID, req.Name)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not create project")
		return
	}
	writeJSON(w, http.StatusCreated, project)
}

// GetProject handles GET /v1/projects/{id}.
func (h *ProjectsHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	orgID, ok := gwmw.OrgIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", "session required")
		return
	}
	project, err := h.Store.GetProject(r.Context(), id)
	if err != nil || project.OrgID != orgID {
		writeError(w, http.StatusNotFound, "not_found", "project not found")
		return
	}
	writeJSON(w, http.StatusOK, project)
}
