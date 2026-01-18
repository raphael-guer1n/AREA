"use client";

import { useEffect, useState, type KeyboardEvent, type MouseEvent, type ReactNode } from "react";

import { cn } from "@/lib/helpers";

type ServiceCardProps = {
  name: string;
  url?: string;
  badge: string;
  logoUrl?: string;
  category?: string;
  gradientFrom?: string;
  gradientTo?: string;
  action?: ReactNode;
  actions?: string[];
  reactions?: string[];
  connected?: boolean;
  className?: string;
  onConnect?: () => void;
  onDisconnect?: () => void;
};

type ServiceDetailsModalProps = {
  open: boolean;
  onClose: () => void;
  name: string;
  category?: string;
  url?: string;
  gradientFrom: string;
  gradientTo: string;
  actions: string[];
  reactions: string[];
  connected: boolean;
};

type ServiceConfirmModalProps = {
  open: boolean;
  mode: "connect" | "disconnect";
  onCancel: () => void;
  onConfirm: () => void;
};

export function ServiceCard({
  name,
  url,
  badge,
  logoUrl,
  category,
  gradientFrom = "#002642",
  gradientTo = "#e59500",
  action,
  actions = [],
  reactions = [],
  connected = false,
  className,
  onConnect,
  onDisconnect,
}: ServiceCardProps) {
  const [isDetailsOpen, setIsDetailsOpen] = useState(false);
  const [confirmAction, setConfirmAction] = useState<"connect" | "disconnect" | null>(null);
  const [isConnected, setIsConnected] = useState(connected);

  useEffect(() => {
    setIsConnected(connected);
  }, [connected]);

  const handleCardClick = () => {
    if (isConnected) {
      setIsDetailsOpen(true);
    } else {
      setConfirmAction("connect");
    }
  };

  const handleKeyDown = (event: KeyboardEvent<HTMLDivElement>) => {
    if (event.key === "Enter" || event.key === " ") {
      event.preventDefault();
      handleCardClick();
    }
  };

  const openDetails = (event: MouseEvent) => {
    event.preventDefault();
    event.stopPropagation();
    setIsDetailsOpen(true);
  };

  const requestDisconnect = (event: MouseEvent) => {
    event.preventDefault();
    event.stopPropagation();
    setConfirmAction("disconnect");
  };

  const handleDisconnectConfirm = () => {
    setIsConnected(false);
    setIsDetailsOpen(false);
    setConfirmAction(null);
    onDisconnect?.();
  };

  const handleConnectConfirm = () => {
    setIsConnected(true);
    setConfirmAction(null);
    onConnect?.();
  };

  const chipLabel = !isConnected ? action ?? "À connecter" : "";

  return (
    <>
      <div
        role="button"
        tabIndex={0}
        onClick={handleCardClick}
        onKeyDown={handleKeyDown}
        className="group block h-full focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-[var(--foreground)]"
        aria-label={`Ouvrir ${name}`}
      >
        <article
          className={cn(
            "relative flex h-full flex-col justify-between gap-3 overflow-hidden rounded-xl px-4 py-4 text-white shadow-[0_10px_30px_rgba(0,0,0,0.08)] transition duration-200 hover:-translate-y-[3px] hover:shadow-[0_16px_36px_rgba(0,0,0,0.14)] aspect-[4/3]",
            className,
          )}
          style={{
            background: `linear-gradient(135deg, ${gradientFrom}, ${gradientTo})`,
          }}
        >
          <div className="flex items-start justify-between">
            <Badge logoUrl={logoUrl}>{badge}</Badge>
            <div className="flex items-center gap-2">
              {chipLabel ? (
                <div className="rounded-full border border-white/25 bg-white/18 px-3 py-1 text-[11px] font-semibold uppercase tracking-wide text-white/90 shadow-[0_0_0_1px_rgba(255,255,255,0.15)]">
                  {chipLabel}
                </div>
              ) : null}
              {isConnected ? (
                <button
                  type="button"
                  onClick={requestDisconnect}
                  className="inline-flex h-9 w-9 items-center justify-center rounded-full border border-white/30 bg-white/18 text-white/90 shadow-[0_0_0_1px_rgba(255,255,255,0.12)] backdrop-blur-[2px] transition hover:scale-[1.03] hover:border-white/50 hover:bg-white/25 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
                  aria-label={`Options pour ${name}`}
                >
                  <OptionsIcon />
                </button>
              ) : null}
              {!isConnected ? (
                <button
                  type="button"
                  onClick={openDetails}
                  className="inline-flex h-9 w-9 items-center justify-center rounded-full border border-white/30 bg-white/18 text-white/90 shadow-[0_0_0_1px_rgba(255,255,255,0.12)] backdrop-blur-[2px] transition hover:scale-[1.03] hover:border-white/50 hover:bg-white/25 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
                  aria-label={`Voir les détails de ${name}`}
                >
                  <EyeIcon />
                </button>
              ) : null}
            </div>
          </div>

          <div className="space-y-1">
            <p className="text-lg font-semibold leading-snug">{name}</p>
            {category ? <p className="text-xs text-white/85">{category}</p> : null}
            <ServiceStatus connected={isConnected} tone="dark" />
          </div>
        </article>
      </div>

      <ServiceDetailsModal
        open={isDetailsOpen}
        onClose={() => setIsDetailsOpen(false)}
        name={name}
        category={category}
        url={url}
        gradientFrom={gradientFrom}
        gradientTo={gradientTo}
        actions={actions}
        reactions={reactions}
        connected={isConnected}
      />

      <ServiceConfirmModal
        open={Boolean(confirmAction)}
        mode={confirmAction ?? "connect"}
        onCancel={() => setConfirmAction(null)}
        onConfirm={confirmAction === "disconnect" ? handleDisconnectConfirm : handleConnectConfirm}
      />
    </>
  );
}

function ServiceDetailsModal({
  open,
  onClose,
  name,
  category,
  url,
  gradientFrom,
  gradientTo,
  actions,
  reactions,
  connected,
}: ServiceDetailsModalProps) {
  if (!open) return null;
  const shouldShowLink = Boolean(url && url !== "#");

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center px-4" role="dialog" aria-modal="true">
      <div className="absolute inset-0 bg-[rgba(4,7,15,0.45)] backdrop-blur-[2px]" aria-hidden onClick={onClose} />
      <div className="relative z-10 w-full max-w-4xl overflow-hidden rounded-2xl border border-[var(--surface-border)] bg-[var(--background)] shadow-[0_20px_60px_rgba(0,0,0,0.18)]">
        <div
          className="h-2 w-full"
          style={{
            background: `linear-gradient(135deg, ${gradientFrom}, ${gradientTo})`,
          }}
        />
        <div className="flex items-start justify-between gap-4 px-6 pb-3 pt-4">
          <div>
            <p className="text-xs font-semibold uppercase tracking-[0.14em] text-[var(--blue-primary-3)]">Détails du service</p>
            <h2 className="text-xl font-semibold text-[var(--foreground)]">{name}</h2>
            <div className="flex flex-wrap items-center gap-2">
              {category ? <p className="text-sm text-[var(--muted)]">Catégorie : {category}</p> : null}
              <ServiceStatus connected={connected} tone="light" />
            </div>
            {shouldShowLink ? (
              <a
                href={url}
                target="_blank"
                rel="noreferrer"
                className="mt-3 inline-flex items-center gap-2 text-sm font-semibold text-[var(--blue-primary-2)] underline-offset-4 transition hover:text-[var(--blue-primary-3)] hover:underline"
              >
                Ouvrir le site
                <ExternalLinkIcon />
              </a>
            ) : null}
          </div>
          <button
            type="button"
            onClick={onClose}
            className="inline-flex h-9 w-9 items-center justify-center rounded-full border border-[var(--surface-border)] bg-[var(--surface)] text-[var(--foreground)] transition hover:bg-[var(--surface-border)]/60 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)]"
            aria-label="Fermer"
          >
            <CloseIcon />
          </button>
        </div>

        <div className="grid gap-6 border-t border-[var(--surface-border)] px-6 py-6 md:grid-cols-2">
          <DetailList title="Actions possibles" items={actions} />
          <DetailList title="Réactions disponibles" items={reactions} />
        </div>
      </div>
    </div>
  );
}

function ServiceConfirmModal({ open, mode, onCancel, onConfirm }: ServiceConfirmModalProps) {
  if (!open) return null;

  const isConnect = mode === "connect";

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center px-4">
      <div className="absolute inset-0 bg-[rgba(4,7,15,0.45)] backdrop-blur-[2px]" aria-hidden onClick={onCancel} />
      <div className="relative z-10 w-full max-w-md overflow-hidden rounded-2xl border border-[var(--surface-border)] bg-[var(--background)] shadow-[0_18px_48px_rgba(0,0,0,0.18)]">
        <div className="flex items-start justify-between gap-4 px-6 pb-2 pt-5">
          <div className="space-y-1">
            <p className="text-xs font-semibold uppercase tracking-[0.14em] text-[var(--blue-primary-3)]">
              {isConnect ? "Connexion" : "Déconnexion"}
            </p>
            <h3 className="text-lg font-semibold text-[var(--foreground)]">
              {isConnect ? "Connecter ce service ?" : "Déconnecter ce service ?"}
            </h3>
            <p className="text-sm text-[var(--muted)]">
              {isConnect
                ? "Vous pourrez utiliser ce service dans vos areas après connexion."
                : "Vous ne pourrez plus utiliser ce service dans vos areas après déconnexion."}
            </p>
          </div>
          <button
            type="button"
            onClick={onCancel}
            className="inline-flex h-9 w-9 items-center justify-center rounded-full border border-[var(--surface-border)] bg-[var(--surface)] text-[var(--foreground)] transition hover:bg-[var(--surface-border)]/60 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)]"
            aria-label="Fermer"
          >
            <CloseIcon />
          </button>
        </div>

        <div className="flex items-center justify-end gap-3 border-t border-[var(--surface-border)] px-6 py-4">
          <button
            type="button"
            onClick={onCancel}
            className="inline-flex items-center justify-center gap-2 rounded-full border border-[var(--surface-border)] bg-[var(--surface)] px-4 py-2 text-sm font-semibold text-[var(--foreground)] shadow-sm transition hover:border-[var(--blue-primary-2)] hover:text-[var(--blue-primary-2)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)]"
          >
            Annuler
          </button>
          <button
            type="button"
            onClick={onConfirm}
            className="inline-flex items-center justify-center gap-2 rounded-full border border-[var(--blue-primary-2)] bg-[var(--blue-primary-2)] px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:border-[var(--blue-primary-3)] hover:bg-[var(--blue-primary-3)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)]"
          >
            {isConnect ? "Connecter" : "Déconnecter"}
          </button>
        </div>
      </div>
    </div>
  );
}

function Badge({ children, logoUrl }: { children: ReactNode; logoUrl?: string }) {
  return (
    <span className="flex h-10 w-10 items-center justify-center rounded-full border border-white/35 bg-white/18 text-base font-semibold uppercase text-white shadow-[0_8px_18px_rgba(0,0,0,0.12)] backdrop-blur-[2px]">
      {logoUrl ? (
        <img src={logoUrl} alt="" className="h-6 w-6 object-contain" loading="lazy" />
      ) : (
        children
      )}
    </span>
  );
}

function OptionsIcon() {
  return (
    <svg viewBox="0 0 24 24" className="h-5 w-5" fill="none" stroke="currentColor" strokeWidth="1.8">
      <circle cx="12" cy="5" r="1.6" />
      <circle cx="12" cy="12" r="1.6" />
      <circle cx="12" cy="19" r="1.6" />
    </svg>
  );
}

function EyeIcon() {
  return (
    <svg viewBox="0 0 24 24" className="h-5 w-5" fill="none" stroke="currentColor" strokeWidth="1.6">
      <path
        d="M3 12s3.5-6 9-6 9 6 9 6-3.5 6-9 6-9-6-9-6Z"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <circle cx="12" cy="12" r="2.5" />
    </svg>
  );
}

function CloseIcon() {
  return (
    <svg viewBox="0 0 24 24" className="h-4 w-4" fill="none" stroke="currentColor" strokeWidth="1.8">
      <path d="M6 6 18 18M6 18 18 6" strokeLinecap="round" />
    </svg>
  );
}

function ExternalLinkIcon() {
  return (
    <svg viewBox="0 0 20 20" className="h-4 w-4" fill="none" stroke="currentColor" strokeWidth="1.6">
      <path d="M11.5 3.5h5v5" strokeLinecap="round" />
      <path d="m9 11 7.5-7.5M16 11v5H4V4h5" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  );
}

function DetailList({ title, items }: { title: string; items: string[] }) {
  return (
    <div className="space-y-3 rounded-2xl border border-[var(--surface-border)] bg-[var(--surface)] px-4 py-4 shadow-[0_10px_30px_rgba(0,0,0,0.04)]">
      <p className="text-sm font-semibold text-[var(--foreground)]">{title}</p>
      {items.length > 0 ? (
        <ul className="space-y-2 text-sm text-[var(--muted)]">
          {items.map((item) => (
            <li
              key={item}
              className="flex items-start gap-2 rounded-xl border border-[var(--surface-border)] bg-[var(--background)] px-3 py-2 text-[var(--foreground)]"
            >
              <span className="mt-0.5 inline-block h-2 w-2 rounded-full bg-[var(--blue-primary-2)]" aria-hidden />
              <span className="leading-snug">{item}</span>
            </li>
          ))}
        </ul>
      ) : (
        <p className="text-sm text-[var(--placeholder)]">Aucun élément disponible pour le moment.</p>
      )}
    </div>
  );
}

function ServiceStatus({ connected, tone = "dark" }: { connected: boolean; tone?: "dark" | "light" }) {
  const color = connected ? "bg-[var(--success,#22c55e)]" : "bg-[var(--danger,#ef4444)]";
  const label = connected ? "Connecté" : "Non connecté";
  const base =
    tone === "dark"
      ? "border-white/25 bg-white/10 text-white/90 shadow-[0_0_0_1px_rgba(255,255,255,0.12)]"
      : "border-[var(--surface-border)] bg-[var(--surface)] text-[var(--foreground)] shadow-[0_1px_2px_rgba(0,0,0,0.08)]";
  const dotShadow =
    tone === "dark" ? "shadow-[0_0_0_2px_rgba(255,255,255,0.2)]" : "shadow-[0_0_0_1px_rgba(0,0,0,0.08)]";

  return (
    <span className={cn("inline-flex items-center gap-2 rounded-full px-3 py-1 text-xs font-semibold uppercase tracking-wide", base)}>
      <span className={cn("inline-block h-2.5 w-2.5 rounded-full", dotShadow, color)} />
      {label}
    </span>
  );
}
