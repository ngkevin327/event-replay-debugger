# Control plane incident response

Runbook for API gateway, metadata Postgres, and deploy pipeline failures affecting all tenants.

## Severity definitions

| SEV | Impact | Examples |
|-----|--------|----------|
| SEV1 | Full platform unavailable | EKS API down, RDS unreachable, all `/health` failing |
| SEV2 | Major feature broken | Incidents cannot be created, replay cannot start |
| SEV3 | Partial degradation | Single AZ node loss, elevated 5xx on one service |
| SEV4 | Minor | Non-critical dashboard lag, single partner webhook delay |

## First 15 minutes

1. Acknowledge PagerDuty / Slack `#incidents`.
2. Check `/status` and Grafana `replay-dashboard`.
3. Identify blast radius: auth, incidents, replay, or infra.
4. Assign roles: **IC**, **comms**, **ops**.

## Rollback (application)

```bash
# Staging or production
helm history replay-platform -n replay
helm rollback replay-platform <revision> -n replay
```

If caused by bad image tag, re-deploy last known good SHA via GitHub Actions re-run or:

```bash
helm upgrade replay-platform deploy/helm/replay-platform \
  -f deploy/helm/replay-platform/values.yaml \
  -f deploy/helm/replay-platform/values-prod.yaml \
  --set apiGateway.image.tag=<good-sha> -n replay
```

## Rollback (infrastructure)

Terraform changes:

```bash
cd infra/terraform/environments/staging
terraform plan  # review
# revert git commit and re-apply, or:
terraform apply -target=module.rds  # surgical only when documented
```

**Never** `terraform destroy` on production without executive approval.

## Communications template

> We are investigating elevated errors on the Replay control plane. Incident creation and replay may be impacted. Updates every 30 minutes at status.replay.example.

## Post-incident

- [ ] Timeline in incident doc
- [ ] Root cause category (code, config, capacity, dependency)
- [ ] Action items with owners
- [ ] Update this runbook if gaps found
