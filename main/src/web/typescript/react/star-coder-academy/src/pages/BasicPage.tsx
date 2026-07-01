import type { ReactNode } from "react";

import { Navbar } from "@/components/Navbar";

type BasicPageProps = {
  eyebrow: string;
  title: string;
  description: string;
  children?: ReactNode;
};

export function BasicPage({ eyebrow, title, description, children }: BasicPageProps) {
  return (
    <div className="min-h-screen bg-background">
      <Navbar />

      <main className="relative min-h-[calc(100vh-3.5rem)] overflow-hidden">
        <div className="absolute inset-0 starfield opacity-35 animate-twinkle pointer-events-none" />
        <div className="absolute inset-0 pointer-events-none bg-[radial-gradient(circle_at_28%_10%,oklch(0.42_0.2_285_/_0.28),transparent_34%),linear-gradient(180deg,oklch(0.13_0.05_270),oklch(0.1_0.04_260))]" />

        <section className="relative mx-auto flex max-w-7xl flex-col justify-center px-6 py-24 md:min-h-[calc(100vh-3.5rem)]">
          <div className="max-w-3xl">
            <p className="text-sm font-semibold uppercase text-accent">{eyebrow}</p>
            <h1 className="mt-4 font-display text-5xl font-bold leading-tight tracking-normal text-foreground md:text-6xl">
              {title}
            </h1>
            <p className="mt-6 max-w-2xl text-lg leading-8 text-muted-foreground">{description}</p>
          </div>

          {children ? <div className="mt-10">{children}</div> : null}
        </section>
      </main>
    </div>
  );
}
