"use client";

import { useCallback, useEffect, useState } from "react";

import { fetchOAuthAuthorizeUrl, loginRequest, logoutRequest, registerRequest } from "@/lib/api/auth";
import type { AuthSession, AuthStatus, SessionStatusResponse } from "@/types/auth";
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
      await fetch("/api/session", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ token }),
      });
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Impossible de sauvegarder la session.";
      setError(message);
    }
  }, []);

  const refreshSession = useCallback(async () => {
    setStatus("loading");
    setError(null);

    try {
      const response = await fetch("/api/session", {
        credentials: "include",
        cache: "no-store",
      });

      if (!response.ok) {
        throw new Error("Impossible de récupérer la session.");
      }

      const data = (await response.json()) as SessionStatusResponse;
      const isAuthenticated = Boolean(data.authenticated);

      setSession({
        token: isAuthenticated ? data.token : null,
      });
      setStatus(isAuthenticated ? "authenticated" : "unauthenticated");

      return isAuthenticated;
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Impossible de récupérer la session.";
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
      const { auth_url } = await fetchOAuthAuthorizeUrl(provider);
      setStatus("idle");
      window.location.href = auth_url;
    } catch (err) {
      const message =
        err instanceof Error
          ? err.message
          : "Impossible de démarrer la connexion OAuth2.";
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
        err instanceof Error ? err.message : "Impossible de se connecter";
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
        err instanceof Error ? err.message : "Impossible de créer le compte";
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
        fetch("/api/session", { method: "DELETE", credentials: "include" }),
      ]);
      setSession({ token: null });
      setUser(null);
      setStatus("unauthenticated");
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Impossible de se déconnecter";
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
