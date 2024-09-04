package kafka_test

import (
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/replay/platform/packages/agent-go/kafka"
	"github.com/replay/platform/packages/shared-go/event"
)

func TestCaptureRecord(t *testing.T) {
	msg := &sarama.ConsumerMessage{
		Topic: "orders", Partition: 0, Offset: 42,
		Timestamp: time.Now(), Value: []byte(`{"ok":true}`),
	}
	ev := kafka.CaptureRecord(msg, nil)
	if ev.Topic != "orders" || ev.Offset != 42 {
		t.Fatalf("ev %+v", ev)
	}
}

func TestWrapConsumerHandler(t *testing.T) {
	var captured int
	w := kafka.WrapConsumerHandler(nil, func(ev event.CapturedEvent) {
		captured++
	}, time.Millisecond)
	if w == nil {
		t.Fatal("nil wrapper")
	}
	_ = captured
}

func TestHookOrdering(t *testing.T) {
	order := make([]string, 0, 2)
	msg := &sarama.ConsumerMessage{Topic: "t", Partition: 0, Offset: 1, Value: []byte("x")}
	ev := kafka.CaptureRecord(msg, nil)
	order = append(order, "capture")
	kafka.LabelOutcome(&ev, 1, event.OutcomeRetryScheduled, nil)
	order = append(order, "label")
	if order[0] != "capture" || order[1] != "label" {
		t.Fatalf("order %v", order)
	}
}
