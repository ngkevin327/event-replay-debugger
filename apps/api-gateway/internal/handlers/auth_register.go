package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/replay/platform/apps/api-gateway/internal/auth"
	"github.com/replay/platform/apps/api-gateway/internal/store"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]string{"error": code, "message": message})
}

// RegisterHandler handles user and organization bootstrap registration.
type RegisterHandler struct {
	Store     *store.Store
	JWTSecret string
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	OrgName  string `json:"org_name"`
}

type authUserPayload struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	OrgID string `json:"org_id"`
}

type registerResponse struct {
	AccessToken string          `json:"access_token"`
	User        authUserPayload `json:"user"`
}

// ServeHTTP implements POST /v1/auth/register.
func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "POST required")
		return
	}

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.OrgName == "" {
		writeError(w, http.StatusBadRequest, "validation_error", "email and org_name required")
		return
	}

	if _, err := h.Store.GetUserByEmail(r.Context(), req.Email); err == nil {
		writeError(w, http.StatusConflict, "email_exists", "email already registered")
		return
	} else if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not look up user")
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_password", err.Error())
		return
	}

	user, err := h.Store.CreateUser(r.Context(), req.Email, hash)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not create user")
		return
	}
	org, err := h.Store.CreateOrganization(r.Context(), req.OrgName, "starter")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not create organization")
		return
	}
	if err := h.Store.CreateMembership(r.Context(), org.ID, user.ID, store.RoleAdmin); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not create membership")
		return
	}

	tokens, err := auth.IssueTokens(h.JWTSecret, user.ID, org.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not issue tokens")
		return
	}

	writeJSON(w, http.StatusCreated, registerResponse{
		AccessToken: tokens.Access,
		User: authUserPayload{
			ID: user.ID, Email: user.Email, OrgID: org.ID,
		},
	})
}
