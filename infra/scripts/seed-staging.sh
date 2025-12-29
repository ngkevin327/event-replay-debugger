#!/usr/bin/env bash
set -euo pipefail

API_URL="${STAGING_API_URL:-https://staging.replay.example}"
ADMIN_TOKEN="${STAGING_ADMIN_TOKEN:-}"
SEED_FILE="$(dirname "$0")/seed-data/demo-incidents.json"

if [[ ! -f "$SEED_FILE" ]]; then
  echo "missing seed file: $SEED_FILE" >&2
  exit 1
fi

auth_header() {
  if [[ -n "$ADMIN_TOKEN" ]]; then
    echo "Authorization: Bearer $ADMIN_TOKEN"
  else
    echo ""
  fi
}

register_demo_user() {
  local email password org_name
  email=$(jq -r '.user.email' "$SEED_FILE")
  password=$(jq -r '.user.password' "$SEED_FILE")
  org_name=$(jq -r '.user.org_name' "$SEED_FILE")
  curl -sf -X POST "$API_URL/v1/auth/register" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$email\",\"password\":\"$password\",\"org_name\":\"$org_name\"}" \
    || true
}

login_and_token() {
  local email password
  email=$(jq -r '.user.email' "$SEED_FILE")
  password=$(jq -r '.user.password' "$SEED_FILE")
  curl -sf -X POST "$API_URL/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$email\",\"password\":\"$password\"}" | jq -r '.access_token'
}

create_project() {
  local token name
  token=$1
  name=$(jq -r '.project.name' "$SEED_FILE")
  curl -sf -X POST "$API_URL/v1/projects" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $token" \
    -d "{\"name\":\"$name\"}" | jq -r '.id'
}

seed_incidents() {
  local token project_id
  token=$1
  project_id=$2
  jq -c '.incidents[]' "$SEED_FILE" | while read -r inc; do
    curl -sf -X POST "$API_URL/v1/projects/$project_id/incidents" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $token" \
      -d "{\"window_start\":$(echo "$inc" | jq '.window_start'),\"window_end\":$(echo "$inc" | jq '.window_end'),\"topic_filters\":$(echo "$inc" | jq '.topic_filters')}" \
      >/dev/null
    echo "seeded incident $(echo "$inc" | jq -r '.label')"
  done
}

main() {
  echo "Seeding staging at $API_URL"
  register_demo_user
  token=$(login_and_token)
  if [[ -z "$token" || "$token" == "null" ]]; then
    echo "failed to obtain access token" >&2
    exit 1
  fi
  project_id=$(create_project "$token")
  seed_incidents "$token" "$project_id"
  echo "Staging seed complete (project=$project_id)"
}

main "$@"
