import { defineStore } from "pinia";
import { ref } from "vue";

import type { LoginResponse } from "~/openapi/models/LoginResponse";
import type { UserResponse } from "~/openapi/models/UserResponse";

/** @description Pinia store for global authentication state. */
const useAuthStore = defineStore("auth", () => {
  /** @description Whether a user is currently authenticated. */
  const isAuthenticated = ref<boolean>(false);

  /** @description The authenticated user's profile, or undefined if not authenticated. */
  const user = ref<{ id: string; email: string } | undefined>(undefined);

  /**
   * @description Authenticates the user by calling the login API. On success the server sets
   * an HttpOnly session cookie and returns the user's profile.
   *
   * @param {string} email The user's email address.
   * @param {string} password The user's password.
   * @throws {Error} If the API returns a non-2xx response.
   */
  const login = async (email: string, password: string): Promise<void> => {
    const response = await $fetch<LoginResponse>("/api/users/login", {
      body: { email, password },
      method: "POST",
    });

    user.value = { email: response.email, id: response.id };
    isAuthenticated.value = true;
  };

  /**
   * @description Restores authentication state by calling the /me endpoint. The browser
   * sends the session cookie automatically. On success the user profile is
   * stored; on failure the state is cleared.
   */
  const fetchUser = async (): Promise<void> => {
    try {
      const response = await $fetch<UserResponse>("/api/users/me");
      user.value = { email: response.email, id: response.id };
      isAuthenticated.value = true;
    } catch {
      user.value = undefined;
      isAuthenticated.value = false;
    }
  };

  return { fetchUser, isAuthenticated, login, user };
});

export default useAuthStore;
