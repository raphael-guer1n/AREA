import Link from "next/link";

import { Card } from "@/components/ui/Card";

type ProfilePageProps = {
  searchParams?: { email?: string; name?: string };
};

export default function ProfilePage({ searchParams }: ProfilePageProps) {
  const email = searchParams?.email ?? "";
  const name =
    searchParams?.name?.trim() ||
    (email ? email.split("@")[0] : "utilisateur");
  const initials = (name || "?").slice(0, 2).toUpperCase();

  return (
    <main className="mx-auto max-w-4xl space-y-6 px-6 py-10">
      <div className="flex items-center justify-between gap-4">
        <div>
          <p className="text-xs uppercase tracking-[0.25em] text-[var(--muted)]">
            Profil
          </p>
          <h1 className="text-3xl font-semibold">
            Bienvenue {name || "utilisateur"}
          </h1>
          <p className="text-sm text-[var(--muted)]">
            Aperçu rapide de vos informations.
          </p>
        </div>
        <div className="flex h-14 w-14 items-center justify-center rounded-full bg-[var(--surface-border)] text-lg font-semibold">
          {initials}
        </div>
      </div>

      <Card title="Informations de connexion">
        {email ? (
          <dl className="space-y-3 text-sm">
            <div>
              <dt className="text-[var(--muted)]">Email</dt>
              <dd className="text-base font-semibold">{email}</dd>
            </div>
          </dl>
        ) : (
          <p className="text-sm text-[var(--muted)]">
            Aucun identifiant fourni. Connectez-vous pour charger vos données.
          </p>
        )}
        <div className="mt-4 flex gap-3">
          <Link href="/login" className="text-sm underline">
            Changer de compte
          </Link>
          <Link href="/areas" className="text-sm underline">
            Voir mes areas
          </Link>
        </div>
      </Card>
    </main>
  );
}
