# Go capture agent SDK

Install and configure the Replay capture agent in under 15 minutes.

## Requirements

- Go 1.25+
- Kafka cluster reachable from the host
- Replay API key with `ingest` scope
- Ingestion service (`POST /v1/ingest/batch`)

## Installation

```bash
go get github.com/replay/platform/packages/agent-go@latest
```

In a Go workspace that includes `packages/agent-go`:

```go
import agentgo "github.com/replay/platform/packages/agent-go"
```

## Configuration

| Variable | Description |
|----------|-------------|
| `REPLAY_PROJECT_ID` | Project UUID |
| `REPLAY_API_KEY` | `rk_live_…` ingest key |
| `REPLAY_KAFKA_BROKERS` | Comma-separated brokers |
| `REPLAY_TOPIC_ALLOWLIST` | Topics to capture |
| `REPLAY_INGEST_URL` | Batch ingest endpoint |
| `REPLAY_BUFFER_DIR` | Offline disk buffer path |
| `REPLAY_BATCH_SIZE` | Events per flush (default 100) |
| `REPLAY_FLUSH_INTERVAL` | Flush interval (default 5s) |

## Wrap a consumer

```go
hook := func(ev event.CapturedEvent) { /* buffer or ship */ }
wrapped := kafka.WrapConsumerHandler(yourHandler, hook, cfg.CaptureTimeout)
group.Consume(ctx, topics, wrapped)
```

## Topic allowlist

Only topics listed in `REPLAY_TOPIC_ALLOWLIST` are captured when the allowlist is non-empty.

## Redaction

Set `REPLAY_REDACT_RULES` to enable default PII masking on payloads before upload.

## Kubernetes

Use the Helm chart at `deploy/helm/replay-agent` (see `docs/deploy/agent-kubernetes.md`).
