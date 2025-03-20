package handlers

import "testing"

func TestIngestSnapshotRequestShape(t *testing.T) {
	req := ingestSnapshotRequest{ConsumerGroup: "cg", Topic: "payments", Partition: 0, Offset: 1}
	if req.Topic != "payments" {
		t.Fatal()
	}
}
