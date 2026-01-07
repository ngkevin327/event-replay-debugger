# Key rotation procedure

Rotate JWT signing keys and webhook HMAC secrets without downtime using overlapping validity windows.

## JWT signing key (`JWT_SECRET`)

1. Generate a new secret: `openssl rand -base64 32`.
2. Set `JWT_KEY_ID=replay-v2` and deploy api-gateway with **both** secrets available:
   - `JWT_SECRET` — current signing key
   - `JWT_SECRET_PREVIOUS` — verifies tokens from prior rotation (optional env)
3. Issue new tokens with `kid` header `replay-v2` (see `apps/api-gateway/internal/auth/jwt.go`).
4. Wait for maximum access token TTL (15 minutes) plus client buffer.
5. Remove `JWT_SECRET_PREVIOUS` and retire `replay-v1` from Secrets Manager.
6. Record rotation in change log / ticket.

## Webhook HMAC secret

1. Add `webhook_secret_next` in project notification preferences (admin API).
2. Notifications worker signs with **both** secrets for 24 hours (partners verify either).
3. Partners update verification secret to `webhook_secret_next`.
4. Promote `webhook_secret_next` → `webhook_secret` via `PUT /v1/projects/{id}/notification-preferences`.
5. Clear `webhook_secret_next`.

## Verification drill

| Step | Owner | Pass criteria |
|------|-------|---------------|
| Staging login after JWT rotation | Backend | New tokens include `kid` |
| Partner webhook receives signed event | Partner | Signature verifies |
| Old JWT rejected after TTL | QA | 401 on expired kid |

## Emergency revoke

If a secret is exposed, rotate immediately, invalidate all refresh tokens in Postgres (`DELETE FROM sessions`), and force partner webhook secret rotation.
