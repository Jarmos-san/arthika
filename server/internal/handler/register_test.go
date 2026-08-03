// Package handler_test contains black-box tests for the handler package.
//
// Tests use a mockQuerier implementing repository.Querier to isolate
// handler logic from database access.
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
	"github.com/Jarmos-san/arthika/server/internal/handler"
	"github.com/Jarmos-san/arthika/server/internal/repository"
)

const (
	testEmail    = "test@example.com"
	testPassword = "supersecret"
)

// mockQuerier implements repository.Querier with configurable function fields.
// Each test sets only the functions it needs; nil functions panic if called.
type mockQuerier struct {
	createUserFn         func(ctx context.Context, arg repository.CreateUserParams) error
	findUserByEmailFn    func(ctx context.Context, email string) (repository.User, error)
	countUsersFn         func(ctx context.Context) (int64, error)
	createAssetClassFn   func(ctx context.Context, arg repository.CreateAssetClassParams) error
	deleteAssetClassFn   func(ctx context.Context, id string) (string, error)
	findAssetClassByIDFn func(ctx context.Context, id string) (repository.AssetClass, error)
	listAssetClassesFn   func(ctx context.Context) ([]repository.AssetClass, error)
	updateAssetClassFn   func(ctx context.Context, arg repository.UpdateAssetClassParams) (repository.AssetClass, error)
}

// CreateUser delegates to m.createUserFn.
func (m *mockQuerier) CreateUser(
	ctx context.Context,
	arg repository.CreateUserParams,
) error {
	return m.createUserFn(ctx, arg)
}

// FindUserByEmail delegates to m.findUserByEmailFn.
func (m *mockQuerier) FindUserByEmail(
	ctx context.Context,
	email string,
) (repository.User, error) {
	return m.findUserByEmailFn(ctx, email)
}

// CountUsers delegates to m.countUsersFn.
func (m *mockQuerier) CountUsers(ctx context.Context) (int64, error) {
	return m.countUsersFn(ctx)
}

// CreateAssetClass delegates to m.createAssetClassFn.
func (m *mockQuerier) CreateAssetClass(
	ctx context.Context,
	arg repository.CreateAssetClassParams,
) error {
	return m.createAssetClassFn(ctx, arg)
}

// DeleteAssetClass delegates to m.deleteAssetClassFn.
func (m *mockQuerier) DeleteAssetClass(ctx context.Context, id string) (string, error) {
	return m.deleteAssetClassFn(ctx, id)
}

// FindAssetClassByID delegates to m.findAssetClassByIDFn.
func (m *mockQuerier) FindAssetClassByID(
	ctx context.Context,
	id string,
) (repository.AssetClass, error) {
	return m.findAssetClassByIDFn(ctx, id)
}

// ListAssetClasses delegates to m.listAssetClassesFn.
func (m *mockQuerier) ListAssetClasses(ctx context.Context) ([]repository.AssetClass, error) {
	return m.listAssetClassesFn(ctx)
}

// UpdateAssetClass delegates to m.updateAssetClassFn.
func (m *mockQuerier) UpdateAssetClass(
	ctx context.Context,
	arg repository.UpdateAssetClassParams,
) (repository.AssetClass, error) {
	return m.updateAssetClassFn(ctx, arg)
}

// TestRegister_Success verifies a valid request returns 201 with the registered email.
func TestRegister_Success(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		createUserFn: func(_ context.Context, _ repository.CreateUserParams) error {
			return nil
		},
		findUserByEmailFn: func(_ context.Context, _ string) (repository.User, error) {
			return repository.User{}, sql.ErrNoRows
		},
		countUsersFn:         nil,
		createAssetClassFn:   nil,
		deleteAssetClassFn:   nil,
		findAssetClassByIDFn: nil,
		listAssetClassesFn:   nil,
		updateAssetClassFn:   nil,
	}

	hdl := handler.NewHandler(slog.Default(), mock)
	req := api.RegisterRequestObject{
		Body: &api.RegisterRequest{
			Email:    testEmail,
			Password: testPassword,
		},
	}

	resp, err := hdl.Register(t.Context(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp == nil {
		t.Fatal("expected non-nil response")
	}

	// Verify the response is not an error type (409 or 422)
	switch resp.(type) {
	case api.Register409JSONResponse:
		t.Fatal("unexpected 409 conflict response")
	case api.Register422JSONResponse:
		t.Fatal("unexpected 422 validation response")
	}
}

// TestRegister_DuplicateEmail verifies that registering an existing email returns 409.
func TestRegister_DuplicateEmail(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		createUserFn: nil,
		findUserByEmailFn: func(_ context.Context, _ string) (repository.User, error) {
			return repository.User{
				ID:           "existing-id",
				Email:        testEmail,
				PasswordHash: "",
			}, nil
		},
		countUsersFn:         nil,
		createAssetClassFn:   nil,
		deleteAssetClassFn:   nil,
		findAssetClassByIDFn: nil,
		listAssetClassesFn:   nil,
		updateAssetClassFn:   nil,
	}

	hdl := handler.NewHandler(slog.Default(), mock)
	req := api.RegisterRequestObject{
		Body: &api.RegisterRequest{
			Email:    testEmail,
			Password: testPassword,
		},
	}

	resp, err := hdl.Register(t.Context(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	conflictResp, ok := resp.(api.Register409JSONResponse)
	if !ok {
		t.Fatalf("expected Register409JSONResponse, got %T", resp)
	}

	if conflictResp.Message != "email already registered" {
		t.Errorf("expected 'email already registered', got %s", conflictResp.Message)
	}
}

// TestRegister_InvalidEmail verifies that a malformed email returns 422.
//
//nolint:dupl // Validation scaffold mirrors the Login counterpart.
func TestRegister_InvalidEmail(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		createUserFn:         nil,
		findUserByEmailFn:    nil,
		countUsersFn:         nil,
		createAssetClassFn:   nil,
		deleteAssetClassFn:   nil,
		findAssetClassByIDFn: nil,
		listAssetClassesFn:   nil,
		updateAssetClassFn:   nil,
	}

	hdl := handler.NewHandler(slog.Default(), mock)
	req := api.RegisterRequestObject{
		Body: &api.RegisterRequest{
			Email:    "not-an-email",
			Password: testPassword,
		},
	}

	resp, err := hdl.Register(t.Context(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	validationResp, ok := resp.(api.Register422JSONResponse)
	if !ok {
		t.Fatalf("expected Register422JSONResponse, got %T", resp)
	}

	if len(validationResp.Errors) == 0 {
		t.Fatal("expected at least one validation error")
	}

	if validationResp.Errors[0].Field != "email" {
		t.Errorf("expected field 'email', got %s", validationResp.Errors[0].Field)
	}
}

// TestRegister_ShortPassword verifies that a password shorter than 8 chars returns 422.
//
//nolint:dupl // Validation scaffold mirrors the Login counterpart.
func TestRegister_ShortPassword(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		createUserFn:         nil,
		findUserByEmailFn:    nil,
		countUsersFn:         nil,
		createAssetClassFn:   nil,
		deleteAssetClassFn:   nil,
		findAssetClassByIDFn: nil,
		listAssetClassesFn:   nil,
		updateAssetClassFn:   nil,
	}

	hdl := handler.NewHandler(slog.Default(), mock)
	req := api.RegisterRequestObject{
		Body: &api.RegisterRequest{
			Email:    testEmail,
			Password: "short",
		},
	}

	resp, err := hdl.Register(t.Context(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	validationResp, ok := resp.(api.Register422JSONResponse)
	if !ok {
		t.Fatalf("expected Register422JSONResponse, got %T", resp)
	}

	if len(validationResp.Errors) == 0 {
		t.Fatal("expected at least one validation error")
	}

	if validationResp.Errors[0].Field != "password" {
		t.Errorf("expected field 'password', got %s", validationResp.Errors[0].Field)
	}
}

// TestRegister_NilBody verifies that a request with no body returns 422.
func TestRegister_NilBody(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		createUserFn:         nil,
		findUserByEmailFn:    nil,
		countUsersFn:         nil,
		createAssetClassFn:   nil,
		deleteAssetClassFn:   nil,
		findAssetClassByIDFn: nil,
		listAssetClassesFn:   nil,
		updateAssetClassFn:   nil,
	}

	hdl := handler.NewHandler(slog.Default(), mock)
	req := api.RegisterRequestObject{
		Body: nil,
	}

	resp, err := hdl.Register(t.Context(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	validationResp, ok := resp.(api.Register422JSONResponse)
	if !ok {
		t.Fatalf("expected Register422JSONResponse, got %T", resp)
	}

	if len(validationResp.Errors) == 0 {
		t.Fatal("expected at least one validation error")
	}
}

// TestRegister_HTTPEndpoint_Success verifies the full HTTP stack returns 201 for valid
// input and sets an HttpOnly JWT cookie.
func TestRegister_HTTPEndpoint_Success(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		createUserFn: func(_ context.Context, _ repository.CreateUserParams) error {
			return nil
		},
		findUserByEmailFn: func(_ context.Context, _ string) (repository.User, error) {
			return repository.User{}, sql.ErrNoRows
		},
		countUsersFn:         nil,
		createAssetClassFn:   nil,
		deleteAssetClassFn:   nil,
		findAssetClassByIDFn: nil,
		listAssetClassesFn:   nil,
		updateAssetClassFn:   nil,
	}

	hdl := handler.NewHandler(slog.Default(), mock)
	strictHandler := api.NewStrictHandler(hdl, nil)

	req := httptest.NewRequestWithContext(
		t.Context(),
		http.MethodPost,
		"/api/users/register",
		strings.NewReader(validLoginBody),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	strictHandler.Register(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	authCookie := findAuthCookie(t, rec)
	if authCookie == nil {
		t.Fatal("expected auth_token cookie")
	}

	if !authCookie.HttpOnly {
		t.Error("expected cookie to have HttpOnly=true")
	}

	if authCookie.Path != "/" {
		t.Errorf("expected cookie path '/', got '%s'", authCookie.Path)
	}

	if authCookie.MaxAge != 86400 {
		t.Errorf("expected cookie MaxAge 86400, got %d", authCookie.MaxAge)
	}

	if authCookie.Value == "" {
		t.Error("expected non-empty cookie value")
	}
}

// findAuthCookie extracts the auth_token cookie from the response.
func findAuthCookie(t *testing.T, rec *httptest.ResponseRecorder) *http.Cookie {
	t.Helper()

	for _, c := range rec.Result().Cookies() {
		if c.Name == "auth_token" {
			return c
		}
	}

	return nil
}

// TestRegister_HTTPEndpoint_DuplicateEmail verifies the full HTTP stack returns 409 for
// a duplicate email.
func TestRegister_HTTPEndpoint_DuplicateEmail(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		createUserFn: nil,
		findUserByEmailFn: func(_ context.Context, _ string) (repository.User, error) {
			return repository.User{
				ID:           "existing-id",
				Email:        "test@example.com",
				PasswordHash: "",
			}, nil
		},
		countUsersFn:         nil,
		createAssetClassFn:   nil,
		deleteAssetClassFn:   nil,
		findAssetClassByIDFn: nil,
		listAssetClassesFn:   nil,
		updateAssetClassFn:   nil,
	}

	hdl := handler.NewHandler(slog.Default(), mock)
	strictHandler := api.NewStrictHandler(hdl, nil)

	req := httptest.NewRequestWithContext(
		t.Context(),
		http.MethodPost,
		"/api/users/register",
		strings.NewReader(validLoginBody),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	strictHandler.Register(rec, req)

	if rec.Code != http.StatusConflict {
		t.Errorf("expected status %d, got %d", http.StatusConflict, rec.Code)
	}
}

// TestRegister_HTTPEndpoint_InvalidBody verifies the full HTTP stack returns 400 for
// non-JSON input.
func TestRegister_HTTPEndpoint_InvalidBody(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		createUserFn:         nil,
		findUserByEmailFn:    nil,
		countUsersFn:         nil,
		createAssetClassFn:   nil,
		deleteAssetClassFn:   nil,
		findAssetClassByIDFn: nil,
		listAssetClassesFn:   nil,
		updateAssetClassFn:   nil,
	}

	hdl := handler.NewHandler(slog.Default(), mock)
	strictHandler := api.NewStrictHandler(hdl, nil)

	req := httptest.NewRequestWithContext(
		t.Context(),
		http.MethodPost,
		"/api/users/register",
		strings.NewReader(`not-json`),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	strictHandler.Register(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

// TestRegister_HTTPEndpoint_CookieAttributes verifies the Secure and SameSite
// attributes of the auth cookie. The default config has CookieSecure=false, so
// Secure should not be set.
func TestRegister_HTTPEndpoint_CookieAttributes(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		createUserFn: func(_ context.Context, _ repository.CreateUserParams) error {
			return nil
		},
		findUserByEmailFn: func(_ context.Context, _ string) (repository.User, error) {
			return repository.User{}, sql.ErrNoRows
		},
		countUsersFn:         nil,
		createAssetClassFn:   nil,
		deleteAssetClassFn:   nil,
		findAssetClassByIDFn: nil,
		listAssetClassesFn:   nil,
		updateAssetClassFn:   nil,
	}

	hdl := handler.NewHandler(slog.Default(), mock)
	strictHandler := api.NewStrictHandler(hdl, nil)

	req := httptest.NewRequestWithContext(
		t.Context(),
		http.MethodPost,
		"/api/users/register",
		strings.NewReader(validLoginBody),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	strictHandler.Register(rec, req)

	authCookie := findAuthCookie(t, rec)
	if authCookie == nil {
		t.Fatal("expected auth_token cookie")
	}

	// COOKIE_SECURE defaults to false in test environment
	if authCookie.Secure {
		t.Error("expected Secure=false when COOKIE_SECURE is not set")
	}

	if authCookie.SameSite != http.SameSiteLaxMode {
		t.Errorf("expected SameSite=Lax, got %v", authCookie.SameSite)
	}
}
