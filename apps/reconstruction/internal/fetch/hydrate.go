package fetch

import (
	"context"
	"sync"
)

// ObjectStore loads payload bytes for an S3 URI.
type ObjectStore interface {
	GetPayload(ctx context.Context, uri string) ([]byte, error)
}

// HydratePool hydrates event payloads concurrently.
type HydratePool struct {
	Store  ObjectStore
	Workers int
}

// Hydrate fills Payload on each row using a worker pool.
func (p *HydratePool) Hydrate(ctx context.Context, rows []EventRow) error {
	if p.Workers <= 0 {
		p.Workers = 4
	}
	sem := make(chan struct{}, p.Workers)
	var wg sync.WaitGroup
	var firstErr error
	var mu sync.Mutex
	for i := range rows {
		if rows[i].S3URI == "" {
			continue
		}
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int) {
			defer wg.Done()
			defer func() { <-sem }()
			b, err := p.Store.GetPayload(ctx, rows[idx].S3URI)
			if err != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
				return
			}
			rows[idx].Payload = b
		}(i)
	}
	wg.Wait()
	return firstErr
}
