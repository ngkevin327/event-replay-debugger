package batch

import (
	"context"
	"sync"

	"github.com/replay/platform/apps/ingestion/internal/payload"
	"github.com/replay/platform/apps/ingestion/internal/storage"
	"github.com/replay/platform/packages/shared-go/event"
)

// Uploader parallelizes payload uploads within a batch.
type Uploader struct {
	S3      *storage.S3Client
	OrgID   string
	Workers int
}

// UploadBatch processes events with a worker pool and returns updated events.
func (u *Uploader) UploadBatch(ctx context.Context, events []event.CapturedEvent) ([]event.CapturedEvent, error) {
	if u.Workers <= 0 {
		u.Workers = 4
	}
	out := make([]event.CapturedEvent, len(events))
	copy(out, events)

	var wg sync.WaitGroup
	errOnce := make(chan error, 1)
	sem := make(chan struct{}, u.Workers)

	for i := range out {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			ev := &out[idx]
			data, hash, truncated := payload.ProcessPayload([]byte(ev.Payload))
			payload.ApplyToEvent(ev, data, hash, truncated)
			key := storage.PayloadKey(u.OrgID, ev.ProjectID, hash)
			uri, err := u.S3.PutObject(ctx, key, data)
			if err != nil {
				select {
				case errOnce <- err:
				default:
				}
				return
			}
			ev.S3URI = uri
		}(i)
	}
	wg.Wait()
	select {
	case err := <-errOnce:
		return nil, err
	default:
	}
	return out, nil
}
