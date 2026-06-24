<script setup lang="ts">
const token = useCookie("token");
if (token.value) {
  const user = useState<User | null>("auth-user", () => null);
  try {
    const segment = token.value.split(".")[1];
    if (segment) {
      const payload = JSON.parse(atob(segment));
      if (typeof payload.sub === "string" && typeof payload.email === "string") {
        user.value = { id: payload.sub, email: payload.email };
      } else {
        token.value = null;
      }
    } else {
      token.value = null;
    }
  } catch {
    token.value = null;
  }
}
</script>

<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>
