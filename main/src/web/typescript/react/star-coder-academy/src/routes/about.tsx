import { createFileRoute } from "@tanstack/react-router";

import { AboutPage } from "@/pages/AboutPage";

export const Route = createFileRoute("/about")({
  head: () => ({
    meta: [
      { title: "About Us - SS Coding" },
      {
        name: "description",
        content: "Learn more about SS Coding and its learning philosophy.",
      },
    ],
  }),
  component: AboutPage,
});
