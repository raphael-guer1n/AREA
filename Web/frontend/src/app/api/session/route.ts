import { NextResponse } from "next/server";

import { clearSessionCookie, hasActiveSession } from "@/lib/session";

export const dynamic = "force-dynamic";

export async function GET() {
  const authenticated = await hasActiveSession();
  return NextResponse.json({ authenticated });
}

export async function DELETE() {
  await clearSessionCookie();
  return NextResponse.json({ success: true });
}
