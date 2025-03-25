package graph

import (
	"testing"

	"github.com/replay/platform/apps/reconstruction/internal/fetch"
)

func TestEdges(t *testing.T) {
	rows := []fetch.EventRow{
		{Topic: "payments", ServiceName: "payment-worker"},
		{Topic: "ledger", ServiceName: "ledger-svc"},
	}
	g := BuildGraph(rows)
	if len(g.Nodes) < 2 {
		t.Fatalf("nodes %d", len(g.Nodes))
	}
}

func TestFailedNodes(t *testing.T) {
	rows := []fetch.EventRow{
		{Topic: "payments", ServiceName: "payment-worker", Outcome: "error"},
	}
	g := BuildGraph(rows)
	MarkFailed(&g, rows)
	if !g.Nodes[0].Failed {
		t.Fatal("expected failed node")
	}
	if len(CascadePaths(g)) == 0 {
		t.Fatal("expected cascade paths")
	}
}

func TestMultiServiceDemoGraph(t *testing.T) {
	rows := []fetch.EventRow{
		{Topic: "payments", ServiceName: "payment-worker"},
		{Topic: "notifications", ServiceName: "notify-svc"},
		{Topic: "ledger", ServiceName: "ledger-svc"},
	}
	g := BuildGraph(rows)
	if len(g.Nodes) < 3 {
		t.Fatalf("got %d nodes", len(g.Nodes))
	}
}
