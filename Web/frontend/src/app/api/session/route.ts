// Internal session endpoint: stores/clears/reads the auth token via HTTP-only cookie
// and validates it against the auth service.
import { NextResponse } from "next/server";

import { fetchAuthenticatedUser } from "@/lib/api/auth";
import { SESSION_COOKIE_NAME, getSessionToken, sessionCookieOptions } from "@/lib/session";

export const dynamic = "force-dynamic";

export async function GET() {
  const token = await getSessionToken();

  if (!token) {
    return NextResponse.json({ authenticated: false, token: null, user: null });
  }

  try {
    const user = await fetchAuthenticatedUser(token);
    return NextResponse.json({ authenticated: true, token, user });
  } catch (error) {
    return NextResponse.json(
      {
        authenticated: false,
        token,
        user: null,
        error: error instanceof Error ? error.message : "Session expired. Please sign in again.",
      },
      { status: 401 },
    );
  }
}

export async function POST(request: Request) {
  const body = (await request.json().catch(() => null)) as { token?: string } | null;

  if (!body?.token) {
    return NextResponse.json({ message: "Missing token." }, { status: 400 });
  }

  const response = NextResponse.json({ success: true });
  response.cookies.set(SESSION_COOKIE_NAME, body.token, sessionCookieOptions());
  return response;
}

export async function DELETE() {
  const response = NextResponse.json({ success: true });
  response.cookies.set(SESSION_COOKIE_NAME, "", sessionCookieOptions({ maxAge: 0 }));
  return response;
}
