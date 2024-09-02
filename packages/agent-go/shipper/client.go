package shipper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/replay/platform/packages/shared-go/event"
)

// IngestClient posts batches to the ingestion service.
type IngestClient struct {
	url    string
	apiKey string
	client *http.Client
}

// NewIngestClient creates a client for POST /v1/ingest/batch.
func NewIngestClient(url, apiKey string) *IngestClient {
	return &IngestClient{
		url:    url,
		apiKey: apiKey,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// PostBatch ships captured events.
func (c *IngestClient) PostBatch(ctx context.Context, events []event.CapturedEvent) error {
	body, err := json.Marshal(map[string]any{"events": events})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Replay-Key", c.apiKey)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ingest status %d", resp.StatusCode)
	}
	return nil
}
