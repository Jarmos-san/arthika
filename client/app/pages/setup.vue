<script setup lang="ts">
/**
 * Setup page — first-time admin account creation.
 *
 * This page is shown when the application has no users and requires initial
 * setup. It presents a registration form (email + password) that calls the
 * `/api/users/register` endpoint via the `useAuth` composable.
 *
 * On success the UI transitions to a success state and auto-redirects to the
 * login page after a brief delay.
 *
 * @see {@link useAuth#register} for the underlying API call and error handling.
 * @see /api/openapi.yml for the POST /api/users/register endpoint contract.
 */

/** Tracks the current UI step: registration form or post-submit success. */
const step = ref<"form" | "success">("form");

/** Bound form model for the email and password fields. */
const state = reactive({ email: "", password: "" });

/** Error message returned from the server on a failed registration attempt. */
const serverError = ref("");

/** Whether the form is currently being submitted (disables the button). */
const submitting = ref(false);

const { register } = useAuth();

/**
 * Client-side form validation for the setup form.
 *
 * Checks that the email is non-empty and matches a basic email regex, and that
 * the password is non-empty and at least 8 characters long. Matches the
 * validation rules expected by `UForm`.
 *
 * @param state - The current form state.
 * @returns An array of error objects, each with a `name` (field name) and
 *          `message`. An empty array means the form is valid.
 */
const validate = (state: { email: string; password: string }) => {
  const errors: { name: string; message: string }[] = [];
  if (!state.email) {
    errors.push({ name: "email", message: "Email is required" });
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(state.email)) {
    errors.push({ name: "email", message: "Invalid email format" });
  }
  if (!state.password) {
    errors.push({ name: "password", message: "Password is required" });
  } else if (state.password.length < 8) {
    errors.push({
      name: "password",
      message: "Password must be at least 8 characters",
    });
  }
  return errors;
};

/**
 * Handle form submission.
 *
 * Calls `useAuth().register()` with the current email and password. On failure
 * the server error message is displayed below the form. On success the UI
 * switches to the success state and auto-redirects to /login after 2 seconds.
 */
const onSubmit = async () => {
  submitting.value = true;
  serverError.value = "";

  const result = await register(state.email, state.password);

  if (!result.success) {
    serverError.value = (result.error! as ApiError).message;
    submitting.value = false;
    return;
  }

  step.value = "success";
  setTimeout(() => navigateTo("/login"), 2000);
  submitting.value = false;

  submitting.value = false;
};
</script>

<template>
  <div class="flex items-center justify-center min-h-screen p-4">
    <UCard v-if="step === 'form'" class="w-full max-w-md">
      <template #header>
        <div class="text-center">
          <h1 class="text-2xl font-bold">Set up your account</h1>
          <p class="text-sm text-muted mt-1">
            Create the first user to get started with Arthika
          </p>
        </div>
      </template>

      <UForm :state="state" :validate="validate" @submit="onSubmit">
        <UFormField label="Email" name="email" required>
          <UInput
            v-model="state.email"
            type="email"
            placeholder="you@example.com"
          />
        </UFormField>

        <UFormField label="Password" name="password" required>
          <UInput
            v-model="state.password"
            type="password"
            placeholder="Minimum 8 characters"
          />
        </UFormField>

        <p class="text-xs text-muted mt-3">
          Your password is securely stored in a local database and it will never
          be shared with anyone else.
        </p>

        <UButton type="submit" block :loading="submitting" class="mt-4">
          Create account
        </UButton>
      </UForm>

      <p v-if="serverError" class="mt-4 text-sm text-error">
        {{ serverError }}
      </p>
    </UCard>

    <div v-else class="text-center">
      <div
        class="w-16 h-16 mx-auto mb-4 rounded-full bg-success/10 text-success flex items-center justify-center animate-success"
      >
        <UIcon name="i-lucide-check" class="w-8 h-8" />
      </div>
      <h2 class="text-2xl font-bold">Account created!</h2>
      <p class="text-muted mt-2">Redirecting to login&hellip;</p>
    </div>
  </div>
</template>

<style scoped>
.animate-success {
  animation: pop 0.4s ease-out;
}

@keyframes pop {
  0% {
    transform: scale(0);
  }
  50% {
    transform: scale(1.15);
  }
  100% {
    transform: scale(1);
  }
}
</style>
