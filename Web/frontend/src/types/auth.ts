import type { User } from "./User";

export type OAuthAuthorizeResponse = {
  auth_url: string;
  provider: string;
};

export type OAuthCallbackPayload = {
  code: string;
  state: string;
};

export type OAuthCallbackOptions = {
  callbackUrl?: string;
};

export type OAuthUserInfo = {
  id?: string;
  email?: string;
  name?: string;
  username?: string;
  raw_data?: Record<string, unknown>;
};

export type OAuthCallbackResponse = {
  provider: string;
  user_info: OAuthUserInfo;
  access_token: string; // OAuth provider token
  token_type: string;
  expires_in: number;
  token?: string; // JWT issued by our backend when using login-with
};

export type AuthSession = {
  token: string | null;
};

export type AuthStatus =
  | "idle"
  | "loading"
  | "authenticated"
  | "unauthenticated"
  | "error";

export type AuthError = {
  message: string;
  reason?: "network" | "server" | "unknown";
};

export type SessionStatusResponse = {
  authenticated: boolean;
  token: string | null;
  user: User | null;
  error?: string;
};
