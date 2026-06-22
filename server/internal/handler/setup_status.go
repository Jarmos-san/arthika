package handler

import (
	"context"
	"fmt"

	"github.com/Jarmos-san/arthika/server/internal/api"
)

// SystemStatus returns whether the application needs initial setup.
func (h *Handler) SystemStatus(
	ctx context.Context,
	_ api.SystemStatusRequestObject,
) (api.SystemStatusResponseObject, error) {
	count, err := h.querier.CountUsers(ctx)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to count users", "error", err)

		return nil, fmt.Errorf("count users: %w", err)
	}

	return api.SystemStatus200JSONResponse{
		NeedsSetup: count == 0,
	}, nil
}
