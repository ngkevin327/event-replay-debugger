package store

import (
	"context"
	"time"
)

// ReplayStatus is the replay run lifecycle state.
type ReplayStatus string

const (
	ReplayPending   ReplayStatus = "pending"
	ReplayRunning   ReplayStatus = "running"
	ReplaySucceeded ReplayStatus = "succeeded"
	ReplayFailed    ReplayStatus = "failed"
	ReplayDiverged  ReplayStatus = "diverged"
	ReplayCancelled ReplayStatus = "cancelled"
)

// ReplayRun is a replay execution record.
type ReplayRun struct {
	ID               string
	IncidentID       string
	OrgID            string
	Status           ReplayStatus
	TimingMode       string
	DivergenceIndex  *int
	ReportURI        *string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// CreateReplayRun inserts a pending replay for an incident.
func (s *Store) CreateReplayRun(ctx context.Context, incidentID, orgID, timingMode string) (ReplayRun, error) {
	var run ReplayRun
	err := s.pool.QueryRow(ctx,
		`INSERT INTO replay_runs (incident_id, org_id, timing_mode, status)
		 VALUES ($1, $2, $3, 'pending')
		 RETURNING id, incident_id, org_id, status, timing_mode, divergence_index, report_uri, created_at, updated_at`,
		incidentID, orgID, timingMode,
	).Scan(
		&run.ID, &run.IncidentID, &run.OrgID, &run.Status, &run.TimingMode,
		&run.DivergenceIndex, &run.ReportURI, &run.CreatedAt, &run.UpdatedAt,
	)
	return run, err
}

// GetReplayRun loads a replay by id scoped to org.
func (s *Store) GetReplayRun(ctx context.Context, orgID, replayID string) (ReplayRun, error) {
	var run ReplayRun
	err := s.pool.QueryRow(ctx,
		`SELECT id, incident_id, org_id, status, timing_mode, divergence_index, report_uri, created_at, updated_at
		 FROM replay_runs WHERE id = $1 AND org_id = $2`,
		replayID, orgID,
	).Scan(
		&run.ID, &run.IncidentID, &run.OrgID, &run.Status, &run.TimingMode,
		&run.DivergenceIndex, &run.ReportURI, &run.CreatedAt, &run.UpdatedAt,
	)
	return run, err
}

// CancelReplayRun marks a replay cancelled.
func (s *Store) CancelReplayRun(ctx context.Context, replayID string) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE replay_runs SET status = 'cancelled', updated_at = NOW()
		 WHERE id = $1 AND status IN ('pending', 'running')`,
		replayID,
	)
	return err
}
