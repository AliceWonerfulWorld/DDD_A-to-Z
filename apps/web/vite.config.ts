import { defineConfig } from "vitest/config";
import react, { reactCompilerPreset } from "@vitejs/plugin-react";
import babel from "@rolldown/plugin-babel";
import tailwindcss from "@tailwindcss/vite";
import VitestUddReporter from "vitest-udd-reporter";
import type { TestModule, SerializedError } from "vitest/node";

// vitest v4 では onFinished が呼ばれないケースがあるため onTestRunEnd も override する
class UddReporterV4 extends VitestUddReporter {
  async onTestRunEnd(
    testModules: ReadonlyArray<TestModule>,
    unhandledErrors: ReadonlyArray<SerializedError>,
    reason: "passed" | "failed" | "interrupted",
  ) {
    await super.onTestRunEnd(testModules, unhandledErrors, reason);
    // onFinished が呼ばれない場合に備えて即出力
    if (this.message) {
      console.log(this.message);
      this.message = "";
    }
  }
}

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
        target: "http://localhost:8082",
        rewrite: (path) => path.replace(/^\/news-api/, ""),
      },
      "/langwar.": {
        target: "http://localhost:8080",
      },
    },
  },
  test: {
    reporters: ["default", new UddReporterV4()],
  },
});
