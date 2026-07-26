/** @description Validation errors for the login form fields. */
interface LoginErrors {
  /** @description Error message for the email field, or `undefined` if valid. */
  email: string | undefined;
  /** @description Error message for the password field, or `undefined` if valid. */
  password: string | undefined;
}

/** @description Validation errors for the register form fields. */
interface RegisterErrors {
  /** @description Error message for the email field, or `undefined` if valid. */
  email: string | undefined;
  /** @description Error message for the password field, or `undefined` if valid. */
  password: string | undefined;
  /** @description Error message for the confirm password field, or `undefined` if valid. */
  confirmPassword: string | undefined;
}

export type { LoginErrors, RegisterErrors };
