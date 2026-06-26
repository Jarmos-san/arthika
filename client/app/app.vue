<script setup lang="ts">
  /**
   * @description Index of the JWT payload segment when splitting by `.`. JWT format is
   * `header.payload.signature`.
   */
  const JWT_PAYLOAD_INDEX = 1;

  const token = useCookie("token");
  if (token.value) {
    const user = useState<User | undefined>("auth-user");
    try {
      const segment = token.value.split(".")[JWT_PAYLOAD_INDEX];
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
