export type OAuthAuthorizeResponse = {
  auth_url: string;
  provider: string;
};

export type OAuthCallbackPayload = {
  code: string;
  state: string;
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
  access_token: string;
  token_type: string;
  expires_in: number;
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
};
