import { pluginTs } from "@kubb/plugin-ts";
import { defineConfig } from "kubb";

export default defineConfig({
  input: {
    path: "../openapi.yml",
  },
  output: {
    clean: true,
    format: "oxfmt",
    lint: "oxlint",
    path: "./app/openapi",
  },
  plugins: [
    pluginTs({
      output: {
        path: "models",
      },
    }),
  ],
});
