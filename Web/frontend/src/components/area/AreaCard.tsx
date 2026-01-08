import Link from "next/link";
import type { ReactNode } from "react";

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
  className,
}: AreaCardProps) {
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
        <div className="flex h-8 w-8 items-center justify-center rounded-full border border-white/30 bg-white/18 text-white/90 text-sm backdrop-blur-[2px]">
          ···
        </div>
      </div>

      <div className="space-y-2">
        <p className="text-lg font-semibold leading-snug">{name}</p>
      </div>

      <div className="mt-auto flex items-center justify-between text-sm text-white/90">
        <div className="flex items-center gap-2">
          <span className="inline-block h-3 w-3 rounded-full bg-white shadow-[0_0_0_2px_rgba(255,255,255,0.18)]" />
          <span>{isActive ? "Actif" : "Inactif"}</span>
        </div>
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
      <button
        type="button"
        onClick={onClick}
        className="group block h-full w-full text-left focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-[var(--foreground)]"
      >
        {card}
      </button>
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
