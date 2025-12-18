# Secret management (staging / production)

Replay never stores secrets in git. Staging and production use **External Secrets Operator (ESO)** to sync AWS Secrets Manager values into Kubernetes.

## Pattern

```yaml
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: api-gateway-secrets
  namespace: replay
spec:
  refreshInterval: 1h
  secretStoreRef:
    name: aws-secrets-manager
    kind: ClusterSecretStore
  target:
    name: api-gateway-secrets
  data:
    - secretKey: DATABASE_URL
      remoteRef:
        key: replay/staging/postgres
        property: url
    - secretKey: JWT_SECRET
      remoteRef:
        key: replay/staging/jwt
        property: signing_key
```

## IRSA for workers

Terraform module `infra/terraform/modules/iam` provisions IRSA roles for:

- `ingestion` — S3 prefix `projects/*`, `incidents/*`
- `replay-worker` — same bucket, read-heavy for replay artifacts

Annotate service accounts in Helm values:

```yaml
serviceAccount:
  annotations:
    eks.amazonaws.com/role-arn: arn:aws:iam::ACCOUNT:role/replay-staging-ingestion-irsa
```

## Rotation

1. Update secret in AWS Secrets Manager.
2. ESO refreshes within `refreshInterval` (default 1h) or force-sync:
   `kubectl annotate externalsecret api-gateway-secrets force-sync=$(date +%s)`
3. Rolling restart: `kubectl rollout restart deployment/api-gateway -n replay`
