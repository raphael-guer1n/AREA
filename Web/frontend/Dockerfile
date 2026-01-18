# syntax=docker/dockerfile:1

FROM node:20-bookworm-slim AS base
WORKDIR /app
ENV NEXT_TELEMETRY_DISABLED=1

FROM base AS deps
ARG NODE_ENV=development
ENV NODE_ENV=${NODE_ENV}
COPY package.json package-lock.json ./
RUN npm ci --include=dev

FROM base AS builder
ARG NODE_ENV=production
ARG API_BASE_URL
ARG NEXT_PUBLIC_API_BASE_URL
ARG AREA_API_BASE_URL
ARG NEXT_PUBLIC_AREA_API_BASE_URL
ARG SERVICES_API_BASE_URL
ARG NEXT_PUBLIC_SERVICES_API_BASE_URL
ENV NODE_ENV=${NODE_ENV}
ENV API_BASE_URL=${API_BASE_URL}
ENV NEXT_PUBLIC_API_BASE_URL=${NEXT_PUBLIC_API_BASE_URL}
ENV AREA_API_BASE_URL=${AREA_API_BASE_URL}
ENV NEXT_PUBLIC_AREA_API_BASE_URL=${NEXT_PUBLIC_AREA_API_BASE_URL}
ENV SERVICES_API_BASE_URL=${SERVICES_API_BASE_URL}
ENV NEXT_PUBLIC_SERVICES_API_BASE_URL=${NEXT_PUBLIC_SERVICES_API_BASE_URL}
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN npm run build

FROM base AS runner
WORKDIR /app
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
ENV PORT=8081
ENV HOSTNAME=0.0.0.0

RUN groupadd -r nextjs && useradd -r -g nextjs nextjs

COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static

USER nextjs
EXPOSE 8081
CMD ["node", "server.js"]
