<script setup lang="ts">
  useHead({ title: "Login" });

  const email = ref<string | undefined>(undefined);
  const password = ref<string | undefined>(undefined);

  /**
   * @description Client-side validation errors for the login form. `undefined` means no
   * error.
   */
  interface LoginErrors {
    email: string | undefined;
    password: string | undefined;
    general: string | undefined;
  }

  /** @description Reactive container for current validation errors. Reset on each submit. */
  const errors = ref<LoginErrors>({
    email: undefined,
    password: undefined,
    general: undefined,
  });

  /**
   * @description Validates email and password fields against the same rules the server
   * enforces (valid email format, password minimum 8 characters). Resets errors
   * before checking and returns `true` only when all fields pass.
   *
   * @returns `true` if the form is valid, `false` otherwise.
   */
  const validate = (): boolean => {
    errors.value = {
      email: undefined,
      password: undefined,
      general: undefined,
    };

    // Email field validation
    if (email.value === undefined) {
      errors.value.email = "Email is required";
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.value)) {
      errors.value.email = "Please enter a valid email address";
    }

    // Password field validation
    if (password.value === undefined) {
      errors.value.password = "Password is required";
    } else if (password.value.length < 8) {
      errors.value.password = "Password must be at least 8 characters";
    }

    // Return true/false based on validation logic above
    return (
      errors.value.email === undefined && errors.value.password === undefined
    );
  };

  /** @description Form submit handler. Validates inputs then navigates to `/dashboard`. */
  const onSubmit = async (): Promise<void> => {
    if (!validate()) {
      return;
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
      <div
        v-if="errors.general"
        class="mb-5 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700"
      >
        {{ errors.general }}
      </div>

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
          to="/setup"
        >
          Set up your account
        </NuxtLink>
      </p>
    </div>
  </div>
</template>
