package incident

import (
	"testing"

	"github.com/replay/platform/apps/api-gateway/internal/queue"
)

func TestEnqueueCollection(t *testing.T) {
	q := queue.NewRedisQueue()
	if err := EnqueueCollection(q, "inc-1"); err != nil {
		t.Fatal(err)
	}
	job, err := q.PopJob()
	if err != nil {
		t.Fatal(err)
	}
	if job.IncidentID != "inc-1" {
		t.Fatalf("got %q", job.IncidentID)
	}
}
