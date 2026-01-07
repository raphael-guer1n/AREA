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

## Signature Support

WebhookService supports multiple signature styles via provider config:
- `hmac` (algorithms: `sha1`, `sha256`, `sha512`; encodings: `hex`, `base64`)
- `header` (token comparison using a header value)
- Legacy types `hmac-sha256` and `hmac-sha1` still work.

For `hmac`, you can customize the signed payload with `signing_string_template`.
Available placeholders:
- `{{body}}` (raw request body)
- `{{headers.<Header-Name>}}`
- `{{method}}`, `{{path}}`, `{{url}}`, `{{query}}`

Example (Slack-style):
```json
{
  "type": "hmac",
  "algorithm": "sha256",
  "encoding": "hex",
  "header": "X-Slack-Signature",
  "prefix": "v0=",
  "timestamp_header": "X-Slack-Request-Timestamp",
  "timestamp_tolerance_seconds": 300,
  "signing_string_template": "v0:{{headers.X-Slack-Request-Timestamp}}:{{body}}",
  "secret_json_path": "secret"
}
```
