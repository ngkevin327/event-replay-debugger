package snapshot

import "testing"

func TestCaptureSnapshots(t *testing.T) {
	events := []struct {
		ConsumerGroup string
		Topic         string
		Partition     uint32
		Offset        uint64
	}{
		{"cg1", "payments", 0, 10},
		{"cg1", "payments", 0, 12},
	}
	snaps := CaptureSnapshots(events)
	if len(snaps) != 1 || snaps[0].Offset != 12 {
		t.Fatalf("got %+v", snaps)
	}
}
