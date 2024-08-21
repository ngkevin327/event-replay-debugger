# Load tests

## Prerequisites

- Ingestion service running on `:8081`
- Valid `X-Replay-Key` with `ingest` scope
- [k6](https://k6.io/) installed

## Ingest burst (Stage 2 exit: 10k events/sec target)

```bash
export INGEST_URL=http://localhost:8081/v1/ingest/batch
export REPLAY_API_KEY=rk_live_your_key_here
export PROJECT_ID=your-project-uuid
k6 run test/load/ingest_burst.js
```

Tune `rate` in `ingest_burst.js` to match available infrastructure.
