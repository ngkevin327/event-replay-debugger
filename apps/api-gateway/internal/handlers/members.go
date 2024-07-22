package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/replay/platform/apps/api-gateway/internal/audit"
	gwmw "github.com/replay/platform/apps/api-gateway/internal/middleware"
	"github.com/replay/platform/apps/api-gateway/internal/store"
)

// MembersHandler manages organization membership.
type MembersHandler struct {
	Store *store.Store
	Audit *audit.Logger
}

type inviteRequest struct {
	Email string             `json:"email"`
	Role  store.MembershipRole `json:"role"`
}

type roleRequest struct {
	Role store.MembershipRole `json:"role"`
}

// InviteMember handles POST /v1/orgs/{id}/members.
func (h *MembersHandler) InviteMember(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "id")
	ctxOrg, ok := gwmw.OrgIDFromContext(r.Context())
	if !ok || orgID != ctxOrg {
		writeError(w, http.StatusForbidden, "forbidden", "org access denied")
		return
	}
	var req inviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	user, err := h.Store.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "user not found")
		return
	}
	if err := h.Store.CreateMembership(r.Context(), orgID, user.ID, req.Role); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not invite member")
		return
	}
	actor, _ := gwmw.UserIDFromContext(r.Context())
	if h.Audit != nil {
		h.Audit.LogAction(actor, "member.invite", "organization", orgID, r.RemoteAddr)
	}
	writeJSON(w, http.StatusCreated, map[string]string{"user_id": user.ID})
}

// UpdateRole handles PUT /v1/orgs/{id}/members/{userId}.
func (h *MembersHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "id")
	userID := chi.URLParam(r, "userId")
	ctxOrg, ok := gwmw.OrgIDFromContext(r.Context())
	if !ok || orgID != ctxOrg {
		writeError(w, http.StatusForbidden, "forbidden", "org access denied")
		return
	}
	var req roleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	if err := h.Store.UpdateMembershipRole(r.Context(), orgID, userID, req.Role); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not update role")
		return
	}
	actor, _ := gwmw.UserIDFromContext(r.Context())
	if h.Audit != nil {
		h.Audit.LogAction(actor, "member.role_change", "user", userID, r.RemoteAddr)
	}
	writeJSON(w, http.StatusOK, map[string]string{"user_id": userID, "role": string(req.Role)})
}

// RemoveMember handles DELETE /v1/orgs/{id}/members/{userId}.
func (h *MembersHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "id")
	userID := chi.URLParam(r, "userId")
	ctxOrg, ok := gwmw.OrgIDFromContext(r.Context())
	if !ok || orgID != ctxOrg {
		writeError(w, http.StatusForbidden, "forbidden", "org access denied")
		return
	}
	if err := h.Store.DeleteMembership(r.Context(), orgID, userID); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not remove member")
		return
	}
	actor, _ := gwmw.UserIDFromContext(r.Context())
	if h.Audit != nil {
		h.Audit.LogAction(actor, "member.remove", "user", userID, r.RemoteAddr)
	}
	w.WriteHeader(http.StatusNoContent)
}
