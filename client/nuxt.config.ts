export default defineNuxtConfig({
  compatibilityDate: "2026-06-23",
  future: {
    compatibilityVersion: 4,
  },
  modules: ["@nuxt/ui", "reka-ui/nuxt"],
  ssr: false,
  typescript: {
    strict: true,
  },
});
