# Quickstart

Install the Replay capture agent, ship your first events, and run an incident replay in under 30 minutes.

## Prerequisites

- Kubernetes 1.27+ or Docker Compose for local demo
- Kafka cluster reachable from your workloads
- Replay project API key from the control plane

## Install the agent (Helm)

```bash
helm repo add replay https://charts.replay.example
helm upgrade --install replay-agent deploy/helm/replay-agent \
  --namespace replay --create-namespace \
  --set projectId=<your-project-id> \
  --set ingest.url=https://api.replay.example/v1/ingest/batch \
  --set apiKeySecretName=replay-api-key
```

## Register and create a project

1. Open the web UI `/register` and create an organization.
2. Create a project named `payments-sandbox`.
3. Generate an API key under **Settings** and store it in the Kubernetes secret referenced by the chart.

## Capture traffic

Confirm the agent heartbeats on the dashboard **Agents** panel. Produce test traffic on your allowlisted topics.

## Create your first incident

1. Navigate to **Incidents** → **Create incident**.
2. Select a 15-minute window and topic filters (e.g. `payments`).
3. Wait for status **ready** — reconstruction builds the timeline artifact.

## Run a replay

1. Open the incident detail page → **Replay** tab.
2. Choose timing mode **strict** for deterministic comparison.
3. Review the divergence report when the replay completes.

## Verify webhook (optional)

Configure notification preferences with a webhook URL. On `incident.ready`, Replay sends a signed JSON payload — see [webhook verification](../docs/integrations/webhooks.md).

## Troubleshooting

| Symptom | Check |
|---------|-------|
| Agent offline | API key secret mounted, `REPLAY_INGEST_URL` reachable |
| Incident stuck collecting | Reconstruction worker logs, Redis queue depth |
| Replay diverged | Golden fixture `duplicate_payout_race` playbook |
