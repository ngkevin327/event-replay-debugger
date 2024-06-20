-- Local seed data for development (run after Postgres is up)
-- Example: psql "$POSTGRES_URL" -f infra/docker-compose/scripts/seed.sql

INSERT INTO organizations (id, name, plan_tier, created_at)
VALUES (
  '11111111-1111-1111-1111-111111111111',
  'Local Dev Org',
  'starter',
  NOW()
) ON CONFLICT DO NOTHING;

-- Placeholder: full schema lands in Stage 1 migrations
