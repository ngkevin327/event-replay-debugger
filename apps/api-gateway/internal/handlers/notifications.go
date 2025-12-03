package handlers

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	gwmw "github.com/replay/platform/apps/api-gateway/internal/middleware"
	"github.com/replay/platform/apps/api-gateway/internal/store"
)

// NotificationPrefs is project-level notification configuration.
type NotificationPrefs struct {
	ProjectID             string   `json:"project_id"`
	WebhookURL            string   `json:"webhook_url,omitempty"`
	WebhookSecret         string   `json:"webhook_secret,omitempty"`
	EmailEnabled          bool     `json:"email_enabled"`
	EmailRecipients       []string `json:"email_recipients"`
	NotifyIncidentReady   bool     `json:"notify_incident_ready"`
	NotifyReplayCompleted bool     `json:"notify_replay_completed"`
}

var (
	prefsMu      sync.RWMutex
	prefsByProj  = map[string]NotificationPrefs{}
	defaultPrefs = NotificationPrefs{EmailEnabled: true, NotifyIncidentReady: true, NotifyReplayCompleted: true}
)

// NotificationsHandler serves notification preference endpoints.
type NotificationsHandler struct {
	Store *store.Store
}

// GetPrefs handles GET /v1/projects/{id}/notification-preferences.
func (h *NotificationsHandler) GetPrefs(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	if !h.authorizeProject(w, r, projectID) {
		return
	}
	prefsMu.RLock()
	p, ok := prefsByProj[projectID]
	prefsMu.RUnlock()
	if !ok {
		p = defaultPrefs
		p.ProjectID = projectID
	}
	writeJSON(w, http.StatusOK, p)
}

// PutPrefs handles PUT /v1/projects/{id}/notification-preferences.
func (h *NotificationsHandler) PutPrefs(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	if !h.authorizeProject(w, r, projectID) {
		return
	}
	var req NotificationPrefs
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid body")
		return
	}
	req.ProjectID = projectID
	prefsMu.Lock()
	prefsByProj[projectID] = req
	prefsMu.Unlock()
	writeJSON(w, http.StatusOK, req)
}

func (h *NotificationsHandler) authorizeProject(w http.ResponseWriter, r *http.Request, projectID string) bool {
	orgID, ok := gwmw.OrgIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", "session required")
		return false
	}
	project, err := h.Store.GetProject(r.Context(), projectID)
	if err != nil || project.OrgID != orgID {
		writeError(w, http.StatusNotFound, "not_found", "project not found")
		return false
	}
	return true
}
