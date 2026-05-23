import { defineConfig } from "vite";
import react, { reactCompilerPreset } from "@vitejs/plugin-react";
import babel from "@rolldown/plugin-babel";
import tailwindcss from "@tailwindcss/vite";

// https://vite.dev/config/
export default defineConfig({
  plugins: [tailwindcss(), react(), babel({ presets: [reactCompilerPreset()] })],
  server: {
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        rewrite: (path) => path.replace(/^\/api/, ""),
      },
      "/news-api": {
        target: "http://localhost:8081",
        rewrite: (path) => path.replace(/^\/news-api/, ""),
      },
      "/langwar.": {
        target: "http://localhost:8080",
      },
    },
  },
});
