"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { FormEvent, useState } from "react";

import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { useAuth } from "@/hooks/useAuth";
import type { LoginPayload } from "@/types/User";

export default function LoginForm() {
  const router = useRouter();
  const { login, isLoading, error } = useAuth();
  const [credentials, setCredentials] = useState<LoginPayload>({
    email: "",
    password: "",
  });
  const [status, setStatus] = useState<string | null>(null);

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const user = await login(credentials);
    if (user) {
      setStatus("Connexion réussie");
      router.push(
        `/profile?email=${encodeURIComponent(user.email)}&name=${encodeURIComponent(user.name ?? "")}`,
      );
    }
  };

  return (
    <Card title="Connexion" subtitle="Accédez à votre espace AREA">
      <form className="space-y-4" onSubmit={handleSubmit}>
        <label className="block text-sm font-medium">
          Email
          <input
            type="email"
            required
            value={credentials.email}
            onChange={(event) =>
              setCredentials((current) => ({
                ...current,
                email: event.target.value,
              }))
            }
            className="mt-2 w-full rounded-lg border border-[var(--surface-border)] bg-[var(--background)] px-3 py-2 focus:outline-none focus:ring-2 focus:ring-[var(--foreground)]"
          />
        </label>
        <label className="block text-sm font-medium">
          Mot de passe
          <input
            type="password"
            required
            value={credentials.password}
            onChange={(event) =>
              setCredentials((current) => ({
                ...current,
                password: event.target.value,
              }))
            }
            className="mt-2 w-full rounded-lg border border-[var(--surface-border)] bg-[var(--background)] px-3 py-2 focus:outline-none focus:ring-2 focus:ring-[var(--foreground)]"
          />
        </label>
        <Button type="submit" disabled={isLoading} className="w-full">
          {isLoading ? "Connexion..." : "Se connecter"}
        </Button>
      </form>
      {error ? <p className="text-sm text-red-500">{error}</p> : null}
      {status ? <p className="text-sm text-emerald-600">{status}</p> : null}
      <p className="text-sm text-[var(--muted)]">
        Pas encore inscrit ?{" "}
        <Link href="/register" className="text-[var(--foreground)] underline">
          Créer un compte
        </Link>
      </p>
    </Card>
  );
}
