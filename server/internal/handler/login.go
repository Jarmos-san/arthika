package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/Jarmos-san/arthika/server/internal/api"
	"github.com/Jarmos-san/arthika/server/internal/auth"
	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
	"golang.org/x/crypto/bcrypt"
)

// Login authenticates a user and returns a signed JWT.
//
// It performs the following steps in order:
//  1. Validates the request body is present.
//  2. Validates email format via net/mail.ParseAddress.
//  3. Looks up the user by email via FindUserByEmail.
//  4. Compares the provided password against the stored bcrypt hash.
//  5. Generates a signed JWT containing the user ID and email.
//
// The generated token is stored in an auth.TokenHolder in the context for the
// cookie-setting strict middleware to consume.
//
// Returns:
//   - Login200JSONResponse on success (200) with user ID and email.
//   - Login401JSONResponse if the email is unknown or the password is wrong (401).
//     The same error message is used for both cases to avoid leaking which
//     credential is incorrect.
//   - Login422JSONResponse if input validation fails (422 Unprocessable Entity).
//   - An unwrapped error with a descriptive message if an internal operation fails,
//     which the strict server middleware maps to 500 Internal Server Error.
func (h *Handler) Login(
	ctx context.Context,
	req api.LoginRequestObject,
) (api.LoginResponseObject, error) {
	// Return a 422 validation error if the request does not have a body
	if req.Body == nil {
		return api.Login422JSONResponse{
			Errors: []api.ValidationError{
				{Field: "body", Message: "request body is required"},
			},
		}, nil
	}

	// Parse the user credentials (email and password) from the request body
	email := strings.TrimSpace(string(req.Body.Email))
	password := req.Body.Password

	// Validate the user credentials and if there are any validation errors then return a
	// 422 HTTP error message
	if errs := validateLoginRequest(email, password); len(errs) > 0 {
		return api.Login422JSONResponse{Errors: errs}, nil
	}

	// Check if the provided email exists in the database, if not then throw a 401 Not
	// Found HTTP status code with an error message and log it on the server
	user, findErr := h.querier.FindUserByEmail(ctx, email)
	if findErr != nil {
		if errors.Is(findErr, sql.ErrNoRows) {
			return api.Login401JSONResponse{
				Message: "invalid email or password",
			}, nil
		}

		h.logger.ErrorContext(ctx, "failed to look up user by email", "error", findErr)

		return nil, fmt.Errorf("find user by email: %w", findErr)
	}

	// Check if the provided password by the user is correct, or throw a 401 HTTP status
	// code with an error message
	if !passwordsMatch(user.PasswordHash, password) {
		return api.Login401JSONResponse{
			Message: "invalid email or password",
		}, nil
	}

	// Generate the JWT if the user-provided credentials are correct and valid
	token, err := auth.GenerateToken(user.ID, user.Email, h.tokenSecret)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to generate JWT", "error", err)

		return nil, fmt.Errorf("generate token: %w", err)
	}

	// Store the token in the context holder for the cookie middleware to consume
	if holder := auth.TokenHolderFromContext(ctx); holder != nil {
		holder.Token = token
	}

	// Check if the provided identifier associated with the user is correct (required to
	// test against tampering)
	userUUID, err := uuid.Parse(user.ID)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to parse user ID as UUID", "error", err)

		return nil, fmt.Errorf("parse user ID: %w", err)
	}

	// Return a valid JSON response if the authentication was successful.
	// Headers are optional; the Set-Cookie header is set by the cookie middleware.
	return api.Login200JSONResponse{ //nolint:exhaustruct
		Body: api.LoginResponse{
			Id:    userUUID,
			Email: types.Email(email),
		},
	}, nil
}

// passwordsMatch compares a plaintext password against a bcrypt hash.
//
// It is extracted as a standalone helper so that the Login handler avoids
// inline error assignment that would trigger the nilerr and noinlineerr
// linters.
func passwordsMatch(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// validateLoginRequest checks the email format.
//
// It returns a slice of api.ValidationError describing each problem found.
// If all checks pass, it returns an empty (nil) slice. The caller is expected
// to return a Login422JSONResponse wrapping these errors when non-empty.
//
// Unlike registration, the password is not validated for minimum length here;
// bcrypt.CompareHashAndPassword is relied upon to reject incorrect passwords.
func validateLoginRequest(email string, password string) []api.ValidationError {
	var errs []api.ValidationError

	// Check if the provided email address is a valid email address, if not then raise
	// an error
	_, parseErr := mail.ParseAddress(email)
	if parseErr != nil {
		errs = append(errs, api.ValidationError{
			Field:   "email",
			Message: "invalid email format",
		})
	}

	// Generate an error response if the password was not provided
	if password == "" {
		errs = append(errs, api.ValidationError{
			Field:   "password",
			Message: "password is required",
		})
	}

	return errs
}
