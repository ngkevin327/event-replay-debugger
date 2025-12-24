# Status page integration

External status providers poll public health endpoints exposed by the API gateway.

## Endpoints

| URL | Purpose | Expected |
|-----|---------|----------|
| `GET /health` | Liveness | `200` `{"status":"ok"}` |
| `GET /ready` | Readiness (DB connected) | `200` `{"status":"ready"}` |
| `GET /status` | Public status page aggregate | `200` with `components[]` |

## Better Uptime / Statuspage mapping

Configure one monitor per component returned by `/status`:

| Component | Monitor type | Target |
|-----------|--------------|--------|
| API Gateway | HTTP | `https://api.replay.example/status` |
| Postgres | Heartbeat | Synthetic via `/ready` failure alerts |
| ClickHouse | HTTP | Internal health port (staging only) |
| Ingestion | HTTP | `https://api.replay.example/health` via ingest path |

## Alerting thresholds

- **Degraded**: any component not `operational` for 2 minutes
- **Major**: `/ready` non-200 for 5 minutes
- **Critical**: `/health` non-200 for 1 minute

## Staging

```
curl -s https://staging.replay.example/status | jq .
```
