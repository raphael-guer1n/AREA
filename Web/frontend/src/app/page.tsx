"use client";

import Image from "next/image";
import Link from "next/link";
import { useEffect, useState } from "react";

import { Button } from "@/components/ui/Button";

const services = [
  { name: "Gmail / Outlook", emoji: "✉️", color: "var(--card-color-3)" },
  { name: "Slack / Discord", emoji: "💬", color: "var(--card-color-2)" },
  { name: "Drive / Notion", emoji: "📁", color: "var(--card-color-1)" },
  { name: "Trello / Jira", emoji: "✅", color: "var(--card-color-4)" },
  { name: "X / LinkedIn", emoji: "🐦", color: "var(--card-color-5)" },
  { name: "Twilio / SMS", emoji: "🎙️", color: "var(--card-color-3)" },
  { name: "Et plus", emoji: "+40", color: "var(--card-color-2)" },
];

const stats = [
  { value: "40+", label: "Connecteurs prêts à l'emploi" },
  { value: "< 5 min", label: "Pour mettre en prod une AREA" },
  { value: "24/7", label: "Hooks surveillés en continu" },
];

const features = [
  {
    title: "Builder no-code guidé",
    body: "Assemblez vos AREAs en if/then, champs préremplis, validations et exemples de payload pour aller vite.",
    icon: "⚡",
  },
  {
    title: "Hooks fiabilisés",
    body: "Webhooks, polling et retries avec backoff pour ne pas rater vos événements critiques.",
    icon: "🛡️",
  },
  {
    title: "Auth centralisée",
    body: "OAuth2, renouvellement des tokens et stockage chiffré pour chaque compte connecté.",
    icon: "🎛️",
  },
  {
    title: "Suivi clair",
    body: "Logs horodatés, statut d'exécution et alertes pour diagnostiquer et corriger rapidement.",
    icon: "🔗",
  },
];

const steps = [
  {
    id: "01",
    title: "Connecter vos services",
    body: "Authentifiez Google, Discord, Trello... via OAuth2 sécurisé et scopes maîtrisés.",
    color: "var(--card-color-1)",
  },
  {
    id: "02",
    title: "Définir l'Action et la REAction",
    body: "Choisissez le déclencheur, mappez les champs et ajoutez filtres ou conditions.",
    color: "var(--card-color-2)",
  },
  {
    id: "03",
    title: "Laisser tourner",
    body: "AREA surveille les hooks, exécute en continu et vous alerte en cas d'échec.",
    color: "var(--card-color-3)",
  },
];

const gridPattern = {
  backgroundImage:
    "linear-gradient(to right, rgba(15,23,42,0.04) 1px, transparent 1px), linear-gradient(to bottom, rgba(15,23,42,0.04) 1px, transparent 1px)",
  backgroundSize: "72px 72px",
};

export default function HomePage() {
  const [scrollDir, setScrollDir] = useState<"up" | "down" | null>(null);
  const [isTop, setIsTop] = useState(true);

  useEffect(() => {
    let lastY = window.scrollY;

    const onScroll = () => {
      const current = window.scrollY;
      if (Math.abs(current - lastY) < 2) return;
      setScrollDir(current > lastY ? "down" : "up");
      const trigger = Math.max(48, window.innerHeight * 0.4);
      setIsTop(current < trigger);
      lastY = current;
    };

    window.addEventListener("scroll", onScroll, { passive: true });
    return () => window.removeEventListener("scroll", onScroll);
  }, []);

  useEffect(() => {
    const revealTargets = Array.from(document.querySelectorAll<HTMLElement>(".reveal"));
    if (revealTargets.length === 0) return;

    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            entry.target.classList.add("revealed");
            observer.unobserve(entry.target);
          }
        });
      },
      { threshold: 0.2 }
    );

    revealTargets.forEach((el) => observer.observe(el));
    return () => observer.disconnect();
  }, []);

  return (
    <div className="relative min-h-screen overflow-hidden bg-[var(--background)] text-[var(--foreground)]">
      <header
        className={`pointer-events-none fixed inset-x-0 top-0 z-50 flex items-center justify-between px-6 py-5 transition-all ${
          isTop ? "opacity-0 -translate-y-2" : "opacity-100 translate-y-0"
        }`}
      >
        <div className="pointer-events-auto flex items-center gap-3">
          <div className="relative h-11 w-11 overflow-hidden rounded-xl bg-white/90 shadow-sm">
            <Image src="/logo.png" alt="Logo AREA" fill className="object-contain p-1.5" priority />
          </div>
          <span className="text-base font-semibold tracking-tight text-black drop-shadow-md">AREA</span>
        </div>
        <div className="pointer-events-auto flex items-center gap-3">
          <Link href="/register">
            <Button className="border-0 bg-[var(--card-color-3)] px-5 py-2 text-white shadow-[0_14px_30px_rgba(2,4,15,0.45)] transition-transform duration-300 hover:-translate-y-1 hover:opacity-95">
              Signup
            </Button>
          </Link>
          <Link href="/login">
            <Button variant="ghost" className="px-4 py-2 transition-transform duration-300 hover:-translate-y-1">
              Login
            </Button>
          </Link>
        </div>
      </header>

      <section className="relative isolate flex min-h-screen flex-col justify-center overflow-hidden">
        <div
          className="absolute inset-0 -z-20 scale-[1.02] bg-cover bg-center bg-no-repeat blur-[4px]"
          style={{ backgroundImage: "var(--hero-bg)" }}
        />
        <div className="absolute inset-0 -z-10 bg-gradient-to-b from-black/32 via-[#0b1b2f]/48 to-[var(--background)]/85" />

        <div className="relative z-10 mx-auto flex max-w-5xl flex-col items-center gap-10 px-6 pt-20 text-center text-white">
          <div className="reveal relative h-20 w-20 overflow-hidden rounded-2xl bg-white/90 shadow-[0_16px_60px_rgba(0,0,0,0.28)]">
            <Image src="/logo.png" alt="Logo AREA" fill className="object-contain p-3" />
          </div>

          <div className="reveal space-y-4">
            <h1 className="text-4xl font-semibold leading-tight sm:text-5xl md:text-[3.6rem]">
              Connectez vos applications,{" "}
              <span className="bg-gradient-to-r from-[var(--card-color-1)] via-[var(--card-color-3)] to-[var(--card-color-2)] bg-clip-text text-transparent">
                simplifiez votre vie
              </span>
            </h1>
            <p className="mx-auto max-w-3xl text-lg leading-relaxed text-white/85 sm:text-xl">
              AREA orchestre vos automatisations puissantes entre vos services favoris. Gagnez du temps,
              gardez le contrôle, et concentrez-vous sur ce qui compte vraiment.
            </p>
          </div>

          <div className="reveal flex flex-wrap items-center justify-center gap-4">
            <Link href="/register">
              <Button className="border-0 bg-white/95 px-7 py-3 text-[var(--card-color-1)] shadow-[0_18px_40px_rgba(15,23,42,0.22)] transition-all duration-300 hover:-translate-y-1 hover:opacity-95">
                S'inscrire
              </Button>
            </Link>
            <Link href="/login">
              <Button
                variant="secondary"
                className="bg-white/15 px-7 py-3 text-white backdrop-blur transition-all duration-300 hover:-translate-y-1 hover:bg-white/25"
              >
                Connexion
              </Button>
            </Link>
          </div>

          <a
            href="#content-start"
            className="reveal mt-6 flex h-12 w-12 items-center justify-center rounded-full border border-white/30 text-white/85 transition hover:-translate-y-1 hover:border-white/60"
            aria-label="Découvrir"
          >
            <svg aria-hidden="true" viewBox="0 0 24 24" className="h-6 w-6" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="m6 9 6 6 6-6" strokeLinecap="round" strokeLinejoin="round" />
            </svg>
          </a>
        </div>
      </section>

      <div id="content-start" className="relative" style={gridPattern}>
        <div className="pointer-events-none absolute inset-0 -z-10">
          <div className="absolute -left-32 top-0 h-96 w-96 rounded-full bg-[var(--card-color-5)] blur-[140px] opacity-65" />
          <div className="absolute right-6 top-10 h-[28rem] w-[28rem] rounded-full bg-[var(--card-color-2)] blur-[160px] opacity-55" />
          <div className="absolute left-1/2 top-1/2 h-[34rem] w-[34rem] -translate-x-1/2 -translate-y-1/2 rounded-full bg-[var(--card-color-1)]/30 blur-[160px]" />
          <div className="absolute inset-0 bg-[radial-gradient(circle_at_20%_20%,rgba(229,218,218,0.08),transparent_35%)]" />
        </div>

        <main className="mx-auto flex max-w-6xl flex-col gap-28 px-6 pb-32 pt-24">
          <section
            id="integrations"
            className="reveal flex flex-col items-center gap-5 rounded-3xl border border-[var(--surface-border)] bg-[var(--surface)]/85 px-6 py-8 shadow-[0_18px_50px_rgba(15,23,42,0.08)] backdrop-blur"
          >
          <span className="text-[11px] font-semibold uppercase tracking-[0.28em] text-[var(--muted)]">
            Connecteurs disponibles
          </span>
          <div className="flex flex-wrap items-center justify-center gap-3">
            {services.map((service) => (
              <div
                key={service.name}
                className="group flex h-12 w-12 items-center justify-center rounded-2xl text-base font-semibold text-white shadow-sm transition duration-300 hover:-translate-y-1 hover:shadow-[0_12px_30px_rgba(0,0,0,0.12)]"
                style={{
                  backgroundColor: service.color,
                }}
                aria-label={service.name}
              >
                <span className="transition duration-300 group-hover:scale-110">{service.emoji}</span>
              </div>
            ))}
          </div>
        </section>

        <section className="reveal grid w-full gap-6 md:grid-cols-3">
          {stats.map((stat) => (
            <div
              key={stat.label}
              className="flex items-center gap-4 rounded-2xl border border-[var(--surface-border)] bg-[var(--surface)]/95 p-5 shadow-[0_20px_40px_rgba(15,23,42,0.08)] backdrop-blur transition duration-300 hover:-translate-y-1 hover:shadow-[0_26px_60px_rgba(15,23,42,0.12)]"
            >
              <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-[var(--card-color-3)] text-white shadow-sm">
                📊
              </div>
              <div className="text-left">
                <div className="text-2xl font-semibold">{stat.value}</div>
                <div className="text-sm text-[var(--muted)]">{stat.label}</div>
              </div>
            </div>
          ))}
        </section>

        <section id="features" className="reveal space-y-8">
          <div className="space-y-3 text-center">
            <h2 className="text-3xl font-semibold sm:text-4xl">Pensée pour les équipes qui livrent</h2>
            <p className="mx-auto max-w-2xl text-lg text-[var(--muted)]">
              Inspirée d'IFTTT/Zapier et développée à Epitech : actions, réactions, hooks monitorés,
              authentification centralisée et supervision prête à l'emploi.
            </p>
          </div>

          <div className="grid gap-5 md:grid-cols-4">
            {features.map((item) => (
              <div
                key={item.title}
                className="flex flex-col gap-3 rounded-2xl border border-[var(--surface-border)] bg-[var(--surface)]/85 p-5 text-left shadow-[0_18px_50px_rgba(15,23,42,0.08)] backdrop-blur transition duration-300 hover:-translate-y-1 hover:shadow-[0_26px_70px_rgba(15,23,42,0.12)]"
              >
                <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-[var(--card-color-2)] text-white shadow-sm">
                  {item.icon}
                </div>
                <h3 className="text-lg font-semibold">{item.title}</h3>
                <p className="text-sm leading-relaxed text-[var(--muted)]">{item.body}</p>
              </div>
            ))}
          </div>
        </section>

        <section className="reveal space-y-8">
          <div className="space-y-3 text-center">
            <h2 className="text-3xl font-semibold sm:text-4xl">Comment ça marche ?</h2>
            <p className="mx-auto max-w-2xl text-lg text-[var(--muted)]">
              3 étapes : connecter vos comptes, choisir le déclencheur et la réaction, laisser AREA
              surveiller et exécuter avec traçabilité.
            </p>
          </div>

          <div className="grid gap-5 md:grid-cols-3">
            {steps.map((step) => (
              <div
                key={step.id}
                className="flex flex-col gap-3 rounded-2xl border border-[var(--surface-border)] bg-[var(--surface)]/85 p-6 shadow-[0_18px_50px_rgba(15,23,42,0.08)] backdrop-blur transition duration-300 hover:-translate-y-1 hover:shadow-[0_26px_70px_rgba(15,23,42,0.12)]"
              >
                <span
                  className="inline-flex h-10 w-10 items-center justify-center rounded-2xl text-base font-semibold text-white shadow-sm"
                  style={{ backgroundColor: step.color }}
                >
                  {step.id}
                </span>
                <h3 className="text-lg font-semibold">{step.title}</h3>
                <p className="text-sm leading-relaxed text-[var(--muted)]">{step.body}</p>
              </div>
            ))}
          </div>
        </section>

        <section className="reveal overflow-hidden rounded-3xl border border-[var(--surface-border)] bg-[var(--surface)]/90 p-12 text-center shadow-[0_24px_80px_rgba(15,23,42,0.14)] backdrop-blur">
          <div className="mx-auto flex max-w-3xl flex-col items-center gap-6">
            <h3 className="text-3xl font-semibold sm:text-4xl text-[var(--foreground)]">
              Prêt à lancer vos automations ?
            </h3>
            <p className="text-lg text-[var(--muted)]">
              Créez votre première AREA, branchez vos services, laissez le moteur d'exécution et de
              supervision travailler pour vous.
            </p>
            <div className="flex flex-wrap items-center justify-center gap-4">
              <Link href="/register">
                <Button className="border-0 bg-[var(--card-color-3)] px-7 py-3 text-[var(--background)] shadow-[0_18px_50px_rgba(15,23,42,0.18)] transition-transform duration-300 hover:-translate-y-1 hover:opacity-90">
                  Signup
                  <svg
                    aria-hidden="true"
                    viewBox="0 0 24 24"
                    className="h-5 w-5"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  >
                    <path d="M5 12h14" />
                    <path d="m13 6 6 6-6 6" />
                  </svg>
                </Button>
              </Link>
              <Link href="/login">
                <Button
                  variant="secondary"
                  className="border border-[var(--surface-border)] bg-transparent px-7 py-3 text-[var(--foreground)] transition-transform duration-300 hover:-translate-y-1 hover:bg-[var(--surface)]/70 hover:text-[var(--foreground)]"
                >
                  Login
                </Button>
              </Link>
            </div>
            <div className="flex flex-wrap items-center justify-center gap-4 text-sm text-[var(--muted)]">
              <span className="flex items-center gap-2">
                <span className="h-2 w-2 rounded-full bg-[var(--card-color-3)]" aria-hidden="true" />
                OAuth2 sécurisé
              </span>
              <span className="flex items-center gap-2">
                <span className="h-2 w-2 rounded-full bg-[var(--card-color-2)]" aria-hidden="true" />
                Hooks monitorés
              </span>
              <span className="flex items-center gap-2">
                <span className="h-2 w-2 rounded-full bg-[var(--card-color-1)]" aria-hidden="true" />
                Logs exploitables
              </span>
            </div>
          </div>
        </section>
      </main>

      </div>
    </div>
  );
}
