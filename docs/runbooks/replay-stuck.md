# Stuck replay run diagnostics

## Symptoms

- Replay remains `running` beyond expected duration.
- Worker job still present in Kubernetes.
- No divergence report or terminal status update.

## Checks

1. `kubectl get jobs -l replay-run-id=<id>`
2. Feeder logs: `kubectl logs job/replay-worker-<id>`
3. Orchestrator state: `GET /v1/replays/{id}` for `divergence_summary`.
4. Sandbox topic lag on Redpanda/MSK.

## Recovery

1. Cancel via `DELETE /v1/replays/{id}`.
2. Tear down worker job via orchestrator admin endpoint.
3. Re-run replay after timeline artifact is confirmed ready.
