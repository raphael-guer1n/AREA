# Backend Architecture

This backend is organized as a gateway plus independent Go microservices. Each service owns its own database and configuration, while the gateway centralizes routing, auth, and policy enforcement.

## System context
```mermaid
flowchart TB
  Client[Web and Mobile clients]
  subgraph Backend
    Gateway[API Gateway]
    Auth[AuthService]
    Services[ServiceService]
    Area[AreaService]
    Webhooks[WebhookService]
    Polling[PollingService]
    Cron[CronService]
    Mail[MailService]
  end
  Client --> Gateway
  Gateway --> Auth
  Gateway --> Services
  Gateway --> Area
  Gateway --> Webhooks
  Gateway --> Polling
  Gateway --> Cron
  Gateway --> Mail
  Auth --> AuthDB[(Postgres)]
  Services --> ServicesDB[(Postgres)]
  Area --> AreaDB[(Postgres)]
  Webhooks --> WebhooksDB[(Postgres)]
  Polling --> PollingDB[(Postgres)]
  Cron --> CronDB[(Postgres)]
  Services --> Configs[(Service, Provider, Webhook JSON configs)]
```

## Runtime components
- **Gateway**: Single entry point. Loads route configs from `Gateway/services-config` and applies middleware (auth, permissions, internal-only, logging).
- **AuthService**: User auth, JWT issuance, OAuth2 token storage and refresh.
- **ServiceService**: Serves provider metadata, OAuth2 config, and action/reaction definitions for the UI.
- **AreaService**: Stores AREA definitions and triggers (minimal stub today).
- **WebhookService**: Manages webhook subscriptions and receives inbound events.
- **PollingService**: Polls provider APIs and RSS feeds based on config.
- **CronService**: Schedules timer-based actions.
- **MailService**: Internal SMTP sender for email reactions.

## Routing model
Routes are defined per service in `Gateway/services-config/**/service.config.json`.

- `name` defines the public prefix: `/{serviceName}{path}`.
- The gateway also accepts unprefixed routes (`{path}`) when there is no conflict.
- Internal-only routes require `X-Internal-Secret`.

For a complete list of routes and flags, see `Backend/GATEWAYS.md`.

## Configuration model
- Gateway: `Gateway/configs/gateway.env` (port, JWT, internal secret, CORS).
- Service env: `.env` per service (ports, DB credentials, upstream URLs).
- Provider, action, reaction metadata: `Services/ServiceService/app/internal/config/`.

## Request flows

### Login and token issuance
```mermaid
sequenceDiagram
  participant C as Client
  participant G as Gateway
  participant A as AuthService
  participant DB as Auth DB

  C->>G: POST /area_auth_api/auth/login
  G->>A: Proxy request
  A->>DB: Verify credentials
  DB-->>A: User record
  A-->>G: JWT + user payload
  G-->>C: Auth response
```

### OAuth connect flow (high level)
```mermaid
sequenceDiagram
  participant C as Client
  participant G as Gateway
  participant A as AuthService
  participant S as ServiceService
  participant P as Provider

  C->>G: GET /area_auth_api/oauth2/authorize?provider=github
  G->>A: Proxy request
  A->>S: Fetch provider config (internal)
  S-->>A: OAuth2 metadata
  A-->>G: Auth URL
  G-->>C: Auth URL
  C->>P: Redirect to provider
  P-->>C: Authorization code
  C->>G: GET /area_auth_api/oauth2/callback?code=...
  G->>A: Exchange code for token
  A-->>G: Token stored
  G-->>C: Success response
```

### Create an AREA
```mermaid
sequenceDiagram
  participant C as Client
  participant G as Gateway
  participant A as AreaService
  participant DB as Area DB

  C->>G: POST /area_area_api/createEvent
  G->>A: Proxy request
  A->>DB: Persist AREA
  DB-->>A: OK
  A-->>G: Success
  G-->>C: Success
```

## Security boundaries
- **JWT validation** occurs at the gateway for routes flagged with `auth_required`.
- **Permissions** (when used) are enforced after JWT validation.
- **Internal-only** routes require `X-Internal-Secret` and are intended for service-to-service calls.
- **CORS** is configured in `gateway.env`.

## Observability
- Gateway logs request details and upstream failures.
- Each service logs independently to stdout; Docker logs are the default sink.
