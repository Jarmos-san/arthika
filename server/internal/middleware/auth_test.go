// Package middleware_test contains black-box tests for the middleware package.
package middleware_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jarmos-san/arthika/server/internal/auth"
	"github.com/Jarmos-san/arthika/server/internal/config"
	"github.com/Jarmos-san/arthika/server/internal/middleware"
)

const (
	testSecret = "test-secret-key"
	cookieName = "token"
)

// testConfig returns a config with the test secret.
func testConfig() config.Config {
	return config.Config{ //nolint:exhaustruct // Only TokenSecret is needed for auth middleware tests.
		TokenSecret: testSecret,
	}
}

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

			authMiddleware := middleware.NewAuthMiddleware(testConfig())
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

// TestAuthMiddleware_MissingToken verifies that a protected endpoint returns
// 401 when no Authorization header or cookie is provided.
func TestAuthMiddleware_MissingToken(t *testing.T) {
	t.Parallel()

	authMiddleware := middleware.NewAuthMiddleware(testConfig())
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

// TestAuthMiddleware_InvalidToken verifies that a protected endpoint returns
// 401 when an invalid token is provided.
func TestAuthMiddleware_InvalidToken(t *testing.T) {
	t.Parallel()

	authMiddleware := middleware.NewAuthMiddleware(testConfig())
	handler := authMiddleware(okHandler{})

	req := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"/api/protected",
		nil,
	)
	req.Header.Set("Authorization", "Bearer invalidtoken")

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

// TestAuthMiddleware_InvalidScheme verifies that a non-Bearer Authorization
// header returns 401.
func TestAuthMiddleware_InvalidScheme(t *testing.T) {
	t.Parallel()

	authMiddleware := middleware.NewAuthMiddleware(testConfig())
	handler := authMiddleware(okHandler{})

	req := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"/api/protected",
		nil,
	)
	req.Header.Set("Authorization", "Basic dXNlcjpwYXNz")

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

// TestAuthMiddleware_ValidToken verifies that a valid token in the
// Authorization header passes through and the user context is populated.
func TestAuthMiddleware_ValidToken(t *testing.T) {
	t.Parallel()

	token, err := auth.GenerateToken("user-1", "user@example.com", testSecret)
	if err != nil {
		t.Fatalf("failed to generate test token: %v", err)
	}

	authMiddleware := middleware.NewAuthMiddleware(testConfig())

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
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

// TestAuthMiddleware_ValidCookie verifies that a valid token in a session
// cookie passes through and the user context is populated.
func TestAuthMiddleware_ValidCookie(t *testing.T) {
	t.Parallel()

	token, err := auth.GenerateToken("user-1", "user@example.com", testSecret)
	if err != nil {
		t.Fatalf("failed to generate test token: %v", err)
	}

	authMiddleware := middleware.NewAuthMiddleware(testConfig())

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
	req.AddCookie(&http.Cookie{ //nolint:exhaustruct,gosec // Test cookie — only Name/Value matter.
		Name:  cookieName,
		Value: token,
	})

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

// TestAuthMiddleware_HeaderTakesPrecedenceOverCookie verifies that when both
// an Authorization header and a cookie are present, the header is used.
func TestAuthMiddleware_HeaderTakesPrecedenceOverCookie(t *testing.T) {
	t.Parallel()

	headerToken, err := auth.GenerateToken("user-header", "header@example.com", testSecret)
	if err != nil {
		t.Fatalf("failed to generate test token: %v", err)
	}

	cookieToken, err := auth.GenerateToken("user-cookie", "cookie@example.com", testSecret)
	if err != nil {
		t.Fatalf("failed to generate test token: %v", err)
	}

	authMiddleware := middleware.NewAuthMiddleware(testConfig())

	checkHandler := http.HandlerFunc(func(responseWriter http.ResponseWriter, r *http.Request) {
		userID := auth.UserIDFromContext(r.Context())
		email := auth.EmailFromContext(r.Context())

		if userID != "user-header" {
			t.Errorf("expected userID 'user-header', got %s", userID)
		}

		if email != "header@example.com" {
			t.Errorf("expected email 'header@example.com', got %s", email)
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
	req.Header.Set("Authorization", "Bearer "+headerToken)
	req.AddCookie(&http.Cookie{ //nolint:exhaustruct,gosec // Test cookie — only Name/Value matter.
		Name:  cookieName,
		Value: cookieToken,
	})

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

// TestAuthMiddleware_InvalidCookie verifies that a protected endpoint returns
// 401 when an invalid token is provided via cookie.
func TestAuthMiddleware_InvalidCookie(t *testing.T) {
	t.Parallel()

	authMiddleware := middleware.NewAuthMiddleware(testConfig())
	handler := authMiddleware(okHandler{})

	req := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"/api/protected",
		nil,
	)
	req.AddCookie(&http.Cookie{ //nolint:exhaustruct,gosec // Test cookie — only Name/Value matter.
		Name:  cookieName,
		Value: "invalidtoken",
	})

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
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

	authMiddleware := middleware.NewAuthMiddleware(testConfig())
	handler := authMiddleware(okHandler{})

	req := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"/api/protected",
		nil,
	)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}
