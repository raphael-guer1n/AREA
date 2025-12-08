"use client";

import Link from "next/link";
import { useEffect, useState } from "react";

import { Button } from "@/components/ui/Button";

const services = [
  { name: "Mail", emoji: "✉️", from: "#4f46e5", to: "#22d3ee" },
  { name: "Chat", emoji: "💬", from: "#0ea5e9", to: "#22c55e" },
  { name: "Docs", emoji: "📁", from: "#7c3aed", to: "#4f46e5" },
  { name: "Tasks", emoji: "✅", from: "#10b981", to: "#14b8a6" },
  { name: "Social", emoji: "🐦", from: "#2563eb", to: "#60a5fa" },
  { name: "Voice", emoji: "🎙️", from: "#f97316", to: "#fb7185" },
  { name: "Plus", emoji: "+40", from: "#0f172a", to: "#1f2937" },
];

const stats = [
  { value: "50+", label: "Lorem ipsum dolor" },
  { value: "10k+", label: "Sit amet consectetur" },
  { value: "99.9%", label: "Adipiscing elit" },
];

const features = [
  {
    title: "Lorem ipsum dolor",
    body: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer ut magna eget lorem aliquet.",
    icon: "⚡",
  },
  {
    title: "Consectetur amet",
    body: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus posuere velit eu lorem viverra.",
    icon: "🛡️",
  },
  {
    title: "Sed do eiusmod",
    body: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Ut vitae lorem nisl nam lacinia.",
    icon: "🎛️",
  },
  {
    title: "Tempor incididunt",
    body: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Fusce eu lacus porttitor, euismod nisl.",
    icon: "🔗",
  },
];

const steps = [
  {
    id: "01",
    title: "Lorem ipsum",
    body: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer vel libero ut arcu.",
    color: "bg-[#f87171]",
  },
  {
    id: "02",
    title: "Dolor sit amet",
    body: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eu lorem ac augue aliquet.",
    color: "bg-[#34d399]",
  },
  {
    id: "03",
    title: "Consectetur elit",
    body: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed vitae justo viverra, aliquet.",
    color: "bg-[#a855f7]",
  },
];

const gridPattern = {
  backgroundImage:
    "linear-gradient(to right, rgba(15,23,42,0.04) 1px, transparent 1px), linear-gradient(to bottom, rgba(15,23,42,0.04) 1px, transparent 1px)",
  backgroundSize: "72px 72px",
};

export default function HomePage() {
  const [scrollDir, setScrollDir] = useState<"up" | "down" | null>(null);

  useEffect(() => {
    let lastY = window.scrollY;

    const onScroll = () => {
      const current = window.scrollY;
      if (Math.abs(current - lastY) < 2) return;
      setScrollDir(current > lastY ? "down" : "up");
      lastY = current;
    };

    window.addEventListener("scroll", onScroll, { passive: true });
    return () => window.removeEventListener("scroll", onScroll);
  }, []);

  return (
    <div
      className="relative min-h-screen overflow-hidden bg-[var(--background)] text-[var(--foreground)]"
      style={gridPattern}
    >
      <div className="pointer-events-none absolute inset-0 -z-10">
        <div className="absolute -left-32 top-0 h-96 w-96 rounded-full bg-gradient-to-br from-[#22d3ee] to-[#4f46e5] blur-[120px] opacity-60" />
        <div className="absolute right-6 top-10 h-[28rem] w-[28rem] rounded-full bg-gradient-to-br from-[#0ea5e9] to-[#22c55e] blur-[140px] opacity-50" />
        <div className="absolute left-1/2 top-1/2 h-[34rem] w-[34rem] -translate-x-1/2 -translate-y-1/2 rounded-full bg-[var(--surface)]/80 blur-[140px]" />
        <div className="absolute inset-0 bg-[radial-gradient(circle_at_20%_20%,rgba(34,197,94,0.08),transparent_35%),radial-gradient(circle_at_80%_10%,rgba(14,165,233,0.08),transparent_30%)]" />
      </div>

      <header
        className={`sticky top-0 z-50 mx-auto flex max-w-6xl items-center justify-between px-6 py-6 transition-all ${
          scrollDir === "down"
            ? "backdrop-blur-3xl bg-[var(--background)]/92 shadow-[0_12px_30px_rgba(0,0,0,0.06)]"
            : "bg-transparent"
        }`}
      >
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-[var(--foreground)] text-[var(--background)] shadow-sm">
            <svg
              aria-hidden="true"
              viewBox="0 0 24 24"
              className="h-6 w-6"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <path d="M13 2 3 14h7l-1 8 10-12h-7l1-8Z" />
            </svg>
          </div>
          <span className="text-lg font-semibold tracking-tight">AREA</span>
        </div>
        <div className="flex items-center gap-3">
          <Link href="/login">
            <Button variant="ghost" className="px-4 py-2 transition-transform duration-300 hover:-translate-y-1">
              Login
            </Button>
          </Link>
          <Link href="/register">
            <Button className="border-0 bg-gradient-to-r from-[#4f46e5] to-[#22d3ee] px-5 py-2 text-white shadow-[0_14px_30px_rgba(79,70,229,0.35)] transition-transform duration-300 hover:-translate-y-1 hover:opacity-95">
              Signup
            </Button>
          </Link>
        </div>
      </header>

      <main className="mx-auto flex max-w-6xl flex-col gap-20 px-6 pb-24 pt-10">
        <section className="flex flex-col items-center gap-9 text-center">
          <div className="inline-flex items-center gap-3 rounded-full border border-[var(--surface-border)] bg-[var(--surface)]/70 px-4 py-2 text-sm font-semibold shadow-sm backdrop-blur">
            <span className="flex h-8 w-8 items-center justify-center rounded-full bg-[var(--foreground)] text-[var(--background)] text-xs">
              ⚡
            </span>
            <span className="text-[var(--muted)]">Lorem ipsum dolor</span>
          </div>

          <div className="space-y-6">
            <h1 className="text-4xl font-semibold leading-tight sm:text-5xl md:text-[3.25rem]">
              Lorem ipsum dolor sit amet,{" "}
              <span className="bg-gradient-to-r from-[#0ea5e9] via-[#22c55e] to-[#7c3aed] bg-clip-text text-transparent">
                consectetur adipiscing
              </span>
            </h1>
            <p className="mx-auto max-w-3xl text-lg leading-relaxed text-[var(--muted)] sm:text-xl">
              Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur pretium, justo at
              suscipit bibendum, magna nunc tristique ante luctus auctor nisl velit in dolor eget
              ligula elementum feugiat.
            </p>
          </div>

          <div className="flex flex-wrap items-center justify-center gap-4">
            <Link href="/register">
              <Button className="border-0 bg-[var(--foreground)] px-7 py-3 text-[var(--background)] shadow-[0_18px_40px_rgba(15,23,42,0.18)] transition-transform duration-300 hover:-translate-y-1 hover:opacity-90">
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
              <Button variant="secondary" className="px-7 py-3 transition-transform duration-300 hover:-translate-y-1">
                Login
              </Button>
            </Link>
          </div>
        </section>

        <section
          id="integrations"
          className="flex flex-col items-center gap-5 rounded-3xl border border-[var(--surface-border)] bg-[var(--surface)]/85 px-6 py-8 shadow-[0_18px_50px_rgba(15,23,42,0.08)] backdrop-blur"
        >
          <span className="text-[11px] font-semibold uppercase tracking-[0.28em] text-[var(--muted)]">
            Lorem ipsum dolor
          </span>
          <div className="flex flex-wrap items-center justify-center gap-3">
            {services.map((service) => (
              <div
                key={service.name}
                className="group flex h-12 w-12 items-center justify-center rounded-2xl text-base font-semibold text-white shadow-sm transition duration-300 hover:-translate-y-1 hover:shadow-[0_12px_30px_rgba(0,0,0,0.12)]"
                style={{
                  backgroundImage: `linear-gradient(135deg, ${service.from}, ${service.to})`,
                }}
                aria-label={service.name}
              >
                <span className="transition duration-300 group-hover:scale-110">{service.emoji}</span>
              </div>
            ))}
          </div>
        </section>

        <section className="grid w-full gap-6 md:grid-cols-3">
          {stats.map((stat) => (
            <div
              key={stat.label}
              className="flex items-center gap-4 rounded-2xl border border-[var(--surface-border)] bg-[var(--surface)]/95 p-5 shadow-[0_20px_40px_rgba(15,23,42,0.08)] backdrop-blur transition duration-300 hover:-translate-y-1 hover:shadow-[0_26px_60px_rgba(15,23,42,0.12)]"
            >
              <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-[#0ea5e9] via-[#22c55e] to-[#4f46e5] text-white shadow-sm">
                📊
              </div>
              <div className="text-left">
                <div className="text-2xl font-semibold">{stat.value}</div>
                <div className="text-sm text-[var(--muted)]">{stat.label}</div>
              </div>
            </div>
          ))}
        </section>

        <section id="features" className="space-y-8">
          <div className="space-y-3 text-center">
            <h2 className="text-3xl font-semibold sm:text-4xl">Pourquoi Lorem Ipsum ?</h2>
            <p className="mx-auto max-w-2xl text-lg text-[var(--muted)]">
              Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris luctus vitae sapien at
              tempus nulla facilisi sed viverra.
            </p>
          </div>

          <div className="grid gap-5 md:grid-cols-4">
            {features.map((item) => (
              <div
                key={item.title}
                className="flex flex-col gap-3 rounded-2xl border border-[var(--surface-border)] bg-[var(--surface)]/85 p-5 text-left shadow-[0_18px_50px_rgba(15,23,42,0.08)] backdrop-blur transition duration-300 hover:-translate-y-1 hover:shadow-[0_26px_70px_rgba(15,23,42,0.12)]"
              >
                <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-[#4f46e5] via-[#22d3ee] to-[#22c55e] text-white shadow-sm">
                  {item.icon}
                </div>
                <h3 className="text-lg font-semibold">{item.title}</h3>
                <p className="text-sm leading-relaxed text-[var(--muted)]">{item.body}</p>
              </div>
            ))}
          </div>
        </section>

        <section className="space-y-8">
          <div className="space-y-3 text-center">
            <h2 className="text-3xl font-semibold sm:text-4xl">Comment ca marche ?</h2>
            <p className="mx-auto max-w-2xl text-lg text-[var(--muted)]">
              Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras id fringilla nisl nec
              imperdiet lorem ipsum.
            </p>
          </div>

          <div className="grid gap-5 md:grid-cols-3">
            {steps.map((step) => (
              <div
                key={step.id}
                className="flex flex-col gap-3 rounded-2xl border border-[var(--surface-border)] bg-[var(--surface)]/85 p-6 shadow-[0_18px_50px_rgba(15,23,42,0.08)] backdrop-blur transition duration-300 hover:-translate-y-1 hover:shadow-[0_26px_70px_rgba(15,23,42,0.12)]"
              >
                <span className={`${step.color} inline-flex h-10 w-10 items-center justify-center rounded-2xl text-base font-semibold text-white shadow-sm`}>
                  {step.id}
                </span>
                <h3 className="text-lg font-semibold">{step.title}</h3>
                <p className="text-sm leading-relaxed text-[var(--muted)]">{step.body}</p>
              </div>
            ))}
          </div>
        </section>

        <section className="overflow-hidden rounded-3xl bg-gradient-to-r from-[#0f172a] via-[#0b1e33] to-[#0a2f4f] p-12 text-center text-white shadow-[0_24px_80px_rgba(0,0,0,0.35)]">
          <div className="mx-auto flex max-w-3xl flex-col items-center gap-6">
            <h3 className="text-3xl font-semibold sm:text-4xl">Pret a lorem ipsum ?</h3>
            <p className="text-lg text-white/90">
              Lorem ipsum dolor sit amet consectetur adipiscing elit. Rejoignez lorem ipsum qui
              gagnent du temps chaque jour.
            </p>
            <div className="flex flex-wrap items-center justify-center gap-4">
              <Link href="/register">
                <Button className="border-0 bg-white px-7 py-3 text-[var(--blue-primary-1)] shadow-[0_18px_50px_rgba(0,0,0,0.25)] transition-transform duration-300 hover:-translate-y-1 hover:bg-white/90">
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
                  className="border border-white/30 bg-transparent px-7 py-3 text-white transition-transform duration-300 hover:-translate-y-1 hover:bg-white/10 hover:text-white"
                >
                  Login
                </Button>
              </Link>
            </div>
            <div className="flex flex-wrap items-center justify-center gap-4 text-sm text-white/80">
              <span className="flex items-center gap-2">
                <span className="h-2 w-2 rounded-full bg-emerald-400" aria-hidden="true" />
                Lorem ipsum
              </span>
              <span className="flex items-center gap-2">
                <span className="h-2 w-2 rounded-full bg-emerald-400" aria-hidden="true" />
                Dolor sit amet
              </span>
              <span className="flex items-center gap-2">
                <span className="h-2 w-2 rounded-full bg-emerald-400" aria-hidden="true" />
                Consectetur elit
              </span>
            </div>
          </div>
        </section>
      </main>

      <footer className="border-t border-[var(--surface-border)] bg-[var(--background)]/80 py-6">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-6">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-[var(--foreground)] text-[var(--background)] shadow-sm">
              <svg
                aria-hidden="true"
                viewBox="0 0 24 24"
                className="h-6 w-6"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <path d="M13 2 3 14h7l-1 8 10-12h-7l1-8Z" />
              </svg>
            </div>
            <span className="text-lg font-semibold tracking-tight">AREA</span>
          </div>
          <span className="text-sm text-[var(--muted)]">© 2025 AREA. Lorem ipsum dolor sit amet.</span>
        </div>
      </footer>
    </div>
  );
}
