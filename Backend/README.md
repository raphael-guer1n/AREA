# AREA Backend

Backend stack composed of the API Gateway plus Go microservices for auth, provider metadata, and automation handling. Use this README as the wiki entry for how the backend is organized and how to run it locally.

## Components
- **Gateway (`Gateway/`)**: JWT verification, permission checks, internal-only protection, CORS, and reverse-proxying to all services via `service.config.json` descriptors.
- **AuthService (`Services/AuthService/`)**: User registration/login, JWT issuance, and OAuth2 flows; PostgreSQL-backed.
- **ServiceService (`Services/ServiceService/`)**: Exposes available providers and OAuth2 configuration to clients.
- **AreaService (`Services/AreaService/`)**: Minimal action/reaction endpoint stub (`/createEvent`) to create automation events.
- **Template (`Template/`)**: Skeleton microservice used as a base for new services (Go + Postgres + Docker).

## Directory Map
```
Backend/
├── Gateway/            # API gateway (Go) + service configs
├── Services/           # Individual microservices
│   ├── AuthService/
│   ├── ServiceService/
│   └── AreaService/
├── Template/           # Microservice starter
├── services.yaml       # Service catalog (high level)
├── start-backend.sh    # Convenience launcher for dev
└── Makefile            # Aggregate helper targets
```

## Quick Start (Docker)
```bash
cd Backend
# Optional: copy env templates inside each service if missing
# AuthService/.env.example, ServiceService/.env.example, AreaService/.env.example

# Launch gateway + services
./start-backend.sh

# Options:
#   --restart-db  restart running DB containers for each service
#   --reset-db    stop services and drop DB volumes before starting (data loss)
```
This script creates the shared Docker network `area_network`, ensures `.env` files exist, starts AuthService and ServiceService (with Postgres), then boots the gateway. Logs can be tailed with `docker compose -f Backend/Services/AuthService/docker-compose.yml logs -f`.

## Running a Single Service Manually
Each service ships with:
- `Dockerfile` and `docker-compose.yml` (map `${SERVER_PORT}` to the host, default 8080; Postgres on `${DB_EXTERNAL_PORT}` default 5432).
- `.env.example` to seed required variables.
- `Makefile` for `make run`, `make test`, `make docker-up`, etc.

Example (AuthService):
```bash
cd Backend/Services/AuthService
cp .env.example .env   # edit secrets and ports
docker compose up -d   # starts API + Postgres on area_network
```

## Gateway Routing Expectations
- External traffic hits the gateway on `GATEWAY_PORT` (default 8080).
- Routes are namespaced: `/auth-service/auth/login`, `/service-service/providers/services`, `/area-service/createEvent`, etc.
- Route behaviors (auth required, permissions, internal-only) are declared in `Gateway/services-config/**/service.config.json`.

## Notes for Mobile/Web Clients
- Auth endpoints exposed via the gateway: `/auth-service/auth/register`, `/auth-service/auth/login`, `/auth-service/auth/me`.
- Provider catalog and OAuth helpers: `/service-service/providers/services`, `/service-service/providers/oauth2-config`, `/service-service/providers/config`.
- Automation creation stub: `/area-service/createEvent` (requires auth).

## Troubleshooting
- **CORS blocked**: Add your origin to `ALLOWED_ORIGINS` in `Gateway/configs/gateway.env`.
- **Invalid JWT**: Ensure `JWT_ALGO` and key/secret in `gateway.env` match the tokens issued by AuthService.
- **Route not found**: Verify the path is namespaced with the service name and that the service config defines the route/method.
