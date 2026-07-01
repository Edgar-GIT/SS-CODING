import { Atom, Cpu, GitBranch, Layers } from "lucide-react";
import type { LucideIcon } from "lucide-react";

export type Concept = {
  id: string;
  title: string;
  description: string;
  status: string;
  icon: LucideIcon;
  searchTerms: string[];
};

export const concepts: Concept[] = [
  {
    id: "recursion",
    title: "Recursion",
    description: "Think in self-referential structures.",
    status: "Coming Soon",
    icon: Atom,
    searchTerms: ["functions", "stack", "self reference"],
  },
  {
    id: "abstraction",
    title: "Abstraction",
    description: "Hide the noise, expose the meaning.",
    status: "Coming Soon",
    icon: Layers,
    searchTerms: ["interfaces", "architecture", "models"],
  },
  {
    id: "control-flow",
    title: "Control Flow",
    description: "Branches, loops and pattern matching.",
    status: "Coming Soon",
    icon: GitBranch,
    searchTerms: ["branches", "loops", "patterns"],
  },
  {
    id: "concurrency",
    title: "Concurrency",
    description: "Threads, async, parallelism - done right.",
    status: "Coming Soon",
    icon: Cpu,
    searchTerms: ["threads", "async", "parallelism"],
  },
];
