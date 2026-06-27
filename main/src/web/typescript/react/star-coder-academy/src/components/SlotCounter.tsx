const DIGITS = "0123456789";

// Each "reel" is a vertical strip that scrolls infinitely upward.
// We don't show a real number — just the illusion of an ever-growing count.
const reels = [
  { width: "ch", speed: 1.1 },
  { width: "ch", speed: 1.4 },
  { width: "ch", speed: 0.9 },
  { width: "0.4ch", isComma: true },
  { width: "ch", speed: 1.7 },
  { width: "ch", speed: 1.2 },
  { width: "ch", speed: 1.5 },
];

export function SlotCounter() {
  return (
    <div
      aria-label="Programmers joining every second"
      className="font-display text-2xl md:text-3xl font-bold text-gradient-cosmic flex items-center leading-none"
      style={{ height: "1.2em" }}
    >
      {reels.map((r, i) =>
        r.isComma ? (
          <span key={i} className="px-[1px]">,</span>
        ) : (
          <span
            key={i}
            className="relative overflow-hidden inline-block"
            style={{ width: "1ch", height: "1em" }}
          >
            <span
              className="absolute left-0 top-0 flex flex-col"
              style={{
                animation: `slot-roll ${r.speed}s linear infinite`,
              }}
            >
              {(DIGITS + DIGITS).split("").map((d, idx) => (
                <span
                  key={idx}
                  style={{ height: "1em", lineHeight: "1em" }}
                  className="block text-center"
                >
                  {d}
                </span>
              ))}
            </span>
          </span>
        ),
      )}
      <span className="ml-1">+</span>
    </div>
  );
}
