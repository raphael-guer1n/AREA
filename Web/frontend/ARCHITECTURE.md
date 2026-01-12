# Frontend Architecture (Next.js App Router)

This frontend keeps UI, client business logic, and network access in separate layers. The key folders are:

- `src/app/`: Routes (server and client) plus API routes.
  - UI pages (`page.tsx`) orchestrate hooks and helpers; avoid inlining HTTP. The `area` and `services` pages currently call `lib/api` directly for data load/creation—wrap in hooks if you want stricter separation.
  - API routes (`app/api/**/route.ts`) act as thin proxies to backend services and manage cookies/sessions server-side.
- `src/components/`: UI-only components.
  - `components/forms/`: forms that delegate logic to hooks such as `useAuth`.
  - `components/ui/`: small reusable primitives.
- `src/hooks/`: reusable client-side business logic.
  - `useAuth` handles login/register/logout, session refresh, and OAuth kick-off using `lib/api`.
  - `useOAuthCallback` processes OAuth2 callbacks and persists the session token.
- `src/lib/`: shared utilities and network calls.
  - `lib/api/**`: the only place with `fetch` to external/back-end services or internal routes (`auth`, `session`, `services`, `area`).
  - `session.ts`: server helpers around `cookies()`; `auth.ts`: localStorage helpers; `helpers.ts`: misc utilities.
- `src/types/`: shared models such as `User` and auth/session contracts.
- `public/`: static assets.

## Data & Auth Flow (simplified)
1. A page or component triggers a hook (`useAuth`, `useOAuthCallback`) or, in some cases, calls `lib/api` helpers directly.
2. `lib/api` functions talk either to backend services (Auth, Services, Area) or to internal routes under `/api`.
3. Internal API routes proxy requests and set/clear the HTTP-only session cookie.
4. Hooks update React state (user, token, status); UI components react to that state.

## Internal API routes
- `api/session`: read/write the session cookie, proxy `/auth/me` to validate the token.
- `api/auth/login` and `api/auth/register`: proxy credential-based auth to the backend.
- `api/auth`: small in-memory demo endpoint (not used for real auth in production).

## Configuration surface
- Auth base URL: `API_BASE_URL` / `NEXT_PUBLIC_API_BASE_URL` (defaults to `.../area_auth_api`).
- Area service: `AREA_API_BASE_URL` / `NEXT_PUBLIC_AREA_API_BASE_URL`.
- Services service: `SERVICES_API_BASE_URL` / `NEXT_PUBLIC_SERVICES_API_BASE_URL`.
- OAuth/site URLs: `NEXT_PUBLIC_SITE_URL`, `NEXT_PUBLIC_OAUTH_CALLBACK_BASE`, `COOKIE_SECURE` flags.
