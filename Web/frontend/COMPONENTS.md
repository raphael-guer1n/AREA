# Component Reference (Web/frontend)

Principles: components are written in TypeScript/React and styled with Tailwind 4 + CSS variables (`globals.css`). Use the `cn` helper (`@/lib/helpers`) to combine classes, and keep network calls out of UI components (pass callbacks/props).

## Navigation and layout
- `AreaNavigation` (`src/components/navigation/AreaNavigation.tsx`)
  Bottom tab navigation (Services / Area / Profile). Highlights the active route via `usePathname`.
- `Card` (`src/components/ui/AreaCard.tsx`)
  Generic container with `title?`, `subtitle?`, `action?`, `tone?: "surface" | "background"`, `className?`.
- `Card` (`src/components/ui/ServiceCard.tsx`)
  Generic container with `title?`, `subtitle?`, `action?`, `className?`.

## Domain cards
- `AreaCard` (`src/components/area/AreaCard.tsx`)
  Gradient card for an AREA. Props: `id`, `name`, `actionLabel`, `reactionLabel`, `actionIcon`, `reactionIcon`, `isActive?`, `gradientFrom?`, `gradientTo?`, `lastRun?`, `href?`, `onClick?`, `className?`.
- `ServiceCard` (`src/components/service/ServiceCard.tsx`)
  Interactive card with detail and confirmation modals. Props: `name`, `url`, `badge`, `category?`, `gradientFrom?`, `gradientTo?`, `action?`, `actions?: string[]`, `reactions?: string[]`, `connected?: boolean`, `className?`, `onConnect?`, `onDisconnect?`.

## Auth and forms
- `LoginWithGoogle` (`src/components/LoginWithGoogle.tsx`)
  OAuth shortcut button using `useAuth.startOAuthLogin("google")`.
- `LoginForm` (`src/components/forms/LoginForm.tsx`)
  Email/username + password form with Zod validation.
- `RegisterForm` (`src/components/forms/RegisterForm.tsx`)
  Minimal signup form that calls `useAuth.register`.

## UI primitives
- `Button` (`src/components/ui/Button.tsx`)
  Variants: `primary | secondary | ghost`.
- `ColorblindToggle` (`src/components/ui/ColorblindToggle.tsx`)
  Toggles tritanopia mode and persists `data-vision="tritanopia"` in `localStorage`.

## Supporting hooks
- `useAuth` (`src/hooks/useAuth.ts`) manages auth, OAuth connect/login, and session cookies via `/api/session`.
- `useOAuthCallback` (`src/hooks/useOAuthCallback.ts`) handles OAuth callback flows and stores the session token.

## Quick examples
```tsx
import { AreaCard } from "@/components/area/AreaCard";
import { ServiceCard } from "@/components/service/ServiceCard";

<AreaCard
  id="demo-1"
  name="Post a message after an event"
  actionLabel="New ticket"
  reactionLabel="Post to Slack"
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
  actions={["New message", "New member"]}
  reactions={["Post a message"]}
  connected
  onDisconnect={() => console.log("disconnect")}
/>
```
