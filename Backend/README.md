# AREA Backend

The backend is a Go microservice stack fronted by a single API gateway. All client traffic (web and mobile) should go through the gateway, which validates auth, permissions, and internal-only routes before proxying to services.

## Architecture at a glance
```mermaid
flowchart LR
  Client[Web and Mobile clients] --> Gateway
  Gateway --> Auth[AuthService]
  Gateway --> Services[ServiceService]
  Gateway --> Area[AreaService]
  Gateway --> Webhooks[WebhookService]
  Gateway --> Polling[PollingService]
  Gateway --> Cron[CronService]
  Gateway --> Mail[MailService]
  Auth --> AuthDB[(Postgres)]
  Services --> ServicesDB[(Postgres)]
  Area --> AreaDB[(Postgres)]
  Webhooks --> WebhooksDB[(Postgres)]
  Polling --> PollingDB[(Postgres)]
  Cron --> CronDB[(Postgres)]
  Services --> Configs[(Service, Provider, Webhook configs)]
```

## Services
| Service | Responsibility | Default port | Notes |
| --- | --- | --- | --- |
| Gateway | Routing, JWT validation, permissions, internal-only protection | 8080 | Loads `Gateway/services-config/**/service.config.json` |
| AuthService | Auth, JWT issuance, OAuth2 token handling | 8083 | Postgres on 5433 by default |
| ServiceService | Provider catalog and action/reaction metadata | 8084 | Postgres on 5434 by default |
| AreaService | AREA creation and automation data | 8085 | Postgres on 5435 by default |
| WebhookService | Webhook subscriptions and inbound webhook receiver | 8086 | Optional, not started by root compose |
| PollingService | Polling-based actions (RSS/APIs) | 8087 | Optional, Postgres on 5437 by default |
| MailService | SMTP email reaction | 8088 | Optional internal-only sender |
| CronService | Timer/cron-based actions | 8089 | Optional, Postgres on 5438 by default |

## Directory layout
```
Backend/
|-- Gateway/            # API gateway (Go) + service configs
|-- Services/           # Individual microservices
|   |-- AuthService/
|   |-- ServiceService/
|   |-- AreaService/
|   `-- WebhookService/
|-- Template/           # Microservice starter
|-- services.yaml       # High-level service catalog (example)
`-- Makefile            # Global helpers for running all services
```

## Quick start (Docker)
```bash
cd Backend

# Create the shared Docker network once
docker network create area_network || true

# Seed env files if missing
cp Services/AuthService/.env.example Services/AuthService/.env
cp Services/ServiceService/.env.example Services/ServiceService/.env
cp Services/AreaService/.env.example Services/AreaService/.env
# Ensure Gateway/configs/gateway.env exists and matches your environment

# Launch gateway + services + their DBs
make docker-up

# Shut everything down
make docker-down
```
Default host ports:
- Gateway: `8080`
- AuthService: `8083` (Postgres `5433`)
- ServiceService: `8084` (Postgres `5434`)
- AreaService: `8085` (Postgres `5435`)

If you want to run optional services (Webhook, Polling, Mail, Cron), start them from their service folders or use `make docker-up` for all of them.

## Running a single service manually
Each service ships with:
- `Dockerfile` and `docker-compose.yml` (maps `${SERVER_PORT}` to the host).
- `.env.example` to seed required variables.
- `Makefile` for `make run`, `make test`, `make docker-up`, etc.

Example (AuthService):
```bash
cd Backend/Services/AuthService
cp .env.example .env

docker compose up -d
```

## Gateway routing basics
- The gateway loads routes from `Gateway/services-config/**/service.config.json`.
- Public routes can be called as either:
  - Namespaced: `/{serviceName}{path}`
  - Direct: `{path}` (no prefix, if no conflicts)
- `serviceName` comes from the `name` field in each service config. Current configs include:
  - `area_auth_api`
  - `area_service_api`
  - `area_area_api`
  - `area_webhook_api`

Internal-only routes require the `X-Internal-Secret` header (see `Gateway/configs/gateway.env`).

All routes and flags are listed in `Backend/GATEWAYS.md`.

## Configuration surface
- Gateway: `Backend/Gateway/configs/gateway.env`
  - `GATEWAY_PORT`, `INTERNAL_SECRET`, `JWT_*`, `ALLOWED_ORIGINS`, timeouts.
- Service env: `.env` files per service under `Backend/Services/*`.
- Service routing: `Backend/Gateway/services-config/**/service.config.json`.
- Provider and action/reaction metadata: `Backend/Services/ServiceService/app/internal/config/**`.

## Documentation
- Backend architecture: `Backend/ARCHITECTURE.md`
- Backend contribution guide: `Backend/HOWTOCONTRIBUTE.md`
- Gateway route catalog: `Backend/GATEWAYS.md`
- Gateway implementation notes: `Backend/Gateway/README.md`
