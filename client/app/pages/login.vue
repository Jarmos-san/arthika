<script setup lang="ts">
  useHead({ title: "Login" });

  const onSubmit = (): void => {};
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
        class="mb-5 hidden rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700"
      >
        Invalid email or password.
      </div>

      <form class="flex flex-col gap-5" @submit.prevent="onSubmit">
        <div>
          <Input
            id="email"
            label="Email"
            type="email"
            placeholder="you@example.com"
          />
          <p class="mt-1.5 hidden text-xs text-red-600">
            Please enter a valid email address.
          </p>
        </div>

        <div>
          <Input
            id="password"
            label="Password"
            type="password"
            placeholder="Enter your password"
          />
          <p class="mt-1.5 hidden text-xs text-red-600">
            Password is required.
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
