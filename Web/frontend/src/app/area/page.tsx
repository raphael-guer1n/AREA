import Link from "next/link";

import { AreaNavigation } from "@/components/navigation/AreaNavigation";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";

export default function AreaPage() {
  return (
    <main className="relative flex min-h-screen justify-center overflow-hidden bg-[var(--surface)] px-6 py-12 pt-10 text-[var(--foreground)]">
      <div className="pointer-events-none absolute inset-0 -z-10 opacity-45">
        <div
          className="absolute -left-10 top-0 h-72 w-72 rounded-full bg-[radial-gradient(circle_at_center,var(--accent)_0,transparent_65%)] blur-3xl"
          aria-hidden
        />
        <div
          className="absolute right-4 top-10 h-72 w-72 rounded-full bg-[radial-gradient(circle_at_center,var(--surface-border)_0,transparent_70%)] blur-3xl"
          aria-hidden
        />
        <div
          className="absolute bottom-[-12%] left-1/3 h-64 w-64 rounded-full bg-[radial-gradient(circle_at_center,var(--surface)_0,transparent_65%)] blur-3xl"
          aria-hidden
        />
      </div>

      <div className="relative w-full max-w-6xl space-y-8">
        <div className="flex justify-center">
          <AreaNavigation />
        </div>

        <section className="relative isolate overflow-hidden rounded-[26px] border border-[var(--surface-border)] bg-white px-6 py-8">
          <div className="flex flex-wrap items-center justify-between gap-8">
            <div className="max-w-2xl space-y-3">
              <p className="text-sm font-semibold uppercase tracking-[0.12em] text-[var(--muted)]">
                Tableau des areas
              </p>
              <h1 className="text-3xl font-semibold leading-tight">Vos automatisations personnelles</h1>
              <p className="text-base text-[var(--muted)]">
                Crée vos automatisations en connectant vos services et en définissant des déclencheurs et des actions
              </p>
            </div>
            <div className="grid w-full max-w-sm grid-cols-2 gap-4 sm:max-w-xs">
              <div className="rounded-2xl border border-[var(--surface-border)] bg-white px-4 py-3">
                <p className="text-xs text-[var(--muted)]">Areas actives</p>
                <p className="text-3xl font-semibold">0</p>
              </div>
              <div className="rounded-2xl border border-[var(--surface-border)] bg-white px-4 py-3">
                <p className="text-xs text-[var(--muted)]">Areas crees</p>
                <p className="text-3xl font-semibold">0</p>
              </div>
            </div>
          </div>
        </section>

        <section className="w-full">
          <Card
            title="Vos areas"
            subtitle="Elles s afficheront ici des qu elles seront creees."
            className="relative w-full overflow-hidden border-[var(--surface-border)] bg-white border-dashed"
          >
            <div className="flex flex-col items-center justify-center gap-5 rounded-2xl border border-dashed border-[var(--surface-border)] bg-white px-6 py-12 text-center">
              <div className="flex h-14 w-14 items-center justify-center rounded-2xl border border-[var(--surface-border)] bg-white">
                <div className="h-6 w-6 rounded-full border-2 border-[var(--surface-border)]" />
              </div>
              <div className="space-y-2">
                <p className="text-lg font-semibold">Pas encore d area</p>
                <p className="text-sm text-[var(--muted)]">
                  Des que vous connecterez vos services et definirez un declencheur, vos areas apparaitront ici avec
                  leur statut en temps reel.
                </p>
              </div>
              <div className="flex flex-wrap justify-center gap-3">
                <Link href="/services">
                  <Button>Commencer avec un service</Button>
                </Link>
                <Button variant="secondary" type="button" className="cursor-not-allowed opacity-70" disabled>
                  Creer une area (bientot)
                </Button>
              </div>
            </div>
          </Card>
        </section>
      </div>
    </main>
  );
}
