package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/replay/platform/apps/api-gateway/internal/config"
	"github.com/replay/platform/apps/api-gateway/internal/db"
	"github.com/replay/platform/apps/api-gateway/internal/logging"
	"github.com/replay/platform/apps/api-gateway/internal/server"
	"github.com/replay/platform/apps/api-gateway/internal/store"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config", "err", err) // before logger init
		os.Exit(1)
	}

	slog.SetDefault(logging.New(cfg.LogLevel))

	deps := server.RouteDeps{JWTSecret: cfg.JWTSecret}
	if cfg.DatabaseURL != "" {
		pool, err := db.OpenPostgres(context.Background(), cfg.DatabaseURL)
		if err != nil {
			slog.Error("database", "err", err)
			os.Exit(1)
		}
		defer pool.Close()
		deps.Store = store.NewStore(pool)
	}

	srv := server.New(cfg, deps)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.Run(); err != nil {
			slog.Error("server stopped", "err", err)
			stop()
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown", "err", err)
		os.Exit(1)
	}
}
