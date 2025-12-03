import type { LoginPayload, RegisterPayload, User } from "@/types/User";
import type {
  GoogleAuthUrlResponse,
  GoogleCallbackPayload,
  GoogleCallbackResponse,
} from "@/types/auth";

const BACKEND_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

type BackendUser = {
  id: number | string;
  email: string;
  username: string;
  created_at?: string;
  updated_at?: string;
};

type BackendAuthSuccess = {
  success: true;
  data: {
    user: BackendUser;
    token: string;
  };
};

type BackendAuthError = {
  success: false;
  error: string;
};

type BackendAuthResponse = BackendAuthSuccess | BackendAuthError;

function mapUser({ user, token }: BackendAuthSuccess["data"]): User {
  return {
    id: String(user.id),
    email: user.email,
    username: user.username,
    name: user.username,
    token,
  };
}

async function handleAuthResponse(response: Response): Promise<User> {
  const body = (await response.json().catch(() => null)) as BackendAuthResponse | null;

  if (!body) {
    throw new Error("Réponse du serveur invalide.");
  }

  if (!body.success) {
    throw new Error(body.error || "Impossible de s'authentifier.");
  }

  return mapUser(body.data);
}

export async function loginRequest(payload: LoginPayload): Promise<User> {
  try {
    const response = await fetch(`${BACKEND_BASE_URL}/auth/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        emailOrUsername: payload.email.trim(),
        password: payload.password.trim(),
      }),
    });
    return handleAuthResponse(response);
  } catch (error) {
    if (error instanceof Error) throw error;
    throw new Error("Impossible de se connecter pour le moment.");
  }
}

export async function registerRequest(payload: RegisterPayload): Promise<User> {
  try {
    const username = (payload.name ?? payload.email).trim();
    const response = await fetch(`${BACKEND_BASE_URL}/auth/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        email: payload.email.trim(),
        username,
        password: payload.password.trim(),
      }),
    });
    return handleAuthResponse(response);
  } catch (error) {
    if (error instanceof Error) throw error;
    throw new Error("Impossible de créer le compte pour le moment.");
  }
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
