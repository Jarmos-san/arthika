<script setup lang="ts">
  import { useHead } from "nuxt/app";
  import { ref } from "vue";

  import useToast from "~/composables/UseToast";
  import type { RegisterErrors } from "~/types/utils/validators";
  import { validateRegister } from "~/utils/validators";

  useHead({
    title: "Register",
  });

  const email = ref<string | undefined>(undefined);
  const password = ref<string | undefined>(undefined);
  const confirmPassword = ref<string | undefined>(undefined);

  const errors = ref<RegisterErrors>({
    confirmPassword: undefined,
    email: undefined,
    password: undefined,
  });

  const toast = useToast();

  const onSubmit = (): void => {
    errors.value = validateRegister(
      email.value,
      password.value,
      confirmPassword.value,
    );

    if (
      errors.value.email ||
      errors.value.password ||
      errors.value.confirmPassword
    ) {
      return;
    }

    toast.publish("Account created successfully!");
  };
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

      <!-- The registration form -->
      <form class="flex flex-col gap-5" @submit.prevent="onSubmit">
        <!-- Email input field -->
        <div>
          <Input
            id="email"
            v-model="email"
            aria-describedby="email-error"
            label="Email"
            placeholder="you@example.com"
            type="email"
          />
          <p
            v-if="errors.email"
            id="email-error"
            class="mt-1.5 text-xs text-red-600"
          >
            {{ errors.email }}
          </p>
        </div>

        <!-- Password input field -->
        <div>
          <Input
            id="password"
            v-model="password"
            aria-describedby="password-error"
            label="Password"
            placeholder="Minimum 8 characters"
            type="password"
          />
          <p
            v-if="errors.password"
            id="password-error"
            class="mt-1.5 text-xs text-red-600"
          >
            {{ errors.password }}
          </p>
        </div>

        <!-- Password confirmation field -->
        <div>
          <Input
            id="confirm-password"
            v-model="confirmPassword"
            aria-describedby="confirm-password-error"
            label="Confirm Password"
            placeholder="Retype your password"
            type="password"
          />
          <p
            v-if="errors.confirmPassword"
            id="confirm-password-error"
            class="mt-1.5 text-xs text-red-600"
          >
            {{ errors.confirmPassword }}
          </p>
        </div>

        <!-- Submission button -->
        <button class="btn-primary mt-1" type="submit">Create Account</button>
      </form>
      <p class="mt-6 text-center text-sm text-stone-400">
        Have an account?
        <NuxtLink
          class="font-medium text-stone-600 hover:text-stone-800"
          to="/login"
        >
          Login in instead
        </NuxtLink>
      </p>
    </div>
  </div>
</template>
