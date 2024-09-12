package kafka

import (
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/replay/platform/packages/shared-go/event"
)

// CaptureHook persists a captured consumer record.
type CaptureHook func(ev event.CapturedEvent)

// ConsumerWrapper adds pre/post capture around a Sarama handler.
type ConsumerWrapper struct {
	inner   sarama.ConsumerGroupHandler
	hook    CaptureHook
	timeout time.Duration
	enqueue chan event.CapturedEvent
}

type noopHandler struct{}

func (noopHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (noopHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (noopHandler) ConsumeClaim(sarama.ConsumerGroupSession, sarama.ConsumerGroupClaim) error {
	return nil
}

// CaptureTimeout is the default async enqueue timeout when none is provided.
const CaptureTimeout = 10 * time.Millisecond

// WrapConsumerHandler returns a handler that captures records around delegate setup.
func WrapConsumerHandler(inner sarama.ConsumerGroupHandler, hook CaptureHook, captureTimeout time.Duration) *ConsumerWrapper {
	if captureTimeout == 0 {
		captureTimeout = CaptureTimeout
	}
	if inner == nil {
		inner = noopHandler{}
	}
	w := &ConsumerWrapper{
		inner:   inner,
		hook:    hook,
		timeout: captureTimeout,
		enqueue: make(chan event.CapturedEvent, 256),
	}
	go w.drainEnqueue()
	return w
}

func (w *ConsumerWrapper) drainEnqueue() {
	for ev := range w.enqueue {
		if w.hook != nil {
			w.hook(ev)
		}
	}
}

// Setup implements sarama.ConsumerGroupHandler.
func (w *ConsumerWrapper) Setup(sess sarama.ConsumerGroupSession) error {
	return w.inner.Setup(sess)
}

// Cleanup implements sarama.ConsumerGroupHandler.
func (w *ConsumerWrapper) Cleanup(sess sarama.ConsumerGroupSession) error {
	return w.inner.Cleanup(sess)
}

// ConsumeClaim captures each message then marks it processed (non-blocking enqueue).
func (w *ConsumerWrapper) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		w.captureAsync(msg, sess)
		sess.MarkMessage(msg, "")
	}
	return nil
}

func (w *ConsumerWrapper) captureAsync(msg *sarama.ConsumerMessage, sess sarama.ConsumerGroupSession) {
	ev := CaptureRecord(msg, sess)
	if w.timeout > 0 {
		select {
		case w.enqueue <- ev:
		default:
		}
		return
	}
	select {
	case w.enqueue <- ev:
	default:
	}
}

var recordPool = sync.Pool{New: func() any { return make([]byte, 0, 4096) }}

// CaptureRecord builds a CapturedEvent from a consumer message.
func CaptureRecord(msg *sarama.ConsumerMessage, sess sarama.ConsumerGroupSession) event.CapturedEvent {
	now := time.Now().UTC()
	group := ""
	if sess != nil {
		group = sess.MemberID()
	}
	buf := recordPool.Get().([]byte)
	buf = append(buf[:0], msg.Value...)
	payload := string(buf)
	recordPool.Put(buf)
	return event.CapturedEvent{
		EventID:       uuid.NewString(),
		CapturedAt:    now,
		Source:        event.SourceConsumer,
		Topic:         msg.Topic,
		Partition:     msg.Partition,
		Offset:        msg.Offset,
		Timestamp:     msg.Timestamp,
		Key:           string(msg.Key),
		Payload:       payload,
		ConsumerGroup: group,
		Outcome:       event.OutcomeSuccess,
	}
}
