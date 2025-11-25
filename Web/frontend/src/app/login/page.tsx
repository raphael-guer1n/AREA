import Link from "next/link";

import LoginForm from "@/components/forms/LoginForm";

export default function LoginPage() {
  return (
    <main className="flex min-h-screen items-center justify-center px-4 py-12">
      <div className="w-full max-w-xl space-y-6">
        <LoginForm />
        <p className="text-center text-sm text-[var(--muted)]">
          Pas encore de compte ?{" "}
          <Link
            href="/register"
            className="font-semibold text-[var(--foreground)] underline"
          >
            Créer un compte
          </Link>
        </p>
        <p className="text-center text-sm">
          <Link href="/" className="text-[var(--foreground)] underline">
            Retour à l&apos;accueil
          </Link>
        </p>
      </div>
    </main>
  );
}
