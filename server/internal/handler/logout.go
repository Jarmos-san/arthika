package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Jarmos-san/arthika/server/internal/api"
	"github.com/Jarmos-san/arthika/server/internal/config"
)

// expiredAuthCookie creates an HttpOnly cookie configured to delete the
// existing auth_token cookie from the client. MaxAge=-1 instructs the browser
// to remove the cookie immediately.
func expiredAuthCookie() *http.Cookie {
	//nolint:gosec,exhaustruct // Secure is configurable via COOKIE_SECURE; only
	// relevant fields set.
	return &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   config.LoadConfig().CookieSecure,
		SameSite: http.SameSiteLaxMode,
	}
}

// logout204WithCookie wraps Logout204Response to clear the auth cookie
// before writing the 204 response.
type logout204WithCookie struct {
	api.Logout204Response
}

// VisitLogoutResponse sets the expired auth cookie and delegates to the
// embedded response writer.
func (r logout204WithCookie) VisitLogoutResponse(
	responseWriter http.ResponseWriter,
) error {
	http.SetCookie(responseWriter, expiredAuthCookie())

	err := r.Logout204Response.VisitLogoutResponse(responseWriter)
	if err != nil {
		return fmt.Errorf("write response: %w", err)
	}

	return nil
}

// Logout clears the authentication cookie.
//
// Returns:
//   - logout204WithCookie on success (204 No Content).
func (h *Handler) Logout(
	_ context.Context,
	_ api.LogoutRequestObject,
) (api.LogoutResponseObject, error) {
	return logout204WithCookie{
		Logout204Response: api.Logout204Response{},
	}, nil
}
