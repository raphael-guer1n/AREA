/**
 * Service-service client (list providers and user connection state).
 * Uses the gateway base so we can share the same API host as auth.
 */
import { BACKEND_BASE_URL } from "./auth";

const DEFAULT_GATEWAY_URL = "http://localhost:8080";

function normalizeGatewayBase(): string {
  const base =
    process.env.API_BASE_URL ??
    process.env.NEXT_PUBLIC_API_BASE_URL ??
    DEFAULT_GATEWAY_URL;

  // Remove trailing slashes and strip service suffixes to get the gateway root.
  const withoutTrailingSlash = base.replace(/\/+$/, "");
  if (withoutTrailingSlash.endsWith("/area_auth_api")) {
    return withoutTrailingSlash.slice(0, -"/area_auth_api".length);
  }
  if (withoutTrailingSlash.endsWith("/auth-service")) {
    return withoutTrailingSlash.slice(0, -"/auth-service".length);
  }
  return withoutTrailingSlash;
}

function normalizeServiceBase(raw?: string): string {
  if (!raw) return `${normalizeGatewayBase()}/area_service_api`;
  const trimmed = raw.replace(/\/+$/, "");
  if (trimmed.endsWith("/area_service_api")) return trimmed;
  if (trimmed.endsWith("/service-service")) {
    return `${trimmed.slice(0, -"/service-service".length)}/area_service_api`;
  }
  if (trimmed.endsWith("/area_auth_api")) {
    return `${trimmed.slice(0, -"/area_auth_api".length)}/area_service_api`;
  }
  if (trimmed.endsWith("/auth-service")) {
    return `${trimmed.slice(0, -"/auth-service".length)}/area_service_api`;
  }
  return `${trimmed}/area_service_api`;
}

export const SERVICE_SERVICE_BASE_URL = normalizeServiceBase(
  process.env.SERVICES_API_BASE_URL ??
    process.env.NEXT_PUBLIC_SERVICES_API_BASE_URL,
);

export type ServiceListResponse = {
  success?: boolean;
  data?: { services?: string[] };
  error?: string;
};

export type ServiceFieldConfig = {
  name: string;
  type: string;
  label: string;
  required?: boolean;
  default?: string | number;
};

export type ServiceActionConfig = {
  title: string;
  label: string;
  type: string;
  fields: ServiceFieldConfig[];
  output_fields?: ServiceFieldConfig[];
};

export type ServiceReactionConfig = {
  title: string;
  label: string;
  url?: string;
  method?: string;
  fields: ServiceFieldConfig[];
};

export type ServiceConfig = {
  provider?: string;
  name: string;
  label?: string;
  icon_url?: string;
  actions: ServiceActionConfig[];
  reactions: ServiceReactionConfig[];
};

type ServiceConfigResponse = {
  success?: boolean;
  data?: ServiceConfig;
  error?: string;
};

export type UserServiceStatus = {
  provider: string;
  is_logged: boolean;
};

export async function fetchServices(): Promise<string[]> {
  const response = await fetch(
    `${SERVICE_SERVICE_BASE_URL}/providers/services`,
    { method: "GET", cache: "no-store" },
  );

  const body = (await response.json().catch(() => null)) as ServiceListResponse | null;

  if (!body?.success || !Array.isArray(body.data?.services)) {
    throw new Error(body?.error ?? "Impossible de récupérer les services.");
  }

  return body.data.services;
}

export async function fetchServiceConfig(serviceName: string): Promise<ServiceConfig> {
  const url = `${SERVICE_SERVICE_BASE_URL}/services/service-config?service=${encodeURIComponent(
    serviceName,
  )}`;
  const response = await fetch(url, { method: "GET", cache: "no-store" });

  const body = (await response.json().catch(() => null)) as ServiceConfigResponse | null;

  if (!response.ok || !body?.success || !body.data) {
    throw new Error(body?.error ?? "Impossible de récupérer la configuration du service.");
  }

  return body.data;
}

export async function fetchUserServiceStatuses(
  token: string,
  userId: string | number,
): Promise<UserServiceStatus[]> {
  const url = `${BACKEND_BASE_URL}/oauth2/providers/${userId}`;
  try {
    const response = await fetch(url, {
      method: "GET",
      cache: "no-store",
      headers: { Authorization: `Bearer ${token}` },
    });

    const body = (await response.json().catch(() => null)) as
      | {
          success?: boolean;
          data?: { providers?: UserServiceStatus[] };
          error?: unknown;
        }
      | null;

    const rawError =
      typeof body?.error === "string"
        ? body.error
        : body?.error
          ? JSON.stringify(body.error)
          : null;

    if (!response.ok || !body?.success || !Array.isArray(body.data?.providers)) {
      return [];
    }

    return body.data.providers;
  } catch {
    return [];
  }
}
