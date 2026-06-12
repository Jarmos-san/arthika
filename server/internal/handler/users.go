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
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Jarmos-san/arthika/server/internal/dto"
	"github.com/Jarmos-san/arthika/server/internal/service"
)

const errInternalServer = "Internal Server Error"

// UserHandler handles HTTP requests related to user resources.
//
// UserHandler acts as the transport layer for user-related endpoints. It delegates
// business logic to the injected UserService and is responsible for constructing HTTP
// responses based on the results.
//
// The handler is safe for concurrent use provided its dependencies are also
// concurrency-safe.
type UserHandler struct {
	service service.UserService
	logger  *slog.Logger
}

// NewUserHandler constructs a new UserHandler with its required dependencies.
//
// Parameters:
//   - service: provides user-related business logic
//   - logger:  used for structured logging within the handler
//
// The returned handler is ready to be registered with an HTTP router.
func NewUserHandler(service service.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

// CreateUser ...
func (u UserHandler) CreateUser( //nolint:funlen
	writer http.ResponseWriter,
	request *http.Request,
) {
	request.Body = http.MaxBytesReader(writer, request.Body, http.DefaultMaxHeaderBytes)

	formErr := request.ParseForm()
	if formErr != nil {
		u.logger.Error("failed to parse form", slog.String("error", formErr.Error()))

		errResp := dto.NewErrorDocument([]dto.ErrorObject{
			{ //nolint:exhaustruct
				Status: errInternalServer,
				Code:   strconv.Itoa(http.StatusInternalServerError),
				Title:  "Form Parsing Failed",
				Detail: "Unable to parse the request form data.",
			},
		})
		writeJSONResponse(writer, errResp, http.StatusInternalServerError, u.logger)

		return
	}

	username := request.FormValue("username")
	email := request.FormValue("email")
	password := request.FormValue("password")

	if username == "" {
		u.logger.Error("missing required field", slog.String("name", "username"))

		errResp := dto.NewErrorDocument(
			[]dto.ErrorObject{validationError("username")},
		)
		writeJSONResponse(
			writer, errResp, http.StatusUnprocessableEntity, u.logger,
		)

		return
	}

	if email == "" {
		u.logger.Error("missing required field", slog.String("name", "email"))

		errResp := dto.NewErrorDocument(
			[]dto.ErrorObject{validationError("email")},
		)
		writeJSONResponse(
			writer, errResp, http.StatusUnprocessableEntity, u.logger,
		)

		return
	}

	if password == "" {
		u.logger.Error("missing required field", slog.String("name", "password"))

		errResp := dto.NewErrorDocument(
			[]dto.ErrorObject{validationError("password")},
		)
		writeJSONResponse(
			writer, errResp, http.StatusUnprocessableEntity, u.logger,
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

		errResp := dto.NewErrorDocument([]dto.ErrorObject{
			{ //nolint:exhaustruct
				Status: errInternalServer,
				Code:   strconv.Itoa(http.StatusInternalServerError),
				Title:  "Failed to Create User",
				Detail: "An unexpected error occurred while creating the user.",
			},
		})
		writeJSONResponse(
			writer, errResp, http.StatusInternalServerError, u.logger,
		)

		return
	}

	writer.Header().Set("Content-Type", "application/vnd.api+json")
	writer.WriteHeader(http.StatusCreated)

	encodingErr := json.NewEncoder(writer).Encode(resp)
	if encodingErr != nil {
		u.logger.Error(
			"JSON encoding failed",
			slog.String("error", encodingErr.Error()),
		)

		errResp := dto.NewErrorDocument([]dto.ErrorObject{
			{ //nolint:exhaustruct
				Status: errInternalServer,
				Code:   strconv.Itoa(http.StatusInternalServerError),
				Title:  "JSON Encoding Failed",
				Detail: "An unexpected error occurred while encoding the response.",
			},
		})
		writeJSONResponse(
			writer, errResp, http.StatusInternalServerError, u.logger,
		)

		return
	}

	u.logger.Info(
		"successfully created new user",
		slog.String("id", resp.ID),
		slog.String("username", resp.Username),
	)
}
