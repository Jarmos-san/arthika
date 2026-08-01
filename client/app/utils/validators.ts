import type {
  AssetClassErrors,
  LoginErrors,
  RegisterErrors,
} from "~/types/utils/validators";

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

/**
 * @description Validates email, password, and confirm password fields. Reuses the same email
 * and password rules as `validateLogin`, and additionally checks that the
 * confirm password matches the password.
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
  const loginErrors = validateLogin(email, password);
  const errors: RegisterErrors = {
    confirmPassword: undefined,
    email: loginErrors.email,
    password: loginErrors.password,
  };

  if (confirmPassword === undefined) {
    errors.confirmPassword = "Please confirm your password";
  } else if (confirmPassword !== password) {
    errors.confirmPassword = "Passwords do not match";
  }

  return errors;
};

/**
 * @description Validates the asset class name field against the rule the server will enforce
 * (name required, non-blank). Returns an object with an error message for the
 * name field, or `undefined` when valid.
 *
 * @param {string | undefined} name - The name value to validate.
 *
 * @returns {AssetClassErrors} An object containing validation errors for the
 *   name field.
 */
const validateAssetClass = (name: string | undefined): AssetClassErrors => {
  if (name === undefined || name.trim() === "") {
    return { name: "Name is required" };
  }

  return { name: undefined };
};

export { validateAssetClass, validateLogin, validateRegister };
