package fetch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// EventRow is a ClickHouse index row with optional hydrated payload.
type EventRow struct {
	EventID         string
	ProjectID       string
	Topic           string
	Partition       uint32
	Offset          uint64
	Timestamp       time.Time
	CapturedAt      time.Time
	CorrelationID   string
	PayloadHash     string
	S3URI           string
	RetryGeneration int
	Outcome         string
	ServiceName     string
	Payload         []byte
	ArrivalIndex    int
}

// CHFetcher loads events for an incident window.
type CHFetcher struct {
	baseURL    string
	httpClient *http.Client
}

// NewCHFetcher creates a ClickHouse HTTP client.
func NewCHFetcher(clickhouseURL string) *CHFetcher {
	return &CHFetcher{
		baseURL: strings.TrimRight(clickhouseURL, "/"),
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// FetchByWindow returns events for project/time range and optional topics.
func (f *CHFetcher) FetchByWindow(ctx context.Context, projectID string, windowStart, windowEnd time.Time, topicFilters []string) ([]EventRow, error) {
	q := fmt.Sprintf(
		`SELECT event_id, topic, partition, offset, timestamp, captured_at, correlation_id, payload_hash, s3_uri, retry_generation, outcome
		 FROM replay.events
		 WHERE project_id = '%s' AND timestamp >= '%s' AND timestamp <= '%s'`,
		projectID,
		windowStart.UTC().Format("2006-01-02 15:04:05"),
		windowEnd.UTC().Format("2006-01-02 15:04:05"),
	)
	if len(topicFilters) > 0 {
		quoted := make([]string, len(topicFilters))
		for i, t := range topicFilters {
			quoted[i] = "'" + strings.ReplaceAll(t, "'", "\\'") + "'"
		}
		q += " AND topic IN (" + strings.Join(quoted, ",") + ")"
	}
	q += " ORDER BY topic, partition, offset FORMAT JSON"
	reqURL := f.baseURL + "/?query=" + url.QueryEscape(q)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("clickhouse fetch: %s %s", resp.Status, string(body))
	}
	var parsed struct {
		Data []map[string]any `json:"data"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, err
	}
	out := make([]EventRow, 0, len(parsed.Data))
	for _, row := range parsed.Data {
		out = append(out, mapRow(projectID, row))
	}
	return out, nil
}

func mapRow(projectID string, row map[string]any) EventRow {
	return EventRow{
		EventID:         str(row["event_id"]),
		ProjectID:       projectID,
		Topic:           str(row["topic"]),
		Partition:       uint32(num(row["partition"])),
		Offset:          uint64(num(row["offset"])),
		Timestamp:       parseTime(str(row["timestamp"])),
		CapturedAt:      parseTime(str(row["captured_at"])),
		CorrelationID:   str(row["correlation_id"]),
		PayloadHash:     str(row["payload_hash"]),
		S3URI:           str(row["s3_uri"]),
		RetryGeneration: int(num(row["retry_generation"])),
		Outcome:         str(row["outcome"]),
	}
}

func str(v any) string {
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	default:
		return fmt.Sprint(t)
	}
}

func num(v any) float64 {
	if v == nil {
		return 0
	}
	switch t := v.(type) {
	case float64:
		return t
	case int:
		return float64(t)
	default:
		return 0
	}
}

func parseTime(s string) time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", s)
	return t
}
