# Replay Platform

Replay captures Kafka consumer and producer traffic, reconstructs incident timelines, and runs deterministic replays so teams can debug async payment and ledger pipelines with confidence.

## What you can do

- **Capture** events with the open-source agent (sidecar or DaemonSet)
- **Create incidents** for time-bounded windows and topics
- **Inspect** virtualized timelines and workflow graphs
- **Replay** with strict or compressed timing and divergence reports

## Architecture overview

| Component | Role |
|-----------|------|
| API Gateway | Auth, projects, incidents, replay control |
| Ingestion | Batch ingest to object storage + ClickHouse |
| Reconstruction | Build timeline and graph artifacts |
| Replay worker | Deterministic replay and divergence detection |
| Web UI | Operator console |

## Next steps

- [Quickstart](./quickstart.md) — install the agent and run your first replay
- [Agent configuration](./agent-config.md) — Helm values reference
- [API reference](./api-reference.md) — control plane REST API
