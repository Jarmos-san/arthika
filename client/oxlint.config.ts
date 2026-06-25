import { defineConfig } from "oxlint";

export default defineConfig({
  categories: {
    correctness: "error",
    nursery: "error",
    pedantic: "warn",
    perf: "warn",
    restriction: "warn",
    style: "warn",
    suspicious: "warn",
  },
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
  overrides: [
    {
      files: ["app/composables/**/*"],
      rules: {
        "unicorn/filename-case": "off",
      },
    },
    {
      files: ["*.config.ts"],
      rules: {
        "import/no-default-export": "off",
      },
    },
  ],
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
  settings: {
    jsdoc: {
      augmentsExtendsReplacesDocs: true,
      exemptDestructuredRootsFromChecks: true,
    },
  },
});
