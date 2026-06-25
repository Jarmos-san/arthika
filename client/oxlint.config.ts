import { defineConfig } from "oxlint";

export default defineConfig({
  plugins: [
    "typescript",
    "unicorn",
    "oxc",
    "vue",
    "jsdoc",
    "vitest",
    "eslint",
    "import",
  ],
  categories: {
    correctness: "error",
  },
  rules: {},
  env: {
    browser: true,
    builtin: true,
    vitest: true,
    vue: true,
  },
  options: {
    typeAware: true,
    typeCheck: true,
  },
  settings: {
    jsdoc: {
      augmentsExtendsReplacesDocs: true,
      exemptDestructuredRootsFromChecks: true,
    },
  },
});
