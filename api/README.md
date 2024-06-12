# API specification

`openapi.yaml` is the contract for the Replay control plane and ingest paths.

## Lint

```bash
npx @redocly/cli lint openapi.yaml
```

## Mock server

Use any OpenAPI mock (Prism, Stoplight) pointed at this file for frontend development before handlers land in `apps/api-gateway`.
