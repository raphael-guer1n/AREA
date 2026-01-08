// Proxy login route: forwards credentials to the auth service and relays its response.
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
      { success: false, error: "Email or username and password are required." },
      { status: 400 },
    );
  }

  try {
    const { status, body } = await authenticateWithCredentials(emailOrUsername, password);

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
