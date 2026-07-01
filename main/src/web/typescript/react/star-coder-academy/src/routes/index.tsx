import { createFileRoute } from "@tanstack/react-router";

import { HomePage } from "@/pages/HomePage";

export const Route = createFileRoute("/")({
  head: () => ({
    meta: [
      { title: "SS Coding - Become an Elite Programmer" },
      {
        name: "description",
        content:
          "Master any programming language through interactive guides, exercises, quizzes, projects, and global rankings. Become an elite developer.",
      },
      { property: "og:title", content: "SS Coding - Become an Elite Programmer" },
      {
        property: "og:description",
        content:
          "Master any programming language through interactive guides, exercises, quizzes, projects, and global rankings.",
      },
    ],
  }),
  component: HomePage,
});
