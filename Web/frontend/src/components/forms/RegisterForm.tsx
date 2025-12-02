"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { FormEvent, useState } from "react";

import { useAuth } from "@/hooks/useAuth";
import type { RegisterPayload } from "@/types/User";

const socialProviders = [
  { key: "google", label: "Continuer avec Google", badge: "G" },
  { key: "apple", label: "Continuer avec Apple", badge: "A" },
  { key: "facebook", label: "Continuer avec Facebook", badge: "F" },
] as const;

export default function RegisterForm() {
  const router = useRouter();
  const { register, isLoading, error } = useAuth();
  const [payload, setPayload] = useState<RegisterPayload>({
    name: "",
    email: "",
    password: "",
  });
  const [confirmPassword, setConfirmPassword] = useState("");
  const [status, setStatus] = useState<string | null>(null);
  const [localError, setLocalError] = useState<string | null>(null);

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setStatus(null);
    setLocalError(null);

    if (payload.password !== confirmPassword) {
      setLocalError("Les mots de passe ne correspondent pas.");
      return;
    }

    const user = await register(payload);
    if (user) {
      setStatus("Compte créé. Vous pouvez vous connecter.");
      router.push("/login");
    }
  };

  return (
    <div className="flex w-full max-w-xl flex-col items-center gap-8">
      <div className="w-full rounded-3xl border border-[var(--surface-border)] bg-[var(--background)] px-8 py-10 shadow-[0_20px_60px_rgba(17,42,70,0.08)] sm:px-10 sm:py-12">
        <div className="mb-8">
          <h1 className="text-2xl font-semibold uppercase tracking-wide text-[var(--foreground)]">
            Inscription
          </h1>
          <div className="mt-2 h-1 w-16 rounded-full bg-[var(--blue-primary-1)]" />
        </div>

        <div className="space-y-3">
          {socialProviders.map((provider) => (
            <button
              key={provider.key}
              type="button"
              className="w-full rounded-xl border border-[var(--surface-border)] bg-[var(--background)] px-4 py-3 text-sm font-semibold text-[var(--foreground)] transition hover:border-[var(--blue-primary-3)] hover:shadow-sm focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)]"
              aria-label={provider.label}
            >
              <span className="flex items-center justify-center gap-3">
                <span className="flex h-9 w-9 items-center justify-center rounded-full border border-[var(--surface-border)] bg-[var(--surface)] text-sm font-semibold text-[var(--blue-soft)]">
                  {provider.badge}
                </span>
                <span>{provider.label}</span>
              </span>
            </button>
          ))}
        </div>

        <div className="my-8 flex items-center gap-3 text-[11px] font-semibold uppercase tracking-[0.22em] text-[var(--muted)]">
          <span className="h-px flex-1 bg-[var(--surface-border)]" aria-hidden />
          <span>Ou par email</span>
          <span className="h-px flex-1 bg-[var(--surface-border)]" aria-hidden />
        </div>

        <form className="space-y-6" onSubmit={handleSubmit}>
          <label
            htmlFor="name"
            className="block text-xs font-semibold uppercase tracking-[0.12em] text-[var(--muted)]"
          >
            Nom complet
            <input
              id="name"
              type="text"
              value={payload.name ?? ""}
              onChange={(event) =>
                setPayload((current) => ({
                  ...current,
                  name: event.target.value,
                }))
              }
              placeholder="Jean Dupont"
              className="mt-2 w-full rounded-xl border border-[var(--surface-border)] bg-[var(--background)] px-4 py-3 text-sm text-[var(--foreground)] placeholder:text-[var(--placeholder)] focus:border-[var(--blue-primary-3)] focus:outline-none focus:ring-2 focus:ring-[var(--blue-primary-3)]/30"
            />
          </label>
          <label
            htmlFor="email"
            className="block text-xs font-semibold uppercase tracking-[0.12em] text-[var(--muted)]"
          >
            Email
            <input
              id="email"
              type="email"
              required
              value={payload.email}
              onChange={(event) =>
                setPayload((current) => ({
                  ...current,
                  email: event.target.value,
                }))
              }
              placeholder="votre@email.com"
              className="mt-2 w-full rounded-xl border border-[var(--surface-border)] bg-[var(--background)] px-4 py-3 text-sm text-[var(--foreground)] placeholder:text-[var(--placeholder)] focus:border-[var(--blue-primary-3)] focus:outline-none focus:ring-2 focus:ring-[var(--blue-primary-3)]/30"
            />
          </label>
          <label
            htmlFor="password"
            className="block text-xs font-semibold uppercase tracking-[0.12em] text-[var(--muted)]"
          >
            Mot de passe
            <input
              id="password"
              type="password"
              required
              value={payload.password}
              onChange={(event) =>
                setPayload((current) => ({
                  ...current,
                  password: event.target.value,
                }))
              }
              placeholder="••••••••"
              className="mt-2 w-full rounded-xl border border-[var(--surface-border)] bg-[var(--background)] px-4 py-3 text-sm text-[var(--foreground)] placeholder:text-[var(--placeholder)] focus:border-[var(--blue-primary-3)] focus:outline-none focus:ring-2 focus:ring-[var(--blue-primary-3)]/30"
            />
          </label>
          <label
            htmlFor="confirmPassword"
            className="block text-xs font-semibold uppercase tracking-[0.12em] text-[var(--muted)]"
          >
            Confirmer le mot de passe
            <input
              id="confirmPassword"
              type="password"
              required
              value={confirmPassword}
              onChange={(event) => setConfirmPassword(event.target.value)}
              placeholder="••••••••"
              className="mt-2 w-full rounded-xl border border-[var(--surface-border)] bg-[var(--background)] px-4 py-3 text-sm text-[var(--foreground)] placeholder:text-[var(--placeholder)] focus:border-[var(--blue-primary-3)] focus:outline-none focus:ring-2 focus:ring-[var(--blue-primary-3)]/30"
            />
          </label>
          <button
            type="submit"
            disabled={isLoading}
            className="mt-2 w-full rounded-xl bg-[var(--blue-primary-1)] px-4 py-3 text-sm font-semibold text-white transition hover:bg-[var(--blue-primary-2)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-2)] disabled:cursor-not-allowed disabled:opacity-70"
          >
            {isLoading ? "Création..." : "Créer mon compte"}
          </button>
          <div className="space-y-1 text-sm">
            {localError ? <p className="text-red-500">{localError}</p> : null}
            {error ? <p className="text-red-500">{error}</p> : null}
            {status ? <p className="text-emerald-600">{status}</p> : null}
          </div>
        </form>

        <div className="mt-3 text-center text-xs font-semibold uppercase tracking-[0.14em] text-[var(--muted)]">
          <Link
            href="/login"
            className="text-[var(--blue-soft)] transition hover:text-[var(--blue-primary-3)]"
          >
            Déjà un compte ? Se connecter
          </Link>
        </div>
      </div>

      <Link
        href="/"
        className="inline-flex items-center gap-2 text-sm font-medium text-[var(--muted)] transition hover:text-[var(--foreground)]"
      >
        <svg
          aria-hidden="true"
          className="h-4 w-4"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
        >
          <path d="M15 18l-6-6 6-6" />
        </svg>
        Retour
      </Link>
    </div>
  );
}
