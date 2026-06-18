package handler

import (
	"log/slog"

	"github.com/Jarmos-san/arthika/server/internal/api"
)

var _ api.StrictServerInterface = (*Handler)(nil)

type Handler struct {
	logger *slog.Logger
}

func NewHandler(logger *slog.Logger) *Handler {
	return &Handler{logger: logger}
}
