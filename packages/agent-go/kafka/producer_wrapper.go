package kafka

import (
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/replay/platform/packages/shared-go/event"
)

// ProducerWrapper captures sync producer sends and delivery outcomes.
type ProducerWrapper struct {
	inner sarama.SyncProducer
	hook  CaptureHook
}

// WrapSyncProducer wraps a sync producer with capture hooks.
func WrapSyncProducer(inner sarama.SyncProducer, hook CaptureHook) *ProducerWrapper {
	return &ProducerWrapper{inner: inner, hook: hook}
}

// SendMessage sends via inner producer and captures the record.
func (p *ProducerWrapper) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	ev := event.CapturedEvent{
		EventID:    uuid.NewString(),
		CapturedAt: time.Now().UTC(),
		Source:     event.SourceProducer,
		Topic:      msg.Topic,
		Timestamp:  time.Now().UTC(),
		Outcome:    event.OutcomeSuccess,
	}
	if msg.Value != nil {
		if b, ok := msg.Value.(sarama.ByteEncoder); ok {
			ev.Payload = string(b)
		}
	}
	partition, offset, err = p.inner.SendMessage(msg)
	if err != nil {
		ev.Outcome = event.OutcomeError
		OnDelivery(&ev, err)
	} else {
		ev.Partition = partition
		ev.Offset = offset
		OnDelivery(&ev, nil)
	}
	if p.hook != nil {
		p.hook(ev)
	}
	return partition, offset, err
}

// OnDelivery labels producer delivery outcome on the event.
func OnDelivery(ev *event.CapturedEvent, err error) {
	if err != nil {
		ev.Outcome = event.OutcomeError
		ev.Error = &event.CaptureError{Type: "delivery", Message: err.Error()}
	} else {
		ev.Outcome = event.OutcomeSuccess
	}
}

func (p *ProducerWrapper) Close() error { return p.inner.Close() }
