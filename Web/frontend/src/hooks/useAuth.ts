"use client";

import { useCallback, useEffect, useState } from "react";

import { fetchOAuthAuthorizeUrl, loginRequest, logoutRequest, registerRequest } from "@/lib/api/auth";
import { clearSession, fetchSessionStatus, persistSessionToken } from "@/lib/api/session";
import type { AuthSession, AuthStatus } from "@/types/auth";
import type { LoginPayload, RegisterPayload, User } from "@/types/User";

type UseAuthOptions = {
  initialUser?: User | null;
  initialSession?: AuthSession | null;
};

export function useAuth(options: UseAuthOptions = {}) {
  const { initialUser = null, initialSession = null } = options;

  const [user, setUser] = useState<User | null>(initialUser);
  const [session, setSession] = useState<AuthSession>({
    token: initialSession?.token ?? null,
  });
  const [status, setStatus] = useState<AuthStatus>(
    initialSession?.token ? "authenticated" : "idle",
  );
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const persistSession = useCallback(async (token: string) => {
    setSession({ token });
    try {
      await persistSessionToken(token);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Unable to save the session.";
      setError(message);
    }
  }, []);

  const refreshSession = useCallback(async () => {
    setStatus("loading");
    setError(null);

    try {
      const data = await fetchSessionStatus();
      const isAuthenticated = Boolean(data.authenticated && data.token);

      if (!isAuthenticated) {
        setSession({ token: null });
        setUser(null);
        setStatus("unauthenticated");
        if (data.error) setError(data.error);
        return false;
      }

      setSession({ token: isAuthenticated ? data.token : null });
      setUser(isAuthenticated ? data.user : null);
      setStatus(isAuthenticated ? "authenticated" : "unauthenticated");

      if (!isAuthenticated && data.error) {
        setError(data.error);
      }

      return isAuthenticated;
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Unable to retrieve the session.";
      setError(message);
      setStatus("error");
      return false;
    }
  }, []);

  useEffect(() => {
    if (initialSession?.token) return;
    void refreshSession();
  }, [initialSession?.token, refreshSession]);

  const startOAuthLogin = useCallback(async (provider: string) => {
    setIsLoading(true);
    setError(null);
    setStatus("loading");

    try {
      const callbackUrl =
        typeof window !== "undefined"
          ? `${window.location.origin}/area`
          : undefined;

      const { auth_url } = await fetchOAuthAuthorizeUrl(provider, {
        mode: "login",
        platform: "web",
        callbackUrl,
      });
      setStatus("idle");
      window.location.href = auth_url;
    } catch (err) {
      const message =
        err instanceof Error
          ? err.message
          : "Unable to start the OAuth2 login.";
      setError(message);
      setStatus("error");
    } finally {
      setIsLoading(false);
    }
  }, []);

  const login = useCallback(async (payload: LoginPayload) => {
    setIsLoading(true);
    setError(null);
    setStatus("loading");

    try {
      const authenticatedUser = await loginRequest(payload);
      if (authenticatedUser.token) {
        await persistSession(authenticatedUser.token);
      }
      setStatus("authenticated");
      setUser(authenticatedUser);
      return authenticatedUser;
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Unable to log in.";
      setError(message);
      setStatus("error");
      return null;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const register = useCallback(async (payload: RegisterPayload) => {
    setIsLoading(true);
    setError(null);
    setStatus("loading");

    try {
      const registeredUser = await registerRequest(payload);
      if (registeredUser.token) {
        await persistSession(registeredUser.token);
      }
      setUser(registeredUser);
      setStatus("authenticated");
      return registeredUser;
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Unable to create the account.";
      setError(message);
      setStatus("error");
      return null;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const logout = useCallback(async () => {
    setIsLoading(true);
    setError(null);

    try {
      await Promise.all([
        logoutRequest(),
        clearSession(),
      ]);
      setSession({ token: null });
      setUser(null);
      setStatus("unauthenticated");
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Unable to log out.";
      setError(message);
      setStatus("error");
    } finally {
      setIsLoading(false);
    }
  }, []);

  return {
    user,
    session,
    token: session.token,
    status,
    isLoading,
    error,
    startOAuthLogin,
    refreshSession,
    login,
    register,
    logout,
  };
}
