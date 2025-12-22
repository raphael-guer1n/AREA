// Proxy registration route: forwards user creation to the auth service and returns the result.
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
      { success: false, error: "Email, username, and password are required." },
      { status: 400 },
    );
  }

  try {
    const { status, body } = await registerWithCredentials(email, username, password);

    return NextResponse.json(
      body ?? { success: false, error: "Invalid server response." },
      { status },
    );
  } catch (error) {
    const message =
      error instanceof Error
        ? error.message
        : "Unable to reach the authentication service.";
    return NextResponse.json({ success: false, error: message }, { status: 502 });
  }
}
