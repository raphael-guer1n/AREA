"use client";

import { useCallback, useEffect, useRef, useState } from "react";

import { AreaNavigation } from "@/components/navigation/AreaNavigation";
import { ServiceCard } from "@/components/service/ServiceCard";
import { Card } from "@/components/ui/AreaCard";
import { useAuth } from "@/hooks/useAuth";
import { fetchOAuthAuthorizeUrl } from "@/lib/api/auth";
import { fetchAvailableProviders, fetchUserProviders } from "@/lib/api/services";
import { normalizeSearchValue } from "@/lib/helpers";

type Service = {
  id: string;
  name: string;
  url?: string;
  badge?: string;
  category?: string;
  gradient?: { from: string; to: string };
  actions?: string[];
  reactions?: string[];
  connected: boolean;
};

const gradients: Array<{ from: string; to: string }> = [
  { from: "#002642", to: "#0b3c5d" },
  { from: "#840032", to: "#a33a60" },
  { from: "#e59500", to: "#f2b344" },
  { from: "#5B834D", to: "#68915a" },
  { from: "#02040f", to: "#1b2640" },
];

function formatProviderName(provider: string) {
  return provider
    .split(/[-_]/)
    .filter(Boolean)
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join(" ");
}

function createService(provider: string, connected: boolean, index: number): Service {
  const gradient = gradients[index % gradients.length];
  const label = formatProviderName(provider) || provider;

  return {
    id: provider,
    name: label,
    badge: label.slice(0, 2).toUpperCase(),
    gradient,
    actions: [],
    reactions: [],
    connected,
  };
}

function matchesSearch(service: Service, term: string) {
  const normalizedTerm = normalizeSearchValue(term);
  if (!normalizedTerm) return true;
  const haystack = normalizeSearchValue(
    [
      service.name,
      service.category ?? "",
      ...(service.actions ?? []),
      ...(service.reactions ?? []),
    ].join(" "),
  );
  return haystack.includes(normalizedTerm);
}

export function ServicesClient() {
  const { user, token } = useAuth();
  const authToken = token ?? user?.token ?? null;
  const [services, setServices] = useState<Service[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [connectError, setConnectError] = useState<string | null>(null);
  const [isConnectModalOpen, setIsConnectModalOpen] = useState(false);
  const [searchTerm, setSearchTerm] = useState("");
  const [modalSearch, setModalSearch] = useState("");
  const isMountedRef = useRef(true);

  const filteredConnected = services.filter(
    (service) => service.connected && matchesSearch(service, searchTerm),
  );
  const connectableServices = services.filter(
    (service) => !service.connected && matchesSearch(service, modalSearch || searchTerm),
  );

  const updateConnection = (id: string, nextState: boolean) => {
    setServices((previous) =>
      previous.map((service) => (service.id === id ? { ...service, connected: nextState } : service)),
    );
  };

  const openConnectModal = () => {
    setConnectError(null);
    setIsConnectModalOpen(true);
  };
  const closeConnectModal = () => {
    setConnectError(null);
    setIsConnectModalOpen(false);
  };

  const handleServiceConnect = useCallback(
    async (providerId: string) => {
      if (!authToken) {
        setConnectError("Vous devez être connecté avant de lier un service.");
        return;
      }

      setConnectError(null);

      try {
        const { auth_url } = await fetchOAuthAuthorizeUrl(providerId, {
          token: authToken,
          mode: "link",
        });
        setIsConnectModalOpen(false);
        window.location.href = auth_url;
      } catch (err) {
        const message =
          err instanceof Error
            ? err.message
            : "Impossible de démarrer la connexion OAuth2.";
        setConnectError(message);
      }
    },
    [authToken],
  );

  useEffect(() => {
    return () => {
      isMountedRef.current = false;
    };
  }, []);

  const loadServices = useCallback(async () => {
    if (!authToken || !isMountedRef.current) {
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      let nextServices: Service[] = [];

      if (user?.id) {
        const providers = await fetchUserProviders(user.id, authToken);
        nextServices = providers.map((provider, index) =>
          createService(provider.provider, provider.is_logged, index),
        );
      } else {
        const providers = await fetchAvailableProviders(authToken);
        nextServices = providers.map((provider, index) => createService(provider, false, index));
      }

      if (isMountedRef.current) {
        setServices(nextServices);
      }
    } catch (err) {
      if (isMountedRef.current) {
        const message =
          err instanceof Error
            ? err.message
            : "Impossible de charger la liste des services pour le moment.";
        setError(message);
        setServices([]);
      }
    } finally {
      if (isMountedRef.current) {
        setIsLoading(false);
      }
    }
  }, [authToken, user?.id]);

  useEffect(() => {
    if (!authToken) return;
    void loadServices();
  }, [authToken, loadServices]);

  const refreshOnFocus = useCallback(() => {
    void loadServices();
  }, [loadServices]);

  useEffect(() => {
    if (!authToken) return;

    const handleVisibilityChange = () => {
      if (document.visibilityState === "visible") {
        void loadServices();
      }
    };

    window.addEventListener("focus", refreshOnFocus);
    document.addEventListener("visibilitychange", handleVisibilityChange);

    return () => {
      window.removeEventListener("focus", refreshOnFocus);
      document.removeEventListener("visibilitychange", handleVisibilityChange);
    };
  }, [authToken, loadServices, refreshOnFocus]);

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

                {connectError ? (
                  <div className="rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
                    {connectError}
                  </div>
                ) : null}

                <div className="rounded-2xl border-2 border border-[var(--surface-border)] bg-[var(--surface)] p-4 sm:p-6">
                  {isLoading ? (
                    <div className="flex min-h-[220px] items-center justify-center text-sm text-[var(--muted)]">
                      Chargement des services disponibles...
                    </div>
                  ) : connectableServices.length ? (
                    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                      {connectableServices.map((service) => (
                        <ServiceCard
                          key={service.id}
                          name={service.name}
                          url={service.url ?? ""}
                          badge={service.badge ?? service.name.slice(0, 2).toUpperCase()}
                          category={service.category}
                          gradientFrom={service.gradient?.from}
                          gradientTo={service.gradient?.to}
                          actions={service.actions ?? []}
                          reactions={service.reactions ?? []}
                          connected={service.connected}
                          action="À connecter"
                          onConnect={() => {
                            void handleServiceConnect(service.id);
                          }}
                        />
                      ))}
                    </div>
                  ) : error ? (
                    <div className="flex min-h-[220px] flex-col items-center justify-center gap-3 text-center">
                      <p className="text-base font-semibold text-[var(--foreground)]">
                        Impossible de charger les services
                      </p>
                      <p className="text-sm text-[var(--muted)]">{error}</p>
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
            {error ? (
              <div className="mb-4 rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
                {error}
              </div>
            ) : null}
            {isLoading ? (
              <div className="flex min-h-[200px] items-center justify-center text-sm text-[var(--muted)]">
                Chargement de vos services...
              </div>
            ) : filteredConnected.length ? (
              <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                {filteredConnected.map((service) => (
                  <ServiceCard
                    key={service.id}
                    name={service.name}
                    url={service.url ?? ""}
                    badge={service.badge ?? service.name.slice(0, 2).toUpperCase()}
                    category={service.category}
                    gradientFrom={service.gradient?.from}
                    gradientTo={service.gradient?.to}
                    actions={service.actions ?? []}
                    reactions={service.reactions ?? []}
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

export default ServicesClient;
