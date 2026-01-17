/**
 * Area-service client (create event/action). Normalizes base URL with sensible defaults
 * when env vars are not provided.
 */
import { BACKEND_BASE_URL as AUTH_BASE } from "./auth";

const DEFAULT_AREA_BASE = "http://localhost:8080/area_area_api";

function normalizeAreaBase(): string {
  const raw =
    process.env.AREA_API_BASE_URL ??
    process.env.NEXT_PUBLIC_AREA_API_BASE_URL ??
    AUTH_BASE ??
    DEFAULT_AREA_BASE;

  const trimmed = raw.replace(/\/+$/, "");
  if (trimmed.endsWith("/area_area_api")) return trimmed;
  if (trimmed.endsWith("/area_auth_api")) {
    return `${trimmed.slice(0, -"/area_auth_api".length)}/area_area_api`;
  }
  if (trimmed.endsWith("/area-service")) {
    return `${trimmed.slice(0, -"/area-service".length)}/area_area_api`;
  }
  if (trimmed.endsWith("/auth-service")) {
    return `${trimmed.slice(0, -"/auth-service".length)}/area_area_api`;
  }
  return `${trimmed}/area_area_api`;
}

export const AREA_SERVICE_BASE_URL = normalizeAreaBase();

export type CreateEventRequest = {
  delay: number;
  event: {
    startTime: string;
    endTime: string;
    summary: string;
    description: string;
  };
};

export type AreaInputField = {
  name: string;
  value: string;
};

export type AreaActionPayload = {
  id?: number;
  service: string;
  provider: string;
  title: string;
  type: string;
  input: AreaInputField[];
};

export type AreaReactionPayload = {
  id?: number;
  service: string;
  provider: string;
  title: string;
  input: AreaInputField[];
};

export type SaveAreaRequest = {
  name: string;
  active: boolean;
  actions: AreaActionPayload[];
  reactions: AreaReactionPayload[];
};

export type BackendArea = {
  id: number;
  user_id?: number;
  name: string;
  active: boolean;
  actions: AreaActionPayload[];
  reactions: AreaReactionPayload[];
};

type GetAreasResponse = {
  success?: boolean;
  data?: BackendArea[];
  error?: string;
};

function toIsoString(value: string): string {
  // Ensure the backend receives RFC3339 (with timezone) for time.Time decoding.
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    throw new Error("Date/heure invalide.");
  }
  return date.toISOString();
}

export async function createEventArea(
  token: string,
  payload: CreateEventRequest,
): Promise<void> {
  const body = {
    ...payload,
    event: {
      ...payload.event,
      startTime: toIsoString(payload.event.startTime),
      endTime: toIsoString(payload.event.endTime),
    },
  };

  const response = await fetch(`${AREA_SERVICE_BASE_URL}/createEvent`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(body),
  });

  if (!response.ok) {
    const body = (await response.json().catch(() => null)) as
      | { error?: string }
      | null;
    const errorMessage =
      body?.error ??
      `Impossible de créer l'area (statut ${response.status}).`;
    throw new Error(errorMessage);
  }
}

export async function saveArea(
  token: string,
  payload: SaveAreaRequest,
): Promise<void> {
  const response = await fetch(`${AREA_SERVICE_BASE_URL}/saveArea`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(payload),
  });

  if (!response.ok) {
    const body = (await response.json().catch(() => null)) as
      | { error?: string }
      | null;
    const errorMessage =
      body?.error ??
      `Impossible de créer l'area (statut ${response.status}).`;
    throw new Error(errorMessage);
  }
}

export async function fetchAreas(token: string): Promise<BackendArea[]> {
  const response = await fetch(`${AREA_SERVICE_BASE_URL}/getAreas`, {
    method: "GET",
    headers: {
      Authorization: `Bearer ${token}`,
    },
    cache: "no-store",
  });

  const body = (await response.json().catch(() => null)) as GetAreasResponse | null;

  if (!response.ok || !body?.success) {
    throw new Error(body?.error ?? "Impossible de récupérer les areas.");
  }

  return Array.isArray(body.data) ? body.data : [];
}

type ToggleAreaPayload = {
  area_id: number;
};

async function toggleArea(
  token: string,
  areaId: number,
  endpoint: "activateArea" | "deactivateArea",
): Promise<void> {
  const response = await fetch(`${AREA_SERVICE_BASE_URL}/${endpoint}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ area_id: areaId } satisfies ToggleAreaPayload),
  });

  const body = (await response.json().catch(() => null)) as
    | { error?: string }
    | { success?: boolean; error?: string }
    | null;

  if (!response.ok || body?.error) {
    const errorMessage =
      body?.error ??
      `Impossible de ${endpoint === "activateArea" ? "activer" : "désactiver"} l'area (statut ${
        response.status
      }).`;
    throw new Error(errorMessage);
  }
}

type DeleteAreaPayload = {
  area_id: number;
};

export async function deleteArea(token: string, areaId: number): Promise<void> {
  const response = await fetch(`${AREA_SERVICE_BASE_URL}/deleteArea`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ area_id: areaId } satisfies DeleteAreaPayload),
  });

  const body = (await response.json().catch(() => null)) as { error?: string } | null;
  if (!response.ok || body?.error) {
    const errorMessage = body?.error ?? `Impossible de supprimer l'area (statut ${response.status}).`;
    throw new Error(errorMessage);
  }
}

export async function activateArea(token: string, areaId: number): Promise<void> {
  return toggleArea(token, areaId, "activateArea");
}

export async function deactivateArea(token: string, areaId: number): Promise<void> {
  return toggleArea(token, areaId, "deactivateArea");
}
