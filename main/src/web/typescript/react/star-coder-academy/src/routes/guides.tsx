import { createFileRoute } from "@tanstack/react-router";

import { GuidesPage } from "@/pages/GuidesPage";

export const Route = createFileRoute("/guides")({
  head: () => ({
    meta: [
      { title: "Guides - SS Coding" },
      {
        name: "description",
        content:
          "Choose a programming language and follow a complete curriculum from fundamentals to advanced engineering.",
      },
    ],
  }),
  component: GuidesPage,
});
