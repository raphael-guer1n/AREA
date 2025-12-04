import Link from "next/link";
import { redirect } from "next/navigation";

import { exchangeOAuthCallback } from "@/lib/api/auth";
import { setSessionCookie } from "@/lib/session";

type CallbackPageProps = {
  searchParams?: {
    code?: string | string[];
    state?: string | string[];
    error?: string;
  };
};

export const dynamic = "force-dynamic";

export default async function AuthCallbackPage({ searchParams }: CallbackPageProps) {
  const codeParam = searchParams?.code;
  const stateParam = searchParams?.state;
  const errorParam = searchParams?.error;

  if (errorParam) {
    return (
      <CallbackError
        title="Connexion annulée"
        message="Le fournisseur OAuth2 a renvoyé une erreur. Merci de réessayer."
      />
    );
  }

  if (!codeParam || Array.isArray(codeParam)) {
    return (
      <CallbackError
        title="Code absent"
        message="Le paramètre 'code' est manquant dans l'URL de retour."
      />
    );
  }

  if (!stateParam || Array.isArray(stateParam)) {
    return (
      <CallbackError
        title="State absent"
        message="Le paramètre 'state' est manquant dans l'URL de retour."
      />
    );
  }

  try {
    const { access_token } = await exchangeOAuthCallback({
      code: codeParam,
      state: stateParam,
    });
    await setSessionCookie(access_token);
    redirect("/dashboard");
  } catch (error) {
    const message =
      error instanceof Error
        ? error.message
        : "Impossible de finaliser la connexion OAuth2.";

    return (
      <CallbackError
        title="Échec de la connexion"
        message={message}
      />
    );
  }
}

function CallbackError({ title, message }: { title: string; message: string }) {
  return (
    <main className="flex min-h-screen items-center justify-center bg-[var(--surface)] px-4 py-12">
      <div className="w-full max-w-lg rounded-3xl border border-[var(--surface-border)] bg-[var(--background)] px-8 py-10 shadow-[0_20px_60px_rgba(17,42,70,0.08)]">
        <div className="space-y-3 text-center">
          <p className="text-xs font-semibold uppercase tracking-[0.14em] text-[var(--muted)]">
            Authentification OAuth2
          </p>
          <h1 className="text-2xl font-semibold text-[var(--foreground)]">{title}</h1>
          <p className="text-sm text-[var(--muted)]">{message}</p>
          <div className="mt-6">
            <Link
              href="/login"
              className="text-sm font-semibold text-[var(--blue-soft)] transition hover:text-[var(--blue-primary-3)]"
            >
              Retour à la connexion
            </Link>
          </div>
        </div>
      </div>
    </main>
  );
}
