/**
 * Service-service client:
 * - Lists available services (actions + reactions)
 * - Retrieves per-service config (fields, body templates)
 * - Fetches user connection state from the auth service
 *
 * Always routes through the gateway using the same host as auth.
 */
import { BACKEND_BASE_URL } from "./auth";

const DEFAULT_GATEWAY_URL = "http://localhost:8080";

function normalizeGatewayBase(): string {
  const base =
    process.env.API_BASE_URL ??
    process.env.NEXT_PUBLIC_API_BASE_URL ??
    DEFAULT_GATEWAY_URL;

  // Remove trailing slashes and strip the auth-service suffix if present to get the gateway root.
  const withoutTrailingSlash = base.replace(/\/+$/, "");
  return withoutTrailingSlash.replace(/\/auth-service$/, "");
}

export const SERVICE_SERVICE_BASE_URL =
  process.env.SERVICES_API_BASE_URL ??
  process.env.NEXT_PUBLIC_SERVICES_API_BASE_URL ??
  `${normalizeGatewayBase()}/service-service`;

export type ServiceFieldType = "text" | "number" | "date" | "json" | "boolean";

export type ServiceField = {
  name: string;
  type: ServiceFieldType;
  label: string;
  required?: boolean;
  default?: string;
  private?: boolean;
};

export type ServiceAction = {
  title: string;
  label?: string;
  type?: string;
  fields: ServiceField[];
  output_fields?: Array<{ name: string; type: string; label: string }>;
};

export type ServiceReaction = {
  title: string;
  label?: string;
  url?: string;
  method?: string;
  fields: ServiceField[];
  bodyType?: string;
  body_struct?: unknown;
};

export type ServiceDefinition = {
  provider?: string;
  name: string;
  label?: string;
  icon_url?: string;
  actions: ServiceAction[];
  reactions: ServiceReaction[];
};

export type ServiceListResponse = {
  success?: boolean;
  data?: { services?: string[] };
  error?: string;
};

type ServiceConfigResponse = {
  success?: boolean;
  data?: ServiceDefinition;
  error?: string;
};

export type UserServiceStatus = {
  provider: string;
  is_logged: boolean;
};

/**
 * List services (from ServiceService service configs, not mocks).
 */
export async function fetchServices(): Promise<string[]> {
  const response = await fetch(`${SERVICE_SERVICE_BASE_URL}/services/services`, {
    method: "GET",
    cache: "no-store",
  });

  const body = (await response.json().catch(() => null)) as ServiceListResponse | null;

  if (!body?.success || !Array.isArray(body.data?.services)) {
    throw new Error(body?.error ?? "Impossible de récupérer les services.");
  }

  return body.data.services;
}

export async function fetchServiceConfig(serviceName: string): Promise<ServiceDefinition> {
  const response = await fetch(
    `${SERVICE_SERVICE_BASE_URL}/services/service-config?service=${encodeURIComponent(serviceName)}`,
    { method: "GET", cache: "no-store" },
  );

  const body = (await response.json().catch(() => null)) as ServiceConfigResponse | null;

  if (!body?.success || !body.data) {
    throw new Error(body?.error ?? `Impossible de récupérer la configuration du service ${serviceName}.`);
  }

  // Normalize fields defaults to string for the UI.
  const normalizeFields = (fields?: ServiceField[]) =>
    (fields ?? []).map((field) => ({
      ...field,
      default: field.default ?? "",
    }));

  return {
    ...body.data,
    actions: (body.data.actions ?? []).map((action) => ({
      ...action,
      fields: normalizeFields(action.fields),
    })),
    reactions: (body.data.reactions ?? []).map((reaction) => ({
      ...reaction,
      fields: normalizeFields(reaction.fields),
    })),
  };
}

export async function fetchServiceCatalog(): Promise<ServiceDefinition[]> {
  const serviceNames = await fetchServices();
  const uniqueNames = Array.from(new Set(serviceNames.filter(Boolean)));

  const results = await Promise.allSettled(
    uniqueNames.map(async (name) => ({
      name,
      config: await fetchServiceConfig(name),
    })),
  );

  const definitions: ServiceDefinition[] = [];
  const errors: string[] = [];

  for (const result of results) {
    if (result.status === "fulfilled") {
      definitions.push(result.value.config);
    } else {
      errors.push(result.reason instanceof Error ? result.reason.message : String(result.reason));
    }
  }

  if (!definitions.length && errors.length) {
    throw new Error(errors[0]);
  }

  return definitions;
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
