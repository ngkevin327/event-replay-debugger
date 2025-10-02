#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SPEC="$ROOT/api/openapi.yaml"
OUT="$ROOT/apps/web/src/api/generated.ts"
echo "Generating types from $SPEC -> $OUT"
# MVP: hand-maintained generated.ts; openapi-typescript optional in CI
test -f "$SPEC"
