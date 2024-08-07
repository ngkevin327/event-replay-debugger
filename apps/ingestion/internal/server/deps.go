package server

import (
	"github.com/replay/platform/apps/ingestion/internal/auth"
	"github.com/replay/platform/apps/ingestion/internal/batch"
	"github.com/replay/platform/apps/ingestion/internal/config"
	"github.com/replay/platform/apps/ingestion/internal/handlers"
	"github.com/replay/platform/apps/ingestion/internal/index"
	"github.com/replay/platform/apps/ingestion/internal/storage"
)

// Deps wires ingestion handlers and infrastructure.
type Deps struct {
	Config    config.Config
	Validator *auth.Validator
	S3        *storage.S3Client
	Uploader  *batch.Uploader
	CH        *index.CHClient
	Writer    *index.Writer
	Ingest    *handlers.IngestHandler
}
