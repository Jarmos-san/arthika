// Package handler provides HTTP transport layer implementations.
//
// It is responsible for handling incoming HTTP requests, delegating business logic to
// the appropriate services, and formatting HTTP responses.
package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Jarmos-san/arthika/server/api"
)

type userService interface {
	GetUser() api.User
	CreateUser(ctx context.Context, username, email, password string) (api.User, error)
	Login(ctx context.Context, email, password string) (api.TokenResponse, error)
}

// Handler handles HTTP requests for API endpoints.
//
// Handler implements the generated api.ServerInterface by delegating business
// logic to the injected userService.
type Handler struct {
	service userService
	logger  *slog.Logger
}

// NewHandler constructs a Handler with its required dependencies.
func NewHandler(service userService, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// GetUser handles GET /users/ requests.
func (h *Handler) GetUser(w http.ResponseWriter, _ *http.Request) {
	user := h.service.GetUser()
	writeJSONResponse(w, user, http.StatusOK, h.logger)
}

// RegisterUser handles POST /users/register requests.
func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) { //nolint:varnamelen
	var body api.RegisterUserJSONRequestBody

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)

		return
	}

	if body.Username == "" || string(body.Email) == "" || body.Password == "" {
		http.Error(w, "missing required fields", http.StatusUnprocessableEntity)

		return
	}

	user, svcErr := h.service.CreateUser(
		r.Context(), body.Username, string(body.Email), body.Password,
	)
	if svcErr != nil {
		h.logger.Error(
			"failed to create user",
			slog.String("error", svcErr.Error()),
		)
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)

		return
	}

	writeJSONResponse(w, user, http.StatusCreated, h.logger)
}

// LoginUser handles POST /login requests.
func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) { //nolint:varnamelen
	var body api.LoginUserJSONRequestBody

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)

		return
	}

	if string(body.Email) == "" || body.Password == "" {
		http.Error(w, "missing required fields", http.StatusUnprocessableEntity)

		return
	}

	token, svcErr := h.service.Login(r.Context(), string(body.Email), body.Password)
	if svcErr != nil {
		h.logger.Error("login failed", slog.String("error", svcErr.Error()))
		http.Error(w, "invalid credentials", http.StatusUnauthorized)

		return
	}

	writeJSONResponse(w, token, http.StatusOK, h.logger)
}
