package jobs

// JobType identifies reconstruction work units.
type JobType string

const (
	JobCollect   JobType = "collect"
	JobTimeline  JobType = "timeline"
	JobBuildGraph JobType = "build_graph"
)

// Job is a typed queue payload.
type Job struct {
	IncidentID string
	Type       JobType
}
