import { NextResponse } from "next/server";

type AuthAction = "login" | "register";

type AuthRequestBody = {
  action?: AuthAction;
  email?: string;
  password?: string;
  name?: string;
};

type StoredUser = {
  id: string;
  email: string;
  name: string;
  password: string;
};

// In-memory store for demo purposes only.
const users = new Map<string, StoredUser>();
const MIN_PASSWORD_LENGTH = 8;

function normalizeEmail(value?: string): string {
  return value?.trim().toLowerCase() ?? "";
}

function buildUserResponse(user: StoredUser) {
  return {
    id: user.id,
    email: user.email,
    name: user.name,
    token: `token-${user.id}-${Math.random().toString(36).slice(2, 8)}`,
  };
}

function jsonError(message: string, status: number) {
  return NextResponse.json({ message }, { status });
}

export async function POST(request: Request) {
  const body = (await request.json().catch(() => null)) as AuthRequestBody | null;

  if (!body) {
    return jsonError("Invalid or missing request body.", 400);
  }

  const { action, email, password, name } = body;

  if (action !== "login" && action !== "register") {
    return jsonError("Invalid action. Use 'login' or 'register'.", 400);
  }

  const normalizedEmail = normalizeEmail(email);
  const cleanPassword = password?.trim() ?? "";

  if (!normalizedEmail || !cleanPassword) {
    return jsonError("Email and password are required.", 400);
  }

  if (cleanPassword.length < MIN_PASSWORD_LENGTH) {
    return jsonError(
      `Password too short (min ${MIN_PASSWORD_LENGTH} characters).`,
      400,
    );
  }

  if (action === "register") {
    if (!name?.trim()) {
      return jsonError("Name is required for registration.", 400);
    }

    if (users.has(normalizedEmail)) {
      return jsonError("An account already exists with this email.", 409);
    }

    const newUser: StoredUser = {
      id: `user-${Date.now()}`,
      email: normalizedEmail,
      name: name.trim(),
      password: cleanPassword,
    };

    users.set(normalizedEmail, newUser);

    return NextResponse.json(buildUserResponse(newUser), { status: 201 });
  }

  const existingUser = users.get(normalizedEmail);

  if (!existingUser || existingUser.password !== cleanPassword) {
    return jsonError("Invalid credentials.", 401);
  }

  return NextResponse.json(buildUserResponse(existingUser), { status: 200 });
}
