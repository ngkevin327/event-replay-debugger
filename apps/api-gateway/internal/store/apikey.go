package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

// APIKey is a hashed project credential.
type APIKey struct {
	ID        string
	ProjectID string
	Name      string
	Prefix    string
	Hash      string
	Scopes    []string
	RevokedAt *time.Time
	CreatedAt time.Time
}

// CreateAPIKey inserts a hashed key row.
func (s *Store) CreateAPIKey(ctx context.Context, projectID, name, prefix, hash string, scopes []string) (APIKey, error) {
	var k APIKey
	err := s.pool.QueryRow(ctx,
		`INSERT INTO api_keys (project_id, name, key_prefix, key_hash, scopes)
		 VALUES ($1, $2, $3, $4, $5::api_key_scope[])
		 RETURNING id, project_id, name, key_prefix, key_hash,
		           ARRAY(SELECT unnest(scopes)::text), revoked_at, created_at`,
		projectID, name, prefix, hash, scopes,
	).Scan(&k.ID, &k.ProjectID, &k.Name, &k.Prefix, &k.Hash, &k.Scopes, &k.RevokedAt, &k.CreatedAt)
	return k, err
}

// GetAPIKeyByPrefix loads an active key by prefix.
func (s *Store) GetAPIKeyByPrefix(ctx context.Context, prefix string) (APIKey, error) {
	var k APIKey
	err := s.pool.QueryRow(ctx,
		`SELECT id, project_id, name, key_prefix, key_hash,
		        ARRAY(SELECT unnest(scopes)::text), revoked_at, created_at
		 FROM api_keys WHERE key_prefix = $1 AND revoked_at IS NULL`,
		prefix,
	).Scan(&k.ID, &k.ProjectID, &k.Name, &k.Prefix, &k.Hash, &k.Scopes, &k.RevokedAt, &k.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return APIKey{}, err
	}
	return k, err
}

// RevokeAPIKey marks a key revoked.
func (s *Store) RevokeAPIKey(ctx context.Context, keyID string) error {
	_, err := s.pool.Exec(ctx, `UPDATE api_keys SET revoked_at = NOW() WHERE id = $1`, keyID)
	return err
}
