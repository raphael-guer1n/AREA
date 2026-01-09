"use client";

import { useEffect, useMemo, useState, type ChangeEvent } from "react";
import { useRouter } from "next/navigation";

import { Card } from "@/components/ui/AreaCard";
import { Button } from "@/components/ui/Button";
import { useAuth } from "@/hooks/useAuth";
import type { AuthStatus, AuthSession } from "@/types/auth";
import type { User } from "@/types/User";

type ConfirmAction = "logout" | "switch" | "delete";
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

type ProfileDraft = {
  avatarUrl: string;
  email: string;
  username: string;
  name: string;
  password: string;
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
          className={`ml-[2px] inline-block h-4 w-4 rounded-full bg-[var(--background)] ring-1 ring-[var(--surface-border)] shadow transition ${checked ? "translate-x-[14px]" : ""}`}
        />
      </span>
    </button>
  );
}

export function ProfileClient({ initialUser, initialSession }: ProfileClientProps) {
  const { user, token, status, logout } = useAuth({
    initialUser,
    initialSession,
  });
  const router = useRouter();

  const [localProfile, setLocalProfile] = useState<User | null>(user);
  const [profileDraft, setProfileDraft] = useState<ProfileDraft>({
    avatarUrl: user?.avatarUrl ?? "",
    email: user?.email ?? "",
    username: user?.username ?? "",
    name: user?.name ?? "",
    password: "",
  });
  const [avatarPreview, setAvatarPreview] = useState<string | null>(user?.avatarUrl ?? null);
  const [isEditOpen, setIsEditOpen] = useState(false);
  const [confirmAction, setConfirmAction] = useState<ConfirmAction | null>(null);
  const [notifyApp, setNotifyApp] = useState(true);
  const [notifyEmail, setNotifyEmail] = useState(true);
  const [notifications, setNotifications] = useState<ProfileNotification[]>([]);

  useEffect(() => {
    setLocalProfile(user);
    setProfileDraft({
      avatarUrl: user?.avatarUrl ?? "",
      email: user?.email ?? "",
      username: user?.username ?? "",
      name: user?.name ?? "",
      password: "",
    });
    setAvatarPreview(user?.avatarUrl ?? null);
  }, [user]);

  useEffect(() => {
    return () => {
      if (avatarPreview && avatarPreview.startsWith("blob:")) {
        URL.revokeObjectURL(avatarPreview);
      }
    };
  }, [avatarPreview]);

  const profile = localProfile ?? user ?? initialUser;
  const isGoogleAccount = useMemo(() => {
    const id = profile?.id?.toString().toLowerCase() ?? "";
    const email = profile?.email?.toLowerCase() ?? "";
    return id.startsWith("google-") || email.includes("googleusercontent") || email.includes("google-oauth");
  }, [profile?.email, profile?.id]);

  const displayName = useMemo(
    () => profile?.name || profile?.username || profile?.email || "Utilisateur",
    [profile?.email, profile?.name, profile?.username],
  );
  const initials = useMemo(() => displayName.slice(0, 2).toUpperCase(), [displayName]);
  const displayEmail = profile?.email || "Email indisponible";
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

  const handleLogout = async () => {
    await logout();
    router.push("/");
  };

  const handleSwitchAccount = async () => {
    await logout();
    router.push("/login");
  };

  const handleDeleteAccount = () => {
    addNotification("Suppression demandée", "Aucune action backend configurée.", "warning");
  };

  const handleAvatarChange = (event: ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    const nextPreview = URL.createObjectURL(file);
    if (avatarPreview && avatarPreview.startsWith("blob:")) {
      URL.revokeObjectURL(avatarPreview);
    }
    setAvatarPreview(nextPreview);
    setProfileDraft((prev) => ({ ...prev, avatarUrl: nextPreview }));
  };

  const resetAvatar = () => {
    if (avatarPreview && avatarPreview.startsWith("blob:")) {
      URL.revokeObjectURL(avatarPreview);
    }
    const fallback = profile?.avatarUrl ?? "";
    setAvatarPreview(fallback || null);
    setProfileDraft((prev) => ({ ...prev, avatarUrl: fallback }));
  };

  const handleSaveProfile = () => {
    const base: User = localProfile ?? user ?? {
      id: initialUser?.id ?? "local-user",
      email: profileDraft.email || "email@exemple.com",
    };

    const nextAvatar = avatarPreview ?? profileDraft.avatarUrl ?? base.avatarUrl ?? "";

    const updatedProfile: User = {
      ...base,
      email: profileDraft.email || base.email,
      username: profileDraft.username || base.username,
      name: profileDraft.name || base.name || base.username || base.email,
      avatarUrl: nextAvatar || undefined,
    };

    setLocalProfile(updatedProfile);
    setProfileDraft({
      avatarUrl: updatedProfile.avatarUrl ?? "",
      email: updatedProfile.email ?? "",
      username: updatedProfile.username ?? "",
      name: updatedProfile.name ?? "",
      password: "",
    });
    setAvatarPreview(nextAvatar || null);

    addNotification("Profil mis à jour", "Modifications enregistrées localement (aucun appel backend).", "info");
    setIsEditOpen(false);
  };

  const confirmCopy: Record<ConfirmAction, { title: string; description: string; confirmLabel: string }> = {
    logout: {
      title: "Déconnexion ?",
      description: "Vous serez redirigé vers l'accueil et la session sera vidée.",
      confirmLabel: "Se déconnecter",
    },
    switch: {
      title: "Changer de compte ?",
      description: "Vous serez déconnecté puis redirigé vers la page de connexion.",
      confirmLabel: "Changer de compte",
    },
    delete: {
      title: "Supprimer le compte ?",
      description: "Aucune suppression réelle n'est déclenchée. Une notification locale sera ajoutée.",
      confirmLabel: "Supprimer",
    },
  };

  const executeConfirmedAction = async (action: ConfirmAction) => {
    setConfirmAction(null);
    if (action === "logout") {
      await handleLogout();
      return;
    }
    if (action === "switch") {
      await handleSwitchAccount();
      return;
    }
    if (action === "delete") {
      handleDeleteAccount();
    }
  };

  const baseActionButton =
    "px-8 py-4 text-lg border border-[var(--surface-border)] bg-[var(--background)] text-[var(--foreground)] hover:border-[var(--blue-primary-2)]";

  return (
    <div className="space-y-7">
      {isEditOpen ? (
        <div className="fixed inset-0 z-50 flex items-center justify-center px-4 py-6">
          <div
            className="absolute inset-0 bg-[rgba(4,7,15,0.45)] backdrop-blur-[2px]"
            aria-hidden
            onClick={() => setIsEditOpen(false)}
          />
          <div
            role="dialog"
            aria-modal="true"
            className="relative z-10 w-full max-w-5xl overflow-hidden rounded-[26px] border border-[var(--surface-border)] bg-[var(--background)] shadow-[0_22px_68px_rgba(0,0,0,0.2)] ring-1 ring-[rgba(28,61,99,0.18)]"
          >
            <div
              className="h-2 w-full"
              style={{
                background: "linear-gradient(135deg, var(--blue-primary-2), var(--card-color-3))",
              }}
            />
            <div className="flex items-start justify-between gap-6 px-8 pb-4 pt-6">
              <div>
                <p className="text-xs font-semibold uppercase tracking-[0.14em] text-[var(--blue-primary-3)]">
                  Edition du profil
                </p>
                <p className="text-xl font-semibold text-[var(--foreground)]">Modifier le compte</p>
                <p className="text-sm text-[var(--muted)]">
                  Modifications locales uniquement. Email et mot de passe restent verrouillés pour les comptes Google.
                </p>
              </div>
              <button
                type="button"
                onClick={() => setIsEditOpen(false)}
                className="inline-flex h-9 w-9 items-center justify-center rounded-full border border-[var(--surface-border)] bg-[var(--surface)] text-[var(--foreground)] transition hover:bg-[var(--surface-border)]/60 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)]"
                aria-label="Fermer la fenêtre de modification"
              >
                ✕
              </button>
            </div>

            <div className="grid gap-8 border-t border-[var(--surface-border)] px-8 py-8 lg:grid-cols-[1fr_1fr]">
              <div className="space-y-4">
                <div className="flex flex-wrap items-center gap-4">
                  <div className="relative h-24 w-24 overflow-hidden rounded-2xl border border-[var(--surface-border)] bg-[var(--surface)] text-lg font-semibold text-[var(--foreground)] ring-1 ring-[var(--surface-border)]">
                    {avatarPreview || profileDraft.avatarUrl || profile?.avatarUrl ? (
                      <img
                        src={avatarPreview ?? profileDraft.avatarUrl ?? profile?.avatarUrl ?? ""}
                        alt={`Avatar de ${displayName}`}
                        className="h-full w-full object-cover"
                      />
                    ) : (
                      <span className="flex h-full w-full items-center justify-center text-xl">{initials}</span>
                    )}
                  </div>
                  <div className="flex flex-wrap gap-2">
                    <label className="inline-flex cursor-pointer items-center gap-2 rounded-full border border-[var(--surface-border)] bg-[var(--surface)] px-4 py-2 text-sm font-semibold text-[var(--foreground)] transition hover:border-[var(--blue-primary-2)] focus-within:outline focus-within:outline-2 focus-within:outline-offset-2 focus-within:outline-[var(--blue-primary-2)]">
                      <input
                        type="file"
                        accept="image/*"
                        className="hidden"
                        onChange={handleAvatarChange}
                      />
                      Changer la photo
                    </label>
                    <button
                      type="button"
                      onClick={resetAvatar}
                      disabled={!avatarPreview && !profileDraft.avatarUrl && !profile?.avatarUrl}
                      className="inline-flex items-center justify-center rounded-full border border-[var(--surface-border)] bg-[var(--surface)] px-4 py-2 text-sm font-semibold text-[var(--foreground)] transition hover:border-[var(--blue-primary-2)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)] disabled:cursor-not-allowed disabled:opacity-60"
                    >
                      Réinitialiser
                    </button>
                  </div>
                </div>

                <div className="space-y-1">
                  <label className="text-xs font-semibold uppercase tracking-[0.1em] text-[var(--muted)]">
                    Nom / affichage
                  </label>
                  <input
                    type="text"
                    value={profileDraft.name}
                    onChange={(event) =>
                      setProfileDraft((prev) => ({ ...prev, name: event.target.value }))
                    }
                    placeholder={displayName}
                    className="w-full rounded-xl border border-[var(--surface-border)] bg-[var(--surface)] px-4 py-3 text-sm font-semibold text-[var(--foreground)] outline-none ring-1 ring-transparent transition focus:border-[var(--blue-primary-2)] focus:ring-[var(--blue-primary-2)]"
                  />
                </div>
                <div className="space-y-1">
                  <label className="text-xs font-semibold uppercase tracking-[0.1em] text-[var(--muted)]">
                    Username
                  </label>
                  <input
                    type="text"
                    value={profileDraft.username}
                    onChange={(event) =>
                      setProfileDraft((prev) => ({ ...prev, username: event.target.value }))
                    }
                    placeholder={profile?.username ?? "username"}
                    className="w-full rounded-xl border border-[var(--surface-border)] bg-[var(--surface)] px-4 py-3 text-sm font-semibold text-[var(--foreground)] outline-none ring-1 ring-transparent transition focus:border-[var(--blue-primary-2)] focus:ring-[var(--blue-primary-2)]"
                  />
                </div>
              </div>

              <div className="space-y-4">
                <div className="space-y-1">
                  <label className="text-xs font-semibold uppercase tracking-[0.1em] text-[var(--muted)]">
                    Adresse mail
                  </label>
                  <input
                    type="email"
                    value={profileDraft.email}
                    onChange={(event) =>
                      setProfileDraft((prev) => ({ ...prev, email: event.target.value }))
                    }
                    placeholder={profile?.email ?? "email@exemple.com"}
                    disabled={isGoogleAccount}
                    className="w-full rounded-xl border border-[var(--surface-border)] bg-[var(--surface)] px-4 py-3 text-sm font-semibold text-[var(--foreground)] outline-none ring-1 ring-transparent transition focus:border-[var(--blue-primary-2)] focus:ring-[var(--blue-primary-2)] disabled:cursor-not-allowed disabled:bg-[var(--surface-border)]/60"
                  />
                  <p className="text-[11px] font-medium text-[var(--muted)]">
                    {isGoogleAccount
                      ? "Email verrouillé car compte créé via Google (aucune requête envoyée)."
                      : "Changement local uniquement : la donnée n'est pas envoyée au backend."}
                  </p>
                </div>

                <div className="space-y-1">
                  <label className="text-xs font-semibold uppercase tracking-[0.1em] text-[var(--muted)]">
                    Mot de passe
                  </label>
                  <input
                    type="password"
                    value={profileDraft.password}
                    onChange={(event) =>
                      setProfileDraft((prev) => ({ ...prev, password: event.target.value }))
                    }
                    placeholder="********"
                    disabled={isGoogleAccount}
                    className="w-full rounded-xl border border-[var(--surface-border)] bg-[var(--surface)] px-4 py-3 text-sm font-semibold text-[var(--foreground)] outline-none ring-1 ring-transparent transition focus:border-[var(--blue-primary-2)] focus:ring-[var(--blue-primary-2)] disabled:cursor-not-allowed disabled:bg-[var(--surface-border)]/60"
                  />
                  <p className="text-[11px] font-medium text-[var(--muted)]">
                    Saisie mémorisée localement pour référence, aucune mise à jour réelle du compte.
                  </p>
                </div>

                <div className="flex flex-wrap items-center justify-end gap-3 pt-2">
                  <Button
                    type="button"
                    variant="ghost"
                    className="px-4 py-2"
                    onClick={() => {
                      setIsEditOpen(false);
                      setProfileDraft((prev) => ({ ...prev, password: "" }));
                    }}
                  >
                    Annuler
                  </Button>
                  <Button type="button" className="px-5 py-2" onClick={handleSaveProfile}>
                    Enregistrer localement
                  </Button>
                </div>
              </div>
            </div>
          </div>
        </div>
      ) : null}

      <section className="relative isolate overflow-hidden rounded-[22px] border border-[var(--surface-border)] bg-[var(--background)] px-6 py-6 shadow-sm ring-1 ring-[rgba(28,61,99,0.14)]">
        <div className="rounded-[18px] border border-[var(--surface-border)] bg-[var(--background)] px-6 py-6 ring-1 ring-[rgba(28,61,99,0.1)]">
          <div className="flex flex-col items-center gap-4 text-center">
            <div className="flex h-16 w-16 items-center justify-center overflow-hidden rounded-2xl bg-[var(--card-color-1)] text-lg font-semibold text-white shadow-sm ring-1 ring-[var(--surface-border)]">
              {profile?.avatarUrl ? (
                <img
                  src={profile.avatarUrl}
                  alt={`Avatar de ${displayName}`}
                  className="h-full w-full object-cover"
                />
              ) : (
                initials
              )}
            </div>
            <div className="space-y-1">
              <p className="text-xs font-semibold uppercase tracking-[0.12em] text-[var(--blue-primary-3)]">
                Profil
              </p>
              <p className="text-xl font-semibold text-[var(--foreground)]">{displayName}</p>
              <p className="text-sm text-[var(--muted)]">{displayEmail}</p>
            </div>
            <div className="grid w-full max-w-2xl gap-3 sm:grid-cols-2">
              <div className="rounded-2xl border border-[var(--surface-border)] bg-[var(--background)] px-4 py-3 text-left">
                <p className="text-xs font-semibold uppercase tracking-[0.1em] text-[var(--muted)]">Username</p>
                <p className="mt-2 text-base font-semibold text-[var(--foreground)]">
                  {profile?.username ?? "N/A"}
                </p>
              </div>
              <div className="rounded-2xl border border-[var(--surface-border)] bg-[var(--background)] px-4 py-3 text-left">
                <p className="text-xs font-semibold uppercase tracking-[0.1em] text-[var(--muted)]">Email</p>
                <p className="mt-2 break-words text-base font-semibold text-[var(--foreground)]">
                  {displayEmail}
                </p>
              </div>
            </div>
          </div>
        </div>
      </section>

      <div className="flex flex-wrap items-center justify-center gap-6 py-3">
        <Button
          type="button"
          variant="ghost"
          className={baseActionButton}
          onClick={() => setIsEditOpen(true)}
        >
          Modifier le compte
        </Button>
        <Button
          type="button"
          variant="ghost"
          className={baseActionButton}
          onClick={() => setConfirmAction("logout")}
        >
          Déconnexion
        </Button>
        <Button
          type="button"
          variant="ghost"
          className={baseActionButton}
          onClick={() => setConfirmAction("switch")}
        >
          Changer de compte
        </Button>
        <Button
          type="button"
          variant="ghost"
          className={baseActionButton}
          onClick={() => setConfirmAction("delete")}
        >
          Supprimer le compte
        </Button>
      </div>

      <Card
        title="Notifications"
        subtitle="Historique local (JSON stocké en localStorage)"
        tone="background"
        className="rounded-[18px] border-[var(--surface-border)] ring-1 ring-[rgba(28,61,99,0.1)]"
      >
        <div className="flex flex-col gap-4">
          <div className="flex flex-wrap items-center justify-between gap-3">
            <div className="flex items-center gap-3">
              <span className="rounded-full bg-[var(--background)] px-3 py-1 text-xs font-semibold text-[var(--muted)] ring-1 ring-[var(--surface-border)]">
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
              <div key={notif.id} className="rounded-2xl border border-[var(--surface-border)] bg-[var(--background)] px-4 py-3 shadow-sm ring-1 ring-[rgba(28,61,99,0.06)]">
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
          tone="background"
          className="rounded-[18px] border-[var(--surface-border)] ring-1 ring-[rgba(28,61,99,0.1)]"
        >
          <div className="space-y-3 text-sm">
            <div className="flex items-center justify-between">
              <span className="text-[var(--muted)]">Statut</span>
              <span className={`inline-flex items-center gap-2 rounded-full px-3 py-1 text-xs font-semibold ${statusMeta.toneClass} bg-[var(--background)] ring-1 ring-[var(--surface-border)]`}>
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
          tone="background"
          className="rounded-[18px] border-[var(--surface-border)] ring-1 ring-[rgba(28,61,99,0.1)]"
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
        tone="background"
        className="rounded-[18px] border-[var(--surface-border)] ring-1 ring-[rgba(28,61,99,0.1)]"
      >
        <div className="space-y-3 text-sm text-[var(--muted)]">
          <p>Consultez les logs dans Area pour repérer les erreurs d&apos;exécution et les payloads échoués.</p>
          <p>Rafraîchissez votre session si vous constatez un statut d&apos;erreur ou une expiration de token.</p>
          <p>Reconnectez les services défaillants depuis l&apos;onglet Services, puis relancez l&apos;automation.</p>
          <p>Si le problème persiste, capturez l&apos;horodatage, le service concerné et le message d&apos;erreur avant d&apos;ouvrir un ticket.</p>
        </div>
      </Card>

      <ProfileConfirmModal
        action={confirmAction}
        copy={confirmCopy}
        onCancel={() => setConfirmAction(null)}
        onConfirm={executeConfirmedAction}
      />
    </div>
  );
}

function ProfileConfirmModal({
  action,
  copy,
  onCancel,
  onConfirm,
}: {
  action: ConfirmAction | null;
  copy: Record<ConfirmAction, { title: string; description: string; confirmLabel: string }>;
  onCancel: () => void;
  onConfirm: (action: ConfirmAction) => void | Promise<void>;
}) {
  if (!action) return null;

  const meta = copy[action];

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center px-4">
      <div className="absolute inset-0 bg-[rgba(4,7,15,0.45)] backdrop-blur-[2px]" aria-hidden onClick={onCancel} />
      <div className="relative z-10 w-full max-w-md overflow-hidden rounded-2xl border border-[var(--surface-border)] bg-[var(--background)] shadow-[0_18px_48px_rgba(0,0,0,0.18)]">
        <div className="flex items-start justify-between gap-4 px-6 pb-2 pt-5">
          <div className="space-y-1">
            <p className="text-xs font-semibold uppercase tracking-[0.14em] text-[var(--blue-primary-3)]">
              Confirmation
            </p>
            <h3 className="text-lg font-semibold text-[var(--foreground)]">{meta.title}</h3>
            <p className="text-sm text-[var(--muted)]">{meta.description}</p>
          </div>
          <button
            type="button"
            onClick={onCancel}
            className="inline-flex h-9 w-9 items-center justify-center rounded-full border border-[var(--surface-border)] bg-[var(--surface)] text-[var(--foreground)] transition hover:bg-[var(--surface-border)]/60 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)]"
            aria-label="Fermer"
          >
            ✕
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
            onClick={() => onConfirm(action)}
            className="inline-flex items-center justify-center gap-2 rounded-full border border-[var(--blue-primary-2)] bg-[var(--blue-primary-2)] px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:border-[var(--blue-primary-3)] hover:bg-[var(--blue-primary-3)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)]"
          >
            {meta.confirmLabel}
          </button>
        </div>
      </div>
    </div>
  );
}
