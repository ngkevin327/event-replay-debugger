package timeline

import (
	"encoding/json"
	"time"
)

// Artifact is the timeline JSON stored in S3.
type Artifact struct {
	Version    int          `json:"version"`
	IncidentID string       `json:"incident_id"`
	BuiltAt    time.Time    `json:"built_at"`
	Events     []EventEntry `json:"events"`
	Gaps       []GapMarker  `json:"gaps,omitempty"`
	Retries    []RetryChain `json:"retry_chains,omitempty"`
}

// EventEntry is a compact timeline row.
type EventEntry struct {
	EventID       string `json:"event_id"`
	Topic         string `json:"topic"`
	Partition     uint32 `json:"partition"`
	Offset        uint64 `json:"offset"`
	ArrivalIndex  int    `json:"arrival_index"`
	RetryGeneration int  `json:"retry_generation"`
	Outcome       string `json:"outcome"`
}

// BuildTimelineArtifact serializes ordered timeline data.
func BuildTimelineArtifact(incidentID string, ordered []EventEntry, gaps []GapMarker, chains []RetryChain) ([]byte, error) {
	art := Artifact{
		Version:    1,
		IncidentID: incidentID,
		BuiltAt:    time.Now().UTC(),
		Events:     ordered,
		Gaps:       gaps,
		Retries:    chains,
	}
	return json.Marshal(art)
}
