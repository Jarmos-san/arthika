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
    node: true,
    "shared-node-browser": true,
    vitest: true,
    vue: true,
  },
  globals: {
    defineNuxtRouteMiddleware: "readonly",
    navigateTo: "readonly",
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
      rules: {
        "import/no-default-export": "allow",
        "no-ternary": "allow",
        "no-undefined": "allow",
        "node/no-process-env": [
          "error",
          { allowedVariables: ["CI", "GITHUB_ACTIONS"] },
        ],
        "typescript/strict-boolean-expressions": "allow",
      },
    },
    {
      files: ["tests/**/*.test.ts"],
      rules: {
        "vitest/no-importing-vitest-globals": "allow",
        "vitest/prefer-strict-boolean-matchers": "allow",
      },
    },
    {
      files: ["app/middleware/**/*.ts"],
      rules: {
        "import/no-default-export": "allow",
        "typescript/prefer-readonly-parameter-types": "allow",
      },
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
  rules: {
    "eslint/sort-imports": "off",
    "oxc/no-async-await": "allow",
  },
  settings: {
    vitest: {
      typecheck: true,
    },
  },
});
