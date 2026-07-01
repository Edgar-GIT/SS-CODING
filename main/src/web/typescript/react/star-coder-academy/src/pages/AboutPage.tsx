import type { ReactNode } from "react";

import { Navbar } from "@/components/Navbar";
import { feats, values } from "@/data/about";

function AboutCard({ children, className = "" }: { children: ReactNode; className?: string }) {
  return (
    <article
      className={`rounded-lg border border-border/60 bg-card/70 p-7 shadow-[0_22px_60px_oklch(0.06_0.03_260/0.22)] backdrop-blur transition-all duration-300 hover:border-accent/80 hover:bg-card/90 hover:shadow-[0_28px_80px_oklch(0.55_0.25_285/0.35)] hover:scale-105 md:p-8 ${className}`}
    >
      {children}
    </article>
  );
}

export function AboutPage() {
  return (
    <div className="min-h-screen bg-background">
      <Navbar />

      <main className="relative min-h-[calc(100vh-3.5rem)] overflow-hidden">
        <div className="absolute inset-0 starfield opacity-45 animate-twinkle pointer-events-none" />
        <div
          className="absolute inset-0 pointer-events-none"
          style={{
            background:
              "radial-gradient(circle at 22% 0%, oklch(0.42 0.21 285 / 0.34), transparent 35%), radial-gradient(circle at 70% 50%, oklch(0.27 0.16 235 / 0.18), transparent 38%), linear-gradient(180deg, oklch(0.13 0.05 270), oklch(0.1 0.04 260))",
          }}
        />

        <section className="relative mx-auto max-w-7xl px-6 pt-5 pb-10 md:pt-8 md:pb-14">
          <div className="max-w-5xl">
            <div className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full border border-border/60 bg-secondary/40 backdrop-blur text-xs font-medium text-muted-foreground">
              <span className="w-2 h-2 rounded-full bg-accent animate-pulse-glow" />
              About the creator
            </div>

            <h1 className="mt-6 font-display text-4xl font-bold leading-tight tracking-normal text-foreground md:text-5xl lg:text-6xl">
              One dev. <span className="text-[oklch(0.64_0.24_280)]">Decades of obsession.</span>
              <br />
              Now yours to inherit.
            </h1>

            <p className="mt-5 max-w-3xl text-base leading-7 text-muted-foreground md:text-lg md:leading-8">
              SS Coding is a solo project — built, taught and maintained by a single programmer who
              has spent a lifetime mastering every corner of computing, and now wants to hand that
              journey to you.
            </p>
          </div>

          <article className="relative mt-8 overflow-hidden rounded-lg border border-border/60 bg-card/70 p-8 shadow-[0_22px_60px_oklch(0.06_0.03_260/0.22)] backdrop-blur md:p-10">
            <div className="absolute inset-0 starfield opacity-30 animate-twinkle" />
            <div className="absolute inset-0 bg-[radial-gradient(circle_at_30%_0%,oklch(0.42_0.21_285/0.24),transparent_42%)]" />

            <div className="relative max-w-5xl space-y-5 text-base leading-7 text-muted-foreground md:text-lg md:leading-8">
              <h2 className="font-display text-2xl font-bold text-foreground md:text-3xl">
                My story.
              </h2>
              <p>
                I wrote my first line of C at{" "}
                <span className="font-semibold text-foreground">11 years old</span> — and from that
                moment, I never stopped. What started as curiosity turned into an obsession, and
                that obsession turned into mastery across every layer of computing and electronics I
                could get my hands on.
              </p>
              <p>
                Today I&apos;m a solo developer, a world-class competitive programmer and a
                professional{" "}
                <span className="font-semibold text-foreground">pentester / gray-hat hacker</span>.
                I&apos;ve won international competitions — including one hosted by{" "}
                <span className="font-semibold text-foreground">NASA</span> — collected trophies
                across algorithms, CTFs and hardware challenges, and shipped things most developers
                only read about: my own CPU and GPU designs, my own operating system, my own
                programming languages, my own PC firmware and BIOS.
              </p>
              <p>
                I know informatics and electronics end-to-end — silicon to shader, bootloader to
                browser, packet to protocol — and I&apos;m fluent in{" "}
                <span className="font-semibold text-foreground">
                  more than 50 programming languages
                </span>
                .
              </p>
              <p>
                But the truth is: I&apos;ve always loved{" "}
                <span className="font-semibold text-foreground">helping people</span>. Watching
                someone go from confused to capable is more rewarding than any trophy. I built SS
                Coding so anyone — regardless of age, background or budget — can walk the same path
                of excellence I did. No gatekeeping. No filler. Just the real journey, from first
                line of code to godlike command over the machine.
              </p>
              <p className="font-semibold text-foreground">
                Welcome aboard. The galaxy of code is yours to explore — and I&apos;ll be your
                guide.
              </p>
            </div>
          </article>

          <div className="mt-14">
            <h2 className="font-display text-3xl font-bold tracking-normal text-foreground md:text-4xl">
              A few <span className="text-[oklch(0.72_0.2_235)]">feats</span>.
            </h2>

            <div className="mt-8 grid gap-6 md:grid-cols-2">
              {feats.map((feat) => (
                <AboutCard key={feat.title} className="min-h-52">
                  <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-linear-to-br from-violet-500 to-sky-500 shadow-[0_14px_30px_oklch(0.12_0.04_260/0.32)]">
                    <feat.icon className="h-6 w-6 text-[oklch(0.1_0.03_260)]" />
                  </div>
                  <h3 className="mt-6 font-display text-xl font-bold text-foreground">
                    {feat.title}
                  </h3>
                  <p className="mt-3 text-base leading-7 text-muted-foreground">
                    {feat.description}
                  </p>
                </AboutCard>
              ))}
            </div>
          </div>

          <div className="mt-14">
            <h2 className="font-display text-3xl font-bold tracking-normal text-foreground md:text-4xl">
              What drives <span className="text-[oklch(0.72_0.2_235)]">this platform</span>.
            </h2>

            <div className="mt-8 grid gap-6 sm:grid-cols-2 xl:grid-cols-4">
              {values.map((value) => (
                <AboutCard key={value.title} className="min-h-48">
                  <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-linear-to-br from-violet-500 to-sky-500 shadow-[0_14px_30px_oklch(0.12_0.04_260/0.32)]">
                    <value.icon className="h-6 w-6 text-[oklch(0.1_0.03_260)]" />
                  </div>
                  <h3 className="mt-6 font-display text-xl font-bold text-foreground">
                    {value.title}
                  </h3>
                  <p className="mt-3 text-sm leading-7 text-muted-foreground">
                    {value.description}
                  </p>
                </AboutCard>
              ))}
            </div>
          </div>
        </section>
      </main>
    </div>
  );
}
