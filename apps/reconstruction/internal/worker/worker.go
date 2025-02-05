package worker

import (
	"context"
	"time"
)

// Job represents a reconstruction queue item.
type Job struct {
	IncidentID string
	JobType    string
}

// JobQueue pops reconstruction jobs (Redis-backed in production).
type JobQueue interface {
	PopJob(ctx context.Context) (Job, error)
}

// Handler runs a single job type.
type Handler interface {
	Handle(ctx context.Context, job Job) error
}

// Worker consumes jobs until context cancellation.
type Worker struct {
	Queue    JobQueue
	Handlers map[string]Handler
	IdleWait time.Duration
}

// ConsumeJobs blocks processing queue items.
func (w *Worker) ConsumeJobs(ctx context.Context) error {
	if w.IdleWait == 0 {
		w.IdleWait = time.Second
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		job, err := w.Queue.PopJob(ctx)
		if err != nil {
			time.Sleep(w.IdleWait)
			continue
		}
		h := w.Handlers[job.JobType]
		if h == nil {
			continue
		}
		_ = h.Handle(ctx, job)
	}
}
