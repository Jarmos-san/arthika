package handler_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Jarmos-san/arthika/server/internal/api"
	"github.com/Jarmos-san/arthika/server/internal/handler"
	"github.com/Jarmos-san/arthika/server/internal/repository"
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

// assertLoginCookie verifies the Set-Cookie header contains the expected
// auth_token cookie with the required security attributes.
func assertLoginCookie(t *testing.T, setCookie string) {
	t.Helper()

	if setCookie == "" {
		t.Fatal("expected Set-Cookie header to be set")
	}

	if !strings.HasPrefix(setCookie, "auth_token=") {
		t.Errorf("expected cookie to start with auth_token=, got %s", setCookie)
	}

	for _, attr := range []string{"HttpOnly", "Path=/", "SameSite=Strict", "Max-Age=86400"} {
		if !strings.Contains(setCookie, attr) {
			t.Errorf("expected cookie to contain %s, got %s", attr, setCookie)
		}
	}
}

// TestLogin_Success verifies valid credentials return 200 with an HttpOnly
// cookie containing the JWT, and a JSON body with only the user's ID and email.
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

	hdl := handler.NewHandler(slog.Default(), mock)
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

	rec := httptest.NewRecorder()

	err = resp.VisitLoginResponse(rec)
	if err != nil {
		t.Fatalf("VisitLoginResponse failed: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	assertLoginCookie(t, rec.Header().Get("Set-Cookie"))

	var body api.LoginResponse

	err = json.NewDecoder(rec.Body).Decode(&body)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body.Id.String() != testUserID {
		t.Errorf("expected ID %s, got %s", testUserID, body.Id)
	}

	if string(body.Email) != testEmail {
		t.Errorf("expected email %s, got %s", testEmail, body.Email)
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

	hdl := handler.NewHandler(slog.Default(), mock)
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

	hdl := handler.NewHandler(slog.Default(), mock)
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

	hdl := handler.NewHandler(slog.Default(), mock)
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

	hdl := handler.NewHandler(slog.Default(), mock)
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

	hdl := handler.NewHandler(slog.Default(), mock)
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

// TestLogin_HTTPEndpoint_Success verifies the full HTTP stack returns 200 with
// a Set-Cookie header for valid credentials.
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

	hdl := handler.NewHandler(slog.Default(), mock)
	strictHandler := api.NewStrictHandler(hdl, nil)

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

	setCookie := rec.Header().Get("Set-Cookie")
	if setCookie == "" {
		t.Fatal("expected Set-Cookie header to be set")
	}

	if !strings.HasPrefix(setCookie, "auth_token=") {
		t.Errorf("expected cookie to start with auth_token=, got %s", setCookie)
	}
}

// TestLogin_HTTPEndpoint_BadCredentials verifies the full HTTP stack returns
// 401 for invalid credentials.
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

	hdl := handler.NewHandler(slog.Default(), mock)
	strictHandler := api.NewStrictHandler(hdl, nil)

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

	hdl := handler.NewHandler(slog.Default(), mock)
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
