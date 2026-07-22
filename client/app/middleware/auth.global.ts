import useAuthStore from "~/stores/auth";
import getRedirectPath from "~/utils/routing";

// oxlint-disable-next-line eslint/require-await
export default defineNuxtRouteMiddleware(async (to) => {
  const { isAuthenticated } = useAuthStore();
  const redirectPath = getRedirectPath(isAuthenticated, to.path);

  if (redirectPath !== undefined) {
    return navigateTo(redirectPath);
  }
});
