<script setup lang="ts">
  import { useHead } from "nuxt/app";
  import { ref } from "vue";

  import useAuthStore from "~/stores/auth";
  import type { LoginErrors } from "~/types/utils/validators";
  import validateLogin from "~/utils/validators";

  useHead({ title: "Login" });

  const email = ref<string | undefined>(undefined);
  const password = ref<string | undefined>(undefined);

  const errors = ref<LoginErrors>({
    email: undefined,
    password: undefined,
  });

  const onSubmit = async (): Promise<void> => {
    errors.value = validateLogin(email.value, password.value);

    if (errors.value.email || errors.value.password) {
      return;
    }

    const { login } = useAuthStore();

    if (email.value !== undefined && password.value !== undefined) {
      login(email.value, password.value);
    }

    await navigateTo("/dashboard");
  };
</script>

<template>
  <div class="flex min-h-screen items-center justify-center">
    <div
      class="w-full max-w-sm rounded-xl border border-stone-200/60 bg-white p-10 shadow-lg"
    >
      <h1 class="text-2xl font-semibold tracking-tight text-stone-800">
        Log in to Arthika
      </h1>
      <p class="mt-1.5 text-sm text-stone-400">
        Enter your credentials to continue.
      </p>
      <div class="my-6 h-px bg-stone-200" />

      <form class="flex flex-col gap-5" @submit.prevent="onSubmit">
        <div>
          <Input
            id="email"
            v-model="email"
            aria-describedby="email-error"
            label="Email"
            type="email"
            placeholder="you@example.com"
          />
          <p
            v-if="errors.email"
            id="email-error"
            class="mt-1.5 text-xs text-red-600"
          >
            {{ errors.email }}
          </p>
        </div>

        <div>
          <Input
            id="password"
            v-model="password"
            aria-describedby="password-error"
            label="Password"
            type="password"
            placeholder="Enter your password"
          />
          <p
            v-if="errors.password"
            id="password-error"
            class="mt-1.5 text-xs text-red-600"
          >
            {{ errors.password }}
          </p>
        </div>

        <button class="btn-primary mt-1" type="submit">
          <span>Log In</span>
          <span class="hidden">Logging in...</span>
        </button>
      </form>

      <p class="mt-6 text-center text-sm text-stone-400">
        Don't have an account?
        <NuxtLink
          class="font-medium text-stone-600 hover:text-stone-800"
          to="/register"
        >
          Set up your account
        </NuxtLink>
      </p>
    </div>
  </div>
</template>
