import { defineStore } from "pinia";
import { ref } from "vue";

import type { RegisterResponse, ValidationErrors } from "~/openapi";

/** @description HTTP status code for Conflict (email already registered). */
const STATUS_CONFLICT = 409;

/** @description HTTP status code for Unprocessable Entity (validation errors). */
const STATUS_UNPROCESSABLE = 422;

/** @description Result of a successful registration attempt. */
interface RegisterResult {
  data: RegisterResponse;
  ok: true;
}

/** @description Result of a failed registration attempt. */
interface RegisterError {
  errors: ValidationErrors["errors"];
  message: string;
  ok: false;
  status: number;
}

/** @description Shape of errors thrown by ofetch on non-2xx responses. */
interface OfetchError {
  data?: unknown;
  status: number;
}

/**
 * @description Checks whether a caught value is an ofetch FetchError with a numeric status.
 *
 * @param {unknown} error - The caught value.
 *
 * @returns {boolean} `true` if the error has a numeric `status` property.
 */
const isOfetchError = (error: unknown): error is OfetchError =>
  typeof error === "object" &&
  error !== null &&
  "status" in error &&
  typeof (error as Record<string, unknown>).status === "number";

/**
 * @description Checks whether a value is an object with a string `message` property.
 *
 * @param {unknown} data - The value to check.
 *
 * @returns {boolean} `true` if the value has a string `message` property.
 */
const hasStringMessage = (data: unknown): data is { message: string } =>
  typeof data === "object" &&
  data !== null &&
  "message" in data &&
  typeof (data as Record<string, unknown>).message === "string";

/**
 * @description Checks whether a value is an object with an `errors` array property.
 *
 * @param {unknown} data - The value to check.
 *
 * @returns {boolean} `true` if the value has an `errors` array property.
 */
const hasErrorsArray = (
  data: unknown,
): data is { errors: ValidationErrors["errors"] } =>
  typeof data === "object" &&
  data !== null &&
  "errors" in data &&
  Array.isArray((data as Record<string, unknown>).errors);

/**
 * @description Extracts the error message from a 409 Conflict response body.
 *
 * @param {unknown} data - The parsed response body.
 *
 * @returns {string} The error message.
 */
const parseConflictMessage = (data: unknown): string => {
  if (hasStringMessage(data)) {
    return data.message;
  }

  return "email already registered";
};

/**
 * @description Extracts field-level errors from a 422 response body.
 *
 * @param {unknown} data - The parsed response body.
 *
 * @returns {ValidationErrors["errors"]} The field-level errors.
 */
const parseValidationErrors = (data: unknown): ValidationErrors["errors"] => {
  if (hasErrorsArray(data)) {
    return data.errors;
  }

  return [];
};

/**
 * @description Maps an ofetch error to a structured RegisterError.
 *
 * @param {number} status - The HTTP status code from the failed response.
 * @param {unknown} data - The parsed response body.
 *
 * @returns {RegisterError} A structured error object.
 */
const handleRegisterError = (status: number, data: unknown): RegisterError => {
  if (status === STATUS_CONFLICT) {
    return {
      errors: [],
      message: parseConflictMessage(data),
      ok: false,
      status: STATUS_CONFLICT,
    };
  }

  if (status === STATUS_UNPROCESSABLE) {
    return {
      errors: parseValidationErrors(data),
      message: "Validation failed",
      ok: false,
      status: STATUS_UNPROCESSABLE,
    };
  }

  return {
    errors: [],
    message: "An unexpected error occurred",
    ok: false,
    status,
  };
};

/**
 * @description Sends a registration request to the server API.
 *
 * @param {string} email - The user's email address.
 * @param {string} password - The user's password.
 * @param {string} apiBase - The base URL of the API server.
 *
 * @returns {Promise<RegisterResponse>} The server response on success.
 */
const sendRegisterRequest = async (
  email: string,
  password: string,
  apiBase: string,
): Promise<RegisterResponse> => {
  const response = await $fetch<RegisterResponse>("/api/users/register", {
    baseURL: apiBase,
    body: { email, password },
    credentials: "include",
    method: "POST",
  });

  return response;
};

/** @description Pinia store for global authentication state. */
const useAuthStore = defineStore("auth", () => {
  /** @description Whether a user is currently authenticated. */
  const isAuthenticated = ref<boolean>(false);

  /**
   * @description Authenticates the user. Stub — will be replaced with an actual API call
   * later.
   *
   * @param {string} _email The user's email address to be used for
   *   authentication.
   * @param {string} _password The user's password to be used for
   *   authentication.
   */
  const login = (_email: string, _password: string): void => {
    isAuthenticated.value = true;
  };

  /**
   * @description Registers a new user by calling the server API. On success the server sets
   * an HttpOnly auth cookie and returns the newly created user.
   *
   * @param {string} email - The user's email address.
   * @param {string} password - The user's password.
   *
   * @returns {Promise<RegisterResult | RegisterError>} The result of the
   *   registration attempt.
   */
  const register = async (
    email: string,
    password: string,
  ): Promise<RegisterResult | RegisterError> => {
    const { apiBase } = useRuntimeConfig().public;

    try {
      const data = await sendRegisterRequest(email, password, apiBase);

      isAuthenticated.value = true;

      return { data, ok: true };
    } catch (error: unknown) {
      if (isOfetchError(error)) {
        return handleRegisterError(error.status, error.data);
      }

      return { errors: [], message: "Network error", ok: false, status: 0 };
    }
  };

  return { isAuthenticated, login, register };
});

export default useAuthStore;
