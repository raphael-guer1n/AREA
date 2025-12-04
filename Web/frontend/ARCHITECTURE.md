# Architecture frontend AREA (Next.js 13+ App Router)

Ce projet suit une séparation stricte entre UI, logique métier (hooks) et accès réseau (lib/api). Résumé des dossiers et de leur rôle :

- `src/app/` : pages et routes API Next.js (App Router).
  - Pages (`page.tsx`) : uniquement UI et orchestration via des hooks, jamais d'appels HTTP directs.
  - Routes API (`app/api/**/route.ts`) : endpoints Next.js côté serveur qui font office de proxy vers le backend et gèrent les cookies/sessions.
- `src/components/` : composants UI **sans logique métier**.
  - `components/forms/` : formulaires qui délèguent la logique aux hooks (`useAuth`, etc.).
  - `components/ui/` : primitives UI réutilisables.
- `src/hooks/` : logique métier réutilisable côté client.
  - `useAuth` gère l’authentification (login/register/logout, refresh session) en s’appuyant sur `lib/api`.
  - `useOAuthCallback` gère le retour OAuth2, l’échange de code et la persistance de session.
- `src/lib/` : utilitaires et accès réseau.
  - `lib/api/` : **seul endroit** avec des appels HTTP (`fetch`).
    - `auth.ts` : login/register/me, OAuth2 authorize/callback, mapping des réponses backend.
    - `session.ts` : appels aux routes internes `/api/session` (persist/clear/fetch).
  - `auth.ts` : helpers localStorage (stockage token côté client si nécessaire).
  - `session.ts` : helpers cookies côté serveur (Next.js `cookies()`).
  - `helpers.ts` : utilitaires génériques (ex. `cn`).
- `src/types/` : modèles et interfaces partagées (ex. `User`, `AuthSession`, payloads).
- `public/` : assets statiques.

## Flux d’authentification (exemple)
1. Une page ou un composant déclenche une action via le hook `useAuth` ou `useOAuthCallback`.
2. Le hook appelle une fonction réseau dans `lib/api/…` :
   - `lib/api/auth.ts` pour login/register/OAuth/me.
   - `lib/api/session.ts` pour créer/supprimer/consulter la session côté Next.
3. Les fonctions `lib/api` parlent soit directement au backend (ex. `/auth/login`), soit aux routes internes Next (`/api/session`).
4. Les routes internes dans `src/app/api/**` agissent comme proxy vers le backend et gèrent les cookies (session HTTP-only).
5. Le hook met à jour l’état React (user, token, status) et la UI réagit.

## Règles essentielles
- Les pages (`app/`) ne font pas de `fetch` direct : elles utilisent des hooks.
- Les composants (`components/`) restent UI-only : pas de logique métier ni d’appels API.
- Toute logique métier est centralisée dans `hooks/`.
- Tous les appels HTTP sortants résident dans `lib/api/**`.
- Tous les modèles et types partagés sont dans `types/`.
