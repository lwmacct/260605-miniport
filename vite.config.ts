import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

const apiTarget = process.env.VITE_API_TARGET ?? "http://127.0.0.1:40238";

export default defineConfig({
  plugins: [react()],
  build: {
    chunkSizeWarningLimit: 5120,
  },
  server: {
    host: "0.0.0.0",
    port: 40239,
    proxy: {
      "/api": {
        target: apiTarget,
        changeOrigin: true,
      },
    },
  },
});
