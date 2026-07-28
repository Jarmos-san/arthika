import getRedirectPath from "~/utils/routing";

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
  const { isAuthenticated } = useAuth();

  const redirectPath = getRedirectPath(isAuthenticated.value, to.path);

  if (redirectPath !== undefined) {
    return navigateTo(redirectPath);
  }
});
