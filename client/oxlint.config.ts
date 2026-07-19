import { defineConfig } from "oxlint";

export default defineConfig({
  categories: {
    correctness: "error",
    nursery: "error",
    pedantic: "warn",
    perf: "error",
    restriction: "error",
    style: "warn",
    suspicious: "error",
  },
  env: {
    amd: true,
    builtin: true,
    "shared-node-browser": true,
    vitest: true,
    vue: true,
  },
  options: {
    maxWarnings: 10,
    reportUnusedDisableDirectives: "error",
    respectEslintDisableDirectives: false,
    typeAware: true,
    typeCheck: true,
  },
  overrides: [
    {
      files: ["generated/models/*.ts"],
      rules: {
        "jsdoc/check-tag-names": "allow",
      },
    },
    {
      files: ["*.config.ts"],
      rules: { "import/no-default-export": "allow" },
    },
  ],
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
  settings: {
    vitest: {
      typecheck: true,
    },
  },
});
