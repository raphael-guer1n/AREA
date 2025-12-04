import "server-only";

import { cookies } from "next/headers";

const SESSION_COOKIE_NAME = "session";

type SessionCookieOptions = {
  httpOnly?: boolean;
  secure?: boolean;
  path?: string;
};

const defaultOptions: SessionCookieOptions = {
  httpOnly: true,
  secure: process.env.NODE_ENV === "production",
  path: "/",
};

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
