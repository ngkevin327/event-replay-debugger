#!/usr/bin/env bash
# Verify local Replay stack health and core API flows
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

API_BASE="${API_BASE_URL:-http://localhost:8080}"
WEB_BASE="${WEB_BASE_URL:-http://localhost:5173}"
FAIL=0

record() {
  local name="$1" ok="$2" detail="$3"
  if [[ "$ok" == "1" ]]; then echo "[PASS] $name - $detail"
  else echo "[FAIL] $name - $detail"; FAIL=1; fi
}

echo "==> Dependency health"
if curl -sf "$API_BASE/health" | grep -q '"ok"'; then
  record api-health 1 "$(curl -sf "$API_BASE/health")"
else
  record api-health 0 "unreachable"
fi

if curl -sf "http://localhost:8081/health" | grep -q '"ok"'; then
  record ingestion-health 1 "$(curl -sf http://localhost:8081/health)"
else
  record ingestion-health 0 "optional: not running"
fi

if curl -sf -o /dev/null "$WEB_BASE"; then
  record web-ui 1 "status ok"
else
  record web-ui 0 "unreachable"
fi

echo "==> Auth + incident flow"
EMAIL="verify-$(date +%s)@example.com"
REG=$(curl -sf -X POST "$API_BASE/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"password-12chars\",\"org_name\":\"Verify Org\"}") || REG=""
TOKEN=$(echo "$REG" | python3 -c "import sys,json; print(json.load(sys.stdin).get('access_token',''))" 2>/dev/null || true)
if [[ -z "$TOKEN" ]]; then
  LOGIN=$(curl -sf -X POST "$API_BASE/v1/auth/login" -H "Content-Type: application/json" \
    -d "{\"email\":\"$EMAIL\",\"password\":\"password-12chars\"}") || LOGIN=""
  TOKEN=$(echo "$LOGIN" | python3 -c "import sys,json; print(json.load(sys.stdin).get('access_token',''))" 2>/dev/null || true)
fi
[[ -n "$TOKEN" ]] && record auth-register 1 "ok" || record auth-register 0 "failed"
if [[ -z "$TOKEN" ]]; then
  echo "Cannot continue without token"
  exit 1
fi

PROJ=$(curl -sf -X POST "$API_BASE/v1/projects" \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"name":"verify-project"}')
PID=$(echo "$PROJ" | python3 -c "import sys,json; print(json.load(sys.stdin).get('id',''))")
record create-project 1 "project $PID"

NOW=$(date -u +%Y-%m-%dT%H:%M:%SZ)
INC=$(curl -sf -X POST "$API_BASE/v1/projects/$PID/incidents" \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d "{\"window_start\":\"$(date -u -d '1 hour ago' +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || date -u -v-1H +%Y-%m-%dT%H:%M:%SZ)\",\"window_end\":\"$NOW\",\"topic_filters\":[\"payments.settlement\"]}")
STATUS=$(echo "$INC" | python3 -c "import sys,json; print(json.load(sys.stdin).get('status',''))")
IID=$(echo "$INC" | python3 -c "import sys,json; print(json.load(sys.stdin).get('id',''))")
[[ "$STATUS" == "ready" ]] && record create-incident 1 "status=$STATUS" || record create-incident 0 "status=$STATUS"

TL=$(curl -sf -H "Authorization: Bearer $TOKEN" "$API_BASE/v1/incidents/$IID/timeline")
EVENTS=$(echo "$TL" | python3 -c "import sys,json; d=json.load(sys.stdin); print(len((d.get('timeline') or {}).get('events') or []))")
[[ "${EVENTS:-0}" -gt 0 ]] && record timeline 1 "events=$EVENTS" || record timeline 0 "no events"

GRAPH=$(curl -sf -H "Authorization: Bearer $TOKEN" "$API_BASE/v1/incidents/$IID/graph")
NODES=$(echo "$GRAPH" | python3 -c "import sys,json; print(len(json.load(sys.stdin).get('nodes') or []))")
[[ "${NODES:-0}" -gt 0 ]] && record graph 1 "nodes=$NODES" || record graph 0 "no nodes"

REPLAY=$(curl -sf -X POST -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  "$API_BASE/v1/incidents/$IID/replays" -d '{"timing_mode":"strict"}')
RID=$(echo "$REPLAY" | python3 -c "import sys,json; print(json.load(sys.stdin).get('id',''))")
[[ -n "$RID" ]] && record create-replay 1 "replay $RID" || record create-replay 0 "failed"

exit $FAIL
