import { NextResponse } from "next/server";

import { registerWithCredentials } from "@/lib/api/auth";

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
    const { status, body } = await registerWithCredentials(email, username, password);

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
