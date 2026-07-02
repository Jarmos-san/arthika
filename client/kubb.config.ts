import { pluginTs } from "@kubb/plugin-ts";
import { defineConfig } from "kubb";

export default defineConfig({
  input: {
    path: "../server/api/openapi.yml",
  },
  output: {
    path: "./generated",
    clean: true,
    format: "oxfmt",
    lint: "oxlint",
  },
  plugins: [
    pluginTs({
      output: {
        path: "models",
      },
    }),
  ],
});
