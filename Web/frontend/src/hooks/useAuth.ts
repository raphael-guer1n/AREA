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
    // If we have both token and user preloaded, skip refresh; otherwise try to refresh.
    if (initialSession?.token && initialUser) return;
    void refreshSession();
  }, [initialSession?.token, initialUser, refreshSession]);

  const resolveCallbackUrl = useCallback((path: string) => {
    const base =
      process.env.NEXT_PUBLIC_SITE_URL ??
      process.env.NEXT_PUBLIC_OAUTH_CALLBACK_BASE ??
      (typeof window !== "undefined" ? window.location.origin : "");
    if (!base) return undefined;
    return `${base.replace(/\/+$/, "")}${path}`;
  }, []);

  const startOAuthLogin = useCallback(async (provider: string) => {
    setIsLoading(true);
    setError(null);
    setStatus("loading");

    try {
      const callbackUrl = resolveCallbackUrl("/area");

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
  }, [resolveCallbackUrl]);

  const startOAuthConnect = useCallback(
    async (provider: string) => {
      if (!session.token) {
        setError("Vous devez être connecté pour lier un service.");
        setStatus("unauthenticated");
        return;
      }

      setIsLoading(true);
      setError(null);
      setStatus("loading");

      try {
        const callbackUrl = resolveCallbackUrl("/services");

        const { auth_url } = await fetchOAuthAuthorizeUrl(provider, {
          mode: "connect",
          platform: "web",
          callbackUrl,
          token: session.token,
        });
        setStatus("idle");
        window.location.href = auth_url;
      } catch (err) {
        const message =
          err instanceof Error
            ? err.message
            : "Impossible de démarrer la connexion du service.";
        setError(message);
        setStatus("error");
        throw err instanceof Error ? err : new Error(message);
      } finally {
        setIsLoading(false);
      }
    },
    [resolveCallbackUrl, session.token],
  );

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
    startOAuthConnect,
    refreshSession,
    login,
    register,
    logout,
  };
}
