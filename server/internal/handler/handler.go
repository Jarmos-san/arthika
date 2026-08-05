package handler

import (
	"log/slog"

	"github.com/Jarmos-san/arthika/server/internal/api"
	"github.com/Jarmos-san/arthika/server/internal/repository"
)

var _ api.StrictServerInterface = (*Handler)(nil)

// Handler serves HTTP requests for the API.
type Handler struct {
	logger  *slog.Logger
	querier repository.Querier
}

// NewHandler creates a new Handler.
func NewHandler(logger *slog.Logger, querier repository.Querier) *Handler {
	return &Handler{
		logger:  logger,
		querier: querier,
	}
}
