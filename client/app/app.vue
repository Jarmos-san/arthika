<script setup lang="ts">
  import useAuthStore from "~/stores/auth";

  const auth = useAuthStore();
  await auth.fetchUser();
</script>

<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>
