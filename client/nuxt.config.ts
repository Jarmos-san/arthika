import tailwindcss from "@tailwindcss/vite";

export default defineNuxtConfig({
  compatibilityDate: "2026-06-23",
  app: {
    head: {
      titleTemplate: "%s | Arthika",
      title: "Arthika",
      meta: [{ name: "description", content: "Personal finance management" }],
    },
  },
  css: ["./app/assets/css/main.css"],
  future: {
    compatibilityVersion: 4,
  },
  modules: ["@pinia/nuxt", "reka-ui/nuxt"],
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
