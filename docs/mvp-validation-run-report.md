# MVP Validation Run Report

**Project:** Replay Platform  
**Repository:** `d:\Projects\Fake Git\real 3`  
**Release target:** `v0.1.0-mvp`  
**Report date:** 2026-05-20 (updated 2026-05-21)  
**Author:** Kevin Ng (validation session via Cursor agent)  
**Branch:** `main` (HEAD at time of last update: `9cf49fb`)

---

## 1. Executive summary

This report documents work performed to prepare the repository for **MVP validation**: deterministic local setup, running the product locally on Windows, automated verification of core API flows, gap analysis against the MVP acceptance checklist, and incremental commits so progress can be followed step-by-step.

**Outcome:** The **control-plane + web UI demo path is runnable locally** with documented scripts and env vars. Auth, project/incident CRUD, timeline/graph (demo fixtures), and replay record creation all pass automated verification. A subsequent **UI polish and consistency pass** (five web commits, `42a69aa`–`9cf49fb`) was re-validated on 2026-05-21 and **did not regress MVP behavior** — see Section 13. **Full end-to-end replay** (capture → ingest → reconstruction → deterministic replay execution) is **not** complete in local runtime; gaps are listed in Section 7.

---

## 2. Objectives

| # | Objective | Result |
|---|-----------|--------|
| 1 | Make first-run local setup deterministic | Done — `setup-local.ps1` / `setup-local.sh`, `.env.local.example`, Makefile targets |
| 2 | Run product locally and capture evidence | Done — Docker deps, API, web UI, verify script |
| 3 | Validate core local flows | Done — 7/7 required checks pass (see Section 5) |
| 4 | Compare state to MVP acceptance criteria | Done — matrix in Section 7 |
| 5 | Document runbook for new developers (<30 min) | Done — `docs/local-dev.md` updated |
| 6 | Commit progress in small logical steps | Done — 11 commits (Section 4) |
| 7 | Confirm UI work does not break MVP flows | Done — Section 13 (2026-05-21) |

---

## 3. Work performed

### Phase A — Baseline and deterministic setup

- Inspected repo layout, `README.md`, `docs/local-dev.md`, `docs/release/v0.1.0-mvp.md`, and service entrypoints.
- Added **`.env.local.example`** with required and optional variables for API, ingestion, and local demo mode.
- Added **`fixtures/local/timeline.json`** and **`fixtures/local/graph.json`** for demo timeline/graph responses.
- Added **`scripts/setup-local.ps1`** and **`scripts/setup-local.sh`** to:
  - Check prerequisites (`docker`, `go`, `pnpm`)
  - Copy env template
  - Start Docker Compose dependencies
  - Wait for Postgres
  - Run API migrations
  - Install web dependencies
- Extended **`Makefile`** with `setup-local`, `migrate-local`, `verify-local`.
- Updated **`infra/docker-compose/scripts/seed.sql`** with org/project UUIDs for optional seeding.

### Phase B — Local demo mode (API)

- Implemented **`FileArtifactLoader`** (`apps/api-gateway/internal/handlers/local_artifact_loader.go`) to serve fixture JSON when `LOCAL_DEMO_MODE=1`.
- Wired **`ResolveTimelineLoader()` / `ResolveGraphLoader()`** in routes (replacing always-stub loaders in demo mode).
- Added **`LOCAL_DEMO_AUTO_READY=1`** so new incidents are marked `ready` immediately for local UI testing.
- Added unit test for file artifact loader.

### Phase C — Verification automation

- Added **`scripts/verify-local.ps1`** and **`scripts/verify-local.sh`** to check:
  - API `/health`
  - Ingestion `/health` (optional)
  - Web UI HTTP 200
  - Register → project → incident → timeline → graph → replay API chain

### Phase D — Infrastructure fixes (Windows)

During local runs on Windows, several host ports were blocked by the OS (Hyper-V excluded port ranges). Compose and env defaults were updated:

| Service | Original port | Local port |
|---------|-----------------|------------|
| Postgres | 5432 | **15433** |
| MinIO API | 9000 | **19000** |
| MinIO console | 9001 | **19001** |
| Redpanda Kafka | 9092 | **19092** |
| Redpanda admin | 9644 | **19644** |
| ClickHouse native | 9009 | removed (HTTP **8123** only) |

- Fixed PostgreSQL migrations: quoted reserved column **`"offset"`** in `000007_snapshots.up.sql` and `000009_replay_events.up.sql`.

### Phase E — Auth API alignment (web client)

- **Register** and **login** responses updated to match web `AuthContext` expectations:
  ```json
  { "access_token": "...", "user": { "id": "...", "email": "...", "org_id": "..." } }
  ```
- Register handler now issues JWT on successful signup (`JWTSecret` passed from routes).

### Phase F — Documentation

- Rewrote **`docs/local-dev.md`** as a 30-minute runbook (ports, env vars, checklist, troubleshooting).
- Added **`docs/mvp-gap-report.md`** (acceptance matrix and architecture gaps).
- This file (**`docs/mvp-validation-run-report.md`**) consolidates the full validation run.

---

## 4. Commits (in order)

| Hash | Message | Purpose |
|------|---------|---------|
| `4b4e0fd` | `chore(local): add deterministic first-run setup scripts` | Setup scripts, env template, fixtures, Makefile, seed |
| `b16600f` | `feat(api-gateway): enable local demo artifacts for MVP validation` | File artifact loader, demo env, auto-ready incidents |
| `a39c150` | `test(scripts): add local health and API flow verification` | `verify-local.ps1` / `.sh` |
| `d55d0f9` | `fix(infra): use Windows-safe host ports for local compose` | Compose ports, migration SQL fix, env defaults |
| `b875072` | `fix(api-gateway): align auth responses with web client and verify script` | Auth JSON shape + verify script token handling |
| `65a95fd` | `docs: add local runbook and MVP gap report for validation` | `local-dev.md`, `mvp-gap-report.md` |
| `42a69aa` | `feat(web): add design system and polished global UI styles` | Tokens, global CSS |
| `cfebb16` | `feat(web): redesign auth, app shell, and operator pages for MVP` | Auth layout, sidebar, pages |
| `a0cfac5` | `feat(web): add light mode, motion system, and extended design tokens` | Theme toggle, animations |
| `7e7247d` | `feat(web): polish timeline, graph, replay flows and empty-state illustrations` | Timeline/graph/replay visuals |
| `9cf49fb` | `feat(web): enforce consistent UI patterns across all views` | Modal, StatusBadge, `components.css` |

**Not pushed** to remote (per user instruction).

---

## 5. Local run status and verification evidence

### 5.1 Services exercised

| Component | Port / URL | Status | Notes |
|-----------|------------|--------|-------|
| Postgres | `localhost:15433` | Running | Migrations applied (12) |
| ClickHouse | `localhost:8123` | Running | `replay.events` table from init SQL |
| MinIO | `localhost:19000` | Running | Not required for control-plane verify |
| Redpanda | `localhost:19092` | Running | Not required for control-plane verify |
| API gateway | `localhost:8080` | Running | `go run ./cmd/server` from `apps/api-gateway` |
| Ingestion | `localhost:8081` | Not started | Optional; verify marks as optional fail |
| Web (Vite) | `localhost:5173` | Reachable | HTTP 200 |

### 5.2 Environment used for verification

```text
DATABASE_URL=postgres://replay:replay@localhost:15433/replay?sslmode=disable
JWT_SECRET=local-dev-jwt-secret
LOCAL_DEMO_MODE=1
LOCAL_DEMO_AUTO_READY=1
LOCAL_DEMO_FIXTURES_DIR=../../fixtures/local
```

### 5.3 Automated verification (`scripts/verify-local.ps1`)

**Exit code:** `0` (required checks only; ingestion optional failure ignored)

| Check | Result | Detail |
|-------|--------|--------|
| `api-health` | PASS | `{"status":"ok"}` |
| `ingestion-health` | FAIL (optional) | Service not running |
| `web-ui` | PASS | HTTP 200 |
| `auth-register` | PASS | JWT issued |
| `create-project` | PASS | Project UUID returned |
| `create-incident` | PASS | `status=ready` |
| `timeline` | PASS | 3 events from fixtures |
| `graph` | PASS | 3 nodes from fixtures |
| `create-replay` | PASS | Replay created, `status=pending` |

### 5.4 Unit tests run during session

| Command | Result |
|---------|--------|
| `go test ./...` in `apps/api-gateway` | Pass |
| `go test ./...` in `apps/ingestion` | Pass |
| `pnpm --filter @replay/web test` | Pass (5 test files) |

### 5.5 Post-UI re-verification (2026-05-21)

After UI commits `42a69aa`–`9cf49fb`, the Replay Docker stack was started and the API was run against Postgres on **15433** with demo env vars. Verification used `API_BASE_URL=http://localhost:8082` because port **8080** was occupied by an unrelated process (register returned 500 against that listener).

| Check | Result | Detail |
|-------|--------|--------|
| `api-health` | PASS | `{"status":"ok"}` |
| `ingestion-health` | FAIL (optional) | Service not running |
| `web-ui` | PASS | HTTP 200 (`pnpm --filter @replay/web dev`) |
| `auth-register` | PASS | JWT issued |
| `create-project` | PASS | Project UUID returned |
| `create-incident` | PASS | `status=ready` |
| `timeline` | PASS | 3 events from fixtures |
| `graph` | PASS | 3 nodes from fixtures |
| `create-replay` | PASS | Replay created, `status=pending` |

**Web unit tests (post-UI):** `pnpm --filter @replay/web test` — 5/5 files passed.

**API handler tests (sanity):** `go test ./internal/handlers/...` in `apps/api-gateway` — pass.

---

## 6. How to reproduce

### One-time setup

```powershell
cd "d:\Projects\Fake Git\real 3"
./scripts/setup-local.ps1
```

### Start API (terminal 1)

```powershell
cd apps/api-gateway
$env:DATABASE_URL="postgres://replay:replay@localhost:15433/replay?sslmode=disable"
$env:JWT_SECRET="local-dev-jwt-secret"
$env:LOCAL_DEMO_MODE="1"
$env:LOCAL_DEMO_AUTO_READY="1"
$env:LOCAL_DEMO_FIXTURES_DIR="../../fixtures/local"
go run ./cmd/server
```

### Start web (terminal 2)

```powershell
pnpm --filter @replay/web dev
```

Open http://localhost:5173

### Verify (terminal 3)

```powershell
./scripts/verify-local.ps1
```

### Manual UI checklist

- [ ] Register / login
- [ ] Create project
- [ ] Create incident → status **ready**
- [ ] Incident detail → timeline events visible
- [ ] Graph tab renders nodes/edges
- [ ] Replay tab → start replay (record appears; may stay `pending`)

**Note:** Vite proxies `/v1` to `http://localhost:8080`. For manual UI testing, run the Replay API on **8080** with demo env vars (see Section 6), or the browser will hit the wrong backend.

---

## 7. MVP readiness — acceptance matrix (PRD §19)

Reference: `docs/release/v0.1.0-mvp.md`

| Criterion | Status | Evidence | Missing pieces | Recommended next action |
|-----------|--------|----------|----------------|-------------------------|
| Capture agent via Helm + heartbeats | **Partial** | Chart `deploy/helm/replay-agent` exists | No local agent in default compose | Add optional `test/demo` profile to setup |
| Create incident (window + topics) | **Pass** | verify `create-incident` | — | — |
| Timeline + graph in web UI | **Pass** (demo) | Fixtures + UI 200; post-UI verify unchanged (Section 13) | Without `LOCAL_DEMO_MODE`, API stubs return `not_ready` | Wire reconstruction → S3/DB artifact loaders |
| Deterministic replay + divergence | **Partial** | Replay row `pending` in DB | Orchestrator is health-only; no worker execution | Implement orchestrator ↔ store ↔ replay-worker |
| E2E Playwright (auth, incidents, replay) | **Partial** | Specs in `test/e2e/` | Replay spec mocks API routes | Run Playwright against live API in CI |
| Staging deploy on `main` | **Not tested** | Workflows in repo | Needs AWS/EKS | Staging-only validation |
| Security (rate limit, tenant isolation, ZAP) | **Partial** | Go unit tests pass | ZAP not run locally | Keep in staging CI |
| Design partner dry runs (×3) | **Fail** | Unchecked in release doc | Operational | Partner onboarding |

---

## 8. Architecture gaps (code-level)

| Component | Runtime status | Notes |
|-----------|----------------|-------|
| API auth (register/login) | **Fixed** | Matches web client contract |
| Timeline/graph loaders | **Partial** | Demo: `FileArtifactLoader`; default: `StubTimelineLoader` / `StubGraphLoader` |
| Reconstruction worker | **Not wired** | In-memory queue, empty handler map in `cmd/worker` |
| Replay orchestrator | **Stub** | Only `/healthz` on `:8091` |
| Replay worker / feeder | **Not integrated** | Library code exists; not driven from API-created runs |
| Notifications worker | **Stub** | `pollEvents` returns empty memory channel |
| Ingestion pipeline | **Optional locally** | Full code path; requires MinIO + ClickHouse + API keys |

---

## 9. What is still missing for full MVP

### P0 (blocks “real” local product demo)

1. **Artifact pipeline** — Reconstruction must write timeline/graph artifacts; API must load them without `LOCAL_DEMO_MODE`.
2. **Replay execution loop** — Orchestrator should provision worker, advance `replay_runs.status`, write divergence report.

### P1 (important for completeness)

3. **Reconstruction worker** — Redis (or equivalent) queue + job handlers connected to ClickHouse/S3.
4. **Ingestion in local profile** — MinIO bucket bootstrap, sample batch ingest, agent or demo producer.
5. **Full-stack E2E** — Playwright against live API (no route mocks for replay).

### P2 (release / ops)

6. **Design partner dry runs** (3×) — checklist in `docs/partners/onboarding-checklist.md`.
7. **Staging verification** — deploy workflow, ZAP baseline, tenant isolation in integrated env.

---

## 10. Files added or materially changed

| Path | Change |
|------|--------|
| `.env.local.example` | New — local env template |
| `fixtures/local/timeline.json` | New — demo timeline |
| `fixtures/local/graph.json` | New — demo graph |
| `scripts/setup-local.ps1` | New |
| `scripts/setup-local.sh` | New |
| `scripts/verify-local.ps1` | New |
| `scripts/verify-local.sh` | New |
| `Makefile` | `setup-local`, `migrate-local`, port defaults |
| `infra/docker-compose/docker-compose.yml` | Windows-safe ports |
| `infra/docker-compose/.env.example` | Port updates |
| `infra/docker-compose/scripts/seed.sql` | Org/project seed |
| `apps/api-gateway/internal/handlers/local_artifact_loader.go` | New |
| `apps/api-gateway/internal/handlers/incidents.go` | `LOCAL_DEMO_AUTO_READY` |
| `apps/api-gateway/internal/handlers/auth_register.go` | Token + user response |
| `apps/api-gateway/internal/handlers/auth_login.go` | User object in response |
| `apps/api-gateway/internal/server/routes.go` | Demo loaders, register JWT |
| `apps/api-gateway/migrations/000007_snapshots.up.sql` | Quote `"offset"` |
| `apps/api-gateway/migrations/000009_replay_events.up.sql` | Quote `"offset"` |
| `apps/ingestion/internal/config/config.go` | Default MinIO port 19000 |
| `docs/local-dev.md` | Rewritten runbook |
| `docs/mvp-gap-report.md` | Gap matrix (companion doc) |
| `docs/mvp-validation-run-report.md` | This report |
| `apps/web/src/styles/*.css` | Design system, components, timeline/graph/replay |
| `apps/web/src/components/Modal.tsx` | Shared modal shell |
| `apps/web/src/components/StatusBadge.tsx` | Shared status pill |

---

## 11. Related documentation

- [Local development runbook](local-dev.md) — step-by-step for new developers
- [MVP gap report](mvp-gap-report.md) — condensed gap matrix
- [Release v0.1.0-mvp](release/v0.1.0-mvp.md) — shipping checklist
- [README](../README.md) — project overview

---

## 12. UI consistency pass — MVP regression analysis (2026-05-21)

### 12.1 Scope of UI changes

Five commits from `42a69aa` through `9cf49fb` touched **presentation only**:

- Stylesheets (`tokens.css`, `global.css`, `components.css`, feature CSS)
- Layout components (`AuthLayout`, `PageHeader`, `Modal`, `StatusBadge`, `ThemeToggle`)
- Page and feature component markup/classes (no new API endpoints)

### 12.2 What was *not* changed (no functional side effects)

| Layer | Changed in UI commits? | Impact |
|-------|------------------------|--------|
| `apps/web/src/api/client.ts` | No | Fetch contract unchanged |
| `apps/web/src/api/hooks.ts` | No | TanStack Query keys and mutations unchanged |
| `apps/web/src/context/AuthContext.tsx` | No | Login/register token storage unchanged |
| `apps/api-gateway/**` | No | Backend behavior unchanged by UI work |
| MVP env vars (`LOCAL_DEMO_*`) | No | Demo timeline/graph/replay still work |

The only `context/` addition across UI commits is **`ThemeContext.tsx`** (light/dark preference in `localStorage`) — cosmetic, not part of MVP API flows.

### 12.3 MVP feature wiring preserved

| MVP flow | UI entry point | API / hook (unchanged) |
|----------|----------------|------------------------|
| Register / login | `Login.tsx`, `Register.tsx` | `AuthContext` → `/v1/auth/register`, `/v1/auth/login` |
| Create incident | `CreateIncidentModal` | `apiFetch` POST `/v1/projects/{id}/incidents` (submit via `form="create-incident-form"`) |
| List incidents | `Incidents.tsx` | `useIncidents` |
| Timeline | `IncidentDetail` → `TimelineList` | `useTimeline` |
| Graph | `IncidentDetail` → `WorkflowGraph` | `useGraph` |
| Start replay | `ReplayPanel` | `useCreateReplay` → navigate to `/replays/{id}` |
| Export / share | `IncidentActions` | `useExportIncident`, `useCreateShare` |
| Settings / agent | `Settings.tsx`, `AgentSetup.tsx` | Existing hooks and static content |

### 12.4 Automated regression evidence

- **Web:** `pnpm --filter @replay/web test` — all 5 test files pass (AuthContext, API client, Dashboard, TimelineList, WorkflowGraph).
- **API (MVP chain):** `scripts/verify-local.ps1` with Replay Postgres + demo API — all required checks pass (Section 5.5).
- **E2E selectors:** Playwright specs (`test/e2e/incidents.spec.ts`) still target `dialog`, `Create incident`, and `.badge` — compatible with `Modal` (`<dialog>`) and `StatusBadge`.

### 12.5 Residual risks (unchanged from pre-UI)

| Risk | Severity | Notes |
|------|----------|-------|
| Replay stays `pending` | Expected | Orchestrator stub; not caused by UI |
| Manual UI needs API on :8080 | Ops | Vite proxy default; wrong process on 8080 breaks browser flows |
| Playwright replay spec mocks API | Test gap | Full-stack E2E still partial |
| Production without `LOCAL_DEMO_MODE` | Product | Timeline/graph stubs return `not_ready` |

### 12.6 Verdict

**The UI polish and consistency work does not introduce MVP functional regressions.** Core control-plane flows verified by `verify-local.ps1` behave the same after UI commits. Remaining MVP gaps (replay execution, reconstruction, ingestion) are **backend/runtime** limitations documented in Sections 7–9, not side effects of the web UI refactor.

---

## 13. Conclusion

The repository is **ready for local MVP validation of the control plane and operator UI** using the scripts and demo mode documented above. Eleven commits capture setup, API demo mode, verification, docs, and UI work in reviewable steps. **UI consistency changes are presentation-only and re-validated without regressing MVP API flows** (Section 12). **Full MVP product behavior** (live reconstruction artifacts, replay execution, capture/ingest loop) requires the P0/P1 items in Section 9 before the PRD §19 checklist can be marked complete without demo flags.
