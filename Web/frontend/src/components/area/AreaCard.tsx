import Link from "next/link";
import {
  useEffect,
  useRef,
  useState,
  type MouseEvent as ReactMouseEvent,
  type ReactNode,
} from "react";

import { cn } from "@/lib/helpers";

type AreaCardProps = {
  id: string | number;
  name: string;
  actionLabel: string;
  reactionLabel: string;
  actionIcon: ReactNode;
  reactionIcon: ReactNode;
  isActive?: boolean;
  gradientFrom?: string;
  gradientTo?: string;
  lastRun?: string;
  href?: string;
  onClick?: () => void;
  onActivate?: () => void;
  onDeactivate?: () => void;
  onDelete?: () => void;
  isBusy?: boolean;
  className?: string;
};

export function AreaCard({
  id,
  name,
  actionLabel: _actionLabel,
  reactionLabel: _reactionLabel,
  actionIcon,
  reactionIcon,
  isActive = false,
  gradientFrom = "#002642",
  gradientTo = "#e59500",
  lastRun,
  href,
  onClick,
  onActivate,
  onDeactivate,
  onDelete,
  isBusy = false,
  className,
}: AreaCardProps) {
  const menuRef = useRef<HTMLDivElement>(null);
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  useEffect(() => {
    if (!isMenuOpen) return;

    const handleClickOutside = (event: MouseEvent) => {
      if (!menuRef.current) return;
      if (event.target instanceof Node && !menuRef.current.contains(event.target)) {
        setIsMenuOpen(false);
      }
    };
    const handleEscape = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        setIsMenuOpen(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    document.addEventListener("keydown", handleEscape);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
      document.removeEventListener("keydown", handleEscape);
    };
  }, [isMenuOpen]);

  const toggleMenu = (event: ReactMouseEvent) => {
    event.preventDefault();
    event.stopPropagation();
    setIsMenuOpen((prev) => !prev);
  };

  const handleMenuAction =
    (callback?: () => void) =>
    (event: ReactMouseEvent) => {
      event.preventDefault();
      event.stopPropagation();
      setIsMenuOpen(false);
      callback?.();
    };

  const card = (
    <article
      className={cn(
        "relative flex h-full flex-col gap-3.5 overflow-hidden rounded-xl px-4 py-4 text-white shadow-[0_10px_35px_rgba(0,0,0,0.08)] transition duration-200 hover:-translate-y-1 hover:shadow-[0_16px_40px_rgba(0,0,0,0.14)] aspect-[16/10]",
        className,
      )}
      style={{
        background: `linear-gradient(135deg, ${gradientFrom}, ${gradientTo})`,
      }}
    >
      <div className="flex items-start justify-between">
        <div className="flex items-center gap-2.5">
          <IconBadge>{actionIcon}</IconBadge>
          <span className="text-white/70">→</span>
          <IconBadge>{reactionIcon}</IconBadge>
        </div>
        <div className="flex items-center gap-2">
          <div className="relative" ref={menuRef}>
            <button
              type="button"
              onClick={toggleMenu}
              className="inline-flex h-9 w-9 items-center justify-center rounded-full border border-white/30 bg-white/18 text-white/90 shadow-[0_0_0_1px_rgba(255,255,255,0.12)] backdrop-blur-[2px] transition hover:scale-[1.03] hover:border-white/50 hover:bg-white/25 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
              aria-haspopup="menu"
              aria-expanded={isMenuOpen}
              aria-label="Actions sur l'area"
            >
              <OptionsIcon />
            </button>
            {isMenuOpen ? (
              <div className="absolute right-0 top-11 z-20 w-44 rounded-xl border border-white/20 bg-[rgba(6,14,25,0.8)] p-1.5 text-sm shadow-[0_16px_40px_rgba(0,0,0,0.25)] backdrop-blur">
                <MenuItem
                  label="Activer"
                  onClick={handleMenuAction(onActivate)}
                  disabled={!onActivate || isActive || isBusy}
                />
                <MenuItem
                  label="Désactiver"
                  onClick={handleMenuAction(onDeactivate)}
                  disabled={!onDeactivate || !isActive || isBusy}
                />
                <MenuItem
                  label="Supprimer"
                  onClick={handleMenuAction(onDelete)}
                  disabled={!onDelete}
                  muted={!onDelete}
                  title={!onDelete ? "Suppression non disponible" : undefined}
                />
              </div>
            ) : null}
          </div>
        </div>
      </div>

      <div className="space-y-2">
        <p className="text-lg font-semibold leading-snug">{name}</p>
      </div>

      <div className="mt-auto flex items-center justify-between text-sm text-white/90">
        <StatusPill isActive={isActive} />
        {lastRun ? (
          <span className="inline-flex items-center gap-1 rounded-full border border-white/30 bg-white/12 px-3 py-1 text-xs font-semibold uppercase tracking-wide text-white shadow-[0_0_0_1px_rgba(255,255,255,0.12)]">
            {lastRun}
          </span>
        ) : null}
      </div>
    </article>
  );

  if (onClick) {
    return (
      <div
        role="button"
        tabIndex={0}
        onClick={onClick}
        onKeyDown={(event) => {
          if (event.key === "Enter" || event.key === " ") {
            event.preventDefault();
            onClick();
          }
        }}
        className="group block h-full w-full text-left focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-[var(--foreground)]"
      >
        {card}
      </div>
    );
  }

  return (
    <Link
      href={href ?? `/area/${id}`}
      className="group block h-full focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-[var(--foreground)]"
    >
      {card}
    </Link>
  );
}

function IconBadge({ children }: { children: ReactNode }) {
  return (
    <span className="flex h-10 w-10 items-center justify-center rounded-full border border-white/35 bg-white/18 text-base font-semibold text-white shadow-[0_8px_18px_rgba(0,0,0,0.12)] backdrop-blur-[2px]">
      {children}
    </span>
  );
}

function StatusPill({ isActive }: { isActive: boolean }) {
  return (
    <span className="inline-flex items-center gap-1.5 rounded-full border border-white/25 bg-white/15 px-2.5 py-1 text-[11px] font-semibold uppercase tracking-wide text-white/90 shadow-[0_0_0_1px_rgba(255,255,255,0.12)]">
      <span
        className={cn(
          "block h-2 w-2 rounded-full shadow-[0_0_0_2px_rgba(255,255,255,0.18)]",
          isActive ? "bg-emerald-300" : "bg-white/70",
        )}
      />
      {isActive ? "Active" : "Inactive"}
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

function MenuItem({
  label,
  onClick,
  disabled = false,
  muted = false,
  title,
}: {
  label: string;
  onClick: (event: ReactMouseEvent) => void;
  disabled?: boolean;
  muted?: boolean;
  title?: string;
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      disabled={disabled}
      title={title}
      className={cn(
        "flex w-full items-center justify-between rounded-lg px-3 py-2 text-left transition",
        muted ? "text-white/50" : "text-white",
        disabled
          ? "cursor-not-allowed opacity-50"
          : "hover:bg-white/10 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white/70",
      )}
    >
      <span>{label}</span>
      {disabled ? <span className="text-[10px] uppercase tracking-wide">Indispo</span> : null}
    </button>
  );
}
