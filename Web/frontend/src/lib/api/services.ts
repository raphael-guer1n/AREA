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

export type ProviderSummary = {
  name: string;
  logo_url?: string;
};

export type ServiceListResponse = {
  success?: boolean;
  data?: { services?: string[]; providers?: ProviderSummary[] };
  error?: string;
};

export type ServiceFieldConfig = {
  name: string;
  type: string;
  label: string;
  required?: boolean;
  default?: string | number;
  selection?: Array<{ value: string; label: string }>;
  multiple?: boolean;
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
  logo_url?: string;
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
  logo_url?: string;
  need_reconnecting?: boolean;
};
async function fetchServiceList(path: string, fallbackError: string): Promise<ServiceListResponse> {
  const response = await fetch(`${SERVICE_SERVICE_BASE_URL}${path}`, {
    method: "GET",
    cache: "no-store",
  });

  const body = (await response.json().catch(() => null)) as ServiceListResponse | null;

  if (!body?.success || !body.data?.services) {
    throw new Error(body?.error ?? fallbackError);
  }

  return body;
}

export async function fetchServices(): Promise<ProviderSummary[]> {
  const body = await fetchServiceList("/providers/services", "Impossible de récupérer les services.");
  const summaries = body.data?.providers;
  if (Array.isArray(summaries) && summaries.length) {
    return summaries.map((provider) => ({
      name: provider.name,
      logo_url: provider.logo_url,
    }));
  }

  return (body.data?.services ?? []).map((name) => ({ name }));
}

export async function fetchServiceNames(): Promise<string[]> {
  const body = await fetchServiceList("/services/services", "Impossible de récupérer la liste des services.");
  return body.data?.services ?? [];
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

export async function disconnectProvider(token: string, provider: string): Promise<void> {
  const response = await fetch(`${BACKEND_BASE_URL}/oauth2/disconnect`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ provider }),
  });

  const body = (await response.json().catch(() => null)) as
    | { success?: boolean; error?: string; message?: string }
    | null;

  if (!response.ok || body?.success === false || body?.error) {
    const errorMessage =
      body?.error ??
      body?.message ??
      `Impossible de déconnecter le service (statut ${response.status}).`;
    throw new Error(errorMessage);
  }
}
