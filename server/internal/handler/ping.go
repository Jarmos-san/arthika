package handler

import (
	"context"

	"github.com/Jarmos-san/arthika/server/internal/api"
)

// Ping responds with a simple health check status.
func (h *Handler) Ping(
	ctx context.Context,
	_ api.PingRequestObject,
) (api.PingResponseObject, error) {
	h.logger.DebugContext(ctx, "ping received")

	return api.Ping200JSONResponse{Status: "ok"}, nil
}
