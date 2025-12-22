// Server-only helpers for reading/writing the session cookie via Next headers API.
import "server-only";

import { cookies } from "next/headers";

export const SESSION_COOKIE_NAME = "session";

type SessionCookieOptions = {
  httpOnly?: boolean;
  secure?: boolean;
  sameSite?: "lax" | "strict" | "none";
  path?: string;
  maxAge?: number;
};

const isSecureCookie =
  process.env.COOKIE_SECURE === "true" ||
  process.env.SESSION_COOKIE_SECURE === "true" ||
  (process.env.NODE_ENV === "production" &&
    Boolean(process.env.NEXT_PUBLIC_SITE_URL?.startsWith("https")));

const defaultOptions: SessionCookieOptions = {
  httpOnly: true,
  secure: isSecureCookie,
  sameSite: "lax",
  path: "/",
  // 7 days
  maxAge: 60 * 60 * 24 * 7,
};

export function sessionCookieOptions(
  overrides: SessionCookieOptions = {},
): SessionCookieOptions {
  return { ...defaultOptions, ...overrides };
}

export async function setSessionCookie(
  token: string,
  options: SessionCookieOptions = defaultOptions,
): Promise<void> {
  const cookieStore = await cookies();
  cookieStore.set(SESSION_COOKIE_NAME, token, { ...defaultOptions, ...options });
}

export async function clearSessionCookie(): Promise<void> {
  const cookieStore = await cookies();
  cookieStore.delete(SESSION_COOKIE_NAME);
}

export async function getSessionToken(): Promise<string | null> {
  const cookieStore = await cookies();
  return cookieStore.get(SESSION_COOKIE_NAME)?.value ?? null;
}

export async function hasActiveSession(): Promise<boolean> {
  return Boolean(await getSessionToken());
}
