# Pages and Routes (Web/frontend)

Overview of Next.js App Router pages and the backend calls they trigger. The `/services`, `/area`, and `/profil` pages share `AreaNavigation`.

## Main routes
- `/` (Home)
  - Component: `HomePage` (`src/app/page.tsx`)
  - UI: static landing with CTA to `/login` and `/register`
  - API calls: none

- `/login`
  - Component: `LoginForm` (`src/components/forms/LoginForm.tsx`)
  - API calls: `POST {API_BASE_URL}/auth/login` via `loginRequest`, then token stored via `/api/session`
  - Actions: email/password login, Google OAuth

- `/register`
  - Component: `RegisterForm` (`src/components/forms/RegisterForm.tsx`)
  - API calls: `POST {API_BASE_URL}/auth/register`, token stored via `/api/session`
  - Actions: account creation, redirect to `/area` on success

- `/auth/callback`
  - Component: `AuthCallbackContent` + `useOAuthCallback`
  - API calls: `GET {API_BASE_URL}/oauth2/callback?code&state`, then `POST /api/session`
  - Actions: handles OAuth callback, displays status, redirects on success

- `/services`
  - Component: `ServicesClient` (`src/app/services/page.tsx`)
  - API calls:
    - `GET {SERVICE_SERVICE_BASE_URL}/providers/services`
    - `GET {API_BASE_URL}/oauth2/providers/{userId}`
    - `GET {API_BASE_URL}/oauth2/authorize?provider=...`
  - Actions: search, connect modal, OAuth window, local disconnect mock

- `/area`
  - Component: `AreaPageContent` (`src/app/area/page.tsx`)
  - API calls:
    - `GET {SERVICE_SERVICE_BASE_URL}/providers/services`
    - `GET {API_BASE_URL}/oauth2/providers/{userId}`
    - `POST {AREA_SERVICE_BASE_URL}/createEvent` via `createEventArea`
  - Actions: 3-step wizard (action, reaction, details), search/filter

- `/profil`
  - Component: `ProfilPage` (SSR) + `ProfileClient` (client)
  - API calls: `GET /api/session`, `GET {API_BASE_URL}/auth/me`
  - Actions: show auth state, local notifications, logout (`useAuth.logout` + `/api/session`)

## Internal API routes (Next)
- `GET/POST/DELETE /api/session` (`src/app/api/session/route.ts`)
  Manages the HTTP-only session cookie and validates `/auth/me`.
- `POST /api/auth` (`src/app/api/auth/route.ts`)
  Demo in-memory auth route (not used in production).

## Integration notes
- Base URLs come from `API_BASE_URL`, `AREA_API_BASE_URL`, and `SERVICES_API_BASE_URL` (plus `NEXT_PUBLIC_*`).
- OAuth flows redirect through `/auth/callback` or stay on the same page depending on context.
