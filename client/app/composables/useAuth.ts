/** A registered user as returned by the API. */
export interface User {
  id: string;
  email: string;
}

/** Payload for POST /api/users/register. */
export interface RegisterRequest {
  email: string;
  password: string;
}

/** Response from POST /api/users/register on success (201). */
export interface RegisterResponse {
  id: string;
  email: string;
}

/** Payload for POST /api/users/login. */
export interface LoginRequest {
  email: string;
  password: string;
}

/** Response from POST /api/users/login on success (200). */
export interface LoginResponse {
  token: string;
  id: string;
  email: string;
}

/** Response from GET /api/setup/status. */
export interface SystemStatusResponse {
  needsSetup: boolean;
}

/** A single field-level validation error returned by the API (422). */
export interface ValidationError {
  field: string;
  message: string;
}

/**
 * A structured error returned by API methods in this composable.
 *
 * `status` distinguishes the error category so callers can handle each case
 * appropriately (e.g. show inline validation errors vs. a generic banner).
 */
export interface ApiError {
  status: "validation" | "conflict" | "unauthorized" | "unknown";
  message: string;
  errors?: ValidationError[];
}

/**
 * Discriminated union for API call results.
 *
 * - `{ success: true; data: T }` — the call succeeded.
 * - `{ success: false; error: ApiError }` — the call failed with a structured
 *   error.
 *
 * Using a discriminated union forces consumers to check `success` before
 * accessing `data`, preventing runtime type errors.
 */
export type ApiResult<T> =
  | { success: true; data: T }
  | { success: false; error: ApiError };

/**
 * Authentication composable for Arthika.
 *
 * Provides reactive auth state and methods for registering, logging in,
 * checking setup status, and logging out. Uses `$fetch` for API calls and
 * `useCookie('token')` for JWT persistence.
 *
 * Auto-imported by Nuxt 4 — no manual import needed in `.vue` files.
 *
 * @returns An object containing:
 *   - `user` — reactive `User | null`, read-only.
 *   - `register(email, password)` — create a new account.
 *   - `login(email, password)` — authenticate and persist the session.
 *   - `checkSetupStatus()` — query whether the app needs first-time setup.
 *   - `logout()` — clear the session.
 *   - `isAuthenticated` — computed boolean derived from `user`.
 *
 * @example
 * ```ts
 * const { register, user, isAuthenticated } = useAuth()
 * const result = await register('a@b.com', 'password123')
 * if (result.success) { /* redirect to login *\/ }
 * ```
 */
export const useAuth = () => {
  const user = ref<User | null>(null);

  /**
   * Register a new user account.
   *
   * POSTs to /api/users/register with the provided credentials. On success the
   * returned user data is **not** automatically persisted — the caller should
   * redirect to /login so the user authenticates with their new account.
   *
   * @param email - User's email address.
   * @param password - User's password (minimum 8 characters).
   * @returns `ApiResult`:
   *   - success: the created user's id and email.
   *   - `conflict` (409): the email is already registered.
   *   - `validation` (422): the email or password failed validation.
   */
  const register = async (email: string, password: string) => {
    try {
      const data = await $fetch<RegisterResponse>("/api/users/register", {
        method: "POST",
        body: { email, password },
      });
      return { success: true, data };
    } catch (err: any) {
      if (err.status === 409) {
        return {
          success: false,
          error: {
            status: "conflict",
            message: err.data?.message ?? "Email already registered",
          },
        };
      }
      if (err.status === 422) {
        return {
          success: false,
          error: {
            status: "validation",
            message: "Validation failed",
            errors: err.data?.errors,
          },
        };
      }
      return {
        success: false,
        error: { status: "unknown", message: err.message },
      };
    }
  };

  /**
   * Authenticate an existing user.
   *
   * POSTs to /api/users/login with the provided credentials. On success the JWT
   * is persisted in a `token` cookie (via `useCookie`) and the user's reactive
   * state is updated.
   *
   * @param email - User's email address.
   * @param password - User's password.
   * @returns `ApiResult`:
   *   - success: the JWT token, user id, and email.
   *   - `unauthorized` (401): the email or password is incorrect.
   *   - `validation` (422): the email or password failed validation.
   */
  const login = async (email: string, password: string) => {
    try {
      const data = await $fetch<LoginResponse>("/api/users/login", {
        method: "POST",
        body: { email, password },
      });
      const token = useCookie("token");
      token.value = data.token;
      user.value = { id: data.id, email: data.email };
      return { success: true, data };
    } catch (err: any) {
      if (err.status === 401) {
        return {
          success: false,
          error: {
            status: "unauthorized",
            message: err.data?.message ?? "Invalid email or password",
          },
        };
      }
      if (err.status === 422) {
        return {
          success: false,
          error: {
            status: "validation",
            message: "Validation failed",
            errors: err.data?.errors,
          },
        };
      }
      return {
        success: false,
        error: { status: "unknown", message: err.message },
      };
    }
  };

  /**
   * Check whether the application needs first-time setup.
   *
   * GETs /api/setup/status (unauthenticated endpoint). Returns whether any
   * users exist in the database. Used by the auth route guard to decide whether
   * to show the registration or login page.
   *
   * @returns `ApiResult`:
   *   - success: `{ needsSetup: boolean }`.
   *   - `unknown`: the request failed (network error, server error).
   */
  const checkSetupStatus = async () => {
    try {
      const data = await $fetch<SystemStatusResponse>("/api/setup/status");
      return { success: true, data };
    } catch (err: any) {
      return {
        success: false,
        error: { status: "unknown", message: err.message },
      };
    }
  };

  /**
   * Log the current user out.
   *
   * Clears the `token` cookie (via `useCookie`) and resets the reactive `user`
   * state to `null`. The calling page is responsible for redirecting (e.g. to
   * /login or /setup).
   */
  const logout = () => {
    const token = useCookie("token");
    token.value = null;
    user.value = null;
  };

  /** Whether a user is currently logged in. Derived from `user !== null`. */
  const isAuthenticated = computed(() => user.value !== null);

  return {
    /** The currently authenticated user, or `null` if not logged in. */
    user: readonly(user),
    register,
    login,
    checkSetupStatus,
    logout,
    isAuthenticated,
  };
};
