# Ingestion lag runbook

## SLOs

| Signal | Target | Alert threshold |
|--------|--------|-----------------|
| Batch p95 latency | < 500ms | > 500ms for 5m |
| Batch error rate | < 1% | > 2% for 5m |
| ClickHouse insert failures | 0 sustained | any 5xx spike |

## Dashboards

Import `infra/grafana/ingestion-dashboard.json` for:

- `ingest_batch_latency_seconds` (histogram)
- `ingest_batch_errors_total` (counter)

## Response

1. Check ingestion pods/process health (`/health`, `/metrics`).
2. Verify ClickHouse and MinIO availability from docker-compose.
3. If `ingest_batch_errors_total` rises with 429 responses, backpressure is active — scale ClickHouse or reduce agent batch rate.
4. Review recent deploys on `apps/ingestion` for regressions in `WriteBatch` or S3 uploads.
