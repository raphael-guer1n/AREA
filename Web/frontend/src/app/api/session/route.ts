import { NextResponse } from "next/server";

import { clearSessionCookie, getSessionToken, hasActiveSession, setSessionCookie } from "@/lib/session";

export const dynamic = "force-dynamic";

export async function GET() {
  const authenticated = await hasActiveSession();
  const token = await getSessionToken();
  return NextResponse.json({ authenticated, token: token ?? null });
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
