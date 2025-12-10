"use client";

import Link from "next/link";

import { useOAuthCallback } from "@/hooks/useOAuthCallback";

export const dynamic = "force-dynamic";

export default function AuthCallbackPage() {
  const { status, error } = useOAuthCallback("/area");
  const isProcessing = status === "idle" || status === "processing" || status === "success";

  if (status === "error") {
    return (
      <CallbackError
        title="Sign-in failed"
        message={error ?? "Unable to finish the OAuth2 login."}
      />
    );
  }

  return (
    <main className="flex min-h-screen items-center justify-center bg-[var(--surface)] px-4 py-12">
      <div className="w-full max-w-lg rounded-3xl border border-[var(--surface-border)] bg-[var(--background)] px-8 py-10 shadow-[0_20px_60px_rgba(17,42,70,0.08)]">
        <div className="space-y-3 text-center">
          <p className="text-xs font-semibold uppercase tracking-[0.14em] text-[var(--muted)]">
            OAuth2 Authentication
          </p>
          <h1 className="text-2xl font-semibold text-[var(--foreground)]">
            {isProcessing ? "Signing you in..." : "Redirecting"}
          </h1>
          <p className="text-sm text-[var(--muted)]">
            {isProcessing
              ? "Please wait while we validate your account."
              : "Redirecting..."}
          </p>
        </div>
      </div>
    </main>
  );
}

function CallbackError({ title, message }: { title: string; message: string }) {
  return (
    <main className="flex min-h-screen items-center justify-center bg-[var(--surface)] px-4 py-12">
      <div className="w-full max-w-lg rounded-3xl border border-[var(--surface-border)] bg-[var(--background)] px-8 py-10 shadow-[0_20px_60px_rgba(17,42,70,0.08)]">
        <div className="space-y-3 text-center">
          <p className="text-xs font-semibold uppercase tracking-[0.14em] text-[var(--muted)]">
            OAuth2 Authentication
          </p>
          <h1 className="text-2xl font-semibold text-[var(--foreground)]">{title}</h1>
          <p className="text-sm text-[var(--muted)]">{message}</p>
          <div className="mt-6">
            <Link
              href="/login"
              className="text-sm font-semibold text-[var(--blue-soft)] transition hover:text-[var(--blue-primary-3)]"
            >
              Back to login
            </Link>
          </div>
        </div>
      </div>
    </main>
  );
}
