import Link from "next/link";

const pageStyle = {
    backgroundColor: "var(--background)",
    color: "var(--foreground)",
};

const outlineButtonStyle = {
    borderColor: "var(--foreground)",
    color: "var(--foreground)",
};

const primaryButtonStyle = {
    backgroundColor: "var(--accent)",
    color: "var(--background)",
};

export default function Home() {
    return (
        <div className="flex min-h-screen flex-col text-sm" style={pageStyle}>
            <header className="flex items-center justify-between px-10 py-6 text-base">
                <span className="text-3xl font-semibold tracking-[0.4em]">AREA</span>
                <Link
                    href="/Login"
                    className="rounded-full border px-6 py-2 font-medium uppercase tracking-wide transition hover:bg-[var(--foreground)] hover:text-[var(--background)]"
                    style={outlineButtonStyle}>
                    Login
                </Link>
            </header>
            <main className="flex flex-1 items-center justify-center px-6 py-12">
                <button
                    className="rounded-full px-12 py-5 text-lg font-semibold uppercase tracking-wide shadow-2xl transition hover:-translate-y-0.5 hover:opacity-90 hover:shadow-[0_20px_45px_rgba(0,0,0,0.35)]"
                    style={primaryButtonStyle}>
                    test backend
                </button>
            </main>
        </div>
    );
}
