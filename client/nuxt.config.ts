import tailwindcss from "@tailwindcss/vite";
import { defineNuxtConfig } from "nuxt/config";

export default defineNuxtConfig({
  app: {
    head: {
      meta: [{ content: "Personal finance management", name: "description" }],
      title: "Arthika",
      titleTemplate: "%s | Arthika",
    },
  },
  compatibilityDate: "2026-06-23",
  css: ["~/assets/css/main.css"],
  future: {
    compatibilityVersion: 4,
  },
  modules: ["@pinia/nuxt", "reka-ui/nuxt"],
  nitro: {
    routeRules: {
      "/api/**": {
        proxy: {
          to: "http://localhost:8000/api/**",
        },
      },
    },
  },
  ssr: true,
  typescript: {
    strict: true,
  },
  vite: {
    optimizeDeps: {
      include: ["@iconify/vue"],
    },
    plugins: [tailwindcss()],
  },
});
