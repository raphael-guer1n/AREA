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
    <main className="flex min-h-screen items-start justify-center bg-[var(--surface)] px-6 py-12 pt-10 text-[var(--foreground)]">
      <div className="flex w-full max-w-5xl flex-col gap-8">
        <div className="flex justify-center">
          <AreaNavigation />
        </div>

        <ProfileClient initialUser={initialUser} initialSession={initialSession} />
      </div>
    </main>
  );
}
