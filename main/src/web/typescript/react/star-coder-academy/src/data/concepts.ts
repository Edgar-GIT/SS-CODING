import { Atom, Cpu, GitBranch, Layers } from "lucide-react";
import type { LucideIcon } from "lucide-react";

export type Concept = {
  id: string;
  title: string;
  description: string;
  status: string;
  icon: LucideIcon;
};

export const concepts: Concept[] = [
  {
    id: "recursion",
    title: "Recursion",
    description: "Think in self-referential structures.",
    status: "Coming Soon",
    icon: Atom,
  },
  {
    id: "abstraction",
    title: "Abstraction",
    description: "Hide the noise, expose the meaning.",
    status: "Coming Soon",
    icon: Layers,
  },
  {
    id: "control-flow",
    title: "Control Flow",
    description: "Branches, loops and pattern matching.",
    status: "Coming Soon",
    icon: GitBranch,
  },
  {
    id: "concurrency",
    title: "Concurrency",
    description: "Threads, async, parallelism - done right.",
    status: "Coming Soon",
    icon: Cpu,
  },
];
