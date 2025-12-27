"use client";

import { useEffect, useMemo, useState } from "react";

import { Card } from "@/components/ui/AreaCard";
import { Button } from "@/components/ui/Button";
import { useAuth } from "@/hooks/useAuth";
import type { AuthStatus, AuthSession } from "@/types/auth";
import type { User } from "@/types/User";

type ProfileClientProps = {
  initialUser: User | null;
  initialSession: AuthSession | null;
};

type StatusStyle = {
  label: string;
  toneClass: string;
  dotClass: string;
  helper: string;
};

type ProfileNotification = {
  id: string;
  title: string;
  detail: string;
  createdAt: string;
  type: "area_created" | "info" | "warning";
};

const statusStyles: Record<AuthStatus, StatusStyle> = {
  authenticated: {
    label: "Authentifié",
    toneClass: "text-[var(--blue-primary-2)]",
    dotClass: "bg-[var(--blue-primary-2)]",
    helper: "Session active. Vous pouvez créer des automatisations.",
  },
  unauthenticated: {
    label: "Non authentifié",
    toneClass: "text-[var(--card-color-2)]",
    dotClass: "bg-[var(--card-color-2)]",
    helper: "Connectez-vous pour accéder à vos services.",
  },
  loading: {
    label: "Chargement...",
    toneClass: "text-[var(--muted)]",
    dotClass: "bg-[var(--muted)]",
    helper: "Mise à jour de la session...",
  },
  idle: {
    label: "En attente",
    toneClass: "text-[var(--muted)]",
    dotClass: "bg-[var(--muted)]",
    helper: "Session en veille, aucune action en cours.",
  },
  error: {
    label: "Erreur",
    toneClass: "text-[var(--card-color-3)]",
    dotClass: "bg-[var(--card-color-3)]",
    helper: "Un incident est survenu lors de la récupération de la session.",
  },
};

function Toggle({ label, checked, onChange }: { label: string; checked: boolean; onChange: () => void }) {
  return (
    <button
      type="button"
      onClick={onChange}
      className="flex w-full items-center justify-between gap-3 rounded-xl border border-[var(--surface-border)] bg-[var(--surface)] px-4 py-3 text-sm font-semibold text-[var(--foreground)] transition hover:border-[var(--foreground)]/60 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--foreground)]"
      aria-pressed={checked}
    >
      <span>{label}</span>
      <span
        className={`inline-flex h-5 w-9 items-center rounded-full transition ${checked ? "bg-[var(--blue-primary-2)]" : "bg-[var(--surface-border)]"}`}
      >
        <span
          className={`ml-[2px] inline-block h-4 w-4 rounded-full bg-white shadow transition ${checked ? "translate-x-[14px]" : ""}`}
        />
      </span>
    </button>
  );
}

export function ProfileClient({ initialUser, initialSession }: ProfileClientProps) {
  const { user, token, status } = useAuth({
    initialUser,
    initialSession,
  });

  const [notifyApp, setNotifyApp] = useState(true);
  const [notifyEmail, setNotifyEmail] = useState(true);
  const [notifications, setNotifications] = useState<ProfileNotification[]>([]);

  const displayName = useMemo(
    () => user?.name || user?.username || user?.email || "Utilisateur",
    [user],
  );
  const initials = useMemo(() => displayName.slice(0, 2).toUpperCase(), [displayName]);
  const displayEmail = user?.email || "Email indisponible";
  const statusMeta = statusStyles[status] ?? statusStyles.idle;
  const maskedToken =
    token && token.length > 18
      ? `${token.slice(0, 10)}…${token.slice(-6)}`
      : token || "Aucun token actif";

  useEffect(() => {
    if (typeof window === "undefined") return;
    const raw = window.localStorage.getItem("area-notifications");
    if (raw) {
      try {
        const parsed = JSON.parse(raw) as ProfileNotification[];
        setNotifications(parsed);
        return;
      } catch {
        // fallback handled below
      }
    }
    const seed: ProfileNotification[] = [
      {
        id: "seed-1",
        title: "Bienvenue sur AREA",
        detail: "Connectez un service pour commencer à créer des automatisations.",
        createdAt: new Date().toISOString(),
        type: "info",
      },
    ];
    setNotifications(seed);
    window.localStorage.setItem("area-notifications", JSON.stringify(seed));
  }, []);

  useEffect(() => {
    if (typeof window === "undefined") return;
    window.localStorage.setItem("area-notifications", JSON.stringify(notifications));
  }, [notifications]);

  const addNotification = (title: string, detail: string, type: ProfileNotification["type"]) => {
    const next: ProfileNotification = {
      id: `notif-${Date.now()}`,
      title,
      detail,
      createdAt: new Date().toISOString(),
      type,
    };
    setNotifications((prev) => [next, ...prev].slice(0, 25));
  };

  const simulateAreaCreation = () => {
    addNotification("AREA créée", "Votre nouvelle automation est prête.", "area_created");
  };

  const clearNotifications = () => {
    setNotifications([]);
  };

  return (
    <div className="space-y-5">
      <section className="relative isolate overflow-hidden rounded-[22px] border border-[var(--surface-border)] bg-white px-6 py-6 shadow-sm ring-1 ring-[rgba(28,61,99,0.14)]">
        <div className="pointer-events-none absolute inset-0 -z-10 opacity-70">
          <div
            className="absolute -right-10 -top-10 h-40 w-40 rounded-full bg-[radial-gradient(circle_at_center,var(--card-color-3)_0,transparent_70%)] blur-3xl"
            aria-hidden
          />
          <div
            className="absolute left-4 top-14 h-32 w-32 rounded-full bg-[radial-gradient(circle_at_center,var(--card-color-4)_0,transparent_70%)] blur-3xl"
            aria-hidden
          />
        </div>

        <div className="rounded-[18px] border border-[var(--surface-border)] bg-white px-6 py-6 ring-1 ring-[rgba(28,61,99,0.1)]">
          <div className="flex flex-col items-center gap-4 text-center">
            <div className="flex h-16 w-16 items-center justify-center rounded-2xl bg-[var(--card-color-1)] text-lg font-semibold text-white shadow-sm ring-1 ring-[var(--surface-border)]">
              {initials}
            </div>
            <div className="space-y-1">
              <p className="text-xs font-semibold uppercase tracking-[0.12em] text-[var(--blue-primary-3)]">
                Profil
              </p>
              <p className="text-xl font-semibold text-[var(--foreground)]">{displayName}</p>
              <p className="text-sm text-[var(--muted)]">{displayEmail}</p>
            </div>
            <div className="grid w-full max-w-2xl gap-3 sm:grid-cols-2">
              <div className="rounded-2xl border border-[var(--surface-border)] bg-white px-4 py-3 text-left">
                <p className="text-xs font-semibold uppercase tracking-[0.1em] text-[var(--muted)]">Username</p>
                <p className="mt-2 text-base font-semibold text-[var(--foreground)]">
                  {user?.username ?? "N/A"}
                </p>
              </div>
              <div className="rounded-2xl border border-[var(--surface-border)] bg-white px-4 py-3 text-left">
                <p className="text-xs font-semibold uppercase tracking-[0.1em] text-[var(--muted)]">Email</p>
                <p className="mt-2 break-words text-base font-semibold text-[var(--foreground)]">
                  {displayEmail}
                </p>
              </div>
            </div>
          </div>
        </div>
      </section>

      <Card
        title="Notifications"
        subtitle="Historique local (JSON stocké en localStorage)"
        className="rounded-[18px] border-[var(--surface-border)] bg-white ring-1 ring-[rgba(28,61,99,0.1)]"
      >
        <div className="flex flex-col gap-4">
          <div className="flex flex-wrap items-center justify-between gap-3">
            <div className="flex items-center gap-3">
              <span className="rounded-full bg-white px-3 py-1 text-xs font-semibold text-[var(--muted)] ring-1 ring-[var(--surface-border)]">
                {notifications.length} notification{notifications.length > 1 ? "s" : ""}
              </span>
              <span className="text-xs text-[var(--muted)]">Key: area-notifications</span>
            </div>
            <div className="flex items-center gap-2">
              <Button
                type="button"
                variant="secondary"
                className="px-4 py-2"
                onClick={simulateAreaCreation}
              >
                Simuler une notif
              </Button>
              <Button
                type="button"
                variant="ghost"
                className="px-3 py-2"
                onClick={clearNotifications}
                disabled={notifications.length === 0}
              >
                Vider
              </Button>
            </div>
          </div>

          <div className="space-y-3">
            {notifications.length === 0 ? (
              <p className="text-sm text-[var(--muted)]">Aucune notification pour le moment.</p>
            ) : null}
            {notifications.map((notif) => (
              <div key={notif.id} className="rounded-2xl border border-[var(--surface-border)] bg-white px-4 py-3 shadow-sm ring-1 ring-[rgba(28,61,99,0.06)]">
                <div className="flex items-start justify-between gap-3">
                  <div className="flex items-center gap-3">
                    <span
                      className={`inline-flex h-7 w-7 items-center justify-center rounded-full text-xs font-semibold ${
                        notif.type === "area_created"
                          ? "bg-[var(--blue-primary-2)] text-white"
                          : notif.type === "warning"
                            ? "bg-[var(--card-color-3)] text-white"
                            : "bg-[var(--surface-border)] text-[var(--foreground)]"
                      }`}
                      aria-hidden
                    >
                      {notif.type === "area_created" ? "A" : notif.type === "warning" ? "!" : "i"}
                    </span>
                    <div>
                      <p className="text-sm font-semibold text-[var(--foreground)]">{notif.title}</p>
                      <p className="text-xs text-[var(--muted)]">{notif.detail}</p>
                    </div>
                  </div>
                  <span className="text-[11px] font-medium uppercase tracking-wide text-[var(--muted)]">
                    {new Date(notif.createdAt).toLocaleString()}
                  </span>
                </div>
              </div>
            ))}
          </div>
        </div>
      </Card>

      <div className="grid gap-4 lg:grid-cols-2">
        <Card
          title="Session & sécurité"
          subtitle={statusMeta.helper}
          className="rounded-[18px] border-[var(--surface-border)] bg-white ring-1 ring-[rgba(28,61,99,0.1)]"
        >
          <div className="space-y-3 text-sm">
            <div className="flex items-center justify-between">
              <span className="text-[var(--muted)]">Statut</span>
              <span className={`inline-flex items-center gap-2 rounded-full px-3 py-1 text-xs font-semibold ${statusMeta.toneClass} bg-white ring-1 ring-[var(--surface-border)]`}>
                <span className={`h-2.5 w-2.5 rounded-full ${statusMeta.dotClass}`} aria-hidden />
                {statusMeta.label}
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-[var(--muted)]">Token masqué</span>
              <span className="font-mono text-xs text-[var(--foreground)]">{maskedToken}</span>
            </div>
          </div>
        </Card>

        <Card
          title="Préférences"
          subtitle="Réglez vos notifications."
          className="rounded-[18px] border-[var(--surface-border)] bg-white ring-1 ring-[rgba(28,61,99,0.1)]"
        >
          <div className="space-y-2">
            <Toggle label="Notifications in-app" checked={notifyApp} onChange={() => setNotifyApp((v) => !v)} />
            <Toggle label="Emails importants" checked={notifyEmail} onChange={() => setNotifyEmail((v) => !v)} />
          </div>
        </Card>
      </div>

      <Card
        title="Support"
        subtitle="Aide et diagnostic rapide"
        className="rounded-[18px] border-[var(--surface-border)] bg-white ring-1 ring-[rgba(28,61,99,0.1)]"
      >
        <div className="space-y-3 text-sm text-[var(--muted)]">
          <p>Consultez les logs dans Area pour repérer les erreurs d&apos;exécution et les payloads échoués.</p>
          <p>Rafraîchissez votre session si vous constatez un statut d&apos;erreur ou une expiration de token.</p>
          <p>Reconnectez les services défaillants depuis l&apos;onglet Services, puis relancez l&apos;automation.</p>
          <p>Si le problème persiste, capturez l&apos;horodatage, le service concerné et le message d&apos;erreur avant d&apos;ouvrir un ticket.</p>
        </div>
      </Card>
    </div>
  );
}
