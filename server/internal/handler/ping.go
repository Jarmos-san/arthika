package handler

import (
	"context"

	"github.com/Jarmos-san/arthika/server/internal/api"
)

func (h *Handler) Ping(
	ctx context.Context,
	_ api.PingRequestObject,
) (api.PingResponseObject, error) {
	h.logger.DebugContext(ctx, "ping received")

	return api.Ping200JSONResponse{Status: "ok"}, nil
}
