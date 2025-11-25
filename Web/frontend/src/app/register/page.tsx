import Link from "next/link";

import RegisterForm from "@/components/forms/RegisterForm";

export default function RegisterPage() {
  return (
    <main className="flex min-h-screen items-center justify-center px-4 py-12">
      <div className="w-full max-w-xl space-y-6">
        <RegisterForm />
        <p className="text-center text-sm">
          <Link href="/" className="text-[var(--foreground)] underline">
            Retour à l&apos;accueil
          </Link>
        </p>
      </div>
    </main>
  );
}
