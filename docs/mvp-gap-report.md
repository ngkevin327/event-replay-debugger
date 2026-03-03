# MVP gap report

**Generated:** 2026-05-20 (local validation run)  
**Release target:** `v0.1.0-mvp` ([docs/release/v0.1.0-mvp.md](release/v0.1.0-mvp.md))

## Local run status (this session)

| Component | Status | Evidence |
|-----------|--------|----------|
| Docker deps (Postgres, ClickHouse, MinIO, Redpanda) | Running | `docker compose -f infra/docker-compose/docker-compose.yml ps` |
| API gateway `:8080` | Running | `GET /health` → `{"status":"ok"}` |
| Ingestion `:8081` | Not started (optional) | verify script: optional fail |
| Web UI `:5173` | Reachable | verify script: HTTP 200 |
| Migrations | Applied (12) | `go run ./cmd/migrate` exit 0 |
| Verify script | Core flows pass | `scripts/verify-local.ps1` exit 0 |

### Verified API flow (2026-05-20)

```
[PASS] api-health
[PASS] auth-register
[PASS] create-project
[PASS] create-incident (status=ready with LOCAL_DEMO_AUTO_READY=1)
[PASS] timeline (3 events from fixtures/local)
[PASS] graph (3 nodes)
[PASS] create-replay (status=pending)
```

## MVP acceptance matrix (PRD §19)

| Criterion | Status | Evidence | Missing / gap | Next action |
|-----------|--------|----------|---------------|-------------|
| Capture agent via Helm + heartbeats | Partial | Helm chart exists; no local agent run in verify | Agent not in local compose; needs K8s or demo worker | Add `test/demo` compose to setup script optional profile |
| Create incident (window + topics) | **Pass** | verify `create-incident` | — | — |
| Timeline + graph in web UI | **Pass** (demo mode) | fixture loaders + UI 200 | Production path uses stub loaders without `LOCAL_DEMO_MODE` | Wire S3/DB artifact loaders from reconstruction output |
| Deterministic replay + divergence | Partial | replay row created `pending` | Orchestrator is health-only; worker not driven from API | Connect replay-orchestrator to store + replay-worker |
| E2E Playwright (auth, incidents, replay) | Partial | specs exist; replay uses mocks | No full-stack e2e without running API | Add `test/e2e` job to verify-local optional step |
| Staging deploy on `main` | Not tested locally | CI workflows in repo | Requires AWS/EKS credentials | Out of scope for local validation |
| Security (rate limit, tenant isolation, ZAP) | Partial | unit tests pass | ZAP/staging not run locally | Document staging-only security verification |
| Design partner dry runs (×3) | Fail | release checklist unchecked | Process, not code | Partner onboarding |

## Architecture gaps (code-level)

| Area | Status | Notes |
|------|--------|-------|
| API auth register/login JSON | **Fixed** | Now returns `{ access_token, user: { id, email, org_id } }` for web client |
| Timeline/graph loaders | Partial | `LOCAL_DEMO_MODE=1` serves `fixtures/local`; default stubs return not_ready |
| Reconstruction worker | Fail (runtime) | In-memory queue, empty handlers in `cmd/worker` |
| Replay orchestrator | Fail (runtime) | Only `/healthz` server |
| Notifications worker | Fail (runtime) | `pollEvents` is no-op memory bus |
| Ingestion pipeline | Partial | Code complete; not required for control-plane UI demo |
| Replay execution | Partial | DB records only; no worker progression to `succeeded`/`diverged` |

## Recommended priority order

1. **P0 — Local full stack:** document ports (15433, 19000, 19092, 8123); keep verify script in CI.
2. **P0 — Artifact pipeline:** reconstruction writes timeline/graph → API loaders read from storage (remove stub default in prod).
3. **P1 — Replay loop:** orchestrator provisions worker, updates `replay_runs.status`, divergence report URI.
4. **P1 — Ingestion in local compose:** optional profile with MinIO bucket bootstrap + sample ingest.
5. **P2 — Design partner checklist:** operational, not blocking local dev.

## Commands to reproduce this validation

```powershell
# Setup
./scripts/setup-local.ps1

# Start API (from repo root, after loading .env.local)
cd apps/api-gateway
$env:DATABASE_URL="postgres://replay:replay@localhost:15433/replay?sslmode=disable"
$env:JWT_SECRET="local-dev-jwt-secret"
$env:LOCAL_DEMO_MODE="1"
$env:LOCAL_DEMO_AUTO_READY="1"
$env:LOCAL_DEMO_FIXTURES_DIR="../../fixtures/local"
go run ./cmd/server

# Web (separate terminal, repo root)
pnpm --filter @replay/web dev

# Verify
./scripts/verify-local.ps1
```
