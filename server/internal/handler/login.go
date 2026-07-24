package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"strings"

	"github.com/Jarmos-san/arthika/server/internal/api"
	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
	"golang.org/x/crypto/bcrypt"
)

// login200WithCookie wraps Login200JSONResponse to set an HttpOnly JWT cookie
// on the response before writing the JSON body.
type login200WithCookie struct {
	api.Login200JSONResponse

	token string
}

// VisitLoginResponse sets the auth cookie and delegates to the embedded JSON
// response writer.
func (r login200WithCookie) VisitLoginResponse(
	responseWriter http.ResponseWriter,
) error {
	// Create the "Set-Cookie" header for the login response
	http.SetCookie(responseWriter, newAuthCookie(r.token))

	// If the login attempt was not successful then raise an error instead
	err := r.Login200JSONResponse.VisitLoginResponse(responseWriter)
	if err != nil {
		return fmt.Errorf("write response: %w", err)
	}

	return nil
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

// passwordsMatch compares a plaintext password against a bcrypt hash.
//
// It is extracted as a standalone helper so that the Login handler avoids
// inline error assignment that would trigger the nilerr and noinlineerr
// linters.
func passwordsMatch(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// Login authenticates a user and sets a signed JWT as an HttpOnly cookie.
//
// It performs the following steps in order:
//  1. Validates the request body is present.
//  2. Validates email format and password presence.
//  3. Looks up the user by email.
//  4. Compares the provided password against the stored bcrypt hash.
//  5. Generates a signed JWT and sets it as an HttpOnly cookie.
//
// Returns:
//   - login200WithCookie on success (200).
//   - Login401JSONResponse if the email is unknown or the password is wrong (401).
//   - Login422JSONResponse if input validation fails (422 Unprocessable Entity).
//   - An unwrapped error if an internal operation fails (500 Internal Server Error).
func (h *Handler) Login(
	ctx context.Context,
	req api.LoginRequestObject,
) (api.LoginResponseObject, error) {
	// Check if the login request has an appropriate body,  else return a validation
	// error response
	if req.Body == nil {
		return api.Login422JSONResponse{
			Errors: []api.ValidationError{
				{
					Field:   "body",
					Message: "request body is required",
				},
			},
		}, nil
	}

	// Serialise the user credentials from the login request and store them in-memory
	email := strings.TrimSpace(string(req.Body.Email))
	password := req.Body.Password

	// Validate the user credentials from the login request
	if errs := validateLoginRequest(email, password); len(errs) > 0 {
		return api.Login422JSONResponse{Errors: errs}, nil
	}

	// Check if the user records exists on the system, else throw a 401 HTTP response
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

	// Check if the provided passwords match, else return a 401 HTTP response with an
	// appropriate message about the wrong password.
	if !passwordsMatch(user.PasswordHash, password) {
		return api.Login401JSONResponse{
			Message: "invalid email or password",
		}, nil
	}

	// Generated the JWT for the authenticated user
	token, err := h.generateAuthToken(ctx, user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	// Attempt to parse the UUID4 assigned the user by the server, if not then raise an
	// error and log it for investigation.
	userUUID, err := uuid.Parse(user.ID)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to parse user ID as UUID", "error", err)

		return nil, fmt.Errorf("parse user ID: %w", err)
	}

	// Create a cookie to be stored on the user's machine and return an appropriate
	// HTTP response after a successful login attempt
	return login200WithCookie{
		Login200JSONResponse: api.Login200JSONResponse{
			Body:    api.LoginResponse{Id: userUUID, Email: types.Email(email)},
			Headers: api.Login200ResponseHeaders{SetCookie: nil},
		},
		token: token,
	}, nil
}
