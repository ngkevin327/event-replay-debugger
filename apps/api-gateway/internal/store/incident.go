package store

import (
	"context"
	"time"
)

// IncidentStatus is the collection lifecycle state.
type IncidentStatus string

const (
	IncidentCollecting IncidentStatus = "collecting"
	IncidentReady      IncidentStatus = "ready"
	IncidentFailed     IncidentStatus = "failed"
)

// Incident is a time-bounded event collection window.
type Incident struct {
	ID               string
	ProjectID        string
	OrgID            string
	Status           IncidentStatus
	WindowStart      time.Time
	WindowEnd        time.Time
	TopicFilters     []string
	EventCount       int64
	CoveragePercent  *float64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// WindowBounds holds incident time range and topic filters.
type WindowBounds struct {
	WindowStart  time.Time
	WindowEnd    time.Time
	TopicFilters []string
}

// InsertIncident creates a new incident in collecting status.
func (s *Store) InsertIncident(ctx context.Context, projectID, orgID string, bounds WindowBounds) (Incident, error) {
	var inc Incident
	err := s.pool.QueryRow(ctx,
		`INSERT INTO incidents (project_id, org_id, window_start, window_end, topic_filters, status)
		 VALUES ($1, $2, $3, $4, $5, 'collecting')
		 RETURNING id, project_id, org_id, status, window_start, window_end, topic_filters,
		           event_count, coverage_percent, created_at, updated_at`,
		projectID, orgID, bounds.WindowStart, bounds.WindowEnd, bounds.TopicFilters,
	).Scan(
		&inc.ID, &inc.ProjectID, &inc.OrgID, &inc.Status,
		&inc.WindowStart, &inc.WindowEnd, &inc.TopicFilters,
		&inc.EventCount, &inc.CoveragePercent, &inc.CreatedAt, &inc.UpdatedAt,
	)
	return inc, err
}

// GetIncident loads an incident by id scoped to org.
func (s *Store) GetIncident(ctx context.Context, orgID, incidentID string) (Incident, error) {
	var inc Incident
	err := s.pool.QueryRow(ctx,
		`SELECT id, project_id, org_id, status, window_start, window_end, topic_filters,
		        event_count, coverage_percent, created_at, updated_at
		 FROM incidents WHERE id = $1 AND org_id = $2`,
		incidentID, orgID,
	).Scan(
		&inc.ID, &inc.ProjectID, &inc.OrgID, &inc.Status,
		&inc.WindowStart, &inc.WindowEnd, &inc.TopicFilters,
		&inc.EventCount, &inc.CoveragePercent, &inc.CreatedAt, &inc.UpdatedAt,
	)
	return inc, err
}

// ListIncidents returns incidents for a project with pagination.
func (s *Store) ListIncidents(ctx context.Context, orgID, projectID string, limit, offset int) ([]Incident, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := s.pool.Query(ctx,
		`SELECT id, project_id, org_id, status, window_start, window_end, topic_filters,
		        event_count, coverage_percent, created_at, updated_at
		 FROM incidents WHERE org_id = $1 AND project_id = $2
		 ORDER BY created_at DESC LIMIT $3 OFFSET $4`,
		orgID, projectID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Incident
	for rows.Next() {
		var inc Incident
		if err := rows.Scan(
			&inc.ID, &inc.ProjectID, &inc.OrgID, &inc.Status,
			&inc.WindowStart, &inc.WindowEnd, &inc.TopicFilters,
			&inc.EventCount, &inc.CoveragePercent, &inc.CreatedAt, &inc.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, inc)
	}
	return out, rows.Err()
}

// UpdateIncidentStatus sets status and optional metadata.
func (s *Store) UpdateIncidentStatus(ctx context.Context, incidentID string, status IncidentStatus, eventCount int64, coverage *float64) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE incidents SET status = $2, event_count = $3, coverage_percent = $4, updated_at = NOW()
		 WHERE id = $1`,
		incidentID, status, eventCount, coverage,
	)
	return err
}
