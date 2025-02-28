package timeline

import "github.com/replay/platform/apps/reconstruction/internal/fetch"

// RetryChain groups related retry attempts.
type RetryChain struct {
	CorrelationID string         `json:"correlation_id"`
	PayloadHash   string         `json:"payload_hash"`
	Events        []fetch.EventRow `json:"events"`
}

// GroupRetries clusters events by correlation id and payload hash.
func GroupRetries(rows []fetch.EventRow) []RetryChain {
	type key struct {
		corr, hash string
	}
	buckets := map[key][]fetch.EventRow{}
	for _, row := range rows {
		k := key{row.CorrelationID, row.PayloadHash}
		buckets[k] = append(buckets[k], row)
	}
	out := make([]RetryChain, 0, len(buckets))
	for k, evs := range buckets {
		out = append(out, RetryChain{
			CorrelationID: k.corr,
			PayloadHash:   k.hash,
			Events:        evs,
		})
	}
	return out
}
