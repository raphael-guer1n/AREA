"use client";

import { Suspense, useCallback, useEffect, useMemo, useState } from "react";
import { useSearchParams } from "next/navigation";

import { AreaNavigation } from "@/components/navigation/AreaNavigation";
import { AreaCard as AreaTileCard } from "@/components/area/AreaCard";
import { Card } from "@/components/ui/AreaCard";
import { cn, normalizeSearchValue } from "@/lib/helpers";
import { useAuth } from "@/hooks/useAuth";
import { useOAuthCallback } from "@/hooks/useOAuthCallback";
import { createEventArea } from "@/lib/api/area";
import { fetchServices, fetchUserServiceStatuses } from "@/lib/api/services";

import { gradients as gradientPalette, mockServices, type MockService } from "../services/mockServices";

type AreaService = MockService;
type AreaGradient = { from: string; to: string };

type CreatedArea = {
  id: string;
  name: string;
  summary: string;
  startTime: string;
  endTime: string;
  delay: number;
  serviceName: string;
  actionService: string;
  reactionService: string;
  actionName: string;
  reactionName: string;
  gradient: AreaGradient;
};

function matchesAreaSearch(area: CreatedArea, normalizedTerm: string) {
  if (!normalizedTerm) return true;
  const haystack = normalizeSearchValue(
    [area.name, area.summary, area.serviceName, area.startTime, area.endTime].join(" "),
  );
  return haystack.includes(normalizedTerm);
}

function pickRandomGradient(): AreaGradient {
  if (!gradientPalette.length) {
    return { from: "#0b3c5d", to: "#e59500" };
  }
  const index = Math.floor(Math.random() * gradientPalette.length);
  return gradientPalette[index];
}

function AreaPageContent() {
  const { token, user } = useAuth();
  const searchParams = useSearchParams();
  const hasOAuthParams = Boolean(searchParams.get("code") && searchParams.get("state"));
  const { status, error } = useOAuthCallback("/area", { enabled: hasOAuthParams });
  const isProcessingOAuth = hasOAuthParams && status !== "error";
  const [areas, setAreas] = useState<CreatedArea[]>([]);
  const [services, setServices] = useState<AreaService[]>([]);
  const [servicesError, setServicesError] = useState<string | null>(null);
  const [isLoadingServices, setIsLoadingServices] = useState(false);
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [selectedService, setSelectedService] = useState<AreaService | null>(null);
  const [delay, setDelay] = useState<number>(0);
  const [startTime, setStartTime] = useState("");
  const [endTime, setEndTime] = useState("");
  const [summary, setSummary] = useState("");
  const [description, setDescription] = useState("");
  const [selectedAreaDetail, setSelectedAreaDetail] = useState<CreatedArea | null>(null);
  const [createError, setCreateError] = useState<string | null>(null);
  const [isCreating, setIsCreating] = useState(false);
  const [searchTerm, setSearchTerm] = useState("");
  const normalizedSearch = normalizeSearchValue(searchTerm);
  const filteredAreas = areas.filter((area) => matchesAreaSearch(area, normalizedSearch));
  const hasSearch = Boolean(normalizedSearch);
  const totalAreas = areas.length;
  const activeCount = areas.length;
  const hasConnectedServices = useMemo(
    () => services.some((service) => service.connected),
    [services],
  );

  const loadServices = useCallback(async () => {
    setIsLoadingServices(true);
    setServicesError(null);

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

      const mappedServices = uniqueServiceIds.map((serviceId, index) => {
        const template = mockServices.find((service) => service.id === serviceId);
        const gradient = template?.gradient ?? gradientPalette[index % gradientPalette.length];
        const connected = statusByService[serviceId];

        return {
          ...(template ?? {
            id: serviceId,
            name: serviceId,
            url: "#",
            badge: serviceId.slice(0, 2).toUpperCase(),
            category: "Service",
            gradient,
            actions: [],
            reactions: [],
            connected: false,
          }),
          gradient,
          // Si le backend ne renvoie pas de statut, on considère le service comme non connecté.
          connected: Boolean(connected),
        } as AreaService;
      });

      setServices(mappedServices);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Impossible de récupérer les services.";
      setServicesError(message);
      setServices([]);
    } finally {
      setIsLoadingServices(false);
    }
  }, [token, user?.id]);

  useEffect(() => {
    void loadServices();
  }, [loadServices]);

  const connectedServices = useMemo(
    () => services.filter((service) => service.connected),
    [services],
  );

  const resetForm = () => {
    setSelectedService(null);
    setDelay(0);
    setStartTime("");
    setEndTime("");
    setSummary("");
    setDescription("");
    setCreateError(null);
  };

  const openModal = () => {
    resetForm();
    setIsCreateModalOpen(true);
  };

  const closeModal = () => {
    setIsCreateModalOpen(false);
  };

  const handleCreateArea = async () => {
    if (!token) {
      setCreateError("Vous devez être connecté pour créer une area.");
      return;
    }
    if (!selectedService) {
      setCreateError("Sélectionnez un service connecté.");
      return;
    }
    if (!startTime || !endTime || !summary) {
      setCreateError("Renseignez au moins la date de début/fin et le résumé.");
      return;
    }

    setIsCreating(true);
    setCreateError(null);

    try {
      const gradient = pickRandomGradient();
      const actionName = selectedService.actions?.[0] ?? "Action sélectionnée";
      const reactionName = selectedService.reactions?.[0] ?? "Réaction sélectionnée";

      await createEventArea(token, {
        delay: Number.isFinite(delay) ? delay : 0,
        event: {
          startTime,
          endTime,
          summary: summary.trim(),
          description: description.trim(),
        },
      });

      const newArea: CreatedArea = {
        id: `area-${Date.now()}`,
        name: `${selectedService.name} → Création d'événement`,
        summary: summary.trim(),
        startTime,
        endTime,
        delay: Number.isFinite(delay) ? delay : 0,
        serviceName: selectedService.name,
        actionService: selectedService.name,
        reactionService: selectedService.name,
        actionName,
        reactionName,
        gradient,
      };

      setAreas((prev) => [newArea, ...prev]);
      setIsCreateModalOpen(false);
      resetForm();
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Impossible de créer l'area.";
      setCreateError(message);
    } finally {
      setIsCreating(false);
    }
  };

  if (hasOAuthParams) {
    return (
      <main className="flex min-h-screen items-center justify-center bg-[var(--surface)] px-4 py-12">
        <div className="w-full max-w-lg rounded-3xl border border-[var(--surface-border)] bg-[var(--background)] px-8 py-10 shadow-[0_20px_60px_rgba(17,42,70,0.08)]">
          <div className="space-y-3 text-center">
            <p className="text-xs font-semibold uppercase tracking-[0.14em] text-[var(--muted)]">
              OAuth2 Authentication
            </p>
            <h1 className="text-2xl font-semibold text-[var(--foreground)]">
              {isProcessingOAuth ? "Signing you in..." : "Sign-in failed"}
            </h1>
            <p className="text-sm text-[var(--muted)]">
              {isProcessingOAuth
                ? "Please wait while we validate your account."
                : error ?? "Unable to finish the OAuth2 login."}
            </p>
          </div>
        </div>
      </main>
    );
  }

  return (
    <main className="relative flex min-h-screen justify-center overflow-hidden bg-[var(--surface)] px-6 py-12 pt-10 text-[var(--foreground)]">
      {isCreateModalOpen ? (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center bg-[rgba(6,14,25,0.55)] px-4 py-10 backdrop-blur-sm"
          onClick={closeModal}
        >
          <div
            className="relative w-full max-w-5xl overflow-hidden rounded-[28px] border border-[var(--surface-border)] bg-[var(--background)] shadow-2xl ring-1 ring-[rgba(28,61,99,0.28)]"
            onClick={(event) => event.stopPropagation()}
            role="dialog"
            aria-modal="true"
            aria-labelledby="create-area-title"
          >
            <button
              type="button"
              className="absolute right-5 top-5 inline-flex h-10 w-10 items-center justify-center rounded-full border border-[var(--surface-border)] bg-[var(--surface)] text-[var(--muted)] transition hover:text-[var(--foreground)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)]"
              onClick={closeModal}
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
                  Création d&apos;area
                </p>
                <div className="flex flex-wrap items-end justify-between gap-3">
                  <div className="space-y-1">
                    <h2 id="create-area-title" className="text-2xl font-semibold text-[var(--foreground)]">
                      Choisissez vos services connectés
                    </h2>
                    <p className="text-sm text-[var(--muted)]">
                      Sélectionnez un service et son action, puis un service et sa réaction.
                    </p>
                  </div>
                </div>
              </div>

              <div className="space-y-5">
                {servicesError ? (
                  <div className="rounded-xl border border-[var(--accent)] bg-[var(--surface)] px-4 py-3 text-sm text-[var(--accent)]">
                    {servicesError}
                  </div>
                ) : null}

                <div className="grid gap-5 lg:grid-cols-2">
                  <div className="rounded-2xl border border-[var(--surface-border)] bg-[var(--surface)] p-4 sm:p-6">
                    <div className="flex items-center justify-between">
                      <div>
                        <p className="text-xs font-semibold uppercase tracking-[0.14em] text-[var(--blue-primary-3)]">
                          Service connecté
                        </p>
                        <p className="text-sm text-[var(--muted)]">Choisissez le service qui créera l&apos;événement</p>
                      </div>
                    </div>
                    <div className="mt-4 space-y-3">
                      {isLoadingServices ? (
                        <p className="text-sm text-[var(--muted)]">Chargement des services...</p>
                      ) : connectedServices.length ? (
                        <div className="grid gap-2 sm:grid-cols-2">
                          {connectedServices.map((service) => (
                            <button
                              key={service.id}
                              type="button"
                              onClick={() => setSelectedService(service)}
                              className={cn(
                                "flex items-center justify-between rounded-xl border px-4 py-3 text-left transition",
                                selectedService?.id === service.id
                                  ? "border-[var(--blue-primary-2)] bg-[var(--blue-primary-2)]/10"
                                  : "border-[var(--surface-border)] hover:border-[var(--blue-primary-2)]",
                              )}
                            >
                              <div className="space-y-1">
                                <p className="text-sm font-semibold">{service.name}</p>
                                <p className="text-xs text-[var(--muted)]">{service.category ?? "Service"}</p>
                              </div>
                              <span className="rounded-full bg-[var(--blue-primary-2)] px-3 py-1 text-[11px] font-semibold uppercase text-white">
                                {service.badge}
                              </span>
                            </button>
                          ))}
                        </div>
                      ) : (
                        <p className="text-sm text-[var(--muted)]">
                          Aucun service connecté. Connectez un service pour commencer.
                        </p>
                      )}
                    </div>
                  </div>

                  <div className="rounded-2xl border border-[var(--surface-border)] bg-[var(--surface)] p-4 sm:p-6">
                    <div className="flex items-center justify-between">
                      <div>
                        <p className="text-xs font-semibold uppercase tracking-[0.14em] text-[var(--blue-primary-3)]">
                          Détails de l&apos;événement
                        </p>
                        <p className="text-sm text-[var(--muted)]">Programmation et contenu</p>
                      </div>
                    </div>
                    <div className="mt-4 space-y-3">
                      <div className="grid gap-3 sm:grid-cols-2">
                        <label className="space-y-1 text-sm">
                          <span className="text-[var(--muted)]">Début</span>
                          <input
                            type="datetime-local"
                            value={startTime}
                            onChange={(e) => setStartTime(e.target.value)}
                            className="w-full rounded-lg border border-[var(--surface-border)] bg-[var(--surface)] px-3 py-2 text-[var(--foreground)] focus:border-[var(--blue-primary-3)] focus:outline-none focus:ring-2 focus:ring-[var(--blue-primary-3)]/25"
                          />
                        </label>
                        <label className="space-y-1 text-sm">
                          <span className="text-[var(--muted)]">Fin</span>
                          <input
                            type="datetime-local"
                            value={endTime}
                            onChange={(e) => setEndTime(e.target.value)}
                            className="w-full rounded-lg border border-[var(--surface-border)] bg-[var(--surface)] px-3 py-2 text-[var(--foreground)] focus:border-[var(--blue-primary-3)] focus:outline-none focus:ring-2 focus:ring-[var(--blue-primary-3)]/25"
                          />
                        </label>
                      </div>

                      <label className="space-y-1 text-sm">
                        <span className="text-[var(--muted)]">Résumé</span>
                        <input
                          type="text"
                          value={summary}
                          onChange={(e) => setSummary(e.target.value)}
                          placeholder="Team Meeting"
                          className="w-full rounded-lg border border-[var(--surface-border)] bg-[var(--surface)] px-3 py-2 text-[var(--foreground)] focus:border-[var(--blue-primary-3)] focus:outline-none focus:ring-2 focus:ring-[var(--blue-primary-3)]/25"
                        />
                      </label>

                      <label className="space-y-1 text-sm">
                        <span className="text-[var(--muted)]">Description</span>
                        <textarea
                          value={description}
                          onChange={(e) => setDescription(e.target.value)}
                          rows={3}
                          placeholder="Weekly sync..."
                          className="w-full rounded-lg border border-[var(--surface-border)] bg-[var(--surface)] px-3 py-2 text-[var(--foreground)] focus:border-[var(--blue-primary-3)] focus:outline-none focus:ring-2 focus:ring-[var(--blue-primary-3)]/25"
                        />
                      </label>

                      <label className="space-y-1 text-sm">
                        <span className="text-[var(--muted)]">Délai avant exécution (secondes)</span>
                        <input
                          type="number"
                          min={0}
                          value={delay}
                          onChange={(e) => setDelay(Number(e.target.value) || 0)}
                          className="w-full rounded-lg border border-[var(--surface-border)] bg-[var(--surface)] px-3 py-2 text-[var(--foreground)] focus:border-[var(--blue-primary-3)] focus:outline-none focus:ring-2 focus:ring-[var(--blue-primary-3)]/25"
                        />
                      </label>
                    </div>
                  </div>
                </div>

                {createError ? (
                  <div className="rounded-xl border border-[var(--accent)] bg-[var(--surface)] px-4 py-3 text-sm text-[var(--accent)]">
                    {createError}
                  </div>
                ) : null}

                <div className="flex flex-wrap justify-end gap-3">
                  <button
                    type="button"
                    onClick={closeModal}
                    className="inline-flex items-center justify-center gap-2 rounded-full border border-[var(--surface-border)] bg-[var(--surface)] px-4 py-2 text-sm font-semibold text-[var(--foreground)] shadow-sm transition hover:border-[var(--blue-primary-2)] hover:text-[var(--blue-primary-2)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)]"
                  >
                    Annuler
                  </button>
                  <button
                    type="button"
                    onClick={handleCreateArea}
                    disabled={
                      isCreating ||
                      !selectedService ||
                      !startTime ||
                      !endTime ||
                      !summary ||
                      !hasConnectedServices
                    }
                    className="inline-flex items-center justify-center gap-2 rounded-full border border-[var(--blue-primary-2)] bg-[var(--blue-primary-2)] px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:border-[var(--blue-primary-3)] hover:bg-[var(--blue-primary-3)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)] disabled:cursor-not-allowed disabled:opacity-50"
                  >
                    {isCreating ? "Création..." : "Créer l'area"}
                  </button>
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
              <button
                type="button"
                onClick={openModal}
                className="inline-flex items-center justify-center gap-2 rounded-full border border-[var(--blue-primary-2)] bg-[var(--blue-primary-2)] px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:border-[var(--blue-primary-3)] hover:bg-[var(--blue-primary-3)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)]"
              >
                Create area
              </button>
            }
            className="relative w-full overflow-hidden rounded-[26px] border-[var(--surface-border)] bg-[var(--background)] ring-1 ring-[rgba(28,61,99,0.15)]"
          >
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
                  value={searchTerm}
                  onChange={(event) => setSearchTerm(event.target.value)}
                  placeholder="Rechercher une area (résumé, service, dates...)"
                  className="w-full rounded-xl border border-[var(--surface-border)] bg-[var(--background)] px-11 py-3 text-sm text-[var(--foreground)] placeholder:text-[var(--placeholder)] focus:border-[var(--blue-primary-3)] focus:outline-none focus:ring-2 focus:ring-[var(--blue-primary-3)]/25"
                />
              </div>

              {filteredAreas.length ? (
                <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                  {filteredAreas.map((area) => {
                    const badge = (area.actionService || area.serviceName || "?").slice(0, 2).toUpperCase();
                    const reactionBadge = (area.reactionService || "?").slice(0, 2).toUpperCase();
                    return (
                      <AreaTileCard
                        key={area.id}
                        id={area.id}
                        name={area.summary || area.name}
                        actionLabel={area.actionName}
                        reactionLabel={area.reactionName}
                        actionIcon={<span>{badge}</span>}
                        reactionIcon={<span>{reactionBadge}</span>}
                        gradientFrom={area.gradient.from}
                        gradientTo={area.gradient.to}
                        isActive
                        onClick={() => setSelectedAreaDetail(area)}
                        className="h-full"
                      />
                    );
                  })}
                </div>
              ) : hasSearch ? (
                <div className="flex flex-col items-center justify-center gap-5 rounded-2xl border border-[var(--surface-border)] bg-[var(--background)] px-6 py-10 text-center">
                  <div className="flex h-14 w-14 items-center justify-center rounded-2xl border border-[var(--surface-border)] bg-[var(--surface)]">
                    <div className="h-6 w-6 rounded-full border-2 border-[var(--surface-border)]" />
                  </div>
                  <div className="space-y-2">
                    <p className="text-lg font-semibold">Aucune area trouvée</p>
                    <p className="text-sm text-[var(--muted)]">
                      Aucun résultat pour cette recherche. Essayez un autre mot-clé ou réinitialisez la recherche.
                    </p>
                  </div>
                  <button
                    type="button"
                    onClick={() => setSearchTerm("")}
                    className="inline-flex items-center justify-center gap-2 rounded-full border border-[var(--surface-border)] bg-[var(--surface)] px-4 py-2 text-sm font-semibold text-[var(--foreground)] shadow-sm transition hover:border-[var(--blue-primary-2)] hover:text-[var(--blue-primary-2)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)]"
                  >
                    Réinitialiser la recherche
                  </button>
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
                <button
                  type="button"
                  onClick={openModal}
                  className="inline-flex items-center justify-center gap-2 rounded-full border border-[var(--blue-primary-2)] bg-[var(--blue-primary-2)] px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:border-[var(--blue-primary-3)] hover:bg-[var(--blue-primary-3)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)]"
                >
                  Créer une area
                </button>
              </div>
            )}
          </div>
        </Card>
        </section>

        {selectedAreaDetail ? (
          <div
            className="fixed inset-0 z-50 flex items-center justify-center bg-[rgba(6,14,25,0.5)] px-4 py-10 backdrop-blur-sm"
            onClick={() => setSelectedAreaDetail(null)}
            role="dialog"
            aria-modal="true"
            aria-label="Détail de l'area"
          >
            <div
              className="relative w-full max-w-2xl overflow-hidden rounded-3xl border border-[var(--surface-border)] bg-[var(--background)] shadow-2xl"
              onClick={(event) => event.stopPropagation()}
            >
              <div
                className="h-2 w-full"
                style={{
                  background: `linear-gradient(90deg, ${selectedAreaDetail.gradient.from}, ${selectedAreaDetail.gradient.to})`,
                }}
              />
              <div className="flex items-start justify-between px-6 py-5">
                <div className="space-y-1">
                  <p className="text-xs font-semibold uppercase tracking-[0.14em] text-[var(--blue-primary-3)]">
                    Détail de l&apos;automation
                  </p>
                  <h3 className="text-2xl font-semibold text-[var(--foreground)]">
                    {selectedAreaDetail.summary || selectedAreaDetail.name}
                  </h3>
                  <p className="text-sm text-[var(--muted)]">{selectedAreaDetail.name}</p>
                </div>
                <button
                  type="button"
                  onClick={() => setSelectedAreaDetail(null)}
                  className="inline-flex h-10 w-10 items-center justify-center rounded-full border border-[var(--surface-border)] bg-[var(--surface)] text-[var(--muted)] transition hover:text-[var(--foreground)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)]"
                  aria-label="Fermer"
                >
                  <svg className="h-5 w-5" viewBox="0 0 20 20" fill="none" aria-hidden>
                    <path d="M15 5 5 15m0-10 10 10" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" />
                  </svg>
                </button>
              </div>

              <div className="grid gap-6 px-6 pb-6 md:grid-cols-2">
                <div className="space-y-3 rounded-2xl border border-[var(--surface-border)] bg-[var(--surface)] p-4">
                  <p className="text-xs font-semibold uppercase tracking-[0.12em] text-[var(--muted)]">
                    Service d&apos;action
                  </p>
                  <p className="text-base font-semibold text-[var(--foreground)]">
                    {selectedAreaDetail.actionService}
                  </p>
                  <p className="text-sm text-[var(--muted)]">{selectedAreaDetail.actionName}</p>
                </div>
                <div className="space-y-3 rounded-2xl border border-[var(--surface-border)] bg-[var(--surface)] p-4">
                  <p className="text-xs font-semibold uppercase tracking-[0.12em] text-[var(--muted)]">
                    Service de réaction
                  </p>
                  <p className="text-base font-semibold text-[var(--foreground)]">
                    {selectedAreaDetail.reactionService}
                  </p>
                  <p className="text-sm text-[var(--muted)]">{selectedAreaDetail.reactionName}</p>
                </div>
                <div className="space-y-2 rounded-2xl border border-[var(--surface-border)] bg-[var(--surface)] p-4 md:col-span-2">
                  <p className="text-xs font-semibold uppercase tracking-[0.12em] text-[var(--muted)]">
                    Paramètres
                  </p>
                  <div className="grid gap-3 sm:grid-cols-3">
                    <div>
                      <p className="text-[var(--muted)] text-xs">Début</p>
                      <p className="text-sm font-semibold text-[var(--foreground)]">{selectedAreaDetail.startTime}</p>
                    </div>
                    <div>
                      <p className="text-[var(--muted)] text-xs">Fin</p>
                      <p className="text-sm font-semibold text-[var(--foreground)]">{selectedAreaDetail.endTime}</p>
                    </div>
                    <div>
                      <p className="text-[var(--muted)] text-xs">Délai</p>
                      <p className="text-sm font-semibold text-[var(--foreground)]">
                        {selectedAreaDetail.delay}s
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        ) : null}
      </div>
    </main>
  );
}

export default function AreaPage() {
  return (
    <Suspense
      fallback={
        <main className="flex min-h-screen items-center justify-center bg-[var(--surface)] px-6 py-12">
          <p className="text-sm text-[var(--muted)]">Chargement...</p>
        </main>
      }
    >
      <AreaPageContent />
    </Suspense>
  );
}
