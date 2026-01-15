"use client";

import { Suspense, useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useSearchParams } from "next/navigation";

import { AreaNavigation } from "@/components/navigation/AreaNavigation";
import { AreaCard as AreaTileCard } from "@/components/area/AreaCard";
import { Card } from "@/components/ui/AreaCard";
import { cn, normalizeSearchValue } from "@/lib/helpers";
import { useAuth } from "@/hooks/useAuth";
import { useOAuthCallback } from "@/hooks/useOAuthCallback";
import { activateArea, deactivateArea, deleteArea, fetchAreas, saveArea, type BackendArea } from "@/lib/api/area";
import {
  fetchServiceConfig,
  fetchServiceNames,
  fetchUserServiceStatuses,
  type ServiceActionConfig,
  type ServiceConfig,
  type ServiceFieldConfig,
  type ServiceReactionConfig,
} from "@/lib/api/services";

type FieldDefinition = ServiceFieldConfig;
type FieldValue = string | string[];
type ActionDefinition = {
  id: string;
  title: string;
  label: string;
  type: string;
  fields: FieldDefinition[];
  output_fields?: FieldDefinition[];
};
type ReactionDefinition = {
  id: string;
  title: string;
  label: string;
  url?: string;
  method?: string;
  fields: FieldDefinition[];
};

type AreaGradient = { from: string; to: string };
type AreaService = {
  id: string;
  name: string;
  provider: string;
  badge: string;
  logoUrl?: string;
  gradient: AreaGradient;
  actions: ActionDefinition[];
  reactions: ReactionDefinition[];
  connected: boolean;
};

type ReactionForm = {
  id: string;
  service: AreaService | null;
  reaction: ReactionDefinition | null;
  fieldValues: Record<string, FieldValue>;
};

const gradientPalette: AreaGradient[] = [
  { from: "#002642", to: "#0b3c5d" },
  { from: "#840032", to: "#a33a60" },
  { from: "#e59500", to: "#f2b344" },
  { from: "#5B834D", to: "#68915a" },
  { from: "#02040f", to: "#1b2640" },
];

type CreatedArea = {
  id: string;
  name: string;
  summary: string;
  startTime: string;
  endTime: string;
  delay: number;
  serviceName: string;
  actionServiceId: string;
  actionService: string;
  primaryReactionServiceId: string;
  reactionServiceIds: string[];
  reactionService: string;
  reactionCount: number;
  actionName: string;
  reactionName: string;
  gradient: AreaGradient;
  active: boolean;
};

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

function mapActionConfig(action: ServiceActionConfig): ActionDefinition {
  return {
    id: action.title,
    title: action.title,
    label: action.label ?? action.title,
    type: action.type,
    fields: action.fields ?? [],
    output_fields: action.output_fields ?? [],
  };
}

function mapReactionConfig(reaction: ServiceReactionConfig): ReactionDefinition {
  return {
    id: reaction.title,
    title: reaction.title,
    label: reaction.label ?? reaction.title,
    url: reaction.url,
    method: reaction.method,
    fields: reaction.fields ?? [],
  };
}

function mapServiceConfig(
  serviceId: string,
  config: ServiceConfig,
  index: number,
): AreaService {
  const gradient = gradientPalette[index % gradientPalette.length] ?? {
    from: "#0b3c5d",
    to: "#e59500",
  };
  const name = config.label ?? formatServiceNameFromId(serviceId) ?? serviceId;
  const actions = (config.actions ?? []).map(mapActionConfig);
  const reactions = (config.reactions ?? []).map(mapReactionConfig);

  return {
    id: serviceId,
    name,
    provider: config.provider ?? "",
    badge: buildBadge(name || serviceId),
    logoUrl: config.logo_url || config.icon_url || "",
    gradient,
    actions,
    reactions,
    connected: false,
  };
}

function matchesAreaSearch(area: CreatedArea, normalizedTerm: string) {
  if (!normalizedTerm) return true;
  const haystack = normalizeSearchValue(
    [area.name, area.summary, area.serviceName, area.startTime, area.endTime].join(" "),
  );
  return haystack.includes(normalizedTerm);
}

function initializeFieldValues(fields: FieldDefinition[]): Record<string, FieldValue> {
  return fields.reduce<Record<string, FieldValue>>((acc, field) => {
    const defaultValue = field.default ?? "";
    if (field.selection?.length && field.multiple) {
      if (typeof defaultValue === "string") {
        const trimmed = defaultValue.trim();
        if (!trimmed) {
          acc[field.name] = [];
        } else if (trimmed.startsWith("[")) {
          try {
            const parsed = JSON.parse(trimmed);
            acc[field.name] = Array.isArray(parsed) ? parsed.map(String) : [trimmed];
          } catch {
            acc[field.name] = trimmed.split(",").map((item) => item.trim()).filter(Boolean);
          }
        } else {
          acc[field.name] = trimmed.split(",").map((item) => item.trim()).filter(Boolean);
        }
      } else if (typeof defaultValue === "number") {
        acc[field.name] = [String(defaultValue)];
      } else {
        acc[field.name] = [];
      }
      return acc;
    }

    acc[field.name] = typeof defaultValue === "number" ? String(defaultValue) : String(defaultValue);
    return acc;
  }, {});
}

function areRequiredFieldsFilled(
  fields: FieldDefinition[],
  values: Record<string, FieldValue>,
): boolean {
  return fields.every((field) => {
    if (!field.required) return true;
    const value = values[field.name];
    if (Array.isArray(value)) {
      return value.length > 0;
    }
    return value !== undefined && String(value).trim().length > 0;
  });
}

function toIsoString(value: string): string {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    throw new Error("Date/heure invalide.");
  }
  return date.toISOString();
}

function formatInputValue(
  field: FieldDefinition,
  rawValue: FieldValue | undefined,
): string {
  const value = rawValue ?? "";
  if (Array.isArray(value)) {
    return value.length ? JSON.stringify(value) : "";
  }
  if (!value) return "";
  if (field.type === "date") {
    return toIsoString(value);
  }
  return String(value);
}

function formatFieldDisplayValue(value: FieldValue | undefined): string {
  if (Array.isArray(value)) {
    return value.join(", ");
  }
  return value ? String(value) : "";
}

function buildInputFields(
  fields: FieldDefinition[],
  values: Record<string, FieldValue>,
): Array<{ name: string; value: string }> {
  return fields.map((field) => ({
    name: field.name,
    value: formatInputValue(field, values[field.name]),
  }));
}

type BackendAction = BackendArea["actions"][number];
type BackendReaction = BackendArea["reactions"][number];

function inputFieldsToRecord(
  input?: Array<{ name: string; value: string }>,
): Record<string, string> {
  if (!input?.length) return {};
  return input.reduce<Record<string, string>>((acc, field) => {
    acc[field.name] = field.value;
    return acc;
  }, {});
}

function resolveServiceById(
  services: AreaService[],
  serviceId?: string,
): AreaService | undefined {
  if (!serviceId) return undefined;
  return services.find((service) => service.id === serviceId);
}

function resolveServiceName(
  services: AreaService[],
  serviceId?: string,
): string {
  if (!serviceId) return "";
  return (
    resolveServiceById(services, serviceId)?.name ??
    formatServiceNameFromId(serviceId) ??
    serviceId
  );
}

function getServiceVisual(
  services: AreaService[],
  serviceId?: string,
  fallbackName?: string,
): { badge: string; logoUrl?: string } {
  const service = resolveServiceById(services, serviceId);
  const badgeSource = service?.name ?? serviceId ?? fallbackName ?? "?";
  return {
    logoUrl: service?.logoUrl,
    badge: service?.badge ?? buildBadge(badgeSource),
  };
}

function resolveActionLabel(
  action: BackendAction | undefined,
  services: AreaService[],
): string {
  if (!action) return "Action";
  const service = resolveServiceById(services, action.service);
  const config = service?.actions.find((item) => item.title === action.title);
  return config?.label ?? action.title ?? "Action";
}

function resolveReactionLabel(
  reaction: BackendReaction | undefined,
  services: AreaService[],
): string {
  if (!reaction) return "Réaction";
  const service = resolveServiceById(services, reaction.service);
  const config = service?.reactions.find((item) => item.title === reaction.title);
  return config?.label ?? reaction.title ?? "Réaction";
}

function pickGradientForKey(key: string): AreaGradient {
  if (!gradientPalette.length) {
    return { from: "#0b3c5d", to: "#e59500" };
  }
  let hash = 0;
  for (let i = 0; i < key.length; i += 1) {
    hash = (hash * 31 + key.charCodeAt(i)) | 0;
  }
  const index = Math.abs(hash) % gradientPalette.length;
  return gradientPalette[index] ?? gradientPalette[0] ?? { from: "#0b3c5d", to: "#e59500" };
}

function resolveGradient(
  services: AreaService[],
  key?: string,
): AreaGradient {
  const service = resolveServiceById(services, key);
  if (service?.gradient) return service.gradient;
  return pickGradientForKey(key ?? "");
}

function mapBackendArea(area: BackendArea, services: AreaService[]): CreatedArea {
  const action = area.actions?.[0];
  const reaction = area.reactions?.[0];
  const actionInputs = inputFieldsToRecord(action?.input);
  const reactionInputs = inputFieldsToRecord(reaction?.input);
  const actionServiceName = resolveServiceName(services, action?.service);
  const reactionServiceName = resolveServiceName(services, reaction?.service);
  const actionServiceId = action?.service ?? "";
  const reactionServiceIds = (area.reactions ?? [])
    .map((item) => item.service)
    .filter((serviceId): serviceId is string => Boolean(serviceId));
  const primaryReactionServiceId = reaction?.service ?? reactionServiceIds[0] ?? "";
  const delayValue = Number.parseInt(actionInputs.delay ?? "0", 10);
  const delay = Number.isFinite(delayValue) ? delayValue : 0;
  const summary = (reactionInputs.summary ?? "").trim() || area.name;
  const reactionCount = area.reactions?.length ?? 0;

  return {
    id: String(area.id),
    name: area.name,
    summary,
    startTime: reactionInputs.start_time ?? "",
    endTime: reactionInputs.end_time ?? "",
    delay,
    serviceName: actionServiceName || reactionServiceName || area.name,
    actionServiceId,
    actionService: actionServiceName,
    primaryReactionServiceId,
    reactionServiceIds,
    reactionService: reactionServiceName,
    reactionCount,
    actionName: resolveActionLabel(action, services),
    reactionName: resolveReactionLabel(reaction, services),
    gradient: resolveGradient(services, action?.service ?? reaction?.service ?? area.name),
    active: area.active,
  };
}

function AreaPageContent() {
  const { token, user } = useAuth();
  const searchParams = useSearchParams();
  const hasOAuthParams = Boolean(searchParams.get("code") && searchParams.get("state"));
  const { status, error } = useOAuthCallback("/area", { enabled: hasOAuthParams });
  const isProcessingOAuth = hasOAuthParams && status !== "error";
  const [rawAreas, setRawAreas] = useState<BackendArea[]>([]);
  const [areasError, setAreasError] = useState<string | null>(null);
  const [isLoadingAreas, setIsLoadingAreas] = useState(false);
  const [services, setServices] = useState<AreaService[]>([]);
  const [servicesError, setServicesError] = useState<string | null>(null);
  const [isLoadingServices, setIsLoadingServices] = useState(false);
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [actionService, setActionService] = useState<AreaService | null>(null);
  const [selectedAction, setSelectedAction] = useState<ActionDefinition | null>(null);
  const [actionFieldValues, setActionFieldValues] = useState<Record<string, FieldValue>>({});
  const [reactions, setReactions] = useState<ReactionForm[]>([]);
  const [areaName, setAreaName] = useState("");
  const [wizardStep, setWizardStep] = useState<"action" | "reactions" | "details">("action");
  const [expandedReactionId, setExpandedReactionId] = useState<string | null>(null);
  const [selectedAreaDetail, setSelectedAreaDetail] = useState<CreatedArea | null>(null);
  const [createError, setCreateError] = useState<string | null>(null);
  const [areaActionError, setAreaActionError] = useState<string | null>(null);
  const [isCreating, setIsCreating] = useState(false);
  const [updatingAreaId, setUpdatingAreaId] = useState<string | null>(null);
  const selectRefs = useRef<Record<string, HTMLDetailsElement | null>>({});
  const reactionIdRef = useRef(0);
  const [searchTerm, setSearchTerm] = useState("");
  const normalizedSearch = normalizeSearchValue(searchTerm);
  const displayAreas = useMemo(() => {
    const sortedAreas = [...rawAreas].sort((a, b) => b.id - a.id);
    return sortedAreas.map((area) => mapBackendArea(area, services));
  }, [rawAreas, services]);
  const filteredAreas = displayAreas.filter((area) => matchesAreaSearch(area, normalizedSearch));
  const hasSearch = Boolean(normalizedSearch);
  const totalAreas = displayAreas.length;
  const activeCount = displayAreas.filter((area) => area.active).length;
  const hasConnectedServices = useMemo(
    () => services.some((service) => service.connected),
    [services],
  );

  const renderFieldInputs = (
    fields: FieldDefinition[],
    values: Record<string, FieldValue>,
    onChange: (name: string, value: FieldValue) => void,
    fieldKeyPrefix = "",
    tokenOptions?: Array<{ token: string; label: string }>,
    onInsertToken?: (name: string, token: string) => void,
  ) => (
    <div className="grid gap-3 sm:grid-cols-2">
      {fields.map((field) => {
        const rawValue = values[field.name];
        const value = Array.isArray(rawValue) ? rawValue.join(", ") : rawValue ?? "";
        const selectKey = fieldKeyPrefix ? `${fieldKeyPrefix}-${field.name}` : field.name;
        const baseClasses =
          "w-full rounded-xl border border-[var(--surface-border)] bg-[var(--background)] px-3 py-2 text-[var(--foreground)] focus:border-[var(--blue-primary-3)] focus:outline-none focus:ring-2 focus:ring-[var(--blue-primary-3)]/25";

        const selectionValues = Array.isArray(rawValue)
          ? rawValue
          : rawValue
            ? [String(rawValue)]
            : [];
        const selectionLabels = selectionValues
          .map((current) => field.selection?.find((option) => option.value === current)?.label ?? current)
          .filter(Boolean);
        const selectionPlaceholder = selectionLabels.length
          ? selectionLabels.join(", ")
          : "Sélectionner...";

        if (field.selection?.length) {
          return (
            <label key={field.name} className="space-y-1 text-sm">
              <span className="text-[var(--muted)]">
                {field.label}
                {field.required ? " *" : ""}
              </span>
              <details
                ref={(node) => {
                  selectRefs.current[selectKey] = node;
                }}
                className="group relative"
              >
                <summary className="flex cursor-pointer list-none items-center justify-between gap-3 rounded-xl border border-[var(--surface-border)] bg-[var(--background)] px-4 py-2.5 text-[var(--foreground)] shadow-[0_1px_0_rgba(0,0,0,0.04)] focus:outline-none focus-visible:ring-2 focus-visible:ring-[var(--blue-primary-3)]/25">
                  <span className={selectionLabels.length ? "text-[var(--foreground)]" : "text-[var(--muted)]"}>
                    {selectionPlaceholder}
                  </span>
                  <span className="text-[var(--muted)] transition group-open:rotate-180">
                    <svg viewBox="0 0 20 20" className="h-4 w-4" fill="none" aria-hidden>
                      <path
                        d="m5 7.5 5 5 5-5"
                        stroke="currentColor"
                        strokeWidth="1.6"
                        strokeLinecap="round"
                        strokeLinejoin="round"
                      />
                    </svg>
                  </span>
                </summary>
                <div className="absolute z-20 mt-2 w-full overflow-hidden rounded-xl border border-[var(--surface-border)] bg-[var(--background)] shadow-[0_12px_30px_rgba(15,25,35,0.12)]">
                  <div className="max-h-56 overflow-y-auto p-1.5">
                    {field.selection.map((option) => {
                      const isSelected = selectionValues.includes(option.value);
                      return (
                        <button
                          key={`${field.name}-${option.value}`}
                          type="button"
                          onClick={() => {
                            if (field.multiple) {
                              const nextValues = isSelected
                                ? selectionValues.filter((item) => item !== option.value)
                                : [...selectionValues, option.value];
                              onChange(field.name, nextValues);
                            } else {
                              onChange(field.name, option.value);
                              selectRefs.current[field.name]?.removeAttribute("open");
                            }
                          }}
                          className={cn(
                            "flex w-full items-center justify-between rounded-lg px-3 py-2 text-left text-sm transition",
                            isSelected
                              ? "bg-[var(--blue-primary-3)]/10 text-[var(--foreground)]"
                              : "text-[var(--foreground)] hover:bg-[var(--surface)]",
                          )}
                        >
                          <span>{option.label}</span>
                          {isSelected ? (
                            <span className="text-[var(--blue-primary-3)]">
                              <svg viewBox="0 0 20 20" className="h-4 w-4" fill="none" aria-hidden>
                                <path
                                  d="m5 10 3 3 7-7"
                                  stroke="currentColor"
                                  strokeWidth="1.8"
                                  strokeLinecap="round"
                                  strokeLinejoin="round"
                                />
                              </svg>
                            </span>
                          ) : null}
                        </button>
                      );
                    })}
                  </div>
                  {field.multiple ? (
                    <div className="border-t border-[var(--surface-border)] px-3 py-2 text-xs text-[var(--muted)]">
                      Cliquez pour ajouter ou retirer des options.
                    </div>
                  ) : null}
                </div>
              </details>
            </label>
          );
        }

        if (field.type === "number") {
          return (
            <label key={field.name} className="space-y-1 text-sm">
              <span className="text-[var(--muted)]">
                {field.label}
                {field.required ? " *" : ""}
              </span>
              <input
                type="number"
                min={0}
                value={value}
                onChange={(e) => onChange(field.name, e.target.value)}
                className={baseClasses}
              />
            </label>
          );
        }

        if (field.type === "date") {
          return (
            <label key={field.name} className="space-y-1 text-sm">
              <span className="text-[var(--muted)]">
                {field.label}
                {field.required ? " *" : ""}
              </span>
              <input
                type="datetime-local"
                value={value}
                onChange={(e) => onChange(field.name, e.target.value)}
                className={baseClasses}
              />
            </label>
          );
        }

        const useTextArea = field.name.toLowerCase().includes("description");

        return (
          <label key={field.name} className="space-y-1 text-sm sm:col-span-2">
            <span className="text-[var(--muted)]">
              {field.label}
              {field.required ? " *" : ""}
            </span>
            {useTextArea ? (
              <>
                <textarea
                  value={value}
                  onChange={(e) => onChange(field.name, e.target.value)}
                  rows={4}
                  className={baseClasses}
                />
                {tokenOptions?.length && onInsertToken ? (
                  <div className="flex flex-wrap gap-2">
                    {tokenOptions.map((token) => (
                      <button
                        key={`${field.name}-${token.token}`}
                        type="button"
                        onClick={() => onInsertToken(field.name, token.token)}
                        className="inline-flex items-center gap-1 rounded-full border border-[var(--surface-border)] bg-[var(--surface)] px-3 py-1 text-xs font-semibold text-[var(--muted)] transition hover:border-[var(--blue-primary-2)] hover:text-[var(--foreground)]"
                      >
                        <span className="text-[var(--blue-primary-3)]">+</span>
                        <span>{token.label}</span>
                      </button>
                    ))}
                  </div>
                ) : null}
              </>
            ) : (
              <>
                <input
                  type="text"
                  value={value}
                  onChange={(e) => onChange(field.name, e.target.value)}
                  className={baseClasses}
                />
                {tokenOptions?.length && onInsertToken ? (
                  <div className="flex flex-wrap gap-2">
                    {tokenOptions.map((token) => (
                      <button
                        key={`${field.name}-${token.token}`}
                        type="button"
                        onClick={() => onInsertToken(field.name, token.token)}
                        className="inline-flex items-center gap-1 rounded-full border border-[var(--surface-border)] bg-[var(--surface)] px-3 py-1 text-xs font-semibold text-[var(--muted)] transition hover:border-[var(--blue-primary-2)] hover:text-[var(--foreground)]"
                      >
                        <span className="text-[var(--blue-primary-3)]">+</span>
                        <span>{token.label}</span>
                      </button>
                    ))}
                  </div>
                ) : null}
              </>
            )}
          </label>
        );
      })}
    </div>
  );

  const loadServices = useCallback(async () => {
    setIsLoadingServices(true);
    setServicesError(null);

    try {
      const serviceIds = await fetchServiceNames();
      const uniqueServiceIds = Array.from(new Set(serviceIds.filter(Boolean)));

      const configResults = await Promise.allSettled(
        uniqueServiceIds.map((serviceId) => fetchServiceConfig(serviceId)),
      );
      const configs = configResults.map((result) =>
        result.status === "fulfilled" ? result.value : null,
      );
      const availableServices = uniqueServiceIds
        .map((serviceId, index) => {
          const config = configs[index];
          if (!config) return null;
          return mapServiceConfig(serviceId, config, index);
        })
        .filter(Boolean) as AreaService[];

      if (!availableServices.length) {
        throw new Error("Impossible de récupérer la configuration des services.");
      }

      if (configResults.some((result) => result.status === "rejected")) {
        setServicesError("Certains services ne sont pas disponibles pour le moment.");
      }

      let statusByService: Record<string, { is_logged: boolean; logo_url?: string }> = {};
      if (token && user?.id) {
        const statuses = await fetchUserServiceStatuses(token, user.id);
        statusByService = statuses.reduce<
          Record<string, { is_logged: boolean; logo_url?: string }>
        >((acc, current) => {
          acc[current.provider] = {
            is_logged: Boolean(current.is_logged),
            logo_url: current.logo_url,
          };
          return acc;
        }, {});
      }

      const mappedServices = availableServices.map((service) => {
        const providerKey = service.provider;
        const providerStatus = providerKey ? statusByService[providerKey] : undefined;
        const isInternalService = !providerKey;
        const isConnected = isInternalService
          ? true
          : Boolean(providerStatus?.is_logged);
        const logoUrl = service.logoUrl || providerStatus?.logo_url;

        return {
          ...service,
          connected: isConnected,
          logoUrl,
        };
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

  const loadAreas = useCallback(async () => {
    if (!token) {
      setRawAreas([]);
      setAreasError(null);
      return;
    }

    setIsLoadingAreas(true);
    setAreasError(null);

    try {
      const areas = await fetchAreas(token);
      setRawAreas(areas);
      setAreaActionError(null);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Impossible de récupérer les areas.";
      setAreasError(message);
    } finally {
      setIsLoadingAreas(false);
    }
  }, [token]);

  useEffect(() => {
    void loadAreas();
  }, [loadAreas]);

  const handleToggleArea = useCallback(
    async (areaId: string, targetActive: boolean) => {
      if (!token) {
        setAreaActionError("Vous devez être connecté pour modifier une area.");
        return;
      }

      const numericId = Number.parseInt(areaId, 10);
      if (!Number.isFinite(numericId)) {
        setAreaActionError("Identifiant d'area invalide.");
        return;
      }

      setAreaActionError(null);
      setUpdatingAreaId(areaId);
      try {
        if (targetActive) {
          await activateArea(token, numericId);
        } else {
          await deactivateArea(token, numericId);
        }
        await loadAreas();
      } catch (err) {
        const message =
          err instanceof Error ? err.message : "Impossible de mettre à jour l'area.";
        setAreaActionError(message);
      } finally {
        setUpdatingAreaId(null);
      }
    },
    [token, loadAreas],
  );

  const handleDeleteArea = useCallback(
    async (areaId: string, areaName: string) => {
      if (!token) {
        setAreaActionError("Vous devez être connecté pour supprimer une area.");
        return;
      }
      const numericId = Number.parseInt(areaId, 10);
      if (!Number.isFinite(numericId)) {
        setAreaActionError("Identifiant d'area invalide.");
        return;
      }

      setAreaActionError(null);
      setUpdatingAreaId(areaId);
      try {
        await deleteArea(token, numericId);
        await loadAreas();
      } catch (err) {
        const message =
          err instanceof Error ? err.message : `Impossible de supprimer "${areaName}".`;
        setAreaActionError(message);
      } finally {
        setUpdatingAreaId(null);
      }
    },
    [token, loadAreas],
  );

  const connectedServices = useMemo(
    () => services.filter((service) => service.connected),
    [services],
  );
  const actionServices = useMemo(
    () => connectedServices.filter((service) => service.actions.length > 0),
    [connectedServices],
  );
  const reactionServices = useMemo(
    () => connectedServices.filter((service) => service.reactions.length > 0),
    [connectedServices],
  );

  useEffect(() => {
    setCreateError(null);
  }, [wizardStep]);

  const reactionIsValid = (reaction: ReactionForm) =>
    Boolean(
      reaction.service &&
        reaction.service.connected &&
        reaction.reaction &&
        areRequiredFieldsFilled(reaction.reaction.fields, reaction.fieldValues),
    );
  const allReactionsValid = reactions.length > 0 && reactions.every((reaction) => reactionIsValid(reaction));
  const actionOutputTokens = useMemo(
    () =>
      (selectedAction?.output_fields ?? []).map((field) => ({
        token: `{{${field.name}}}`,
        label: field.label ?? field.name,
      })),
    [selectedAction],
  );
  const canProceedAction =
    Boolean(actionService && actionService.connected && selectedAction) &&
    areRequiredFieldsFilled(selectedAction?.fields ?? [], actionFieldValues);
  const canProceedReactions = allReactionsValid;
  const wizardSteps: Array<{ id: "action" | "reactions" | "details"; title: string; description: string }> = [
    { id: "action", title: "Action", description: "Déclencheur" },
    { id: "reactions", title: "Réactions", description: "Actions exécutées" },
    { id: "details", title: "Détails", description: "Récapitulatif" },
  ];
  const currentStepIndex = wizardSteps.findIndex((step) => step.id === wizardStep);

  const canCreate =
    Boolean(
      actionService &&
        actionService.connected &&
        selectedAction &&
        areaName.trim() &&
        areRequiredFieldsFilled(selectedAction.fields, actionFieldValues) &&
        allReactionsValid,
    ) && hasConnectedServices;

  const goToNextStep = () => {
    setCreateError(null);
    if (wizardStep === "action") {
      if (canProceedAction) {
        setWizardStep("reactions");
      } else {
        setCreateError("Complétez le déclencheur avant de continuer.");
      }
      return;
    }
    if (wizardStep === "reactions") {
      if (canProceedReactions) {
        setWizardStep("details");
      } else {
        setCreateError("Ajoutez au moins une réaction complète avant de continuer.");
      }
    }
  };

  const goToPreviousStep = () => {
    setCreateError(null);
    if (wizardStep === "reactions") {
      setWizardStep("action");
      return;
    }
    if (wizardStep === "details") {
      setWizardStep("reactions");
    }
  };

  const resetForm = () => {
    setActionService(null);
    setSelectedAction(null);
    setActionFieldValues({});
    setReactions([]);
    reactionIdRef.current = 0;
    setAreaName("");
    setCreateError(null);
    setWizardStep("action");
    setExpandedReactionId(null);
  };

  const openModal = () => {
    resetForm();
    setIsCreateModalOpen(true);
  };

  const closeModal = () => {
    setIsCreateModalOpen(false);
  };

  const clearCreateError = () => setCreateError(null);

  const addReaction = () => {
    const newId = `reaction-${reactionIdRef.current++}`;
    setReactions((prev) => [
      ...prev,
      { id: newId, service: null, reaction: null, fieldValues: {} },
    ]);
    setExpandedReactionId(newId);
    clearCreateError();
  };

  const removeReaction = (reactionId: string) => {
    setReactions((prev) => {
      const nextReactions = prev.filter((reaction) => reaction.id !== reactionId);
      setExpandedReactionId((current) => (current === reactionId ? nextReactions[0]?.id ?? null : current));
      return nextReactions;
    });
    clearCreateError();
  };

  const updateReactionService = (reactionId: string, service: AreaService) => {
    setReactions((prev) =>
      prev.map((reaction) =>
        reaction.id === reactionId
          ? { ...reaction, service, reaction: null, fieldValues: {} }
          : reaction,
      ),
    );
    clearCreateError();
  };

  const updateReactionSelection = (reactionId: string, reactionDef: ReactionDefinition) => {
    setReactions((prev) =>
      prev.map((reaction) =>
        reaction.id === reactionId
          ? { ...reaction, reaction: reactionDef, fieldValues: initializeFieldValues(reactionDef.fields) }
          : reaction,
      ),
    );
    clearCreateError();
  };

  const updateReactionField = (reactionId: string, name: string, value: FieldValue) => {
    setReactions((prev) =>
      prev.map((reaction) =>
        reaction.id === reactionId
          ? { ...reaction, fieldValues: { ...reaction.fieldValues, [name]: value } }
          : reaction,
      ),
    );
  };

  const insertTokenIntoReactionField = (reactionId: string, name: string, token: string) => {
    setReactions((prev) =>
      prev.map((reaction) => {
        if (reaction.id !== reactionId) return reaction;
        const current = reaction.fieldValues[name];
        const asString = Array.isArray(current) ? current.join(", ") : current ? String(current) : "";
        const nextValue = asString ? `${asString} ${token}` : token;
        return { ...reaction, fieldValues: { ...reaction.fieldValues, [name]: nextValue } };
      }),
    );
    clearCreateError();
  };

  const handleCreateArea = async () => {
    if (!token) {
      setCreateError("Vous devez être connecté pour créer une area.");
      return;
    }
    if (!actionService) {
      setCreateError("Sélectionnez un service pour le déclencheur.");
      return;
    }
    if (!selectedAction) {
      setCreateError("Choisissez un déclencheur.");
      return;
    }
    if (!reactions.length) {
      setCreateError("Ajoutez au moins une réaction.");
      return;
    }
    if (!areaName.trim()) {
      setCreateError("Renseignez au moins le nom de l'area.");
      return;
    }
    if (
      !areRequiredFieldsFilled(selectedAction.fields, actionFieldValues) ||
      reactions.some((reaction) => !reactionIsValid(reaction))
    ) {
      setCreateError("Complétez tous les champs obligatoires du déclencheur et des réactions.");
      return;
    }

    setIsCreating(true);
    setCreateError(null);

    try {
      const actionInputs = buildInputFields(selectedAction.fields, actionFieldValues);
      const reactionPayloads = reactions.map((reaction) => {
        const reactionInputs = buildInputFields(reaction.reaction?.fields ?? [], reaction.fieldValues);
        return {
          service: reaction.service?.id ?? "",
          provider: reaction.service?.provider || reaction.service?.id || "",
          title: reaction.reaction?.title ?? "",
          input: reactionInputs,
        };
      });

      await saveArea(token, {
        name: areaName.trim(),
        active: true,
        actions: [
          {
            service: actionService.id,
            provider: actionService.provider || actionService.id,
            title: selectedAction.title,
            type: selectedAction.type,
            input: actionInputs,
          },
        ],
        reactions: reactionPayloads,
      });

      await loadAreas();
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
          className="fixed inset-0 z-50 flex items-center justify-center bg-[rgba(6,14,25,0.55)] px-6 py-12 backdrop-blur-sm"
          onClick={closeModal}
        >
          <div
            className="relative w-full max-w-6xl max-h-[90vh] overflow-y-auto rounded-[28px] border border-[var(--surface-border)] bg-[var(--background)] shadow-2xl ring-1 ring-[rgba(28,61,99,0.28)]"
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

            <div className="space-y-7 px-10 pb-11 pt-9">
              <div className="space-y-2">
                <p className="text-xs font-semibold uppercase tracking-[0.14em] text-[var(--blue-primary-3)]">
                  Création d&apos;area
                </p>
                <div className="flex flex-wrap items-end justify-between gap-3">
                  <div className="space-y-1">
                    <h2 id="create-area-title" className="text-2xl font-semibold text-[var(--foreground)]">
                      Composez votre automation
                    </h2>
                    <p className="text-sm text-[var(--muted)]">Sélectionnez un déclencheur et une réaction pour structurer l&apos;area.</p>
                  </div>
                </div>
              </div>

              <div className="space-y-5">
                {servicesError ? (
                  <div className="rounded-xl border border-[var(--accent)] bg-[var(--surface)] px-4 py-3 text-sm text-[var(--accent)]">
                    {servicesError}
                  </div>
                ) : null}

                <div className="flex flex-wrap gap-3 rounded-2xl border border-[var(--surface-border)] bg-[var(--surface)] p-4 sm:p-5">
                  {wizardSteps.map((step, index) => {
                    const isActive = step.id === wizardStep;
                    const isDone = currentStepIndex > index;
                    return (
                      <div
                        key={step.id}
                        className={cn(
                          "flex items-center gap-3 rounded-xl border px-3 py-2 text-sm transition",
                          isActive
                            ? "border-[var(--blue-primary-3)] bg-[var(--blue-primary-3)]/10 text-[var(--foreground)] shadow-[0_0_0_2px_rgba(28,61,99,0.12)]"
                            : "border-[var(--surface-border)] bg-[var(--background)] text-[var(--muted)]",
                        )}
                      >
                        <span
                          className={cn(
                            "flex h-7 w-7 items-center justify-center rounded-full border text-xs font-semibold",
                            isActive || isDone
                              ? "border-[var(--blue-primary-3)] bg-[var(--blue-primary-3)]/10 text-[var(--blue-primary-3)]"
                              : "border-[var(--surface-border)] text-[var(--muted)]",
                          )}
                        >
                          {index + 1}
                        </span>
                        <div className="leading-tight">
                          <p className="font-semibold">{step.title}</p>
                          <p className="text-[11px] text-[var(--muted)]">{step.description}</p>
                        </div>
                      </div>
                    );
                  })}
                </div>

                {wizardStep === "action" ? (
                  <div className="rounded-2xl border border-[var(--surface-border)] bg-[var(--surface)] p-5 sm:p-6">
                    <div className="space-y-1">
                      <h3 className="text-lg font-semibold text-[var(--foreground)]">Action (Déclencheur)</h3>
                      <p className="text-sm text-[var(--muted)]">
                        Choisissez d&apos;abord le service, puis l&apos;événement qui déclenchera l&apos;automation.
                      </p>
                    </div>

                    <div className="mt-5 space-y-5">
                      <div className="space-y-3">
                        <p className="text-xs font-semibold uppercase tracking-[0.14em] text-[var(--muted)]">Service</p>
                        {isLoadingServices ? (
                          <p className="text-sm text-[var(--muted)]">Chargement des services...</p>
                        ) : actionServices.length ? (
                          <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5">
                            {actionServices.map((service) => (
                              <button
                                key={service.id}
                                type="button"
                                onClick={() => {
                                clearCreateError();
                                setActionService(service);
                                setSelectedAction(null);
                                setActionFieldValues({});
                              }}
                              className={cn(
                                  "flex h-14 items-center justify-start gap-3 rounded-xl border px-4 text-sm font-semibold transition",
                                  actionService?.id === service.id
                                    ? "border-[var(--blue-primary-3)] bg-[var(--blue-primary-3)]/10 text-[var(--foreground)] shadow-[0_0_0_2px_rgba(28,61,99,0.12)]"
                                    : "border-[var(--surface-border)] bg-[var(--background)] hover:-translate-y-0.5 hover:border-[var(--blue-primary-2)]",
                                )}
                              >
                                <span className="flex h-10 w-10 items-center justify-center rounded-full border border-[var(--surface-border)] bg-[var(--surface)] text-xs font-semibold uppercase text-[var(--foreground)]">
                                  {service.logoUrl ? (
                                    <img
                                      src={service.logoUrl}
                                      alt=""
                                      className="h-6 w-6 object-contain"
                                      loading="lazy"
                                    />
                                  ) : (
                                    service.badge
                                  )}
                                </span>
                                <span className="truncate">{service.name}</span>
                              </button>
                            ))}
                          </div>
                        ) : (
                          <p className="text-sm text-[var(--muted)]">
                            Aucun service avec déclencheur disponible. Connectez un service compatible pour commencer.
                          </p>
                        )}
                      </div>

                      <div className="space-y-3">
                        <p className="text-xs font-semibold uppercase tracking-[0.14em] text-[var(--muted)]">
                          Déclencheur
                        </p>
                        {actionService && actionService.actions?.length ? (
                          <div className="space-y-3">
                          {actionService.actions.map((action) => {
                            const isSelected = selectedAction?.id === action.id;
                            return (
                              <button
                                key={action.id}
                                type="button"
                                onClick={() => {
                                  clearCreateError();
                                  setSelectedAction(action);
                                  setActionFieldValues(initializeFieldValues(action.fields));
                                  setExpandedReactionId(null);
                                }}
                                  className={cn(
                                    "flex w-full items-center justify-between rounded-xl border px-4 py-3 text-left text-sm transition",
                                    isSelected
                                      ? "border-[var(--blue-primary-3)] bg-[var(--blue-primary-3)]/10 text-[var(--foreground)] shadow-[0_0_0_2px_rgba(28,61,99,0.12)]"
                                      : "border-[var(--surface-border)] bg-[var(--background)] hover:border-[var(--blue-primary-2)]",
                                  )}
                                >
                                  <span>{action.label}</span>
                                  {isSelected ? <span className="text-[var(--blue-primary-3)]">●</span> : null}
                                </button>
                              );
                            })}
                          </div>
                        ) : (
                          <div className="rounded-xl border border-dashed border-[var(--surface-border)] bg-[var(--background)] px-4 py-3 text-sm text-[var(--muted)]">
                            Choisissez un service pour afficher ses déclencheurs.
                          </div>
                        )}
                      </div>
                      {selectedAction ? (
                        <div className="space-y-3">
                          <p className="text-xs font-semibold uppercase tracking-[0.14em] text-[var(--muted)]">
                            Paramètres du déclencheur
                          </p>
                          <div className="rounded-xl border border-[var(--surface-border)] bg-[var(--background)] p-4">
                            {renderFieldInputs(
                              selectedAction.fields,
                              actionFieldValues,
                              (name, value) => {
                                clearCreateError();
                                setActionFieldValues((prev) => ({ ...prev, [name]: value }));
                              },
                              "action",
                            )}
                          </div>
                        </div>
                      ) : null}
                    </div>
                  </div>
                ) : null}

                {wizardStep === "reactions" ? (
                  <div className="rounded-2xl border border-[var(--surface-border)] bg-[var(--surface)] p-5 sm:p-6">
                    <div className="space-y-1">
                      <h3 className="text-lg font-semibold text-[var(--foreground)]">Réactions</h3>
                      <p className="text-sm text-[var(--muted)]">
                        Ajoutez une ou plusieurs actions exécutées après le déclencheur.
                      </p>
                    </div>

                    <div className="mt-5 space-y-4">
                  {isLoadingServices ? (
                    <p className="text-sm text-[var(--muted)]">Chargement des services...</p>
                  ) : null}
                  {!isLoadingServices && !reactionServices.length ? (
                    <p className="text-sm text-[var(--muted)]">
                          Aucun service avec réaction disponible. Connectez un service compatible pour commencer.
                        </p>
                      ) : null}
                      {!isLoadingServices && reactionServices.length ? (
                        reactions.length ? (
                          <div className="space-y-4">
                            {reactions.map((reaction, index) => (
                              <div
                                key={reaction.id}
                                className="space-y-3 rounded-xl border border-[var(--surface-border)] bg-[var(--background)] p-4 shadow-[0_6px_18px_rgba(12,24,46,0.06)]"
                              >
                                <div className="flex flex-wrap items-center justify-between gap-2">
                                  <button
                                    type="button"
                                    onClick={() =>
                                      setExpandedReactionId((current) =>
                                        current === reaction.id ? null : reaction.id,
                                      )
                                    }
                                    className="flex flex-1 items-start gap-3 text-left"
                                  >
                                    <div className="min-w-0">
                                      <p className="text-xs font-semibold uppercase tracking-[0.12em] text-[var(--muted)]">
                                        Réaction {index + 1}
                                      </p>
                                      <p className="text-sm font-semibold text-[var(--foreground)] truncate">
                                        {reaction.service?.name ?? "Sélectionnez un service"}
                                      </p>
                                      <p className="text-[var(--muted)] text-xs truncate">
                                        {reaction.reaction?.label ?? "Aucune action sélectionnée"}
                                      </p>
                                    </div>
                                    <span
                                      className={cn(
                                        "rounded-full px-3 py-1 text-xs font-semibold",
                                        reactionIsValid(reaction)
                                          ? "bg-[var(--success,#22c55e)]/15 text-[var(--success,#22c55e)]"
                                          : "bg-[var(--danger,#ef4444)]/15 text-[var(--danger,#ef4444)]",
                                      )}
                                    >
                                      {reactionIsValid(reaction) ? "Prête" : "Incomplète"}
                                    </span>
                                  </button>
                                  <div className="flex items-center gap-2">
                                    <button
                                      type="button"
                                      onClick={() =>
                                        setExpandedReactionId((current) =>
                                          current === reaction.id ? null : reaction.id,
                                        )
                                      }
                                      className="inline-flex h-9 w-9 items-center justify-center rounded-full border border-[var(--surface-border)] bg-[var(--surface)] text-[var(--muted)] transition hover:text-[var(--foreground)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)]"
                                      aria-label={`Réduire ou développer la réaction ${index + 1}`}
                                    >
                                      <svg
                                        className={cn(
                                          "h-4 w-4 transition",
                                          expandedReactionId === reaction.id ? "rotate-180" : "",
                                        )}
                                        viewBox="0 0 20 20"
                                        fill="none"
                                        aria-hidden
                                      >
                                        <path
                                          d="m5 7.5 5 5 5-5"
                                          stroke="currentColor"
                                          strokeWidth="1.6"
                                          strokeLinecap="round"
                                          strokeLinejoin="round"
                                        />
                                      </svg>
                                    </button>
                                    <button
                                      type="button"
                                      onClick={() => removeReaction(reaction.id)}
                                      className="inline-flex h-9 w-9 items-center justify-center rounded-full border border-[var(--surface-border)] bg-[var(--surface)] text-[var(--muted)] transition hover:text-[var(--foreground)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)]"
                                      aria-label={`Supprimer la réaction ${index + 1}`}
                                    >
                                      <svg className="h-4 w-4" viewBox="0 0 20 20" fill="none" aria-hidden>
                                        <path d="m5 5 10 10M15 5 5 15" stroke="currentColor" strokeWidth="1.6" />
                                      </svg>
                                    </button>
                                  </div>
                                </div>

                                {expandedReactionId === reaction.id ? (
                                  <div className="space-y-3">
                                    <div className="space-y-2 rounded-xl border border-[var(--surface-border)] bg-[var(--surface)] p-3">
                                      <p className="text-xs font-semibold uppercase tracking-[0.12em] text-[var(--muted)]">
                                        Service
                                      </p>
                                      <div className="grid grid-cols-2 gap-2 sm:grid-cols-3 lg:grid-cols-4">
                                        {reactionServices.map((service) => (
                                          <button
                                            key={`${reaction.id}-${service.id}`}
                                            type="button"
                                            onClick={() => updateReactionService(reaction.id, service)}
                                            className={cn(
                                              "flex h-12 items-center justify-start gap-3 rounded-lg border px-3 text-xs font-semibold transition sm:text-sm",
                                              reaction.service?.id === service.id
                                                ? "border-[var(--blue-primary-3)] bg-[var(--blue-primary-3)]/10 text-[var(--foreground)] shadow-[0_0_0_2px_rgba(28,61,99,0.12)]"
                                                : "border-[var(--surface-border)] bg-[var(--background)] hover:border-[var(--blue-primary-2)]",
                                            )}
                                          >
                                            <span className="flex h-9 w-9 items-center justify-center rounded-full border border-[var(--surface-border)] bg-[var(--surface)] text-[11px] font-semibold uppercase text-[var(--foreground)]">
                                              {service.logoUrl ? (
                                                <img
                                                  src={service.logoUrl}
                                                  alt=""
                                                  className="h-5 w-5 object-contain"
                                                  loading="lazy"
                                                />
                                              ) : (
                                                service.badge
                                              )}
                                            </span>
                                            <span className="truncate">{service.name}</span>
                                          </button>
                                        ))}
                                      </div>
                                    </div>

                                    <div className="space-y-2 rounded-xl border border-[var(--surface-border)] bg-[var(--surface)] p-3">
                                      <p className="text-xs font-semibold uppercase tracking-[0.12em] text-[var(--muted)]">
                                        Action
                                      </p>
                                  {reaction.service && reaction.service.reactions?.length ? (
                                    <div className="space-y-2">
                                      {reaction.service.reactions.map((reactionDef) => {
                                        const isSelected = reaction.reaction?.id === reactionDef.id;
                                        return (
                                              <button
                                                key={`${reaction.id}-${reactionDef.id}`}
                                                type="button"
                                                onClick={() => updateReactionSelection(reaction.id, reactionDef)}
                                                className={cn(
                                                  "flex w-full items-center justify-between rounded-xl border px-4 py-3 text-left text-sm transition",
                                                  isSelected
                                                    ? "border-[var(--blue-primary-3)] bg-[var(--blue-primary-3)]/10 text-[var(--foreground)] shadow-[0_0_0_2px_rgba(28,61,99,0.12)]"
                                                    : "border-[var(--surface-border)] bg-[var(--background)] hover:border-[var(--blue-primary-2)]",
                                                )}
                                              >
                                                <span>{reactionDef.label}</span>
                                                {isSelected ? <span className="text-[var(--blue-primary-3)]">●</span> : null}
                                              </button>
                                            );
                                          })}
                                        </div>
                                      ) : (
                                        <div className="rounded-xl border border-dashed border-[var(--surface-border)] bg-[var(--background)] px-4 py-3 text-sm text-[var(--muted)]">
                                          Choisissez un service pour afficher les actions disponibles.
                                        </div>
                                      )}
                                    </div>

                                    {reaction.reaction ? (
                                      <div className="space-y-2 rounded-xl border border-[var(--surface-border)] bg-[var(--surface)] p-3">
                                        <p className="text-xs font-semibold uppercase tracking-[0.12em] text-[var(--muted)]">
                                          Paramètres de la réaction
                                        </p>
                                        {actionOutputTokens.length ? (
                                          <p className="text-xs text-[var(--muted)]">
                                            Utilisez les boutons + sous les champs pour insérer les données du déclencheur.
                                          </p>
                                        ) : null}
                                        <div className="rounded-xl border border-[var(--surface-border)] bg-[var(--background)] p-4 space-y-3">
                                          {renderFieldInputs(
                                            reaction.reaction.fields,
                                            reaction.fieldValues,
                                            (name, value) => {
                                              clearCreateError();
                                              updateReactionField(reaction.id, name, value);
                                            },
                                            reaction.id,
                                            actionOutputTokens,
                                            (name, token) => insertTokenIntoReactionField(reaction.id, name, token),
                                          )}
                                        </div>
                                      </div>
                                    ) : null}
                                  </div>
                                ) : null}
                              </div>
                            ))}
                          </div>
                        ) : (
                          <div className="rounded-xl border border-dashed border-[var(--surface-border)] bg-[var(--background)] px-4 py-3 text-sm text-[var(--muted)]">
                            Cliquez sur &quot;Ajouter une réaction&quot; pour définir l&apos;action à exécuter après le déclencheur.
                          </div>
                        )
                      ) : null}
                      <div className="flex justify-end">
                        <button
                          type="button"
                          onClick={addReaction}
                          disabled={isLoadingServices || !reactionServices.length}
                          className="inline-flex items-center justify-center gap-2 rounded-full border border-[var(--blue-primary-2)] bg-[var(--blue-primary-2)] px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:border-[var(--blue-primary-3)] hover:bg-[var(--blue-primary-3)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)] disabled:cursor-not-allowed disabled:opacity-50"
                        >
                          Ajouter une réaction
                        </button>
                      </div>
                    </div>
                  </div>
                ) : null}

                {wizardStep === "details" ? (
                  <div className="space-y-4">
                    <div className="space-y-4 rounded-2xl border border-[var(--surface-border)] bg-[var(--surface)] p-5 sm:p-6">
                      <div className="space-y-1">
                        <h3 className="text-lg font-semibold text-[var(--foreground)]">Détails de l&apos;area</h3>
                        <p className="text-sm text-[var(--muted)]">Nom et validation finale.</p>
                      </div>
                      <label className="block space-y-2 text-sm">
                        <span className="text-[var(--muted)]">Nom de l&apos;area</span>
                        <input
                          type="text"
                          value={areaName}
                          onChange={(e) => {
                            clearCreateError();
                            setAreaName(e.target.value);
                          }}
                          placeholder="Démo marketing"
                          className="w-full rounded-xl border border-[var(--surface-border)] bg-[var(--background)] px-4 py-3 text-[var(--foreground)] shadow-sm focus:border-[var(--blue-primary-3)] focus:outline-none focus:ring-2 focus:ring-[var(--blue-primary-3)]/25"
                        />
                      </label>
                      <div className="rounded-xl border border-[var(--surface-border)] bg-[var(--background)] px-4 py-3 text-sm text-[var(--muted)]">
                        Les paramètres du déclencheur et des réactions ont été saisis ci-dessus. Vérifiez-les avant de créer l&apos;area.
                      </div>
                    </div>

                    <div className="rounded-2xl border border-[var(--surface-border)] bg-[var(--surface)] p-5 sm:p-6 shadow-[0_10px_40px_rgba(10,25,45,0.05)] ring-1 ring-[rgba(28,61,99,0.18)]">
                      <div className="space-y-1">
                        <p className="text-xs font-semibold uppercase tracking-[0.14em] text-[var(--blue-primary-3)]">
                          Récapitulatif
                        </p>
                        <h3 className="text-lg font-semibold text-[var(--foreground)]">Vue globale</h3>
                        <p className="text-sm text-[var(--muted)]">Action, réactions et paramètres essentiels regroupés ici.</p>
                      </div>
                      <div className="mt-4 space-y-4 text-sm">
                        <div className="rounded-xl border border-[var(--surface-border)] bg-[var(--background)] p-3">
                          <p className="text-[var(--muted)] text-xs">Déclencheur</p>
                          <p className="text-base font-semibold text-[var(--foreground)]">
                            {actionService?.name ?? "Non défini"}
                          </p>
                          <p className="text-[var(--muted)] text-xs">
                            {selectedAction?.label ?? "Aucun déclencheur sélectionné"}
                          </p>
                          {selectedAction?.fields.map((field) => (
                            <p key={`action-${field.name}`} className="text-xs text-[var(--muted)]">
                              {field.label}:{" "}
                              <span className="text-[var(--foreground)]">
                                {formatFieldDisplayValue(actionFieldValues[field.name]) || "—"}
                              </span>
                            </p>
                          ))}
                        </div>

                        <div className="space-y-2">
                          <p className="text-[var(--muted)] text-xs">Réactions</p>
                          {reactions.length ? (
                            <div className="space-y-2">
                              {reactions.map((reaction, index) => (
                                <div
                                  key={`recap-${reaction.id}`}
                                  className="rounded-xl border border-[var(--surface-border)] bg-[var(--background)] p-3"
                                >
                                  <div className="flex items-start justify-between gap-2">
                                    <div>
                                      <p className="text-xs font-semibold text-[var(--muted)]">#{index + 1}</p>
                                      <p className="text-sm font-semibold text-[var(--foreground)]">
                                        {reaction.service?.name ?? "Service non défini"}
                                      </p>
                                      <p className="text-[var(--muted)] text-xs">
                                        {reaction.reaction?.label ?? "Aucune action sélectionnée"}
                                      </p>
                                    </div>
                                    <span className="text-[var(--muted)] text-xs">
                                      {reactionIsValid(reaction) ? "Prête" : "Incomplète"}
                                    </span>
                                  </div>
                                  {reaction.reaction?.fields.map((field) => (
                                    <p key={`recap-${reaction.id}-${field.name}`} className="text-xs text-[var(--muted)]">
                                      {field.label}:{" "}
                                      <span className="text-[var(--foreground)]">
                                        {formatFieldDisplayValue(reaction.fieldValues[field.name]) || "—"}
                                      </span>
                                    </p>
                                  ))}
                                </div>
                              ))}
                            </div>
                          ) : (
                            <div className="rounded-xl border border-dashed border-[var(--surface-border)] bg-[var(--background)] px-3 py-2 text-[var(--muted)]">
                              Ajoutez au moins une réaction pour finaliser.
                            </div>
                          )}
                        </div>

                        <div className="rounded-xl border border-[var(--surface-border)] bg-[var(--background)] px-3 py-2">
                          <p className="text-[var(--muted)] text-xs">Nom de l&apos;area</p>
                          <p className="text-base font-semibold text-[var(--foreground)]">
                            {areaName.trim() || "Non renseigné"}
                          </p>
                        </div>
                      </div>
                    </div>
                  </div>
                ) : null}

                {createError ? (
                  <div className="rounded-xl border border-[var(--accent)] bg-[var(--surface)] px-4 py-3 text-sm text-[var(--accent)]">
                    {createError}
                  </div>
                ) : null}

                <div className="flex flex-wrap items-center justify-between gap-3">
                  <div className="flex gap-2">
                    {wizardStep !== "action" ? (
                      <button
                        type="button"
                        onClick={goToPreviousStep}
                        className="inline-flex items-center justify-center gap-2 rounded-full border border-[var(--surface-border)] bg-[var(--surface)] px-4 py-2 text-sm font-semibold text-[var(--foreground)] shadow-sm transition hover:border-[var(--blue-primary-2)] hover:text-[var(--blue-primary-2)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)]"
                      >
                        Étape précédente
                      </button>
                    ) : null}
                  </div>
                  <div className="flex gap-3">
                    <button
                      type="button"
                      onClick={() => {
                        closeModal();
                        resetForm();
                      }}
                      className="inline-flex items-center justify-center gap-2 rounded-full border border-[var(--surface-border)] bg-[var(--surface)] px-4 py-2 text-sm font-semibold text-[var(--foreground)] shadow-sm transition hover:border-[var(--blue-primary-2)] hover:text-[var(--blue-primary-2)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)]"
                    >
                      Annuler
                    </button>
                    {wizardStep !== "details" ? (
                      <button
                        type="button"
                        onClick={goToNextStep}
                        disabled={wizardStep === "action" ? !canProceedAction : !canProceedReactions}
                        className="inline-flex items-center justify-center gap-2 rounded-full border border-[var(--blue-primary-2)] bg-[var(--blue-primary-2)] px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:border-[var(--blue-primary-3)] hover:bg-[var(--blue-primary-3)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)] disabled:cursor-not-allowed disabled:opacity-50"
                      >
                        Continuer
                      </button>
                    ) : (
                      <button
                        type="button"
                        onClick={handleCreateArea}
                        disabled={isCreating || !canCreate}
                        className="inline-flex items-center justify-center gap-2 rounded-full border border-[var(--blue-primary-2)] bg-[var(--blue-primary-2)] px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:border-[var(--blue-primary-3)] hover:bg-[var(--blue-primary-3)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)] disabled:cursor-not-allowed disabled:opacity-50"
                      >
                        {isCreating ? "Création..." : "Finaliser et créer"}
                      </button>
                    )}
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

              {areasError ? (
                <div className="rounded-xl border border-[var(--accent)] bg-[var(--surface)] px-4 py-3 text-sm text-[var(--accent)]">
                  {areasError}
                </div>
              ) : null}
              {areaActionError ? (
                <div className="rounded-xl border border-[var(--accent)] bg-[var(--surface)] px-4 py-3 text-sm text-[var(--accent)]">
                  {areaActionError}
                </div>
              ) : null}

              {isLoadingAreas ? (
                <div className="rounded-xl border border-[var(--surface-border)] bg-[var(--surface)] px-4 py-3 text-sm text-[var(--muted)]">
                  Chargement des areas...
                </div>
              ) : null}

              {filteredAreas.length ? (
                <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                  {filteredAreas.map((area) => {
                    const actionVisual = getServiceVisual(
                      services,
                      area.actionServiceId || area.serviceName,
                      area.actionService || area.serviceName,
                    );
                    const primaryReactionServiceId =
                      area.primaryReactionServiceId || area.reactionServiceIds[0] || area.reactionService;
                    const reactionVisual = getServiceVisual(
                      services,
                      primaryReactionServiceId,
                      area.reactionService || area.actionService || area.serviceName,
                    );
                    const actionIcon = actionVisual.logoUrl ? (
                      <img src={actionVisual.logoUrl} alt="" className="h-6 w-6 object-contain" loading="lazy" />
                    ) : (
                      <span>{actionVisual.badge}</span>
                    );
                    const reactionIcon =
                      area.reactionCount > 1 ? (
                        <span className="text-lg font-bold leading-none">+</span>
                      ) : reactionVisual.logoUrl ? (
                        <img src={reactionVisual.logoUrl} alt="" className="h-6 w-6 object-contain" loading="lazy" />
                      ) : (
                        <span>{reactionVisual.badge}</span>
                      );
                    return (
                      <AreaTileCard
                        key={area.id}
                        id={area.id}
                        name={area.summary || area.name}
                        actionLabel={area.actionName}
                        reactionLabel={area.reactionName}
                        actionIcon={actionIcon}
                        reactionIcon={reactionIcon}
                        gradientFrom={area.gradient.from}
                        gradientTo={area.gradient.to}
                        isActive={area.active}
                        onActivate={() => handleToggleArea(area.id, true)}
                        onDeactivate={() => handleToggleArea(area.id, false)}
                        onDelete={() => handleDeleteArea(area.id, area.name)}
                        isBusy={updatingAreaId === area.id}
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
            className="relative w-full max-w-4xl overflow-hidden rounded-3xl border border-[var(--surface-border)] bg-[var(--background)] shadow-2xl"
            onClick={(event) => event.stopPropagation()}
          >
              <div
                className="h-2 w-full"
                style={{
                  background: `linear-gradient(90deg, ${selectedAreaDetail.gradient.from}, ${selectedAreaDetail.gradient.to})`,
                }}
              />
              <div className="space-y-4 px-7 pb-7 pt-6">
                <div className="flex items-start justify-between gap-4">
                  <div className="space-y-1">
                    <p className="text-xs font-semibold uppercase tracking-[0.14em] text-[var(--blue-primary-3)]">
                      Détail de l&apos;automation
                    </p>
                    <h3 className="text-2xl font-semibold text-[var(--foreground)]">
                      {selectedAreaDetail.summary || selectedAreaDetail.name}
                    </h3>
                    <p className="text-sm text-[var(--muted)]">{selectedAreaDetail.name}</p>
                  </div>
                  <div className="flex flex-col items-end gap-2">
                    <span
                      className={cn(
                        "rounded-full px-3 py-1 text-xs font-semibold",
                        selectedAreaDetail.active
                          ? "bg-[var(--success,#22c55e)]/15 text-[var(--success,#22c55e)]"
                          : "bg-[var(--danger,#ef4444)]/15 text-[var(--danger,#ef4444)]",
                      )}
                    >
                      {selectedAreaDetail.active ? "Active" : "Inactive"}
                    </span>
                    <p className="text-[var(--muted)] text-xs">ID: {selectedAreaDetail.id}</p>
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
                </div>

                <div className="space-y-4 rounded-2xl border border-[var(--surface-border)] bg-[var(--surface)] p-4">
                  <p className="text-xs font-semibold uppercase tracking-[0.12em] text-[var(--muted)]">
                    Vue d&apos;ensemble
                  </p>
                  <div className="grid gap-4 lg:grid-cols-3">
                    <div className="rounded-xl border border-[var(--surface-border)] bg-[var(--background)] p-3 text-sm">
                      <p className="text-[var(--muted)] text-xs">Déclencheur</p>
                      <p className="font-semibold text-[var(--foreground)]">
                        {selectedAreaDetail.actionService || "—"}
                      </p>
                      <p className="text-[var(--muted)] text-xs">{selectedAreaDetail.actionName || "—"}</p>
                    </div>
                    <div className="rounded-xl border border-[var(--surface-border)] bg-[var(--background)] p-3 text-sm">
                      <p className="text-[var(--muted)] text-xs">Réaction principale</p>
                      <p className="font-semibold text-[var(--foreground)]">
                        {selectedAreaDetail.reactionService || "—"}
                      </p>
                      <p className="text-[var(--muted)] text-xs">{selectedAreaDetail.reactionName || "—"}</p>
                    </div>
                    <div className="rounded-xl border border-[var(--surface-border)] bg-[var(--background)] p-3 text-sm">
                      <p className="text-[var(--muted)] text-xs">Service associé</p>
                      <p className="font-semibold text-[var(--foreground)]">
                        {selectedAreaDetail.serviceName || selectedAreaDetail.actionService || "—"}
                      </p>
                      <p className="text-[var(--muted)] text-xs">Résumé: {selectedAreaDetail.summary || "—"}</p>
                    </div>
                  </div>
                  <div className="grid gap-3 sm:grid-cols-3">
                    <div className="rounded-xl border border-[var(--surface-border)] bg-[var(--background)] p-3 text-sm">
                      <p className="text-[var(--muted)] text-xs">Début</p>
                      <p className="text-sm font-semibold text-[var(--foreground)]">{selectedAreaDetail.startTime || "—"}</p>
                    </div>
                    <div className="rounded-xl border border-[var(--surface-border)] bg-[var(--background)] p-3 text-sm">
                      <p className="text-[var(--muted)] text-xs">Fin</p>
                      <p className="text-sm font-semibold text-[var(--foreground)]">{selectedAreaDetail.endTime || "—"}</p>
                    </div>
                    <div className="rounded-xl border border-[var(--surface-border)] bg-[var(--background)] p-3 text-sm">
                      <p className="text-[var(--muted)] text-xs">Délai</p>
                      <p className="text-sm font-semibold text-[var(--foreground)]">{selectedAreaDetail.delay}s</p>
                    </div>
                  </div>
                </div>

                <div className="rounded-2xl border border-[var(--surface-border)] bg-[var(--surface)] p-4">
                  <p className="text-xs font-semibold uppercase tracking-[0.12em] text-[var(--muted)]">
                    Récapitulatif
                  </p>
                  <div className="mt-3 grid gap-3 sm:grid-cols-2">
                    <div className="rounded-xl border border-[var(--surface-border)] bg-[var(--background)] p-3 text-sm">
                      <p className="text-[var(--muted)] text-xs">Titre</p>
                      <p className="font-semibold text-[var(--foreground)]">
                        {selectedAreaDetail.summary || selectedAreaDetail.name}
                      </p>
                    </div>
                    <div className="rounded-xl border border-[var(--surface-border)] bg-[var(--background)] p-3 text-sm">
                      <p className="text-[var(--muted)] text-xs">Nom interne</p>
                      <p className="font-semibold text-[var(--foreground)]">{selectedAreaDetail.name}</p>
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
