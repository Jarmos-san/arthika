<script setup lang="ts">
  import useToast from "~/composables/useToast";

  const { dismiss, messages } = useToast();
</script>

<template>
  <ToastProvider>
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
    <AppToast
      v-for="(msg, i) in messages"
      :key="i"
      :duration="3000"
      :message="msg"
      @dismiss="dismiss(i)"
    />
    <ToastViewport
      class="fixed right-0 bottom-0 z-50 flex flex-col gap-2 p-4"
    />
  </ToastProvider>
</template>
