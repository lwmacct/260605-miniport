import react from "@vitejs/plugin-react";
import path from "node:path";
import { defineConfig } from "vite";

const apiTarget = process.env.API_PROXY_TARGET ?? "http://localhost:40238";
const rootDir = import.meta.dirname;

export default defineConfig({
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
      "/authme": {
        target: apiTarget,
      },
      "/api": {
        target: apiTarget,
        ws: true,
      },
    },
  },
});
