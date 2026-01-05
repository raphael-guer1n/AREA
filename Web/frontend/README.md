# AREA – Frontend (Next.js)

Interface web pour créer et piloter les AREAs (automatisations type IFTTT/Zapier). Stack moderne, prête pour la prod et documentée.

## En bref
- Next.js 16 (App Router) + React 19 + TypeScript.
- Tailwind CSS 4 (variables CSS pour le thème + mode tritanopie via `ColorblindToggle`).
- Auth/OAuth2 centralisée (`useAuth`, proxys `/api/session`).
- Pages clés : landing, login/register, catalogue services, création d’AREAs, profil.
- Dockerfile + `docker-compose.yml` pour builder/servir en conteneur.

## Arborescence rapide
```
src/
  app/                  # Routes Next (App Router) + API routes internes
    area/               # Tableau + création d’AREAs (wizard)
    services/           # Catalogue/connexion de services
    login, register     # Auth email+mdp
    profil/             # Profil (SSR + client)
    auth/callback/      # Callback OAuth2
    api/auth, api/session
    layout.tsx, globals.css
  components/           # UI réutilisable (cartes, formulaires, navigation)
  hooks/                # `useAuth`, `useOAuthCallback`, etc.
  lib/                  # Appels API (`lib/api/*`), helpers, session
  types/                # Modèles partagés (User, Auth, Services)
public/                 # Assets statiques
docker-compose.yml, Dockerfile, eslint.config.mjs, tsconfig.json
```

## Variables d’environnement (mettes-les dans `.env.local`)
```bash
# Auth (gateway)
API_BASE_URL=http://localhost:8080/auth-service
NEXT_PUBLIC_API_BASE_URL=$API_BASE_URL

# Area service
AREA_API_BASE_URL=http://localhost:8080/area-service
NEXT_PUBLIC_AREA_API_BASE_URL=$AREA_API_BASE_URL

# Services service
SERVICES_API_BASE_URL=http://localhost:8080/service-service
NEXT_PUBLIC_SERVICES_API_BASE_URL=$SERVICES_API_BASE_URL

# Site/OAuth
NEXT_PUBLIC_SITE_URL=http://localhost:3000
NEXT_PUBLIC_OAUTH_CALLBACK_BASE=http://localhost:3000
COOKIE_SECURE=false        # true en prod si HTTPS
```
Notes : valeurs par défaut alignées avec le gateway exposé par `docker-compose` (backend sur `8080`). En conteneur, le front utilise `host.docker.internal` pour joindre le backend.

## Démarrage rapide (local)
```bash
cd Web/frontend
npm install
npm run dev        # http://localhost:3000

# Vérifier la qualité / build
npm run lint
npm run build
npm run start      # sert le build
```

Scripts principaux :
- `npm run dev` : dev server avec HMR.
- `npm run lint` : ESLint (Next config).
- `npm run build` : build de prod.
- `npm run start` : sert le build.

## Docker / Compose
- Build : `docker build -t area-frontend .`
- Run : `docker run --rm -p 3000:3000 -e NEXT_PUBLIC_API_BASE_URL=http://host.docker.internal:8080/auth-service area-frontend`
- Compose (depuis ce dossier) : `docker compose up --build`
  - Navigateur → `http://localhost:8080/{service}`
  - Serveur Next (SSRs/fetch server) → `http://host.docker.internal:8080/{service}`
  - Override en exportant les variables ci-dessus si le backend est ailleurs.

## Pages & flux (résumé)
- `/` : landing statique (CTA login/register).
- `/login`, `/register` : formulaires email+mdp + OAuth Google (redirection).
- `/services` : liste des services (GET providers), état de connexion par utilisateur (GET oauth2/providers/{id}), démarrage OAuth “connect”.
- `/area` : wizard 3 étapes (action → réaction → détails) ; POST `createEvent` vers l’area-service.
- `/profil` : profil utilisateur, logout, notifications locales.
- Routes internes Next : `/api/session` (cookie HTTP-only, validation `auth/me`), `/api/auth` (démo en mémoire).

## Style & accessibilité
- Tailwind 4 via `@import` dans `globals.css` + variables CSS (`--background`, `--foreground`, etc.).
- Palette alternative tritanopie via `ColorblindToggle` (attribut `data-vision="tritanopia"` sur `<html>`).
- Composants UI réutilisables : `AreaCard`, `ServiceCard`, `Button`, `Card`, `AreaNavigation`, formulaires auth.

## Liens utiles
- Architecture : `ARCHITECTURE.md`
- Composants : `COMPONENTS.md`
- Pages et routes : `PAGES_ROUTES.md`
- Contribution front : `HOWTOCONTRIBUTE.md`
