# Contribution Guide

This repository includes multiple components (backend services, web frontend, mobile app). Use the component-specific guides for detailed workflows, and follow the sections below for backend additions.

## Prerequisites
- Go 1.22+, Docker, Docker Compose, Make (optional)
- Node.js + npm (web)
- Flutter SDK (mobile)

## Component Guides
- Backend: `Backend/HOWTOCONTRIBUTE.md`
- Web: `Web/frontend/HOWTOCONTRIBUTE.md`
- Mobile: `Mobile/area_mobile/HOWTOCONTRIBUTE.md`

## Add a Service (Backend)
Follow the full steps in `Backend/HOWTOCONTRIBUTE.md`.
Quick outline:
1. Copy the template: `Backend/Template/Microservice`.
2. Update module path, OpenAPI, and `.env.example`.
3. Add to Docker and gateway config.
4. Update service configs in ServiceService (if needed).
5. Update docs and route catalogs.

## Add an Action
Actions are defined in ServiceService configs:
- Path: `Backend/Services/ServiceService/app/internal/config/services/*.json`
- Add action metadata and fields.
- If type is `webhook` or `polling`, also add configs in `webhooks/` or `polling/`.

## Add a Reaction
Reactions are defined in the same ServiceService configs:
- Add reaction metadata, endpoint URL, method, and fields.
- Ensure provider OAuth config exists in `Backend/Services/ServiceService/app/internal/config/providers/`.
- Ensure AuthService has matching env vars for provider credentials.
