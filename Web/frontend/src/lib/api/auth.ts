import { sleep } from '@/lib/helpers';
import type { LoginPayload, RegisterPayload, User } from '@/types/User';

export async function loginRequest(payload: LoginPayload): Promise<User> {
  await sleep(300);
  if (!payload.email || !payload.password) {
    throw new Error('Email et mot de passe requis');
  }

  return {
    id: `user-${payload.email}`,
    email: payload.email,
    name: payload.email.split('@')[0],
    token: 'demo-token',
  };
}

export async function registerRequest(payload: RegisterPayload): Promise<User> {
  await sleep(350);
  if (!payload.email || !payload.password) {
    throw new Error('Email et mot de passe requis');
  }

  return {
    id: `user-${Date.now()}`,
    email: payload.email,
    name: payload.name ?? payload.email.split('@')[0],
    token: 'demo-token',
  };
}

export async function logoutRequest(): Promise<void> {
  await sleep(80);
}
