import { defineVitestProject } from "@nuxt/test-utils/config";
import {
  configDefaults,
  coverageConfigDefaults,
  defineConfig,
} from "vitest/config";

const MIN_BAIL = 0;
const MAX_FAILURES = 3;
const MAX_RETRIES_IN_CI = 3;
const MAX_RETRIES_IN_LOCAL = 0;
const MAX_WORKERS = 4;

export default defineConfig({
  test: {
    bail: process.env.CI ? MAX_FAILURES : MIN_BAIL,
    coverage: {
      enabled: true,
      exclude: [
        "*.config.ts",
        "**/app.vue",
        "**/types/*.ts",
        "**/pages/**/*.vue",
        "**/layouts/**/*.vue",
        ...coverageConfigDefaults.exclude,
      ],
      reporter: ["text", "json"],
    },
    fileParallelism: process.env.CI ? true : undefined,
    logHeapUsage: true,
    maxWorkers: process.env.CI ? MAX_WORKERS : "100%",
    pool: "threads",
    printConsoleTrace: true,
    projects: [
      await defineVitestProject({
        test: {
          environment: "nuxt",
          include: ["tests/nuxt/**/*.test.ts"],
          isolate: false,
          mockReset: true,
          name: "nuxt",
        },
      }),
      {
        test: {
          environment: "node",
          include: ["tests/unittests/**/*.test.ts"],
          mockReset: true,
          name: "unit",
        },
      },
    ],
    reporters: [
      ...configDefaults.reporters,
      ...(process.env.GITHUB_ACTIONS ? ["github-actions"] : []),
    ],
    retry: process.env.CI ? MAX_RETRIES_IN_CI : MAX_RETRIES_IN_LOCAL,
    silent: process.env.CI ? false : "passed-only",
    watch: false,
  },
});
