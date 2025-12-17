#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
NETWORK_NAME="area_network"
RESTART_DBS="false"
RESET_DBS="false"

log() {
  printf '%s\n' "$*"
}

usage() {
  cat <<EOF
Usage: $(basename "$0") [options]

Options:
  --restart-db      Restart the database containers for each service before starting
  --reset-db        Stop services and remove DB volumes before starting (data loss)
  -h, --help        Show this help message
EOF
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

restart_db() {
  local dir="$1"
  log "Restarting database for $(basename "$dir")..."
  if (cd "$dir" && "${COMPOSE_CMD[@]}" restart db >/dev/null 2>&1); then
    log "  - db restarted."
  else
    log "  - db not running, starting it..."
    (cd "$dir" && "${COMPOSE_CMD[@]}" up -d db)
  fi
}

reset_db() {
  local dir="$1"
  log "Resetting database (down -v) for $(basename "$dir")..."
  (cd "$dir" && "${COMPOSE_CMD[@]}" down -v || true)
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

start_gateway() {
  local dir="$ROOT_DIR/Gateway"
  local env_file="$dir/configs/gateway.env"

  if [[ ! -f "$env_file" ]]; then
    log "Gateway env file missing at $env_file"
    exit 1
  fi

  start_compose "$dir"
}

parse_args() {
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --restart-db|--restart-dbs)
        RESTART_DBS="true"
        ;;
      --reset-db|--reset-dbs)
        RESET_DBS="true"
        ;;
      -h|--help)
        usage
        exit 0
        ;;
      *)
        log "Unknown option: $1"
        usage
        exit 1
        ;;
    esac
    shift
  done
}

main() {
  parse_args "$@"
  ensure_network

  ensure_env "$ROOT_DIR/Services/ServiceService"
  ensure_env "$ROOT_DIR/Services/AuthService"

  if [[ "$RESET_DBS" == "true" ]]; then
    RESTART_DBS="false" # reset has priority
    reset_db "$ROOT_DIR/Services/ServiceService"
    reset_db "$ROOT_DIR/Services/AuthService"
  fi

  start_compose "$ROOT_DIR/Services/ServiceService"
  start_compose "$ROOT_DIR/Services/AuthService"

  if [[ "$RESTART_DBS" == "true" ]]; then
    restart_db "$ROOT_DIR/Services/ServiceService"
    restart_db "$ROOT_DIR/Services/AuthService"
  fi

  start_gateway

  log ""
  log "All backend containers are starting."
  log "Use 'docker ps' or '${COMPOSE_CMD[*]} -f Backend/Services/AuthService/docker-compose.yml logs -f' to inspect."
}

main "$@"
