"use client";

import { useCallback, useState } from "react";

import { loginRequest, logoutRequest, registerRequest } from "@/lib/api/auth";
import { clearToken, getStoredToken, storeToken } from "@/lib/auth";
import type { LoginPayload, RegisterPayload, User } from "@/types/User";

export function useAuth(initialUser?: User | null) {
  const [user, setUser] = useState<User | null>(initialUser ?? null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const login = useCallback(async (payload: LoginPayload) => {
    setIsLoading(true);
    setError(null);
    try {
      const authenticatedUser = await loginRequest(payload);
      if (authenticatedUser.token) {
        storeToken(authenticatedUser.token);
      }
      setUser(authenticatedUser);
      return authenticatedUser;
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Impossible de se connecter";
      setError(message);
      return null;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const register = useCallback(async (payload: RegisterPayload) => {
    setIsLoading(true);
    setError(null);
    try {
      const registeredUser = await registerRequest(payload);
      if (registeredUser.token) {
        storeToken(registeredUser.token);
      }
      setUser(registeredUser);
      return registeredUser;
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Impossible de créer le compte";
      setError(message);
      return null;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const logout = useCallback(async () => {
    await logoutRequest();
    clearToken();
    setUser(null);
  }, []);

  return {
    user,
    token: getStoredToken(),
    isLoading,
    error,
    login,
    register,
    logout,
  };
}
