# API reference

Generated from `api/openapi.yaml`. Regenerate with:

```bash
bash scripts/generate-docs-api.sh
```

## Authentication

- **Session**: `Authorization: Bearer <access_token>` from `/v1/auth/login`
- **API key**: `X-Replay-Key` header for ingest and agent endpoints

## Core resources

| Method | Path | Summary |
|--------|------|---------|
| POST | `/v1/auth/register` | Register user and org |
| POST | `/v1/auth/login` | Obtain access token |
| GET | `/v1/projects` | List projects |
| POST | `/v1/projects/{projectId}/incidents` | Create incident |
| GET | `/v1/incidents/{incidentId}/timeline` | Fetch timeline |
| POST | `/v1/incidents/{incidentId}/replays` | Start replay |
| GET | `/v1/projects/{id}/notification-preferences` | Get notification prefs |

_Full OpenAPI spec: `api/openapi.yaml`_
