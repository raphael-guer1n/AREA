"use client";

import type { ReactNode } from "react";
import { useEffect } from "react";
import { usePathname, useRouter } from "next/navigation";

import { useAuth } from "@/hooks/useAuth";

const PUBLIC_ROUTES = new Set(["/", "/login", "/register", "/auth/callback"]);

function isPublicRoute(pathname: string | null) {
  if (!pathname) return true;
  if (PUBLIC_ROUTES.has(pathname)) return true;

  return Array.from(PUBLIC_ROUTES)
    .filter((route) => route !== "/")
    .some((route) => pathname.startsWith(`${route}/`));
}

type AuthGuardProps = {
  children: ReactNode;
};

export function AuthGuard({ children }: AuthGuardProps) {
  const pathname = usePathname();
  const router = useRouter();
  const { status } = useAuth();

  const isPublic = isPublicRoute(pathname);
  const isAuthenticated = status === "authenticated";

  useEffect(() => {
    if (isPublic || isAuthenticated) return;

    if (status === "unauthenticated" || status === "error") {
      router.replace("/login");
    }
  }, [isAuthenticated, isPublic, router, status]);

  if (!isPublic && !isAuthenticated) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-[var(--surface)] text-[var(--foreground)]">
        <p>Chargement de la session...</p>
      </div>
    );
  }

  return <>{children}</>;
}
