/** @description Validation errors for the login form fields. */
interface LoginErrors {
  /** @description Error message for the email field, or `undefined` if valid. */
  email: string | undefined;
  /** @description Error message for the password field, or `undefined` if valid. */
  password: string | undefined;
}

/** @description Minimum number of characters required for the password field. */
const MIN_PASSWORD_LENGTH = 8;

/**
 * @description Validates email and password fields against the same rules the server
 * enforces (valid email format, password minimum 8 characters). Returns an
 * object with error messages for each invalid field, or `undefined` if the
 * field is valid.
 *
 * @param {string | undefined} email - The email value to validate.
 * @param {string | undefined} password - The password value to validate.
 *
 * @returns {LoginErrors} An object containing validation errors for each field.
 */
const validateLogin = (
  email: string | undefined,
  password: string | undefined,
): LoginErrors => {
  const errors: LoginErrors = { email: undefined, password: undefined };

  if (email === undefined) {
    errors.email = "Email is required";
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/u.test(email)) {
    errors.email = "Please enter a valid email address";
  }

  if (password === undefined) {
    errors.password = "Password is required";
  } else if (password.length < MIN_PASSWORD_LENGTH) {
    errors.password = "Password must be at least 8 characters";
  }

  return errors;
};

export { type LoginErrors, validateLogin };
