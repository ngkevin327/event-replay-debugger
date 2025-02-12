package fetch

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetchByWindow(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"data":[{"event_id":"e1","topic":"payments","partition":0,"offset":1,"timestamp":"2026-05-19 10:00:00","captured_at":"2026-05-19 10:00:01","correlation_id":"c1","payload_hash":"h1","s3_uri":"s3://b/k","retry_generation":0,"outcome":"success"}]}`))
	}))
	defer srv.Close()
	f := NewCHFetcher(srv.URL)
	rows, err := f.FetchByWindow(context.Background(), "proj-1", time.Now().Add(-time.Hour), time.Now(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Topic != "payments" {
		t.Fatalf("got %+v", rows)
	}
}
