import { Link } from "@tanstack/react-router";
import { User } from "lucide-react";
import { useEffect, useState } from "react";

const links = [
  { to: "/", label: "Home" },
  { to: "/exercises", label: "Exercises" },
  { to: "/guides", label: "Guides" },
  { to: "/concepts", label: "Concepts" },
  { to: "/ranking", label: "Ranking" },
  { to: "/about", label: "About Us" },
] as const;

function useScrollProgress() {
  const [p, setP] = useState(0);
  useEffect(() => {
    const onScroll = () => {
      const h = document.documentElement;
      const total = h.scrollHeight - h.clientHeight;
      setP(total > 0 ? (h.scrollTop / total) * 100 : 0);
    };
    onScroll();
    window.addEventListener("scroll", onScroll, { passive: true });
    window.addEventListener("resize", onScroll);
    return () => {
      window.removeEventListener("scroll", onScroll);
      window.removeEventListener("resize", onScroll);
    };
  }, []);
  return p;
}

export function Navbar() {
  const progress = useScrollProgress();
  return (
    <header className="sticky top-0 z-50 backdrop-blur-xl bg-background/70 border-b border-border/60">
      <nav className="mx-auto max-w-7xl px-6 h-14 flex items-center justify-between">
        <Link to="/" className="group shrink-0 select-none">
          <span className="font-display text-2xl font-bold tracking-tight inline-flex items-baseline gap-2">
            <span className="text-gradient-cosmic animate-brand-shimmer drop-shadow-[0_0_12px_oklch(0.55_0.25_285_/_0.5)]">
              SS
            </span>
            <span className="text-foreground transition-all duration-300 group-hover:tracking-widest group-hover:text-stardust">
              CODING
            </span>
          </span>
        </Link>

        <ul className="hidden md:flex items-center gap-1">
          {links.map((l) => (
            <li key={l.to}>
              <Link
                to={l.to}
                activeOptions={{ exact: l.to === "/" }}
                activeProps={{
                  className:
                    "text-foreground bg-secondary/60 shadow-[0_0_20px_oklch(0.55_0.25_285_/_0.35)]",
                }}
                inactiveProps={{ className: "text-muted-foreground hover:text-foreground" }}
                className="relative px-3 py-2 rounded-lg text-sm font-medium transition-all hover:bg-secondary/40"
              >
                {l.label}
              </Link>
            </li>
          ))}
        </ul>

        <Link
          to="/profile"
          aria-label="Profile"
          className="hidden md:inline-flex relative items-center justify-center h-10 w-10 rounded-full bg-cosmic shadow-cosmic hover:scale-110 transition-transform ring-2 ring-primary/40 hover:ring-primary"
        >
          <span className="absolute inset-0 rounded-full bg-cosmic opacity-60 blur-md animate-pulse-glow -z-10" />
          <User className="w-5 h-5 text-primary-foreground" strokeWidth={2.5} />
        </Link>
      </nav>

      <div className="absolute bottom-0 left-0 right-0 h-[2px] bg-border/30">
        <div
          className="h-full bg-cosmic shadow-[0_0_10px_oklch(0.55_0.25_285_/_0.9)] transition-[width] duration-100 ease-out"
          style={{ width: `${progress}%` }}
        />
      </div>
    </header>
  );
}
