package collector

import (
	"context"
	"time"
)

// IncidentStatus is the collection lifecycle state.
type IncidentStatus string

const (
	StatusCollecting IncidentStatus = "collecting"
	StatusReady      IncidentStatus = "ready"
	StatusFailed     IncidentStatus = "failed"
)

// IncidentMeta describes the collection window.
type IncidentMeta struct {
	ID           string
	ProjectID    string
	WindowStart  time.Time
	WindowEnd    time.Time
	TopicFilters []string
}

// StatusUpdater persists incident lifecycle changes.
type StatusUpdater interface {
	UpdateStatus(ctx context.Context, incidentID string, status IncidentStatus, eventCount int64, coverage *float64) error
}

// Collector runs incident event collection.
type Collector struct {
	Counter *Counter
	Store   StatusUpdater
}

// RunCollection counts events and transitions incident status.
func (c *Collector) RunCollection(ctx context.Context, meta IncidentMeta) error {
	count, err := c.Counter.CountEventsInWindow(ctx, meta.ProjectID, meta.WindowStart, meta.WindowEnd, meta.TopicFilters)
	if err != nil {
		_ = c.Store.UpdateStatus(ctx, meta.ID, StatusFailed, 0, nil)
		return err
	}
	coverage := coveragePercent(count)
	status := StatusReady
	if count == 0 {
		status = StatusFailed
	}
	return c.Store.UpdateStatus(ctx, meta.ID, status, count, &coverage)
}

// UpdateStatus writes the incident status via the configured updater.
func (c *Collector) UpdateStatus(ctx context.Context, incidentID string, status IncidentStatus, eventCount int64, coverage *float64) error {
	return c.Store.UpdateStatus(ctx, incidentID, status, eventCount, coverage)
}

func coveragePercent(eventCount int64) float64 {
	if eventCount == 0 {
		return 0
	}
	return 100
}
