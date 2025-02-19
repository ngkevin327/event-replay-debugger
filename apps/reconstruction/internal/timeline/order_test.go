package timeline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/replay/platform/apps/reconstruction/internal/fetch"
)

func TestPrimaryOrder(t *testing.T) {
	rows := []fetch.EventRow{
		{Topic: "a", Partition: 0, Offset: 2},
		{Topic: "a", Partition: 0, Offset: 1},
	}
	out := OrderByPartitionOffset(rows)
	if out[0].Offset != 1 {
		t.Fatal()
	}
}

func TestArrivalAnnotation(t *testing.T) {
	now := time.Now()
	rows := []fetch.EventRow{
		{CapturedAt: now.Add(time.Second)},
		{CapturedAt: now},
	}
	out := ArrivalOrderIndex(rows)
	if out[0].ArrivalIndex != 0 || out[1].ArrivalIndex != 1 {
		t.Fatalf("%+v", out)
	}
}

func TestOutOfOrderCaptureFixture(t *testing.T) {
	path := filepath.Join("..", "..", "..", "test", "golden", "out_of_order_capture.json")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Skip("fixture missing")
	}
	var rows []fetch.EventRow
	if err := json.Unmarshal(b, &rows); err != nil {
		t.Fatal(err)
	}
	out := OrderByPartitionOffset(rows)
	if len(out) != len(rows) {
		t.Fatal()
	}
}
