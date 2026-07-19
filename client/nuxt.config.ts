import tailwindcss from "@tailwindcss/vite";

export default defineNuxtConfig({
  compatibilityDate: "2026-06-23",
  css: ["./app/assets/css/main.css"],
  future: {
    compatibilityVersion: 4,
  },
  modules: ["reka-ui/nuxt"],
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
