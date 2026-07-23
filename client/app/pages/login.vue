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

  const serverError = ref<string | undefined>(undefined);
  const isSubmitting = ref<boolean>(false);

  const handleSubmit = async (
    loginEmail: string,
    loginPassword: string,
  ): Promise<void> => {
    const { login } = useAuthStore();
    await login(loginEmail, loginPassword);
    await navigateTo("/dashboard");
  };

  const formatError = (error: unknown): string => {
    const message =
      error instanceof Error ? error.message : "An unexpected error occurred";

    if (
      message.includes("401") ||
      message.toLowerCase().includes("unauthorized")
    ) {
      return "Invalid email or password";
    }

    return "Something went wrong. Please try again.";
  };

  const onSubmit = async (): Promise<void> => {
    serverError.value = undefined;
    errors.value = validateLogin(email.value, password.value);

    if (errors.value.email || errors.value.password) {
      return;
    }

    if (email.value === undefined || password.value === undefined) {
      return;
    }

    isSubmitting.value = true;

    try {
      await handleSubmit(email.value, password.value);
    } catch (error: unknown) {
      serverError.value = formatError(error);
    } finally {
      isSubmitting.value = false;
    }
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
        <p
          v-if="serverError"
          class="rounded-md bg-red-50 p-3 text-sm text-red-600"
        >
          {{ serverError }}
        </p>

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

        <button class="btn-primary mt-1" :disabled="isSubmitting" type="submit">
          <span v-if="!isSubmitting">Log In</span>
          <span v-else>Logging in...</span>
        </button>
      </form>

      <p class="mt-6 text-center text-sm text-stone-400">
        Don't have an account?
        <NuxtLink
          class="font-medium text-stone-600 hover:text-stone-800"
          to="/setup"
        >
          Set up your account
        </NuxtLink>
      </p>
    </div>
  </div>
</template>
