"use client";

import { AreaNavigation } from "@/components/navigation/AreaNavigation";
import { useAuth } from "@/hooks/useAuth";

export default function ProfilPage() {
  const { user, token } = useAuth();

  return (
    <main className="flex min-h-screen flex-col items-center gap-10 bg-[var(--surface)] px-6 py-12 pt-10 text-[var(--foreground)]">
      <AreaNavigation />

      <section className="w-full max-w-xl rounded-[26px] border border-[var(--surface-border)] bg-[var(--background)] p-6 shadow-[0_10px_35px_rgba(9,18,36,0.18)]">
        <div className="space-y-4">
          <div>
            <p className="text-sm font-semibold uppercase tracking-[0.18em] text-[var(--blue-primary-2)]">
              Profil
            </p>
            <h1 className="text-2xl font-semibold text-[var(--foreground)]">Informations du compte</h1>
          </div>

          {user ? (
            <dl className="space-y-3 text-sm">
              <div>
                <dt className="text-[var(--muted)]">Nom</dt>
                <dd className="text-base font-medium text-[var(--foreground)]">{user.name ?? user.username}</dd>
              </div>
              <div>
                <dt className="text-[var(--muted)]">Email</dt>
                <dd className="text-base font-medium text-[var(--foreground)]">{user.email}</dd>
              </div>
              <div>
                <dt className="text-[var(--muted)]">Auth token</dt>
                <dd className="mt-1 rounded-xl border border-[var(--surface-border)] bg-[var(--surface)] px-3 py-2 font-mono text-xs break-all text-[var(--foreground)]">
                  {token ?? user.token ?? "Non disponible"}
                </dd>
              </div>
            </dl>
          ) : (
            <p className="text-sm text-[var(--muted)]">Aucune information disponible.</p>
          )}
        </div>
      </section>
    </main>
  );
}
