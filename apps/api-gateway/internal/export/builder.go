package export

import (
	"time"

	"github.com/replay/platform/apps/api-gateway/internal/store"
)

// Summary is the JSON export payload for an incident.
type Summary struct {
	IncidentID      string    `json:"incident_id"`
	ProjectID       string    `json:"project_id"`
	Status          string    `json:"status"`
	WindowStart     time.Time `json:"window_start"`
	WindowEnd       time.Time `json:"window_end"`
	EventCount      int64     `json:"event_count"`
	CoveragePercent float64   `json:"coverage_percent,omitempty"`
	TopicFilters    []string  `json:"topic_filters"`
	ExportedAt      time.Time `json:"exported_at"`
	TimelineStats   Stats     `json:"timeline_stats"`
}

// Stats summarizes timeline metadata.
type Stats struct {
	EventCount int64 `json:"event_count"`
	GapCount   int   `json:"gap_count"`
}

// BuildSummary constructs export JSON from incident metadata.
func BuildSummary(inc store.Incident) Summary {
	coverage := 0.0
	if inc.CoveragePercent != nil {
		coverage = *inc.CoveragePercent
	}
	return Summary{
		IncidentID:      inc.ID,
		ProjectID:       inc.ProjectID,
		Status:          string(inc.Status),
		WindowStart:     inc.WindowStart,
		WindowEnd:       inc.WindowEnd,
		EventCount:      inc.EventCount,
		CoveragePercent: coverage,
		TopicFilters:    inc.TopicFilters,
		ExportedAt:      time.Now().UTC(),
		TimelineStats: Stats{
			EventCount: inc.EventCount,
			GapCount:   0,
		},
	}
}

// MaxExportBytes limits export payload size.
const MaxExportBytes = 1 << 20
