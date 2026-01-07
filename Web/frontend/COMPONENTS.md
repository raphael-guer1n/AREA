# Documentation des composants (Web/frontend)

Principes : composants présentés ci-dessous sont écrits en TypeScript/React et stylés avec Tailwind 4 + variables CSS (`globals.css`). Utiliser le helper `cn` (`@/lib/helpers`) pour combiner des classes, et garder les appels réseau hors des composants UI (passer des callbacks/props).

## Navigation & layout
- `AreaNavigation` (`src/components/navigation/AreaNavigation.tsx`)  
  Barre d’onglets (Services / Area / Profil) qui surligne l’onglet actif via `usePathname`. Props : aucune.
- `Card` (`src/components/ui/AreaCard.tsx`)  
  Conteneur générique avec `title?`, `subtitle?`, `action?`, `tone?: "surface" | "background"`, `className?`. Sert de wrapper pour les sections principales.
- `Card` (`src/components/ui/ServiceCard.tsx`)  
  Variante identique pour les cartes de services ; même API de props que ci-dessus.

## Cartes métier
- `AreaCard` (`src/components/area/AreaCard.tsx`)  
  Carte gradient pour une AREA. Props : `id`, `name`, `actionLabel`, `reactionLabel`, `actionIcon`, `reactionIcon`, `isActive?`, `gradientFrom?`, `gradientTo?`, `lastRun?`, `href?`, `onClick?`, `className?`. Si `onClick` est fourni, la carte rend un bouton ; sinon un lien vers `href` ou `/area/{id}`.
- `ServiceCard` (`src/components/service/ServiceCard.tsx`)  
  Carte interactive (modale de détails + confirm connect/disconnect). Props : `name`, `url`, `badge`, `category?`, `gradientFrom?`, `gradientTo?`, `action?` (label “À connecter” par défaut), `actions?: string[]`, `reactions?: string[]`, `connected?: boolean`, `className?`, `onConnect?`, `onDisconnect?`. Gère l’état local `isConnected` et affiche deux modales intégrées (détails + confirmation).
- `LoginWithGoogle` (`src/components/LoginWithGoogle.tsx`)  
  Bouton OAuth rapide. Props : `label?` et `className?`. Utilise `useAuth.startOAuthLogin("google")`.

## Auth & formulaires
- `LoginForm` (`src/components/forms/LoginForm.tsx`)  
  Formulaire email/mdp avec validation Zod, déclenche `useAuth.login` puis redirige vers `/area`. Inclut bouton Google via `startOAuthLogin`.
- `RegisterForm` (`src/components/forms/RegisterForm.tsx`)  
  Formulaire d’inscription minimal (email, username, password) qui appelle `useAuth.register`, puis redirige sur succès.

## Primitives UI
- `Button` (`src/components/ui/Button.tsx`)  
  Variantes `primary | secondary | ghost`. Accepte toutes les props HTML de bouton + `className`.
- `ColorblindToggle` (`src/components/ui/ColorblindToggle.tsx`)  
  Bascule tritanopie. Persiste l’état en `localStorage` et positionne `data-vision="tritanopia"` sur `<html>`.

## Hooks support utilisés par les composants
- `useAuth` (`src/hooks/useAuth.ts`) : machine à états client pour auth, login/register, OAuth connect/login, gestion token (HTTP-only cookie via `/api/session`).
- `useOAuthCallback` (`src/hooks/useOAuthCallback.ts`) : traite le retour `code/state` OAuth2, persiste le token via `/api/session`, puis redirige.

## Exemples rapides
```tsx
import { AreaCard } from "@/components/area/AreaCard";
import { ServiceCard } from "@/components/service/ServiceCard";

<AreaCard
  id="demo-1"
  name="Poster un message après un événement"
  actionLabel="Nouveau ticket"
  reactionLabel="Post Slack"
  actionIcon={<span>TI</span>}
  reactionIcon={<span>SL</span>}
  isActive
  onClick={() => setOpen(true)}
/>

<ServiceCard
  name="Slack"
  url="https://slack.com"
  badge="SL"
  category="Messaging"
  actions={["Nouveau message", "Nouveau membre"]}
  reactions={["Poster un message"]}
  connected={true}
  onDisconnect={() => console.log("disconnect")}
/>
```
