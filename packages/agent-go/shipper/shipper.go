package shipper

import (
	"context"
	"time"

	"github.com/replay/platform/packages/agent-go/buffer"
)

// BatchThreshold configures flush triggers.
type BatchThreshold struct {
	MaxEvents int
	Interval  time.Duration
}

// Shipper runs the flush loop from disk buffer to ingest.
type Shipper struct {
	buf       *buffer.DiskBuffer
	client    *IngestClient
	threshold BatchThreshold
}

// NewShipper wires buffer and ingest client.
func NewShipper(buf *buffer.DiskBuffer, client *IngestClient, threshold BatchThreshold) *Shipper {
	return &Shipper{buf: buf, client: client, threshold: threshold}
}

// FlushLoop periodically drains the buffer and posts batches.
func (s *Shipper) FlushLoop(ctx context.Context) error {
	ticker := time.NewTicker(s.threshold.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := s.flush(ctx); err != nil {
				return err
			}
		}
	}
}

func (s *Shipper) flush(ctx context.Context) error {
	events, err := s.buf.Drain()
	if err != nil || len(events) == 0 {
		return err
	}
	return s.client.PostBatch(ctx, events)
}
