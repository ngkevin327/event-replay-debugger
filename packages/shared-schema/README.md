# shared-schema

Canonical JSON Schema for events captured from customer Kafka workloads.

## Consumers

- Capture agent (`packages/agent-go`) — serializes outbound batches
- Ingestion service — validates inbound batches before persistence
- Reconstruction — reads indexed metadata aligned with this shape

## Versioning

- Schema file is versioned implicitly via `$id` path (`v1`)
- Breaking field changes require a new schema file and coordinated rollout across agent + ingestion

## Validation

Run tests:

```bash
go test ./packages/shared-schema/...
```
