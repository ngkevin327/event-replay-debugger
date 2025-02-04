#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SPEC="$ROOT/api/openapi.yaml"
BASE="${1:-origin/main}"
if ! git -C "$ROOT" rev-parse --verify "$BASE" >/dev/null 2>&1; then
  echo "skip oas-diff: base $BASE not found"
  exit 0
fi
if git -C "$ROOT" diff --quiet "$BASE" -- "$SPEC"; then
  echo "openapi unchanged vs $BASE"
  exit 0
fi
echo "openapi changed vs $BASE; review for breaking API changes" >&2
git -C "$ROOT" diff --stat "$BASE" -- "$SPEC"
exit 0
