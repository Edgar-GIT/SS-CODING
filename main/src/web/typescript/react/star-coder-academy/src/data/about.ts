import { Code2, Cpu, Globe, Heart, Rocket, ShieldAlert, Sparkles, Trophy } from "lucide-react";
import type { LucideIcon } from "lucide-react";

export type AboutItem = {
  title: string;
  description: string;
  icon: LucideIcon;
};

export const feats: AboutItem[] = [
  {
    icon: Trophy,
    title: "World-class competitor",
    description:
      "Countless international programming competitions won — including a NASA challenge — plus trophies across CTFs, algorithms and hardware contests.",
  },
  {
    icon: Code2,
    title: "50+ programming languages",
    description:
      "Fluent in over fifty languages, from low-level assembly and C to modern systems, functional and esoteric languages.",
  },
  {
    icon: Cpu,
    title: "Built my own stack, top to bottom",
    description:
      "Designed my own CPU and GPU, wrote my own operating system, invented my own programming languages, and forged custom PC firmware and BIOS from scratch.",
  },
  {
    icon: ShieldAlert,
    title: "Pentester & gray-hat hacker",
    description:
      "Professional offensive security work — breaking systems to make them stronger, across web, network, hardware and firmware.",
  },
];

export const values: AboutItem[] = [
  {
    icon: Rocket,
    title: "Ambition",
    description: "I believe anyone can become elite with the right guidance and enough obsession.",
  },
  {
    icon: Heart,
    title: "Craft",
    description:
      "Every guide, exercise and quiz is hand-tuned by me for real learning — no filler, no fluff.",
  },
  {
    icon: Globe,
    title: "Accessible",
    description: "Top-tier education shouldn't be locked behind a price tag or a campus.",
  },
  {
    icon: Sparkles,
    title: "Wonder",
    description: "Programming feels like magic. This platform is built to keep that spark alive.",
  },
];
