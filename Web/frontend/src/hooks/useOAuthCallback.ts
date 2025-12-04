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
        error: "Le fournisseur OAuth2 a renvoyé une erreur. Merci de réessayer.",
      });
      return;
    }

    if (!code) {
      setCallbackState({
        status: "error",
        error: "Le paramètre 'code' est manquant dans l'URL de retour.",
      });
      return;
    }

    if (!stateParam) {
      setCallbackState({
        status: "error",
        error: "Le paramètre 'state' est manquant dans l'URL de retour.",
      });
      return;
    }

    const handleCallback = async () => {
      setCallbackState({ status: "processing" });

      try {
        const { access_token } = await exchangeOAuthCallback({
          code,
          state: stateParam,
        });

        await persistSessionToken(access_token);
        setCallbackState({ status: "success" });
        router.replace(redirectTo);
      } catch (err) {
        const message =
          err instanceof Error
            ? err.message
            : "Impossible de finaliser la connexion OAuth2.";
        setCallbackState({ status: "error", error: message });
      }
    };

    void handleCallback();
  }, [code, errorParam, redirectTo, router, stateParam, status]);

  return callbackState;
}
