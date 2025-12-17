"use client";

import { AreaNavigation } from "@/components/navigation/AreaNavigation";
import { Card } from "@/components/ui/AreaCard";
import { useAuth } from "@/hooks/useAuth";
import type { AuthSession } from "@/types/auth";
import type { User } from "@/types/User";

type ProfileClientProps = {
  initialUser: User | null;
  initialSession: AuthSession | null;
};

export function ProfileClient({ initialUser, initialSession }: ProfileClientProps) {
  const { user, token, status, refreshSession } = useAuth({
    initialUser,
    initialSession,
  });

  const displayName = user?.name || user?.username || "Unknown user";
  const displayEmail = user?.email || "Email unavailable";
  const displayToken = token || "No token in session";

  return (
    <Card
      title="Profil"
      action={
        <button
          type="button"
          onClick={() => refreshSession()}
          className="rounded-full border border-[var(--surface-border)] bg-[var(--surface)] px-3 py-1 text-xs font-semibold text-[var(--foreground)] transition hover:border-[var(--blue-primary-3)] hover:text-[var(--blue-primary-3)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)]"
        >
          Rafraîchir
        </button>
      }
      className="rounded-[26px] border-[var(--surface-border)] bg-[var(--background)] ring-1 ring-[rgba(28,61,99,0.15)]"
    >
      <div className="space-y-4 text-sm">
        <div className="flex items-center gap-3">
          <span className="inline-flex h-10 w-10 items-center justify-center rounded-full bg-[var(--blue-primary-2)] text-sm font-semibold text-white">
            {displayName.slice(0, 2).toUpperCase()}
          </span>
          <div>
            <p className="text-base font-semibold text-[var(--foreground)]">{displayName}</p>
            <p className="text-[var(--muted)]">{displayEmail}</p>
            <p className="text-xs text-[var(--muted)]">Status: {status}</p>
          </div>
        </div>

        <div className="rounded-xl border border-[var(--surface-border)] bg-[var(--surface)] px-4 py-3">
          <p className="text-xs font-semibold uppercase tracking-[0.12em] text-[var(--muted)]">
            Session Token
          </p>
          <p className="mt-2 break-words font-mono text-xs text-[var(--foreground)]">{displayToken}</p>
        </div>
      </div>
    </Card>
  );
}
