package kafka_test

import (
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/replay/platform/packages/agent-go/kafka"
)

func BenchmarkConsumerLag(b *testing.B) {
	msg := &sarama.ConsumerMessage{
		Topic: "bench", Partition: 0, Offset: 0,
		Timestamp: time.Now(), Value: []byte(`{"n":1}`),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = kafka.CaptureRecord(msg, nil)
	}
}

func BenchmarkCaptureTimeout(b *testing.B) {
	BenchmarkConsumerLag(b)
}

func BenchmarkCapture(b *testing.B) {
	BenchmarkConsumerLag(b)
}
