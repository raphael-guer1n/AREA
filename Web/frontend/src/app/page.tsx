import Link from "next/link";

import { Button } from "@/components/ui/Button";

export default function HomePage() {
  return (
    <main className="flex min-h-screen items-center justify-center px-6 py-12">
      <div className="flex gap-4">
        <Link href="/login">
          <Button variant="secondary">Se connecter</Button>
        </Link>
        <Link href="/register">
          <Button>Créer un compte</Button>
        </Link>
      </div>
    </main>
  );
}
