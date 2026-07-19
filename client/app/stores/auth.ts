/** @description Pinia store for global authentication state. */
export const useAuthStore = defineStore("auth", () => {
  /** @description Whether a user is currently authenticated. */
  const isAuthenticated = ref<boolean>(false);

  /**
   * @description Authenticates the user. Stub — will be replaced with an actual API call
   * later.
   */
  const login = (email: string, _password: string): void => {
    console.log("Logging in with:", email);
    isAuthenticated.value = true;
  };

  return { isAuthenticated, login };
});
