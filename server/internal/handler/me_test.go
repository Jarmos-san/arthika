package handler_test

import (
	"log/slog"
	"testing"

	"github.com/Jarmos-san/arthika/server/internal/api"
	"github.com/Jarmos-san/arthika/server/internal/auth"
	"github.com/Jarmos-san/arthika/server/internal/handler"
	"github.com/google/uuid"
)

// TestGetCurrentUser_Success verifies that a valid context returns the user's
// ID and email.
func TestGetCurrentUser_Success(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		createUserFn:      nil,
		findUserByEmailFn: nil,
		countUsersFn:      nil,
	}

	hdl := handler.NewHandler(slog.Default(), mock, testTokenKey)

	ctx := auth.NewContext(t.Context(), testUserID, testEmail)

	resp, err := hdl.GetCurrentUser(ctx, api.GetCurrentUserRequestObject{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	userResp, ok := resp.(api.GetCurrentUser200JSONResponse)
	if !ok {
		t.Fatalf("expected GetCurrentUser200JSONResponse, got %T", resp)
	}

	expectedUUID := uuid.MustParse(testUserID)
	if userResp.Id != expectedUUID {
		t.Errorf("expected ID %s, got %s", testUserID, userResp.Id)
	}

	if string(userResp.Email) != testEmail {
		t.Errorf("expected email %s, got %s", testEmail, userResp.Email)
	}
}

// TestGetCurrentUser_EmptyContext verifies that a context without auth claims
// causes a parse error (mapped to 500 by the strict middleware).
func TestGetCurrentUser_EmptyContext(t *testing.T) {
	t.Parallel()

	mock := &mockQuerier{
		createUserFn:      nil,
		findUserByEmailFn: nil,
		countUsersFn:      nil,
	}

	hdl := handler.NewHandler(slog.Default(), mock, testTokenKey)

	_, err := hdl.GetCurrentUser(t.Context(), api.GetCurrentUserRequestObject{})
	if err == nil {
		t.Fatal("expected error for empty context, got nil")
	}
}
