# Staging deployment

Staging runs on AWS EKS created by `infra/terraform/environments/staging`. Application workloads ship via the `replay-platform` Helm chart.

## Prerequisites

- AWS CLI v2 with `staging` profile
- `kubectl` configured for `replay-staging` EKS cluster
- `helm` 3.14+
- Terraform 1.7+

## Bootstrap infrastructure

```bash
cd infra/terraform/environments/staging
terraform init
terraform plan -out=staging.tfplan
terraform apply staging.tfplan
```

Remote state lives in S3 (`replay-terraform-state-staging`) with DynamoDB locking (`replay-terraform-locks`). See `backend.tf`.

## Configure kubectl

```bash
aws eks update-kubeconfig --name replay-staging --region us-east-1 --profile staging
kubectl get nodes
```

## Deploy application chart

Images are built by `.github/workflows/build-images.yml`. Staging deploy runs on merge to `main` via `.github/workflows/deploy-staging.yml`.

Manual install:

```bash
helm upgrade --install replay-platform deploy/helm/replay-platform \
  -f deploy/helm/replay-platform/values.yaml \
  -f deploy/helm/replay-platform/values-staging.yaml \
  --namespace replay --create-namespace
```

## Secrets

Application secrets (database URL, JWT signing key, webhook HMAC secrets) are **not** committed. Use External Secrets Operator — see [secrets.md](./secrets.md).

| Secret | Kubernetes key | Source |
|--------|----------------|--------|
| Postgres DSN | `DATABASE_URL` | AWS Secrets Manager `replay/staging/postgres` |
| JWT secret | `JWT_SECRET` | AWS Secrets Manager `replay/staging/jwt` |
| SES credentials | `AWS_*` | IRSA role on notifications worker |

## Verify health

```bash
kubectl -n replay port-forward svc/api-gateway 8080:8080
curl -s http://localhost:8080/health
curl -s http://localhost:8080/status
```

Public staging URL (after ingress): `https://staging.replay.example`

## Seed demo data

```bash
export STAGING_API_URL=https://staging.replay.example
export STAGING_ADMIN_TOKEN="<from-secrets>"
bash infra/scripts/seed-staging.sh
```

## Rollback

```bash
helm rollback replay-platform -n replay
```

Or re-deploy a previous image tag via `values-staging.yaml` `image.tag` overrides.
