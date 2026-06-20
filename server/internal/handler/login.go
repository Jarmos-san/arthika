package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
	"golang.org/x/crypto/bcrypt"

	"github.com/Jarmos-san/arthika/server/internal/api"
	"github.com/Jarmos-san/arthika/server/internal/config"
)

// jwtClaims holds the custom claims embedded in the JWT issued during login.
//
// Embedding jwt.RegisteredClaims provides standard fields (sub, iat, exp, etc.)
// while the Email field carries the authenticated user's email address.
// This type is defined here for co-location with the Login handler but is
// designed to be extracted into a shared auth package when the JWT middleware
// (task 1c) is implemented.
type jwtClaims struct {
	jwt.RegisteredClaims

	Email string `json:"email"`
}

// tokenExpiry is the lifetime of issued JWTs.
const tokenExpiry = 24 * time.Hour

// generateToken creates a signed JWT for the given user.
//
// The token embeds the user's ID (in the Subject claim) and email address
// (in a custom claim), signed with HMAC-SHA256 using the provided secret.
// The token expires after tokenExpiry from the time of issuance.
//
// This function is a standalone building block so that it can be extracted
// into a shared auth package when the middleware layer (task 1c) needs it.
func generateToken(userID, email, secret string) (string, error) {
	// Timestamp required to create the token
	now := time.Now()

	// Create a "claims" object
	claims := jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Arthika API",
			Subject:   userID,
			Audience:  jwt.ClaimStrings{},
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenExpiry)),
			NotBefore: nil,
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.NewString(),
		},
		Email: email,
	}

	// Create a token signed using the HMAC-SHA256 algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create a bytes string from the token, or return an error if it failed
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signedToken, nil
}

// Login authenticates a user and returns a signed JWT.
//
// It performs the following steps in order:
//  1. Validates the request body is present.
//  2. Validates email format via net/mail.ParseAddress.
//  3. Looks up the user by email via FindUserByEmail.
//  4. Compares the provided password against the stored bcrypt hash.
//  5. Generates a signed JWT containing the user ID and email.
//
// Returns:
//   - Login200JSONResponse on success (200) with a signed JWT, user ID and email.
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

	// Validate the user credentials and if there are any validation erros then return a
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
				Message: "provided email address was not found",
			}, nil
		}

		h.logger.ErrorContext(ctx, "failed to look up user by email", "error", findErr)

		return nil, fmt.Errorf("find user by email: %w", findErr)
	}

	// Check if the provided password by the user is correct, or throw a 401 HTTP status
	// code with an error message
	if !passwordsMatch(user.PasswordHash, password) {
		return api.Login401JSONResponse{
			Message: "invalid or wrong password",
		}, nil
	}

	// Load the secret token (used to generate the JWT) from the configuration variables
	tokenSecret := config.LoadConfig().TokenSecret

	// Generate the JWT if the user-provided credentials are correct and valid
	token, err := generateToken(user.ID, user.Email, tokenSecret)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to generate JWT", "error", err)

		return nil, fmt.Errorf("generate token: %w", err)
	}

	// Check if the provided identifier associated with the user is correct (required to
	// test against tampering)
	userUUID, err := uuid.Parse(user.ID)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to parse user ID as UUID", "error", err)

		return nil, fmt.Errorf("parse user ID: %w", err)
	}

	// Return a valid JSON response if the authentication was successful
	return api.Login200JSONResponse{
		Token: token,
		Id:    userUUID,
		Email: types.Email(email),
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
