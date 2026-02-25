#!/usr/bin/env bash
# Deterministic first-run local setup (bash)
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

echo "==> Checking prerequisites..."
for cmd in docker go pnpm; do
  command -v "$cmd" >/dev/null || { echo "Missing: $cmd"; exit 1; }
done

if [[ ! -f .env.local ]]; then
  cp .env.local.example .env.local
  echo "Created .env.local from .env.local.example"
fi

echo "==> Starting Docker dependencies..."
docker compose -f infra/docker-compose/docker-compose.yml up -d

echo "==> Waiting for Postgres..."
for i in $(seq 1 30); do
  if docker compose -f infra/docker-compose/docker-compose.yml exec -T postgres pg_isready -U replay >/dev/null 2>&1; then
    break
  fi
  sleep 2
  if [[ $i -eq 30 ]]; then
    echo "Postgres not ready"
    exit 1
  fi
done

echo "==> Applying API migrations..."
export DATABASE_URL="${DATABASE_URL:-postgres://replay:replay@localhost:5432/replay?sslmode=disable}"
(cd apps/api-gateway && go run ./cmd/migrate -dir ./migrations)

echo "==> Installing web dependencies..."
pnpm install --filter @replay/web

echo "==> Setup complete."
