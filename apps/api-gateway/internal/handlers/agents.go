package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	gwmw "github.com/replay/platform/apps/api-gateway/internal/middleware"
	"github.com/replay/platform/apps/api-gateway/internal/store"
)

// AgentsHandler manages agent registration and heartbeats.
type AgentsHandler struct {
	Store *store.Store
}

type registerAgentRequest struct {
	ProjectID string   `json:"project_id"`
	Hostname  string   `json:"hostname"`
	Version   string   `json:"version"`
	Topics    []string `json:"topics"`
}

// RegisterAgent handles POST /v1/agents/register.
func (h *AgentsHandler) RegisterAgent(w http.ResponseWriter, r *http.Request) {
	var req registerAgentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid body")
		return
	}
	agent, err := h.Store.RegisterAgent(r.Context(), req.ProjectID, req.Hostname, req.Version, req.Topics)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not register agent")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"agent_id": agent.ID})
}

type heartbeatRequest struct {
	AgentID       string `json:"agent_id"`
	EventsShipped int    `json:"events_shipped"`
}

// HeartbeatAgent handles POST /v1/agents/heartbeat.
func (h *AgentsHandler) HeartbeatAgent(w http.ResponseWriter, r *http.Request) {
	var req heartbeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid body")
		return
	}
	agent, err := h.Store.UpsertHeartbeat(r.Context(), req.AgentID, req.EventsShipped)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "agent not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"agent_id": agent.ID,
		"status":   store.AgentStatus(agent.LastHeartbeatAt),
	})
}

type agentView struct {
	ID        string   `json:"id"`
	ProjectID string   `json:"project_id"`
	Hostname  string   `json:"hostname"`
	Version   string   `json:"version"`
	Topics    []string `json:"topics"`
	Status    string   `json:"status"`
}

// ListAgents handles GET /v1/agents.
func (h *AgentsHandler) ListAgents(w http.ResponseWriter, r *http.Request) {
	projectID := r.URL.Query().Get("project_id")
	if projectID == "" {
		writeError(w, http.StatusBadRequest, "validation_error", "project_id required")
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
	agents, err := h.Store.ListAgentsByProject(r.Context(), projectID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not list agents")
		return
	}
	views := make([]agentView, 0, len(agents))
	for _, a := range agents {
		views = append(views, agentView{
			ID: a.ID, ProjectID: a.ProjectID, Hostname: a.Hostname,
			Version: a.Version, Topics: a.Topics, Status: store.AgentStatus(a.LastHeartbeatAt),
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"agents": views, "as_of": time.Now().UTC()})
}


