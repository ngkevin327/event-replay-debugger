package lifecycle

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Register posts agent metadata to the control plane.
func Register(ctx context.Context, baseURL, projectID, hostname, version string, topics []string) (string, error) {
	body, _ := json.Marshal(map[string]any{
		"project_id": projectID,
		"hostname":   hostname,
		"version":    version,
		"topics":     topics,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/v1/agents/register", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("register status %d", resp.StatusCode)
	}
	var out struct {
		AgentID string `json:"agent_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	return out.AgentID, nil
}

// HeartbeatLoop sends periodic heartbeats until context cancel.
func HeartbeatLoop(ctx context.Context, baseURL, agentID string, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			body, _ := json.Marshal(map[string]any{"agent_id": agentID, "events_shipped": 0})
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/v1/agents/heartbeat", bytes.NewReader(body))
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}
			resp.Body.Close()
		}
	}
}
