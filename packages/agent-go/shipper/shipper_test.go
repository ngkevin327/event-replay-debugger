package shipper_test

import (
	"testing"
	"time"

	"github.com/replay/platform/packages/agent-go/shipper"
)

func TestBatchThresholdDefaults(t *testing.T) {
	th := shipper.BatchThreshold{MaxEvents: 100, Interval: time.Second}
	if th.MaxEvents != 100 {
		t.Fatal("unexpected threshold")
	}
}
