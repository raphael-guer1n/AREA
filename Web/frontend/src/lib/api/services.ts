type ProvidersResponse<T> = {
  success?: boolean;
  data?: {
    providers?: T;
  };
  error?: string;
};

export type ProviderStatus = {
  provider: string;
  is_logged: boolean;
};

type FetchOptions = RequestInit & { token?: string };

async function parseProvidersResponse<T>(
  response: Response,
): Promise<ProvidersResponse<T>["data"]["providers"]> {
  const body = (await response.json().catch(() => null)) as ProvidersResponse<T> | null;

  if (!body?.success || !body.data?.providers) {
    throw new Error(body?.error ?? "Unable to retrieve the services list.");
  }

  return body.data.providers;
}

export async function fetchUserProviders(
  userId: string | number,
  token?: string,
): Promise<ProviderStatus[]> {
  try {
    const response = await fetch(`/api/auth-providers/${userId}`, {
      method: "GET",
      cache: "no-store",
      headers: token ? { Authorization: `Bearer ${token}` } : undefined,
    });

    const providers = await parseProvidersResponse<ProviderStatus[]>(response);

    return providers.map((provider) => ({
      provider: provider.provider,
      is_logged: Boolean(provider.is_logged),
    }));
  } catch (error) {
    const fallbackProviders = await fetchAvailableProviders(token);
    return fallbackProviders.map((provider) => ({
      provider,
      is_logged: false,
    }));
  }
}

export async function fetchAvailableProviders(token?: string): Promise<string[]> {
  const response = await fetch("/api/auth-providers", {
    method: "GET",
    cache: "no-store",
    headers: token ? { Authorization: `Bearer ${token}` } : undefined,
  });

  return parseProvidersResponse<string[]>(response);
}
