"use client";

import { useState } from "react";
import Link from "next/link";

import { AreaNavigation } from "@/components/navigation/AreaNavigation";
import { Card } from "@/components/ui/Card";


export default function ServicesPage() {
  const [isConnectModalOpen, setIsConnectModalOpen] = useState(false);

  const openConnectModal = () => setIsConnectModalOpen(true);
  const closeConnectModal = () => setIsConnectModalOpen(false);

  return (
    <main className="relative flex min-h-screen justify-center overflow-hidden bg-[var(--surface)] px-6 py-12 pt-10 text-[var(--foreground)]">
      {isConnectModalOpen ? (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center bg-[rgba(6,14,25,0.55)] px-4 py-10 backdrop-blur-sm"
          onClick={closeConnectModal}
        >
          <div
            className="relative w-full max-w-5xl overflow-hidden rounded-[28px] border border-[var(--surface-border)] bg-[var(--background)] shadow-2xl ring-1 ring-[rgba(28,61,99,0.28)]"
            onClick={(event) => event.stopPropagation()}
            role="dialog"
            aria-modal="true"
            aria-labelledby="connect-service-title"
          >
            <button
              type="button"
              className="absolute right-5 top-5 inline-flex h-10 w-10 items-center justify-center rounded-full border border-[var(--surface-border)] bg-[var(--surface)] text-[var(--muted)] transition hover:text-[var(--foreground)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)]"
              onClick={closeConnectModal}
              aria-label="Fermer la fenêtre"
            >
              <svg className="h-5 w-5" viewBox="0 0 20 20" fill="none" aria-hidden>
                <path
                  d="M15 5 5 15m0-10 10 10"
                  stroke="currentColor"
                  strokeWidth="1.6"
                  strokeLinecap="round"
                />
              </svg>
            </button>

            <div className="space-y-6 px-8 pb-9 pt-7">
              <div className="space-y-2">
                <p className="text-xs font-semibold uppercase tracking-[0.14em] text-[var(--blue-primary-3)]">
                  Connexion de service
                </p>
                <div className="flex flex-wrap items-end justify-between gap-3">
                  <div className="space-y-1">
                    <h2
                      id="connect-service-title"
                      className="text-2xl font-semibold text-[var(--foreground)]"
                    >
                      Connectez un nouveaux service
                    </h2>
                    <p className="text-sm text-[var(--muted)]">
                      Recherchez votre outil et ajoutez-le à votre espace.
                    </p>
                  </div>
    
                </div>
              </div>

              <div className="space-y-5">
                <div className="relative">
                  <span className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-[var(--muted)]">
                    <svg className="h-5 w-5" viewBox="0 0 20 20" fill="none" aria-hidden>
                      <path
                        d="m15.5 15.5-3.5-3.5m-6-3a5 5 0 1 0 10 0 5 5 0 0 0-10 0Z"
                        stroke="currentColor"
                        strokeWidth="1.4"
                        strokeLinecap="round"
                        strokeLinejoin="round"
                      />
                    </svg>
                  </span>
                  <input
                    type="search"
                    placeholder="Rechercher un service (Slack, Gmail, Discord...)"
                    className="w-full rounded-xl border border-[var(--surface-border)] bg-[var(--surface)] px-11 py-3 text-sm text-[var(--foreground)] placeholder:text-[var(--placeholder)] focus:border-[var(--blue-primary-3)] focus:outline-none focus:ring-2 focus:ring-[var(--blue-primary-3)]/25"
                  />
                </div>

                <div className="min-h-[320px] rounded-2xl border-2 border-dashed border-[var(--surface-border)] bg-[var(--surface)] p-8">
                  <div className="mx-auto flex h-full max-w-3xl flex-col items-center justify-center gap-4 text-center">
                    <div className="flex h-16 w-16 items-center justify-center rounded-2xl border border-[var(--surface-border)] bg-[var(--background)] shadow-sm">
                      <svg className="h-7 w-7 text-[var(--muted)]" viewBox="0 0 24 24" fill="none" aria-hidden>
                        <path
                          d="M12 7v10m-5-5h10"
                          stroke="currentColor"
                          strokeWidth="1.6"
                          strokeLinecap="round"
                        />
                        <rect
                          x="4.5"
                          y="4.5"
                          width="15"
                          height="15"
                          rx="4"
                          stroke="currentColor"
                          strokeWidth="1.6"
                          strokeLinecap="round"
                          strokeDasharray="3 3"
                        />
                      </svg>
                    </div>
                    <div className="space-y-1">
                      <p className="text-base font-semibold text-[var(--foreground)]">
                        Les cartes de connexion apparaîtront ici
                      </p>
                      <p className="text-sm text-[var(--muted)]">
                        Ajoutez prochainement les intégrations disponibles pour continuer la configuration.
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      ) : null}

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

        <section className="relative isolate overflow-hidden rounded-[26px] border border-[var(--surface-border)] bg-[var(--background)] px-6 py-7 ring-1 ring-[rgba(28,61,99,0.18)]">
          <div className="flex flex-wrap items-center justify-between gap-4">
            <div className="space-y-2">
              <p className="text-xs font-semibold uppercase tracking-[0.14em] text-[var(--blue-primary-3)]">
                Catalogue des services
              </p>
              <h1 className="text-2xl font-semibold text-[var(--foreground)]">
                Recherchez et filtrez simplement
              </h1>
              <p className="text-sm text-[var(--muted)]">
                Trouvez rapidement un service grâce à la recherche.
              </p>
            </div>
          </div>

          <div className="mt-6 space-y-5">
            <div className="relative">
              <span className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-[var(--muted)]">
                <svg className="h-5 w-5" viewBox="0 0 20 20" fill="none" aria-hidden>
                  <path
                    d="m15.5 15.5-3.5-3.5m-6-3a5 5 0 1 0 10 0 5 5 0 0 0-10 0Z"
                    stroke="currentColor"
                    strokeWidth="1.4"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                </svg>
              </span>
              <input
                type="search"
                placeholder="Rechercher un service (Slack, Gmail, Discord...)"
                className="w-full rounded-xl border border-[var(--surface-border)] bg-[var(--background)] px-11 py-3 text-sm text-[var(--foreground)] placeholder:text-[var(--placeholder)] focus:border-[var(--blue-primary-3)] focus:outline-none focus:ring-2 focus:ring-[var(--blue-primary-3)]/25"
              />
            </div>
          </div>
        </section>

        <section className="w-full">
          <Card
            title="Vos services"
            subtitle="Connectez vos applications pour commencer à créer vos automatisations."
            action={
              <button
                type="button"
                onClick={openConnectModal}
                className="inline-flex items-center justify-center gap-2 rounded-full border border-[var(--blue-primary-2)] bg-[var(--blue-primary-2)] px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:border-[var(--blue-primary-3)] hover:bg-[var(--blue-primary-3)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)]"
              >
                Connecter un service
              </button>
            }
            className="relative w-full overflow-hidden rounded-[26px] border-[var(--surface-border)] bg-[var(--background)] ring-1 ring-[rgba(28,61,99,0.15)]"
          >
            <div className="flex flex-col items-center justify-center gap-5 rounded-2xl border border-[var(--surface-border)] bg-[var(--background)] px-6 py-10 text-center">
              <div className="flex h-14 w-14 items-center justify-center rounded-2xl border border-[var(--surface-border)] bg-[var(--surface)]">
                <div className="h-6 w-6 rounded-full border-2 border-[var(--surface-border)]" />
              </div>
              <div className="space-y-2">
                <p className="text-lg font-semibold">Pas encore de service</p>
                <p className="text-sm text-[var(--muted)]">
                  Dès que vous connecterez vos applications, elles apparaitront ici et pourront être utilisées dans vos
                  areas.
                </p>
              </div>
              <button
                type="button"
                onClick={openConnectModal}
                className="inline-flex items-center justify-center gap-2 rounded-full border border-[var(--blue-primary-2)] bg-[var(--blue-primary-2)] px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:border-[var(--blue-primary-3)] hover:bg-[var(--blue-primary-3)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)]"
              >
                Connecter un service
              </button>
            </div>
          </Card>
        </section>
      </div>
    </main>
  );
}
