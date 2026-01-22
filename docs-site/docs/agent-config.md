# Agent configuration

Helm values for `deploy/helm/replay-agent`.

## Required values

| Value | Description |
|-------|-------------|
| `projectId` | Replay project UUID |
| `ingest.url` | Ingest batch endpoint (`/v1/ingest/batch`) |
| `controlPlane.url` | API gateway base URL for register/heartbeat |
| `apiKeySecretName` | Kubernetes secret containing `api-key` |

## Example `values.yaml`

```yaml
projectId: "660e8400-e29b-41d4-a716-446655440001"
image:
  repository: ghcr.io/replay-platform/replay-agent
  tag: "0.1.0"
ingest:
  url: https://api.replay.example/v1/ingest/batch
controlPlane:
  url: https://api.replay.example
apiKeySecretName: replay-api-key
resources:
  requests:
    cpu: 100m
    memory: 128Mi
```

## Topic allowlist

Topics must match the project allowlist configured in the control plane. The agent drops events for non-allowlisted topics and increments a metric `replay_agent_events_dropped_total`.

## Upgrades

Pin `image.tag` to the platform release (`0.1.0` for MVP). Roll the DaemonSet with `helm upgrade` — no full restart of captured workloads required.
