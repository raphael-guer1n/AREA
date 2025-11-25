import Link from "next/link";

import { Card } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";

const shortcuts = [
  {
    title: "Services",
    description: "Connectez et explorez les services disponibles.",
    href: "/services",
  },
  {
    title: "Areas",
    description: "Gérez vos automatisations actives.",
    href: "/areas",
  },
  {
    title: "Profil",
    description: "Consultez vos informations de compte.",
    href: "/profile",
  },
];

export default function HomePage() {
  return (
    <main className="min-h-screen px-6 py-12">
      <div className="mx-auto flex max-w-5xl flex-col gap-12">
        <header className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div className="space-y-2">
            <p className="text-xs uppercase tracking-[0.25em] text-[var(--muted)]">
              Automatisation
            </p>
            <h1 className="text-4xl font-semibold">AREA</h1>
            <p className="text-sm text-[var(--muted)]">
              Reliez vos services pour créer des actions et réactions sans effort.
            </p>
          </div>
          <div className="flex gap-3">
            <Link href="/login">
              <Button variant="secondary">Se connecter</Button>
            </Link>
            <Link href="/register">
              <Button>Créer un compte</Button>
            </Link>
          </div>
        </header>

        <Card title="Démarrer rapidement" subtitle="Choisissez où aller">
          <div className="grid gap-4 md:grid-cols-3">
            {shortcuts.map((shortcut) => (
              <Link
                key={shortcut.href}
                href={shortcut.href}
                className="group rounded-xl border border-[var(--surface-border)] bg-[var(--background)] p-4 transition hover:-translate-y-0.5 hover:border-[var(--foreground)]"
              >
                <h3 className="text-lg font-semibold group-hover:text-[var(--foreground)]">
                  {shortcut.title}
                </h3>
                <p className="mt-1 text-sm text-[var(--muted)]">
                  {shortcut.description}
                </p>
              </Link>
            ))}
          </div>
        </Card>
      </div>
    </main>
  );
}
