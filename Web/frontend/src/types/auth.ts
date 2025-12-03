export type GoogleAuthUrlResponse = {
  auth_url: string;
};

export type GoogleCallbackPayload = {
  code: string;
};

export type GoogleCallbackResponse = {
  token: string;
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
};
