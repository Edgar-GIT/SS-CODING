import { createFileRoute } from "@tanstack/react-router";

import { RankingPage } from "@/pages/RankingPage";

export const Route = createFileRoute("/ranking")({
  head: () => ({
    meta: [
      { title: "Ranking - SS Coding" },
      {
        name: "description",
        content: "Follow rankings, milestones, and competitive coding progress.",
      },
    ],
  }),
  component: RankingPage,
});
