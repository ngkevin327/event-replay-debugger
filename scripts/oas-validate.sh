#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SPEC="$ROOT/api/openapi.yaml"
if [[ ! -f "$SPEC" ]]; then
  echo "missing $SPEC" >&2
  exit 1
fi
python3 - "$SPEC" <<'PY'
import sys
try:
    import yaml
except ImportError:
    print("PyYAML required: pip install pyyaml", file=sys.stderr)
    sys.exit(1)
with open(sys.argv[1]) as f:
    doc = yaml.safe_load(f)
assert doc.get("openapi", "").startswith("3."), "openapi version required"
assert "paths" in doc and doc["paths"], "paths required"
print("openapi ok:", sys.argv[1])
PY
