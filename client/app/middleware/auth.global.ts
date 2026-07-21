import useAuthStore from "~/stores/auth";
import getRedirectPath from "~/utils/routing";

export default defineNuxtRouteMiddleware(async (to) => {
  const { isAuthenticated } = useAuthStore();
  const redirectPath = getRedirectPath(isAuthenticated, to.path);

  if (redirectPath !== undefined) {
    await navigateTo(redirectPath);
  }
});
