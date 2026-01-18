# AreaService

AreaService is the orchestration layer that stores AREA definitions, activates or deactivates them, and triggers reactions when actions fire. It connects to AuthService for user context, ServiceService for action/reaction metadata, and action engines (Polling, Webhook, Cron, Mail).

## Responsibilities
- Persist AREA definitions and user ownership.
- Create, update, and delete action subscriptions via downstream services.
- Trigger reactions when an action fires.
- Handle activation and deactivation flows.

## Quick Start
```bash
cd Backend/Services/AreaService
cp .env.example .env

docker compose up -d
```

Access it through the gateway at `http://localhost:8080/area_area_api`.
Direct access (no gateway) is `http://localhost:8085`.

## API Endpoints
Public (auth required unless noted):
- **GET** `/health` - Health check
- **POST** `/createEvent` - Create a calendar event (OAuth2 required)
- **POST** `/saveArea` - Save an AREA definition
- **GET** `/getAreas` - List user AREAs
- **POST** `/activateArea` - Activate an AREA
- **POST** `/deactivateArea` - Deactivate an AREA
- **POST** `/deleteArea` - Delete an AREA

Internal-only (gateway requires `X-Internal-Secret`):
- **POST** `/triggerArea` - Trigger an AREA when an action fires
- **POST** `/deactivateAreasByProvider` - Deactivate all AREAs for a provider

## Configuration
`.env` variables (see `.env.example`):
```env
DB_HOST=localhost
DB_EXTERNAL_PORT=5435
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=area_area_db
SERVER_PORT=8085

AUTH_SERVICE_URL=http://gateway:8080/area_auth_api
SERVICE_SERVICE_URL=http://gateway:8080/area_service_api
AREA_SERVICE_URL=http://gateway:8080/area_area_api
INTERNAL_SECRET=secret123

CREATE_ACTIONS_URLS='{...}'
DEL_ACTIONS_URLS='{...}'
ACTIVATE_ACTIONS_URLS='{...}'
DEACTIVATE_ACTIONS_URLS='{...}'
```

The `*_ACTIONS_URLS` maps tell AreaService where to create, delete, activate, and deactivate action subscriptions (webhook, polling, cron).

## How It Works (High Level)
1. **Save AREA**: `/saveArea` validates provider connections (AuthService) and action/reaction configs (ServiceService).
2. **Action setup**: AreaService calls the configured action engine (Polling/Webhook/Cron) to create subscriptions.
3. **Trigger**: When an action fires, the engine calls `/triggerArea` (internal) to dispatch reactions.
4. **Reactions**: AreaService executes configured reactions (e.g., SMTP email) and updates status.

## OpenAPI
The OpenAPI specification is in `openapi.yaml`.
