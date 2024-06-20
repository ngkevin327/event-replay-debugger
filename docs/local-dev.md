# Local development

## Prerequisites

- Docker Desktop or compatible engine with Compose v2
- Go 1.22+
- Node 20+ and pnpm 9+ (web workspace)
- Make

## Start dependencies

From repository root:

```bash
make dev-deps
```

This starts:

| Service | Port | Purpose |
|---------|------|---------|
| Postgres | 5432 | Control plane metadata |
| MinIO | 9000 (API), 9001 (console) | Payload object storage |
| Redpanda | 9092 | Kafka-compatible broker |
| ClickHouse | 8123 (HTTP), 9009 (native) | Event index |

Copy environment template:

```bash
cp infra/docker-compose/.env.example infra/docker-compose/.env
```

## Verify health

```bash
docker compose -f infra/docker-compose/docker-compose.yml ps
```

ClickHouse tables:

```bash
curl -s 'http://localhost:8123/?query=SHOW+TABLES+FROM+replay'
```

## Seed data

After Stage 1 migrations exist, apply seed:

```bash
psql "postgres://replay:replay@localhost:5432/replay?sslmode=disable" \
  -f infra/docker-compose/scripts/seed.sql
```

Until migrations land, seed file documents intended test org/project IDs only.

## Run services

Services start individually once implemented:

```bash
# API gateway (Stage 1)
# go run ./apps/api-gateway/cmd/server

# Ingestion (Stage 2)
# go run ./apps/ingestion/cmd/server
```

## Troubleshooting

- **Port 5432 in use** — stop local Postgres or change published port in compose override file.
- **ClickHouse on Apple Silicon** — if official image fails, pin an arm64-compatible tag in a local compose override (not committed).
- **Redpanda not ready** — wait for health; first start may take 30–60s on cold pull.

## Stop stack

```bash
docker compose -f infra/docker-compose/docker-compose.yml down
```

Add `-v` to drop volumes when resetting local data.
