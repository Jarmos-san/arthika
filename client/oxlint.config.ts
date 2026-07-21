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
        "import/no-relative-parent-imports": "allow",
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
    {
      files: ["generated/**/*.ts"],
      rules: {
        "import/group-exports": "allow",
        "import/no-named-export": "allow",
        "unicorn/filename-case": ["error", { case: "pascalCase" }],
      },
    },
    {
      files: ["app/**/*.vue"],
      rules: {
        "import/consistent-type-specifier-style": "allow",
        "vitest/require-hook": "off",
        "vue/define-props-destructuring": "allow",
      },
    },
    {
      files: ["app/components/**/*.vue"],
      rules: { "unicorn/filename-case": ["error", { case: "pascalCase" }] },
    },
    {
      files: ["app/**/*.ts"],
      rules: { "import/no-named-export": "allow" },
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
    "no-default-export": "allow",
    "no-undefined": "allow",
    "oxc/no-async-await": "allow",
  },
  settings: {
    vitest: {
      typecheck: true,
    },
  },
});
