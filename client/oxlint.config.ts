import { defineConfig } from "oxlint";

export default defineConfig({
  plugins: ["typescript", "unicorn", "oxc", "vue", "jsdoc", "vitest", "eslint"],
  categories: { correctness: "error" },
  rules: {},
  env: {
    builtin: true,
  },
  options: {
    typeAware: true,
    typeCheck: true,
  },
});
