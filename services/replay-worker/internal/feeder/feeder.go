package feeder

import (
	"context"
	"time"
)

// Feeder publishes timeline events to sandbox Kafka.
type Feeder struct {
	IdleWait time.Duration
}

// NewFeeder creates a feeder with defaults.
func NewFeeder() *Feeder {
	return &Feeder{IdleWait: time.Second}
}

// RunLoop drives replay publishing until cancelled.
func (f *Feeder) RunLoop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			time.Sleep(f.IdleWait)
		}
	}
}
