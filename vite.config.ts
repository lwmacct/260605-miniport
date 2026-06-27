import react from "@vitejs/plugin-react";
import fs from "node:fs";
import path from "node:path";
import { pathToFileURL } from "node:url";
import { defineConfig, mergeConfig, type UserConfig } from "vite";

const apiTarget = process.env.API_PROXY_TARGET ?? "http://localhost:40238";
const rootDir = import.meta.dirname;

async function loadLocalConfig(): Promise<UserConfig> {
  const localConfigPath = path.resolve(rootDir, "vite.local.ts");

  if (!fs.existsSync(localConfigPath)) {
    return {};
  }

  const mod = await import(/* @vite-ignore */ pathToFileURL(localConfigPath).href);
  return mod.default ?? {};
}

export default defineConfig(async () =>
  mergeConfig(
    {
      plugins: [react()],
      resolve: {
        alias: {
          "@": path.resolve(rootDir, "src"),
        },
      },
      build: {
        chunkSizeWarningLimit: 5120,
      },
      server: {
        host: "0.0.0.0",
        port: 40239,
        strictPort: true,
        proxy: {
          "/api": {
            target: apiTarget,
            changeOrigin: true,
            ws: true,
          },
        },
      },
    },
    await loadLocalConfig(),
  ),
);
