import { defineConfig } from "oxfmt";

export default defineConfig({
  printWidth: 80,
  sortImports: true,
  sortTailwindcss: true,
  sortPackageJson: true,
  htmlWhitespaceSensitivity: "strict",
  jsdoc: {
    descriptionTag: true,
    descriptionWithDot: true,
    preferCodeFences: true,
    separateReturnsFromParam: true,
  },
  proseWrap: "always",
  vueIndentScriptAndStyle: true,
});
