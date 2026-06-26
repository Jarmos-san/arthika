<script setup lang="ts">
  const token = useCookie("token");
  if (token.value) {
    const user = useState<User | undefined>("auth-user");
    try {
      const segment = token.value.split(".")[1];
      if (segment) {
        const payload = JSON.parse(atob(segment));
        if (
          typeof payload.sub === "string" &&
          typeof payload.email === "string"
        ) {
          user.value = { email: payload.email, id: payload.sub };
        } else {
          // oxlint-disable-next-line unicorn/no-null
          token.value = null;
        }
      } else {
        // oxlint-disable-next-line unicorn/no-null
        token.value = null;
      }
    } catch {
      // oxlint-disable-next-line unicorn/no-null
      token.value = null;
    }
  }
</script>

<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>
