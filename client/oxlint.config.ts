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
    $fetch: "readonly",
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
      files: ["app/openapi/models/*.ts"],
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
        "node/no-top-level-await": "allow",
        "typescript/strict-boolean-expressions": "allow",
      },
    },
    {
      files: ["tests/**/*.test.ts"],
      rules: {
        "import/no-relative-parent-imports": "allow",
        "max-lines-per-function": "allow",
        "max-statements": "allow",
        "vitest/no-hooks": ["warn", { allow: ["afterEach", "beforeEach"] }],
        "vitest/no-importing-vitest-globals": "allow",
        "vitest/prefer-strict-boolean-matchers": "allow",
      },
    },
    {
      files: ["app/middleware/**/*.ts"],
      rules: {
        "eslint/require-await": "allow",
        "import/no-default-export": "allow",
        "typescript/prefer-readonly-parameter-types": "allow",
      },
    },
    {
      files: ["app/openapi/**/*.ts"],
      rules: {
        "import/group-exports": "allow",
        "import/no-named-export": "allow",
        "typescript/no-explicit-any": "allow",
        "unicorn/filename-case": ["error", { case: "pascalCase" }],
      },
    },
    {
      files: ["app/types/components/**/*.ts"],
      rules: {
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
      files: ["app/composables/**/*.ts", "tests/nuxt/composables/**/*.test.ts"],
      rules: {
        "import/no-named-export": "allow",
        "unicorn/filename-case": ["error", { case: "camelCase" }],
      },
    },
    {
      files: ["app/**/*.ts"],
      rules: { "import/no-named-export": "allow" },
    },
    {
      files: ["tests/nuxt/components/**/*.test.ts"],
      rules: {
        "typescript/ban-ts-comment": "off",
        "typescript/prefer-ts-expect-error": "allow",
        "unicorn/filename-case": ["error", { case: "pascalCase" }],
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
    "import/consistent-type-specifier-style": [
      "warn",
      "prefer-top-level-if-only-type-imports",
    ],
    "init-declarations": ["warn", "never"],
    "max-lines-per-function": [
      "warn",
      {
        skipBlankLines: true,
        skipComments: true,
      },
    ],
    "no-default-export": "allow",
    "no-undefined": "allow",
    "no-void": ["warn", { allowAsStatement: true }],
    "node/no-process-env": [
      "error",
      { allowedVariables: ["CI", "GITHUB_ACTIONS", "NODE_ENV"] },
    ],

    "oxc/no-async-await": "allow",
  },
  settings: {
    vitest: {
      typecheck: true,
    },
  },
});
