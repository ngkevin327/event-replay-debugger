package incident

import (
	"context"
	"testing"
)

type fakePatch struct {
	status IncidentStatus
}

func (f *fakePatch) PatchIncident(_ context.Context, _ string, status IncidentStatus, _ int64, _ *float64) error {
	f.status = status
	return nil
}

func TestOnJobComplete(t *testing.T) {
	p := &fakePatch{}
	_ = OnJobComplete(context.Background(), p, "inc-1", "collect", true, 10, nil)
	if p.status != StatusReady {
		t.Fatalf("got %s", p.status)
	}
}
