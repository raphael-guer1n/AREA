# Gateway Route Catalog

This document lists every route currently loaded by the API gateway from `Gateway/services-config/**/service.config.json`.

## How to read this catalog
- The gateway accepts **namespaced** routes (`/{serviceName}{path}`) and **direct** routes (`{path}`).
- Tables below show the namespaced form. To call a direct route, drop the prefix.
- Internal-only routes require the `X-Internal-Secret` header (see `Gateway/configs/gateway.env`).

## Service summary
| Service name | Config file | Base URL | Gateway prefix |
| --- | --- | --- | --- |
| area_auth_api | Gateway/services-config/auth-service/service.config.json | http://area_auth_api:8083 | /area_auth_api |
| area_service_api | Gateway/services-config/service-service/service.config.json | http://area_service_api:8084 | /area_service_api |
| area_area_api | Gateway/services-config/area-service/service.config.json | http://area_area_api:8085 | /area_area_api |
| area_webhook_api | Gateway/services-config/webhook-service/service.config.json | http://area_webhook_api:8086 | /area_webhook_api |
| example-service | Gateway/services-config/example-service/service.config.json | http://localhost:9999/api | /example-service |

## area_auth_api
| Route (prefixed) | Methods | Auth | Internal | Permissions | Notes |
| --- | --- | --- | --- | --- | --- |
| /area_auth_api/health | GET | no | no | none | Health check |
| /area_auth_api/auth/register | POST | no | no | none | Register user |
| /area_auth_api/auth/login | POST | no | no | none | Login |
| /area_auth_api/auth/me | GET, DELETE | yes | no | none | Get or delete current user |
| /area_auth_api/oauth2/providers | GET | no | no | none | List OAuth providers |
| /area_auth_api/oauth2/authorize | GET | yes | no | none | Build OAuth authorize URL |
| /area_auth_api/oauth2/callback | GET | no | no | none | OAuth callback |
| /area_auth_api/oauth2/store | POST | no | no | none | Store OAuth tokens |
| /area_auth_api/oauth2/providers/{userId} | GET | no | no | none | List providers for user |
| /area_auth_api/oauth2/provider/token/ | GET | no | yes | none | Internal token fetch |
| /area_auth_api/oauth2/provider/profile/ | GET | no | yes | none | Internal profile fetch |
| /area_auth_api/loginwith | GET | no | no | none | OAuth login without user context |

## area_service_api
| Route (prefixed) | Methods | Auth | Internal | Permissions | Notes |
| --- | --- | --- | --- | --- | --- |
| /area_service_api/health | GET | no | no | none | Health check |
| /area_service_api/providers/services | GET | no | no | none | List provider names |
| /area_service_api/providers/oauth2-config | GET | no | yes | none | OAuth2 config lookup |
| /area_service_api/providers/config | GET | no | yes | none | Provider config lookup |
| /area_service_api/webhooks/providers | GET | no | yes | none | Webhook providers list |
| /area_service_api/webhooks/providers/config | GET | no | yes | none | Webhook provider config |
| /area_service_api/services/services | GET | no | no | none | List service names |
| /area_service_api/services/service-config | GET | no | no | none | Action/reaction metadata |

## area_area_api
| Route (prefixed) | Methods | Auth | Internal | Permissions | Notes |
| --- | --- | --- | --- | --- | --- |
| /area_area_api/health | GET | no | no | none | Health check |
| /area_area_api/createEvent | POST | yes | no | none | Create event (stub) |
| /area_area_api/saveArea | POST | yes | no | none | Save an AREA |
| /area_area_api/getAreas | GET | yes | no | none | List AREAs |
| /area_area_api/triggerArea | POST | no | yes | none | Internal trigger |

## area_webhook_api
| Route (prefixed) | Methods | Auth | Internal | Permissions | Notes |
| --- | --- | --- | --- | --- | --- |
| /area_webhook_api/health | GET | no | no | none | Health check |
| /area_webhook_api/actions | POST, PUT | yes | no | none | Create or update webhook subscriptions |
| /area_webhook_api/actions/{actionId} | GET, DELETE | yes | no | none | Fetch or delete subscription |
| /area_webhook_api/webhooks/{provider}/{hookId} | GET, POST | no | no | none | Inbound webhook receiver |

## example-service
| Route (prefixed) | Methods | Auth | Internal | Permissions | Notes |
| --- | --- | --- | --- | --- | --- |
| /example-service/test | GET | no | no | none | Example route |
