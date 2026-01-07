# Guide de contribution (Frontend)

## Pré-requis
- Node.js 18.18+ (20 recommandé) + npm (lockfile présent).
- Backend/gateway accessible (par défaut `localhost:8080`).
- Installe les deps : `npm install`.

## Boucle de dev
1. `npm run dev` (http://localhost:3000) pour tester en local.
2. Mets à jour `.env.local` si tu ajoutes/modifies des variables (miroir `NEXT_PUBLIC_*`).
3. Avant de pousser : `npm run lint` (et idéalement `npm run build`).

## Ajouter ou modifier un composant
- Range le composant dans `src/components/...` (UI seule) ; tape toutes les props et garde les effets réseau dans les hooks/lib.
- Réutilise les primitives (`Button`, `Card`, `AreaCard`, `ServiceCard`, `AreaNavigation`) et le helper `cn`.
- Styles : Tailwind 4 + variables CSS définies dans `globals.css`. Respecte le thème et la variation tritanopie (`ColorblindToggle`).
- Si besoin d’assets, place-les dans `public/` et référence-les via un chemin relatif (`/docs/...`).
- Ajoute une entrée dans `COMPONENTS.md` si le composant est “public”.

## Ajouter une page/route
- Crée un dossier sous `src/app/<route>/page.tsx` (App Router).
- Si la page doit partager la navigation principale, importe `AreaNavigation`.
- Place les appels réseau dans `src/lib/api/*` ou un hook dédié (`src/hooks/*`) plutôt que dans le JSX.
- Pense au rendu serveur vs client : ajoute `"use client";` si tu utilises `useState/useEffect` ou des hooks de navigation.
- Renseigne les API utilisées dans `PAGES_ROUTES.md` (composant principal + endpoints).

## Conventions de code
- Typage strict TypeScript ; évite `any`.
- Validation côté client avec Zod (`zod`), gestion d’erreurs utilisateur (messages clairs).
- Auth/OAuth : passe toujours par `useAuth` + `/api/session` pour gérer le cookie HTTP-only.
- Nom des classes Tailwind : privilégie la lisibilité (groupes logiques) ; factorise via des composants si besoin.

## Tests manuels rapides
- Auth : login/register -> redirection `/area`.
- Services : ouvrir la modale, lancer un connect OAuth, vérifier l’état “Connecté”.
- Area : créer une area via le wizard (dates/summary requis).
- Profil : vérifier affichage du token masqué et le logout.

## Documentation & captures
- Mets à jour `README.md`, `COMPONENTS.md`, `PAGES_ROUTES.md` si tu modifies la structure ou les API.
- Dépose/actualise les captures dans `public/docs/` pour garder les visuels à jour.
