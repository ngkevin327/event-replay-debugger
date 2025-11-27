# Webhook verification

Replay signs outbound webhook payloads with HMAC-SHA256. Partners verify deliveries using the shared secret configured in project notification preferences.

## Headers

| Header | Description |
|--------|-------------|
| `Content-Type` | `application/json` |
| `X-Replay-Signature` | `sha256=<hex>` HMAC of raw body |

## Verification (Go)

```go
import "github.com/replay/platform/apps/notifications/internal/webhook"

func verify(secret string, body []byte, header string) error {
    return webhook.VerifySignature(secret, body, header)
}
```

## Verification (Node.js)

```javascript
const crypto = require("crypto");

function verify(secret, body, header) {
  const expected =
    "sha256=" +
    crypto.createHmac("sha256", secret).update(body).digest("hex");
  return crypto.timingSafeEqual(Buffer.from(expected), Buffer.from(header));
}
```

## Events

| Event | When |
|-------|------|
| `incident.ready` | Incident reconstruction finished |
| `replay.completed` | Replay run reached terminal state |

## Example payload

```json
{
  "event": "incident.ready",
  "incident_id": "550e8400-e29b-41d4-a716-446655440000",
  "project_id": "660e8400-e29b-41d4-a716-446655440001"
}
```
