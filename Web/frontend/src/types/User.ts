export type User = {
  id: string;
  email: string;
  username?: string;
  name?: string;
  avatarUrl?: string;
  token?: string;
};

export type LoginPayload = {
  email: string;
  password: string;
};

export type RegisterPayload = LoginPayload & {
  name?: string;
};

export type AuthState = {
  user: User | null;
  token: string | null;
};
