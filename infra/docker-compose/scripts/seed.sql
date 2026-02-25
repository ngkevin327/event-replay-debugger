-- Local seed data (run after migrations)
-- Usage: psql "postgres://replay:replay@localhost:5432/replay?sslmode=disable" -f infra/docker-compose/scripts/seed.sql

INSERT INTO organizations (id, name, plan_tier, created_at)
VALUES (
  '11111111-1111-1111-1111-111111111111',
  'Local Dev Org',
  'starter',
  NOW()
) ON CONFLICT (id) DO NOTHING;

INSERT INTO projects (id, org_id, name, created_at)
VALUES (
  '22222222-2222-2222-2222-222222222222',
  '11111111-1111-1111-1111-111111111111',
  'payments-sandbox',
  NOW()
) ON CONFLICT (id) DO NOTHING;
