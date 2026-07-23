// Package middleware provides HTTP middleware for the Chi router.
package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/Jarmos-san/arthika/server/internal/auth"
	"github.com/Jarmos-san/arthika/server/internal/config"
)

// writeUnauthorized writes a standard 401 Unauthorised JSON response.
func writeUnauthorized(responseWriter http.ResponseWriter) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusUnauthorized)

	_, _ = responseWriter.Write([]byte(`{"message":"unauthorized"}`))
}

// extractToken returns the JWT from either the Authorization header or the
// session cookie. The Authorization header takes precedence.
func extractToken(request *http.Request) string {
	authHeader := request.Header.Get("Authorization")
	if token, found := strings.CutPrefix(authHeader, "Bearer "); found {
		return token
	}

	cookie, err := request.Cookie("token")
	if err == nil && cookie.Value != "" {
		return cookie.Value
	}

	return ""
}

// NewAuthMiddleware returns a Chi-compatible middleware that validates JWT
// tokens on protected routes. Public routes are passed through without
// authentication.
//
// The middleware accepts a JWT from either the Authorization header (Bearer
// scheme) or an HttpOnly session cookie named "token". The header takes
// precedence when both are present. On success it injects the authenticated
// user's ID and email into the request context via auth.NewContext.
func NewAuthMiddleware(cfg config.Config) func(http.Handler) http.Handler {
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

				// Extract the JWT from the Authorization header or session cookie
				token := extractToken(request)
				if token == "" {
					writeUnauthorized(responseWriter)

					return
				}

				// Check if the provided JWT is valid, if not, then write a 401
				// Unauthorised HTTP status code and return a message
				claims, err := auth.ValidateToken(token, cfg.TokenSecret)
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
