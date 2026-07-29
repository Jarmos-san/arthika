import { type ComputedRef, computed } from "vue";

import { useAsyncData, useState } from "#app";
import type {
  CurrentUserResponse,
  ErrorResponse,
  ValidationErrors,
} from "~/openapi";

/** @description HTTP 200 OK. */
const STATUS_OK = 200;

/** @description HTTP 201 Created. */
const STATUS_CREATED = 201;

/** @description HTTP 401 Unauthorized. */
const STATUS_UNAUTHORIZED = 401;

/** @description HTTP 409 Conflict. */
const STATUS_CONFLICT = 409;

/** @description HTTP 422 Unprocessable Entity. */
const STATUS_UNPROCESSABLE = 422;

/** @description Returned when no HTTP status is available. */
const STATUS_NETWORK_ERROR = 0;

/**
 * @description Extracts the HTTP status code from a `$fetch` error. Returns
 * {@link STATUS_NETWORK_ERROR} when no status is available.
 *
 * @param {unknown} error - The error thrown by `$fetch`.
 *
 * @returns {number} The HTTP status code, or `0` for network errors.
 */
const getErrorStatus = (error: unknown): number => {
  if (typeof error === "object" && error !== null && "status" in error) {
    const { status } = error as { status: unknown };
    if (typeof status === "number") {
      return status;
    }
  }

  return STATUS_NETWORK_ERROR;
};

/**
 * @description Type guard that checks whether a value matches the `ErrorResponse` shape (`{
 * message: string }`).
 *
 * @param {unknown} value - The value to check.
 *
 * @returns {boolean} `true` if the value is an `ErrorResponse`.
 */
const isErrorResponse = (value: unknown): value is ErrorResponse =>
  typeof value === "object" &&
  value !== null &&
  "message" in value &&
  typeof (value as { message: unknown }).message === "string";

/**
 * @description Type guard that checks whether a value matches the `ValidationErrors` shape
 * (`{ errors: Array }`).
 *
 * @param {unknown} value - The value to check.
 *
 * @returns {boolean} `true` if the value is a `ValidationErrors`.
 */
const isValidationErrors = (value: unknown): value is ValidationErrors =>
  typeof value === "object" &&
  value !== null &&
  "errors" in value &&
  Array.isArray((value as { errors: unknown }).errors);

/**
 * @description Extracts the parsed error body from a `$fetch` error. Returns `undefined` for
 * network errors where no response body is available.
 *
 * @param {unknown} error - The error thrown by `$fetch`.
 *
 * @returns {ErrorResponse | ValidationErrors | undefined} The parsed error
 *   body, or `undefined`.
 */
const getErrorBody = (
  error: unknown,
): ErrorResponse | ValidationErrors | undefined => {
  if (typeof error === "object" && error !== null && "data" in error) {
    const { data } = error as { data: unknown };
    if (isErrorResponse(data)) {
      return data;
    }

    if (isValidationErrors(data)) {
      return data;
    }
  }

  return undefined;
};

/** @description Result returned by {@link UseAuth.login} and {@link UseAuth.register}. */
interface ApiResult<SuccessData = void> {
  /** @description Response body on success. */
  data?: SuccessData;

  /** @description Whether the request succeeded. */
  ok: boolean;

  /** @description Parsed error body from the server. */
  body?: ErrorResponse | ValidationErrors;

  /** @description HTTP status code, or `0` for network errors. */
  status: number;
}

/** @description Return type for the {@link useAuth} composable. */
interface UseAuth {
  /** @description Whether a user is currently authenticated. */
  isAuthenticated: ComputedRef<boolean>;

  /** @description Authenticates the user with email and password. */
  login: (email: string, password: string) => Promise<ApiResult>;

  /** @description Registers a new user with email and password. */
  register: (email: string, password: string) => Promise<ApiResult>;
}

/** @description Nuxt state key used to persist the authenticated user across SSR. */
const USER_STATE_KEY = "auth-user";

/**
 * @description Composable for managing authentication state. Uses `useState` to persist the
 * user across the SSR boundary and `useAsyncData` to fetch the current session
 * on mount. Provides `login` and `register` methods that update the shared
 * state and return structured results.
 *
 * @returns {UseAuth} An object with `isAuthenticated`, `login`, and `register`.
 */
const useAuth = (): UseAuth => {
  const user = useState<CurrentUserResponse | undefined>(USER_STATE_KEY);

  void useAsyncData(`${USER_STATE_KEY}:session`, async () => {
    try {
      const data = await $fetch<CurrentUserResponse>("/api/users/current-user");
      user.value = data;
    } catch {
      user.value = undefined;
    }

    // oxlint-disable-next-line unicorn/no-null
    return null;
  });

  const isAuthenticated = computed(() => user.value !== undefined);

  const submitAuth = async (
    url: string,
    credentials: Readonly<{ email: string; password: string }>,
    successStatus: number,
  ): Promise<ApiResult> => {
    try {
      const data = await $fetch<CurrentUserResponse>(url, {
        body: { email: credentials.email, password: credentials.password },
        method: "POST",
      });

      user.value = data;

      return { ok: true, status: successStatus };
    } catch (error: unknown) {
      return {
        body: getErrorBody(error),
        ok: false,
        status: getErrorStatus(error),
      };
    }
  };

  const login = async (email: string, password: string): Promise<ApiResult> => {
    const result = await submitAuth(
      "/api/users/login",
      { email, password },
      STATUS_OK,
    );
    return result;
  };

  const register = async (
    email: string,
    password: string,
  ): Promise<ApiResult> => {
    const result = await submitAuth(
      "/api/users/register",
      { email, password },
      STATUS_CREATED,
    );
    return result;
  };

  return { isAuthenticated, login, register };
};

export type { ApiResult };
export {
  STATUS_CONFLICT,
  STATUS_NETWORK_ERROR,
  STATUS_OK,
  STATUS_UNAUTHORIZED,
  STATUS_UNPROCESSABLE,
};
export default useAuth;
