"use client";

import { useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";

import { exchangeOAuthCallback } from "@/lib/api/auth";
import { persistSessionToken } from "@/lib/api/session";

type CallbackState =
  | { status: "idle" | "processing" | "success"; error?: undefined }
  | { status: "error"; error: string };

export function useOAuthCallback(redirectTo = "/dashboard") {
  const router = useRouter();
  const searchParams = useSearchParams();
  const code = searchParams.get("code");
  const stateParam = searchParams.get("state");
  const errorParam = searchParams.get("error");

  const [callbackState, setCallbackState] = useState<CallbackState>({
    status: "idle",
  });

  const { status } = callbackState;

  useEffect(() => {
    if (status !== "idle") return;

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
        const { token } = await exchangeOAuthCallback({
          code,
          state: stateParam,
        });

        if (!token) {
          throw new Error("OAuth2 callback did not return a session token.");
        }

        await persistSessionToken(token);
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
  }, [code, errorParam, redirectTo, router, stateParam, status]);

  return callbackState;
}
