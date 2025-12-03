"use client";

import Link from "next/link";
import { FormEvent, useState } from "react";

import { useAuth } from "@/hooks/useAuth";
import type { LoginPayload } from "@/types/User";

type ButtonState = "idle" | "success" | "error";

export default function LoginForm() {
  const { login, isLoading, error } = useAuth();
  const [credentials, setCredentials] = useState<LoginPayload>({
    email: "",
    password: "",
  });
  const [feedback, setFeedback] = useState<
    { message: string; tone: "success" | "error" } | null
  >(null);
  const [buttonState, setButtonState] = useState<ButtonState>("idle");

  const buttonVariants: Record<ButtonState, string> = {
    idle:
      "bg-[var(--blue-primary-1)] hover:bg-[var(--blue-primary-2)] focus-visible:outline-[var(--blue-primary-2)]",
    success:
      "bg-emerald-600 hover:bg-emerald-500 focus-visible:outline-emerald-500",
    error: "bg-red-600 hover:bg-red-500 focus-visible:outline-red-500",
  };

  const resetFeedback = () => {
    setFeedback(null);
    setButtonState("idle");
  };

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setFeedback(null);
    setButtonState("idle");

    if (!credentials.email || !credentials.password) {
      setFeedback({ message: "Email et mot de passe requis.", tone: "error" });
      setButtonState("error");
      return;
    }

    const user = await login(credentials);
    if (user) {
      setFeedback({ message: "Connexion réussie.", tone: "success" });
      setButtonState("success");
      return;
    }

    setFeedback({ message: "Échec de la connexion.", tone: "error" });
    setButtonState("error");
  };

  return (
    <div className="flex w-full max-w-xl flex-col items-center gap-8">
      <div className="w-full rounded-3xl border border-[var(--surface-border)] bg-[var(--background)] px-8 py-10 shadow-[0_20px_60px_rgba(17,42,70,0.08)] sm:px-10 sm:py-12">
        <div className="mb-8">
          <h1 className="text-2xl font-semibold uppercase tracking-wide text-[var(--foreground)]">
            Connexion
          </h1>
          <div className="mt-2 h-1 w-16 rounded-full bg-[var(--blue-primary-1)]" />
        </div>

        <p className="mb-6 text-sm text-[var(--muted)]">
          Utilisez vos identifiants pour accéder à votre espace AREA.
        </p>

        <form className="space-y-6" onSubmit={handleSubmit}>
          <label
            htmlFor="email"
            className="block text-xs font-semibold uppercase tracking-[0.12em] text-[var(--muted)]"
          >
            Email
            <input
              id="email"
              type="email"
              required
              value={credentials.email}
              onChange={(event) => {
                setCredentials((current) => ({
                  ...current,
                  email: event.target.value,
                }));
                resetFeedback();
              }}
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
              value={credentials.password}
              onChange={(event) => {
                setCredentials((current) => ({
                  ...current,
                  password: event.target.value,
                }));
                resetFeedback();
              }}
              placeholder="••••••••"
              className="mt-2 w-full rounded-xl border border-[var(--surface-border)] bg-[var(--background)] px-4 py-3 text-sm text-[var(--foreground)] placeholder:text-[var(--placeholder)] focus:border-[var(--blue-primary-3)] focus:outline-none focus:ring-2 focus:ring-[var(--blue-primary-3)]/30"
            />
          </label>
          <button
            type="submit"
            disabled={isLoading}
            className={`mt-2 w-full rounded-xl px-4 py-3 text-sm font-semibold text-white transition focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 disabled:cursor-not-allowed disabled:opacity-70 ${buttonVariants[buttonState]}`}
          >
            {isLoading ? "Connexion..." : "Se connecter"}
          </button>
          <div className="space-y-1 text-sm">
            {feedback ? (
              <p
                className={
                  feedback.tone === "success" ? "text-emerald-600" : "text-red-500"
                }
              >
                {feedback.message}
              </p>
            ) : null}
            {error ? <p className="text-red-500">{error}</p> : null}
          </div>
        </form>

        <div className="mt-8 text-center text-xs font-semibold uppercase tracking-[0.14em] text-[var(--muted)]">
          <Link
            href="/register"
            className="text-[var(--blue-soft)] transition hover:text-[var(--blue-primary-3)]"
          >
            Créer un compte
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
