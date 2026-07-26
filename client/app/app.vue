<script setup lang="ts">
  import { useCookie } from "nuxt/app";

  import useToast from "~/composables/useToast";
  import useAuthStore from "~/stores/auth";

  const auth = useAuthStore();
  const { dismiss, messages } = useToast();
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
