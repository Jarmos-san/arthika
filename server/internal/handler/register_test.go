package handler_test

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Jarmos-san/arthika/server/internal/api"
	"github.com/Jarmos-san/arthika/server/internal/handler"
)

func TestRegister_ReturnsNotImplemented(t *testing.T) {
	t.Parallel()

	hdl := handler.NewHandler(slog.Default())
	req := api.RegisterRequestObject{
		Body: &api.RegisterRequest{
			Email:    "test@example.com",
			Password: "supersecret",
		},
	}

	resp, err := hdl.Register(t.Context(), req)

	if resp != nil {
		t.Errorf("expected nil response, got %v", resp)
	}

	if !errors.Is(err, handler.ErrNotImplemented) {
		t.Errorf("expected ErrNotImplemented, got %v", err)
	}
}

func TestRegister_ReturnsNotImplemented_WhenBodyIsNil(t *testing.T) {
	t.Parallel()

	hdl := handler.NewHandler(slog.Default())
	req := api.RegisterRequestObject{
		Body: nil,
	}

	resp, err := hdl.Register(t.Context(), req)

	if resp != nil {
		t.Errorf("expected nil response, got %v", resp)
	}

	if !errors.Is(err, handler.ErrNotImplemented) {
		t.Errorf("expected ErrNotImplemented, got %v", err)
	}
}

func TestRegister_ErrNotImplemented_IsSentinel(t *testing.T) {
	t.Parallel()

	hdl := handler.NewHandler(slog.Default())
	req := api.RegisterRequestObject{
		Body: &api.RegisterRequest{
			Email:    "test@example.com",
			Password: "supersecret",
		},
	}

	_, err := hdl.Register(t.Context(), req)
	if err == nil {
		t.Fatal("expected non-nil error")
	}

	if !strings.Contains(err.Error(), "not implemented") {
		t.Errorf("expected error to contain 'not implemented', got %q", err.Error())
	}
}

func TestRegister_HTTPEndpoint_Returns500(t *testing.T) {
	t.Parallel()

	hdl := handler.NewHandler(slog.Default())
	strictHandler := api.NewStrictHandler(hdl, nil)

	body := `{"email":"test@example.com","password":"supersecret"}`
	req := httptest.NewRequestWithContext(
		t.Context(),
		http.MethodPost,
		"/api/users/register",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	strictHandler.Register(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestRegister_HTTPEndpoint_Returns400_WhenBodyInvalid(t *testing.T) {
	t.Parallel()

	hdl := handler.NewHandler(slog.Default())
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
