import type { LoginPayload, RegisterPayload, User } from "@/types/User";
import type {
  OAuthAuthorizeResponse,
  OAuthCallbackPayload,
  OAuthCallbackResponse,
} from "@/types/auth";

export const BACKEND_BASE_URL =
  process.env.API_BASE_URL ??
  process.env.NEXT_PUBLIC_API_BASE_URL ??
  "http://localhost:8080";

export type BackendUser = {
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

export function mapUser({ user, token }: BackendAuthSuccess["data"]): User {
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
    const response = await fetch("/api/auth/login", {
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
    const response = await fetch("/api/auth/register", {
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

export async function fetchOAuthAuthorizeUrl(
  provider: string,
): Promise<OAuthAuthorizeResponse> {
  const response = await fetch(
    `${BACKEND_BASE_URL}/auth/oauth2/authorize?provider=${encodeURIComponent(provider)}`,
    {
      method: "GET",
      credentials: "include",
      cache: "no-store",
    },
  );

  const body = (await response.json().catch(() => null)) as
    | {
        success?: boolean;
        data?: { auth_url?: string; provider?: string };
        error?: string;
      }
    | null;

  if (!body?.success || !body.data?.auth_url) {
    throw new Error(
      body?.error ?? "Impossible de récupérer l'URL d'autorisation OAuth2.",
    );
  }

  return {
    auth_url: body.data.auth_url,
    provider: body.data.provider ?? provider,
  };
}

export async function exchangeOAuthCallback(
  payload: OAuthCallbackPayload,
): Promise<OAuthCallbackResponse> {
  const url = `${BACKEND_BASE_URL}/auth/oauth2/callback?code=${encodeURIComponent(
    payload.code,
  )}&state=${encodeURIComponent(payload.state)}`;

  const response = await fetch(url, {
    method: "GET",
    credentials: "include",
    cache: "no-store",
  });

  const body = (await response.json().catch(() => null)) as
    | {
        success?: boolean;
        data?: Partial<OAuthCallbackResponse>;
        error?: string;
      }
    | null;

  if (!body?.success || !body.data?.access_token) {
    throw new Error(
      body?.error ?? "Impossible de finaliser l'authentification OAuth2.",
    );
  }

  return {
    provider: body.data.provider ?? "unknown",
    user_info: body.data.user_info ?? {},
    access_token: body.data.access_token,
    token_type: body.data.token_type ?? "Bearer",
    expires_in: body.data.expires_in ?? 0,
  };
}
