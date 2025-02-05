package worker

import (
	"context"
	"testing"
	"time"
)

type memQueue struct {
	jobs []Job
}

func (q *memQueue) PopJob(ctx context.Context) (Job, error) {
	if len(q.jobs) == 0 {
		return Job{}, context.DeadlineExceeded
	}
	j := q.jobs[0]
	q.jobs = q.jobs[1:]
	return j, nil
}

type stubHandler struct {
	called bool
}

func (h *stubHandler) Handle(ctx context.Context, job Job) error {
	h.called = true
	return nil
}

func TestConsumeJobs(t *testing.T) {
	h := &stubHandler{}
	w := &Worker{
		Queue: &memQueue{jobs: []Job{{IncidentID: "inc-1", JobType: "collect"}}},
		Handlers: map[string]Handler{
			"collect": h,
		},
		IdleWait: time.Millisecond,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_ = w.ConsumeJobs(ctx)
	if !h.called {
		t.Fatal("handler not called")
	}
}
