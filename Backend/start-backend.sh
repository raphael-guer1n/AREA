#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
NETWORK_NAME="area_network"

log() {
  printf '%s\n' "$*"
}

if docker compose version >/dev/null 2>&1; then
  COMPOSE_CMD=(docker compose)
elif docker-compose version >/dev/null 2>&1; then
  COMPOSE_CMD=(docker-compose)
else
  log "Docker Compose v2 or v1 is required (install Docker Desktop or docker-compose)."
  exit 1
fi

ensure_network() {
  if ! docker network inspect "$NETWORK_NAME" >/dev/null 2>&1; then
    log "Creating Docker network '$NETWORK_NAME'..."
    docker network create "$NETWORK_NAME" >/dev/null
  fi
}

ensure_env() {
  local dir="$1"
  local env_file="$dir/.env"
  local example_file="$dir/.env.example"

  if [[ ! -f "$env_file" && -f "$example_file" ]]; then
    cp "$example_file" "$env_file"
    log "  - $(basename "$dir"): created .env from .env.example (review secrets/ports if needed)"
  elif [[ ! -f "$env_file" && ! -f "$example_file" ]]; then
    log "  - $(basename "$dir"): no .env or .env.example found"
  fi
}

start_compose() {
  local dir="$1"
  log "Starting $(basename "$dir")..."
  (cd "$dir" && "${COMPOSE_CMD[@]}" up -d --build)
}

reset_databases() {
  local dir="$1"
  log "Resetting databases for $(basename "$dir")..."
  # Bring down stack and drop volumes (to reset Postgres data)
  (cd "$dir" && "${COMPOSE_CMD[@]}" down -v || true)
}

start_gateway() {
  local dir="$ROOT_DIR/Gateway"
  local env_file="$dir/configs/gateway.env"

  if [[ ! -f "$env_file" ]]; then
    log "Gateway env file missing at $env_file"
    exit 1
  fi

  start_compose "$dir"
}

main() {
  ensure_network

  ensure_env "$ROOT_DIR/Services/ServiceService"
  ensure_env "$ROOT_DIR/Services/AuthService"

  reset_databases "$ROOT_DIR/Services/ServiceService"
  reset_databases "$ROOT_DIR/Services/AuthService"

  start_compose "$ROOT_DIR/Services/ServiceService"
  start_compose "$ROOT_DIR/Services/AuthService"
  start_gateway

  log ""
  log "All backend containers are starting."
  log "Use 'docker ps' or '${COMPOSE_CMD[*]} -f Backend/Services/AuthService/docker-compose.yml logs -f' to inspect."
}

main "$@"
