/**
 * Area-service client (create event/action). Normalizes base URL with sensible defaults
 * when env vars are not provided.
 */
import { BACKEND_BASE_URL as AUTH_BASE } from "./auth";

const DEFAULT_AREA_BASE = "http://localhost:8085/area-service";

function normalizeAreaBase(): string {
  const raw =
    process.env.AREA_API_BASE_URL ??
    process.env.NEXT_PUBLIC_AREA_API_BASE_URL ??
    AUTH_BASE.replace(/\/auth-service$/, "/area-service") ??
    DEFAULT_AREA_BASE;

  const trimmed = raw.replace(/\/+$/, "");
  if (/\/area-service$/.test(trimmed)) return trimmed;
  return `${trimmed}/area-service`;
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
