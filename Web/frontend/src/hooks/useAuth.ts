"use client";

import { useCallback, useEffect, useState } from "react";

import {
  fetchGoogleAuthUrl,
  loginRequest,
  logoutRequest,
  registerRequest,
} from "@/lib/api/auth";
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

      setSession((current) => ({
        token: isAuthenticated ? current.token : null,
      }));
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

  const startGoogleLogin = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    setStatus("loading");

    try {
      const { auth_url } = await fetchGoogleAuthUrl();
      setStatus("idle");
      window.location.href = auth_url;
    } catch (err) {
      const message =
        err instanceof Error
          ? err.message
          : "Impossible de démarrer la connexion Google.";
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
        setSession({ token: authenticatedUser.token });
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
        setSession({ token: registeredUser.token });
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
    startGoogleLogin,
    refreshSession,
    login,
    register,
    logout,
  };
}
