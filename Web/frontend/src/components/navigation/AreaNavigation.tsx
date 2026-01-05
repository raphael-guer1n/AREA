"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

import { cn } from "@/lib/helpers";

const navItems = [
  { label: "Services", href: "/services" },
  { label: "Area", href: "/area" },
  { label: "Profil", href: "/profil" },
];

export function AreaNavigation() {
  const pathname = usePathname();

  return (
    <div className="relative w-full max-w-[min(88vw,36rem)]">
      <div
        className="absolute inset-0 -z-10 mx-auto max-w-[min(88vw,36rem)] rounded-3xl bg-[var(--surface)]"
        aria-hidden
      />
      <nav className="flex items-center justify-center gap-3 sm:gap-4 rounded-3xl border border-[var(--surface-border)] bg-[var(--background)] px-[clamp(8px,2.5vw,14px)] py-[clamp(8px,2vw,12px)] shadow-[0_10px_30px_rgba(17,42,70,0.07)] ring-1 ring-[rgba(28,61,99,0.08)]">
        {navItems.map((item) => {
          const isActive = pathname === item.href;
          return (
            <Link
              key={item.href}
              href={item.href}
              aria-current={isActive ? "page" : undefined}
              className={cn(
                "relative inline-flex min-w-[clamp(84px,24vw,132px)] items-center justify-center gap-2 rounded-2xl border px-[clamp(12px,3vw,18px)] py-[clamp(9px,2.4vw,12px)] text-[clamp(0.92rem,1.9vw,1.02rem)] font-semibold transition focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--blue-primary-3)]",
                isActive
                  ? "border-[var(--blue-primary-3)] bg-[var(--blue-primary-3)] text-white shadow-[0_0_0_2px_rgba(28,61,99,0.12)]"
                  : "border-[var(--surface-border)] text-[var(--foreground)] bg-[var(--surface)] hover:border-[var(--blue-primary-2)] hover:bg-[var(--surface)]",
              )}
            >
              <span>{item.label}</span>
            </Link>
          );
        })}
      </nav>
    </div>
  );
}