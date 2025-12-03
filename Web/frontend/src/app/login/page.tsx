import LoginWithGoogle from "@/components/LoginWithGoogle";
import LoginForm from "@/components/forms/LoginForm";
import { Card } from "@/components/ui/Card";

export default function LoginPage() {
  return (
    <main className="flex min-h-screen items-center justify-center bg-[var(--surface)] px-4 py-12">
      <div className="flex w-full max-w-3xl flex-col gap-6">
        <Card
          title="Connexion Google"
          subtitle="Nous récupérons l'URL OAuth2 auprès du backend avant de vous rediriger."
        >
          <LoginWithGoogle />
        </Card>
        <LoginForm />
      </div>
    </main>
  );
}
