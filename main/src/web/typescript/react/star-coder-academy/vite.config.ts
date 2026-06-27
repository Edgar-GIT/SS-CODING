import path from "node:path";
import { fileURLToPath } from "node:url";
import tailwindcss from "@tailwindcss/vite";
import { tanstackStart } from "@tanstack/react-start/plugin/vite";
import viteReact from "@vitejs/plugin-react";
import { nitro } from "nitro/vite";
import { defineConfig } from "vite";

const appRoot = fileURLToPath(new URL(".", import.meta.url));
const resourcesDir = path.resolve(appRoot, "../../../../../resources");

export default defineConfig({
  plugins: [
    tanstackStart({
      server: { entry: "server" },
    }),
    viteReact(),
    tailwindcss(),
    nitro(),
  ],
  resolve: {
    alias: {
      "@resources": resourcesDir,
    },
    dedupe: ["react", "react-dom"],
    tsconfigPaths: true,
  },
  server: {
    fs: {
      allow: [appRoot, resourcesDir],
    },
  },
});
