/** @description Validation errors for the login form fields. */
interface LoginErrors {
  /** @description Error message for the email field, or `undefined` if valid. */
  email: string | undefined;
  /** @description Error message for the password field, or `undefined` if valid. */
  password: string | undefined;
}

/** @description Validation errors for the registration form fields. */
interface RegisterErrors {
  /** @description Error message for the confirm password field, or `undefined` if valid. */
  confirmPassword: string | undefined;
  /** @description Error message for the email field, or `undefined` if valid. */
  email: string | undefined;
  /** @description Error message for the password field, or `undefined` if valid. */
  password: string | undefined;
  /**
   * @description Server-level error message (e.g. "email already registered"), or
   * `undefined`.
   */
  server: string | undefined;
}

export type { LoginErrors, RegisterErrors };
