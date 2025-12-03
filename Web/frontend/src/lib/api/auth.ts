import type { LoginPayload, RegisterPayload, User } from "@/types/User";
import type {
  GoogleAuthUrlResponse,
  GoogleCallbackPayload,
  GoogleCallbackResponse,
} from "@/types/auth";

type AuthAction = "login" | "register";
const BACKEND_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

async function postAuth(
  action: AuthAction,
  payload: LoginPayload | RegisterPayload,
): Promise<User> {
  try {
    const response = await fetch("/api/auth", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ action, ...payload }),
    });

    if (!response.ok) {
      const errorBody = (await response.json().catch(() => null)) as
        | { message?: string }
        | null;
      const message =
        errorBody?.message ?? "Impossible de s'authentifier pour le moment.";
      throw new Error(message);
    }

    return (await response.json()) as User;
  } catch (error) {
    if (error instanceof Error) {
      throw error;
    }
    throw new Error("Erreur réseau pendant l'authentification.");
  }
}

export function loginRequest(payload: LoginPayload): Promise<User> {
  return postAuth("login", payload);
}

export function registerRequest(payload: RegisterPayload): Promise<User> {
  return postAuth("register", payload);
}

export async function logoutRequest(): Promise<void> {
  return Promise.resolve();
}

export async function fetchGoogleAuthUrl(): Promise<GoogleAuthUrlResponse> {
  const response = await fetch(`${BACKEND_BASE_URL}/auth/google/url`, {
    method: "GET",
    credentials: "include",
    cache: "no-store",
  });

  if (!response.ok) {
    throw new Error("Impossible de récupérer l'URL de connexion Google.");
  }

  const data = (await response.json()) as Partial<GoogleAuthUrlResponse>;
  if (!data.auth_url) {
    throw new Error("Réponse du serveur invalide (auth_url manquant).");
  }

  return { auth_url: data.auth_url };
}

export async function exchangeGoogleCode(
  payload: GoogleCallbackPayload,
): Promise<GoogleCallbackResponse> {
  const response = await fetch(`${BACKEND_BASE_URL}/auth/google/callback`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
    credentials: "include",
    cache: "no-store",
  });

  if (!response.ok) {
    const errorBody = (await response.json().catch(() => null)) as
      | { message?: string }
      | null;
    const message =
      errorBody?.message ??
      "Impossible d'échanger le code d'autorisation Google.";
    throw new Error(message);
  }

  const data = (await response.json()) as Partial<GoogleCallbackResponse>;
  if (!data.token) {
    throw new Error("Jeton de session manquant dans la réponse du serveur.");
  }

  return { token: data.token };
}
