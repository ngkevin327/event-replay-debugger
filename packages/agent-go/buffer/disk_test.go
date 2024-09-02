package buffer_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/replay/platform/packages/agent-go/buffer"
	"github.com/replay/platform/packages/shared-go/event"
)

func TestDiskBufferDropOldest(t *testing.T) {
	dir := t.TempDir()
	b, err := buffer.NewDiskBuffer(dir, 2)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		ev := event.CapturedEvent{
			EventID:    fmt.Sprintf("event-%d", i),
			ProjectID:  "p",
			CapturedAt: time.Now(),
			Source:     event.SourceConsumer,
			Topic:      "t",
		}
		if err := b.Append(ev); err != nil {
			t.Fatal(err)
		}
	}
	got, err := b.Drain()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("len %d", len(got))
	}
}
