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

export type ServiceListResponse = {
  success?: boolean;
  data?: { services?: string[] };
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
    const errorMessage =
      rawError ??
      `Impossible de récupérer l'état de connexion des services (statut ${response.status}).`;

    if (!response.ok || !body?.success || !Array.isArray(body.data?.providers)) {
      console.warn("fetchUserServiceStatuses response error:", errorMessage);
      return [];
    }

    return body.data.providers;
  } catch (error) {
    console.error("fetchUserServiceStatuses error:", error);
    return [];
  }
}
