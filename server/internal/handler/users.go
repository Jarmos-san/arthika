// Package handler provides HTTP transport layer implementations.
//
// It is responsible for handling incoming HTTP requests, delegating business logic to
// the appropriate services, and formatting HTTP responses.
//
// Handlers should remain thin and only deal with HTTP-specific concerns such as:
//   - request parsing
//   - response encoding
//   - status code handling
//
// Business logic must be delegated to the service layer.
package handler

import (
	"log/slog"
	"net/http"

	"github.com/Jarmos-san/arthika/server/internal/service"
)

// UserHandler handles HTTP requests related to user resources.
type UserHandler struct {
	service service.UserService
	logger  *slog.Logger
}

// NewUserHandler constructs a new UserHandler with its required dependencies.
func NewUserHandler(service service.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

// CreateUser handles POST /users/register.
func (u UserHandler) CreateUser( //nolint:funlen
	writer http.ResponseWriter,
	request *http.Request,
) {
	request.Body = http.MaxBytesReader(writer, request.Body, http.DefaultMaxHeaderBytes)

	formErr := request.ParseForm()
	if formErr != nil {
		u.logger.Error("failed to parse form", slog.String("error", formErr.Error()))

		writeJSONError(writer, []map[string]any{
			{
				errKeyStatus: http.StatusInternalServerError,
				errKeyTitle:  "Form Parsing Failed",
				errKeyDetail: "Unable to parse the request form data.",
			},
		}, http.StatusInternalServerError, u.logger)

		return
	}

	username := request.FormValue("username")
	email := request.FormValue("email")
	password := request.FormValue("password")

	if username == "" {
		u.logger.Error("missing required field", slog.String("name", "username"))

		writeJSONError(
			writer, []map[string]any{validationError("username")},
			http.StatusUnprocessableEntity, u.logger,
		)

		return
	}

	if email == "" {
		u.logger.Error("missing required field", slog.String("name", "email"))

		writeJSONError(
			writer, []map[string]any{validationError("email")},
			http.StatusUnprocessableEntity, u.logger,
		)

		return
	}

	if password == "" {
		u.logger.Error("missing required field", slog.String("name", "password"))

		writeJSONError(
			writer, []map[string]any{validationError("password")},
			http.StatusUnprocessableEntity, u.logger,
		)

		return
	}

	resp, serviceErr := u.service.CreateUser(
		request.Context(),
		username,
		email,
		password,
	)
	if serviceErr != nil {
		u.logger.Error("failed to create user",
			slog.String("error", serviceErr.Error()),
		)

		writeJSONError(writer, []map[string]any{
			{
				errKeyStatus: http.StatusInternalServerError,
				errKeyTitle:  "Failed to Create User",
				errKeyDetail: "An unexpected error occurred while creating the user.",
			},
		}, http.StatusInternalServerError, u.logger)

		return
	}

	writeJSONResponse(writer, map[string]any{
		"id":       resp.ID,
		"username": resp.Username,
		"email":    resp.Email,
	}, http.StatusCreated, u.logger)

	u.logger.Info(
		"successfully created new user",
		slog.String("id", resp.ID),
		slog.String("username", resp.Username),
	)
}
