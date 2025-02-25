package timeline

import (
	"testing"

	"github.com/replay/platform/apps/reconstruction/internal/fetch"
)

func TestDetectGaps(t *testing.T) {
	rows := []fetch.EventRow{
		{Topic: "payments", Partition: 0, Offset: 1},
		{Topic: "payments", Partition: 0, Offset: 3},
	}
	gaps := DetectGaps(rows)
	if len(gaps) != 1 || gaps[0].Start != 2 {
		t.Fatalf("got %+v", gaps)
	}
}
