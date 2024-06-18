# ADR-003: Shared sandbox Kafka for replay (MVP)

## Status

Accepted (2026-05-19)

## Context

Replay runs need an isolated Kafka target to republish timeline events without touching production topics. Dedicated MSK cluster per replay run is safest but costly and slow to provision for Starter tier.

## Decision

MVP uses a **shared internal Kafka cluster** (or Redpanda in local dev) with topic namespaced per replay run (`replay.{run_id}.*`). ACLs restrict producers/consumers to replay prefixes. Enterprise may add dedicated cells later.

## Consequences

### Positive

- Faster replay startup (no per-run cluster provisioning in MVP)
- Lower infrastructure cost for design partner phase
- Sufficient isolation when combined with strict ACLs and ephemeral workers

### Negative

- Noisy-neighbor risk if replay load spikes concurrently
- Requires careful topic cleanup TTL jobs

## Alternatives considered

- **Dedicated MSK per replay** — strongest isolation, MVP too slow/expensive
- **In-memory bus only** — fast tests, diverges from production Kafka semantics
