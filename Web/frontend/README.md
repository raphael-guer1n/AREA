# AREA – Frontend (Next.js)

Interface web pour créer et piloter les AREAs. Stack moderne (Next.js 16, React 19, TypeScript, Tailwind CSS 4) avec des hooks d’authentification, des proxys API internes et des composants UI réutilisables.

## Pile et dépendances
- Next.js 16 (App Router) + React 19 + TypeScript.
- Tailwind CSS 4 via `@import` dans `globals.css` ; palette gérée par variables CSS (vision tritanopie activable via `ColorblindToggle`).
- Zod pour la validation des formulaires.
- Dockerfile + `docker-compose.yml` pour builder/servir en conteneur.
- Pré-requis : Node.js 18.18+ (recommandé 20), npm (lockfile présent), accès au backend/gateway (par défaut sur `localhost:8080`).

## Arborescence rapide
```
src/
  app/                # Routes Next (App Router) + API routes
    area/             # Tableau + création d’AREAs
    services/         # Catalogue/connexion de services
    login, register   # Auth email+mdp
    profil/           # Profil utilisateur (SSR + client)
    auth/callback/    # Callback OAuth2
    api/auth, api/session # Proxys internes
    layout.tsx, globals.css
  components/         # UI réutilisable (cartes, formulaires, navigation)
  hooks/              # `useAuth`, `useOAuthCallback`, etc.
  lib/                # Appels API (`lib/api/*`), helpers, session
  types/              # Modèles partagés (User, Auth, Services)
public/               # Assets statiques (placeholders, logos)
docker-compose.yml, Dockerfile, eslint.config.mjs, tsconfig.json
```

## Configuration (backend/gateway)
Créer un `.env.local` à la racine du front (variables en miroir côté serveur et client) :
```bash
API_BASE_URL=http://localhost:8080/area_auth_api
NEXT_PUBLIC_API_BASE_URL=$API_BASE_URL

AREA_API_BASE_URL=http://localhost:8080/area_area_api
NEXT_PUBLIC_AREA_API_BASE_URL=$AREA_API_BASE_URL

SERVICES_API_BASE_URL=http://localhost:8080/area_service_api
NEXT_PUBLIC_SERVICES_API_BASE_URL=$SERVICES_API_BASE_URL

NEXT_PUBLIC_SITE_URL=http://localhost:3000
NEXT_PUBLIC_OAUTH_CALLBACK_BASE=http://localhost:3000
COOKIE_SECURE=false # true en prod si HTTPS
```
Notes :
- Les valeurs par défaut se basent sur le gateway exposé par `docker-compose` (backend sur `8080`).
- Côté Docker, `host.docker.internal` est utilisé pour joindre le backend depuis le container.

## Installation et lancement local
```bash
cd Web/frontend
npm install
npm run dev      # http://localhost:3000

# Production locale
npm run build
npm run start

# Qualité
npm run lint
```

## Docker
- Build : `docker build -t area-frontend .`
- Run : `docker run --rm -p 3000:3000 -e NEXT_PUBLIC_API_BASE_URL=http://host.docker.internal:8080/area_auth_api area-frontend`
- Compose (depuis ce dossier) : `docker compose up --build`
  - Variables par défaut : navigateur → `http://localhost:8080/{service}`, serveur Next → `http://host.docker.internal:8080/{service}`.
  - Override en exportant les variables avant le `compose` si le backend est ailleurs.

## Captures / GIFs
Dépose tes captures ou GIFs dans `public/docs/` (par ex. `home.png`, `services.png`, `area.png`). Les vignettes ci-dessous afficheront un placeholder tant que les fichiers ne sont pas fournis.

![Accueil](data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' width='640' height='360'><rect width='640' height='360' fill='%23f5f7fa'/><rect x='18' y='18' width='604' height='324' rx='18' fill='%23e2e8f0' stroke='%23112a46' stroke-width='3'/><text x='50%25' y='50%25' dominant-baseline='middle' text-anchor='middle' fill='%23112a46' font-family='Arial' font-size='16'>Place ta capture accueil dans public/docs/home.png</text></svg>)

![Services](data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' width='640' height='360'><rect width='640' height='360' fill='%23f5f7fa'/><rect x='18' y='18' width='604' height='324' rx='18' fill='%23e2e8f0' stroke='%23112a46' stroke-width='3'/><text x='50%25' y='50%25' dominant-baseline='middle' text-anchor='middle' fill='%23112a46' font-family='Arial' font-size='16'>Place ta capture services dans public/docs/services.png</text></svg>)

![Création d'area](data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' width='640' height='360'><rect width='640' height='360' fill='%23f5f7fa'/><rect x='18' y='18' width='604' height='324' rx='18' fill='%23e2e8f0' stroke='%23112a46' stroke-width='3'/><text x='50%25' y='50%25' dominant-baseline='middle' text-anchor='middle' fill='%23112a46' font-family='Arial' font-size='16'>Place ta capture creation dans public/docs/area.png</text></svg>)

## Liens utiles
- Architecture détaillée : `ARCHITECTURE.md`
- Composants : `COMPONENTS.md`
- Pages et routes : `PAGES_ROUTES.md`
- Contribution front : `HOWTOCONTRIBUTE.md`
