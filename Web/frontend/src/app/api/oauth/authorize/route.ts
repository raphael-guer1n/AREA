import { NextResponse } from "next/server";

import { AUTH_SERVICE_BASE_URL } from "@/lib/api/auth";

export async function GET(request: Request) {
  const { searchParams } = new URL(request.url);
  const provider = searchParams.get("provider");
  const mode = searchParams.get("mode") ?? "login";
  const callbackUrl = searchParams.get("callback_url") ?? "";
  const platform = searchParams.get("platform") ?? "";

  if (!provider) {
    return NextResponse.json({ success: false, error: "provider is required" }, { status: 400 });
  }

  // Choose upstream path based on mode
  const basePath =
    mode === "link" ? "/auth/oauth2/authorize" : "/auth/oauth2/login";

  const upstreamParams = new URLSearchParams({ provider });
  if (callbackUrl) upstreamParams.set("callback_url", callbackUrl);
  if (platform) upstreamParams.set("platform", platform);

  const upstreamUrl = `${AUTH_SERVICE_BASE_URL}${basePath}?${upstreamParams.toString()}`;

  try {
    const response = await fetch(upstreamUrl, {
      method: "GET",
      headers: (() => {
        const headers: Record<string, string> = {};
        const auth = request.headers.get("authorization");
        if (auth) headers["Authorization"] = auth;
        return headers;
      })(),
      cache: "no-store",
      redirect: "manual",
    });

    const body = await response.json().catch(() => null);

    return NextResponse.json(body ?? { success: false, error: "Invalid upstream response" }, {
      status: response.status,
    });
  } catch (error) {
    const message = error instanceof Error ? error.message : "Unable to reach auth service.";
    return NextResponse.json({ success: false, error: message }, { status: 502 });
  }
}
