"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { FormEvent, useState } from "react";

export default function LoginPage() {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [status, setStatus] = useState<string | null>(null);
    const router = useRouter();

    const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        if (email === "test@test.com" && password === "0000") {
            setStatus("Connexion réussie, redirection vers le profil...");
            const query = new URLSearchParams({ email, password }).toString();
            router.push(`/profil?${query}`);
        } else {
            setStatus("Identifiants incorrects");
        }
    };

    return (
        <main
            className="flex min-h-screen items-center justify-center px-4"
            style={{ backgroundColor: "var(--background)" }}>
            <div
                className="w-full max-w-md rounded-2xl p-8 shadow"
                style={{ backgroundColor: "var(--surface)", color: "var(--foreground)" }}>
                <h1 className="text-2xl font-bold">Connexion</h1>
                <p className="mt-2 text-sm" style={{ color: "var(--muted)" }}>
                    Entrez vos identifiants pour accéder à Area
                </p>
                <form className="mt-7 space-y-6" onSubmit={handleSubmit}>
                    <label className="block text-sm font-medium" style={{ color: "var(--foreground)" }}>
                        Email
                        <input
                            type="email"
                            required
                            value={email}
                            onChange={(event) => setEmail(event.target.value)}
                            className="mt-2 w-full rounded-lg border px-4 py-2 text-base focus:outline-none focus:ring-2 focus:ring-[var(--foreground)]"
                            style={{
                                backgroundColor: "var(--surface)",
                                borderColor: "var(--surface-border)",
                                color: "var(--foreground)",
                            }}/>
                    </label>
                    <label className="block text-sm font-medium" style={{ color: "var(--foreground)" }}>
                        Mot de passe
                        <input
                            type="password"
                            required
                            value={password}
                            onChange={(event) => setPassword(event.target.value)}
                            className="mt-2 w-full rounded-lg border px-4 py-2 text-base focus:outline-none focus:ring-2 focus:ring-[var(--foreground)]"
                            style={{
                                backgroundColor: "var(--surface)",
                                borderColor: "var(--surface-border)",
                                color: "var(--foreground)",
                            }}
                        />
                    </label>
                    <button
                        type="submit"
                        className="w-full rounded-full px-6 py-3 text-base font-semibold transition hover:opacity-90"
                        style={{ backgroundColor: "var(--accent)", color: "var(--background)" }}>
                        Se connecter
                    </button>
                </form>
                {status && (
                    <p
                        className="mt-4 rounded-md px-4 py-2 text-sm"
                        style={{ backgroundColor: "var(--surface-border)", color: "var(--foreground)" }}>
                        {status}
                    </p>
                )}
                <div className="mt-6 text-center text-sm" style={{ color: "var(--muted)" }}>
                    <Link href="/" className="font-semibold text-[var(--foreground)] hover:underline">
                        Retour à l&apos;accueil
                    </Link>
                </div>
            </div>
        </main>
    );
}
