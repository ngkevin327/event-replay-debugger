# Replay determinism contract

## Goal

Replay runs must reproduce captured outcome chains with ≥99% match on golden incidents.

## Inputs

- Timeline artifact (S3) with ordered events.
- Captured outcomes per offset checkpoint.
- Sandbox Kafka topics isolated from production (`replay-sandbox.*`).

## Modes

- **strict** — inter-event delay matches original capture timestamps.
- **compressed** — delays scaled down for faster validation runs.

## Divergence

When `CompareChain` fails, orchestrator persists `reports/{replay_id}/divergence.json` with `first_mismatch_index`.

## MVP limitations

- External HTTP calls are stubbed in demo scenarios.
- Echo consumer validates topic/payload presence only.
- MSK module provisions cluster skeleton; local dev uses docker-compose Redpanda.

## Status transitions

`pending` → `running` → (`succeeded` | `failed` | `diverged` | `cancelled`).
