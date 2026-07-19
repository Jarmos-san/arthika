<script setup lang="ts">
  const auth = useAuthStore();
  const token = useCookie("token");
  if (token.value) {
    try {
      const segment = token.value.split(".")[1];
      if (segment) {
        const payload = JSON.parse(atob(segment));
        if (
          typeof payload.sub === "string" &&
          typeof payload.email === "string"
        ) {
          auth.isAuthenticated = true;
        } else {
          token.value = null;
        }
      } else {
        token.value = null;
      }
    } catch {
      token.value = null;
    }
  }
</script>

<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>
