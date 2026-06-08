package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Jarmos-san/arthika/server/internal/config"
	"github.com/Jarmos-san/arthika/server/internal/dto"
	jwt "github.com/golang-jwt/jwt/v5"
)

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

// writeJSONResponse serialises the provided payload into a JSON document and writes it
// to the HTTP response writer using the supplied status code.
//
// The helper centralises JSON response handling for HTTP handlers by:
//
// 1. Encoding the provided payload into JSON.
// 2. Setting the appropriate JSON:API content type header.
// 3. Writing the HTTP status code.
// 4. Writing the serialised JSON response body.
// 5. Logging any encoding or write failures.
//
// The response uses the JSON:API media type:
//
//	Content-Type: application/vnd.api+json
//
// If JSON encoding fails, the function logs the error and attempts to return an HTTP
// 500 Internal Server Error response.
//
// Parameters:
//
//	writer - The HTTP response writer used to construct the response.
//	payload - The response payload to serialise into JSON.
//
// statusCode - The HTTP status code to send with the response.
// logger - The structured logger used for reporting response errors.
func writeJSONResponse(
	writer http.ResponseWriter,
	payload any,
	statusCode int,
	logger *slog.Logger,
) {
	writer.Header().Set("Content-Type", "application/vnd.api+json")
	writer.WriteHeader(statusCode)

	err := json.NewEncoder(writer).Encode(payload)
	if err != nil {
		logger.Error("failed to encode JSON", slog.Any("error", err))
	}
}
