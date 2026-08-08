/** @description Routes accessible without authentication. */
const publicRoutes = new Set(["/login", "/register"]);

/**
 * @description Determines the redirect path based on authentication state and the target
 * route. Authenticated users visiting `/login` are sent to `/dashboard`.
 * Unauthenticated users on protected routes are sent to `/login`.
 *
 * @param {boolean} isAuthenticated - Whether the user is currently
 *   authenticated.
 * @param {string} targetPath - The path of the route being navigated to.
 *
 * @returns {string | undefined} The path to redirect to, or `undefined` if no
 *   redirect is needed.
 */
const getRedirectPath = (
  isAuthenticated: boolean,
  targetPath: string,
): string | undefined => {
  if (targetPath === "/") {
    return "/dashboard";
  }

  if (isAuthenticated && targetPath === "/login") {
    return "/dashboard";
  }

  if (!isAuthenticated && !publicRoutes.has(targetPath)) {
    return "/login";
  }

  return undefined;
};

export default getRedirectPath;
