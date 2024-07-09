package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/replay/platform/apps/api-gateway/internal/auth"
	"github.com/replay/platform/apps/api-gateway/internal/store"
)

// LoginHandler authenticates email/password and issues session cookies.
type LoginHandler struct {
	Store     *store.Store
	JWTSecret string
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	UserID      string `json:"user_id"`
	OrgID       string `json:"org_id"`
}

// ServeHTTP implements POST /v1/auth/login.
func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "POST required")
		return
	}
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	user, err := h.Store.GetUserByEmail(r.Context(), req.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusUnauthorized, "invalid_credentials", "invalid email or password")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not load user")
		return
	}
	if err := auth.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		writeError(w, http.StatusUnauthorized, "invalid_credentials", "invalid email or password")
		return
	}

	orgID, err := h.Store.GetPrimaryOrgForUser(r.Context(), user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not resolve organization")
		return
	}

	tokens, err := auth.IssueTokens(h.JWTSecret, user.ID, orgID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not issue tokens")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "replay_access",
		Value:    tokens.Access,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((15 * time.Minute).Seconds()),
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "replay_refresh",
		Value:    tokens.Refresh,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((7 * 24 * time.Hour).Seconds()),
	})

	writeJSON(w, http.StatusOK, loginResponse{
		AccessToken: tokens.Access,
		UserID:      user.ID,
		OrgID:       orgID,
	})
}
