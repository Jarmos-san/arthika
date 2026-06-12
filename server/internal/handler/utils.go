package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

const (
	errKeyStatus = "status"
	errKeyTitle  = "title"
	errKeyDetail = "detail"
)

func validationError(field string) map[string]any {
	return map[string]any{
		errKeyStatus: http.StatusUnprocessableEntity,
		errKeyTitle:  "Missing Required Field",
		errKeyDetail: fmt.Sprintf("Missing required field: '%s' cannot be empty", field),
	}
}

func writeJSONError(
	writer http.ResponseWriter,
	errs []map[string]any,
	statusCode int,
	logger *slog.Logger,
) {
	writeJSONResponse(writer, map[string]any{"errors": errs}, statusCode, logger)
}

func writeJSONResponse(
	writer http.ResponseWriter,
	payload any,
	statusCode int,
	logger *slog.Logger,
) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)

	err := json.NewEncoder(writer).Encode(payload)
	if err != nil {
		logger.Error("failed to encode JSON", slog.Any("error", err))
	}
}
