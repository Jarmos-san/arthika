<script setup lang="ts">
  import { ref } from "vue";

  const email = ref("");
  const password = ref("");
</script>

<template>
  <div
    class="flex min-h-screen items-center justify-center from-amber-100/40 via-amber-50 to-amber-50"
  >
    <div
      class="w-full max-w-sm rounded-xl border border-stone-200/60 bg-white p-10 shadow-lg"
    >
      <h1 class="text-2xl font-semibold tracking-tight text-stone-800">
        Welcome to Arthika
      </h1>
      <p class="mt-1.5 text-sm text-stone-400">
        Let's set up your account to get started.
      </p>

      <div class="my-6 h-px bg-stone-200" />

      <form class="flex flex-col gap-5" @submit.prevent>
        <Input
          v-model="email"
          id="email"
          label="Email"
          type="email"
          placeholder="you@example.com"
        />

        <Input
          v-model="password"
          id="password"
          label="Password"
          type="password"
          placeholder="Minimum 8 characters"
        />

        <button type="submit" class="btn-primary mt-1"> Create Account </button>
      </form>
    </div>
  </div>
</template>
