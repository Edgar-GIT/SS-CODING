import { createFileRoute } from "@tanstack/react-router";

import { ProfilePage } from "@/pages/ProfilePage";

export const Route = createFileRoute("/profile")({
  head: () => ({
    meta: [
      { title: "Profile - SS Coding" },
      {
        name: "description",
        content: "View coding progress, completed chapters, and personal learning stats.",
      },
    ],
  }),
  component: ProfilePage,
});
