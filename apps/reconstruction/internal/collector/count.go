package collector

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Counter queries ClickHouse for incident windows.
type Counter struct {
	baseURL    string
	httpClient *http.Client
}

// NewCounter creates a ClickHouse HTTP counter.
func NewCounter(clickhouseURL string) *Counter {
	return &Counter{
		baseURL: strings.TrimRight(clickhouseURL, "/"),
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// CountEventsInWindow returns events matching project, window, and optional topics.
func (c *Counter) CountEventsInWindow(ctx context.Context, projectID string, windowStart, windowEnd time.Time, topicFilters []string) (int64, error) {
	q := fmt.Sprintf(
		`SELECT count() FROM replay.events WHERE project_id = '%s' AND timestamp >= '%s' AND timestamp <= '%s'`,
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
	reqURL := c.baseURL + "/?query=" + url.QueryEscape(q)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return 0, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode >= 300 {
		return 0, fmt.Errorf("clickhouse count: %s %s", resp.Status, string(body))
	}
	var count int64
	_, err = fmt.Sscanf(strings.TrimSpace(string(body)), "%d", &count)
	return count, err
}
