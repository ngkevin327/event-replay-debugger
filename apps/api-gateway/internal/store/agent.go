package store

import (
	"context"
	"encoding/json"
	"time"
)

// Agent is a registered capture agent.
type Agent struct {
	ID              string
	ProjectID       string
	Hostname        string
	Version         string
	Topics          []string
	LastHeartbeatAt *time.Time
	CreatedAt       time.Time
}

// RegisterAgent inserts a new agent row.
func (s *Store) RegisterAgent(ctx context.Context, projectID, hostname, version string, topics []string) (Agent, error) {
	topicsJSON, _ := json.Marshal(topics)
	var a Agent
	var topicsRaw []byte
	err := s.pool.QueryRow(ctx,
		`INSERT INTO agents (project_id, hostname, version, topics)
		 VALUES ($1, $2, $3, $4::jsonb)
		 RETURNING id, project_id, hostname, version, topics, last_heartbeat_at, created_at`,
		projectID, hostname, version, topicsJSON,
	).Scan(&a.ID, &a.ProjectID, &a.Hostname, &a.Version, &topicsRaw, &a.LastHeartbeatAt, &a.CreatedAt)
	if err != nil {
		return Agent{}, err
	}
	_ = json.Unmarshal(topicsRaw, &a.Topics)
	return a, nil
}

// UpsertHeartbeat updates agent heartbeat metadata.
func (s *Store) UpsertHeartbeat(ctx context.Context, agentID string, eventsShipped int) (Agent, error) {
	var a Agent
	var topicsRaw []byte
	err := s.pool.QueryRow(ctx,
		`UPDATE agents SET last_heartbeat_at = NOW(), updated_at = NOW()
		 WHERE id = $1
		 RETURNING id, project_id, hostname, version, topics, last_heartbeat_at, created_at`,
		agentID,
	).Scan(&a.ID, &a.ProjectID, &a.Hostname, &a.Version, &topicsRaw, &a.LastHeartbeatAt, &a.CreatedAt)
	if err != nil {
		return Agent{}, err
	}
	_ = json.Unmarshal(topicsRaw, &a.Topics)
	return a, nil
}

// ListAgentsByProject returns agents for a project.
func (s *Store) ListAgentsByProject(ctx context.Context, projectID string) ([]Agent, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, project_id, hostname, version, topics, last_heartbeat_at, created_at
		 FROM agents WHERE project_id = $1 ORDER BY created_at DESC`, projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Agent
	for rows.Next() {
		var a Agent
		var topicsRaw []byte
		if err := rows.Scan(&a.ID, &a.ProjectID, &a.Hostname, &a.Version, &topicsRaw, &a.LastHeartbeatAt, &a.CreatedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(topicsRaw, &a.Topics)
		out = append(out, a)
	}
	return out, rows.Err()
}

// AgentStatus computes health from last heartbeat.
func AgentStatus(last *time.Time) string {
	if last == nil {
		return "offline"
	}
	if time.Since(*last) > 15*time.Minute {
		return "offline"
	}
	if time.Since(*last) > 5*time.Minute {
		return "degraded"
	}
	return "healthy"
}
