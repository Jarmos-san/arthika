package handler_test

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Jarmos-san/arthika/server/internal/api"
	"github.com/Jarmos-san/arthika/server/internal/auth"
	"github.com/Jarmos-san/arthika/server/internal/config"
	"github.com/Jarmos-san/arthika/server/internal/handler"
	"github.com/Jarmos-san/arthika/server/internal/middleware"
	"github.com/Jarmos-san/arthika/server/internal/repository"
	"github.com/google/uuid"
)

// testUserID is a well-known UUID used in login test expectations.
const testUserID = "550e8400-e29b-41d4-a716-446655440000"

// testPasswordHash is a bcrypt hash of testPassword ("supersecret") at MinCost,
// pre-computed so that login tests do not pay the bcrypt cost on every run.
//
//nolint:gosec // This is a test fixture, not a real credential.
const testPasswordHash = "$2a$04$vVBiX25rd3eL4C1Sp0TOy.mlm/jT9SnI7qERMHDlTEfp.mh28RBwi"

// validLoginBody is a JSON string accepted by the login endpoint for tests
// that need a well-formed request body.
const validLoginBody = `{"email":"test@example.com","password":"supersecret"}`

// testCookieConfig returns a config suitable for cookie middleware tests.
func testCookieConfig() config.Config {
	return config.Config{ //nolint:exhaustruct // Only cookie-related fields are relevant for these tests.
		TokenSecret:    testTokenKey,
		CookieSecure:   false,
		CookieSameSite: http.SameSiteLaxMode,
	}
}

// TestLogin_Success verifies valid credentials return 200 with the user's
// UUID and email. The JWT is delivered via a Set-Cookie header, not the body.
func TestLogin_Success(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		createUserFn: nil,
		findUserByEmailFn: func(_ context.Context, email string) (repository.User, error) {
			return repository.User{
				ID:           testUserID,
				Email:        email,
				PasswordHash: testPasswordHash,
			}, nil
		},
		countUsersFn: nil,
	}

	hdl := handler.NewHandler(slog.Default(), mock, testTokenKey)
	req := api.LoginRequestObject{
		Body: &api.LoginRequest{
			Email:    testEmail,
			Password: testPassword,
		},
	}

	resp, err := hdl.Login(t.Context(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loginResp, ok := resp.(api.Login200JSONResponse)
	if !ok {
		t.Fatalf("expected Login200JSONResponse, got %T", resp)
	}

	expectedUUID := uuid.MustParse(testUserID)
	if loginResp.Body.Id != expectedUUID {
		t.Errorf("expected ID %s, got %s", testUserID, loginResp.Body.Id)
	}

	if string(loginResp.Body.Email) != testEmail {
		t.Errorf("expected email %s, got %s", testEmail, loginResp.Body.Email)
	}
}

// TestLogin_Success_TokenInHolder verifies that a successful login stores the
// generated JWT in the auth.TokenHolder accessible via context.
func TestLogin_Success_TokenInHolder(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		createUserFn: nil,
		findUserByEmailFn: func(_ context.Context, email string) (repository.User, error) {
			return repository.User{
				ID:           testUserID,
				Email:        email,
				PasswordHash: testPasswordHash,
			}, nil
		},
		countUsersFn: nil,
	}

	hdl := handler.NewHandler(slog.Default(), mock, testTokenKey)

	// Simulate the cookie middleware by injecting a TokenHolder into the context.
	ctx, holder := auth.WithTokenHolder(t.Context())

	req := api.LoginRequestObject{
		Body: &api.LoginRequest{
			Email:    testEmail,
			Password: testPassword,
		},
	}

	_, err := hdl.Login(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if holder.Token == "" {
		t.Error("expected token to be stored in holder, got empty string")
	}
}

// TestLogin_WrongPassword verifies that an incorrect password returns 401
// without leaking whether the email exists.
func TestLogin_WrongPassword(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		createUserFn: nil,
		findUserByEmailFn: func(_ context.Context, email string) (repository.User, error) {
			return repository.User{
				ID:           testUserID,
				Email:        email,
				PasswordHash: testPasswordHash,
			}, nil
		},
		countUsersFn: nil,
	}

	hdl := handler.NewHandler(slog.Default(), mock, testTokenKey)
	req := api.LoginRequestObject{
		Body: &api.LoginRequest{
			Email:    testEmail,
			Password: "wrongpassword",
		},
	}

	resp, err := hdl.Login(t.Context(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	unauthorizedResp, ok := resp.(api.Login401JSONResponse)
	if !ok {
		t.Fatalf("expected Login401JSONResponse, got %T", resp)
	}

	if unauthorizedResp.Message != "invalid email or password" {
		t.Errorf(
			"expected 'invalid email or password', got %s",
			unauthorizedResp.Message,
		)
	}
}

// TestLogin_EmailNotFound verifies that an unknown email returns 401 with the
// same message as a wrong password, preventing email enumeration.
func TestLogin_EmailNotFound(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		createUserFn: nil,
		findUserByEmailFn: func(_ context.Context, _ string) (repository.User, error) {
			return repository.User{}, sql.ErrNoRows
		},
		countUsersFn: nil,
	}

	hdl := handler.NewHandler(slog.Default(), mock, testTokenKey)
	req := api.LoginRequestObject{
		Body: &api.LoginRequest{
			Email:    "unknown@example.com",
			Password: testPassword,
		},
	}

	resp, err := hdl.Login(t.Context(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	unauthorizedResp, ok := resp.(api.Login401JSONResponse)
	if !ok {
		t.Fatalf("expected Login401JSONResponse, got %T", resp)
	}

	if unauthorizedResp.Message != "invalid email or password" {
		t.Errorf(
			"expected 'invalid email or password', got %s",
			unauthorizedResp.Message,
		)
	}
}

// TestLogin_InvalidEmail verifies that a malformed email returns 422.
func TestLogin_InvalidEmail(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		createUserFn:      nil,
		findUserByEmailFn: nil,
		countUsersFn:      nil,
	}

	hdl := handler.NewHandler(slog.Default(), mock, testTokenKey)
	req := api.LoginRequestObject{
		Body: &api.LoginRequest{
			Email:    "not-an-email",
			Password: testPassword,
		},
	}

	resp, err := hdl.Login(t.Context(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	validationResp, ok := resp.(api.Login422JSONResponse)
	if !ok {
		t.Fatalf("expected Login422JSONResponse, got %T", resp)
	}

	if len(validationResp.Errors) == 0 {
		t.Fatal("expected at least one validation error")
	}

	if validationResp.Errors[0].Field != "email" {
		t.Errorf("expected field 'email', got %s", validationResp.Errors[0].Field)
	}
}

// TestLogin_NilBody verifies that a request with no body returns 422.
func TestLogin_NilBody(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		createUserFn:      nil,
		findUserByEmailFn: nil,
		countUsersFn:      nil,
	}

	hdl := handler.NewHandler(slog.Default(), mock, testTokenKey)
	req := api.LoginRequestObject{
		Body: nil,
	}

	resp, err := hdl.Login(t.Context(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	validationResp, ok := resp.(api.Login422JSONResponse)
	if !ok {
		t.Fatalf("expected Login422JSONResponse, got %T", resp)
	}

	if len(validationResp.Errors) == 0 {
		t.Fatal("expected at least one validation error")
	}
}

// TestLogin_EmptyPassword verifies that an empty password returns 422.
func TestLogin_EmptyPassword(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		createUserFn:      nil,
		findUserByEmailFn: nil,
		countUsersFn:      nil,
	}

	hdl := handler.NewHandler(slog.Default(), mock, testTokenKey)
	req := api.LoginRequestObject{
		Body: &api.LoginRequest{
			Email:    testEmail,
			Password: "",
		},
	}

	resp, err := hdl.Login(t.Context(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	validationResp, ok := resp.(api.Login422JSONResponse)
	if !ok {
		t.Fatalf("expected Login422JSONResponse, got %T", resp)
	}

	if len(validationResp.Errors) == 0 {
		t.Fatal("expected at least one validation error")
	}

	if validationResp.Errors[0].Field != "password" {
		t.Errorf("expected field 'password', got %s", validationResp.Errors[0].Field)
	}
}

// TestLogin_HTTPEndpoint_Success verifies the full HTTP stack returns 200 for
// valid credentials and sets a Set-Cookie header.
func TestLogin_HTTPEndpoint_Success(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		createUserFn: nil,
		findUserByEmailFn: func(_ context.Context, email string) (repository.User, error) {
			return repository.User{
				ID:           testUserID,
				Email:        email,
				PasswordHash: testPasswordHash,
			}, nil
		},
		countUsersFn: nil,
	}

	hdl := handler.NewHandler(slog.Default(), mock, testTokenKey)
	cookieMW := middleware.NewCookieMiddleware(testCookieConfig())
	strictHandler := api.NewStrictHandler(hdl, []api.StrictMiddlewareFunc{cookieMW})

	req := httptest.NewRequestWithContext(
		t.Context(),
		http.MethodPost,
		"/api/users/login",
		strings.NewReader(validLoginBody),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	strictHandler.Login(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	cookieHeader := rec.Header().Get("Set-Cookie")
	if cookieHeader == "" {
		t.Fatal("expected Set-Cookie header to be set")
	}

	if !strings.Contains(cookieHeader, "token=") {
		t.Errorf("expected cookie to contain 'token=', got %q", cookieHeader)
	}

	if !strings.Contains(cookieHeader, "HttpOnly") {
		t.Errorf("expected cookie to be HttpOnly, got %q", cookieHeader)
	}

	if !strings.Contains(cookieHeader, "Path=/") {
		t.Errorf("expected cookie Path to be /, got %q", cookieHeader)
	}
}

// TestLogin_HTTPEndpoint_BadCredentials verifies the full HTTP stack returns
// 401 for invalid credentials and does not set a cookie.
func TestLogin_HTTPEndpoint_BadCredentials(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		createUserFn: nil,
		findUserByEmailFn: func(_ context.Context, email string) (repository.User, error) {
			return repository.User{
				ID:           testUserID,
				Email:        email,
				PasswordHash: testPasswordHash,
			}, nil
		},
		countUsersFn: nil,
	}

	hdl := handler.NewHandler(slog.Default(), mock, testTokenKey)
	cookieMW := middleware.NewCookieMiddleware(testCookieConfig())
	strictHandler := api.NewStrictHandler(hdl, []api.StrictMiddlewareFunc{cookieMW})

	body := `{"email":"test@example.com","password":"wrongpassword"}`
	req := httptest.NewRequestWithContext(
		t.Context(),
		http.MethodPost,
		"/api/users/login",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	strictHandler.Login(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}

	if cookieHeader := rec.Header().Get("Set-Cookie"); cookieHeader != "" {
		t.Errorf("expected no Set-Cookie header on 401, got %q", cookieHeader)
	}
}

// TestLogin_HTTPEndpoint_InvalidBody verifies the full HTTP stack returns 400
// for non-JSON input.
func TestLogin_HTTPEndpoint_InvalidBody(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		createUserFn:      nil,
		findUserByEmailFn: nil,
		countUsersFn:      nil,
	}

	hdl := handler.NewHandler(slog.Default(), mock, testTokenKey)
	strictHandler := api.NewStrictHandler(hdl, nil)

	req := httptest.NewRequestWithContext(
		t.Context(),
		http.MethodPost,
		"/api/users/login",
		strings.NewReader(`not-json`),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	strictHandler.Login(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
