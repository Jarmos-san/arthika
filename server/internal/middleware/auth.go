// Package middleware provides HTTP middleware for the Chi router.
package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/Jarmos-san/arthika/server/internal/auth"
)

// writeUnauthorized writes a standard 401 Unauthorised JSON response.
func writeUnauthorized(responseWriter http.ResponseWriter) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusUnauthorized)

	_, _ = responseWriter.Write([]byte(`{"message":"unauthorized"}`))
}

// NewAuthMiddleware returns a Chi-compatible middleware that validates JWT
// tokens on protected routes. Public routes are passed through without
// authentication.
//
// The middleware expects a Bearer token in the Authorization header and
// injects the authenticated user's ID and email into the request context
// via auth.NewContext on success.
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

				// Parse the HTTP "Authorization" header, and if it does not exist then
				// write a 401 Unauthorised HTTP status code and a message
				authHeader := request.Header.Get("Authorization")
				if !strings.HasPrefix(authHeader, "Bearer ") {
					writeUnauthorized(responseWriter)

					return
				}

				// Parse the JWT from the "Authorization" token
				token := strings.TrimPrefix(authHeader, "Bearer ")

				// Check if the provided JWT is valid, if not, then write a 401
				// Unauthorised HTTP status code and return a message
				claims, err := auth.ValidateToken(token, secret)
				if err != nil {
					writeUnauthorized(responseWriter)

					return
				}

				// Continue serving the request with a context if the token verification
				// was a success
				ctx := auth.NewContext(request.Context(), claims.Subject, claims.Email)
				next.ServeHTTP(responseWriter, request.WithContext(ctx))
			},
		)
	}
}
