<script setup lang="ts">
  import { useHead } from "nuxt/app";

  useHead({ title: "Dashboard" });
</script>

<template>
  <div>
    <AppNavbar />
    <div class="p-6">Hello World!</div>
  </div>
</template>
