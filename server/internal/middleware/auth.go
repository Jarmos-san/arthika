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

// NewAuthMiddleware returns a Chi-compatible middleware that validates JWT
// tokens from an HttpOnly cookie on protected routes. Public routes are
// passed through without authentication.
//
// The middleware reads the auth_token cookie, validates the JWT, and injects
// the authenticated user's ID and email into the request context via
// auth.NewContext on success.
func NewAuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(responseWriter http.ResponseWriter, request *http.Request) {
				// Check if the request URL is allowed to be unprotected
				if slices.Contains([]string{
					"/api/ping",
					"/api/users/register",
					"/api/users/login",
					"/api/setup/status",
				}, request.URL.Path) {
					next.ServeHTTP(responseWriter, request)

					return
				}

				// Read the JWT from the auth_token cookie. If the cookie is
				// missing or empty, return 401 Unauthorised.
				cookie, err := request.Cookie("auth_token")
				if err != nil || cookie.Value == "" {
					writeUnauthorized(responseWriter)

					return
				}

				// Validate the JWT from the cookie. If the token is invalid
				// (bad signature, expired, etc.), return 401 Unauthorised.
				claims, err := auth.ValidateToken(cookie.Value, secret)
				if err != nil {
					writeUnauthorized(responseWriter)

					return
				}

				// Continue serving the request with a context containing the
				// authenticated user's claims.
				ctx := auth.NewContext(request.Context(), claims.Subject, claims.Email)
				next.ServeHTTP(responseWriter, request.WithContext(ctx))
			},
		)
	}
}
