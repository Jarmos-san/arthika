package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/Jarmos-san/arthika/server/internal/api"
	"github.com/Jarmos-san/arthika/server/internal/repository"
	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
	"golang.org/x/crypto/bcrypt"
)

// minPasswordLength is the minimum number of characters required for a user's
// password during registration. Passwords shorter than this are rejected with a
// validation error.
const minPasswordLength = 8

// Register creates a new user account.
//
// It performs the following steps in order:
//  1. Validates the request body is present.
//  2. Validates email format via net/mail.ParseAddress.
//  3. Validates password meets the minimum length requirement (see minPasswordLength).
//  4. Checks for an existing user with the same email via FindUserByEmail.
//  5. Hashes the password with bcrypt (DefaultCost).
//  6. Generates a UUID for the new user.
//  7. Persists the user via CreateUser.
//
// Returns:
//   - Register201JSONResponse on success (201 Created).
//   - Register409JSONResponse if the email is already registered (409 Conflict).
//   - Register422JSONResponse if input validation fails (422 Unprocessable Entity).
//   - An unwrapped error with a descriptive message if an internal operation fails,
//     which the strict server middleware maps to 500 Internal Server Error.
func (h *Handler) Register(
	ctx context.Context,
	req api.RegisterRequestObject,
) (api.RegisterResponseObject, error) {
	// Return a 422 Unprocessable Content HTTP status code if the request has no body
	if req.Body == nil {
		return api.Register422JSONResponse{
			Errors: []api.ValidationError{
				{Field: "body", Message: "request body is required"},
			},
		}, nil
	}

	// Parse the request body for the user credentials (email and password)
	email := strings.TrimSpace(string(req.Body.Email))
	password := req.Body.Password

	// Create an error response if the validation of the user credentials failed
	if errs := validateRegisterRequest(email, password); len(errs) > 0 {
		return api.Register422JSONResponse{Errors: errs}, nil
	}

	// Query the database to check if a user with the same email credential
	// already exists in the database records, if it exists then return an error
	// response
	_, err := h.querier.FindUserByEmail(ctx, email)
	if err == nil {
		return api.Register409JSONResponse{
			Message: "email already registered",
		}, nil
	}

	// Throw an if the database query failed for whatever reason
	if !errors.Is(err, sql.ErrNoRows) {
		h.logger.ErrorContext(ctx, "failed to check for existing user", "error", err)

		return nil, fmt.Errorf("find user by email: %w", err)
	}

	// Create a hashed string for the password, or throw an error if the hashing failed
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to hash password", "error", err)

		return nil, fmt.Errorf("hash password: %w", err)
	}

	// Create a UUID4 object to be assigned to the new registered user
	userID := uuid.New()

	// Query the database to insert the new user records in the database with their
	// generated identifier, email and hashed password, or throw an error if it failed
	err = h.querier.CreateUser(ctx, repository.CreateUserParams{
		ID:           userID.String(),
		Email:        email,
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to create user", "error", err)

		return nil, fmt.Errorf("create user: %w", err)
	}

	// Return a response containing the new user's identifier and the email address if
	// records were successfully inserted in to the database
	return api.Register201JSONResponse{
		Id:    userID,
		Email: types.Email(email),
	}, nil
}

// validateRegisterRequest checks the email format and password length.
//
// It returns a slice of api.ValidationError describing each problem found.
// If all checks pass, it returns an empty (nil) slice. The caller is expected
// to return a Register422JSONResponse wrapping these errors when non-empty.
func validateRegisterRequest(email string, password string) []api.ValidationError {
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
