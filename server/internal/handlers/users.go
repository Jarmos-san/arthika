// Package handlers provides HTTP transport layer implementations.
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
package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/Jarmos-san/arthika/server/internal/config"
	"github.com/Jarmos-san/arthika/server/internal/dto"
	"github.com/Jarmos-san/arthika/server/internal/services"
)

// UserHandler handles HTTP requests related to user resources.
//
// UserHandler acts as the transport layer for user-related endpoints. It delegates
// business logic to the injected UserService and is responsible for constructing HTTP
// responses based on the results.
//
// The handler is safe for concurrent use provided its dependencies are also
// concurrency-safe.
type UserHandler struct {
	service services.UserService
	logger  *slog.Logger
}

// NewUserHandler constructs a new UserHandler with its required dependencies.
//
// Parameters:
//   - service: provides user-related business logic
//   - logger:  used for structured logging within the handler
//
// The returned handler is ready to be registered with an HTTP router.
func NewUserHandler(service services.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

// UserResponse represents the JSON response payload for a user resource.
//
// This type is a transport-layer DTO and defines the external API contract. It should
// remain decoupled from internal domain models to avoid leaking implementation details.
type UserResponse struct {
	Name string `json:"name"`
}

// GetUser handles HTTP GET requests for retrieving a user.
//
// It invokes the underlying UserService to fetch user data and returns a JSON response
// to the client.
//
// Success Response:
//   - Status: 200 OK
//   - Content-Type: application/vnd.api+json
//   - Body: JSON-encoded UserResponse
//
// Error Response:
//   - Status: 500 Internal Server Error
//   - Body: plain text error message
//
// Notes:
//   - The method uses a value receiver; this is acceptable since the handler
//     struct contains only references. Pointer receivers are still preferred
//     for consistency across handler methods.
//   - JSON encoding errors are logged but cannot alter the response once headers
//     have been written.
//   - A new logger instance is created during encoding failure, which is
//     inefficient and should be avoided in favor of the injected logger.
func (u UserHandler) GetUser(writer http.ResponseWriter, _ *http.Request) {
	user, serviceErr := u.service.GetUser()
	if serviceErr != nil {
		u.logger.Error(
			"failed to fetch user",
			slog.String("error", serviceErr.Error()),
		)
		http.Error(writer, "internal server error", http.StatusInternalServerError)

		return
	}

	resp := UserResponse{
		Name: user.Name,
	}

	writer.Header().Set("Content-Type", "application/vnd.api+json")
	writer.WriteHeader(http.StatusOK)

	encodingErr := json.NewEncoder(writer).Encode(resp)
	if encodingErr != nil {
		u.logger.Error(
			"JSON encoding failed",
			slog.String("error", encodingErr.Error()),
		)
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

		http.Error(writer, "form parsing failed", http.StatusInternalServerError)

		return
	}

	username := request.FormValue("username")
	email := request.FormValue("email")
	password := request.FormValue("password")

	if username == "" {
		u.logger.Error("missing required field", slog.String("name", "username"))

		http.Error(
			writer,
			"missing required field: 'username'",
			http.StatusUnprocessableEntity,
		)

		return
	}

	if email == "" {
		u.logger.Error("missing required field", slog.String("name", "email"))

		http.Error(
			writer,
			"missing required field: 'email'",
			http.StatusUnprocessableEntity,
		)

		return
	}

	if password == "" {
		u.logger.Error("missing required field", slog.String("name", "password"))

		http.Error(
			writer,
			"missing required field: 'password'",
			http.StatusUnprocessableEntity,
		)

		return
	}

	resp, serviceErr := u.service.CreateUser(username, password)
	if serviceErr != nil {
		u.logger.Error("failed to create user")

		http.Error(writer, "internal server error", http.StatusInternalServerError)

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

		http.Error(writer, "internal server error", http.StatusInternalServerError)

		return
	}

	u.logger.Info(
		"successfully created new user",
		slog.String("id", resp.ID),
		slog.String("username", resp.Name),
		slog.String("password", resp.PasswordHash),
	)
}

// validationError creates a JSON:API compliant error object for a missing required
// field validation failure.
//
// The returned dto.ErrorObject contains metadata describing the validation error
// including the HTTP status code, a short title and a detailed human-readable message
// indicating which field is missing.
//
// Parameters:
//
//	field - The name of the field that failed the validation check.
//
// Returns:
//
//	A dto.ErrorObject representing the validation failure.
//
// TODO: Move this function elsewhere perhaps into a different module
func validationError(field string) dto.ErrorObject {
	return dto.ErrorObject{
		Status: "Validation Error",
		Code:   strconv.Itoa(http.StatusUnprocessableEntity),
		Title:  "Missing Required Field",
		Detail: fmt.Sprintf("Missing required field: '%s' cannot be empty", field),
	}
}

// createJWT generates and signs a new JSON Web Token (JWT) using the application's
// configured secret key.
//
// The JWT is signed using the HS512 signing algorithm. The secret key used for signing
// is loaded frfom the application configuration.
//
// Returns:
//
//	string - The signed JWT.
//	error - An error if token signing fails.
func createJWT() (string, error) {
	// Load the app config, required for fetching the secret key for the JWT
	config := config.LoadConfig()

	// Setup the JWT generation process
	key := []byte(config.TokenSecret)    // Use the secret token loaded from the configs
	t := jwt.New(jwt.SigningMethodHS512) // Use HMAC-SHA256 algorithm for token signing

	return t.SignedString(key) // Create a signed JWT (or throw an error)
}

// LoginUser authenticates a user login request and returns a JSON:API compliant
// response containing a newly generated JWT access token.
//
// The handler performs the following operations:
//
// 1. Reads the username and password from the incoming request.
// 2. Validates that all required fields are present.
// 3. Returns validation errors if required fields are missing.
// 4. Generates a signed JWT access token.
// 5. Wraps the token in a JSON:API resource object.
// 6. Serialises and returns the response to the client.
//
// Validation failures return HTTP 422 Unprocessable Entity responses. Successful
// authentication returns HTTP 200 OK.
//
// Parameters:
//
//	writer - The HTTP response writer used to construct the response.
//	request - The incoming HTTP request containing login credentials.
func (u UserHandler) LoginUser(writer http.ResponseWriter, request *http.Request) {
	// Read the form-data from the request
	username := request.FormValue("username") // Read the username
	password := request.FormValue("password") // Read the password

	// Initialise a slice to store the validation errors which will be serialised as a
	// single object
	validationErrs := make([]dto.ErrorObject, 0)

	// Validate the username, or return an error response if invalid
	if username == "" {
		u.logger.Error("missing required field", slog.String("field", "username"))
		validationErrs = append(validationErrs, validationError("username"))
	}

	// Validate the password, or return an error response if invalid
	if password == "" {
		u.logger.Error("missing required field", slog.String("field", "password"))
		validationErrs = append(validationErrs, validationError("password"))
	}

	// Return a error response if the user credentials are invalid
	if len(validationErrs) > 0 {
		writer.Header().Set("Content-Type", "application/vnd.api+json")
		writer.WriteHeader(http.StatusUnprocessableEntity)
		resp := dto.NewErrorDocument(validationErrs)

		if err := json.NewEncoder(writer).Encode(resp); err != nil {
			u.logger.Error(
				"server error",
				slog.String("type", "failed to encode JSON data"),
			)

			http.Error(
				writer,
				"failed to encode to JSON",
				http.StatusInternalServerError,
			)
		}

		return
	}

	// Generate the JWT, or return an error response if it failed
	token, tokenErr := createJWT()
	if tokenErr != nil {
		writer.Header().Set("Content-Type", "application/vnd.api+json")
		writer.WriteHeader(http.StatusUnprocessableEntity)
		resp := dto.NewErrorDocument(validationErrs)
		if err := json.NewEncoder(writer).Encode(resp); err != nil {
			u.logger.Error(
				"server error",
				slog.String("type", "failed to encode JSON data"),
			)

			http.Error(
				writer,
				"failed to encode to JSON",
				http.StatusInternalServerError,
			)
		}

		return
	}

	// Create a "Resource Object" (according to the JSON:API spec) to respond to the
	// client request with the JWT
	resourceObject := dto.ResourceObject{
		Type: "tokens",
		ID:   uuid.NewString(),
		Attributes: map[string]any{
			"accessToken": token,
			"tokenType":   "Bearer",
			"expiresIn":   time.Hour.Seconds(), // 3600 seconds
		},
	}

	// Create a JSON:API compliant "document" to be serialised into a JSON response
	resp := dto.NewSingleDocument(resourceObject)

	// Set the HTTP headers for the response
	writer.Header().Set("Content-Type", "appplication/vnd.api+json")
	writer.WriteHeader(http.StatusOK)

	// Create an appropriate log statement before serializing the response
	u.logger.Info("new tokens generated", slog.String("username", username))

	// Serialise the objects and write a JSON response for the client
	encodingErr := json.NewEncoder(writer).Encode(resp)
	if encodingErr != nil {
		http.Error(writer, "failed to encode to JSON", http.StatusInternalServerError)
	}
}
