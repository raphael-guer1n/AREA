# Pages & routes (Web/frontend)

Vue d’ensemble des routes Next (App Router) et des appels backend associés. Les pages `/services`, `/area` et `/profil` partagent la barre `AreaNavigation`.

## Routes principales
- `/` (Home)  
  - Composant : `HomePage` (`src/app/page.tsx`).  
  - UI : landing statique avec CTA vers `/login` et `/register`.  
  - Appels API : aucun.

- `/login`  
  - Composant : `LoginForm` (`src/components/forms/LoginForm`).  
  - Appels API : `POST {API_BASE_URL}/auth/login` via `loginRequest`, puis persistance du token via `/api/session`.  
  - Actions : login email+mdp, OAuth Google (redirection).

- `/register`  
  - Composant : `RegisterForm` (`src/components/forms/RegisterForm`).  
  - Appels API : `POST {API_BASE_URL}/auth/register`, persistance du token via `/api/session`.  
  - Actions : création de compte, redirection vers `/area` si succès.

- `/auth/callback`  
  - Composant : `AuthCallbackContent` + `useOAuthCallback`.  
  - Appels API : `GET {API_BASE_URL}/oauth2/callback?code&state`, puis `POST /api/session` si un `token` d’app est retourné.  
  - Actions : traite le retour OAuth2, affiche l’état (processing/erreur), redirige après succès.

- `/services`  
  - Composant : `ServicesClient` (`src/app/services/page.tsx`).  
  - Appels API :  
    - `GET {SERVICE_SERVICE_BASE_URL}/providers/services` (liste des services).  
    - `GET {API_BASE_URL}/oauth2/providers/{userId}` (statut connexion par service).  
    - `GET {API_BASE_URL}/oauth2/authorize?provider=...` (démarrage OAuth “connect”).  
  - Actions : recherche, modal “connecter un service”, ouverture de la fenêtre OAuth, mock de déconnexion locale via `onDisconnect`.

- `/area`  
  - Composant : `AreaPageContent` (`src/app/area/page.tsx`).  
  - Appels API :  
    - `GET {SERVICE_SERVICE_BASE_URL}/providers/services` + `GET {API_BASE_URL}/oauth2/providers/{userId}` pour peupler les services connectés.  
    - `POST {AREA_SERVICE_BASE_URL}/createEvent` via `createEventArea` (création d’une area).  
  - Actions : wizard 3 étapes (action → réaction → détails), filtrage/recherche, modale de détail.

- `/profil`  
  - Composant : `ProfilPage` (SSR) + `ProfileClient` (client).  
  - Appels API : `GET /api/session` (lecture cookie), `GET {API_BASE_URL}/auth/me` pour précharger l’utilisateur.  
  - Actions : affichage du statut d’auth, notifications locales (stockées dans `localStorage`), bouton logout (`useAuth.logout` + `/api/session`).

## API routes internes (Next)
- `GET/POST/DELETE /api/session` (`src/app/api/session/route.ts`)  
  Stocke/lit/supprime le token en cookie HTTP-only et valide la session via `auth/me`.
- `POST /api/auth` (`src/app/api/auth/route.ts`)  
  Route démo en mémoire (non utilisée en prod) pour login/register de test.

## Notes d’intégration
- Les bases d’URL sont construites à partir des variables `API_BASE_URL`, `AREA_API_BASE_URL`, `SERVICES_API_BASE_URL` et leurs équivalents `NEXT_PUBLIC_*`.  
- Les flows OAuth (login et connect) redirigent vers `/area` ou `/services` selon le contexte, en passant par `/auth/callback` ou la même page avec `code/state`.
