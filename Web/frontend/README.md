This is a [Next.js](https://nextjs.org) project bootstrapped with [`create-next-app`](https://nextjs.org/docs/app/api-reference/cli/create-next-app).

## Getting Started

First, run the development server:

```bash
npm run dev
# or
yarn dev
# or
pnpm dev
# or
bun dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

You can start editing the page by modifying `app/page.tsx`. The page auto-updates as you edit the file.

This project uses [`next/font`](https://nextjs.org/docs/app/building-your-application/optimizing/fonts) to automatically optimize and load [Geist](https://vercel.com/font), a new font family for Vercel.

## Docker

Build the production image from this folder:

```bash
docker build -t area-frontend .
```

Run it locally:

```bash
docker run --rm -p 3000:3000 area-frontend
```

Pass any required `NEXT_PUBLIC_*` environment variables with `-e` flags. The app listens on port 3000 inside the container; you can change the published port by adjusting the `-p` flag.

Run everything in one command (build + run) with Docker Compose from this folder:

```bash
docker compose up --build
```

Par défaut, Compose compile le frontend avec des URLs accessibles depuis ton navigateur et le serveur Next (backend en 8080 via le gateway) :
- Côté navigateur (`NEXT_PUBLIC_*`): `http://localhost:8080/{service}`
- Côté serveur Next (`API_*`): `http://host.docker.internal:8080/{service}`

Override them by exporting env vars before running (for example if your backend is elsewhere):

```bash
export API_BASE_URL=http://your-api-host:8080/auth-service
export NEXT_PUBLIC_API_BASE_URL=$API_BASE_URL
export AREA_API_BASE_URL=http://your-api-host:8080/area-service
export NEXT_PUBLIC_AREA_API_BASE_URL=$AREA_API_BASE_URL
export SERVICES_API_BASE_URL=http://your-api-host:8080/service-service
export NEXT_PUBLIC_SERVICES_API_BASE_URL=$SERVICES_API_BASE_URL
docker compose up --build
```

If your backend runs outside the container, set the API base when building/running so the frontend can reach it (example for a backend on host port 8080):

```bash
docker build \
  -t area-frontend \
  --build-arg API_BASE_URL=http://host.docker.internal:8080/auth-service \
  --build-arg NEXT_PUBLIC_API_BASE_URL=http://host.docker.internal:8080/auth-service \
  .

docker run --rm -p 3000:3000 \
  -e API_BASE_URL=http://host.docker.internal:8080/auth-service \
  -e NEXT_PUBLIC_API_BASE_URL=http://host.docker.internal:8080/auth-service \
  area-frontend
```

## Learn More

To learn more about Next.js, take a look at the following resources:

- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API.
- [Learn Next.js](https://nextjs.org/learn) - an interactive Next.js tutorial.

You can check out [the Next.js GitHub repository](https://github.com/vercel/next.js) - your feedback and contributions are welcome!

## Deploy on Vercel

The easiest way to deploy your Next.js app is to use the [Vercel Platform](https://vercel.com/new?utm_medium=default-template&filter=next.js&utm_source=create-next-app&utm_campaign=create-next-app-readme) from the creators of Next.js.

Check out our [Next.js deployment documentation](https://nextjs.org/docs/app/building-your-application/deploying) for more details.
