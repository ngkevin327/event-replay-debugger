package kafka_test

import (
	"errors"
	"testing"

	"github.com/replay/platform/packages/agent-go/kafka"
	"github.com/replay/platform/packages/shared-go/event"
)

func TestLabelOutcome(t *testing.T) {
	ev := &event.CapturedEvent{}
	kafka.LabelOutcome(ev, 1, event.OutcomeRetryScheduled, errors.New("retry"))
	if ev.RetryGeneration != 1 || ev.Outcome != event.OutcomeRetryScheduled {
		t.Fatalf("ev %+v", ev)
	}
}
