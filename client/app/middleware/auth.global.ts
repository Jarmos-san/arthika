import useAuthStore from "~/stores/auth";

/**
 * @description Global route guard that enforces authentication. Unauthenticated users are
 * redirected to `/login` except on public routes (`/login`, `/setup`).
 * Authenticated users visiting `/login` are redirected to `/dashboard`.
 *
 * @param {RouteLocation} to - The target route location being navigated to.
 */
export default defineNuxtRouteMiddleware(async (to) => {
  const { isAuthenticated } = useAuthStore();

  const publicRoutes = new Set(["/login", "/setup"]);

  if (isAuthenticated && to.path === "/login") {
    await navigateTo("/dashboard");
  }

  if (!isAuthenticated && !publicRoutes.has(to.path)) {
    await navigateTo("/login");
  }
});
