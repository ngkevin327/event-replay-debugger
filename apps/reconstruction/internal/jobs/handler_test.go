package jobs

import (
	"context"
	"testing"
)

func TestDispatch(t *testing.T) {
	var got Job
	d := NewDispatcher()
	d.Register(JobCollect, func(ctx context.Context, job Job) error {
		got = job
		return nil
	})
	_ = d.Dispatch(context.Background(), Job{IncidentID: "inc-1", Type: JobCollect})
	if got.IncidentID != "inc-1" {
		t.Fatalf("got %+v", got)
	}
}
