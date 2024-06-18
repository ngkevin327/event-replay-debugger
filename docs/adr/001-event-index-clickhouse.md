# ADR-001: ClickHouse for event index

## Status

Accepted (2026-05-19)

## Context

Captured Kafka events produce high-volume metadata (topic, partition, offset, timestamps, outcomes) and large payloads stored separately in S3. Incident reconstruction queries need fast time-range scans per `project_id` without loading full payloads into Postgres.

## Decision

Use **ClickHouse** for the `events` index table. Store payload bytes in **S3** (content-addressed by SHA-256). Postgres remains the system of record for orgs, projects, incidents, and replay runs.

## Consequences

### Positive

- Efficient time-range and topic filters for reconstruction jobs
- Keeps RDS costs predictable for relational metadata
- Aligns with columnar analytics patterns for large incident windows

### Negative

- Additional operational component in local and staging stacks
- Cross-store consistency (CH row + S3 object) requires careful ingest ack semantics

## Alternatives considered

- **Postgres only** — simpler ops, poor scan performance at MVP target scale
- **Elasticsearch** — stronger full-text search, heavier ops burden for MVP team size
