// Package middleware_test contains black-box tests for the middleware package.
package middleware_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Jarmos-san/arthika/server/internal/auth"
	"github.com/Jarmos-san/arthika/server/internal/middleware"
	jwt "github.com/golang-jwt/jwt/v5"
)

const (
	testSecret     = "test-secret-key"
	testCookieName = "auth_token"
)

// okHandler is a simple handler that returns 200 with a JSON body for testing.
type okHandler struct{}

func (okHandler) ServeHTTP(responseWriter http.ResponseWriter, _ *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)

	_, _ = responseWriter.Write([]byte(`{"status":"ok"}`))
}

// errorResponse is used to decode 401 JSON responses.
type errorResponse struct {
	Message string `json:"message"`
}

// testCookie creates an auth_token cookie for testing with the given value.
func testCookie(value string) *http.Cookie {
	return &http.Cookie{
		Name:        testCookieName,
		Value:       value,
		Path:        "/",
		HttpOnly:    true,
		Secure:      true,
		SameSite:    http.SameSiteStrictMode,
		MaxAge:      0,
		Quoted:      false,
		Domain:      "",
		Expires:     time.Time{},
		RawExpires:  "",
		Raw:         "",
		Unparsed:    nil,
		Partitioned: false,
	}
}

// TestAuthMiddleware_PublicPaths verifies that public paths pass through
// without requiring a token.
func TestAuthMiddleware_PublicPaths(t *testing.T) {
	t.Parallel()

	publicPaths := []string{
		"/api/ping",
		"/api/users/register",
		"/api/users/login",
		"/api/setup/status",
	}

	for _, path := range publicPaths {
		t.Run(path, func(t *testing.T) {
			t.Parallel()

			authMiddleware := middleware.NewAuthMiddleware(testSecret)
			handler := authMiddleware(okHandler{})

			req := httptest.NewRequestWithContext(
				context.Background(),
				http.MethodGet,
				path,
				nil,
			)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Errorf("expected 200, got %d", rec.Code)
			}
		})
	}
}

// TestAuthMiddleware_MissingCookie verifies that a protected endpoint returns
// 401 when no auth_token cookie is provided.
func TestAuthMiddleware_MissingCookie(t *testing.T) {
	t.Parallel()

	authMiddleware := middleware.NewAuthMiddleware(testSecret)
	handler := authMiddleware(okHandler{})

	req := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"/api/protected",
		nil,
	)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}

	var resp errorResponse

	err := json.NewDecoder(rec.Body).Decode(&resp)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Message != "unauthorized" {
		t.Errorf("expected 'unauthorized', got %s", resp.Message)
	}
}

// TestAuthMiddleware_EmptyCookie verifies that a protected endpoint returns
// 401 when the auth_token cookie is present but empty.
func TestAuthMiddleware_EmptyCookie(t *testing.T) {
	t.Parallel()

	authMiddleware := middleware.NewAuthMiddleware(testSecret)
	handler := authMiddleware(okHandler{})

	req := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"/api/protected",
		nil,
	)
	req.AddCookie(testCookie(""))

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

// TestAuthMiddleware_InvalidToken verifies that a protected endpoint returns
// 401 when an invalid token is in the cookie.
func TestAuthMiddleware_InvalidToken(t *testing.T) {
	t.Parallel()

	authMiddleware := middleware.NewAuthMiddleware(testSecret)
	handler := authMiddleware(okHandler{})

	req := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"/api/protected",
		nil,
	)
	req.AddCookie(testCookie("invalidtoken"))

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

// TestAuthMiddleware_ValidToken verifies that a valid token in the cookie
// passes through and the user context is populated.
func TestAuthMiddleware_ValidToken(t *testing.T) {
	t.Parallel()

	token, err := auth.GenerateToken("user-1", "user@example.com", testSecret)
	if err != nil {
		t.Fatalf("failed to generate test token: %v", err)
	}

	authMiddleware := middleware.NewAuthMiddleware(testSecret)

	checkHandler := http.HandlerFunc(func(responseWriter http.ResponseWriter, r *http.Request) {
		userID := auth.UserIDFromContext(r.Context())
		email := auth.EmailFromContext(r.Context())

		if userID != "user-1" {
			t.Errorf("expected userID 'user-1', got %s", userID)
		}

		if email != "user@example.com" {
			t.Errorf("expected email 'user@example.com', got %s", email)
		}

		responseWriter.WriteHeader(http.StatusOK)
	})

	handler := authMiddleware(checkHandler)

	req := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"/api/protected",
		nil,
	)
	req.AddCookie(testCookie(token))

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

// TestAuthMiddleware_WrongSecret verifies that a token signed with a different
// secret is rejected.
func TestAuthMiddleware_WrongSecret(t *testing.T) {
	t.Parallel()

	token, err := auth.GenerateToken("user-1", "user@example.com", "different-secret")
	if err != nil {
		t.Fatalf("failed to generate test token: %v", err)
	}

	authMiddleware := middleware.NewAuthMiddleware(testSecret)
	handler := authMiddleware(okHandler{})

	req := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"/api/protected",
		nil,
	)
	req.AddCookie(testCookie(token))

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

// TestAuthMiddleware_ExpiredToken verifies that an expired token in the cookie
// is rejected with 401.
func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	t.Parallel()

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Arthika API",
			Subject:   "user-1",
			Audience:  jwt.ClaimStrings{},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			NotBefore: nil,
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			ID:        "expired-id",
		},
		Email: "user@example.com",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(testSecret))
	if err != nil {
		t.Fatalf("failed to sign expired token: %v", err)
	}

	authMiddleware := middleware.NewAuthMiddleware(testSecret)
	handler := authMiddleware(okHandler{})

	req := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"/api/protected",
		nil,
	)
	req.AddCookie(testCookie(signedToken))

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}
