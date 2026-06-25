import { defineConfig } from "oxfmt";

export default defineConfig({
  htmlWhitespaceSensitivity: "strict",
  jsdoc: {
    descriptionTag: true,
    descriptionWithDot: true,
    preferCodeFences: true,
    separateReturnsFromParam: true,
  },
  printWidth: 80,
  proseWrap: "always",
  sortImports: true,
  sortPackageJson: true,
  sortTailwindcss: true,
  vueIndentScriptAndStyle: true,
});
