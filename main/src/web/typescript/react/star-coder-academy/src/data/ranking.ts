export type LeaderboardPodiumEntry = {
  rank: number;
  username: string;
  title: string;
  tier: string;
  xp: number;
  solved: number;
  initials: string;
  accentClassName: string;
  trophyClassName: string;
};

export type LeaderboardRow = {
  rank: number;
  username: string;
  tier: string;
  xp: number;
};

export const podium: LeaderboardPodiumEntry[] = [
  {
    rank: 2,
    username: "stellar_byte",
    title: "Galactic",
    tier: "Nebula",
    xp: 172640,
    solved: 1284,
    initials: "SB",
    accentClassName: "from-slate-200 to-sky-300 text-slate-950",
    trophyClassName: "text-slate-300",
  },
  {
    rank: 1,
    username: "nova.dev",
    title: "Galactic Overlord",
    tier: "Supernova",
    xp: 194320,
    solved: 1612,
    initials: "ND",
    accentClassName: "from-yellow-300 to-amber-500 text-slate-950",
    trophyClassName: "text-yellow-300",
  },
  {
    rank: 3,
    username: "quantum.fox",
    title: "Galactic",
    tier: "Nebula",
    xp: 158900,
    solved: 1107,
    initials: "QF",
    accentClassName: "from-orange-300 to-rose-500 text-slate-950",
    trophyClassName: "text-orange-300",
  },
];

export const leaderboardRows: LeaderboardRow[] = [
  { rank: 4, username: "void.runner", tier: "Nebula", xp: 150000 },
  { rank: 5, username: "byteweaver", tier: "Nebula", xp: 145800 },
  { rank: 6, username: "stack_wraith", tier: "Quantum", xp: 139420 },
  { rank: 7, username: "syntax.sage", tier: "Quantum", xp: 132900 },
  { rank: 8, username: "logicforge", tier: "Orbit", xp: 126300 },
];
