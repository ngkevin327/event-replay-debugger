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
