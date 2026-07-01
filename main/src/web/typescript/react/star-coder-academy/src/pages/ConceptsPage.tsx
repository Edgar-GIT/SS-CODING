import { Filter, Search } from "lucide-react";
import { useMemo, useState } from "react";

import { Navbar } from "@/components/Navbar";
import { concepts } from "@/data/concepts";

function normalize(value: string) {
  return value.toLowerCase().trim();
}

export function ConceptsPage() {
  const [query, setQuery] = useState("");
  const normalizedQuery = normalize(query);

  const filteredConcepts = useMemo(() => {
    if (!normalizedQuery) {
      return concepts;
    }

    return concepts.filter((concept) => {
      const searchableValues = [
        concept.title,
        concept.description,
        concept.id,
        ...concept.searchTerms,
      ].map(normalize);
      return searchableValues.some((value) => value.includes(normalizedQuery));
    });
  }, [normalizedQuery]);

  return (
    <div className="min-h-screen bg-background">
      <Navbar />

      <main className="relative min-h-[calc(100vh-3.5rem)] overflow-hidden">
        <div className="absolute inset-0 starfield opacity-45 animate-twinkle pointer-events-none" />
        <div
          className="absolute inset-0 pointer-events-none"
          style={{
            background:
              "radial-gradient(circle at 25% 0%, oklch(0.42 0.21 285 / 0.34), transparent 34%), radial-gradient(circle at 62% 46%, oklch(0.27 0.16 235 / 0.22), transparent 42%), linear-gradient(180deg, oklch(0.13 0.05 270), oklch(0.1 0.04 260))",
          }}
        />

        <section className="relative mx-auto max-w-7xl px-6 pt-7 pb-12 md:pt-10 md:pb-16">
          <div className="max-w-5xl">
            <h1 className="font-display text-5xl font-bold leading-none tracking-normal text-foreground md:text-6xl lg:text-[64px]">
              Concepts that work in{" "}
              <span className="whitespace-nowrap text-[oklch(0.64_0.24_280)]">any</span>
              <br />
              <span className="text-[oklch(0.72_0.2_235)]">language</span>.
            </h1>

            <p className="mt-7 max-w-3xl text-lg leading-8 text-muted-foreground md:text-xl">
              Learn the timeless ideas behind great software. Each concept is paired with examples
              in your favorite language.
            </p>
          </div>

          <div className="mt-12 flex flex-col gap-3 sm:flex-row">
            <div className="min-w-0 flex-1">
              <label className="sr-only" htmlFor="concept-search">
                Search concepts
              </label>
              <div className="relative">
                <Search className="pointer-events-none absolute left-4 top-1/2 h-5 w-5 -translate-y-1/2 text-muted-foreground" />
                <input
                  id="concept-search"
                  type="search"
                  value={query}
                  onChange={(event) => setQuery(event.target.value)}
                  placeholder="Search concepts by name, language, topic..."
                  className="h-12 w-full rounded-lg border border-border/70 bg-card/70 pl-12 pr-4 text-sm text-foreground outline-none backdrop-blur transition-colors placeholder:text-muted-foreground focus:border-accent focus:ring-2 focus:ring-accent/20"
                />
              </div>
            </div>

            <button
              type="button"
              className="inline-flex h-12 shrink-0 items-center justify-center gap-2 rounded-lg border border-border/70 bg-card/70 px-5 text-sm font-semibold text-foreground backdrop-blur transition-colors hover:border-accent/60 hover:bg-card/90"
            >
              <Filter className="h-4 w-4" />
              Filters
            </button>
          </div>

          <div className="mt-12 grid gap-6 sm:grid-cols-2 xl:grid-cols-4">
            {filteredConcepts.map((concept) => (
              <article
                key={concept.id}
                className="group min-h-[192px] rounded-lg border border-border/60 bg-card/70 p-6 shadow-[0_22px_60px_oklch(0.06_0.03_260_/_0.22)] backdrop-blur transition-all duration-300 hover:-translate-y-1 hover:border-accent/60 hover:bg-card/90"
              >
                <div className="flex items-start justify-between gap-4">
                  <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg bg-gradient-to-br from-violet-500 to-sky-500 shadow-[0_14px_30px_oklch(0.12_0.04_260_/_0.32)]">
                    <concept.icon className="h-6 w-6 text-[oklch(0.1_0.03_260)]" />
                  </div>

                  <span className="inline-flex h-6 shrink-0 items-center gap-2 rounded-full border border-accent/40 bg-accent/10 px-3 text-[10px] font-semibold uppercase text-accent">
                    <span className="h-1.5 w-1.5 rounded-full bg-accent" />
                    {concept.status}
                  </span>
                </div>

                <h2 className="mt-6 font-display text-xl font-bold leading-tight tracking-normal text-foreground">
                  {concept.title}
                </h2>
                <p className="mt-2 max-w-[15rem] text-sm leading-6 text-muted-foreground">
                  {concept.description}
                </p>
              </article>
            ))}
          </div>

          {filteredConcepts.length === 0 ? (
            <div className="mt-8 rounded-lg border border-border/60 bg-card/70 px-6 py-8 text-center text-muted-foreground backdrop-blur">
              No concepts match "{query}".
            </div>
          ) : null}
        </section>
      </main>
    </div>
  );
}
