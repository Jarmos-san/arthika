import tailwindcss from "@tailwindcss/vite";

export default defineNuxtConfig({
  compatibilityDate: "2026-06-23",
  css: ["./app/assets/css/main.css"],
  future: {
    compatibilityVersion: 4,
  },
  modules: ["@nuxt/ui", "reka-ui/nuxt"],
  ssr: false,
  typescript: {
    strict: true,
  },
  vite: {
    plugins: [tailwindcss()],
  },
});
