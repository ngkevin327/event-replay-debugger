package collector

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type fakeUpdater struct {
	status     IncidentStatus
	eventCount int64
	coverage   *float64
}

func (f *fakeUpdater) UpdateStatus(_ context.Context, _ string, status IncidentStatus, eventCount int64, coverage *float64) error {
	f.status = status
	f.eventCount = eventCount
	f.coverage = coverage
	return nil
}

func TestStatusTransitionEmptyWindow(t *testing.T) {
	srv := httptestCountServer("0\n")
	defer srv.Close()
	up := &fakeUpdater{}
	c := &Collector{Counter: NewCounter(srv.URL), Store: up}
	meta := IncidentMeta{
		ID: "inc-1", ProjectID: "proj-1",
		WindowStart: time.Now().Add(-time.Hour), WindowEnd: time.Now(),
	}
	if err := c.RunCollection(context.Background(), meta); err != nil {
		t.Fatal(err)
	}
	if up.status != StatusFailed {
		t.Fatalf("got %s", up.status)
	}
}

func TestStatusTransitionPartialCoverage(t *testing.T) {
	srv := httptestCountServer("5\n")
	defer srv.Close()
	up := &fakeUpdater{}
	c := &Collector{Counter: NewCounter(srv.URL), Store: up}
	meta := IncidentMeta{
		ID: "inc-2", ProjectID: "proj-1",
		WindowStart: time.Now().Add(-2 * time.Hour), WindowEnd: time.Now(),
		TopicFilters: []string{"payments"},
	}
	if err := c.RunCollection(context.Background(), meta); err != nil {
		t.Fatal(err)
	}
	if up.status != StatusReady {
		t.Fatalf("got %s", up.status)
	}
	if up.eventCount != 5 {
		t.Fatalf("count %d", up.eventCount)
	}
	if up.coverage == nil || *up.coverage != 100 {
		t.Fatalf("coverage %+v", up.coverage)
	}
}

func httptestCountServer(body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(body))
	}))
}
