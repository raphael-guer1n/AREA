import { NextResponse } from "next/server";

import { fetchAuthenticatedUser } from "@/lib/api/auth";
import { clearSessionCookie, getSessionToken, setSessionCookie } from "@/lib/session";

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
    await clearSessionCookie();
    const message =
      error instanceof Error ? error.message : "Session expirée. Merci de vous reconnecter.";
    return NextResponse.json(
      { authenticated: false, token: null, user: null, error: message },
      { status: 401 },
    );
  }
}

export async function POST(request: Request) {
  const body = (await request.json().catch(() => null)) as { token?: string } | null;

  if (!body?.token) {
    return NextResponse.json({ message: "Token manquant." }, { status: 400 });
  }

  await setSessionCookie(body.token);
  return NextResponse.json({ success: true });
}

export async function DELETE() {
  await clearSessionCookie();
  return NextResponse.json({ success: true });
}
