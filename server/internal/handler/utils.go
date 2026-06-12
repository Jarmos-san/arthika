package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

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
