package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"strings"

	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
	"golang.org/x/crypto/bcrypt"

	"github.com/Jarmos-san/arthika/server/internal/api"
	"github.com/Jarmos-san/arthika/server/internal/auth"
	"github.com/Jarmos-san/arthika/server/internal/config"
	"github.com/Jarmos-san/arthika/server/internal/repository"
)

// minPasswordLength is the minimum number of characters required for a user's
// password during registration. Passwords shorter than this are rejected with a
// validation error.
const minPasswordLength = 8

// cookieMaxAge is the lifetime of the auth cookie in seconds (24 hours).
const cookieMaxAge = 86400

// errUserAlreadyExists is returned when a registration attempt uses an email
// that is already registered.
var errUserAlreadyExists = errors.New("email already registered")

// newAuthCookie creates an HttpOnly JWT cookie with the configured security
// attributes.
func newAuthCookie(token string) *http.Cookie {
	//nolint:gosec,exhaustruct // Secure is configurable via COOKIE_SECURE; only
	// relevant fields set.
	return &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		MaxAge:   cookieMaxAge,
		HttpOnly: true,
		Secure:   config.LoadConfig().CookieSecure,
		SameSite: http.SameSiteLaxMode,
	}
}

// register201WithCookie wraps Register201JSONResponse to set an HttpOnly JWT
// cookie on the response before writing the JSON body.
type register201WithCookie struct {
	api.Register201JSONResponse

	token string
}

// VisitRegisterResponse sets the auth cookie and delegates to the embedded
// JSON response writer.
func (r register201WithCookie) VisitRegisterResponse(
	responseWriter http.ResponseWriter,
) error {
	// Prepare the Set-Cookie header for the response
	http.SetCookie(responseWriter, newAuthCookie(r.token))

	// Write the response along with the cookie header, or return an error
	err := r.Register201JSONResponse.VisitRegisterResponse(responseWriter)
	if err != nil {
		return fmt.Errorf("write response: %w", err)
	}

	return nil
}

// validateRegisterRequest checks the email format and password length.
//
// It returns a slice of api.ValidationError describing each problem found.
// If all checks pass, it returns an empty (nil) slice. The caller is expected
// to return a Register422JSONResponse wrapping these errors when non-empty.
func validateRegisterRequest(email string, password string) []api.ValidationError {
	// An array of the validation errors
	var errs []api.ValidationError

	// Attempt to parse the provided user's email, or throw an error if parsing failed.
	// The returned reference to *mail.Address is thrown away since it is unnecessary.
	_, parseErr := mail.ParseAddress(email)
	if parseErr != nil {
		errs = append(errs, api.ValidationError{
			Field:   "email",
			Message: "invalid email format",
		})
	}

	// Check if the minimum length of the provided password is secure enough
	if len(password) < minPasswordLength {
		errs = append(errs, api.ValidationError{
			Field:   "password",
			Message: "password must be at least 8 characters",
		})
	}

	return errs
}

// hashPassword accepts a byte slice as input for the user's password to be hashed
// persisting it in to the database.
//
// It returns a stringified version of the hashed password or a wrapped error if the
// password hashing failed for some reason.
func hashPassword(password []byte) (string, error) {
	// Generate a password hash as a byte slice which will be converted into a string
	// later and then returned
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return string(hashedPassword), nil
}

// createUser validates that the email is not taken, hashes the password,
// and persists the new user. Returns the user ID on success.
func (h *Handler) createUser(
	ctx context.Context,
	email, password string,
) (uuid.UUID, error) {
	// Query the database and attempt to find a user with the provided email address. If
	// a user record with the provided email already exist, then return an error.
	_, err := h.querier.FindUserByEmail(ctx, email)
	if err == nil {
		return uuid.Nil, errUserAlreadyExists
	}

	// Return a SQL database error if no rows were found in the database for the
	// provided user credentials
	if !errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, fmt.Errorf("find user by email: %w", err)
	}

	// Create a password hash for database persistence
	hashedPassword, err := hashPassword([]byte(password))
	if err != nil {
		return uuid.Nil, fmt.Errorf("hash password: %w", err)
	}

	// Create a new UUID4 on the server to be assigned to the user
	userID := uuid.New()

	// Persist the user's credentials on the database with the freshly generated data
	// (the ID, the email and the hashed password).
	err = h.querier.CreateUser(ctx, repository.CreateUserParams{
		ID:           userID.String(),
		Email:        email,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("create user: %w", err)
	}

	return userID, nil
}

// Register creates a new user account.
//
// It performs the following steps in order:
//  1. Validates the request body is present.
//  2. Validates email format and password length.
//  3. Checks for an existing user, hashes the password, and persists the new user.
//  4. Generates a signed JWT and sets it as an HttpOnly cookie.
//
// Returns:
//   - register201WithCookie on success (201 Created).
//   - Register409JSONResponse if the email is already registered (409 Conflict).
//   - Register422JSONResponse if input validation fails (422 Unprocessable Entity).
//   - An unwrapped error if an internal operation fails (500 Internal Server Error).
func (h *Handler) Register(
	ctx context.Context,
	req api.RegisterRequestObject,
) (api.RegisterResponseObject, error) {
	// Return a JSON response containing the message for an invalid request body
	if req.Body == nil {
		return api.Register422JSONResponse{
			Errors: []api.ValidationError{
				{
					Field:   "body",
					Message: "request body is required",
				},
			},
		}, nil
	}

	// Store the user credentials in memory for persistence
	email := strings.TrimSpace(string(req.Body.Email))
	password := req.Body.Password

	// Validate the password's length to ensure it is secure enough
	if errs := validateRegisterRequest(email, password); len(errs) > 0 {
		return api.Register422JSONResponse{Errors: errs}, nil
	}

	// Check if the provided email address of the user already exists in the database
	userID, err := h.createUser(ctx, email, password)
	if errors.Is(err, errUserAlreadyExists) {
		return api.Register409JSONResponse{Message: "email already registered"}, nil
	}

	// Return an error response from the server, in case of a server error
	if err != nil {
		return nil, err
	}

	// Generate a JWT for the user, or return an error if it was not created
	token, err := h.generateAuthToken(ctx, userID.String(), email)
	if err != nil {
		return nil, err
	}

	// Return a valid response with a cookie header and the JWT in it, if the logic
	// above was successful.
	return register201WithCookie{
		Register201JSONResponse: api.Register201JSONResponse{
			Body:    api.RegisterResponse{Id: userID, Email: types.Email(email)},
			Headers: api.Register201ResponseHeaders{SetCookie: nil},
		},
		token: token,
	}, nil
}

// generateAuthToken creates a signed JWT for the given user.
func (h *Handler) generateAuthToken(
	ctx context.Context,
	userID, email string,
) (string, error) {
	// Fetch the secret token which is used to sign the JWT
	tokenSecret := config.LoadConfig().TokenSecret

	// Generated a signed JWT based on the email address and the secret token or return
	// an error if the server failed to produce the JWT
	token, err := auth.GenerateToken(userID, email, tokenSecret)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to generate JWT", "error", err)

		return "", fmt.Errorf("generate token: %w", err)
	}

	// Return the JWT if the token was successfully generated
	return token, nil
}
