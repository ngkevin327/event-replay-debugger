package webhook

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPDoer performs outbound HTTP requests.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Sender delivers webhook payloads with retry backoff.
type Sender struct {
	Client  HTTPDoer
	Timeout time.Duration
	Retries int
}

// Deliver posts payload to url until success or retries exhausted.
func (s *Sender) Deliver(ctx context.Context, url string, body []byte) error {
	return s.DeliverSigned(ctx, url, body, "")
}

// DeliverSigned posts payload with optional HMAC signature header.
func (s *Sender) DeliverSigned(ctx context.Context, url string, body []byte, secret string) error {
	if s.Client == nil {
		s.Client = &http.Client{Timeout: s.Timeout}
	}
	retries := s.Retries
	if retries <= 0 {
		retries = 3
	}
	var lastErr error
	backoff := 200 * time.Millisecond
	for attempt := 0; attempt <= retries; attempt++ {
		if err := s.sendOnce(ctx, url, body, secret); err != nil {
			lastErr = err
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
				backoff *= 2
				continue
			}
		}
		return nil
	}
	return fmt.Errorf("webhook delivery failed after %d retries: %w", retries, lastErr)
}

func (s *Sender) sendOnce(ctx context.Context, url string, body []byte, secret string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if secret != "" {
		req.Header.Set(signatureHeader, SignPayload(secret, body))
	}
	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	return fmt.Errorf("unexpected status %d", resp.StatusCode)
}
