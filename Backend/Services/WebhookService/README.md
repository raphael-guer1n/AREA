# Webhook Service

Webhook receiver and subscription service for AREA. It stores webhook subscriptions in PostgreSQL, validates incoming webhooks using provider configs defined in ServiceService, and can auto-register webhooks for OAuth2 providers.

## Quick Start

```bash
# Build and start containers
docker-compose up -d
```

The API should be accessed through the gateway at `http://localhost:8080/webhook-service`.

## API Endpoints

- **GET** `/health` - Check service health
- **POST** `/subscriptions` - Create a webhook subscription (intended for internal use via AreaService)
- **GET** `/subscriptions/{hookId}` - Fetch a subscription
- **POST** `/webhooks/{provider}/{hookId}` - Receive a webhook

## Configuration

Environment variables can be configured in `.env`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=webhook_service_db
SERVER_PORT=8085
SERVICE_SERVICE_URL=http://gateway:8080/service-service
AUTH_SERVICE_URL=http://gateway:8080/auth-service
AREA_SERVICE_URL=http://gateway:8080/area-service
PUBLIC_BASE_URL=https://api.example.com/webhook-service
```

## Provider Configs

Webhook providers (signature rules, event headers, mappings, setup templates) are configured in `Services/ServiceService/app/internal/config/webhooks/*.json` and served by ServiceService. If a provider defines a `setup` block with OAuth2 auth, WebhookService will create the webhook automatically using the user's OAuth2 token from AuthService.

When running behind the gateway, `PUBLIC_BASE_URL` must include the `/webhook-service` prefix so providers call the gateway route rather than the service directly.
