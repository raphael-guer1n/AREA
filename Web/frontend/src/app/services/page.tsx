"use client";

import { Suspense, useCallback, useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";

import { AreaNavigation } from "@/components/navigation/AreaNavigation";
import { ServiceCard } from "@/components/service/ServiceCard";
import { Card } from "@/components/ui/AreaCard";
import { fetchServices, fetchUserServiceStatuses } from "@/lib/api/services";
import { useAuth } from "@/hooks/useAuth";
import { useOAuthCallback } from "@/hooks/useOAuthCallback";
import { normalizeSearchValue } from "@/lib/helpers";
import { gradients as gradientPalette, mockServices, type MockService } from "./mockServices";

const serviceTemplatesById = new Map(mockServices.map((service) => [service.id, service]));

function matchesSearch(service: MockService, term: string) {
  const normalizedTerm = normalizeSearchValue(term);
  if (!normalizedTerm) return true;
  const haystack = normalizeSearchValue(
    [service.name, service.category ?? "", ...service.actions, ...service.reactions].join(" "),
  );
  return haystack.includes(normalizedTerm);
}

function formatServiceNameFromId(serviceId: string) {
  return serviceId
    .split(/[-_]/)
    .filter(Boolean)
    .map((segment) => segment[0]?.toUpperCase() + segment.slice(1))
    .join(" ");
}

function buildBadge(name: string) {
  const parts = name.split(" ").filter(Boolean);
  const letters =
    parts.length >= 2
      ? `${parts[0]?.[0] ?? ""}${parts[1]?.[0] ?? ""}`
      : name.slice(0, 2);
  return letters.toUpperCase();
}

function mapBackendService(serviceId: string, index: number): MockService {
  const template = serviceTemplatesById.get(serviceId);
  const formattedName = formatServiceNameFromId(serviceId);
  const name = template?.name ?? (formattedName || serviceId);
  const badge = template?.badge ?? buildBadge(name || serviceId);
  const gradient = template?.gradient ?? gradientPalette[index % gradientPalette.length];

  return {
    id: serviceId,
    name,
    url: template?.url ?? "#",
    badge,
    category: template?.category,
    gradient,
    actions: template?.actions ?? [],
    reactions: template?.reactions ?? [],
    connected: template?.connected ?? false,
  };
}

export function ServicesClient() {
  const { token, user, startOAuthConnect } = useAuth();
  const searchParams = useSearchParams();
  const hasOAuthParams = Boolean(searchParams.get("code") && searchParams.get("state"));
  const { status: oauthStatus, error: oauthError } = useOAuthCallback("/services", {
    enabled: hasOAuthParams,
  });
  const isProcessingOAuth = hasOAuthParams && oauthStatus !== "error";
  const [services, setServices] = useState<MockService[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isConnectModalOpen, setIsConnectModalOpen] = useState(false);
  const [searchTerm, setSearchTerm] = useState("");
  const [modalSearch, setModalSearch] = useState("");

  const loadServices = useCallback(async () => {
    setIsLoading(true);
    setError(null);

    try {
      const serviceIds = await fetchServices();
      const uniqueServiceIds = Array.from(new Set(serviceIds.filter(Boolean)));

      let statusByService: Record<string, boolean> = {};
      if (token && user?.id) {
        try {
          const statuses = await fetchUserServiceStatuses(token, user.id);
          statusByService = statuses.reduce<Record<string, boolean>>((acc, current) => {
            acc[current.provider] = Boolean(current.is_logged);
            return acc;
          }, {});
        } catch (statusError) {
          console.error(statusError);
        }
      }

      const mappedServices = uniqueServiceIds.map((serviceId, index) =>
        ({
          ...mapBackendService(serviceId, index),
          connected: Boolean(statusByService[serviceId]),
        }),
      );
      setServices(mappedServices);
    } catch (err) {
      const message =
        err instanceof Error
          ? err.message
          : "Impossible de récupérer les services.";
      setError(message);
      setServices([]);
    } finally {
      setIsLoading(false);
    }
  }, [token, user?.id]);

  useEffect(() => {
    void loadServices();
  }, [loadServices]);

  useEffect(() => {
    if (oauthError) {
      setError(oauthError);
    }
  }, [oauthError]);

  const filteredConnected = services.filter(
    (service) => service.connected && matchesSearch(service, searchTerm),
  );
  const connectableServices = services.filter(
    (service) => !service.connected && matchesSearch(service, modalSearch || searchTerm),
  );

  const handleConnect = async (serviceId: string) => {
    if (!token) {
      setError("Vous devez être connecté pour lier un service.");
      return;
    }
    try {
      setError(null);
      await startOAuthConnect(serviceId);
    } catch (err) {
      const message =
        err instanceof Error
          ? err.message
          : "Une erreur est survenue lors de la connexion du service.";
      setError(message);
    }
  };

  const updateConnection = (id: string, nextState: boolean) => {
    setServices((previous) =>
      previous.map((service) => (service.id === id ? { ...service, connected: nextState } : service)),
    );
  };

  const openConnectModal = () => setIsConnectModalOpen(true);
  const closeConnectModal = () => setIsConnectModalOpen(false);

  return (
    <main className="relative flex min-h-screen justify-center overflow-hidden bg-[var(--surface)] px-6 py-12 pt-10 text-[var(--foreground)]">
      {isProcessingOAuth ? (
        <div className="absolute inset-0 z-50 flex items-center justify-center bg-[rgba(6,14,25,0.35)] backdrop-blur-sm">
          <div className="rounded-xl border border-[var(--surface-border)] bg-[var(--background)] px-4 py-3 text-sm text-[var(--muted)] shadow-lg">
            Connexion du service en cours...
          </div>
        </div>
      ) : null}
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
                      Connectez un nouveau service
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
                    value={modalSearch}
                    onChange={(event) => setModalSearch(event.target.value)}
                    placeholder="Rechercher un service (Slack, Gmail, Discord...)"
                    className="w-full rounded-xl border border-[var(--surface-border)] bg-[var(--surface)] px-11 py-3 text-sm text-[var(--foreground)] placeholder:text-[var(--placeholder)] focus:border-[var(--blue-primary-3)] focus:outline-none focus:ring-2 focus:ring-[var(--blue-primary-3)]/25"
                  />
                </div>

                <div className="rounded-2xl border-2 border border-[var(--surface-border)] bg-[var(--surface)] p-4 sm:p-6">
                  {isLoading ? (
                    <div className="flex min-h-[220px] flex-col items-center justify-center gap-3 text-center">
                      <p className="text-base font-semibold text-[var(--foreground)]">
                        Chargement des services...
                      </p>
                    </div>
                  ) : error ? (
                    <div className="flex min-h-[220px] flex-col items-center justify-center gap-3 text-center">
                      <p className="text-base font-semibold text-[var(--foreground)]">
                        Impossible de charger les services
                      </p>
                      <p className="text-sm text-[var(--muted)]">{error}</p>
                      <button
                        type="button"
                        onClick={loadServices}
                        className="inline-flex items-center justify-center gap-2 rounded-full border border-[var(--blue-primary-2)] bg-[var(--blue-primary-2)] px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:border-[var(--blue-primary-3)] hover:bg-[var(--blue-primary-3)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)]"
                      >
                        Réessayer
                      </button>
                    </div>
                  ) : connectableServices.length ? (
                    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                      {connectableServices.map((service) => (
                        <ServiceCard
                          key={service.id}
                          name={service.name}
                          url={service.url}
                          badge={service.badge}
                          category={service.category}
                          gradientFrom={service.gradient.from}
                          gradientTo={service.gradient.to}
                          actions={service.actions}
                          reactions={service.reactions}
                          connected={service.connected}
                          action="À connecter"
                          onConnect={() => handleConnect(service.id)}
                        />
                      ))}
                    </div>
                  ) : (
                    <div className="flex min-h-[220px] flex-col items-center justify-center gap-3 text-center">
                      <p className="text-base font-semibold text-[var(--foreground)]">
                        Aucun service à connecter trouvé
                      </p>
                      <p className="text-sm text-[var(--muted)]">
                        Ajustez votre recherche ou revenez plus tard pour de nouvelles intégrations.
                      </p>
                    </div>
                  )}
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
                value={searchTerm}
                onChange={(event) => setSearchTerm(event.target.value)}
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
            {isLoading ? (
              <div className="flex flex-col items-center justify-center gap-3 rounded-2xl border border-[var(--surface-border)] bg-[var(--background)] px-6 py-10 text-center">
                <p className="text-lg font-semibold text-[var(--foreground)]">Chargement des services...</p>
                <p className="text-sm text-[var(--muted)]">Nous récupérons vos services disponibles.</p>
              </div>
            ) : error ? (
              <div className="flex flex-col items-center justify-center gap-5 rounded-2xl border border-[var(--surface-border)] bg-[var(--background)] px-6 py-10 text-center">
                <div className="flex h-14 w-14 items-center justify-center rounded-2xl border border-[var(--surface-border)] bg-[var(--surface)]">
                  <div className="h-6 w-6 rounded-full border-2 border-[var(--blue-primary-2)]" />
                </div>
                <div className="space-y-2">
                  <p className="text-lg font-semibold text-[var(--foreground)]">
                    Impossible de charger vos services
                  </p>
                  <p className="text-sm text-[var(--muted)]">{error}</p>
                </div>
                <button
                  type="button"
                  onClick={loadServices}
                  className="inline-flex items-center justify-center gap-2 rounded-full border border-[var(--blue-primary-2)] bg-[var(--blue-primary-2)] px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:border-[var(--blue-primary-3)] hover:bg-[var(--blue-primary-3)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)]"
                >
                  Réessayer
                </button>
              </div>
            ) : filteredConnected.length ? (
              <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                {filteredConnected.map((service) => (
                  <ServiceCard
                    key={service.id}
                    name={service.name}
                    url={service.url}
                    badge={service.badge}
                    category={service.category}
                    gradientFrom={service.gradient.from}
                    gradientTo={service.gradient.to}
                    actions={service.actions}
                    reactions={service.reactions}
                    connected={service.connected}
                    action="Connecté"
                    onDisconnect={() => updateConnection(service.id, false)}
                  />
                ))}
              </div>
            ) : (
              <div className="flex flex-col items-center justify-center gap-5 rounded-2xl border border-[var(--surface-border)] bg-[var(--background)] px-6 py-10 text-center">
                <div className="flex h-14 w-14 items-center justify-center rounded-2xl border border-[var(--surface-border)] bg-[var(--surface)]">
                  <div className="h-6 w-6 rounded-full border-2 border-[var(--surface-border)]" />
                </div>
                <div className="space-y-2">
                  <p className="text-lg font-semibold">Pas encore de service connecté</p>
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
            )}
          </Card>
        </section>
      </div>
    </main>
  );
}

export default function ServicesPage() {
  return (
    <Suspense
      fallback={
        <main className="flex min-h-screen items-center justify-center bg-[var(--surface)] px-6 py-12">
          <p className="text-sm text-[var(--muted)]">Chargement des services...</p>
        </main>
      }
    >
      <ServicesClient />
    </Suspense>
  );
}
