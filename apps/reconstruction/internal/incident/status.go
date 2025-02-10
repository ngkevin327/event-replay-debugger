package incident

import (
	"context"
)

// IncidentStatus is the lifecycle state written back to Postgres.
type IncidentStatus string

const (
	StatusCollecting IncidentStatus = "collecting"
	StatusReady      IncidentStatus = "ready"
	StatusFailed     IncidentStatus = "failed"
)

// Patcher updates incident rows after jobs complete.
type Patcher interface {
	PatchIncident(ctx context.Context, incidentID string, status IncidentStatus, eventCount int64, coverage *float64) error
}

// OnJobComplete updates incident status after a reconstruction job finishes.
func OnJobComplete(ctx context.Context, patch Patcher, incidentID string, jobType string, ok bool, eventCount int64, coverage *float64) error {
	status := StatusReady
	if !ok {
		status = StatusFailed
	}
	if jobType == "collect" && eventCount == 0 {
		status = StatusFailed
	}
	return patch.PatchIncident(ctx, incidentID, status, eventCount, coverage)
}
