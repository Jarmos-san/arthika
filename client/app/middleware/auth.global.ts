import { computed } from "vue";

import { useRequestHeaders, useState } from "#app";
import type { CurrentUserResponse } from "~/openapi";
import getRedirectPath from "~/utils/routing";

/** @description Shared state key for the authenticated user (matches useAuth). */
const USER_STATE_KEY = "auth-user";

/**
 * @description Global route middleware that enforces authentication boundaries. Redirects
 * authenticated users away from public routes (e.g. `/login`) and
 * unauthenticated users away from protected routes to `/login`. Navigation
 * logic is delegated to `getRedirectPath`.
 *
 * @param {RouteLocationNormalized} to - The target route being navigated to.
 *
 * @returns {void | ReturnType<typeof navigateTo>} `void` when no redirect is
 *   needed, or the result of `navigateTo` to perform the redirect.
 */
export default defineNuxtRouteMiddleware(async (to) => {
  // Initialise the state of user when the middleware is fired
  const user = useState<CurrentUserResponse | undefined>(USER_STATE_KEY);

  // Check the authentication state of the user
  if (user.value === undefined) {
    try {
      const data = await $fetch<CurrentUserResponse>(
        "/api/users/current-user",
        { headers: useRequestHeaders(["cookie"]) },
      );
      user.value = data;
    } catch {
      user.value = undefined;
    }
  }

  // Boolean flag to set the authentication state of the user
  const isAuthenticated = computed(() => user.value !== undefined);

  // Create the redirection path based on the authentication state
  const redirectPath = getRedirectPath(isAuthenticated.value, to.path);

  // Redirect the unauthenticated user to the login page
  if (redirectPath !== undefined) {
    return navigateTo(redirectPath);
  }
});
