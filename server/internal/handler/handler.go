package handler

import (
	"log/slog"

	"github.com/Jarmos-san/arthika/server/internal/api"
)

var _ api.StrictServerInterface = (*Handler)(nil)

// Handler serves HTTP requests for the API.
type Handler struct {
	logger *slog.Logger
}

// NewHandler creates a new Handler.
func NewHandler(logger *slog.Logger) *Handler {
	return &Handler{logger: logger}
}
