/**
 * @description Global route guard that enforces authentication. Unauthenticated users are
 * redirected to `/login` except on public routes (`/login`, `/setup`).
 * Authenticated users visiting `/login` are redirected to `/dashboard`.
 */
export default defineNuxtRouteMiddleware((to) => {
  const { isAuthenticated } = useAuthStore();

  const publicRoutes = new Set(["/login", "/setup"]);

  if (isAuthenticated && to.path === "/login") {
    return navigateTo("/dashboard");
  }

  if (!isAuthenticated && !publicRoutes.has(to.path)) {
    return navigateTo("/login");
  }
});
