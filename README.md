# Replay

Production incident replay for async backends. Replay captures Kafka consumer/producer behavior, reconstructs event timelines, and runs deterministic replays so teams can see how failures unfolded—not just that something failed.

## Target users

Backend engineers on teams running Kafka (or similar) with event-driven workflows, retries, and microservices. Typical customers are small-to-mid SaaS and fintech backends (roughly 5–50 engineers) debugging production incidents.

## Architecture overview

| Component | Role |
|-----------|------|
| **Capture agent** (`packages/agent-go`) | SDK hooks on Kafka clients; ships batches to ingest |
| **Ingestion** (`apps/ingestion`) | Validates and stores payloads (S3) + index (ClickHouse) |
| **API gateway** (`apps/api-gateway`) | Auth, orgs/projects, incidents, replay control plane |
| **Reconstruction** (`apps/reconstruction`) | Builds timelines and workflow graphs from stored events |
| **Replay orchestrator** (`apps/replay-orchestrator`) | Manages replay run lifecycle |
| **Replay worker** (`services/replay-worker`) | Feeds sandbox Kafka from timeline artifacts |
| **Web** (`apps/web`) | Incident timeline, graph, replay UI |

```
Customer Kafka ──► Agent ──► Ingestion ──► S3 + ClickHouse
                              ▲
Control plane (API) ◄── Incidents / Replay orchestration
                              │
                         Web UI
```

## Requirements

- Go 1.22+
- Node.js 20+ and pnpm 9+
- Docker and Docker Compose (local dependencies)
- Make

## Local development

```bash
# Start Postgres, MinIO, Redpanda, ClickHouse
make dev-deps

# Run lint and tests
make lint
make test
```

See [docs/local-dev.md](docs/local-dev.md) for ports, seed data, and troubleshooting.

## Build and test

```bash
make build    # build Go services
make test     # unit tests
make lint     # Go + web lint
```

## Deployment overview

- **Control plane:** containerized services on Kubernetes (EKS), configured via Helm under `deploy/helm/`.
- **Capture agent:** customer-cluster Helm chart (`deploy/helm/replay-agent`) or Go SDK in application images.
- **Infrastructure:** Terraform modules under `infra/terraform/` (staging/prod).

## Repository layout

- `apps/` — Go services and React web app
- `packages/` — shared schema, Go libraries, capture agent
- `services/` — ephemeral replay worker
- `infra/` — docker-compose for local dev, Terraform for cloud
- `api/` — OpenAPI specification
- `docs/adr/` — architecture decision records

## License

Proprietary. All rights reserved.
