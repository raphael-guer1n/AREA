# Polling Service

Polling service for AREA. It stores polling subscriptions in PostgreSQL, periodically calls provider APIs or RSS feeds using provider configs from ServiceService, filters the results, and dispatches events to AreaService.

## Quick Start

```bash
docker-compose up -d
```

The API should be accessed through the gateway at `http://localhost:8080/area_polling_api`.

## API Endpoints

- **GET** `/health` - Check service health
- **POST** `/actions` - Create polling subscriptions for AREA actions (requires Authorization)
- **PUT** `/actions` - Update polling subscriptions (requires Authorization)
- **GET** `/actions/{actionId}` - Fetch polling info by action_id (requires Authorization)
- **DELETE** `/actions/{actionId}` - Delete polling subscription by action_id (requires Authorization)

## Configuration

Environment variables can be configured in `.env`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=polling_service_db
SERVER_PORT=8087
SERVICE_SERVICE_URL=http://gateway:8080/area_service_api
AUTH_SERVICE_URL=http://gateway:8080/area_auth_api
AREA_SERVICE_URL=http://gateway:8080/area_area_api
POLLING_TICK_SECONDS=60
```

## Provider Configs

Polling providers are configured in `Services/ServiceService/app/internal/config/polling/*.json`. Each provider defines:
- the polling request (URL template, headers, auth, body)
- payload format (`json` or `xml`)
- how to extract and filter items
- output field mappings
- polling interval

The PollingService fetches those configs through the gateway and uses the logged-in user's OAuth2 token if a provider requires `oauth2` auth.
