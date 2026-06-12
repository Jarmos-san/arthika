package handler_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Jarmos-san/arthika/server/api"
	"github.com/Jarmos-san/arthika/server/internal/handler"
)

const testContentTypeJSON = "application/json"

type mockUserService struct {
	user      api.User
	createErr error
	token     api.TokenResponse
	loginErr  error
}

func (m mockUserService) GetUser() api.User {
	return m.user
}

func (m mockUserService) CreateUser(_ context.Context, _, _, _ string) (api.User, error) {
	return api.User{}, m.createErr
}

func (m mockUserService) Login(_ context.Context, _, _ string) (api.TokenResponse, error) {
	return m.token, m.loginErr
}

func TestGetUser_Success(t *testing.T) {
	t.Parallel()

	mockSvc := mockUserService{ //nolint:exhaustruct
		user: api.User{
			Name: "Test User",
		},
	}

	logger := slog.Default()

	userHandler := handler.NewHandler(mockSvc, logger)

	req := httptest.NewRequestWithContext(
		context.TODO(),
		http.MethodGet,
		"/users/",
		nil,
	)

	recorder := httptest.NewRecorder()

	userHandler.GetUser(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	if ct := recorder.Header().Get("Content-Type"); ct != testContentTypeJSON {
		t.Fatalf("unexpected content-type: %s", ct)
	}

	var resp api.User

	err := json.NewDecoder(recorder.Body).Decode(&resp)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Name != "Test User" {
		t.Errorf("expected name %q, got %q", "Test User", resp.Name)
	}
}

type assertError struct{}

func (assertError) Error() string { return "test error" }

func TestRegisterUser_Success(t *testing.T) {
	t.Parallel()

	mockSvc := mockUserService{} //nolint:exhaustruct

	logger := slog.Default()
	userHandler := handler.NewHandler(mockSvc, logger)

	body := `{"username":"alice","email":"alice@example.com","password":"secret"}`
	req := httptest.NewRequestWithContext(
		context.TODO(),
		http.MethodPost,
		"/users/register",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", testContentTypeJSON)

	recorder := httptest.NewRecorder()
	userHandler.RegisterUser(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, recorder.Code)
	}

	if ct := recorder.Header().Get("Content-Type"); ct != testContentTypeJSON {
		t.Fatalf("unexpected content-type: %s", ct)
	}
}

func TestRegisterUser_ValidationError(t *testing.T) {
	t.Parallel()

	mockSvc := mockUserService{} //nolint:exhaustruct

	logger := slog.Default()
	userHandler := handler.NewHandler(mockSvc, logger)

	body := `{"username":"","email":"a@b.com","password":"secret"}`
	req := httptest.NewRequestWithContext(
		context.TODO(),
		http.MethodPost,
		"/users/register",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", testContentTypeJSON)

	recorder := httptest.NewRecorder()
	userHandler.RegisterUser(recorder, req)

	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status %d, got %d", http.StatusUnprocessableEntity, recorder.Code)
	}
}

func TestLoginUser_Success(t *testing.T) {
	t.Parallel()

	mockSvc := mockUserService{ //nolint:exhaustruct
		token: api.TokenResponse{
			AccessToken: "test-token",
			TokenType:   "Bearer",
			ExpiresIn:   3600,
		},
	}

	logger := slog.Default()
	userHandler := handler.NewHandler(mockSvc, logger)

	body := `{"email":"alice@example.com","password":"secret"}`
	req := httptest.NewRequestWithContext(
		context.TODO(),
		http.MethodPost,
		"/login",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", testContentTypeJSON)

	recorder := httptest.NewRecorder()
	userHandler.LoginUser(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	if ct := recorder.Header().Get("Content-Type"); ct != testContentTypeJSON {
		t.Fatalf("unexpected content-type: %s", ct)
	}

	var resp api.TokenResponse

	err := json.NewDecoder(recorder.Body).Decode(&resp)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.AccessToken != "test-token" {
		t.Errorf("expected token %q, got %q", "test-token", resp.AccessToken)
	}
}

func TestLoginUser_ValidationError(t *testing.T) {
	t.Parallel()

	mockSvc := mockUserService{} //nolint:exhaustruct

	logger := slog.Default()
	userHandler := handler.NewHandler(mockSvc, logger)

	body := `{"email":"a@b.com","password":""}`
	req := httptest.NewRequestWithContext(
		context.TODO(),
		http.MethodPost,
		"/login",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", testContentTypeJSON)

	recorder := httptest.NewRecorder()
	userHandler.LoginUser(recorder, req)

	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status %d, got %d", http.StatusUnprocessableEntity, recorder.Code)
	}
}

func TestLoginUser_InvalidCredentials(t *testing.T) {
	t.Parallel()

	mockSvc := mockUserService{ //nolint:exhaustruct
		loginErr: assertError{},
	}

	logger := slog.Default()
	userHandler := handler.NewHandler(mockSvc, logger)

	body := `{"email":"alice@example.com","password":"wrong"}`
	req := httptest.NewRequestWithContext(
		context.TODO(),
		http.MethodPost,
		"/login",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", testContentTypeJSON)

	recorder := httptest.NewRecorder()
	userHandler.LoginUser(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, recorder.Code)
	}
}
