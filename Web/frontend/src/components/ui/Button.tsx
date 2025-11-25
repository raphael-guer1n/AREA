"use client";

import type { ButtonHTMLAttributes } from "react";

import { cn } from "@/lib/helpers";

type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: "primary" | "secondary" | "ghost";
};

export function Button({ variant = "primary", className, ...props }: ButtonProps) {
  const baseClasses =
    "inline-flex items-center justify-center gap-2 rounded-full px-4 py-2 text-sm font-semibold transition focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2";
  const variants = {
    primary:
      "bg-[var(--foreground)] text-[var(--background)] hover:opacity-90 focus-visible:outline-[var(--foreground)]",
    secondary:
      "border border-[var(--surface-border)] text-[var(--foreground)] hover:bg-[var(--surface-border)] focus-visible:outline-[var(--foreground)]",
    ghost:
      "text-[var(--foreground)] hover:bg-[var(--surface-border)] focus-visible:outline-[var(--foreground)]",
  };

  return (
    <button className={cn(baseClasses, variants[variant], className)} {...props} />
  );
}
