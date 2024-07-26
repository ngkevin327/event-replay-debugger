package auth

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	sharedauth "github.com/replay/platform/packages/shared-go/auth"
)

var (
	ErrInvalidAPIKey = errors.New("invalid api key")
	ErrMissingScope  = errors.New("missing ingest scope")
)

// Validator checks API keys against control-plane Postgres metadata.
type Validator struct {
	pool *pgxpool.Pool
}

// NewValidator connects to DATABASE_URL when set.
func NewValidator(ctx context.Context) (*Validator, error) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		return &Validator{}, nil
	}
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	return &Validator{pool: pool}, nil
}

// ValidateAPIKey resolves project binding and ingest scope for a plaintext key.
func (v *Validator) ValidateAPIKey(ctx context.Context, plain string) (projectID string, err error) {
	if v.pool == nil {
		return "", fmt.Errorf("database not configured")
	}
	prefix, err := sharedauth.ParseKeyPrefix(plain)
	if err != nil {
		return "", ErrInvalidAPIKey
	}
	hash, err := sharedauth.HashPrefix(plain)
	if err != nil {
		return "", ErrInvalidAPIKey
	}
	var scopes []string
	err = v.pool.QueryRow(ctx,
		`SELECT project_id, ARRAY(SELECT unnest(scopes)::text)
		 FROM api_keys WHERE key_prefix = $1 AND key_hash = $2 AND revoked_at IS NULL`,
		prefix, hash,
	).Scan(&projectID, &scopes)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrInvalidAPIKey
	}
	if err != nil {
		return "", err
	}
	for _, s := range scopes {
		if s == "ingest" {
			return projectID, nil
		}
	}
	return "", ErrMissingScope
}
