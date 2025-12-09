"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { FormEvent, useState } from "react";

import { useAuth } from "@/hooks/useAuth";
import type { LoginPayload } from "@/types/User";
import { z } from "zod";

const loginSchema = z.object({
  email: z.string().trim().email("Please enter a valid email address."),
  password: z.string().min(8, "Password must be at least 8 characters."),
});

type ButtonState = "idle" | "success" | "error";
type SocialProvider = {
  key: "google";
  label: string;
  badge: string;
  onClick?: () => void;
};

export default function LoginForm() {
  const router = useRouter();
  const { login, startOAuthLogin, isLoading } = useAuth();
  const [credentials, setCredentials] = useState<LoginPayload>({
    email: "",
    password: "",
  });
  const [buttonState, setButtonState] = useState<ButtonState>("idle");

  const socialProviders: SocialProvider[] = [
    {
      key: "google",
      label: "Continue with Google",
      badge: "G",
      onClick: () => startOAuthLogin("google"),
    },
  ];

  const buttonVariants: Record<ButtonState, string> = {
    idle:
      "bg-[var(--blue-primary-1)] hover:bg-[var(--blue-primary-2)] focus-visible:outline-[var(--blue-primary-2)]",
    success:
      "bg-emerald-600 hover:bg-emerald-500 focus-visible:outline-emerald-500",
    error: "bg-red-600 hover:bg-red-500 focus-visible:outline-red-500",
  };

  const resetFeedback = () => {
    setButtonState("idle");
  };

  const validateCredentials = () => {
    const result = loginSchema.safeParse(credentials);
    if (!result.success) {
      setButtonState("error");
      return null;
    }
    return result.data;
  };

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setButtonState("idle");

    const validated = validateCredentials();
    if (!validated) return;

    const user = await login(validated);
    if (user) {
      setButtonState("success");
      router.push("/area");
      return;
    }

    setButtonState("error");
  };

  return (
    <div className="flex w-full max-w-xl flex-col items-center gap-8">
      <div className="w-full rounded-3xl border border-[var(--surface-border)] bg-[var(--background)] px-8 py-10 shadow-[0_20px_60px_rgba(17,42,70,0.08)] sm:px-10 sm:py-12">
        <div className="mb-8">
          <h1 className="text-2xl font-semibold uppercase tracking-wide text-[var(--foreground)]">
            Login
          </h1>
          <div className="mt-2 h-1 w-16 rounded-full bg-[var(--blue-primary-1)]" />
        </div>

        <div className="space-y-3">
          {socialProviders.map((provider) => (
            <button
              key={provider.key}
              type="button"
              onClick={() => provider.onClick?.()}
              disabled={isLoading && provider.key === "google"}
              className="w-full rounded-xl border border-[var(--surface-border)] bg-[var(--background)] px-4 py-3 text-sm font-semibold text-[var(--foreground)] transition hover:border-[var(--blue-primary-3)] hover:shadow-sm focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)] disabled:cursor-not-allowed disabled:opacity-70"
              aria-label={provider.label}
            >
              <span className="flex items-center justify-center gap-3">
                <span className="flex h-9 w-9 items-center justify-center rounded-full border border-[var(--surface-border)] bg-[var(--surface)] text-sm font-semibold text-[var(--blue-soft)]">
                  {provider.badge}
                </span>
                <span>
                  {isLoading && provider.key === "google"
                    ? "Redirecting..."
                    : provider.label}
                </span>
              </span>
            </button>
          ))}
        </div>

        <div className="my-8 flex items-center gap-3 text-[11px] font-semibold uppercase tracking-[0.22em] text-[var(--muted)]">
          <span className="h-px flex-1 bg-[var(--surface-border)]" aria-hidden />
          <span>Or with email</span>
          <span className="h-px flex-1 bg-[var(--surface-border)]" aria-hidden />
        </div>

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
              placeholder="you@example.com"
              className="mt-2 w-full rounded-xl border border-[var(--surface-border)] bg-[var(--background)] px-4 py-3 text-sm text-[var(--foreground)] placeholder:text-[var(--placeholder)] focus:border-[var(--blue-primary-3)] focus:outline-none focus:ring-2 focus:ring-[var(--blue-primary-3)]/30"
            />
          </label>
          <label
            htmlFor="password"
            className="block text-xs font-semibold uppercase tracking-[0.12em] text-[var(--muted)]"
          >
            Password
            <input
              id="password"
              type="password"
              minLength={8}
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
            {isLoading ? "Logging in..." : "Log in"}
          </button>
        </form>

        <div className="mt-1 text-center text-xs font-semibold uppercase tracking-[0.14em] text-[var(--muted)]">
          <Link
            href="/register"
            className="text-[var(--blue-soft)] transition hover:text-[var(--blue-primary-3)]"
          >
            Create an account
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
        Back
      </Link>
    </div>
  );
}
