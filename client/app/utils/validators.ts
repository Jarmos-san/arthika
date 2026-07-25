import type { LoginErrors, RegisterErrors } from "~/types/utils/validators";

/** @description Minimum number of characters required for the password field. */
const MIN_PASSWORD_LENGTH = 8;

const EMAIL_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/u;

/**
 * @description Validates an email value against the server's rules.
 *
 * @param {string | undefined} email - The email value to validate.
 *
 * @returns {string | undefined} An error message if invalid, or `undefined` if
 *   valid.
 */
const validateEmailField = (email: string | undefined): string | undefined => {
  if (email === undefined) {
    return "Email is required";
  }

  if (!EMAIL_REGEX.test(email)) {
    return "Please enter a valid email address";
  }

  return undefined;
};

/**
 * @description Validates a password value against the server's rules.
 *
 * @param {string | undefined} password - The password value to validate.
 *
 * @returns {string | undefined} An error message if invalid, or `undefined` if
 *   valid.
 */
const validatePasswordField = (
  password: string | undefined,
): string | undefined => {
  if (password === undefined) {
    return "Password is required";
  }

  if (password.length < MIN_PASSWORD_LENGTH) {
    return "Password must be at least 8 characters";
  }

  return undefined;
};

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
  const errors: LoginErrors = {
    email: validateEmailField(email),
    password: validatePasswordField(password),
  };

  return errors;
};

/**
 * @description Validates email, password, and confirm password fields for the registration
 * form. Returns an object with error messages for each invalid field, or
 * `undefined` for valid fields.
 *
 * @param {string | undefined} email - The email value to validate.
 * @param {string | undefined} password - The password value to validate.
 * @param {string | undefined} confirmPassword - The confirm password value to
 *   validate.
 *
 * @returns {RegisterErrors} An object containing validation errors for each
 *   field.
 */
const validateRegister = (
  email: string | undefined,
  password: string | undefined,
  confirmPassword: string | undefined,
): RegisterErrors => {
  let confirmPasswordError: string | undefined = undefined;

  if (confirmPassword === undefined) {
    confirmPasswordError = "Please confirm your password";
  } else if (password !== undefined && confirmPassword !== password) {
    confirmPasswordError = "Passwords do not match";
  }

  return {
    confirmPassword: confirmPasswordError,
    email: validateEmailField(email),
    password: validatePasswordField(password),
    server: undefined,
  };
};

export default validateLogin;
export { validateRegister };
