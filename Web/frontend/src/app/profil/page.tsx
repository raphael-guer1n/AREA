import { AreaNavigation } from "@/components/navigation/AreaNavigation";
import { getSessionToken } from "@/lib/session";
import { fetchAuthenticatedUser } from "@/lib/api/auth";
import type { User } from "@/types/User";
import type { AuthSession } from "@/types/auth";
import { ProfileClient } from "./profile-client";

export const dynamic = "force-dynamic";

export default async function ProfilPage() {
  let initialUser: User | null = null;
  let initialSession: AuthSession | null = null;

  const token = await getSessionToken();
  if (token) {
    initialSession = { token };
    try {
      initialUser = await fetchAuthenticatedUser(token);
    } catch {
      // keep token; client will validate/refresh
    }
  }

  return (
    <main className="relative flex min-h-screen justify-center overflow-hidden bg-[var(--surface)] px-6 py-12 pt-10 text-[var(--foreground)]">
      <div className="pointer-events-none absolute inset-0 -z-10 opacity-45">
        <div
          className="absolute -left-12 top-10 h-64 w-64 rounded-full bg-[radial-gradient(circle_at_center,var(--accent)_0,transparent_65%)] blur-3xl"
          aria-hidden
        />
        <div
          className="absolute right-0 top-20 h-64 w-64 rounded-full bg-[radial-gradient(circle_at_center,var(--blue-primary-2)_0,transparent_70%)] blur-3xl"
          aria-hidden
        />
        <div
          className="absolute bottom-[-14%] left-1/3 h-60 w-60 rounded-full bg-[radial-gradient(circle_at_center,var(--blue-primary-1)_0,transparent_65%)] blur-3xl"
          aria-hidden
        />
      </div>

      <div className="relative w-full max-w-5xl space-y-8">
        <div className="flex justify-center">
          <AreaNavigation />
        </div>

        <ProfileClient initialUser={initialUser} initialSession={initialSession} />
      </div>
    </main>
  );
}
