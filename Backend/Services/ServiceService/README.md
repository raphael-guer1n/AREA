# ServiceService

The ServiceService exposes the catalog of providers, services, actions, and reactions used by the AREA platform. It serves configuration JSON to the UI and to backend services (AuthService, PollingService, WebhookService) so they can build OAuth2 flows and interpret action/reaction metadata.

## Responsibilities
- Serve the list of available providers/services.
- Expose OAuth2 configuration metadata to AuthService (internal routes).
- Expose provider configs for polling and webhooks (internal routes).
- Serve action/reaction metadata for the UI and AreaService.

## Quick Start
```bash
cd Backend/Services/ServiceService
cp .env.example .env

docker compose up -d
```

Access it through the gateway at `http://localhost:8080/area_service_api`.
Direct access (no gateway) is `http://localhost:8084`.

## API Endpoints
Public:
- **GET** `/health` - Health check
- **GET** `/providers/services` - List provider names (and logos)
- **GET** `/services/services` - List service names
- **GET** `/services/service-config?service=github` - Action/reaction metadata for a service

Internal-only (gateway requires `X-Internal-Secret`):
- **GET** `/providers/oauth2-config?service=google` - OAuth2 config for AuthService
- **GET** `/providers/config?service=google` - Full provider config
- **GET** `/webhooks/providers` - Webhook provider names
- **GET** `/webhooks/providers/config?provider=github` - Webhook provider config
- **GET** `/polling/providers` - Polling provider names
- **GET** `/polling/providers/config?provider=rss` - Polling provider config

## Configuration
`.env` variables (see `.env.example`):
```env
DB_HOST=localhost
DB_EXTERNAL_PORT=5434
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=service_service_db
SERVER_PORT=8084
```

## Config Files Layout
ServiceService loads static JSON configs from `app/internal/config/`:
- `services/` - Actions/reactions and UI metadata per service.
- `providers/` - OAuth2 provider configuration (auth URLs, scopes, tokens).
- `webhooks/` - Webhook provider rules (signature verification, setup templates).
- `polling/` - Polling provider rules (requests, parsing, filters).

Adding a new provider/service:
1. Add a `services/<service>.json` definition.
2. Add a matching OAuth2 provider config in `providers/<service>.json`.
3. If it uses polling/webhooks, add configs in `polling/` or `webhooks/`.
4. Ensure AuthService has matching client ID/secret env vars.

## OpenAPI
The OpenAPI specification is in `openapi.yaml`.
