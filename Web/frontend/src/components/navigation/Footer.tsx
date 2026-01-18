import Image from "next/image";

const legalLinks = [
  "Mentions légales",
  "CGU",
  "Politique de confidentialité",
  "Cookies",
  "Contact",
];

export function Footer() {
  return (
    <footer className="border-t border-[var(--surface-border)] bg-[var(--background)]/80">
      <div className="mx-auto flex w-full max-w-6xl flex-col gap-6 px-6 py-8 md:flex-row md:items-center md:justify-between">
        <div className="flex items-center gap-3">
          <div className="relative h-10 w-10 overflow-hidden rounded-xl bg-white/90 shadow-sm">
            <Image src="/logo.png" alt="Logo AREA" fill className="object-contain p-1" />
          </div>
          <div className="space-y-1">
            <span className="text-lg font-semibold tracking-tight">AREA</span>
            <p className="text-xs text-[var(--muted)]">
              Automatisation open-source portée par l&apos;équipe Epitech.
            </p>
          </div>
        </div>
        <div className="flex flex-wrap gap-x-4 gap-y-2 text-xs text-[var(--muted)]">
          {legalLinks.map((label) => (
            <span
              key={label}
              className="rounded-full border border-[var(--surface-border)] bg-[var(--surface)] px-3 py-1.5 text-[11px] font-semibold uppercase tracking-wide text-[var(--muted)]"
            >
              {label}
            </span>
          ))}
        </div>
      </div>
      <div className="border-t border-[var(--surface-border)]/80 py-4">
        <div className="mx-auto flex w-full max-w-6xl items-center justify-between px-6 text-xs text-[var(--muted)]">
          <span>© 2025 AREA</span>
          <span>Tous droits réservés.</span>
        </div>
      </div>
    </footer>
  );
}
