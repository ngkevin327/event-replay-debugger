package event_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/replay/platform/packages/shared-go/event"
)

func TestCapturedEventRoundTrip(t *testing.T) {
	orig := event.CapturedEvent{
		EventID:         "550e8400-e29b-41d4-a716-446655440000",
		ProjectID:       "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		CapturedAt:      time.Date(2026, 5, 19, 12, 0, 0, 0, time.UTC),
		Source:          event.SourceConsumer,
		Topic:           "payments.settle",
		Partition:       3,
		Offset:          128901,
		Timestamp:       time.Date(2026, 5, 19, 11, 59, 58, 123000000, time.UTC),
		ConsumerGroup:   "settlement-workers",
		RetryGeneration: 1,
		Outcome:         event.OutcomeSuccess,
	}
	data, err := json.Marshal(orig)
	if err != nil {
		t.Fatal(err)
	}
	var decoded event.CapturedEvent
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Topic != orig.Topic || decoded.Offset != orig.Offset {
		t.Fatalf("mismatch: %+v vs %+v", decoded, orig)
	}
	if decoded.RetryGeneration != 1 {
		t.Fatalf("retry_generation: got %d", decoded.RetryGeneration)
	}
}
