package middleware

import (
	"context"
	"net/http"

	"github.com/Jarmos-san/arthika/server/internal/api"
	"github.com/Jarmos-san/arthika/server/internal/auth"
	"github.com/Jarmos-san/arthika/server/internal/config"
)

// cookieMaxAge is the lifetime of the session cookie in seconds. It mirrors
// the JWT expiry defined in the auth package.
const cookieMaxAge = 24 * 60 * 60

// NewCookieMiddleware returns a strict-server middleware that sets an HttpOnly
// session cookie containing the signed JWT after a successful login.
//
// The handler stores the generated token in an auth.TokenHolder via the
// request context. This middleware reads it back and writes the Set-Cookie
// header before the response is flushed to the client.
func NewCookieMiddleware(cfg config.Config) api.StrictMiddlewareFunc {
	return func(
		handler api.StrictHandlerFunc,
		operationID string,
	) api.StrictHandlerFunc {
		return func(
			ctx context.Context,
			responseWriter http.ResponseWriter,
			request *http.Request,
			requestBody any,
		) (any, error) {
			// Create a TokenHolder and inject it into the context so the login
			// handler can store the generated JWT there.
			ctx, holder := auth.WithTokenHolder(ctx)

			response, err := handler(ctx, responseWriter, request, requestBody)
			if err != nil {
				return response, err
			}

			// Only set the cookie for the Login operation on success.
			if operationID == "Login" && holder.Token != "" {
				cookie := &http.Cookie{ //nolint:exhaustruct,gosec // Cookie fields are set explicitly; Secure is config-driven.
					Name:     "token",
					Value:    holder.Token,
					Path:     "/",
					MaxAge:   cookieMaxAge,
					HttpOnly: true,
					Secure:   cfg.CookieSecure,
					SameSite: cfg.CookieSameSite,
					Domain:   cfg.CookieDomain,
				}
				http.SetCookie(responseWriter, cookie)
			}

			return response, err
		}
	}
}
