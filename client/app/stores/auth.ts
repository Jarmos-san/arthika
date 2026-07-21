import { defineStore } from "pinia";
import { ref } from "vue";

/** @description Pinia store for global authentication state. */
const useAuthStore = defineStore("auth", () => {
  /** @description Whether a user is currently authenticated. */
  const isAuthenticated = ref<boolean>(false);

  /**
   * @description Authenticates the user. Stub — will be replaced with an actual API call
   * later.
   *
   * @param {string} _email The user's email address to be used for
   *   authentication.
   * @param{string} _password The user's password to be used for authentication.
   */
  const login = (_email: string, _password: string): void => {
    isAuthenticated.value = true;
  };

  return { isAuthenticated, login };
});

export default useAuthStore;
