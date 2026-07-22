/** @description Validation errors for the login form fields. */
export interface LoginErrors {
  /** @description Error message for the email field, or `undefined` if valid. */
  email: string | undefined;
  /** @description Error message for the password field, or `undefined` if valid. */
  password: string | undefined;
}
