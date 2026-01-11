# AREA Gateway

Core API Gateway for routing, auth, permissions, and internal service communication.

## Base URL
- Default: `http://localhost:8080`
- Gateway accepts both the namespaced form (`/{service-name}{path}`) and the raw route path.
- Clients should always use the namespaced form to avoid collisions.

## Routes

### auth-service
- `GET /auth-service/health` (auth: no)
- `POST /auth-service/auth/register` (auth: no)
- `POST /auth-service/auth/login` (auth: no)
- `GET /auth-service/auth/me` (auth: yes)
- `DELETE /auth-service/auth/me` (auth: yes)
- `GET /auth-service/oauth2/providers` (auth: no)
- `GET /auth-service/oauth2/authorize` (auth: yes)
- `GET /auth-service/oauth2/callback` (auth: no)
- `POST /auth-service/oauth2/store` (auth: no)
- `GET /auth-service/oauth2/providers/{userId}` (auth: no)
- `GET /auth-service/oauth2/provider/token/` (auth: no)
- `GET /auth-service/loginwith` (auth: no)

### service-service
- `GET /service-service/health` (auth: no)
- `GET /service-service/providers/services` (auth: no)
- `GET /service-service/providers/oauth2-config` (auth: no)
- `GET /service-service/providers/config` (auth: no)
- `GET /service-service/services/services` (auth: no)
- `GET /service-service/services/service-config` (auth: no, query: `service`)
- `GET /service-service/webhooks/providers` (auth: no)
- `GET /service-service/webhooks/providers/config` (auth: no)

### area-service
- `GET /area-service/health` (auth: no)
- `POST /area-service/createEvent` (auth: yes)

### webhook-service
- `GET /webhook-service/health` (auth: no)
- `GET /webhook-service/subscriptions` (auth: no)
- `POST /webhook-service/subscriptions` (auth: no)
- `GET /webhook-service/subscriptions/{hookId}` (auth: no)
- `DELETE /webhook-service/subscriptions/{hookId}` (auth: no)
- `GET /webhook-service/webhooks/{provider}/{hookId}` (auth: no)
- `POST /webhook-service/webhooks/{provider}/{hookId}` (auth: no)

### example-service
- `GET /example-service/test` (auth: no)
