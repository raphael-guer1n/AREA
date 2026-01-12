# Mobile Architecture (Flutter)

The mobile app is built around a clean separation between UI screens, providers, and API services. Network calls are isolated in service classes, while UI screens remain focused on rendering and user interaction.

## Layered structure
- `lib/screens/`: UI pages (auth, services, profile, main shell).
- `lib/providers/`: state management (`AuthProvider`).
- `lib/services/`: HTTP and OAuth helpers (`AuthService`, `ServiceConnector`).
- `lib/models/`: client-side models.
- `lib/theme/`: colors, typography, and theme.

## Data flow
```mermaid
flowchart LR
  Screen[Screen] --> Provider[AuthProvider]
  Screen --> Service[ServiceConnector]
  Provider --> AuthSvc[AuthService]
  AuthSvc --> Gateway
  Service --> Gateway
```

## Auth and session flow
```mermaid
sequenceDiagram
  participant U as User
  participant UI as Flutter UI
  participant A as AuthService (mobile)
  participant G as Gateway
  participant B as AuthService (backend)

  U->>UI: Submit credentials
  UI->>A: loginWithEmail
  A->>G: POST /auth-service/auth/login
  G->>B: Proxy request
  B-->>G: JWT + user
  G-->>A: Response
  A-->>UI: Save token to secure storage
```

## OAuth deep links
- The app listens for deep links with the scheme/host `area://auth`.
- Android intent filter is defined in `android/app/src/main/AndroidManifest.xml`.
- OAuth login uses `app_links` to capture the redirect and store the token.

## Configuration
- Base URL comes from `.env` (`BASE_URL`).
- Most routes are built as `$BASE_URL/auth-service/...` in `lib/services/*`.

## Error handling
- Service methods throw exceptions with readable messages.
- UI layers catch and display errors via banners or dialogs.
