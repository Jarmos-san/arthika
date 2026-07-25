package handler_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jarmos-san/arthika/server/internal/api"
	"github.com/Jarmos-san/arthika/server/internal/auth"
	"github.com/Jarmos-san/arthika/server/internal/handler"
)

// TestCurrentUser_Success verifies that CurrentUser returns the authenticated
// user's ID and email when valid auth context is present.
func TestCurrentUser_Success(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		createUserFn:      nil,
		findUserByEmailFn: nil,
		countUsersFn:      nil,
	}

	hdl := handler.NewHandler(slog.Default(), mock)

	ctx := auth.NewContext(t.Context(), testUserID, testEmail)
	req := api.CurrentUserRequestObject{}

	resp, err := hdl.CurrentUser(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp == nil {
		t.Fatal("expected non-nil response")
	}

	jsonResp, ok := resp.(api.CurrentUser200JSONResponse)
	if !ok {
		t.Fatalf("expected CurrentUser200JSONResponse, got %T", resp)
	}

	if jsonResp.Id.String() != testUserID {
		t.Errorf("expected user ID %s, got %s", testUserID, jsonResp.Id)
	}

	if string(jsonResp.Email) != testEmail {
		t.Errorf("expected email %s, got %s", testEmail, jsonResp.Email)
	}
}

// TestCurrentUser_InvalidUUID verifies that CurrentUser returns an error when
// the user ID in context is not a valid UUID.
func TestCurrentUser_InvalidUUID(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		createUserFn:      nil,
		findUserByEmailFn: nil,
		countUsersFn:      nil,
	}

	hdl := handler.NewHandler(slog.Default(), mock)

	ctx := auth.NewContext(t.Context(), "not-a-uuid", testEmail)
	req := api.CurrentUserRequestObject{}

	_, err := hdl.CurrentUser(ctx, req)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

// TestCurrentUser_HTTPEndpoint_Success verifies the full HTTP stack returns
// 200 with the authenticated user's ID and email.
func TestCurrentUser_HTTPEndpoint_Success(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		createUserFn:      nil,
		findUserByEmailFn: nil,
		countUsersFn:      nil,
	}

	hdl := handler.NewHandler(slog.Default(), mock)
	strictHandler := api.NewStrictHandler(hdl, nil)

	req := httptest.NewRequestWithContext(
		auth.NewContext(t.Context(), testUserID, testEmail),
		http.MethodGet,
		"/api/users/current-user",
		nil,
	)

	rec := httptest.NewRecorder()
	strictHandler.CurrentUser(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}
}
