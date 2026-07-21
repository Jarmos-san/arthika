<script setup lang="ts">
  import { useCookie } from "nuxt/app";

  import useAuthStore from "~/stores/auth";

  const auth = useAuthStore();
  const token = useCookie("token");
  if (token.value) {
    try {
      const [segment] = token.value.split(".");
      if (segment) {
        const payload = JSON.parse(atob(segment));
        if (
          typeof payload.sub === "string" &&
          typeof payload.email === "string"
        ) {
          auth.isAuthenticated = true;
        } else {
          token.value = undefined;
        }
      } else {
        token.value = undefined;
      }
    } catch {
      token.value = undefined;
    }
  }
</script>

<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>
