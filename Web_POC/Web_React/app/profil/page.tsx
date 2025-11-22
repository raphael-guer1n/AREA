import Link from "next/link";

type ProfilPageProps = {
    searchParams?: Promise<{
        email?: string;
        password?: string;
    }>;
};

export default async function ProfilPage({ searchParams }: ProfilPageProps) {
    const params = await searchParams;
    const email = params?.email || "";
    const password = params?.password || "";
    const hasCredentials = email !== "" && password !== "";
    const initial = email ? email.charAt(0).toUpperCase() : "?";

    return (
        <div className="flex min-h-screen flex-col" style={{ backgroundColor: "var(--background)" }}>
            <header className="flex items-center justify-between px-10 py-6 text-base">
                <span className="text-3xl font-semibold tracking-[0.4em]">AREA</span>
                <Link
                    href="/"
                    className="rounded-full border px-6 py-2 font-medium uppercase tracking-wide transition hover:bg-[var(--foreground)] hover:text-[var(--background)]"
                    style={{ borderColor: "var(--foreground)", color: "var(--foreground)" }}>
                    Accueil
                </Link>
            </header>

            <main className="flex flex-1 items-center justify-center px-4 pb-12">
                <div
                    className="w-full max-w-lg rounded-3xl p-10 shadow"
                    style={{ backgroundColor: "var(--surface)", color: "var(--foreground)" }}>
                <div className="flex items-center justify-between gap-4">
                    <div>
                        <p className="text-xs uppercase tracking-[0.2em]" style={{ color: "var(--muted)" }}>
                            Espace profil
                        </p>
                        <h1 className="mt-1 text-3xl font-bold">Bienvenue</h1>
                        <p className="mt-2 text-sm" style={{ color: "var(--muted)" }}>
                            Aperçu rapide avec les identifiants saisis.
                        </p>
                    </div>
                    <div
                        className="flex h-14 w-14 items-center justify-center rounded-full text-lg font-semibold"
                        style={{ backgroundColor: "var(--surface-border)", color: "var(--foreground)" }}>
                        {initial}
                    </div>
                </div>

                {hasCredentials ? (
                    <dl className="mt-8 space-y-5 text-sm">
                        <div>
                            <dt className="font-medium text-[var(--muted)]">Email</dt>
                            <dd className="text-lg font-semibold">{email}</dd>
                        </div>
                        <div>
                            <dt className="font-medium text-[var(--muted)]">Mot de passe</dt>
                            <dd className="text-lg font-semibold">{password}</dd>
                        </div>
                    </dl>
                ) : (
                    <p
                        className="mt-8 rounded-lg px-4 py-3 text-sm"
                        style={{ backgroundColor: "var(--surface-border)", color: "var(--foreground)" }}>
                        Aucun identifiant reçu. Retournez au formulaire pour vous connecter.
                    </p>
                )}
                </div>
            </main>
        </div>
    );
}
