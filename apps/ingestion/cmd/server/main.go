package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/replay/platform/apps/ingestion/internal/auth"
	"github.com/replay/platform/apps/ingestion/internal/batch"
	"github.com/replay/platform/apps/ingestion/internal/config"
	"github.com/replay/platform/apps/ingestion/internal/handlers"
	"github.com/replay/platform/apps/ingestion/internal/index"
	"github.com/replay/platform/apps/ingestion/internal/quota"
	"github.com/replay/platform/apps/ingestion/internal/server"
	"github.com/replay/platform/apps/ingestion/internal/storage"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	deps := server.Deps{Config: cfg}
	if v, err := auth.NewValidator(ctx); err == nil {
		deps.Validator = v
	}
	if s3, err := storage.NewS3Client(cfg); err == nil {
		deps.S3 = s3
		deps.Uploader = &batch.Uploader{S3: s3, Workers: 4}
	}
	ch := index.NewCHClient(cfg.ClickHouseURL)
	deps.CH = ch
	deps.Writer = index.NewWriter(ch)
	deps.Dedup = index.NewDedup()
	deps.Quota = quota.NewEnforcer(0)
	if deps.Validator != nil && deps.Uploader != nil {
		deps.Ingest = &handlers.IngestHandler{
			Validator: deps.Validator,
			Uploader:  deps.Uploader,
			Writer:    deps.Writer,
			Dedup:     deps.Dedup,
			Quota:     deps.Quota,
			OrgID:     os.Getenv("DEFAULT_ORG_ID"),
		}
	}

	srv := server.New(cfg.HTTPAddr, deps)
	runCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("ingestion listening", "addr", cfg.HTTPAddr)
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			slog.Error("server stopped", "err", err)
			stop()
		}
	}()

	<-runCtx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown", "err", err)
		os.Exit(1)
	}
}
