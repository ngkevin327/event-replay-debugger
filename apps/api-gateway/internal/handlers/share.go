package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	gwmw "github.com/replay/platform/apps/api-gateway/internal/middleware"
	"github.com/replay/platform/apps/api-gateway/internal/store"
)

type shareRecord struct {
	IncidentID string
	OrgID      string
	Scope      string
	ExpiresAt  time.Time
}

var (
	shareMu    sync.RWMutex
	shareByTok = map[string]shareRecord{}
)

// ShareHandler manages expiring read-only share links.
type ShareHandler struct {
	Store *store.Store
}

type createShareRequest struct {
	TTLHours int `json:"ttl_hours"`
}

type createShareResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	URL       string    `json:"url"`
}

// CreateShareToken handles POST /v1/incidents/{incidentId}/share-tokens.
func (h *ShareHandler) CreateShareToken(w http.ResponseWriter, r *http.Request) {
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
	var req createShareRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	if req.TTLHours <= 0 {
		req.TTLHours = 72
	}
	raw := make([]byte, 32)
	_, _ = rand.Read(raw)
	token := hex.EncodeToString(raw)
	expires := time.Now().UTC().Add(time.Duration(req.TTLHours) * time.Hour)
	shareMu.Lock()
	shareByTok[token] = shareRecord{
		IncidentID: incidentID,
		OrgID:      inc.OrgID,
		Scope:      gwmw.ShareScopeRead,
		ExpiresAt:  expires,
	}
	shareMu.Unlock()
	writeJSON(w, http.StatusCreated, createShareResponse{
		Token:     token,
		ExpiresAt: expires,
		URL:       "/v1/shared/incidents/" + incidentID + "?token=" + token,
	})
}

// GetSharedIncident handles GET /v1/shared/incidents/{incidentId}.
func (h *ShareHandler) GetSharedIncident(w http.ResponseWriter, r *http.Request) {
	incidentID := chi.URLParam(r, "incidentId")
	shareInc, ok := gwmw.ShareIncidentIDFromContext(r.Context())
	if !ok || shareInc != incidentID {
		writeError(w, http.StatusForbidden, "forbidden", "token not valid for incident")
		return
	}
	orgID, ok := gwmw.ShareOrgIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", "share context required")
		return
	}
	inc, err := h.Store.GetIncident(r.Context(), orgID, incidentID)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "incident not found")
		return
	}
	writeJSON(w, http.StatusOK, inc)
}

// LookupShareToken resolves a token for middleware.
func LookupShareToken(token string) (incidentID, orgID, scope string, expired bool, ok bool) {
	shareMu.RLock()
	rec, found := shareByTok[token]
	shareMu.RUnlock()
	if !found {
		return "", "", "", false, false
	}
	return rec.IncidentID, rec.OrgID, rec.Scope, time.Now().UTC().After(rec.ExpiresAt), true
}

// HashShareToken hashes tokens for persistence.
func HashShareToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
