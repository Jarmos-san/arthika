import { defineConfig } from "oxlint";

export default defineConfig({
  plugins: [
    "typescript",
    "unicorn",
    "oxc",
    "eslint",
    "import",
    "jsdoc",
    "vitest",
    "promise",
    "node",
    "vue",
  ],
  categories: {
    correctness: "error",
    // nursery: "error",
    // pedantic: "warn",
    // perf: "warn",
    // restriction: "error",
    // style: "warn",
    // suspicious: "error",
  },
  rules: {},
  env: {
    amd: true,
    builtin: true,
    "shared-node-browser": true,
    vitest: true,
    vue: true,
  },
  options: {
    typeAware: true,
    typeCheck: true,
    maxWarnings: 10,
    reportUnusedDisableDirectives: "error",
    respectEslintDisableDirectives: false,
  },
  settings: {
    vitest: {
      typecheck: true,
    },
  },
});
