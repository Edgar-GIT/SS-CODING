const DIGITS = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9];

const REEL_SPEEDS = [0.07, 0.05, 0.09, 0.06, 0.08, 0.04, 0.07, 0.05];

function Reel({ speed }: { speed: number }) {
  return (
    <span className="relative overflow-hidden inline-block w-[0.62em] h-[1em] align-top">
      <span
        className="absolute left-0 top-0 flex flex-col animate-slot-roll"
        style={{ animationDuration: `${speed}s` }}
      >
        {DIGITS.map((digit, idx) => (
          <span
            key={idx}
            className="h-[1em] leading-none flex items-center justify-center text-gradient-cosmic"
          >
            {digit}
          </span>
        ))}
      </span>
    </span>
  );
}

export function SlotCounter() {
  return (
    <div
      aria-label="Programmers joining every second"
      className="font-display text-2xl md:text-3xl font-bold flex items-center leading-none tabular-nums min-h-[1.15em]"
    >
      {REEL_SPEEDS.slice(0, 3).map((speed, i) => (
        <Reel key={`a-${i}`} speed={speed} />
      ))}
      <span className="text-gradient-cosmic w-[0.3em] text-center">,</span>
      {REEL_SPEEDS.slice(3).map((speed, i) => (
        <Reel key={`b-${i}`} speed={speed} />
      ))}
      <span className="text-gradient-cosmic ml-0.5">+</span>
    </div>
  );
}
