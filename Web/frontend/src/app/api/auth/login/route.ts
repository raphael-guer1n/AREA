import { NextResponse } from "next/server";

import { authenticateWithCredentials } from "@/lib/api/auth";

type LoginRequestBody = {
  email?: string;
  emailOrUsername?: string;
  password?: string;
};

export async function POST(request: Request) {
  const body = (await request.json().catch(() => null)) as LoginRequestBody | null;

  const emailOrUsername =
    body?.emailOrUsername?.trim() ?? body?.email?.trim() ?? "";
  const password = body?.password?.trim() ?? "";

  if (!emailOrUsername || !password) {
    return NextResponse.json(
      { success: false, error: "Email ou nom d'utilisateur et mot de passe requis." },
      { status: 400 },
    );
  }

  try {
    const { status, body } = await authenticateWithCredentials(emailOrUsername, password);

    return NextResponse.json(
      body ?? { success: false, error: "Réponse du serveur invalide." },
      { status },
    );
  } catch (error) {
    const message =
      error instanceof Error
        ? error.message
        : "Impossible de contacter le service d'authentification.";
    return NextResponse.json({ success: false, error: message }, { status: 502 });
  }
}
