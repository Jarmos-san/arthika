/**
 * @description Pinia store for global authentication state. In development mode
 * `isAuthenticated` is always `true` to avoid repeated logins during DX.
 */
export const useAuthStore = defineStore("auth", () => {
  /** @description Whether a user is currently authenticated. Always `true` in dev mode. */
  const isAuthenticated = ref<boolean>(import.meta.dev);

  return {
    isAuthenticated,
  };
});
