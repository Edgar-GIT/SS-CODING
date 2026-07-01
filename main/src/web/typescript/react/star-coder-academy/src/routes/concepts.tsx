import { createFileRoute } from "@tanstack/react-router";

import { ConceptsPage } from "@/pages/ConceptsPage";

export const Route = createFileRoute("/concepts")({
  head: () => ({
    meta: [
      { title: "Concepts - SS Coding" },
      {
        name: "description",
        content: "Study programming concepts, patterns, and engineering fundamentals.",
      },
    ],
  }),
  component: ConceptsPage,
});
