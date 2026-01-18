# AREA Project

AREA is an automation platform inspired by IFTTT/Zapier. It includes a Go microservice backend (gateway + services), a Next.js web client, and a Flutter mobile app.

## Repository Map
- `Backend/` - Go gateway + microservices
- `Web/` - Web client (Next.js)
- `Mobile/` - Mobile client (Flutter)
- `POC/` - Proofs of concept and experiments

## Architecture
High-level architecture and service responsibilities live here:
- Backend architecture: `Backend/ARCHITECTURE.md`
- Gateway routes: `Backend/GATEWAYS.md`
- Class diagrams: `CLASS_DIAGRAMS.md`
- Sequence diagram: `SEQUENCE_DIAGRAM.md`

## Installation
Requirements by component:
- Backend: Go 1.22+, Docker, Docker Compose, Make (optional)
- Web: Node.js + npm
- Mobile: Flutter SDK + platform toolchain

## Launch (Backend + Clients)
```bash
# Backend
cd Backend
docker network create area_network || true
cp Services/AuthService/.env.example Services/AuthService/.env
cp Services/ServiceService/.env.example Services/ServiceService/.env
cp Services/AreaService/.env.example Services/AreaService/.env
make docker-up

# Web
cd ../Web/frontend
npm install
npm run dev

# Mobile
cd ../../Mobile/area_mobile
flutter pub get
flutter run
```

Default ports:
- Gateway: `8080`
- Web frontend: `8081`
- AuthService: `8083`
- ServiceService: `8084`
- AreaService: `8085`

## API Server
All client traffic goes through the gateway:
- Base URL: `http://localhost:8080`
- Route catalog: `Backend/GATEWAYS.md`
- Service OpenAPI specs: `Backend/Services/*/openapi.yaml`

## Documentation
- Backend overview: `Backend/README.md`
- Backend architecture: `Backend/ARCHITECTURE.md`
- Gateway routes: `Backend/GATEWAYS.md`
- Web frontend: `Web/frontend/README.md`
- Mobile app: `Mobile/area_mobile/README.md`
- Contribution guide: `HOWTOCONTRIBUTE.md`
- Class diagrams: `CLASS_DIAGRAMS.md`
- Sequence diagram: `SEQUENCE_DIAGRAM.md`
