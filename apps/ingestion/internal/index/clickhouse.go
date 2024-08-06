package index

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// CHClient wraps ClickHouse HTTP access.
type CHClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewCHClient creates a client for the given HTTP endpoint.
func NewCHClient(url string) *CHClient {
	return &CHClient{
		baseURL: url,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Ping checks ClickHouse availability.
func (c *CHClient) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/ping", nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("clickhouse ping: %s", resp.Status)
	}
	return nil
}
