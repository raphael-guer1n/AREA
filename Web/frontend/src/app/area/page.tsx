import Link from "next/link";

import { AreaCard } from "@/components/area/AreaCard";
import { AreaNavigation } from "@/components/navigation/AreaNavigation";
import { Card } from "@/components/ui/Card";
import { cn } from "@/lib/helpers";

import { mockAreas } from "./mockAreas";

type ServiceIconProps = {
  label: string;
  colorClass?: string;
};

function ServiceIcon({ label, colorClass }: ServiceIconProps) {
  return (
    <span
      className={cn(
        "flex h-8 w-8 items-center justify-center rounded-full bg-white/18 text-[10px] font-semibold uppercase tracking-wide text-white",
        colorClass,
      )}
    >
      {label.slice(0, 2)}
    </span>
  );
}

export default function AreaPage() {
  const totalAreas = mockAreas.length;
  const activeCount = mockAreas.filter((area) => area.active).length;
  const pausedCount = totalAreas - activeCount;

  return (
    <main className="relative flex min-h-screen justify-center overflow-hidden bg-[var(--surface)] px-6 py-12 pt-10 text-[var(--foreground)]">
      <div className="pointer-events-none absolute inset-0 -z-10 opacity-45">
        <div
          className="absolute -left-10 top-0 h-72 w-72 rounded-full bg-[radial-gradient(circle_at_center,var(--accent)_0,transparent_65%)] blur-3xl"
          aria-hidden
        />
        <div
          className="absolute right-4 top-10 h-72 w-72 rounded-full bg-[radial-gradient(circle_at_center,var(--blue-primary-2)_0,transparent_70%)] blur-3xl"
          aria-hidden
        />
        <div
          className="absolute bottom-[-12%] left-1/3 h-64 w-64 rounded-full bg-[radial-gradient(circle_at_center,var(--blue-primary-1)_0,transparent_65%)] blur-3xl"
          aria-hidden
        />
      </div>

      <div className="relative w-full max-w-6xl space-y-8">
        <div className="flex justify-center">
          <AreaNavigation />
        </div>

        <section className="relative isolate overflow-hidden rounded-[26px] border border-[var(--surface-border)] bg-[var(--background)] px-6 py-8 ring-1 ring-[rgba(28,61,99,0.2)]">
          <div className="flex flex-wrap items-center justify-between gap-8">
            <div className="max-w-2xl space-y-3">
              <p className="text-sm font-semibold uppercase tracking-[0.12em] text-[var(--blue-primary-3)]">
                Tableau des areas
              </p>
              <h1 className="text-3xl font-semibold leading-tight text-[var(--foreground)]">
                Vos automatisations personnelles
              </h1>
              <p className="text-base text-[var(--muted)]">
                Crée vos automatisations en connectant vos services et en définissant des déclencheurs et des actions
              </p>
            </div>
            <div className="grid w-full max-w-sm grid-cols-2 gap-4 sm:max-w-xs">
              <div className="rounded-2xl border border-[var(--surface-border)] bg-[var(--background)] px-4 py-3 ring-1 ring-[rgba(28,61,99,0.15)]">
                <p className="text-xs text-[var(--muted)]">Areas actives</p>
                <p className="text-3xl font-semibold text-[var(--blue-primary-2)]">{activeCount}</p>
              </div>
              <div className="rounded-2xl border border-[var(--surface-border)] bg-[var(--background)] px-4 py-3 ring-1 ring-[rgba(28,61,99,0.1)]">
                <p className="text-xs text-[var(--muted)]">Areas crees</p>
                <p className="text-3xl font-semibold text-[var(--blue-primary-2)]">{totalAreas}</p>
              </div>
            </div>
          </div>
        </section>

        <section className="w-full">
          <Card
            title="Vos areas"
            action={
              <Link
                href="/area/create"
                className="inline-flex items-center justify-center gap-2 rounded-full border border-[var(--blue-primary-2)] bg-[var(--blue-primary-2)] px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:border-[var(--blue-primary-3)] hover:bg-[var(--blue-primary-3)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)]"
              >
                Create area
              </Link>
            }
            className="relative w-full overflow-hidden rounded-[26px] border-[var(--surface-border)] bg-[var(--background)] ring-1 ring-[rgba(28,61,99,0.15)]"
          >
            {totalAreas > 0 ? (
              <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                {mockAreas.map((area) => (
                  <AreaCard
                    key={area.id}
                    id={area.id}
                    name={area.name}
                    actionLabel={area.action.label}
                    reactionLabel={area.reaction.label}
                    actionIcon={<ServiceIcon label={area.action.label} colorClass={area.action.colorClass} />}
                    reactionIcon={<ServiceIcon label={area.reaction.label} colorClass={area.reaction.colorClass} />}
                    isActive={area.active}
                    gradientFrom={area.gradient?.from}
                    gradientTo={area.gradient?.to}
                  />
                ))}
              </div>
            ) : (
              <div className="flex flex-col items-center justify-center gap-5 rounded-2xl border border-[var(--surface-border)] bg-[var(--background)] px-6 py-10 text-center">
                <div className="flex h-14 w-14 items-center justify-center rounded-2xl border border-[var(--surface-border)] bg-[var(--surface)]">
                  <div className="h-6 w-6 rounded-full border-2 border-[var(--surface-border)]" />
                </div>
                <div className="space-y-2">
                  <p className="text-lg font-semibold">Pas encore d&apos;area</p>
                  <p className="text-sm text-[var(--muted)]">
                    Dès que vous connecterez vos services et définirez un déclencheur, vos areas apparaitront ici avec
                    leur statut en temps réel.
                  </p>
                </div>
                <Link
                  href="/area/create"
                  className="inline-flex items-center justify-center gap-2 rounded-full border border-[var(--blue-primary-2)] bg-[var(--blue-primary-2)] px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:border-[var(--blue-primary-3)] hover:bg-[var(--blue-primary-3)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)]"
                >
                  Créer une area
                </Link>
              </div>
            )}
          </Card>
        </section>
      </div>
    </main>
  );
}
