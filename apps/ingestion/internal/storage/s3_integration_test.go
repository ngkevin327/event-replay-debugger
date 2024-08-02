//go:build integration

package storage_test

import (
	"context"
	"os"
	"testing"

	"github.com/replay/platform/apps/ingestion/internal/config"
	"github.com/replay/platform/apps/ingestion/internal/storage"
)

func TestPayloadRoundtrip(t *testing.T) {
	if os.Getenv("S3_ENDPOINT") == "" {
		t.Skip("S3_ENDPOINT not set")
	}
	cfg := config.Load()
	client, err := storage.NewS3Client(cfg)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	key := storage.PayloadKey("org", "proj", "deadbeef")
	body := []byte("payload-bytes")
	_, err = client.PutObject(ctx, key, body)
	if err != nil {
		t.Fatal(err)
	}
	got, err := client.GetObject(ctx, key)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(body) {
		t.Fatalf("got %q", got)
	}
}
