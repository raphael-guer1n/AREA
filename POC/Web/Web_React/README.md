# AREA – React + Next.js Frontend POC

## POC Introduction
A lightweight prototype to validate an initial Area user experience using React/Next.js: homepage, login form, and redirection to a profile displaying submitted credentials.

## Objectives
- Verify App Router routing and navigation between pages.
- Illustrate client-side form handling and local state.
- Prepare for future integration of server-side authentication.
- Test a light/dark theme based on CSS variables.

## Implemented Features
- Minimal homepage with a test CTA and login link.
- Login page with basic validation and authentication simulation (email `test@test.com`, password `0000`).
- Redirection to the profile page with credentials passed in the URL.
- Profile page displaying the received email/password or a message indicating no data.
- Clean styling using CSS variables, Tailwind 4, and Geist fonts.

## Installation Instructions
1. Prerequisites: Node.js 18.18+ and npm.
2. Install the dependencies:

```bash
npm install
```

## Commands to launch the project
- Development: `npm run dev` then open http://localhost:3000
- Lint: `npm run lint`
- Production build: `npm run build`
- Production start: `npm run start`

## Technologies used
- Next.js 16 (App Router)
- React 19 + React DOM
- TypeScript
- Tailwind CSS v4 + PostCSS
- Fonts Geist (Next Font)

## Strengths of React/Next.js
* Simple and intuitive file tree-based routing.
* Hybrid rendering (SSR/CSR) with hot refresh for rapid development.
* Built-in optimizations (bundling, TypeScript, font and image handling).
* Mature ecosystem, numerous complementary libraries (e.g., Tailwind for styling).

## Limitations observed (specific to the POC)
* Intentionally simplified authentication: hard-coded credentials on the client side only for demonstration purposes.
* No API calls or backend in this POC: all logic is executed on the front end.
* No global state management solution or automated tests yet (these can be added later).
* Accessibility and client-side validation aspects will be explored further in a production version.

## Quick conclusion
The POC validates navigation, form handling, and the presentation of a basic theme using Next.js. The next key steps are connecting to an authentication backend, securing data transport (sessions/tokens), and adding automated tests.