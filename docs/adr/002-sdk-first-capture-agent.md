# ADR-002: SDK-first capture agent

## Status

Accepted (2026-05-19)

## Context

Customers run Kafka consumers and producers in their own clusters. We need capture fidelity (headers, outcomes, retry signals) without requiring broker-side plugins or mirror-maker deployments for MVP.

## Decision

Ship a **Go SDK** that wraps official Kafka client handlers (Sarama) with pre/post capture hooks. Helm chart provides configuration injection; **sidecar** capture is optional and limited in v1.

## Consequences

### Positive

- Highest fidelity for consumer latency, errors, and retry paths
- Minimal intrusion: no broker config changes required
- Natural fit for Go-heavy ICP backends

### Negative

- Requires code change in customer services (handler wrapper)
- Java/Python agents deferred post-MVP

## Alternatives considered

- **Sidecar-only** — lower integration friction, weaker outcome/retry visibility
- **Broker interceptors** — operationally heavy, not portable across hosted Kafka vendors
