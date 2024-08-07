package index

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/replay/platform/packages/shared-go/event"
)

const latencyThreshold = 500 * time.Millisecond

// Writer inserts event batches into ClickHouse.
type Writer struct {
	client   *CHClient
	lastSlow atomic.Bool
}

// NewWriter binds a ClickHouse client.
func NewWriter(client *CHClient) *Writer {
	return &Writer{client: client}
}

// Backpressure reports whether recent inserts exceeded latency threshold.
func (w *Writer) Backpressure() bool {
	return w.lastSlow.Load()
}

// WriteBatch inserts rows via HTTP JSONEachRow.
func (w *Writer) WriteBatch(ctx context.Context, orgID string, events []event.CapturedEvent) error {
	start := time.Now()
	defer func() {
		if time.Since(start) > latencyThreshold {
			w.lastSlow.Store(true)
		} else {
			w.lastSlow.Store(false)
		}
	}()

	var b strings.Builder
	for _, ev := range events {
		trunc := 0
		if ev.PayloadTruncated {
			trunc = 1
		}
		outcome := ev.Outcome
		if outcome == "" {
			outcome = event.OutcomeSuccess
		}
		b.WriteString(fmt.Sprintf(
			`{"event_id":"%s","project_id":"%s","org_id":"%s","captured_at":"%s","source":"%s","topic":"%s","partition":%d,"offset":%d,"timestamp":"%s","consumer_group":"%s","retry_generation":%d,"outcome":"%s","payload_hash":"%s","payload_truncated":%d,"s3_uri":"%s"}`,
			ev.EventID, ev.ProjectID, orgID, ev.CapturedAt.Format(time.RFC3339Nano),
			ev.Source, ev.Topic, ev.Partition, ev.Offset, ev.Timestamp.Format(time.RFC3339Nano),
			ev.ConsumerGroup, ev.RetryGeneration, outcome, ev.PayloadHash, trunc, ev.S3URI,
		))
		b.WriteByte('\n')
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.client.baseURL+"/?query="+urlQueryEscape("INSERT INTO replay.events FORMAT JSONEachRow"), bytes.NewReader([]byte(b.String())))
	if err != nil {
		return err
	}
	resp, err := w.client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("clickhouse insert: %s %s", resp.Status, string(body))
	}
	return nil
}

func urlQueryEscape(q string) string {
	return strings.ReplaceAll(q, " ", "+")
}
