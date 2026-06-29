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
  globals: {
    $fetch: "readonly",
    computed: "readonly",
    navigateTo: "readonly",
    reactive: "readonly",
    readonly: "readonly",
    ref: "readonly",
    useAuth: "readonly",
    useCookie: "readonly",
    useState: "readonly",
  },
  options: {
    typeAware: true,
    typeCheck: true,
  },
  overrides: [
    {
      files: ["app/composables/**/*"],
      rules: {
        "import/no-named-export": "off",
        "typescript/explicit-module-boundary-types": "off",
        "unicorn/filename-case": "off",
        "unicorn/prefer-ternary": "off",
      },
    },
    {
      files: ["*.config.ts"],
      rules: {
        "import/no-default-export": "off",
      },
    },
    {
      files: ["*.vue"],
      rules: {
        "import/unambiguous": "off",
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
  rules: {
    "eslint/no-ternary": "off",
    "eslint/no-undefined": "off",
    "oxc/no-async-await": "off",
    "typescript/explicit-function-return-type": "off",
    "unicorn/no-useless-undefined": "off",
  },
  settings: {
    jsdoc: {
      augmentsExtendsReplacesDocs: true,
      exemptDestructuredRootsFromChecks: true,
    },
  },
});
