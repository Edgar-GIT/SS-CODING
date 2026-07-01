import { Code2, Target, Zap } from "lucide-react";
import type { LucideIcon } from "lucide-react";

export type ExerciseTrack = {
  id: string;
  title: string;
  count: number;
  status: string;
  icon: LucideIcon;
  searchTerms: string[];
};

export const exerciseTracks: ExerciseTrack[] = [
  {
    id: "fundamentals",
    title: "Fundamentals",
    count: 240,
    status: "Coming Soon",
    icon: Code2,
    searchTerms: ["basics", "syntax", "language"],
  },
  {
    id: "algorithms",
    title: "Algorithms",
    count: 510,
    status: "Coming Soon",
    icon: Zap,
    searchTerms: ["logic", "problem solving", "topic"],
  },
  {
    id: "data-structures",
    title: "Data Structures",
    count: 320,
    status: "Coming Soon",
    icon: Target,
    searchTerms: ["arrays", "trees", "graphs", "topic"],
  },
];
