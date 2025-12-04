"use client";

import { useMemo } from "react";

import { Button } from "@/components/ui/Button";
import { useAuth } from "@/hooks/useAuth";
import { cn } from "@/lib/helpers";

type LoginWithGoogleProps = {
  label?: string;
  className?: string;
};

export default function LoginWithGoogle({
  label = "Continuer avec Google",
  className,
}: LoginWithGoogleProps) {
  const { startOAuthLogin, isLoading, error } = useAuth();

  const buttonLabel = useMemo(
    () => (isLoading ? "Redirection vers Google..." : label),
    [isLoading, label],
  );

  return (
    <div className={cn("space-y-2", className)}>
      <Button
        type="button"
        variant="secondary"
        disabled={isLoading}
        onClick={() => startOAuthLogin("google")}
        className="w-full rounded-xl border border-[var(--surface-border)] bg-[var(--background)] px-4 py-3 text-sm font-semibold text-[var(--foreground)] hover:border-[var(--blue-primary-3)] hover:shadow-sm focus-visible:outline-[var(--blue-primary-3)]"
      >
        <span className="flex items-center justify-center gap-3">
          <span className="flex h-9 w-9 items-center justify-center rounded-full border border-[var(--surface-border)] bg-white text-sm font-semibold text-[var(--blue-soft)]">
            G
          </span>
          <span>{buttonLabel}</span>
        </span>
      </Button>
      {error ? <p className="text-sm text-red-500">{error}</p> : null}
    </div>
  );
}
