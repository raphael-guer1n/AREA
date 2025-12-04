import type { SessionStatusResponse } from "@/types/auth";

export async function persistSessionToken(token: string): Promise<void> {
  const response = await fetch("/api/session", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    credentials: "include",
    body: JSON.stringify({ token }),
  });

  if (!response.ok) {
    throw new Error("Impossible de sauvegarder la session.");
  }
}

export async function clearSession(): Promise<void> {
  const response = await fetch("/api/session", {
    method: "DELETE",
    credentials: "include",
  });

  if (!response.ok) {
    throw new Error("Impossible de fermer la session.");
  }
}

export async function fetchSessionStatus(): Promise<SessionStatusResponse> {
  const response = await fetch("/api/session", {
    method: "GET",
    credentials: "include",
    cache: "no-store",
  });

  const body = (await response.json().catch(() => null)) as SessionStatusResponse | null;

  if (!body) {
    throw new Error("Réponse de session invalide.");
  }

  return body;
}
