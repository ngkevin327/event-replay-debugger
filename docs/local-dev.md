# Local development runbook

Target: a new developer can run and verify the control plane + web UI in **under 30 minutes**.

## Prerequisites

| Tool | Version |
|------|---------|
| Docker Desktop (Compose v2) | latest |
| Go | 1.22+ (module uses 1.25) |
| Node.js | 20+ |
| pnpm | 9+ |
| PowerShell or bash | for scripts |

## Quick start (< 30 min)

### 1. First-time setup

**Windows:**

```powershell
./scripts/setup-local.ps1
```

**macOS / Linux:**

```bash
chmod +x scripts/setup-local.sh scripts/verify-local.sh
./scripts/setup-local.sh
```

This will:

- Copy `.env.local.example` → `.env.local` (if missing)
- Start Docker services (Postgres, MinIO, Redpanda, ClickHouse)
- Apply API migrations
- Install web dependencies

### 2. Load environment

Copy and edit `.env.local` from `.env.local.example`, then export variables.

**PowerShell:**

```powershell
Get-Content .env.local | ForEach-Object {
  if ($_ -match '^([^#][^=]+)=(.*)$') {
    Set-Item -Path "env:$($matches[1].Trim())" -Value $matches[2].Trim()
  }
}
```

### 3. Start services (three terminals)

**Terminal A — API gateway**

```powershell
cd apps/api-gateway
go run ./cmd/server
```

**Terminal B — Web UI**

```powershell
pnpm --filter @replay/web dev
```

Open http://localhost:5173

**Terminal C — Ingestion (optional)**

```powershell
cd apps/ingestion
go run ./cmd/server
```

### 4. Verify working local product

```powershell
./scripts/verify-local.ps1
```

Expected: all **required** checks pass (`api-health`, `auth-register`, `create-project`, `create-incident`, `timeline`, `graph`, `create-replay`). Ingestion is optional.

## Docker services and ports

Host ports avoid common Windows reserved ranges (Hyper-V / excluded port blocks):

| Service | Host port | Purpose |
|---------|-----------|---------|
| Postgres | **15433** | Control plane DB |
| MinIO API | **19000** | Payload storage |
| MinIO console | **19001** | Admin UI |
| Redpanda (Kafka) | **19092** | Message broker |
| ClickHouse HTTP | **8123** | Event index |

```bash
docker compose -f infra/docker-compose/docker-compose.yml ps
```

ClickHouse sanity check:

```bash
curl -s "http://localhost:8123/?query=SHOW+TABLES+FROM+replay"
```

## Environment variables

### Required (API gateway)

| Variable | Example | Purpose |
|----------|---------|---------|
| `DATABASE_URL` | `postgres://replay:replay@localhost:15433/replay?sslmode=disable` | Postgres connection |
| `JWT_SECRET` | `local-dev-jwt-secret-change-me` | Session JWT signing |

### Required for local MVP demo (timeline/graph + ready incidents)

| Variable | Example | Purpose |
|----------|---------|---------|
| `LOCAL_DEMO_MODE` | `1` | Serve `fixtures/local` timeline/graph |
| `LOCAL_DEMO_AUTO_READY` | `1` | Mark new incidents `ready` immediately |
| `LOCAL_DEMO_FIXTURES_DIR` | `../../fixtures/local` | Path from `apps/api-gateway` cwd |

### Optional (ingestion service)

| Variable | Default | Purpose |
|----------|---------|---------|
| `HTTP_ADDR` | `:8081` | Ingestion listen address |
| `S3_ENDPOINT` | `http://localhost:19000` | MinIO |
| `S3_ACCESS_KEY` / `S3_SECRET_KEY` | `minio` / `minioadmin` | MinIO credentials |
| `S3_BUCKET` | `replay-payloads` | Payload bucket |
| `CLICKHOUSE_URL` | `http://localhost:8123` | Event index |
| `DEFAULT_ORG_ID` | — | Org for ingest quota |

### Optional (web / e2e)

| Variable | Default | Purpose |
|----------|---------|---------|
| `API_BASE_URL` | `http://localhost:8080` | verify script override |
| `WEB_BASE_URL` | `http://localhost:5173` | verify script override |

## Manual UI verification checklist

- [ ] Open http://localhost:5173 — login/register page loads
- [ ] Register a new org/user — redirects to dashboard
- [ ] Create a project (Settings or dashboard flow)
- [ ] Create incident — status shows **ready** (with `LOCAL_DEMO_AUTO_READY=1`)
- [ ] Open incident detail — timeline shows events, graph renders
- [ ] Start replay — replay row appears (may stay `pending` without orchestrator)

## Makefile shortcuts

```bash
make dev-deps        # docker compose up
make migrate-local   # apply API migrations
make setup-local     # deps + migrate + pnpm install
```

## Seed data (optional)

```bash
psql "postgres://replay:replay@localhost:15433/replay?sslmode=disable" \
  -f infra/docker-compose/scripts/seed.sql
```

Provides org `11111111-1111-1111-1111-111111111111` and project `22222222-2222-2222-2222-222222222222`.

## Troubleshooting

| Symptom | Fix |
|---------|-----|
| `ports are not available` on Windows | Use documented high ports; run `netsh interface ipv4 show excludedportrange protocol=tcp` |
| Migration dirty version | `docker compose ... down -v` and re-run setup |
| Timeline/graph 409 not_ready | Set `LOCAL_DEMO_MODE=1` and correct `LOCAL_DEMO_FIXTURES_DIR` |
| Register works but UI login fails | Ensure API returns `access_token` + `user` (v0.1.0-mvp+ auth shape) |
| ClickHouse fails on ARM Mac | Use compose override with arm64-compatible image (local only) |

## Stop stack

```bash
docker compose -f infra/docker-compose/docker-compose.yml down
```

Add `-v` to reset volumes.

## Related docs

- [MVP gap report](mvp-gap-report.md) — acceptance matrix and known gaps
- [Release v0.1.0-mvp](release/v0.1.0-mvp.md) — shipping checklist
- [Quickstart (product)](site/docs/quickstart.md) — agent install (Helm/K8s)
