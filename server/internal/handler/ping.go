package handler

import (
	"context"
	"log/slog"

	"github.com/Jarmos-san/arthika/server/internal/api"
)

var _ api.StrictServerInterface = (*PingHandler)(nil)

// PingHandler handles health check requests.
type PingHandler struct {
	logger *slog.Logger
}

// NewPingHandler creates a new PingHandler.
func NewPingHandler(logger *slog.Logger) *PingHandler {
	return &PingHandler{logger: logger}
}

// Ping responds with a status of "ok" to confirm the server is running.
func (h *PingHandler) Ping(
	ctx context.Context,
	_ api.PingRequestObject,
) (api.PingResponseObject, error) {
	h.logger.DebugContext(ctx, "ping received")

	return api.Ping200JSONResponse{Status: "ok"}, nil
}
