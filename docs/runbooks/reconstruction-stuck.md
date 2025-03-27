# Reconstruction stuck jobs

## Symptoms

- Incident remains in `collecting` longer than 15 minutes.
- Reconstruction worker logs show repeated idle pops or job failures.
- Timeline API returns `409` (not ready).

## Diagnostics

1. Check worker process: `kubectl logs deploy/reconstruction-worker`.
2. Inspect queue depth (Redis `LLEN reconstruction:jobs` or in-memory queue in local dev).
3. Verify ClickHouse returns rows for the incident window:
   `SELECT count() FROM replay.events WHERE project_id = '<id>' AND timestamp BETWEEN ...`.
4. Confirm S3/MinIO credentials if hydration errors appear in worker logs.

## Recovery

1. Re-enqueue collection: call internal enqueue with `incident_id`.
2. If ClickHouse is empty, mark incident `failed` after validating ingest lag.
3. Restart worker deployment after fixing upstream dependency outages.

## Prevention

- Alert on `reconstruction_job_failures_total` growth.
- Track `reconstruction_job_duration_ms` p95 against SLO (15m collection bound).
