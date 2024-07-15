package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	gwmw "github.com/replay/platform/apps/api-gateway/internal/middleware"
	"github.com/replay/platform/apps/api-gateway/internal/store"
)

// OrgsHandler serves organization CRUD.
type OrgsHandler struct {
	Store *store.Store
}

type orgRequest struct {
	Name     string `json:"name"`
	PlanTier string `json:"plan_tier"`
}

// CreateOrg handles POST /v1/orgs.
func (h *OrgsHandler) CreateOrg(w http.ResponseWriter, r *http.Request) {
	var req orgRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "validation_error", "name required")
		return
	}
	plan := req.PlanTier
	if plan == "" {
		plan = "starter"
	}
	org, err := h.Store.CreateOrganization(r.Context(), req.Name, plan)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not create org")
		return
	}
	userID, _ := gwmw.UserIDFromContext(r.Context())
	_ = h.Store.CreateMembership(r.Context(), org.ID, userID, store.RoleAdmin)
	writeJSON(w, http.StatusCreated, org)
}

// GetOrg handles GET /v1/orgs/{id}.
func (h *OrgsHandler) GetOrg(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	orgID, ok := gwmw.OrgIDFromContext(r.Context())
	if !ok || id != orgID {
		writeError(w, http.StatusForbidden, "forbidden", "org access denied")
		return
	}
	org, err := h.Store.GetOrganization(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "organization not found")
		return
	}
	writeJSON(w, http.StatusOK, org)
}

// UpdateOrg handles PUT /v1/orgs/{id}.
func (h *OrgsHandler) UpdateOrg(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	orgID, ok := gwmw.OrgIDFromContext(r.Context())
	if !ok || id != orgID {
		writeError(w, http.StatusForbidden, "forbidden", "org access denied")
		return
	}
	var req orgRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	org, err := h.Store.UpdateOrganization(r.Context(), id, req.Name, req.PlanTier)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not update org")
		return
	}
	writeJSON(w, http.StatusOK, org)
}
