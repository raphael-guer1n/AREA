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

export type BackendAuthResponse = BackendAuthSuccess | BackendAuthError;

type BackendResponse<T> = {
  status: number;
  body: T | null;
};

export function mapUser({ user, token }: BackendAuthSuccess["data"]): User {
  return {
    id: String(user.id),
    email: user.email,
    username: user.username,
    name: user.username,
    token,
  };
}

function parseAuthResponse(body: BackendAuthResponse | null): User {
  if (!body) {
    throw new Error("Invalid server response.");
  }

  if (!body.success) {
    throw new Error(body.error || "Unable to authenticate.");
  }

  return mapUser(body.data);
}

async function postAuth(
  path: string,
  payload: Record<string, unknown>,
): Promise<BackendResponse<BackendAuthResponse>> {
  try {
    const response = await fetch(`${BACKEND_BASE_URL}${path}`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
      cache: "no-store",
    });

    const body = (await response.json().catch(() => null)) as BackendAuthResponse | null;

    return { status: response.status, body };
  } catch (error) {
    throw new Error("Unable to reach the authentication service.");
  }
}

export async function authenticateWithCredentials(
  emailOrUsername: string,
  password: string,
): Promise<BackendResponse<BackendAuthResponse>> {
  return postAuth("/auth/login", { emailOrUsername, password });
}

export async function registerWithCredentials(
  email: string,
  username: string,
  password: string,
): Promise<BackendResponse<BackendAuthResponse>> {
  return postAuth("/auth/register", { email, username, password });
}

export async function fetchAuthenticatedUser(token: string): Promise<User> {
  const response = await fetch(`${BACKEND_BASE_URL}/auth/me`, {
    method: "GET",
    headers: { Authorization: `Bearer ${token}` },
    cache: "no-store",
  });

  const body = (await response.json().catch(() => null)) as
    | {
        success?: boolean;
        data?: { user?: BackendUser };
        error?: string;
      }
    | null;

  if (!body?.success || !body.data?.user) {
    throw new Error(body?.error ?? "Invalid session.");
  }

  return mapUser({ user: body.data.user, token });
}

export async function loginRequest(payload: LoginPayload): Promise<User> {
  const { body } = await authenticateWithCredentials(
    payload.email.trim(),
    payload.password.trim(),
  );
  return parseAuthResponse(body);
}

export async function registerRequest(payload: RegisterPayload): Promise<User> {
  const username = (payload.name ?? payload.email).trim();
  const { body } = await registerWithCredentials(
    payload.email.trim(),
    username,
    payload.password.trim(),
  );
  return parseAuthResponse(body);
}

export async function logoutRequest(): Promise<void> {
  return Promise.resolve();
}

type OAuthAuthorizeMode = "login" | "link";

export async function fetchOAuthAuthorizeUrl(
  provider: string,
  options: {
    token?: string | null;
    mode?: OAuthAuthorizeMode;
    callbackUrl?: string;
    platform?: string;
  } = {},
): Promise<OAuthAuthorizeResponse> {
  const { mode = "login", callbackUrl, platform } = options;

  if (mode === "link") {
    const searchParams = new URLSearchParams({ provider });
    if (callbackUrl) searchParams.set("callback_url", callbackUrl);
    if (platform) searchParams.set("platform", platform);

    const response = await fetch(
      `${BACKEND_BASE_URL}/auth/oauth2/authorize?${searchParams.toString()}`,
      {
        method: "GET",
        credentials: "include",
        cache: "no-store",
        headers: options.token
          ? { Authorization: `Bearer ${options.token}` }
          : undefined,
      },
    );

    const body = (await response.json().catch(() => null)) as
      | {
          success?: boolean;
          data?: { auth_url?: string; provider?: string };
          error?: string;
        }
      | null;

    if (!response.ok || !body?.success || !body.data?.auth_url) {
      throw new Error(
        body?.error ?? "Unable to retrieve the OAuth2 authorization URL.",
      );
    }

    return {
      auth_url: body.data.auth_url,
      provider: body.data.provider ?? provider,
    };
  }

  const response = await fetch(
    `${BACKEND_BASE_URL}/auth/oauth2/login?provider=${encodeURIComponent(provider)}`,
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
      body?.error ?? "Unable to retrieve the OAuth2 authorization URL.",
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
      body?.error ?? "Unable to complete OAuth2 authentication.",
    );
  }

  return {
    provider: body.data.provider ?? "unknown",
    user_info: body.data.user_info ?? {},
    access_token: body.data.access_token,
    token_type: body.data.token_type ?? "Bearer",
    expires_in: body.data.expires_in ?? 0,
    token: body.data.token,
    user: body.data.user,
  };
}
