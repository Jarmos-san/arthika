// Package middleware provides HTTP middleware for the Chi router.
package middleware

import (
	"net/http"
	"slices"

	"github.com/Jarmos-san/arthika/server/internal/auth"
)

// writeUnauthorized writes a standard 401 Unauthorised JSON response.
func writeUnauthorized(responseWriter http.ResponseWriter) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusUnauthorized)

	_, _ = responseWriter.Write([]byte(`{"message":"unauthorized"}`))
}

// publicPaths contains routes that bypass authentication.
//
//nolint:gochecknoglobals // Fixed set of paths, safe as package-level.
var publicPaths = []string{
	"/api/ping",
	"/api/users/register",
	"/api/users/login",
	"/api/setup/status",
}

// authHandler is an HTTP handler that validates JWT tokens from cookies.
type authHandler struct {
	next   http.Handler
	secret string
}

func (h *authHandler) ServeHTTP(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	// Check if the requested path is unprotected and does not require authentication
	if slices.Contains(publicPaths, request.URL.Path) {
		h.next.ServeHTTP(responseWriter, request)

		return
	}

	// Read the cookie from the request, if it was not found then raise a
	// 403 Unauthorised HTTP status code with a message
	cookie, err := request.Cookie("auth_token")
	if err != nil || cookie.Value == "" {
		writeUnauthorized(responseWriter)

		return
	}

	// Validate the JWT (provided in the cookie), if not then raise a 403 Unauthorised
	// HTTP status code with a message
	claims, err := auth.ValidateToken(cookie.Value, h.secret)
	if err != nil {
		writeUnauthorized(responseWriter)

		return
	}

	// Create a new request context with the JWT and serve the request
	ctx := auth.NewContext(request.Context(), claims.Subject, claims.Email)
	h.next.ServeHTTP(responseWriter, request.WithContext(ctx))
}

// NewAuthMiddleware returns a Chi-compatible middleware that validates JWT
// tokens on protected routes. Public routes are passed through without
// authentication.
//
// The middleware reads the JWT from an HttpOnly cookie named "auth_token" and
// injects the authenticated user's ID and email into the request context via
// auth.NewContext on success.
func NewAuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return &authHandler{next: next, secret: secret}
	}
}
