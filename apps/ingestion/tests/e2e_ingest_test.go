//go:build e2e

package tests_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/replay/platform/apps/ingestion/internal/auth"
	"github.com/replay/platform/apps/ingestion/internal/batch"
	"github.com/replay/platform/apps/ingestion/internal/config"
	"github.com/replay/platform/apps/ingestion/internal/handlers"
	"github.com/replay/platform/apps/ingestion/internal/index"
	"github.com/replay/platform/apps/ingestion/internal/quota"
	"github.com/replay/platform/apps/ingestion/internal/server"
	"github.com/replay/platform/apps/ingestion/internal/storage"
	"github.com/replay/platform/packages/shared-go/event"
)

func TestE2EIngestBatch(t *testing.T) {
	if os.Getenv("E2E_INGEST") == "" {
		t.Skip("set E2E_INGEST=1 with docker-compose services running")
	}
	cfg := config.Load()
	ctx := context.Background()
	validator, err := auth.NewValidator(ctx)
	if err != nil {
		t.Fatal(err)
	}
	s3, err := storage.NewS3Client(cfg)
	if err != nil {
		t.Fatal(err)
	}
	ch := index.NewCHClient(cfg.ClickHouseURL)
	if err := ch.Ping(ctx); err != nil {
		t.Skip("clickhouse not available:", err)
	}
	handler := &handlers.IngestHandler{
		Validator: validator,
		Uploader:  &batch.Uploader{S3: s3, Workers: 2, OrgID: "00000000-0000-0000-0000-000000000001"},
		Writer:    index.NewWriter(ch),
		Dedup:     index.NewDedup(),
		Quota:     quota.NewEnforcer(1_000_000),
		OrgID:     "00000000-0000-0000-0000-000000000001",
	}
	srv := server.New(":0", server.Deps{Ingest: handler})
	body, _ := json.Marshal(map[string]any{
		"events": []event.CapturedEvent{{
			EventID:    "11111111-1111-1111-1111-111111111111",
			ProjectID:  os.Getenv("E2E_PROJECT_ID"),
			CapturedAt: time.Now().UTC(),
			Source:     event.SourceConsumer,
			Topic:      "orders",
			Partition:  0,
			Offset:     1,
			Timestamp:  time.Now().UTC(),
			Payload:    `{"ok":true}`,
		}},
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/ingest/batch", bytes.NewReader(body))
	req.Header.Set("X-Replay-Key", os.Getenv("E2E_API_KEY"))
	rec := httptest.NewRecorder()
	srv.Router().ServeHTTP(rec, req)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("status %d body %s", rec.Code, rec.Body.String())
	}
}
