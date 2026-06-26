import { FetchError } from "ofetch";

/** @description A registered user as returned by the API. */
interface User {
  id: string;
  email: string;
}

/** @description Payload for POST /api/users/register. */
interface RegisterRequest {
  email: string;
  password: string;
}

/** @description Response from POST /api/users/register on success (201). */
interface RegisterResponse {
  id: string;
  email: string;
}

/** @description Payload for POST /api/users/login. */
interface LoginRequest {
  email: string;
  password: string;
}

/** @description Response from POST /api/users/login on success (200). */
interface LoginResponse {
  token: string;
  id: string;
  email: string;
}

/** @description Response from GET /api/setup/status. */
interface SystemStatusResponse {
  needsSetup: boolean;
}

/** @description Shape of the JSON body returned by the API on error. */
interface ApiErrorBody {
  message?: string;
  errors?: ValidationError[];
}

/** @description A field-level validation error returned by the API (422). */
interface ValidationError {
  field: string;
  message: string;
}

/**
 * @description A structured error returned by API methods in this composable. `status`
 * distinguishes the error category so callers can handle each case
 * appropriately (e.g. show inline validation errors vs. a generic banner).
 */
interface ApiError {
  status: "validation" | "conflict" | "unauthorized" | "unknown";
  message: string;
  errors?: ValidationError[];
}

/**
 * @description Discriminated union for API call results. - `{ success: true; data: T }` —
 * the call succeeded. - `{ success: false; error: ApiError }` — the call failed
 * with a structured error. Using a discriminated union forces consumers to
 * check `success` before accessing `data`, preventing runtime type errors.
 */
type ApiResult<Data> =
  | { success: true; data: Data }
  | { success: false; error: ApiError };

/**
 * @description Authentication composable for Arthika. Provides reactive auth state and
 * methods for registering, logging in, checking setup status, and logging out.
 * Uses `$fetch` for API calls and `useCookie('token')` for JWT persistence.
 * Auto-imported by Nuxt 4 — no manual import needed in `.vue` files.
 *
 * @example
 *   ```ts
 *   const { register, user, isAuthenticated } = useAuth()
 *   const result = await register('a@b.com', 'password123')
 *   if (result.success) { /* redirect to login *\/ }
 *   ```;
 *
 * @returns {object} An object containing:
 *
 *   - `user` — reactive `User | undefined`, read-only.
 *   - `register(email, password)` — create a new account.
 *   - `login(email, password)` — authenticate and persist the session.
 *   - `checkSetupStatus()` — query whether the app needs first-time setup.
 *   - `logout()` — clear the session.
 *   - `isAuthenticated` — computed boolean derived from `user`.
 */
const useAuth = () => {
  const user = useState<User | undefined>("auth-user", () => undefined);

  /**
   * @description Register a new user account. POSTs to /api/users/register with the provided
   * credentials. On success the returned user data is **not** automatically
   * persisted — the caller should redirect to /login so the user authenticates
   * with their new account.
   *
   * @param {string} email - User's email address.
   * @param {string} password - User's password (minimum 8 characters).
   *
   * @returns {ApiResult} `ApiResult`:
   *
   *   - success: the created user's id and email.
   *   - `conflict` (409): the email is already registered.
   *   - `validation` (422): the email or password failed validation.
   */
  const register = async (email: string, password: string) => {
    try {
      const data = await $fetch<RegisterResponse>("/api/users/register", {
        body: { email, password },
        method: "POST",
      });
      return { data, success: true };
    } catch (error: unknown) {
      if (error instanceof FetchError) {
        // oxlint-disable-next-line typescript/no-unsafe-type-assertion
        const body = error.data as ApiErrorBody | undefined;
        if (error.status === 409) {
          return {
            error: {
              message: body?.message ?? "Email already registered",
              status: "conflict",
            },
            success: false,
          };
        }
        if (error.status === 422) {
          return {
            error: {
              errors: body?.errors,
              message: "Validation failed",
              status: "validation",
            },
            success: false,
          };
        }
      }
      let message: string;
      if (error instanceof Error) {
        ({ message } = error);
      } else {
        message = "An unexpected error occurred";
      }
      return {
        error: {
          message,
          status: "unknown",
        },
        success: false,
      };
    }
  };

  /**
   * @description Authenticate an existing user. POSTs to /api/users/login with the provided
   * credentials. On success the JWT is persisted in a `token` cookie (via
   * `useCookie`) and the user's reactive state is updated.
   *
   * @param {string} email - User's email address.
   * @param {string} password - User's password.
   *
   * @returns {ApiResult} `ApiResult`:
   *
   *   - success: the JWT token, user id, and email.
   *   - `unauthorized` (401): the email or password is incorrect.
   *   - `validation` (422): the email or password failed validation.
   */
  const login = async (email: string, password: string) => {
    try {
      const data = await $fetch<LoginResponse>("/api/users/login", {
        body: { email, password },
        method: "POST",
      });
      const token = useCookie("token");
      token.value = data.token;
      user.value = { email: data.email, id: data.id };
      return { data, success: true };
    } catch (error: unknown) {
      if (error instanceof FetchError) {
        // oxlint-disable-next-line typescript/no-unsafe-type-assertion
        const body = error.data as ApiErrorBody | undefined;
        if (error.status === 401) {
          return {
            error: {
              message: body?.message ?? "Invalid email or password",
              status: "unauthorized",
            },
            success: false,
          };
        }
        if (error.status === 422) {
          return {
            error: {
              errors: body?.errors,
              message: "Validation failed",
              status: "validation",
            },
            success: false,
          };
        }
      }
      let message: string;
      if (error instanceof Error) {
        ({ message } = error);
      } else {
        message = "An unexpected error occurred";
      }
      return {
        error: {
          message,
          status: "unknown",
        },
        success: false,
      };
    }
  };

  /**
   * @description Check whether the application needs first-time setup. GETs
   * /api/setup/status (unauthenticated endpoint). Returns whether any users
   * exist in the database. Used by the auth route guard to decide whether to
   * show the registration or login page.
   *
   * @returns {ApiResult} `ApiResult`:
   *
   *   - success: `{ needsSetup: boolean }`.
   *   - `unknown`: the request failed (network error, server error).
   */
  const checkSetupStatus = async () => {
    try {
      const data = await $fetch<SystemStatusResponse>("/api/setup/status");
      return { data, success: true };
    } catch (error: unknown) {
      let message: string;
      if (error instanceof Error) {
        ({ message } = error);
      } else {
        message = "An unexpected error occurred";
      }
      return {
        error: {
          message,
          status: "unknown",
        },
        success: false,
      };
    }
  };

  /**
   * @description Log the current user out. Clears the `token` cookie (via `useCookie`) and
   * resets the reactive `user` state to `undefined`. The calling page is
   * responsible for redirecting (e.g. to /login or /setup).
   */
  const logout = (): void => {
    const token = useCookie("token");
    // oxlint-disable-next-line unicorn/no-null
    token.value = null;
    user.value = undefined;
  };

  /** @description Whether a user is currently logged in. Derived from `user !== undefined`. */
  const isAuthenticated = computed(() => user.value !== undefined);

  return {
    checkSetupStatus,
    isAuthenticated,
    login,
    logout,
    register,
    /** @description The currently authenticated user, or `undefined` if not logged in. */
    user: readonly(user),
  };
};

export type {
  ApiError,
  ApiResult,
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  RegisterResponse,
  SystemStatusResponse,
  User,
  ValidationError,
};

export { useAuth };
