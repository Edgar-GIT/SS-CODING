import { Flame, Sparkles, Star, Trophy } from "lucide-react";

import { Navbar } from "@/components/Navbar";
import { leaderboardRows, podium } from "@/data/ranking";

function formatNumber(value: number) {
  return new Intl.NumberFormat("en-US").format(value);
}

function getPodiumLayout(rank: number) {
  if (rank === 1) {
    return {
      wrapper: "order-1 md:order-2 md:-mt-6",
      visual: "h-48",
      platform: "h-44",
      platformWidth: "max-w-[320px]",
      avatar: "h-[88px] w-[88px] text-2xl",
      trophy: "h-28 w-28",
      rankText: "text-6xl",
    };
  }

  if (rank === 2) {
    return {
      wrapper: "order-2 md:order-1 md:mt-24",
      visual: "h-32",
      platform: "h-36",
      platformWidth: "max-w-[320px]",
      avatar: "h-16 w-16 text-base",
      trophy: "h-14 w-14",
      rankText: "text-5xl",
    };
  }

  return {
    wrapper: "order-3 md:mt-28",
    visual: "h-32",
    platform: "h-32",
    platformWidth: "max-w-[320px]",
    avatar: "h-16 w-16 text-base",
    trophy: "h-14 w-14",
    rankText: "text-5xl",
  };
}

export function RankingPage() {
  return (
    <div className="min-h-screen bg-background">
      <Navbar />

      <main className="relative min-h-[calc(100vh-3.5rem)] overflow-hidden">
        <div className="absolute inset-0 starfield opacity-45 animate-twinkle pointer-events-none" />
        <div
          className="absolute inset-0 pointer-events-none"
          style={{
            background:
              "radial-gradient(circle at 32% 0%, oklch(0.42 0.21 285 / 0.34), transparent 34%), radial-gradient(circle at 55% 44%, oklch(0.39 0.2 80 / 0.14), transparent 22%), linear-gradient(180deg, oklch(0.13 0.05 270), oklch(0.1 0.04 260))",
          }}
        />

        <section className="relative mx-auto max-w-7xl px-6 pt-7 pb-12 md:pt-10 md:pb-16">
          <div className="max-w-4xl">
            <h1 className="font-display text-5xl font-bold leading-none tracking-normal text-foreground md:text-6xl lg:text-7xl">
              The <span className="text-[oklch(0.64_0.24_280)]">galactic</span>{" "}
              <span className="text-[oklch(0.72_0.2_235)]">leaderboard</span>.
            </h1>

            <p className="mt-6 max-w-xl text-lg leading-8 text-muted-foreground">
              Solve exercises, ship projects and win quizzes to earn XP. The top of the galaxy is
              reserved for the relentless.
            </p>
          </div>

          <div className="relative mx-auto mt-14 max-w-6xl">
            <div className="absolute left-1/2 top-10 hidden h-72 w-72 -translate-x-1/2 rounded-full bg-yellow-400/20 blur-3xl md:block" />
            <div className="grid gap-8 md:grid-cols-3 md:items-end">
              {podium.map((entry) => {
                const layout = getPodiumLayout(entry.rank);
                const isChampion = entry.rank === 1;

                return (
                  <article
                    key={entry.rank}
                    className={`relative flex flex-col items-center ${layout.wrapper}`}
                  >
                    <div className={`relative mb-2 flex ${layout.visual} items-end justify-center`}>
                      {isChampion ? (
                        <>
                          <Sparkles className="absolute -left-8 top-2 h-5 w-5 text-yellow-300/70" />
                          <Sparkles className="absolute -right-8 top-5 h-5 w-5 text-yellow-300/70" />
                        </>
                      ) : null}
                      <Trophy className={`${layout.trophy} ${entry.trophyClassName}`} />
                    </div>

                    <div className="relative z-10">
                      <div
                        className={`flex ${layout.avatar} items-center justify-center rounded-full bg-gradient-to-br ${entry.accentClassName} font-display font-bold shadow-[0_0_50px_oklch(0.69_0.19_85_/_0.35)] ring-2 ring-background`}
                      >
                        {entry.initials}
                      </div>
                      <span className="absolute -bottom-1 -right-1 flex h-6 w-6 items-center justify-center rounded-full bg-gradient-to-br from-yellow-300 to-orange-500 text-xs font-bold text-slate-950 ring-2 ring-background">
                        {entry.rank}
                      </span>
                    </div>

                    <div
                      className={`mt-3 flex w-full ${layout.platformWidth} flex-col items-center rounded-t-lg border border-border/60 bg-card/70 px-7 py-7 text-center shadow-[0_24px_70px_oklch(0.06_0.03_260_/_0.28)] backdrop-blur ${layout.platform}`}
                    >
                      <div
                        className={`font-display ${layout.rankText} font-bold leading-none text-[oklch(0.72_0.2_235)]`}
                      >
                        #{entry.rank}
                      </div>
                      <h2 className="mt-4 text-lg font-bold text-foreground">{entry.username}</h2>
                      <p className="mt-2 text-sm text-muted-foreground">{entry.title}</p>
                      <div className="mt-5 flex items-center justify-center gap-6 text-xs font-semibold uppercase text-accent">
                        <span className="inline-flex items-center gap-1">
                          <Flame className="h-3.5 w-3.5" />
                          {formatNumber(entry.xp)} XP
                        </span>
                        <span className="inline-flex items-center gap-1">
                          <Star className="h-3.5 w-3.5" />
                          {formatNumber(entry.solved)} solved
                        </span>
                      </div>
                    </div>
                  </article>
                );
              })}
            </div>
          </div>

          <div className="mt-16 overflow-hidden rounded-lg border border-border/60 bg-card/70 backdrop-blur">
            <div className="grid grid-cols-[80px_1.5fr_1fr_1fr] border-b border-border/50 px-5 py-5 text-xs font-semibold uppercase text-muted-foreground md:grid-cols-[120px_2fr_1fr_1fr]">
              <span>Rank</span>
              <span>Coder</span>
              <span>Tier</span>
              <span className="text-right">XP</span>
            </div>

            {leaderboardRows.map((row) => (
              <div
                key={row.rank}
                className="grid grid-cols-[80px_1.5fr_1fr_1fr] items-center border-b border-border/30 px-5 py-5 text-base last:border-b-0 md:grid-cols-[120px_2fr_1fr_1fr]"
              >
                <span className="text-muted-foreground">#{row.rank}</span>
                <span className="flex min-w-0 items-center gap-3 font-semibold text-foreground">
                  <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-gradient-to-br from-violet-500 to-sky-500 text-slate-950">
                    <Star className="h-3.5 w-3.5" />
                  </span>
                  <span className="truncate">{row.username}</span>
                </span>
                <span className="text-muted-foreground">{row.tier}</span>
                <span className="text-right font-bold text-foreground">{formatNumber(row.xp)}</span>
              </div>
            ))}
          </div>
        </section>
      </main>
    </div>
  );
}
