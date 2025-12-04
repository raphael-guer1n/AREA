import { NextResponse } from "next/server";

import { BACKEND_BASE_URL } from "@/lib/api/auth";

type RegisterRequestBody = {
  email?: string;
  username?: string;
  name?: string;
  password?: string;
};

export async function POST(request: Request) {
  const body = (await request.json().catch(() => null)) as RegisterRequestBody | null;

  const email = body?.email?.trim() ?? "";
  const username =
    body?.username?.trim() ?? body?.name?.trim() ?? body?.email?.trim() ?? "";
  const password = body?.password?.trim() ?? "";

  if (!email || !username || !password) {
    return NextResponse.json(
      { success: false, error: "Email, nom d'utilisateur et mot de passe requis." },
      { status: 400 },
    );
  }

  try {
    const backendResponse = await fetch(`${BACKEND_BASE_URL}/auth/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, username, password }),
      cache: "no-store",
    });

    const backendBody = (await backendResponse.json().catch(() => null)) as
      | Record<string, unknown>
      | null;

    if (!backendBody) {
      return NextResponse.json(
        { success: false, error: "Réponse du serveur invalide." },
        { status: 502 },
      );
    }

    return NextResponse.json(backendBody, { status: backendResponse.status });
  } catch (error) {
    const message =
      error instanceof Error
        ? error.message
        : "Impossible de contacter le service d'authentification.";
    return NextResponse.json({ success: false, error: message }, { status: 502 });
  }
}
