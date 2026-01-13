/**
 * Handles OAuth2 callback query params, exchanges the code, persists the app session token,
 * then redirects to the intended page. Designed to run on client-only routes.
 */
"use client";

import { useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";

import { exchangeOAuthCallback } from "@/lib/api/auth";
import { persistSessionToken } from "@/lib/api/session";

type CallbackState =
  | { status: "idle" | "processing" | "success"; error?: undefined }
  | { status: "error"; error: string };

type UseOAuthCallbackOptions = {
  enabled?: boolean;
  callbackPath?: string;
};

export function useOAuthCallback(redirectTo = "/area", options: UseOAuthCallbackOptions = {}) {
  const router = useRouter();
  const searchParams = useSearchParams();
  const code = searchParams.get("code");
  const stateParam = searchParams.get("state");
  const errorParam = searchParams.get("error");
  const enabled = options.enabled ?? true;
  const callbackPath = options.callbackPath ?? redirectTo;

  const [callbackState, setCallbackState] = useState<CallbackState>({
    status: "idle",
  });

  const { status } = callbackState;

  useEffect(() => {
    if (!enabled || status !== "idle") return;

    if (errorParam) {
      setCallbackState({
        status: "error",
        error: "The OAuth2 provider returned an error. Please try again.",
      });
      return;
    }

    if (!code) {
      setCallbackState({
        status: "error",
        error: "The 'code' parameter is missing from the callback URL.",
      });
      return;
    }

    if (!stateParam) {
      setCallbackState({
        status: "error",
        error: "The 'state' parameter is missing from the callback URL.",
      });
      return;
    }

    const handleCallback = async () => {
      setCallbackState({ status: "processing" });

      try {
        const callbackUrl =
          typeof window !== "undefined"
            ? `${window.location.origin}${callbackPath}`
            : undefined;

        const { access_token, token } = await exchangeOAuthCallback(
          {
            code,
            state: stateParam,
          },
          callbackUrl ? { callbackUrl } : {},
        );

        // For "connect" flows, backend may not return an app JWT. Only persist when we actually
        // receive one to avoid clobbering the existing session with a provider access token.
        if (token) {
          await persistSessionToken(token);
          // Force a reload so server components re-read the session cookie.
          if (typeof window !== "undefined") {
            window.location.assign(redirectTo);
            return;
          }
        }

        setCallbackState({ status: "success" });
        router.replace(redirectTo);
      } catch (err) {
        const message =
          err instanceof Error
            ? err.message
            : "Unable to finish the OAuth2 login.";
        setCallbackState({ status: "error", error: message });
      }
    };

    void handleCallback();
  }, [code, enabled, errorParam, redirectTo, router, stateParam, status]);

  return callbackState;
}
