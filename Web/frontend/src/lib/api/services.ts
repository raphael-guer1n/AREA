import { BACKEND_BASE_URL } from "@/lib/api/auth";

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

async function parseProvidersResponse<T>(
  response: Response,
): Promise<ProvidersResponse<T>["data"]["providers"]> {
  const body = (await response.json().catch(() => null)) as ProvidersResponse<T> | null;

  if (!body?.success || !body.data?.providers) {
    throw new Error(body?.error ?? "Unable to retrieve the services list.");
  }

  return body.data.providers;
}

async function fetchWithFallback(
  paths: string[],
  options?: RequestInit,
): Promise<Response> {
  let lastError: Error | null = null;

  for (const path of paths) {
    try {
      const response = await fetch(`${BACKEND_BASE_URL}${path}`, options);

      if (response.ok) {
        return response;
      }

      if (response.status === 404) {
        continue;
      }

      const fallbackMessage = await response
        .json()
        .then((data) => (data?.error as string | undefined) ?? response.statusText)
        .catch(() => response.statusText);

      throw new Error(fallbackMessage || "Unable to reach the services endpoint.");
    } catch (error) {
      lastError =
        error instanceof Error ? error : new Error("Unable to reach the services endpoint.");
    }
  }

  throw lastError ?? new Error("Unable to reach the services endpoint.");
}

export async function fetchUserProviders(userId: string | number): Promise<ProviderStatus[]> {
  const response = await fetchWithFallback(
    [`/auth/oauth2/providers/${userId}`, `/oauth2/providers/${userId}`],
    {
      method: "GET",
      cache: "no-store",
    },
  );

  const providers = await parseProvidersResponse<ProviderStatus[]>(response);

  return providers.map((provider) => ({
    provider: provider.provider,
    is_logged: Boolean(provider.is_logged),
  }));
}

export async function fetchAvailableProviders(): Promise<string[]> {
  const response = await fetchWithFallback(
    ["/auth/oauth2/providers", "/oauth2/providers"],
    {
      method: "GET",
      cache: "no-store",
    },
  );

  return parseProvidersResponse<string[]>(response);
}
