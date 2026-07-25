<script setup lang="ts">
  import { useHead } from "nuxt/app";
  import { ref } from "vue";

  import useAuthStore from "~/stores/auth";
  import type { RegisterErrors } from "~/types/utils/validators";
  import { validateRegister } from "~/utils/validators";

  useHead({ title: "Register" });

  /** @description HTTP status code for Conflict (email already registered). */
  const STATUS_CONFLICT = 409;

  /** @description HTTP status code for Unprocessable Entity (validation errors). */
  const STATUS_UNPROCESSABLE = 422;

  const email = ref<string | undefined>(undefined);
  const password = ref<string | undefined>(undefined);
  const confirmPassword = ref<string | undefined>(undefined);

  const errors = ref<RegisterErrors>({
    confirmPassword: undefined,
    email: undefined,
    password: undefined,
    server: undefined,
  });

  /**
   * @description Maps server-side errors from the register API response onto the form's
   * field-level error state.
   *
   * @param {number} status - The HTTP status code from the failed response.
   * @param {string | undefined} message - The error message from the server.
   * @param {{ field: string; message: string }[] | undefined} fieldErrors -
   *   Field-level validation errors from the server.
   */
  const applyServerErrors = (
    status: number,
    message: string | undefined,
    fieldErrors: { field: string; message: string }[] | undefined,
  ): void => {
    if (status === STATUS_CONFLICT) {
      errors.value.server = message;

      return;
    }

    if (status === STATUS_UNPROCESSABLE && fieldErrors) {
      for (const err of fieldErrors) {
        if (err.field === "email") {
          errors.value.email = err.message;
        }

        if (err.field === "password") {
          errors.value.password = err.message;
        }
      }
    }
  };

  const onSubmit = async (): Promise<void> => {
    errors.value = validateRegister(
      email.value,
      password.value,
      confirmPassword.value,
    );

    if (
      errors.value.email ||
      errors.value.password ||
      errors.value.confirmPassword ||
      email.value === undefined ||
      password.value === undefined
    ) {
      return;
    }

    const { register } = useAuthStore();
    const result = await register(email.value, password.value);

    if (result.ok) {
      await navigateTo("/dashboard");

      return;
    }

    applyServerErrors(result.status, result.message, result.errors);
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

      <form class="flex flex-col gap-5" @submit.prevent="onSubmit">
        <p
          v-if="errors.server"
          class="rounded-lg bg-red-50 p-3 text-sm text-red-600"
        >
          {{ errors.server }}
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
            placeholder="Minimum 8 characters"
          />
          <p
            v-if="errors.password"
            id="password-error"
            class="mt-1.5 text-xs text-red-600"
          >
            {{ errors.password }}
          </p>
        </div>

        <div>
          <Input
            id="confirm-password"
            v-model="confirmPassword"
            aria-describedby="confirm-password-error"
            label="Confirm Password"
            type="password"
            placeholder="Re-enter your password"
          />
          <p
            v-if="errors.confirmPassword"
            id="confirm-password-error"
            class="mt-1.5 text-xs text-red-600"
          >
            {{ errors.confirmPassword }}
          </p>
        </div>

        <button type="submit" class="btn-primary mt-1">Create Account</button>
      </form>

      <p class="mt-6 text-center text-sm text-stone-400">
        Already have an account?
        <NuxtLink
          class="font-medium text-stone-600 hover:text-stone-800"
          to="/login"
        >
          Log in
        </NuxtLink>
      </p>
    </div>
  </div>
</template>
