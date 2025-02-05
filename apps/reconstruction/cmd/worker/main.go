package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/replay/platform/apps/reconstruction/internal/worker"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	q := newMemoryQueue()
	w := &worker.Worker{
		Queue:    q,
		Handlers: map[string]worker.Handler{},
	}
	log.Println("reconstruction worker started")
	if err := w.ConsumeJobs(ctx); err != nil && err != context.Canceled {
		log.Fatal(err)
	}
}

type memoryQueue struct {
	jobs []worker.Job
}

func newMemoryQueue() *memoryQueue {
	return &memoryQueue{}
}

func (q *memoryQueue) PopJob(ctx context.Context) (worker.Job, error) {
	select {
	case <-ctx.Done():
		return worker.Job{}, ctx.Err()
	default:
	}
	if len(q.jobs) == 0 {
		return worker.Job{}, context.DeadlineExceeded
	}
	j := q.jobs[0]
	q.jobs = q.jobs[1:]
	return j, nil
}
