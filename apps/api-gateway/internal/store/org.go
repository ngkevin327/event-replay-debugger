package store

import (
	"context"

	"github.com/replay/platform/apps/api-gateway/internal/db"
)

// Store provides Postgres-backed metadata access.
type Store struct {
	pool *db.Pool
}

// NewStore wraps a connection pool.
func NewStore(pool *db.Pool) *Store {
	return &Store{pool: pool}
}

// CreateOrganization inserts a new org row.
func (s *Store) CreateOrganization(ctx context.Context, name, planTier string) (Organization, error) {
	var o Organization
	err := s.pool.QueryRow(ctx,
		`INSERT INTO organizations (name, plan_tier) VALUES ($1, $2)
		 RETURNING id, name, plan_tier, created_at`,
		name, planTier,
	).Scan(&o.ID, &o.Name, &o.PlanTier, &o.CreatedAt)
	return o, err
}

// GetOrganization loads an org by id.
func (s *Store) GetOrganization(ctx context.Context, id string) (Organization, error) {
	var o Organization
	err := s.pool.QueryRow(ctx,
		`SELECT id, name, plan_tier, created_at FROM organizations WHERE id = $1`, id,
	).Scan(&o.ID, &o.Name, &o.PlanTier, &o.CreatedAt)
	return o, err
}

// UpdateOrganization updates org metadata.
func (s *Store) UpdateOrganization(ctx context.Context, id, name, planTier string) (Organization, error) {
	var o Organization
	err := s.pool.QueryRow(ctx,
		`UPDATE organizations SET name = COALESCE(NULLIF($2, ''), name),
		                          plan_tier = COALESCE(NULLIF($3, ''), plan_tier)
		 WHERE id = $1
		 RETURNING id, name, plan_tier, created_at`,
		id, name, planTier,
	).Scan(&o.ID, &o.Name, &o.PlanTier, &o.CreatedAt)
	return o, err
}
