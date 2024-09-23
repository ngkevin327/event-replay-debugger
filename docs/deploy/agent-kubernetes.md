# Kubernetes agent deployment

Deploy the Replay capture agent as a sidecar or standalone Deployment.

## Prerequisites

- Helm 3
- API key secret with `ingest` scope
- Ingestion and API gateway reachable from the cluster

## Install

```bash
kubectl create secret generic replay-agent-api-key \
  --from-literal=api-key=rk_live_your_key

helm template replay-agent deploy/helm/replay-agent \
  --set projectId=<project-uuid> \
  --set ingest.url=http://ingestion.replay.svc:8081/v1/ingest/batch
```

## Upgrade

```bash
helm upgrade replay-agent deploy/helm/replay-agent -f values-prod.yaml
```

## Configuration

- `values.yaml` — image, ingest URL, topic allowlist, resource limits
- ConfigMap `agent.yaml` — allowlist and batch thresholds mounted at `/etc/replay`

## Secrets

Never commit API keys. Reference `apiKeySecretName` in values and mount via `REPLAY_API_KEY` env.
