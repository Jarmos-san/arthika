<script setup>
useHead({
  title: "Arthika - Portfolio Tracker",
});

useSeoMeta;
</script>

<template>
  <UApp>
    <UHeader>
      <template #left>
        <NuxtLink to="/" class="font-bold text-lg"> Arthika </NuxtLink>
      </template>
    </UHeader>

    <UMain>
      <NuxtPage />
    </UMain>

    <UFooter>
      <p class="text-sm text-muted">
        Arthika — &copy; {{ new Date().getFullYear() }}
      </p>
    </UFooter>
  </UApp>
</template>
