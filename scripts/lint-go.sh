#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

if ! command -v golangci-lint >/dev/null 2>&1; then
  echo "golangci-lint not installed; skipping Go lint"
  exit 0
fi

while IFS= read -r -d '' mod; do
  dir="$(dirname "$mod")"
  echo "linting $dir"
  (cd "$dir" && golangci-lint run ./...)
done < <(find apps packages services -name go.mod -print0 2>/dev/null)

echo "go lint complete"
