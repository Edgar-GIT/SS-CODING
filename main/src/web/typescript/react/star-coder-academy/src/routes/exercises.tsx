import { createFileRoute } from "@tanstack/react-router";

import { ExercisesPage } from "@/pages/ExercisesPage";

export const Route = createFileRoute("/exercises")({
  head: () => ({
    meta: [
      { title: "Exercises - SS Coding" },
      {
        name: "description",
        content: "Practice programming with focused exercises and project checkpoints.",
      },
    ],
  }),
  component: ExercisesPage,
});
