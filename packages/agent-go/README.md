# Replay Go capture agent

Customer-deployable SDK that wraps Sarama consumers/producers, buffers offline, and ships batches to Replay ingestion.

## Quick start

```bash
export REPLAY_PROJECT_ID=<uuid>
export REPLAY_API_KEY=rk_live_...
export REPLAY_KAFKA_BROKERS=localhost:9092
export REPLAY_TOPIC_ALLOWLIST=orders,payments.settlement
export REPLAY_INGEST_URL=http://localhost:8081/v1/ingest/batch
```

```go
cfg, _ := agentgo.LoadFromEnv()
ag := agentgo.New(cfg)
_ = ag.Start(context.Background())
```

See `docs/agent-go-sdk.md` for full installation and configuration.
