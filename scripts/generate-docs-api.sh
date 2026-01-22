#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
OUT="$ROOT/docs-site/docs/api-reference.md"
SPEC="$ROOT/api/openapi.yaml"

if command -v npx >/dev/null 2>&1; then
  npx --yes @redocly/cli@1 build-docs "$SPEC" -o /tmp/replay-api.html 2>/dev/null || true
fi

cat >"$OUT" <<'HDR'
# API reference

Generated from `api/openapi.yaml`. Regenerate with:

```bash
bash scripts/generate-docs-api.sh
```

## Authentication

- **Session**: `Authorization: Bearer <access_token>` from `/v1/auth/login`
- **API key**: `X-Replay-Key` header for ingest and agent endpoints

## Core resources

HDR

grep -E '^  /' "$SPEC" | head -40 >>"$OUT" || true
echo "" >>"$OUT"
echo "_Full OpenAPI spec: \`api/openapi.yaml\`_" >>"$OUT"
echo "Generated API reference written to $OUT"
