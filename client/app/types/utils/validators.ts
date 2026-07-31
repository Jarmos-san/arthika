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

/** @description Validation errors for the asset class form fields. */
interface AssetClassFormErrors {
  /** @description Error message for the name field, or `undefined` if valid. */
  name: string | undefined;
}

export type { AssetClassFormErrors, LoginErrors, RegisterErrors };
