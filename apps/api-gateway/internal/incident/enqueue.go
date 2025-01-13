package incident

import (
	"github.com/replay/platform/apps/api-gateway/internal/queue"
)

// EnqueueCollection schedules background collection for an incident.
func EnqueueCollection(q *queue.RedisQueue, incidentID string) error {
	return q.PushJob(queue.CollectionJob{IncidentID: incidentID})
}
