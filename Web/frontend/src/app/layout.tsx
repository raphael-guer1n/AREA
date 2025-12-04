import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";

import { ColorblindToggle } from "@/components/ui/ColorblindToggle";
import "./globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "AREA",
  description: "Automate and orchestrate actions across services.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="fr">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased bg-[var(--background)] text-[var(--foreground)]`}
      >
        <ColorblindToggle />
        {children}
      </body>
    </html>
  );
}
