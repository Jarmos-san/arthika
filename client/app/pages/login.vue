<script setup lang="ts">
  import { useHead } from "nuxt/app";
  import { ref } from "vue";

  import useAuth, {
    STATUS_NETWORK_ERROR,
    STATUS_UNAUTHORIZED,
    STATUS_UNPROCESSABLE,
  } from "~/composables/useAuth";
  import useToast from "~/composables/useToast";
  import type { LoginErrors } from "~/types/utils/validators";
  import { validateLogin } from "~/utils/validators";

  // Add a title for the page
  useHead({ title: "Login" });

  // Initialise the state of the login form
  const email = ref<string | undefined>(undefined);
  const password = ref<string | undefined>(undefined);
  const isSubmitting = ref(false);

  // Initialise the errors to be rendered on the login form
  const errors = ref<LoginErrors>({
    email: undefined,
    password: undefined,
  });

  // Create an instance of the "toast" component from Reka UI
  const toast = useToast();

  // Utility function to render error message and update the error state
  const handleLoginError = (status: number, body: unknown): boolean => {
    if (status === STATUS_UNAUTHORIZED) {
      toast.publish("Invalid email or password");
      return true;
    }

    if (
      status === STATUS_UNPROCESSABLE &&
      body &&
      typeof body === "object" &&
      "errors" in body
    ) {
      for (const validationError of (
        body as { errors: { field: string; message: string }[] }
      ).errors) {
        if (validationError.field in errors.value) {
          errors.value[validationError.field as keyof LoginErrors] =
            validationError.message;
        }
      }
      return true;
    }

    return false;
  };

  // Utility wrapper to check if the user credentials are valid.
  const submitLogin = async (): Promise<void> => {
    const auth = useAuth();
    const emailValue = email.value ?? "";
    const passwordValue = password.value ?? "";
    const result = await auth.login(emailValue, passwordValue);

    if (result.ok) {
      await navigateTo("/dashboard");
      return;
    }

    if (
      !handleLoginError(result.status, result.body) &&
      result.status === STATUS_NETWORK_ERROR
    ) {
      toast.publish("Something went wrong. Please try again.");
    }
  };

  // Event handler for the submit button
  const onSubmit = async (): Promise<void> => {
    errors.value = validateLogin(email.value, password.value);

    if (errors.value.email || errors.value.password) {
      return;
    }

    isSubmitting.value = true;
    try {
      await submitLogin();
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
          to="/register"
        >
          Set up your account
        </NuxtLink>
      </p>
    </div>
  </div>
</template>
