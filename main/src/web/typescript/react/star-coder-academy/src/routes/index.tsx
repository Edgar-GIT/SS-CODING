import { createFileRoute, Link } from "@tanstack/react-router";
import { Navbar } from "@/components/Navbar";
import { SlotCounter } from "@/components/SlotCounter";
import homeLogo from "@resources/img/logo/home_logo.png";
import { Code2, Trophy, Zap, BookOpen, Rocket, Target } from "lucide-react";

export const Route = createFileRoute("/")({
  head: () => ({
    meta: [
      { title: "SS Coding — Become an Elite Programmer" },
      {
        name: "description",
        content:
          "Master any programming language through interactive guides, exercises, quizzes, projects, and global rankings. Become an elite developer.",
      },
      { property: "og:title", content: "SS Coding — Become an Elite Programmer" },
      {
        property: "og:description",
        content:
          "Master any programming language through interactive guides, exercises, quizzes, projects, and global rankings.",
      },
    ],
  }),
  component: Home,
});

const features = [
  {
    icon: BookOpen,
    title: "Deep Guides",
    desc: "From first variable to advanced architecture — pick a language, master every concept.",
  },
  {
    icon: Code2,
    title: "Live Exercises",
    desc: "Thousands of hands-on challenges across every paradigm. Code, run, level up.",
  },
  {
    icon: Zap,
    title: "Quizzes",
    desc: "Rapid-fire knowledge checks that lock concepts into long-term memory.",
  },
  {
    icon: Rocket,
    title: "Real Projects",
    desc: "Ship portfolio-grade builds with structured roadmaps and reviewable milestones.",
  },
  {
    icon: Trophy,
    title: "Global Ranking",
    desc: "Climb the leaderboard. Compete with elite coders across the galaxy.",
  },
  {
    icon: Target,
    title: "Skill Trees",
    desc: "Personalized paths that adapt to your pace, goals, and chosen languages.",
  },
];

const stats: { value: string; label: string; live?: boolean }[] = [
  { value: "40+", label: "Languages" },
  { value: "5,000+", label: "Exercises" },
  { value: "", label: "Coders", live: true },
  { value: "∞", label: "Possibilities" },
];

function Home() {
  return (
    <div className="min-h-screen">
      <Navbar />

      {/* HERO */}
      <section className="relative overflow-hidden">
        <div className="absolute inset-0 starfield opacity-40 animate-twinkle pointer-events-none" />
        <div className="absolute inset-0 pointer-events-none"
             style={{ background: "var(--gradient-nebula)" }} />

        <div className="relative mx-auto max-w-7xl px-6 pt-6 pb-16 lg:pt-8 lg:pb-24 grid lg:grid-cols-2 gap-8 lg:gap-10 items-start">
          <div className="space-y-8">
            <div className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full border border-border/60 bg-secondary/40 backdrop-blur text-xs font-medium text-muted-foreground">
              <span className="w-2 h-2 rounded-full bg-accent animate-pulse-glow" />
              Launch your coding journey
            </div>

            <h1 className="font-display text-5xl md:text-6xl lg:text-7xl font-bold leading-[1.05] tracking-tight">
              Become an{" "}
              <span className="text-gradient-cosmic">elite programmer</span>
              <br />
              capable of building{" "}
              <span className="relative">
                anything
                <span className="absolute -bottom-2 left-0 right-0 h-1 bg-cosmic rounded-full opacity-70" />
              </span>
              .
            </h1>

            <p className="text-lg md:text-xl text-muted-foreground max-w-xl leading-relaxed">
              Master any language, conquer every concept. Interactive guides,
              real exercises, quizzes, projects and a global ranking — all in one
              cosmic learning platform.
            </p>

            <div className="flex flex-wrap gap-4">
              <Link
                to="/guides"
                className="inline-flex items-center gap-2 px-7 py-3.5 rounded-xl bg-cosmic text-primary-foreground font-semibold shadow-cosmic hover:scale-105 transition-transform"
              >
                Start Learning
                <Rocket className="w-4 h-4" />
              </Link>
              <Link
                to="/exercises"
                className="inline-flex items-center gap-2 px-7 py-3.5 rounded-xl border border-border bg-secondary/40 backdrop-blur font-semibold text-foreground hover:bg-secondary/70 transition-colors"
              >
                Try an Exercise
              </Link>
            </div>

            <div className="grid grid-cols-4 gap-4 pt-6 max-w-lg">
              {stats.map((s) => (
                <div key={s.label}>
                  {s.live ? (
                    <SlotCounter />
                  ) : (
                    <div className="font-display text-2xl md:text-3xl font-bold text-gradient-cosmic">
                      {s.value}
                    </div>
                  )}
                  <div className="text-xs uppercase tracking-wider text-muted-foreground mt-1">
                    {s.label}
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* HERO IMAGE */}
          <div className="relative flex justify-center lg:justify-start items-start -mt-2 lg:-mt-6 lg:-ml-8 xl:-ml-12">
            <div className="absolute top-1/4 left-1/4 w-3/4 h-3/4 bg-cosmic opacity-30 blur-3xl rounded-full animate-pulse-glow" />
            <img
              src={homeLogo}
              alt="SS Coding"
              className="relative w-full max-w-[480px] lg:max-w-[520px] animate-float-slow drop-shadow-[0_0_60px_oklch(0.55_0.25_285_/_0.6)] translate-x-[-4%] lg:translate-x-[-8%] -translate-y-2 lg:-translate-y-6"
            />
          </div>
        </div>
      </section>

      {/* FEATURES */}
      <section className="relative py-24 border-t border-border/40">
        <div className="mx-auto max-w-7xl px-6">
          <div className="max-w-2xl mb-16">
            <div className="text-sm font-semibold text-accent uppercase tracking-widest mb-3">
              The Arsenal
            </div>
            <h2 className="font-display text-4xl md:text-5xl font-bold leading-tight">
              Everything you need to{" "}
              <span className="text-gradient-cosmic">level up</span>.
            </h2>
            <p className="text-muted-foreground mt-4 text-lg">
              A complete ecosystem engineered to turn curious beginners into
              senior engineers.
            </p>
          </div>

          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
            {features.map((f) => (
              <div
                key={f.title}
                className="group relative p-7 rounded-2xl bg-card/60 backdrop-blur border border-border/60 hover:border-primary/60 transition-all hover:-translate-y-1 hover:shadow-cosmic overflow-hidden"
              >
                <div className="absolute inset-0 opacity-0 group-hover:opacity-100 transition-opacity"
                     style={{ background: "var(--gradient-nebula)" }} />
                <div className="relative">
                  <div className="w-12 h-12 rounded-xl bg-cosmic flex items-center justify-center mb-5 shadow-glow">
                    <f.icon className="w-6 h-6 text-primary-foreground" />
                  </div>
                  <h3 className="font-display text-xl font-bold mb-2">
                    {f.title}
                  </h3>
                  <p className="text-muted-foreground leading-relaxed">{f.desc}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA */}
      <section className="relative py-24">
        <div className="mx-auto max-w-5xl px-6">
          <div className="relative rounded-3xl overflow-hidden border border-border/60 bg-card/60 backdrop-blur p-12 md:p-16 text-center">
            <div className="absolute inset-0 starfield opacity-30 animate-twinkle" />
            <div className="absolute inset-0"
                 style={{ background: "var(--gradient-nebula)" }} />
            <div className="relative">
              <h2 className="font-display text-4xl md:text-5xl font-bold mb-5">
                Your launchpad is{" "}
                <span className="text-gradient-cosmic">ready</span>.
              </h2>
              <p className="text-muted-foreground text-lg max-w-xl mx-auto mb-8">
                Join thousands of coders building the future — one commit at a time.
              </p>
              <Link
                to="/guides"
                className="inline-flex items-center gap-2 px-8 py-4 rounded-xl bg-cosmic text-primary-foreground font-semibold shadow-cosmic hover:scale-105 transition-transform"
              >
                Begin Your Mission
                <Rocket className="w-5 h-5" />
              </Link>
            </div>
          </div>
        </div>
      </section>

      <footer className="border-t border-border/40 py-8">
        <div className="mx-auto max-w-7xl px-6 text-center text-sm text-muted-foreground">
          © {new Date().getFullYear()} SS Coding — Forge your path among the stars.
        </div>
      </footer>
    </div>
  );
}
