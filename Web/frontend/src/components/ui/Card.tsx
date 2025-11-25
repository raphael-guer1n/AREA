"use client";

import type { ReactNode } from "react";

import { cn } from "@/lib/helpers";

type CardProps = {
  title?: string;
  subtitle?: string;
  action?: ReactNode;
  children: ReactNode;
  className?: string;
};

export function Card({
  title,
  subtitle,
  action,
  children,
  className,
}: CardProps) {
  return (
    <div
      className={cn(
        "rounded-2xl border border-[var(--surface-border)] bg-[var(--surface)] p-6 shadow-sm",
        className,
      )}
    >
      {(title || subtitle || action) && (
        <div className="mb-4 flex items-start justify-between gap-4">
          <div className="space-y-1">
            {title ? <h2 className="text-xl font-semibold">{title}</h2> : null}
            {subtitle ? (
              <p className="text-sm text-[var(--muted)]">{subtitle}</p>
            ) : null}
          </div>
          {action ? <div className="shrink-0">{action}</div> : null}
        </div>
      )}
      <div className="space-y-3 text-sm leading-relaxed">{children}</div>
    </div>
  );
}
