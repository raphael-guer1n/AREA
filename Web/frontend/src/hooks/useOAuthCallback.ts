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
};

export function useOAuthCallback(redirectTo = "/area", options: UseOAuthCallbackOptions = {}) {
  const router = useRouter();
  const searchParams = useSearchParams();
  const code = searchParams.get("code");
  const stateParam = searchParams.get("state");
  const errorParam = searchParams.get("error");
  const enabled = options.enabled ?? true;

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
        const { access_token, token } = await exchangeOAuthCallback({
          code,
          state: stateParam,
        });

        const sessionToken = token ?? access_token;
        if (!sessionToken) {
          throw new Error("The authentication server did not return a session token.");
        }

        await persistSessionToken(sessionToken);
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
