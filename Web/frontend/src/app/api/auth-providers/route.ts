import { NextResponse } from "next/server";

import { AUTH_SERVICE_BASE_URL } from "@/lib/api/auth";

export async function GET(request: Request) {
  const authHeader = request.headers.get("authorization") ?? undefined;

  try {
    const response = await fetch(`${AUTH_SERVICE_BASE_URL}/oauth2/providers`, {
      method: "GET",
      headers: authHeader ? { Authorization: authHeader } : undefined,
      cache: "no-store",
    });

    const body = await response.json().catch(() => null);

    if (!response.ok) {
      return NextResponse.json(body ?? { error: "Upstream error" }, {
        status: response.status,
      });
    }

    return NextResponse.json(body);
  } catch (error) {
    const message = error instanceof Error ? error.message : "Unable to reach auth service.";
    return NextResponse.json({ success: false, error: message }, { status: 502 });
  }
}
