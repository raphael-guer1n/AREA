"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { FormEvent, useState } from "react";

import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { useAuth } from "@/hooks/useAuth";
import type { RegisterPayload } from "@/types/User";

export default function RegisterForm() {
  const router = useRouter();
  const { register, isLoading, error } = useAuth();
  const [payload, setPayload] = useState<RegisterPayload>({
    email: "",
    password: "",
    name: "",
  });
  const [status, setStatus] = useState<string | null>(null);

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const user = await register(payload);
    if (user) {
      setStatus("Compte créé. Vous pouvez vous connecter.");
      router.push("/login");
    }
  };

  return (
    <Card title="Inscription" subtitle="Créez un compte pour démarrer vos automatisations">
      <form className="space-y-4" onSubmit={handleSubmit}>
        <label className="block text-sm font-medium">
          Nom
          <input
            type="text"
            value={payload.name ?? ""}
            onChange={(event) =>
              setPayload((current) => ({ ...current, name: event.target.value }))
            }
            className="mt-2 w-full rounded-lg border border-[var(--surface-border)] bg-[var(--background)] px-3 py-2 focus:outline-none focus:ring-2 focus:ring-[var(--foreground)]"
            placeholder="Votre nom complet"
          />
        </label>
        <label className="block text-sm font-medium">
          Email
          <input
            type="email"
            required
            value={payload.email}
            onChange={(event) =>
              setPayload((current) => ({ ...current, email: event.target.value }))
            }
            className="mt-2 w-full rounded-lg border border-[var(--surface-border)] bg-[var(--background)] px-3 py-2 focus:outline-none focus:ring-2 focus:ring-[var(--foreground)]"
          />
        </label>
        <label className="block text-sm font-medium">
          Mot de passe
          <input
            type="password"
            required
            value={payload.password}
            onChange={(event) =>
              setPayload((current) => ({
                ...current,
                password: event.target.value,
              }))
            }
            className="mt-2 w-full rounded-lg border border-[var(--surface-border)] bg-[var(--background)] px-3 py-2 focus:outline-none focus:ring-2 focus:ring-[var(--foreground)]"
          />
        </label>
        <Button type="submit" disabled={isLoading} className="w-full">
          {isLoading ? "Création..." : "Créer un compte"}
        </Button>
      </form>
      {error ? <p className="text-sm text-red-500">{error}</p> : null}
      {status ? <p className="text-sm text-emerald-600">{status}</p> : null}
      <p className="text-sm text-[var(--muted)]">
        Déjà inscrit ?{" "}
        <Link href="/login" className="text-[var(--foreground)] underline">
          Se connecter
        </Link>
      </p>
    </Card>
  );
}
