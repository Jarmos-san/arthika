package handler_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jarmos-san/arthika/server/internal/api"
	"github.com/Jarmos-san/arthika/server/internal/handler"
)

// TestLogout_Success verifies that Logout returns a non-nil response without error.
func TestLogout_Success(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		createUserFn:           nil,
		findUserByEmailFn:      nil,
		countUsersFn:           nil,
		createAssetClassFn:     nil,
		deleteAssetClassFn:     nil,
		findAssetClassByIDFn:   nil,
		findAssetClassByNameFn: nil,
		listAssetClassesFn:     nil,
		updateAssetClassFn:     nil,
	}

	hdl := handler.NewHandler(slog.Default(), mock)
	req := api.LogoutRequestObject{}

	resp, err := hdl.Logout(t.Context(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp == nil {
		t.Fatal("expected non-nil response")
	}
}

// TestLogout_HTTPEndpoint_Success verifies the full HTTP stack returns 204 and
// sets an expired auth cookie to clear the session.
func TestLogout_HTTPEndpoint_Success(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		createUserFn:           nil,
		findUserByEmailFn:      nil,
		countUsersFn:           nil,
		createAssetClassFn:     nil,
		deleteAssetClassFn:     nil,
		findAssetClassByIDFn:   nil,
		findAssetClassByNameFn: nil,
		listAssetClassesFn:     nil,
		updateAssetClassFn:     nil,
	}

	hdl := handler.NewHandler(slog.Default(), mock)
	strictHandler := api.NewStrictHandler(hdl, nil)

	req := httptest.NewRequestWithContext(
		t.Context(),
		http.MethodPost,
		"/api/users/logout",
		nil,
	)

	rec := httptest.NewRecorder()
	strictHandler.Logout(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}

	authCookie := findAuthCookie(t, rec)
	if authCookie == nil {
		t.Fatal("expected auth_token cookie in response")
	}

	if authCookie.Value != "" {
		t.Errorf("expected empty cookie value, got %q", authCookie.Value)
	}

	if authCookie.MaxAge != -1 {
		t.Errorf("expected cookie MaxAge -1, got %d", authCookie.MaxAge)
	}

	if !authCookie.HttpOnly {
		t.Error("expected cookie to have HttpOnly=true")
	}

	if authCookie.Path != "/" {
		t.Errorf("expected cookie path '/', got '%s'", authCookie.Path)
	}
}

// TestLogout_HTTPEndpoint_CookieAttributes verifies the Secure and SameSite
// attributes of the expired auth cookie.
func TestLogout_HTTPEndpoint_CookieAttributes(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		createUserFn:           nil,
		findUserByEmailFn:      nil,
		countUsersFn:           nil,
		createAssetClassFn:     nil,
		deleteAssetClassFn:     nil,
		findAssetClassByIDFn:   nil,
		findAssetClassByNameFn: nil,
		listAssetClassesFn:     nil,
		updateAssetClassFn:     nil,
	}

	hdl := handler.NewHandler(slog.Default(), mock)
	strictHandler := api.NewStrictHandler(hdl, nil)

	req := httptest.NewRequestWithContext(
		t.Context(),
		http.MethodPost,
		"/api/users/logout",
		nil,
	)

	rec := httptest.NewRecorder()
	strictHandler.Logout(rec, req)

	authCookie := findAuthCookie(t, rec)
	if authCookie == nil {
		t.Fatal("expected auth_token cookie in response")
	}

	// COOKIE_SECURE defaults to false in test environment
	if authCookie.Secure {
		t.Error("expected Secure=false when COOKIE_SECURE is not set")
	}

	if authCookie.SameSite != http.SameSiteLaxMode {
		t.Errorf("expected SameSite=Lax, got %v", authCookie.SameSite)
	}
}
